#!/usr/bin/env python3
"""Gateway tests use local sockets only; optional bubblewrap probe makes no inference."""
import contextlib
import http.client
import importlib.util
import io
import json
import os
from pathlib import Path
import shutil
import socket
import subprocess
import tempfile
import threading
import time
import unittest
from unittest.mock import Mock, patch

import network


def fixture():
    return {'providers': {'fixture': {
        'package': '@opencode/ai/providers/openai-compatible',
        'settings': {'baseURL': 'https://provider.example/v1', 'apiKey': '{env:MODEL_KEY}'},
        'models': {'model': {'modelID': 'actual-model', 'headers': {'X-Tenant': '{env:TENANT}'}}}}}}


class UnixHTTP(http.client.HTTPConnection):
    def __init__(self, path):
        super().__init__('localhost', timeout=5)
        self.path = path

    def connect(self):
        self.sock = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
        self.sock.settimeout(self.timeout)
        self.sock.connect(self.path)


class NetworkTests(unittest.TestCase):
    def setUp(self):
        self.env = {'MODEL_KEY': 'host-only-secret', 'TENANT': 'host-only-tenant'}
        self.temp = tempfile.TemporaryDirectory(prefix='eval-network-test-')
        self.addCleanup(self.temp.cleanup)
        self.root = Path(self.temp.name)

    def test_fixed_provider_and_credentials(self):
        safe, upstream = network.provider_setup(fixture(), 'fixture/model', self.env)
        self.assertEqual(upstream['host'], 'provider.example')
        self.assertEqual(upstream['path'], '/v1/chat/completions')
        self.assertEqual(upstream['model'], 'actual-model')
        self.assertEqual(upstream['headers']['Authorization'], 'Bearer host-only-secret')
        self.assertNotIn('host-only', json.dumps(safe))
        self.assertNotIn('{env:', json.dumps(safe))
        self.assertNotIn('https://provider.example', json.dumps(safe))
        self.assertIn('http://127.0.0.1:18080/v1', json.dumps(safe))

    def test_unknown_adapter_or_unsafe_endpoint_fails_closed(self):
        for url in ['http://provider.example', 'https://provider.example/?next=evil',
                    'https://user:password@provider.example', 'https://provider.example:8443']:
            config = fixture()
            config['providers']['fixture']['settings']['baseURL'] = url
            with self.assertRaises(ValueError):
                network.provider_setup(config, 'fixture/model', self.env)
        config = fixture()
        config['providers']['fixture']['package'] = 'unknown'
        with self.assertRaises(ValueError):
            network.provider_setup(config, 'fixture/model', self.env)

    @contextlib.contextmanager
    def gateway(self):
        _, upstream = network.provider_setup(fixture(), 'fixture/model', self.env)
        server = network.UnixServer(str(self.root / 'model.sock'), network.Gateway)
        server.upstream = upstream
        with network.serving(server):
            yield str(self.root / 'model.sock')

    def request(self, path, method, route, body=None, headers=None):
        with contextlib.closing(UnixHTTP(path)) as client:
            client.request(method, route, body=body, headers=headers or {})
            response = client.getresponse()
            return response.status, response.read()

    def test_gateway_blocks_fetching_and_bad_payloads_without_upstream(self):
        valid = {'model': 'actual-model', 'messages': [{'role': 'user', 'content': 'hello'}]}
        with self.gateway() as path, patch.object(network.http.client, 'HTTPSConnection') as upstream:
            self.assertEqual(self.request(path, 'GET', '/health')[0], 204)
            for method, route in [('CONNECT', 'github.com:443'), ('GET', 'https://github.com/'),
                                  ('POST', '/v1/chat/completions?url=https://github.com'),
                                  ('POST', '//github.com'), ('POST', '/v1/../models')]:
                self.assertEqual(self.request(path, method, route)[0], 403)
            for body in [dict(valid, model='another-model'), dict(valid, web_search_options={}),
                         dict(valid, tools=[{'type': 'web_search'}]),
                         dict(valid, messages=[{'content': [{'type': 'image_url', 'image_url': {'url': 'https://github.com'}}]}]),
                         [], {'model': 'actual-model'}, dict(valid, messages=[{'content': {'url': 'https://github.com'}}])]:
                self.assertEqual(self.request(path, 'POST', '/v1/chat/completions', json.dumps(body))[0], 400)
            upstream.assert_not_called()

    def test_streamed_completion_uses_host_credentials_and_fixed_destination(self):
        response = Mock(status=200)
        response.getheader.return_value = 'text/event-stream'
        response.read1.side_effect = [b'data: {"text":"ok"}\n\n', b'data: [DONE]\n\n', b'']
        connection = Mock()
        connection.getresponse.return_value = response
        body = json.dumps({'model': 'actual-model', 'messages': [{'role': 'user', 'content': 'hello'}], 'stream': True})
        with self.gateway() as path, patch.object(network.http.client, 'HTTPSConnection', return_value=connection) as upstream:
            status, data = self.request(path, 'POST', '/v1/chat/completions', body,
                                        {'Authorization': 'attacker', 'Host': 'github.com', 'X-Tenant': 'attacker',
                                         'x-opencode-session': 'ses_test_conversation', 'User-Agent': 'opencode/test',
                                         'X-Forwarded-Host': 'github.com', 'Proxy-Authorization': 'attacker'})
        self.assertEqual(status, 200)
        self.assertIn(b'[DONE]', data)
        upstream.assert_called_once_with('provider.example', timeout=900)
        args, kw = connection.request.call_args
        self.assertEqual(args, ('POST', '/v1/chat/completions'))
        self.assertEqual(kw['headers']['Authorization'], 'Bearer host-only-secret')
        self.assertEqual(kw['headers']['X-Tenant'], 'host-only-tenant')
        self.assertEqual(kw['headers']['x-opencode-session'], 'ses_test_conversation')
        self.assertEqual(kw['headers']['user-agent'], 'opencode/test')
        self.assertNotIn('host', {key.lower() for key in kw['headers']})
        self.assertNotIn('x-forwarded-host', {key.lower() for key in kw['headers']})
        self.assertNotIn('proxy-authorization', {key.lower() for key in kw['headers']})
        connection.close.assert_called_once()

    def test_client_disconnect_before_headers_and_during_stream(self):
        for stage in ('headers', 'stream', 'redirect', 'upstream-error'):
            with self.subTest(stage=stage), tempfile.TemporaryDirectory(dir=self.root) as directory:
                _, upstream = network.provider_setup(fixture(), 'fixture/model', self.env)
                waiting, release, closed = threading.Event(), threading.Event(), threading.Event()
                response = Mock(status=307 if stage == 'redirect' else 200)
                response.getheader.return_value = 'text/event-stream'

                def wait_for_disconnect():
                    waiting.set()
                    if not release.wait(5):
                        raise RuntimeError('disconnect test timed out')

                def getresponse():
                    if stage != 'stream':
                        wait_for_disconnect()
                    if stage == 'upstream-error':
                        raise TimeoutError('provider timed out')
                    return response

                def read_chunk(size):
                    wait_for_disconnect()
                    return b'data: chunk\n\n'

                response.read1.side_effect = read_chunk
                connection = Mock()
                connection.getresponse.side_effect = getresponse
                connection.close.side_effect = closed.set
                server = network.UnixServer(directory + '/gateway.sock', network.Gateway)
                server.upstream = upstream
                server.handle_error = Mock()
                body = json.dumps({'model': 'actual-model', 'messages': []}).encode()
                request = (b'POST /v1/chat/completions HTTP/1.0\r\nContent-Length: '
                           + str(len(body)).encode() + b'\r\n\r\n' + body)
                errors = []
                send_error = network.Gateway.send_error

                def record_error(handler, *args, **kwargs):
                    errors.append(args[0])
                    return send_error(handler, *args, **kwargs)

                with network.serving(server), patch.object(network.http.client, 'HTTPSConnection', return_value=connection), patch.object(network.Gateway, 'send_error', record_error):
                    with socket.socket(socket.AF_UNIX, socket.SOCK_STREAM) as client:
                        client.connect(directory + '/gateway.sock')
                        client.sendall(request)
                        self.assertTrue(waiting.wait(5))
                        client.shutdown(socket.SHUT_RDWR)
                    release.set()
                    self.assertTrue(closed.wait(5))
                server.handle_error.assert_not_called()
                connection.close.assert_called_once()
                self.assertEqual(errors, [502] if stage in ('redirect', 'upstream-error') else [])

    def test_upstream_failure_returns_502_to_connected_client(self):
        connection = Mock()
        connection.getresponse.side_effect = TimeoutError('upstream failed')
        with self.gateway() as path, patch.object(network.http.client, 'HTTPSConnection', return_value=connection):
            status, _ = self.request(path, 'POST', '/v1/chat/completions',
                                     json.dumps({'model': 'actual-model', 'messages': []}))
        self.assertEqual(status, 502)
        connection.close.assert_called_once()

    @unittest.skipUnless(os.environ.get('EVAL_TEST_OPENCODE') == '1', 'set EVAL_TEST_OPENCODE=1 for real CLI with simulated provider')
    def test_real_opencode_sends_session_header_through_gateway(self):
        spec = importlib.util.spec_from_file_location('runner_cli_probe', Path(__file__).with_name('run.py'))
        runner = importlib.util.module_from_spec(spec)
        spec.loader.exec_module(runner)
        workspace = self.root / 'workspace'
        workspace.mkdir()
        subprocess.run(['git', 'init', '-q', str(workspace)], check=True)
        config = fixture()
        provider = config['providers'].pop('fixture')
        definition = provider['models'].pop('model')
        definition.update(modelID='deepseek-v4.1-flash', name='DeepSeek test transport',
                          limit={'context': 128000, 'output': 4096})
        provider['models']['deepseek-v4.1-flash'] = definition
        config['providers']['opencode-go'] = provider
        path = self.root / 'provider.json'
        runner.write_json(path, runner.configure_provider(config, 2))
        env = runner.isolated_env(self.root, path)
        env.update(self.env)
        env.update(json.loads(subprocess.check_output(['go', 'env', '-json', 'GOROOT', 'GOMODCACHE'])))
        env.update(EVAL_LIBRARY_ROOT=str(self.root / 'libraries'), EVAL_OPENCODE_MODEL='opencode-go/deepseek-v4.1-flash')
        requests = []

        def connect(host, timeout):
            connection = Mock()
            response = Mock(status=200)
            response.getheader.return_value = 'text/event-stream'
            chunk = {'id': 'chatcmpl-probe', 'object': 'chat.completion.chunk', 'created': 1,
                     'model': 'deepseek-v4.1-flash',
                     'choices': [{'index': 0, 'delta': {'role': 'assistant', 'content': 'NETWORK_OK'}, 'finish_reason': None}]}
            finish = {**chunk, 'choices': [{'index': 0, 'delta': {}, 'finish_reason': 'stop'}],
                      'usage': {'prompt_tokens': 1, 'completion_tokens': 1, 'total_tokens': 2}}
            response.read1.side_effect = [f'data: {json.dumps(chunk)}\n\n'.encode(),
                                         f'data: {json.dumps(finish)}\n\n'.encode(), b'data: [DONE]\n\n', b'']
            connection.getresponse.return_value = response
            connection.request.side_effect = lambda *args, **kw: requests.append(kw)
            return connection

        # Only the host gateway's HTTPS connection is replaced. The installed CLI,
        # sandbox, HTTP request serialization and both Unix bridges are real.
        with patch.object(network.http.client, 'HTTPSConnection', side_effect=connect):
            result = runner.opencode_command(shutil.which('opencode'),
                ['run', '--agent', 'build', '--model', env['EVAL_OPENCODE_MODEL'], '--format', 'json',
                 'Reply exactly NETWORK_OK. Do not use tools.'],
                workspace, self.root / 'events.jsonl', env, 2, timeout=45)
        self.assertEqual(result['status'], 'completed')
        self.assertTrue(result.get('completion', {}).get('verified'))
        self.assertTrue(requests)
        for request in requests:
            headers = {key.lower(): value for key, value in request['headers'].items()}
            self.assertTrue(headers.get('x-opencode-session'))
            self.assertTrue(headers.get('user-agent'))
            self.assertEqual(headers['authorization'], 'Bearer host-only-secret')
        self.assertIn('NETWORK_OK', (self.root / 'events.jsonl').read_text())

    def test_redirect_is_not_followed(self):
        response = Mock(status=307)
        connection = Mock()
        connection.getresponse.return_value = response
        with self.gateway() as path, patch.object(network.http.client, 'HTTPSConnection', return_value=connection) as upstream:
            status, _ = self.request(path, 'POST', '/v1/chat/completions', json.dumps({'model': 'actual-model', 'messages': []}))
        self.assertEqual(status, 502)
        self.assertEqual(upstream.call_count, 1)
        response.read1.assert_not_called()

    def test_server_environment_has_no_provider_credentials_or_proxy(self):
        spec = importlib.util.spec_from_file_location('runner_network_test', Path(__file__).with_name('run.py'))
        runner = importlib.util.module_from_spec(spec)
        spec.loader.exec_module(runner)
        config = self.root / 'provider.json'
        config.write_text(json.dumps(fixture()))
        env = {**self.env, 'OPENCODE_CONFIG': str(config), 'EVAL_OPENCODE_MODEL': 'fixture/model',
               'HTTPS_PROXY': 'http://host-proxy', 'EVAL_OPENCODE_TOKEN': 'another-secret'}
        @contextlib.contextmanager
        def server(binary, cwd, log, worker_env, safe, port):
            self.assertNotIn('MODEL_KEY', worker_env)
            self.assertNotIn('TENANT', worker_env)
            self.assertNotIn('HTTPS_PROXY', worker_env)
            self.assertNotIn('EVAL_OPENCODE_TOKEN', worker_env)
            self.assertEqual(worker_env['GOPROXY'], 'off')
            self.assertNotIn('host-only', json.dumps(safe))
            yield 'local', worker_env
        with patch.object(runner, 'isolated_v2_server', side_effect=server):
            with runner.v2_server('opencode', self.root, self.root / 'events', env):
                pass

    @unittest.skipUnless(os.environ.get('EVAL_TEST_NETWORK_SANDBOX') == '1', 'set EVAL_TEST_NETWORK_SANDBOX=1 for bubblewrap probe')
    def test_real_sandbox_blocks_network_and_keeps_gateway_control_and_go(self):
        spec = importlib.util.spec_from_file_location('runner_network_probe', Path(__file__).with_name('run.py'))
        runner = importlib.util.module_from_spec(spec)
        spec.loader.exec_module(runner)
        workspace = self.root / 'workspace'
        profile = self.root / 'profile'
        library = self.root / 'library'
        for directory in (workspace, profile, library): directory.mkdir()
        (workspace / 'go.mod').write_text('module probe\n\ngo 1.24.0\n')
        (workspace / 'probe.go').write_text('package probe\nfunc Answer() int { return 42 }\n')
        (workspace / 'probe_test.go').write_text('package probe\nimport "testing"\nfunc TestAnswer(t *testing.T) { if Answer()!=42 { t.Fatal("bad") } }\n')
        script = workspace / 'probe.py'
        script.write_text('''import http.client, http.server, json, os, socket, subprocess, sys
assert 'MODEL_KEY' not in os.environ
assert 'HTTPS_PROXY' not in os.environ
for host, port in [('1.1.1.1',443), ('github.com',443), ('127.0.0.1',int(sys.argv[1]))]:
    try:
        socket.create_connection((host, port), timeout=1).close()
    except OSError: pass
    else: raise RuntimeError('external connection succeeded')
for method, route in [('CONNECT','github.com:443'),('GET','https://github.com/'),('GET','/v1/models')]:
    c=http.client.HTTPConnection('127.0.0.1',18080,timeout=5)
    c.request(method,route)
    assert c.getresponse().status==403
    c.close()
subprocess.run(['go','test','./...'],check=True)
class Handler(http.server.BaseHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200); self.end_headers(); self.wfile.write(b'isolated-control-ok')
http.server.HTTPServer(('127.0.0.1',18081),Handler).handle_request()
''')
        go = json.loads(subprocess.check_output(['go','env','-json','GOROOT','GOMODCACHE']))
        with socket.socket() as host_service:
            host_service.bind(('127.0.0.1',0)); host_service.listen()
            host_port = host_service.getsockname()[1]
            with network.network_runtime(fixture(), 'fixture/model', self.env) as (_, directory, port):
                env = {**go, 'EVAL_NETWORK_DIR': directory, 'EVAL_LIBRARY_ROOT': str(library),
                       'PATH': os.environ['PATH'], 'HOME': str(profile), 'GOCACHE': str(profile/'go-cache'),
                       'GOPROXY':'off', 'GOSUMDB':'off', 'GOTOOLCHAIN':'local', 'RES_OPTIONS':'attempts:1 timeout:1'}
                argv = runner.sandbox_command('/usr/bin/python3',workspace,profile,env)
                process = subprocess.Popen(argv+[str(script),str(host_port)],env=env,stdout=subprocess.PIPE,stderr=subprocess.PIPE)
                try:
                    deadline=time.monotonic()+90
                    while time.monotonic()<deadline:
                        if process.poll() is not None:
                            out,err=process.communicate(); self.fail((out+err).decode())
                        try:
                            c=http.client.HTTPConnection('127.0.0.1',port,timeout=1)
                            c.request('GET','/'); response=c.getresponse(); result=response.read(); c.close()
                            self.assertEqual(result,b'isolated-control-ok'); break
                        except (OSError,http.client.HTTPException): time.sleep(.1)
                    else: self.fail('control bridge timeout')
                    out,err=process.communicate(timeout=10)
                    self.assertEqual(process.returncode,0,(out+err).decode())
                finally:
                    if process.poll() is None: process.kill();process.wait()


if __name__ == '__main__':
    unittest.main()
