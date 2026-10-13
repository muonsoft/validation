"""Restricted model gateway and Unix-socket transport for offline eval workers."""

from contextlib import contextmanager
import copy
import http.client
import http.server
import json
import os
import re
import select
import socket
import socketserver
import subprocess
import sys
import tempfile
import threading
import urllib.parse


class TCPServer(socketserver.ThreadingTCPServer):
    allow_reuse_address = True
    daemon_threads = True


class UnixServer(socketserver.ThreadingUnixStreamServer):
    daemon_threads = True


@contextmanager
def serving(server):
    thread = threading.Thread(target=server.serve_forever, daemon=True)
    thread.start()
    try:
        yield server
    finally:
        server.shutdown()
        server.server_close()
        thread.join()


class Bridge(socketserver.BaseRequestHandler):
    def handle(self):
        family, address = self.server.destination
        with socket.socket(family, socket.SOCK_STREAM) as peer:
            peer.settimeout(10)
            try:
                peer.connect(address)
                peer.settimeout(None)
                readers = [self.request, peer]
                while readers:
                    ready, _, _ = select.select(readers, [], [], 60)
                    for source in ready:
                        target = peer if source is self.request else self.request
                        data = source.recv(65536)
                        if data:
                            target.sendall(data)
                        else:
                            readers.remove(source)
                            target.shutdown(socket.SHUT_WR)
            except OSError:
                return


def resolve(value, env):
    if isinstance(value, str):
        match = re.fullmatch(r"\{env:([A-Z_][A-Z_0-9]*)\}", value)
        if match:
            if not env.get(match[1]):
                raise ValueError(f"missing provider credential variable: {match[1]}")
            return env[match[1]]
    if isinstance(value, dict):
        return {key: resolve(item, env) for key, item in value.items()}
    if isinstance(value, list):
        return [resolve(item, env) for item in value]
    return value


def provider_setup(config, model, env):
    """Fail closed for adapters whose request protocol has not been restricted."""
    provider_id, model_id = model.split("/", 1)
    provider = config.get("providers", {}).get(provider_id, {})
    definition = provider.get("models", {}).get(model_id)
    if not definition:
        raise ValueError("network isolation requires an explicit selected model definition")
    package = definition.get("package", provider.get("package"))
    if package != "@opencode/ai/providers/openai-compatible":
        raise ValueError("network isolation currently supports only the V2 openai-compatible adapter")
    settings = resolve({**provider.get("settings", {}), **definition.get("settings", {})}, env)
    endpoint = urllib.parse.urlsplit(settings.get("baseURL", ""))
    if (endpoint.scheme != "https" or not endpoint.hostname or endpoint.username or endpoint.password
            or endpoint.query or endpoint.fragment or endpoint.port not in (None, 443)):
        raise ValueError("model gateway requires a fixed HTTPS baseURL on port 443")
    if not settings.get("apiKey"):
        raise ValueError("model gateway requires an apiKey credential")
    headers = resolve({**provider.get("headers", {}), **definition.get("headers", {})}, env)
    if any(key.lower() in {"host", "content-length", "transfer-encoding", "connection"} for key in headers):
        raise ValueError("provider headers cannot override HTTP framing or destination")
    headers = {**headers, "Authorization": "Bearer " + settings["apiKey"], "Content-Type": "application/json"}
    upstream = {"host": endpoint.hostname, "path": endpoint.path.rstrip("/") + "/chat/completions",
                "headers": headers, "model": definition.get("modelID", model_id)}
    # Construct a new provider instead of retaining credentials or alternate URLs.
    safe_definition = {key: copy.deepcopy(definition[key]) for key in
                       ("modelID", "name", "capabilities", "limit", "compatibility", "variants", "cost")
                       if key in definition}
    safe_settings = {"baseURL": "http://127.0.0.1:18080/v1", "apiKey": "eval-local-gateway"}
    if "provider" in settings:
        safe_settings["provider"] = settings["provider"]
    safe_definition.update(package=package, settings=safe_settings, headers={})
    safe = {key: copy.deepcopy(value) for key, value in config.items() if key != "providers"}
    safe["providers"] = {provider_id: {"package": package, "settings": safe_settings,
                                        "models": {model_id: safe_definition}}}
    if re.search(r"\{env:", json.dumps(safe)):
        raise ValueError("unexpected environment reference outside provider credentials")
    return safe, upstream


class Gateway(http.server.BaseHTTPRequestHandler):
    """An inference endpoint, never a general-purpose forward proxy."""
    protocol_version = "HTTP/1.0"

    def log_message(self, *args):
        pass  # Never log credentials, prompts, completions, or upstream errors.

    def do_GET(self):
        if self.path != "/health":
            self.send_error(403)
            return
        self.send_response(204)
        self.end_headers()

    def do_CONNECT(self):
        self.send_error(403)

    def do_POST(self):
        if self.path != "/v1/chat/completions" or self.headers.get("Transfer-Encoding"):
            self.send_error(403)
            return
        try:
            size = int(self.headers.get("Content-Length", "0"))
            if not 0 < size <= 16 * 1024 * 1024:
                raise ValueError()
            self.connection.settimeout(60)
            raw = self.rfile.read(size)
            body = json.loads(raw)
            if (not isinstance(body, dict) or body.get("model") != self.server.upstream["model"]
                    or not isinstance(body.get("messages"), list)):
                raise ValueError()
            allowed = {"model", "messages", "stream", "stream_options", "max_tokens", "max_completion_tokens",
                       "temperature", "top_p", "presence_penalty", "frequency_penalty", "seed", "stop",
                       "tools", "tool_choice", "parallel_tool_calls", "response_format", "reasoning_effort",
                       "thinking", "reasoning", "verbosity", "metadata", "user", "service_tier",
                       "logprobs", "top_logprobs", "store"}
            if set(body) - allowed:
                raise ValueError()
            # Remote search/file/image tools would reintroduce a fetching channel.
            if any(tool.get("type") != "function" for tool in body.get("tools", [])):
                raise ValueError()
            for message in body["messages"]:
                content = message.get("content")
                if content is not None and not isinstance(content, (str, list)):
                    raise ValueError()
                if isinstance(content, list) and any(part.get("type") != "text" for part in content):
                    raise ValueError()
        except (ValueError, TypeError, AttributeError, OSError):
            self.send_error(400)
            return
        upstream = self.server.upstream
        connection = http.client.HTTPSConnection(upstream["host"], timeout=900)
        started = False
        try:
            connection.request("POST", upstream["path"], body=raw, headers=upstream["headers"])
            response = connection.getresponse()
            # Never follow redirects, including redirects to the same host.
            if 300 <= response.status < 400:
                self.send_error(502, "upstream redirect refused")
                return
            self.send_response(response.status)
            self.send_header("Content-Type", response.getheader("Content-Type", "application/json"))
            self.end_headers()
            started = True
            while chunk := response.read1(65536):
                self.wfile.write(chunk)
                self.wfile.flush()
        except (OSError, http.client.HTTPException):
            if not started:
                self.send_error(502, "model gateway request failed")
        finally:
            connection.close()


@contextmanager
def network_runtime(config, model, env):
    safe, upstream = provider_setup(config, model, env)
    # Short paths are required by AF_UNIX; this directory holds sockets only.
    with tempfile.TemporaryDirectory(prefix="eval-net-") as directory:
        gateway = UnixServer(directory + "/model.sock", Gateway)
        gateway.upstream = upstream
        control = TCPServer(("127.0.0.1", 0), Bridge)
        control.destination = (socket.AF_UNIX, directory + "/control.sock")
        with serving(gateway), serving(control):
            yield safe, directory, control.server_address[1]


def worker_main(directory, executable, args):
    """Runs inside the network namespace; only loopback and Unix bridges exist."""
    interfaces = {name for _, name in socket.if_nameindex()}
    if interfaces != {"lo"}:
        raise RuntimeError("worker network namespace is not isolated")
    for address in ("1.1.1.1", "2606:4700:4700::1111"):
        family = socket.AF_INET6 if ":" in address else socket.AF_INET
        with socket.socket(family, socket.SOCK_STREAM) as probe:
            probe.settimeout(1)
            try:
                probe.connect((address, 443))
            except OSError:
                pass
            else:
                raise RuntimeError("worker has external network connectivity")
    model = TCPServer(("127.0.0.1", 18080), Bridge)
    model.destination = (socket.AF_UNIX, directory + "/model.sock")
    control = UnixServer(directory + "/control.sock", Bridge)
    control.destination = (socket.AF_INET, ("127.0.0.1", 18081))
    with serving(model), serving(control):
        connection = http.client.HTTPConnection("127.0.0.1", 18080, timeout=5)
        connection.request("GET", "/health")
        if connection.getresponse().status != 204:
            raise RuntimeError("model gateway unavailable")
        connection.close()
        for method, route, body, status in (
                ("CONNECT", "github.com:443", None, 403),
                ("GET", "https://github.com/", None, 403),
                ("POST", "/v1/chat/completions", "{}", 400)):
            connection = http.client.HTTPConnection("127.0.0.1", 18080, timeout=5)
            connection.request(method, route, body=body)
            response = connection.getresponse()
            if response.status != status:
                raise RuntimeError("model gateway restriction check failed")
            response.read()
            connection.close()
        result = subprocess.run([executable, *args], env=os.environ)
        return result.returncode


if __name__ == "__main__":
    sys.exit(worker_main(sys.argv[1], sys.argv[2], sys.argv[3:]))
