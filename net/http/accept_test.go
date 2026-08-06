package http_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/media"
	"github.com/stretchr/testify/require"
)

func TestAcceptItems(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		header   string
		expected []string
	}{
		{name: "single item", header: "application/json", expected: []string{"application/json"}},
		{name: "list", header: "application/json, text/html", expected: []string{"application/json", "text/html"}},
		{
			name:     "quoted comma",
			header:   `application/yaml; profile="a,b", application/toml`,
			expected: []string{`application/yaml; profile="a,b"`, "application/toml"},
		},
		{
			name:     "escaped quoted comma",
			header:   `application/yaml; profile="a\",b", application/toml`,
			expected: []string{`application/yaml; profile="a\",b"`, "application/toml"},
		},
		{name: "malformed quoted parameter", header: `application/yaml; profile="a,b`, expected: []string{`application/yaml; profile="a,b`}},
		{name: "empty", header: "", expected: []string{""}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.expected, http.AcceptItems(tt.header))
		})
	}
}

func TestFirstAcceptItem(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		header   string
		expected string
	}{
		{name: "single item", header: "application/json", expected: "application/json"},
		{name: "list", header: "application/json, text/html", expected: "application/json"},
		{
			name:     "quoted comma",
			header:   `application/yaml; profile="a,b", application/toml`,
			expected: `application/yaml; profile="a,b"`,
		},
		{name: "empty", header: "", expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.expected, http.FirstAcceptItem(tt.header))
		})
	}
}

func TestIsAcceptZeroQuality(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		item     string
		expected bool
	}{
		{name: "no q parameter", item: "application/json", expected: false},
		{name: "zero quality", item: "application/json;q=0", expected: true},
		{name: "zero quality with decimals", item: "application/json;q=0.000", expected: true},
		{name: "nonzero quality", item: "application/json;q=0.8", expected: false},
		{name: "full quality", item: "application/json;q=1", expected: false},
		{name: "unparsable", item: "not-a-media-type", expected: false},
		{name: "unparsable q value", item: "application/json;q=abc", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.expected, http.IsAcceptZeroQuality(tt.item))
		})
	}
}

func TestIsAcceptWildcard(t *testing.T) {
	t.Parallel()

	target := media.MustParse(media.NDJSON)

	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{name: "bare wildcard", value: "*/*", expected: true},
		{name: "matching major type wildcard", value: "application/*", expected: true},
		{name: "mismatched major type wildcard", value: "text/*", expected: false},
		{name: "concrete exact match", value: media.NDJSON, expected: false},
		{name: "concrete unrelated type", value: media.JSON, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			value, err := media.Parse(tt.value)
			require.NoError(t, err)

			require.Equal(t, tt.expected, http.IsAcceptWildcard(value, target))
		})
	}
}
