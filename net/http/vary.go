package http

import "github.com/alexfalkowski/go-service/v2/strings"

// AddVary appends the given request header fields to Vary without duplicating existing fields.
//
// Vary field names are compared case-insensitively. A Vary value of "*" already covers all request
// headers, so it is left unchanged.
func AddVary(header Header, fields ...string) {
	existing := map[string]struct{}{}
	for _, value := range header.Values(VaryKey) {
		for field := range strings.SplitSeq(value, ",") {
			existing[strings.ToLower(strings.TrimSpace(field))] = struct{}{}
		}
	}

	if _, ok := existing["*"]; ok {
		return
	}

	for _, field := range fields {
		key := strings.ToLower(strings.TrimSpace(field))
		if strings.IsEmpty(key) {
			continue
		}
		if _, ok := existing[key]; ok {
			continue
		}

		header.Add(VaryKey, field)
		existing[key] = struct{}{}
	}
}
