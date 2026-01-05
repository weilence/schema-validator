package validators

import (
	"net"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/miekg/dns"
	"github.com/weilence/schema-validator/schema"
)

func registerNetwork(r *Registry) {
	// ------------------------ workaround from go-playground/validator ------------------------
	// CIDR validators
	r.Register("cidr", func(ctx *schema.Context) error {
		val := ctx.Value().String()
		_, _, err := net.ParseCIDR(val)
		if err != nil {
			return schema.ErrCheckFailed
		}
		return nil
	})

	r.Register("cidrv4", func(ctx *schema.Context) error {
		val := ctx.Value().String()
		ip, _, err := net.ParseCIDR(val)
		if err != nil || ip.To4() == nil {
			return schema.ErrCheckFailed
		}
		return nil
	})

	r.Register("cidrv6", func(ctx *schema.Context) error {
		val := ctx.Value().String()
		ip, _, err := net.ParseCIDR(val)
		if err != nil || ip.To4() != nil {
			return schema.ErrCheckFailed
		}
		return nil
	})

	// Data URI
	var dataURIRegex = regexp.MustCompile(`^data:[^;]+(;base64)?,.*$`)
	r.Register("datauri", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if !dataURIRegex.MatchString(str) {
			return schema.ErrCheckFailed
		}
		return nil
	})

	// FQDN
	r.Register("fqdn", func(ctx *schema.Context) error {
		val := ctx.Value().String()
		if dns.IsFqdn(val) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	// Hostname validators
	var hostnameRegex = regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-.]{0,61}[a-zA-Z0-9])?$`)
	r.Register("hostname", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if hostnameRegex.MatchString(str) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	var hostnameRFC1123Regex = regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-.]{0,61}[a-zA-Z0-9])?$`)
	r.Register("hostname_rfc1123", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if hostnameRFC1123Regex.MatchString(str) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	r.Register("hostname_port", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		host, portStr, err := net.SplitHostPort(str)
		if err != nil {
			return schema.ErrCheckFailed
		}
		if net.ParseIP(host) == nil && !hostnameRegex.MatchString(host) {
			return schema.ErrCheckFailed
		}

		port, err := strconv.Atoi(portStr)
		if err != nil || port < 1 || port > 65535 {
			return schema.ErrCheckFailed
		}

		return nil
	})

	r.Register("port", func(ctx *schema.Context) error {
		field := ctx.Value()

		port, err := field.IntE()
		if err != nil {
			return err
		}

		if port < 1 || port > 65535 {
			return schema.ErrCheckFailed
		}

		return nil
	})

	r.Register("ip", func(ctx *schema.Context) error {
		field := ctx.Value()
		val := field.String()
		if net.ParseIP(val) != nil {
			return nil
		}

		return schema.ErrCheckFailed
	})

	// IP address validators
	r.Register("ip4_addr", func(ctx *schema.Context) error {
		val := ctx.Value().String()
		ip := net.ParseIP(val)
		if ip != nil && ip.To4() != nil {
			return nil
		}
		return schema.ErrCheckFailed
	})

	r.Register("ip6_addr", func(ctx *schema.Context) error {
		val := ctx.Value().String()
		ip := net.ParseIP(val)
		if ip != nil && ip.To4() == nil {
			return nil
		}
		return schema.ErrCheckFailed
	})

	r.Register("ip_addr", func(ctx *schema.Context) error {
		val := ctx.Value().String()
		if net.ParseIP(val) != nil {
			return nil
		}
		return schema.ErrCheckFailed
	})

	r.Register("ipv4", func(ctx *schema.Context) error {
		val := ctx.Value().String()
		ip := net.ParseIP(val)
		if ip != nil && ip.To4() != nil {
			return nil
		}
		return schema.ErrCheckFailed
	})

	r.Register("ipv6", func(ctx *schema.Context) error {
		val := ctx.Value().String()
		ip := net.ParseIP(val)
		if ip != nil && ip.To4() == nil {
			return nil
		}
		return schema.ErrCheckFailed
	})

	// MAC address
	var macRegex = regexp.MustCompile(`^([0-9a-fA-F]{2}[:-]){5}[0-9a-fA-F]{2}$`)
	r.Register("mac", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if macRegex.MatchString(str) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	// TCP/UDP address validators
	r.Register("tcp4_addr", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		host, port, err := net.SplitHostPort(str)
		if err != nil {
			return schema.ErrCheckFailed
		}
		ip := net.ParseIP(host)
		if ip == nil || ip.To4() == nil {
			return schema.ErrCheckFailed
		}
		if _, err := net.LookupPort("tcp", port); err != nil {
			return schema.ErrCheckFailed
		}
		return nil
	})

	r.Register("tcp6_addr", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		host, port, err := net.SplitHostPort(str)
		if err != nil {
			return schema.ErrCheckFailed
		}
		if host[0] == '[' && host[len(host)-1] == ']' {
			host = host[1 : len(host)-1]
		}
		ip := net.ParseIP(host)
		if ip == nil || ip.To4() != nil {
			return schema.ErrCheckFailed
		}
		if _, err := net.LookupPort("tcp", port); err != nil {
			return schema.ErrCheckFailed
		}
		return nil
	})

	r.Register("tcp_addr", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		_, err := net.ResolveTCPAddr("tcp", str)
		if err != nil {
			return schema.ErrCheckFailed
		}
		return nil
	})

	r.Register("udp4_addr", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		host, port, err := net.SplitHostPort(str)
		if err != nil {
			return schema.ErrCheckFailed
		}
		ip := net.ParseIP(host)
		if ip == nil || ip.To4() == nil {
			return schema.ErrCheckFailed
		}
		if _, err := net.LookupPort("udp", port); err != nil {
			return schema.ErrCheckFailed
		}
		return nil
	})

	r.Register("udp6_addr", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		host, port, err := net.SplitHostPort(str)
		if err != nil {
			return schema.ErrCheckFailed
		}
		if host[0] == '[' && host[len(host)-1] == ']' {
			host = host[1 : len(host)-1]
		}
		ip := net.ParseIP(host)
		if ip == nil || ip.To4() != nil {
			return schema.ErrCheckFailed
		}
		if _, err := net.LookupPort("udp", port); err != nil {
			return schema.ErrCheckFailed
		}
		return nil
	})

	r.Register("udp_addr", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		_, err := net.ResolveUDPAddr("udp", str)
		if err != nil {
			return schema.ErrCheckFailed
		}
		return nil
	})

	// Unix address
	r.Register("unix_addr", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if strings.HasPrefix(str, "/") || strings.HasPrefix(str, "@") {
			return nil
		}
		return schema.ErrCheckFailed
	})

	r.Register("uds_exists", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if strings.HasPrefix(str, "@") {
			// Abstract socket, assume exists
			return nil
		}
		if _, err := os.Stat(str); os.IsNotExist(err) {
			return schema.ErrCheckFailed
		}
		return nil
	})

	// URI/URL validators
	r.Register("uri", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if _, err := url.ParseRequestURI(str); err != nil {
			return schema.ErrCheckFailed
		}
		return nil
	})

	var urlRegex = regexp.MustCompile(`^https?://[^\s]+$`)
	r.Register("url", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if !urlRegex.MatchString(str) {
			return schema.ErrCheckFailed
		}

		return nil
	})

	var httpURLRegex = regexp.MustCompile(`^https?://[^\s]+$`)
	r.Register("http_url", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if !httpURLRegex.MatchString(str) {
			return schema.ErrCheckFailed
		}
		return nil
	})

	var httpsURLRegex = regexp.MustCompile(`^https://[^\s]+$`)
	r.Register("https_url", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if !httpsURLRegex.MatchString(str) {
			return schema.ErrCheckFailed
		}
		return nil
	})

	r.Register("url_encoded", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if strings.Contains(str, " ") {
			return schema.ErrCheckFailed
		}

		if _, err := url.QueryUnescape(str); err != nil {
			return schema.ErrCheckFailed
		}

		return nil
	})

	var urnRFC2141Regex = regexp.MustCompile(`^urn:[a-zA-Z0-9][a-zA-Z0-9-]{0,31}:[a-zA-Z0-9()+,.:=@;$_!*'-]+$`)
	r.Register("urn_rfc2141", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if urnRFC2141Regex.MatchString(str) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	// ------------------------ end of workaround ------------------------

	r.Register("domain", func(ctx *schema.Context) error {
		field := ctx.Value()
		val := field.String()
		_, ok := dns.IsDomainName(val)
		if ok {
			return nil
		}

		return schema.ErrCheckFailed
	})

}
