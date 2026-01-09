package rule

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/weilence/schema-validator/data"
	"github.com/weilence/schema-validator/schema"
)

func TestFormatValidators(t *testing.T) {
	r := NewRegistry()
	registerFormat(r)

	tests := []struct {
		name     string
		ruleName string
		value    string
		wantErr  bool
	}{
		// base64
		{"base64 valid", "base64", "SGVsbG8=", false},
		{"base64 invalid", "base64", "invalid", true},
		// base64url
		{"base64url valid", "base64url", "SGVsbG8=", false},
		{"base64url invalid", "base64url", "invalid", true},
		// base64rawurl
		{"base64rawurl valid", "base64rawurl", "SGVsbG8", false},
		{"base64rawurl invalid", "base64rawurl", "SGVsbG8=", true},
		// bic
		{"bic valid", "bic", "DEUTDEFF", false},
		{"bic invalid", "bic", "invalid", true},
		// bcp47_language_tag
		{"bcp47_language_tag valid", "bcp47_language_tag", "en-US", false},
		{"bcp47_language_tag invalid", "bcp47_language_tag", "invalid_tag", true},
		// btc_addr
		{"btc_addr valid", "btc_addr", "1BvBMSEYstWetqTFn5Au4m4GFg7xJaNVN2", false},
		{"btc_addr invalid", "btc_addr", "invalid", true},
		// btc_addr_bech32
		{"btc_addr_bech32 valid", "btc_addr_bech32", "bc1qw508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t4", false},
		{"btc_addr_bech32 invalid", "btc_addr_bech32", "invalid", true},
		// credit_card
		{"credit_card valid", "credit_card", "4111111111111111", false},
		{"credit_card invalid", "credit_card", "1234567890123456", true},
		// mongodb
		{"mongodb valid", "mongodb", "507f1f77bcf86cd799439011", false},
		{"mongodb invalid", "mongodb", "invalid", true},
		// cron
		{"cron valid", "cron", "* * * * *", false},
		{"cron invalid", "cron", "invalid", true},
		// datetime
		{"datetime valid", "datetime", "2023-01-01T00:00:00Z", false},
		{"datetime invalid", "datetime", "invalid", true},
		// e164
		{"e164 valid", "e164", "+1234567890", false},
		{"e164 invalid", "e164", "1234567890", true},
		// ein
		{"ein valid", "ein", "12-3456789", false},
		{"ein invalid", "ein", "invalid", true},
		// email
		{"email valid", "email", "test@example.com", false},
		{"email invalid", "email", "invalid", true},
		// eth_addr
		{"eth_addr valid", "eth_addr", "0x742d35Cc6634C0532925a3b844Bc454e4438f44e", false},
		{"eth_addr invalid", "eth_addr", "invalid", true},
		// hexadecimal
		{"hexadecimal valid", "hexadecimal", "1a2b3c", false},
		{"hexadecimal invalid", "hexadecimal", "1g2h3i", true},
		// hexcolor
		{"hexcolor valid", "hexcolor", "#ffffff", false},
		{"hexcolor invalid", "hexcolor", "#gggggg", true},
		// hsl
		{"hsl valid", "hsl", "hsl(120, 50%, 50%)", false},
		{"hsl invalid", "hsl", "invalid", true},
		// hsla
		{"hsla valid", "hsla", "hsla(120, 50%, 50%, 0.5)", false},
		{"hsla invalid", "hsla", "invalid", true},
		// html
		{"html valid", "html", "<p>hello</p>", false},
		{"html invalid", "html", "hello", true},
		// html_encoded
		{"html_encoded valid", "html_encoded", "hello&amp;world", false},
		{"html_encoded invalid", "html_encoded", "hello", true},
		// isbn10
		{"isbn10 valid", "isbn10", "0306406152", false},
		{"isbn10 invalid", "isbn10", "invalid", true},
		// isbn13
		{"isbn13 valid", "isbn13", "9780306406157", false},
		{"isbn13 invalid", "isbn13", "invalid", true},
		// issn
		{"issn valid", "issn", "2049-3630", false},
		{"issn invalid", "issn", "invalid", true},
		// iso3166_1_alpha2
		{"iso3166_1_alpha2 valid", "iso3166_1_alpha2", "US", false},
		{"iso3166_1_alpha2 invalid", "iso3166_1_alpha2", "invalid", true},
		// iso3166_1_alpha3
		{"iso3166_1_alpha3 valid", "iso3166_1_alpha3", "USA", false},
		{"iso3166_1_alpha3 invalid", "iso3166_1_alpha3", "invalid", true},
		// iso3166_1_alpha_numeric
		{"iso3166_1_alpha_numeric valid", "iso3166_1_alpha_numeric", "840", false},
		{"iso3166_1_alpha_numeric invalid", "iso3166_1_alpha_numeric", "invalid", true},
		// iso3166_2
		{"iso3166_2 valid", "iso3166_2", "US-CA", false},
		{"iso3166_2 invalid", "iso3166_2", "invalid", true},
		// iso4217
		{"iso4217 valid", "iso4217", "USD", false},
		{"iso4217 invalid", "iso4217", "invalid", true},
		// json
		{"json valid", "json", "{\"key\": \"value\"}", false},
		{"json invalid", "json", "invalid", true},
		// jwt
		{"jwt valid", "jwt", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c", false},
		{"jwt invalid", "jwt", "invalid", true},
		// latitude
		{"latitude valid", "latitude", "45.0", false},
		{"latitude invalid", "latitude", "91.0", true},
		// longitude
		{"longitude valid", "longitude", "90.0", false},
		{"longitude invalid", "longitude", "181.0", true},
		// luhn_checksum
		{"luhn_checksum valid", "luhn_checksum", "4111111111111111", false},
		{"luhn_checksum invalid", "luhn_checksum", "1234567890123456", true},
		// rgb
		{"rgb valid", "rgb", "rgb(255, 0, 0)", false},
		{"rgb invalid", "rgb", "invalid", true},
		// rgba
		{"rgba valid", "rgba", "rgba(255, 0, 0, 0.5)", false},
		{"rgba invalid", "rgba", "invalid", true},
		// ssn
		{"ssn valid", "ssn", "123-45-6789", false},
		{"ssn invalid", "ssn", "invalid", true},
		// timezone
		{"timezone valid", "timezone", "America/New_York", false},
		{"timezone invalid", "timezone", "invalid", true},
		// uuid
		{"uuid valid", "uuid", "550e8400-e29b-41d4-a716-446655440000", false},
		{"uuid invalid", "uuid", "invalid", true},
		// uuid4
		{"uuid4 valid", "uuid4", "550e8400-e29b-41d4-a716-446655440000", false},
		{"uuid4 invalid", "uuid4", "550e8400-e29b-11d4-a716-446655440000", true}, // not version 4
		// md5
		{"md5 valid", "md5", "9e107d9d372bb6826bd81d3542a419d6", false},
		{"md5 invalid", "md5", "invalid", true},
		// sha256
		{"sha256 valid", "sha256", "a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3", false},
		{"sha256 invalid", "sha256", "invalid", true},
		// semver
		{"semver valid", "semver", "1.0.0", false},
		{"semver invalid", "semver", "invalid", true},
		// ulid
		{"ulid valid", "ulid", "01ARZ3NDEKTSV4RRFFQ69G5FAV", false},
		{"ulid invalid", "ulid", "invalid", true},
		// cve
		{"cve valid", "cve", "CVE-2023-1234", false},
		{"cve invalid", "cve", "invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := schema.NewObject().
				AddField("test", schema.NewField().AddValidator(r.NewValidator(tt.ruleName)))
			ctx := schema.NewContext(s, data.New(map[string]any{"test": tt.value}))
			err := s.Validate(ctx)
			assert.NoError(t, err)
			assert.Equal(t, ctx.Errors().HasErrorCode(tt.ruleName), tt.wantErr)
		})
	}
}
