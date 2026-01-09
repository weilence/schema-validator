package rule

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/weilence/schema-validator/data"
	"github.com/weilence/schema-validator/schema"
)

func TestNetworkValidators(t *testing.T) {
	r := NewRegistry()
	registerNetwork(r)

	tests := []struct {
		name     string
		ruleName string
		value    string
		params   []any
		wantErr  bool
	}{
		// ip
		{"ip valid", "ip", "192.168.1.1", nil, false},
		{"ip invalid", "ip", "invalid", nil, true},
		// port
		{"port valid", "port", "8080", nil, false},
		{"port invalid", "port", "99999", nil, true},
		// domain
		{"domain valid", "domain", "example.com", nil, false},
		{"domain invalid", "domain", "invalid..com", nil, true},
		// url
		{"url valid", "url", "http://example.com", nil, false},
		{"url invalid", "url", "not a url", nil, true},
		// cidr
		{"cidr valid", "cidr", "192.168.1.0/24", nil, false},
		{"cidr invalid", "cidr", "invalid", nil, true},
		// cidrv4
		{"cidrv4 valid", "cidrv4", "192.168.1.0/24", nil, false},
		{"cidrv4 invalid", "cidrv4", "2001:db8::/32", nil, true},
		// cidrv6
		{"cidrv6 valid", "cidrv6", "2001:db8::/32", nil, false},
		{"cidrv6 invalid", "cidrv6", "192.168.1.0/24", nil, true},
		// datauri
		{"datauri valid", "datauri", "data:text/plain;base64,SGVsbG8=", nil, false},
		{"datauri invalid", "datauri", "invalid", nil, true},
		// fqdn
		{"fqdn valid", "fqdn", "example.com.", nil, false},
		{"fqdn invalid", "fqdn", "invalid..com", nil, true},
		// hostname
		{"hostname valid", "hostname", "localhost", nil, false},
		{"hostname invalid", "hostname", "invalid%host", nil, true},
		// hostname_rfc1123
		{"hostname_rfc1123 valid", "hostname_rfc1123", "example", nil, false},
		{"hostname_rfc1123 invalid", "hostname_rfc1123", "invalid_host", nil, true},
		// hostname_port
		{"hostname_port valid", "hostname_port", "example.com:8080", nil, false},
		{"hostname_port invalid", "hostname_port", "invalid:99999", nil, true},
		// ip4_addr
		{"ip4_addr valid", "ip4_addr", "192.168.1.1", nil, false},
		{"ip4_addr invalid", "ip4_addr", "2001:db8::1", nil, true},
		// ip6_addr
		{"ip6_addr valid", "ip6_addr", "2001:db8::1", nil, false},
		{"ip6_addr invalid", "ip6_addr", "192.168.1.1", nil, true},
		// ip_addr
		{"ip_addr valid", "ip_addr", "192.168.1.1", nil, false},
		{"ip_addr invalid", "ip_addr", "invalid", nil, true},
		// ipv4
		{"ipv4 valid", "ipv4", "192.168.1.1", nil, false},
		{"ipv4 invalid", "ipv4", "2001:db8::1", nil, true},
		// ipv6
		{"ipv6 valid", "ipv6", "2001:db8::1", nil, false},
		{"ipv6 invalid", "ipv6", "192.168.1.1", nil, true},
		// mac
		{"mac valid", "mac", "00:11:22:33:44:55", nil, false},
		{"mac invalid", "mac", "invalid", nil, true},
		// tcp4_addr
		{"tcp4_addr valid", "tcp4_addr", "192.168.1.1:8080", nil, false},
		{"tcp4_addr invalid", "tcp4_addr", "2001:db8::1:8080", nil, true},
		// tcp6_addr
		{"tcp6_addr valid", "tcp6_addr", "[2001:db8::1]:8080", nil, false},
		{"tcp6_addr invalid", "tcp6_addr", "192.168.1.1:8080", nil, true},
		// tcp_addr
		{"tcp_addr valid", "tcp_addr", "192.168.1.1:8080", nil, false},
		{"tcp_addr invalid", "tcp_addr", "invalid:8080", nil, true},
		// udp4_addr
		{"udp4_addr valid", "udp4_addr", "192.168.1.1:8080", nil, false},
		{"udp4_addr invalid", "udp4_addr", "2001:db8::1:8080", nil, true},
		// udp6_addr
		{"udp6_addr valid", "udp6_addr", "[2001:db8::1]:8080", nil, false},
		{"udp6_addr invalid", "udp6_addr", "192.168.1.1:8080", nil, true},
		// udp_addr
		{"udp_addr valid", "udp_addr", "192.168.1.1:8080", nil, false},
		{"udp_addr invalid", "udp_addr", "invalid:8080", nil, true},
		// unix_addr
		{"unix_addr valid", "unix_addr", "/tmp/socket", nil, false},
		{"unix_addr invalid", "unix_addr", "invalid", nil, true},
		// uds_exists (assuming /tmp/test.sock exists or not)
		{"uds_exists valid", "uds_exists", "@abstract", nil, false}, // abstract socket
		{"uds_exists invalid", "uds_exists", "/nonexistent", nil, true},
		// uri
		{"uri valid", "uri", "http://example.com", nil, false},
		{"uri invalid", "uri", "invalid uri", nil, true},
		// http_url
		{"http_url valid", "http_url", "http://example.com", nil, false},
		{"http_url invalid", "http_url", "ftp://example.com", nil, true},
		// https_url
		{"https_url valid", "https_url", "https://example.com", nil, false},
		{"https_url invalid", "https_url", "http://example.com", nil, true},
		// url_encoded
		{"url_encoded valid", "url_encoded", "hello%20world", nil, false},
		{"url_encoded invalid", "url_encoded", "hello world", nil, true},
		// urn_rfc2141
		{"urn_rfc2141 valid", "urn_rfc2141", "urn:ietf:rfc:2648", nil, false},
		{"urn_rfc2141 invalid", "urn_rfc2141", "invalid", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := schema.NewObject().
				AddField("test", schema.NewField().AddValidator(r.NewValidator(tt.ruleName, tt.params...)))
			ctx := schema.NewContext(s, data.New(map[string]any{"test": tt.value}))
			err := s.Validate(ctx)
			assert.NoError(t, err)
			assert.Equal(t, ctx.Errors().HasErrorCode(tt.ruleName), tt.wantErr)
		})
	}
}
