package url_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v3/bytes"
	"github.com/alexfalkowski/go-service/v3/encoding/hjson"
	"github.com/alexfalkowski/go-service/v3/encoding/json"
	"github.com/alexfalkowski/go-service/v3/encoding/toml"
	"github.com/alexfalkowski/go-service/v3/encoding/yaml"
	"github.com/alexfalkowski/go-service/v3/net/url"
	"github.com/stretchr/testify/require"
)

// TestEscapeAndUnescape checks the distinct rules for path and query components.
func TestEscapeAndUnescape(t *testing.T) {
	t.Parallel()

	const raw = "a/b c+"
	tests := []struct {
		name     string
		escaped  string
		escape   func(string) string
		unescape func(string) (string, error)
	}{
		{name: "path segment", escaped: "a%2Fb%20c+", escape: url.PathEscape, unescape: url.PathUnescape},
		{name: "query component", escaped: "a%2Fb+c%2B", escape: url.QueryEscape, unescape: url.QueryUnescape},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.escaped, tt.escape(raw))
			decoded, err := tt.unescape(tt.escaped)
			require.NoError(t, err)
			require.Equal(t, raw, decoded)

			_, err = tt.unescape("%ZZ")
			require.ErrorIs(t, err, url.EscapeError("%ZZ"))
		})
	}
}

// TestJoinPathPreservesEscapedSegment checks that an escaped slash stays within one segment.
func TestJoinPathPreservesEscapedSegment(t *testing.T) {
	t.Parallel()

	joined, err := url.JoinPath("https://example.com/items/", url.PathEscape("a/b"))
	require.NoError(t, err)
	require.Equal(t, "https://example.com/items/a%2Fb", joined)
}

// TestParseQuery checks that valid parameters survive a malformed escape.
func TestParseQuery(t *testing.T) {
	t.Parallel()

	const query = "tag=a&tag=b&q=a+b"
	want := url.Values{"tag": {"a", "b"}, "q": {"a b"}}

	values, err := url.ParseQuery(query)
	require.NoError(t, err)
	require.Equal(t, want, values)

	values, err = url.ParseQuery(query + "&bad=%ZZ")
	require.ErrorIs(t, err, url.EscapeError("%ZZ"))
	require.Equal(t, want, values)
}

// TestParseRequestURI checks request-target parsing and typed URL errors.
func TestParseRequestURI(t *testing.T) {
	t.Parallel()

	const target = "/items/a%2Fb?q=1"
	parsed, err := url.ParseRequestURI(target)
	require.NoError(t, err)
	require.Equal(t, target, parsed.RequestURI())

	_, err = url.ParseRequestURI("items/a")
	var parseErr *url.Error
	require.ErrorAs(t, err, &parseErr)
	require.Equal(t, "items/a", parseErr.URL)

	_, err = url.ParseRequestURI("https://example.com|/items")
	require.ErrorIs(t, err, url.InvalidHostError("|"))
}

// TestUserinfo checks credentials constructed for a URL.
func TestUserinfo(t *testing.T) {
	t.Parallel()

	address := url.URL{Scheme: "https", Host: "example.com", User: url.User("alice")}
	require.Equal(t, "https://alice@example.com", address.String())

	address.User = url.UserPassword("alice", "p@ss")
	require.Equal(t, "https://alice:p%40ss@example.com", address.String())
}

func TestSplitPathSeparatesBaseAndTail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		path  string
		first string
		rest  string
		ok    bool
	}{
		{name: "pair", path: "/service/method", first: "service", rest: "method", ok: true},
		{name: "deep rest", path: "/service/users/123", first: "service", rest: "users/123", ok: true},
		{name: "missing leading slash", path: "service/method"},
		{name: "missing rest", path: "/service"},
		{name: "empty first", path: "//method"},
		{name: "empty rest", path: "/service/"},
		{name: "root", path: "/"},
		{name: "empty", path: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			first, rest, ok := url.SplitPath(tt.path)
			require.Equal(t, tt.first, first)
			require.Equal(t, tt.rest, rest)
			require.Equal(t, tt.ok, ok)
		})
	}
}

func TestParseSecureURL(t *testing.T) {
	t.Parallel()

	const raw = "https://[::1]:8443/a%2Fb?x=1#section"
	secure, err := url.ParseSecureURL(raw)
	require.NoError(t, err)
	require.Equal(t, raw, secure.String())

	parsed := secure.URL()
	require.Equal(t, "::1", parsed.Hostname())
	parsed.Scheme = "http"
	require.Equal(t, raw, secure.String())
}

func TestParseSecureURLAcceptsLoopbackHTTP(t *testing.T) {
	t.Parallel()

	tests := []string{
		"http://LOCALHOST:8080/health",
		"http://127.0.0.2:8080/health",
		"http://[::1]:8080/health",
		"http://[::ffff:127.0.0.1]:8080/health",
	}
	for _, raw := range tests {
		t.Run(raw, func(t *testing.T) {
			t.Parallel()

			secure, err := url.ParseSecureURL(raw)
			require.NoError(t, err)
			require.Equal(t, raw, secure.String())
		})
	}
}

func TestParseSecureURLRejectsInvalidInput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		raw      string
		insecure bool
	}{
		{name: "empty"},
		{name: "cleartext hostname", raw: "http://example.com", insecure: true},
		{name: "cleartext private IP", raw: "http://192.168.1.10", insecure: true},
		{name: "cleartext nonlocal IP", raw: "http://203.0.113.10", insecure: true},
		{name: "cleartext wildcard IP", raw: "http://0.0.0.0", insecure: true},
		{name: "localhost in userinfo", raw: "http://localhost@example.com", insecure: true},
		{name: "localhost in domain", raw: "http://localhost.example.com", insecure: true},
		{name: "different scheme on loopback", raw: "ftp://localhost"},
		{name: "relative", raw: "//example.com/path"},
		{name: "opaque", raw: "https:example.com"},
		{name: "missing host", raw: "https:///path"},
		{name: "port only", raw: "https://:443/path"},
		{name: "invalid escape", raw: "https://example.com/%ZZ?token=example"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			secure, err := url.ParseSecureURL(tt.raw)
			require.ErrorIs(t, err, url.ErrInvalidSecureURL)
			require.Nil(t, secure)
			if tt.insecure {
				require.EqualError(t, err, url.ErrInsecureURL.Error())
				require.ErrorIs(t, err, url.ErrInsecureURL)
				return
			}
			require.EqualError(t, err, url.ErrInvalidSecureURL.Error())
		})
	}
}

func TestMustParseSecureURLPanicsForInvalidInput(t *testing.T) {
	t.Parallel()

	require.PanicsWithValue(t, url.ErrInsecureURL, func() { url.MustParseSecureURL("http://example.com") })
}

func TestSecureURLUnmarshalTextCopiesInput(t *testing.T) {
	t.Parallel()

	const raw = "https://example.com/path"
	text := bytes.NewBufferString(raw).Bytes()
	var secure url.SecureURL
	require.NoError(t, secure.UnmarshalText(text))

	text[len(text)-1] = 'X'
	require.Equal(t, raw, secure.String())
}

func TestSecureURLConfigRoundTrip(t *testing.T) {
	t.Parallel()

	const raw = "https://example.com/a%2Fb?x=1&y=2"
	tests := []struct {
		name      string
		valid     string
		invalid   string
		marshal   func(any) ([]byte, error)
		unmarshal func([]byte, any) error
	}{
		{name: "json", valid: `{"url":"https:\/\/example.com/a%2Fb?x=1&y=2"}`, invalid: `{"url":{}}`, marshal: json.Marshal, unmarshal: json.Unmarshal},
		{name: "hjson", valid: `url: "` + raw + `"`, invalid: `url: "http://example.com"`, marshal: hjson.Marshal, unmarshal: hjson.Unmarshal},
		{name: "toml", valid: `url = "` + raw + `"`, invalid: `url = ""`, marshal: toml.Marshal, unmarshal: toml.Unmarshal},
		{name: "yaml", valid: `url: ` + raw, invalid: `url: {}`, marshal: yaml.Marshal, unmarshal: yaml.Unmarshal},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var cfg struct {
				URL *url.SecureURL `json:"url" toml:"url" yaml:"url"`
			}
			require.NoError(t, tt.unmarshal([]byte(tt.valid), &cfg))
			require.NotNil(t, cfg.URL)
			require.Equal(t, raw, cfg.URL.String())

			encoded, err := tt.marshal(cfg)
			require.NoError(t, err)
			decoded := cfg
			decoded.URL = nil
			require.NoError(t, tt.unmarshal(encoded, &decoded))
			require.Equal(t, cfg, decoded)

			require.Error(t, tt.unmarshal([]byte(tt.invalid), &cfg))
			require.Equal(t, raw, cfg.URL.String())
		})
	}
}

func TestSecureURLZeroValueIsUnset(t *testing.T) {
	t.Parallel()

	var secure url.SecureURL
	require.True(t, secure.IsZero())
	require.Empty(t, secure.String())
	require.Nil(t, secure.URL())

	_, err := secure.MarshalText()
	require.ErrorIs(t, err, url.ErrInvalidSecureURL)
}
