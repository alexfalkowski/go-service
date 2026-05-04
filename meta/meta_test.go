package meta_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/stretchr/testify/require"
)

func TestSnakeCase(t *testing.T) {
	ctx := t.Context()
	ctx = meta.WithAttributes(ctx,
		meta.NewPair("testId", meta.String("1")),
		meta.NewPair("see", meta.Ignored("secret")),
		meta.NewPair("redacted", meta.Redacted("2")),
	)

	require.Equal(t, meta.Map{"test_id": "1", "redacted": "*"}, meta.SnakeStrings(ctx, meta.NoPrefix))
}

func TestCamelCase(t *testing.T) {
	ctx := t.Context()
	ctx = meta.WithAttributes(ctx,
		meta.NewPair("testId", meta.String("1")),
		meta.NewPair("see", meta.Ignored("secret")),
		meta.NewPair("redacted", meta.Redacted("2")),
	)

	require.Equal(t, meta.Map{"testId": "1", "redacted": "*"}, meta.CamelStrings(ctx, meta.NoPrefix))
}

func TestNoneCase(t *testing.T) {
	ctx := t.Context()
	ctx = meta.WithAttributes(ctx,
		meta.NewPair("testId", meta.String("1")),
		meta.NewPair("see", meta.Ignored("secret")),
		meta.NewPair("redacted", meta.Redacted("2")),
	)

	require.Equal(t, meta.Map{"testId": "1", "redacted": "*"}, meta.Strings(ctx, meta.NoPrefix))
}

func TestPrefix(t *testing.T) {
	ctx := t.Context()
	ctx = meta.WithAttributes(ctx,
		meta.NewPair("testId", meta.String("1")),
		meta.NewPair("see", meta.Ignored("secret")),
		meta.NewPair("redacted", meta.Redacted("2")),
	)

	require.Equal(t, meta.Map{"test.testId": "1", "test.redacted": "*"}, meta.Strings(ctx, "test."))
}

func TestWithAttributesReturnsSameContextWithoutPairs(t *testing.T) {
	ctx := t.Context()

	require.Same(t, ctx, meta.WithAttributes(ctx))
}

func TestPairHelpers(t *testing.T) {
	for _, test := range []struct {
		pair  func(meta.Value) meta.Pair
		name  string
		key   string
		value string
	}{
		{name: "request id", key: meta.RequestIDKey, value: "request-id", pair: meta.WithRequestID},
		{name: "system", key: meta.SystemKey, value: "system", pair: meta.WithSystem},
		{name: "service", key: meta.ServiceKey, value: "service", pair: meta.WithService},
		{name: "method", key: meta.MethodKey, value: "method", pair: meta.WithMethod},
		{name: "code", key: meta.CodeKey, value: "code", pair: meta.WithCode},
		{name: "duration", key: meta.DurationKey, value: "duration", pair: meta.WithDuration},
		{name: "user agent", key: meta.UserAgentKey, value: "user-agent", pair: meta.WithUserAgent},
		{name: "user id", key: meta.UserIDKey, value: "user-id", pair: meta.WithUserID},
		{name: "ip addr", key: meta.IPAddrKey, value: "ip-addr", pair: meta.WithIPAddr},
		{name: "ip addr kind", key: meta.IPAddrKindKey, value: "ip-addr-kind", pair: meta.WithIPAddrKind},
		{name: "authorization", key: meta.AuthorizationKey, value: "authorization", pair: meta.WithAuthorization},
		{name: "geolocation", key: meta.GeolocationKey, value: "geolocation", pair: meta.WithGeolocation},
	} {
		t.Run(test.name, func(t *testing.T) {
			pair := test.pair(meta.String(test.value))

			require.Equal(t, test.key, pair.Key)
			require.Equal(t, meta.String(test.value), pair.Value)
		})
	}
}

func TestWithAttributesKeepsParentContextIsolatedWithSinglePairs(t *testing.T) {
	parent := meta.WithAttributes(t.Context(), meta.WithRequestID(meta.String("parent")))
	child := meta.WithAttributes(parent, meta.WithUserID(meta.String("child")))

	require.Equal(t, meta.String("parent"), meta.Attribute(parent, "requestId"))
	require.True(t, meta.Attribute(parent, "userId").IsEmpty())
	require.Equal(t, meta.String("parent"), meta.Attribute(child, "requestId"))
	require.Equal(t, meta.String("child"), meta.Attribute(child, "userId"))
}

func TestWithAttributesKeepsParentContextIsolated(t *testing.T) {
	parent := meta.WithAttributes(t.Context(),
		meta.WithRequestID(meta.String("parent")),
		meta.WithUserAgent(meta.String("test-agent")),
	)
	child := meta.WithAttributes(parent,
		meta.WithRequestID(meta.String("child")),
		meta.WithUserID(meta.String("user")),
	)

	require.Equal(t, meta.String("parent"), meta.Attribute(parent, "requestId"))
	require.Equal(t, meta.String("test-agent"), meta.Attribute(parent, "userAgent"))
	require.True(t, meta.Attribute(parent, "userId").IsEmpty())
	require.Equal(t, meta.String("child"), meta.Attribute(child, "requestId"))
	require.Equal(t, meta.String("test-agent"), meta.Attribute(child, "userAgent"))
	require.Equal(t, meta.String("user"), meta.Attribute(child, "userId"))
}

func TestUserID(t *testing.T) {
	ctx := meta.WithAttributes(t.Context(), meta.WithUserID(meta.String("user-id")))
	require.Equal(t, meta.String("user-id"), meta.UserID(ctx))
}

func TestGeolocation(t *testing.T) {
	ctx := meta.WithAttributes(t.Context(), meta.WithGeolocation(meta.String("geo:47,11")))
	require.Equal(t, meta.String("geo:47,11"), meta.Geolocation(ctx))
}
