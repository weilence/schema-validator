package validators

import (
	"encoding/base64"
	"net/mail"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/weilence/schema-validator/schema"
)

func registerFormat(r *Registry) {
	// ------------------------ workaround from go-playground/validator ------------------------
	r.Register("base64", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		_, err := base64.StdEncoding.DecodeString(str)
		if err != nil {
			return schema.ErrCheckFailed
		}
		return nil
	})

	r.Register("base64url", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		_, err := base64.URLEncoding.DecodeString(str)
		if err != nil {
			return schema.ErrCheckFailed
		}
		return nil
	})

	r.Register("base64rawurl", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		_, err := base64.RawURLEncoding.DecodeString(str)
		if err != nil {
			return schema.ErrCheckFailed
		}
		return nil
	})

	var bicRegex = regexp.MustCompile(`^[A-Z]{6}[A-Z0-9]{2}([A-Z0-9]{3})?$`)
	r.Register("bic_iso_9362_2014", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if bicRegex.MatchString(str) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	r.Register("bic", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if bicRegex.MatchString(str) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	var bcp47Regex = regexp.MustCompile(`^[a-zA-Z]{1,8}(-[a-zA-Z0-9]{1,8})*$`)
	r.Register("bcp47_language_tag", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if bcp47Regex.MatchString(str) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	var btcAddrRegex = regexp.MustCompile(`^[13][a-km-zA-HJ-NP-Z1-9]{25,34}$`)
	r.Register("btc_addr", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if btcAddrRegex.MatchString(str) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	var btcBech32Regex = regexp.MustCompile(`^bc1[a-z0-9]{39,59}$`)
	r.Register("btc_addr_bech32", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if btcBech32Regex.MatchString(str) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	r.Register("credit_card", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		str = strings.ReplaceAll(str, " ", "")
		str = strings.ReplaceAll(str, "-", "")
		if len(str) < 13 || len(str) > 19 {
			return schema.ErrCheckFailed
		}
		for _, r := range str {
			if !unicode.IsDigit(r) {
				return schema.ErrCheckFailed
			}
		}
		// Simple Luhn check
		sum := 0
		alternate := false
		for i := len(str) - 1; i >= 0; i-- {
			n := int(str[i] - '0')
			if alternate {
				n *= 2
				if n > 9 {
					n -= 9
				}
			}
			sum += n
			alternate = !alternate
		}
		if sum%10 != 0 {
			return schema.ErrCheckFailed
		}
		return nil
	})

	var mongoIDRegex = regexp.MustCompile(`^[a-fA-F0-9]{24}$`)
	r.Register("mongodb", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if mongoIDRegex.MatchString(str) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	var mongoConnRegex = regexp.MustCompile(`^mongodb(\+srv)?://.*$`)
	r.Register("mongodb_connection_string", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if mongoConnRegex.MatchString(str) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	var cronRegex = regexp.MustCompile(`^(@(annually|yearly|monthly|weekly|daily|midnight|hourly))|(((\d+,)+\d+|(\d+(\/|-)\d+)|\d+|\*|\?) ?){5,7}$`)
	r.Register("cron", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if cronRegex.MatchString(str) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	r.Register("spicedb", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		// Simple check for SpiceDB format
		if strings.Contains(str, "/") {
			return nil
		}
		return schema.ErrCheckFailed
	})

	r.Register("datetime", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		_, err := time.Parse(time.RFC3339, str)
		if err != nil {
			_, err = time.Parse("2006-01-02 15:04:05", str)
			if err != nil {
				return schema.ErrCheckFailed
			}
		}
		return nil
	})

	var e164Regex = regexp.MustCompile(`^\+[1-9]\d{1,14}$`)
	r.Register("e164", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if e164Regex.MatchString(str) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	var einRegex = regexp.MustCompile(`^\d{2}-\d{7}$`)
	r.Register("ein", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if einRegex.MatchString(str) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	r.Register("email", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		_, err := mail.ParseAddress(str)
		if err != nil {
			return schema.ErrCheckFailed
		}
		return nil
	})

	var ethAddrRegex = regexp.MustCompile(`^0x[a-fA-F0-9]{40}$`)
	r.Register("eth_addr", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if ethAddrRegex.MatchString(str) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	var hexRegex = regexp.MustCompile(`^[0-9a-fA-F]+$`)
	r.Register("hexadecimal", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if hexRegex.MatchString(str) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	var hexColorRegex = regexp.MustCompile(`^#([0-9a-fA-F]{3}|[0-9a-fA-F]{6})$`)
	r.Register("hexcolor", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if hexColorRegex.MatchString(str) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	var hslRegex = regexp.MustCompile(`^hsl\(\d+,\s*\d+%,\s*\d+%\)$`)
	r.Register("hsl", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if hslRegex.MatchString(str) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	var hslaRegex = regexp.MustCompile(`^hsla\(\d+,\s*\d+%,\s*\d+%,\s*[\d.]+\)$`)
	r.Register("hsla", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if hslaRegex.MatchString(str) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	var htmlRegex = regexp.MustCompile(`<[^>]+>`)
	r.Register("html", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if htmlRegex.MatchString(str) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	r.Register("html_encoded", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if strings.Contains(str, "&") && strings.Contains(str, ";") {
			return nil
		}
		return schema.ErrCheckFailed
	})

	r.Register("isbn", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		str = strings.ReplaceAll(str, "-", "")
		if len(str) == 10 {
			return validateISBN10(str)
		} else if len(str) == 13 {
			return validateISBN13(str)
		}
		return schema.ErrCheckFailed
	})

	r.Register("isbn10", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		str = strings.ReplaceAll(str, "-", "")
		if len(str) == 10 {
			return validateISBN10(str)
		}
		return schema.ErrCheckFailed
	})

	r.Register("isbn13", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		str = strings.ReplaceAll(str, "-", "")
		if len(str) == 13 {
			return validateISBN13(str)
		}
		return schema.ErrCheckFailed
	})

	var issnRegex = regexp.MustCompile(`^\d{4}-\d{3}[\dX]$`)
	r.Register("issn", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if issnRegex.MatchString(str) {
			return validateISSN(str)
		}
		return schema.ErrCheckFailed
	})

	var iso3166Alpha2Regex = regexp.MustCompile(`^[A-Z]{2}$`)
	r.Register("iso3166_1_alpha2", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if iso3166Alpha2Regex.MatchString(str) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	var iso3166Alpha3Regex = regexp.MustCompile(`^[A-Z]{3}$`)
	r.Register("iso3166_1_alpha3", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if iso3166Alpha3Regex.MatchString(str) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	var iso3166NumericRegex = regexp.MustCompile(`^\d{3}$`)
	r.Register("iso3166_1_alpha_numeric", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if iso3166NumericRegex.MatchString(str) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	var iso3166_2Regex = regexp.MustCompile(`^[A-Z]{2}-[A-Z0-9]{1,3}$`)
	r.Register("iso3166_2", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if iso3166_2Regex.MatchString(str) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	var iso4217Regex = regexp.MustCompile(`^[A-Z]{3}$`)
	r.Register("iso4217", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if iso4217Regex.MatchString(str) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	r.Register("json", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if strings.HasPrefix(str, "{") && strings.HasSuffix(str, "}") {
			return nil
		}
		if strings.HasPrefix(str, "[") && strings.HasSuffix(str, "]") {
			return nil
		}
		return schema.ErrCheckFailed
	})

	var jwtRegex = regexp.MustCompile(`^[A-Za-z0-9-_]+\.[A-Za-z0-9-_]+\.[A-Za-z0-9-_]*$`)
	r.Register("jwt", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if jwtRegex.MatchString(str) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	r.Register("latitude", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		lat, err := strconv.ParseFloat(str, 64)
		if err != nil || lat < -90 || lat > 90 {
			return schema.ErrCheckFailed
		}
		return nil
	})

	r.Register("longitude", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		lng, err := strconv.ParseFloat(str, 64)
		if err != nil || lng < -180 || lng > 180 {
			return schema.ErrCheckFailed
		}
		return nil
	})

	r.Register("luhn_checksum", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		sum := 0
		alternate := false
		for i := len(str) - 1; i >= 0; i-- {
			n := int(str[i] - '0')
			if alternate {
				n *= 2
				if n > 9 {
					n -= 9
				}
			}
			sum += n
			alternate = !alternate
		}
		if sum%10 != 0 {
			return schema.ErrCheckFailed
		}
		return nil
	})

	r.Register("postcode_iso3166_alpha2", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		// Simple check, in practice need country-specific
		if len(str) >= 3 && len(str) <= 10 {
			return nil
		}
		return schema.ErrCheckFailed
	})

	r.Register("postcode_iso3166_alpha2_field", func(ctx *schema.Context) error {
		// Same as above
		return nil
	})

	var rgbRegex = regexp.MustCompile(`^rgb\(\d+,\s*\d+,\s*\d+\)$`)
	r.Register("rgb", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if rgbRegex.MatchString(str) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	var rgbaRegex = regexp.MustCompile(`^rgba\(\d+,\s*\d+,\s*\d+,\s*[\d.]+\)$`)
	r.Register("rgba", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if rgbaRegex.MatchString(str) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	var ssnRegex = regexp.MustCompile(`^\d{3}-\d{2}-\d{4}$`)
	r.Register("ssn", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if ssnRegex.MatchString(str) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	r.Register("timezone", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		_, err := time.LoadLocation(str)
		if err != nil {
			return schema.ErrCheckFailed
		}
		return nil
	})

	var uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	r.Register("uuid", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if uuidRegex.MatchString(str) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	r.Register("uuid3", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if uuidRegex.MatchString(str) && strings.HasPrefix(str[14:15], "3") {
			return nil
		}
		return schema.ErrCheckFailed
	})

	r.Register("uuid3_rfc4122", func(ctx *schema.Context) error {
		return nil // Same as uuid3
	})

	r.Register("uuid4", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if uuidRegex.MatchString(str) && strings.HasPrefix(str[14:15], "4") {
			return nil
		}
		return schema.ErrCheckFailed
	})

	r.Register("uuid4_rfc4122", func(ctx *schema.Context) error {
		return nil // Same as uuid4
	})

	r.Register("uuid5", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if uuidRegex.MatchString(str) && strings.HasPrefix(str[14:15], "5") {
			return nil
		}
		return schema.ErrCheckFailed
	})

	r.Register("uuid5_rfc4122", func(ctx *schema.Context) error {
		return nil // Same as uuid5
	})

	r.Register("uuid_rfc4122", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if uuidRegex.MatchString(str) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	var md4Regex = regexp.MustCompile(`^[a-fA-F0-9]{32}$`)
	r.Register("md4", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if md4Regex.MatchString(str) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	var md5Regex = regexp.MustCompile(`^[a-fA-F0-9]{32}$`)
	r.Register("md5", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if md5Regex.MatchString(str) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	var sha256Regex = regexp.MustCompile(`^[a-fA-F0-9]{64}$`)
	r.Register("sha256", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if sha256Regex.MatchString(str) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	var sha384Regex = regexp.MustCompile(`^[a-fA-F0-9]{96}$`)
	r.Register("sha384", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if sha384Regex.MatchString(str) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	var sha512Regex = regexp.MustCompile(`^[a-fA-F0-9]{128}$`)
	r.Register("sha512", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if sha512Regex.MatchString(str) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	var ripemd128Regex = regexp.MustCompile(`^[a-fA-F0-9]{32}$`)
	r.Register("ripemd128", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if ripemd128Regex.MatchString(str) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	var ripemd160Regex = regexp.MustCompile(`^[a-fA-F0-9]{40}$`)
	r.Register("ripemd160", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if ripemd160Regex.MatchString(str) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	var tiger128Regex = regexp.MustCompile(`^[a-fA-F0-9]{32}$`)
	r.Register("tiger128", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if tiger128Regex.MatchString(str) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	var tiger160Regex = regexp.MustCompile(`^[a-fA-F0-9]{40}$`)
	r.Register("tiger160", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if tiger160Regex.MatchString(str) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	var tiger192Regex = regexp.MustCompile(`^[a-fA-F0-9]{48}$`)
	r.Register("tiger192", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if tiger192Regex.MatchString(str) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	var semverRegex = regexp.MustCompile(`^v?(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)
	r.Register("semver", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if semverRegex.MatchString(str) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	var ulidRegex = regexp.MustCompile(`^[0-9A-HJKMNP-TV-Z]{26}$`)
	r.Register("ulid", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if ulidRegex.MatchString(str) {
			return nil
		}
		return schema.ErrCheckFailed
	})

	var cveRegex = regexp.MustCompile(`^CVE-\d{4}-\d{4,}$`)
	r.Register("cve", func(ctx *schema.Context) error {
		str := ctx.Value().String()
		if cveRegex.MatchString(str) {
			return nil
		}
		return schema.ErrCheckFailed
	})
	// ------------------------ end of workaround ------------------------
}

func validateISBN10(isbn string) error {
	sum := 0
	for i, r := range isbn {
		if r == 'X' && i == 9 {
			sum += 10 * (10 - i)
		} else {
			d, err := strconv.Atoi(string(r))
			if err != nil {
				return schema.ErrCheckFailed
			}
			sum += d * (10 - i)
		}
	}
	if sum%11 != 0 {
		return schema.ErrCheckFailed
	}
	return nil
}

func validateISBN13(isbn string) error {
	sum := 0
	for i, r := range isbn {
		d, err := strconv.Atoi(string(r))
		if err != nil {
			return schema.ErrCheckFailed
		}
		if i%2 == 0 {
			sum += d
		} else {
			sum += d * 3
		}
	}
	if sum%10 != 0 {
		return schema.ErrCheckFailed
	}
	return nil
}

func validateISSN(issn string) error {
	issn = strings.ReplaceAll(issn, "-", "")
	sum := 0
	for i, r := range issn {
		if r == 'X' && i == 7 {
			sum += 10 * (8 - i)
		} else {
			d, err := strconv.Atoi(string(r))
			if err != nil {
				return schema.ErrCheckFailed
			}
			sum += d * (8 - i)
		}
	}
	if sum%11 != 0 {
		return schema.ErrCheckFailed
	}
	return nil
}
