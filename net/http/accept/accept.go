package accept

import (
	"mime"
	"strconv"

	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// Items splits header, an HTTP list-valued header value, into its items.
//
// It splits on the first unquoted comma per RFC 9110 list syntax (https://www.rfc-editor.org/rfc/rfc9110#section-5.6.1),
// so a comma inside a quoted parameter value does not split the list. A backslash inside a quoted value
// escapes the next character, including a quote or comma, per the quoted-pair rule.
//
// Surrounding whitespace around each returned item is trimmed. A malformed value (an unterminated quote)
// is returned as a single item, so Items always returns at least one item.
//
// Examples:
//
//	Items(`application/json`)                                  // ["application/json"]
//	Items(`application/json, text/html`)                       // ["application/json", "text/html"]
//	Items(`application/yaml; profile="a,b", application/toml`) // [`application/yaml; profile="a,b"`, "application/toml"]
func Items(header string) []string {
	var items []string

	quoted := false
	escaped := false
	start := 0

	for index := range header {
		if escaped {
			// The preceding backslash quotes this character, so it cannot change the list state.
			escaped = false
			continue
		}

		if quoted && header[index] == '\\' {
			// A quoted-pair protects the next character, including a quote or comma.
			escaped = true
			continue
		}

		if header[index] == '"' {
			// Only an unescaped quote can enter or leave a quoted parameter value.
			quoted = !quoted
			continue
		}

		if header[index] == ',' && !quoted {
			// An HTTP list comma ends an item only outside a quoted string.
			items = append(items, strings.TrimSpace(header[start:index]))
			start = index + 1
		}
	}

	// Leave a single or malformed trailing item intact so the caller can accept or reject it.
	return append(items, strings.TrimSpace(header[start:]))
}

// First returns the first item of header; see [Items] for the splitting rules.
func First(header string) string {
	return Items(header)[0]
}

// IsZeroQuality reports whether item's q parameter is present and exactly zero — the RFC 9110 §12.4.2
// signal that item is not acceptable at all, distinct from a q value used only to order preference among
// otherwise-acceptable items. A missing, unparsable, or nonzero q is not zero quality; this only answers
// the exclusion question, not full quality-value ordering.
func IsZeroQuality(item string) bool {
	_, params, err := mime.ParseMediaType(item)
	if err != nil {
		return false
	}

	q, ok := params["q"]
	if !ok {
		return false
	}

	value, err := strconv.ParseFloat(q, 64)

	return err == nil && value == 0
}

// IsWildcard reports whether value is a wildcard media range — "*/*" or "type/*" — that is satisfied by
// target, a concrete media type a caller can produce: "*" matches any major type, and a concrete major
// type only matches target's own major type, so "text/*" does not satisfy an "application/x-ndjson"
// target the way "*/*" or "application/*" does.
func IsWildcard(value, target media.Type) bool {
	if value.Subtype() != "*" {
		return false
	}

	major := value.Major()

	return major == "*" || major == target.Major()
}
