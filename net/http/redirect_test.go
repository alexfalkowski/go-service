package http_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/url"
	"github.com/stretchr/testify/require"
)

func TestSameOriginRedirectMatchesSchemeHostAndPort(t *testing.T) {
	t.Parallel()

	tests := []struct {
		want error
		name string
		next string
	}{
		{name: "same origin", next: "https://example.com/next", want: nil},
		{name: "same origin host case", next: "https://EXAMPLE.com/next", want: nil},
		{name: "same origin default port", next: "https://example.com:443/next", want: nil},
		{name: "different host", next: "https://other.example.com/next", want: http.ErrUseLastResponse},
		{name: "different scheme", next: "http://example.com/next", want: http.ErrUseLastResponse},
		{name: "different port", next: "https://example.com:444/next", want: http.ErrUseLastResponse},
	}

	prev, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "https://example.com/start", http.NoBody)
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			next, err := http.NewRequestWithContext(t.Context(), http.MethodGet, tt.next, http.NoBody)
			require.NoError(t, err)
			err = http.SameOriginRedirect(next, []*http.Request{prev})
			if tt.want == nil {
				require.NoError(t, err)
				return
			}

			require.ErrorIs(t, err, tt.want)
		})
	}
}

func TestSameOriginRedirectStopsAfterTenRedirects(t *testing.T) {
	t.Parallel()

	via := make([]*http.Request, 10)
	for i := range via {
		request, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "https://example.com/redirect", http.NoBody)
		require.NoError(t, err)

		via[i] = request
	}

	next, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "https://example.com/redirect", http.NoBody)
	require.NoError(t, err)

	require.NoError(t, http.SameOriginRedirect(next, via[:9]))

	err = http.SameOriginRedirect(next, via)
	require.Error(t, err)
	require.NotErrorIs(t, err, http.ErrUseLastResponse)

	crossOrigin, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "https://other.example.com/redirect", http.NoBody)
	require.NoError(t, err)
	require.ErrorIs(t, http.SameOriginRedirect(crossOrigin, via), http.ErrUseLastResponse)
}

func TestSameOriginComparesSchemeHostAndPort(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		prev string
		next string
		want bool
	}{
		{name: "same origin", prev: "https://example.com/start", next: "https://example.com/next", want: true},
		{name: "same origin host case", prev: "https://example.com/start", next: "https://EXAMPLE.com/next", want: true},
		{name: "same origin https default port", prev: "https://example.com/start", next: "https://example.com:443/next", want: true},
		{name: "same origin http default port", prev: "http://example.com/start", next: "http://example.com:80/next", want: true},
		{name: "different host", prev: "https://example.com/start", next: "https://other.example.com/next", want: false},
		{name: "different scheme", prev: "https://example.com/start", next: "http://example.com/next", want: false},
		{name: "different port", prev: "https://example.com/start", next: "https://example.com:444/next", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			prev, err := url.Parse(tt.prev)
			require.NoError(t, err)

			next, err := url.Parse(tt.next)
			require.NoError(t, err)

			require.Equal(t, tt.want, http.SameOrigin(prev, next))
		})
	}

	prev, err := url.Parse("https://example.com/start")
	require.NoError(t, err)

	next, err := url.Parse("https://example.com/next")
	require.NoError(t, err)

	require.False(t, http.SameOrigin(nil, next))
	require.False(t, http.SameOrigin(prev, nil))
}

func TestIsCrossOriginRedirect(t *testing.T) {
	t.Parallel()

	prev, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "https://example.com/start", http.NoBody)
	require.NoError(t, err)

	same, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "https://example.com/next", http.NoBody)
	require.NoError(t, err)
	same.Response = &http.Response{Request: prev}

	different, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "https://other.example.com/next", http.NoBody)
	require.NoError(t, err)
	different.Response = &http.Response{Request: prev}

	require.False(t, http.IsCrossOriginRedirect(nil))
	require.False(t, http.IsCrossOriginRedirect(prev))
	require.False(t, http.IsCrossOriginRedirect(same))
	require.True(t, http.IsCrossOriginRedirect(different))
}
