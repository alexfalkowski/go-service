package content

import "github.com/alexfalkowski/go-service/v2/strings"

// firstListItem returns the first item of an HTTP list-valued header value.
//
// It splits on the first unquoted comma per RFC 9110 list syntax (https://www.rfc-editor.org/rfc/rfc9110#section-5.6.1),
// so a comma inside a quoted parameter value does not end the first item. A backslash inside a quoted value
// escapes the next character, including a quote or comma, per the quoted-pair rule.
//
// Surrounding whitespace around the returned item is trimmed.
//
// Examples:
//
//	firstListItem(`application/json`)                                  // "application/json"
//	firstListItem(`application/json, text/html`)                       // "application/json"
//	firstListItem(`application/yaml; profile="a,b", application/toml`) // `application/yaml; profile="a,b"`
//	firstListItem(`application/yaml; profile="a\",b", application/toml`) // `application/yaml; profile="a\",b"`
//	firstListItem(`application/yaml; profile="a,b`)                    // `application/yaml; profile="a,b` (malformed, returned as-is)
func firstListItem(value string) string {
	quoted := false
	escaped := false

	for index := range value {
		if escaped {
			// The preceding backslash quotes this character, so it cannot change the list state.
			escaped = false
			continue
		}

		if quoted && value[index] == '\\' {
			// A quoted-pair protects the next character, including a quote or comma.
			escaped = true
			continue
		}

		if value[index] == '"' {
			// Only an unescaped quote can enter or leave a quoted parameter value.
			quoted = !quoted
			continue
		}

		if value[index] == ',' && !quoted {
			// An HTTP list comma ends the first item only outside a quoted string.
			return strings.TrimSpace(value[:index])
		}
	}

	// Leave a single or malformed item intact so the caller can accept or reject it.
	return strings.TrimSpace(value)
}
