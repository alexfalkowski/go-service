package http_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/stretchr/testify/require"
)

func TestAddVary(t *testing.T) {
	tests := []struct {
		name     string
		existing []string
		fields   []string
		want     []string
	}{
		{
			name:   "adds fields",
			fields: []string{http.AcceptKey, http.ContentTypeKey},
			want:   []string{http.AcceptKey, http.ContentTypeKey},
		},
		{
			name:   "ignores blank fields",
			fields: []string{"", " \t ", http.AcceptKey},
			want:   []string{http.AcceptKey},
		},
		{
			name:     "preserves comma separated fields",
			existing: []string{"Accept-Encoding, Accept"},
			fields:   []string{http.AcceptKey, http.ContentTypeKey},
			want:     []string{"Accept-Encoding, Accept", http.ContentTypeKey},
		},
		{
			name:     "compares fields case insensitively",
			existing: []string{"accept, CONTENT-TYPE"},
			fields:   []string{http.AcceptKey, http.ContentTypeKey},
			want:     []string{"accept, CONTENT-TYPE"},
		},
		{
			name:     "preserves wildcard",
			existing: []string{"*"},
			fields:   []string{http.AcceptKey, http.ContentTypeKey},
			want:     []string{"*"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header := http.Header{http.VaryKey: tt.existing}
			http.AddVary(header, tt.fields...)

			require.Equal(t, tt.want, header.Values(http.VaryKey))
		})
	}
}
