package validate_test

import (
	"testing"

	"github.com/muonsoft/validation/validate"
	"github.com/stretchr/testify/assert"
)

func TestURL_WhenDefaultSchemasAndValidURL_ExpectNoError(t *testing.T) {
	for _, url := range validURLs() {
		t.Run(url, func(t *testing.T) {
			err := validate.URL(url)

			assert.NoError(t, err)
		})
	}
}

func TestURL_WhenDefaultSchemasAndInvalidURL_ExpectError(t *testing.T) {
	for _, url := range invalidURLs() {
		t.Run(url, func(t *testing.T) {
			err := validate.URL(url)

			assert.Error(t, err)
		})
	}
}

func TestURL_WhenCustomSchemasAndValidURL_ExpectNoError(t *testing.T) {
	for _, url := range validURLsWithCustomSchemas() {
		t.Run(url, func(t *testing.T) {
			err := validate.URL(url, "ftp", "file", "git")

			assert.NoError(t, err)
		})
	}
}

func TestURL_WhenDefaultSchemasAndRelativeSchemaAndValidURL_ExpectNoError(t *testing.T) {
	urls := append(validURLs(), validRelativeURLs()...)
	for _, url := range urls {
		t.Run(url, func(t *testing.T) {
			err := validate.URL(url, "http", "https", "")

			assert.NoError(t, err)
		})
	}
}

func TestRelativeURL_WhenDefaultSchemasAndRelativeSchemaAndInvalidURL_ExpectError(t *testing.T) {
	urls := append(invalidURLs(), invalidRelativeURLs()...)
	for _, url := range urls {
		t.Run(url, func(t *testing.T) {
			err := validate.URL(url, "http", "https", "")

			assert.Error(t, err)
		})
	}
}

func validURLs() []string {
	return []string{
		"http://a.pl",
		"http://www.example.com",
		"HTTP://WWW.EXAMPLE.COM",
		"http://www.example.com.",
		"http://www.example.museum",
		"https://example.com/",
		"https://example.com:80/",
		"http://examp_le.com",
		"http://www.sub_domain.examp_le.com",
		"http://www.example.coop/",
		"http://www.test-example.com/",
		"http://www.example.com/",
		"http://example.fake/blog/",
		"http://example.com/?",
		"http://example.com/search?type=&q=url+validator",
		"http://example.com/#",
		"http://example.com/#?",
		"http://www.example.com/doc/current/book/validation.html#supported-constraints",
		"http://very.long.domain.name.com/",
		"http://localhost/",
		"http://myhost123/",
		"http://127.0.0.1/",
		"http://127.0.0.1:80/",
		"http://[::1]/",
		"http://[::1]:80/",
		"http://[1:2:3::4:5:6:7]/",
		"http://sãopaulo.com/",
		"http://xn--sopaulo-xwa.com/",
		"http://sãopaulo.com.br/",
		"http://xn--sopaulo-xwa.com.br/",
		"http://пример.испытание/",
		"http://xn--e1afmkfd.xn--80akhbyknj4f/",
		"http://مثال.إختبار/",
		"http://xn--mgbh0fb.xn--kgbechtv/",
		"http://例子.测试/",
		"http://xn--fsqu00a.xn--0zwm56d/",
		"http://例子.測試/",
		"http://xn--fsqu00a.xn--g6w251d/",
		"http://例え.テスト/",
		"http://xn--r8jz45g.xn--zckzah/",
		"http://مثال.آزمایشی/",
		"http://xn--mgbh0fb.xn--hgbk6aj7f53bba/",
		"http://실례.테스트/",
		"http://xn--9n2bp8q.xn--9t4b11yi5a/",
		"http://العربية.idn.icann.org/",
		"http://xn--ogb.idn.icann.org/",
		"http://xn--e1afmkfd.xn--80akhbyknj4f.xn--e1afmkfd/",
		"http://xn--espaa-rta.xn--ca-ol-fsay5a/",
		"http://xn--d1abbgf6aiiy.xn--p1ai/",
		"http://☎.com/",
		"http://username:password@example.com",
		"http://user.name:password@example.com",
		"http://user_name:pass_word@example.com",
		"http://username:pass.word@example.com",
		"http://user.name:pass.word@example.com",
		"http://user-name@example.com",
		"http://user_name@example.com",
		"http://u%24er:password@example.com",
		"http://user:pa%24%24word@example.com",
		"http://example.com?",
		"http://example.com?query=1",
		"http://example.com/?query=1",
		"http://example.com/?querie%24=1",
		"http://example.com/path%24/?query=1",
		"http://example.com#",
		"http://example.com#fragment",
		"http://example.com/#fragment",
		"http://example.com/#one_more%20test",
		"http://example.com/exploit.html?hello[0]=test",
	}
}

func validRelativeURLs() []string {
	return []string{
		"//example.com",
		"//examp_le.com",
		"//example.fake/blog/",
		"//example.com/search?type=&q=url+validator",
	}
}

func invalidURLs() []string {
	return []string{
		"",
		"example.com",
		"://example.com",
		"http ://example.com",
		"http:/example.com",
		"http://example.com::aa",
		"http://example.com:aa",
		"ftp://example.fr",
		"faked://example.fr",
		"http://127.0.0.1:aa/",
		"ftp://[::1]/",
		"http://[::1",
		"http://hello.☎/",
		"http://:password@example.com",
		"http://:password@@example.com",
		"http://username:passwordexample.com",
		"http://usern@me:password@example.com",
		"http://nota%hex:password@example.com",
		"http://username:nota%hex@example.com",
		"http://example.com/exploit.html?<script>alert(1);</script>",
		"http://example.com/exploit.html?hel lo",
		"http://example.com/exploit.html?not_a%hex",
		"http://",
	}
}

func invalidRelativeURLs() []string {
	return []string{
		"/example.com",
		"//example.com::aa",
		"//example.com:aa",
		"//127.0.0.1:aa/",
		"//[::1",
		"//hello.☎/",
		"//:password@example.com",
		"//:password@@example.com",
		"//username:passwordexample.com",
		"//usern@me:password@example.com",
		"//example.com/exploit.html?<script>alert(1);</script>",
		"//example.com/exploit.html?hel lo",
		"//example.com/exploit.html?not_a%hex",
		"//",
	}
}

func validURLsWithCustomSchemas() []string {
	return []string{
		"ftp://example.com",
		"file://127.0.0.1",
		"git://[::1]/",
	}
}

func TestIP_WhenValidIP_ExpectNoError(t *testing.T) {
	ips := append(validIPsV4(), validIPsV6()...)
	for _, ip := range ips {
		t.Run(ip, func(t *testing.T) {
			err := validate.IP(ip)

			assert.NoError(t, err)
		})
	}
}

func TestIP_WhenInvalidIP_ExpectError(t *testing.T) {
	ips := append(invalidIPsV4(), invalidIPsV6()...)
	for _, ip := range ips {
		t.Run(ip, func(t *testing.T) {
			err := validate.IP(ip)

			assert.Error(t, err)
		})
	}
}

func TestIPv4_WhenValidIP_ExpectNoError(t *testing.T) {
	for _, ip := range validIPsV4() {
		t.Run(ip, func(t *testing.T) {
			err := validate.IPv4(ip)

			assert.NoError(t, err)
		})
	}
}

func TestIPv4_WhenInvalidIP_ExpectError(t *testing.T) {
	ips := append(append(invalidIPsV4(), validIPsV6()...), invalidIPsV6()...)
	for _, ip := range ips {
		t.Run(ip, func(t *testing.T) {
			err := validate.IPv4(ip)

			assert.Error(t, err)
		})
	}
}

func TestIPv6_WhenValidIP_ExpectNoError(t *testing.T) {
	for _, ip := range validIPsV6() {
		t.Run(ip, func(t *testing.T) {
			err := validate.IPv6(ip)

			assert.NoError(t, err)
		})
	}
}

func TestIPv6_WhenInvalidIP_ExpectError(t *testing.T) {
	ips := append(append(invalidIPsV4(), validIPsV4()...), invalidIPsV6()...)
	for _, ip := range ips {
		t.Run(ip, func(t *testing.T) {
			err := validate.IPv6(ip)

			assert.Error(t, err)
		})
	}
}

func TestIP_WhenDenyPrivateAndInvalidIP_ExpectError(t *testing.T) {
	ips := append(invalidPrivateIPsV4(), invalidIPsV6()...)
	for _, ip := range ips {
		t.Run(ip, func(t *testing.T) {
			err := validate.IP(ip, validate.DenyPrivateIP())

			assert.Error(t, err)
		})
	}
}

func TestIPv4_WhenDenyPrivateAndInvalidIP_ExpectError(t *testing.T) {
	for _, ip := range invalidPrivateIPsV4() {
		t.Run(ip, func(t *testing.T) {
			err := validate.IPv4(ip, validate.DenyPrivateIP())

			assert.Error(t, err)
		})
	}
}

func TestIPv6_WhenDenyPrivateAndInvalidIP_ExpectError(t *testing.T) {
	for _, ip := range invalidPrivateIPsV6() {
		t.Run(ip, func(t *testing.T) {
			err := validate.IPv6(ip, validate.DenyPrivateIP())

			assert.Error(t, err)
		})
	}
}

func validIPsV4() []string {
	return []string{
		"0.0.0.0",
		"10.0.0.0",
		"123.45.67.178",
		"172.16.0.0",
		"192.168.1.0",
		"224.0.0.1",
		"255.255.255.255",
		"127.0.0.0",
	}
}

func validIPsV6() []string {
	return []string{
		"2001:0db8:85a3:0000:0000:8a2e:0370:7334",
		"2001:0DB8:85A3:0000:0000:8A2E:0370:7334",
		"2001:0Db8:85a3:0000:0000:8A2e:0370:7334",
		"fdfe:dcba:9876:ffff:fdc6:c46b:bb8f:7d4c",
		"fdc6:c46b:bb8f:7d4c:fdc6:c46b:bb8f:7d4c",
		"fdc6:c46b:bb8f:7d4c:0000:8a2e:0370:7334",
		"fe80:0000:0000:0000:0202:b3ff:fe1e:8329",
		"fe80:0:0:0:202:b3ff:fe1e:8329",
		"fe80::202:b3ff:fe1e:8329",
		"0:0:0:0:0:0:0:0",
		"::",
		"0::",
		"::0",
		"0::0",
		// IPv4 mapped to IPv6
		"2001:0db8:85a3:0000:0000:8a2e:0.0.0.0",
		"::0.0.0.0",
		"::255.255.255.255",
		"::123.45.67.178",
	}
}

func invalidIPsV6() []string {
	return []string{
		"z001:0db8:85a3:0000:0000:8a2e:0370:7334",
		"fe80",
		"fe80:8329",
		"fe80:::202:b3ff:fe1e:8329",
		"fe80::202:b3ff::fe1e:8329",
		// IPv4 mapped to IPv6
		"2001:0db8:85a3:0000:0000:8a2e:0370:0.0.0.0",
		"::0.0",
		"::0.0.0",
		"::256.0.0.0",
		"::0.256.0.0",
		"::0.0.256.0",
		"::0.0.0.256",
	}
}

func invalidPrivateIPsV6() []string {
	return []string{
		"fdfe:dcba:9876:ffff:fdc6:c46b:bb8f:7d4c",
		"fdc6:c46b:bb8f:7d4c:fdc6:c46b:bb8f:7d4c",
		"fdc6:c46b:bb8f:7d4c:0000:8a2e:0370:7334",
	}
}

func invalidIPsV4() []string {
	return []string{
		"0",
		"0.0",
		"0.0.0",
		"256.0.0.0",
		"0.256.0.0",
		"0.0.256.0",
		"0.0.0.256",
		"-1.0.0.0",
		"foobar",
	}
}

func invalidPrivateIPsV4() []string {
	return []string{
		"10.0.0.0",
		"172.16.0.0",
		"192.168.1.0",
	}
}
