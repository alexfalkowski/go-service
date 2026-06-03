package meta_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/stretchr/testify/require"
)

func TestStrings(t *testing.T) {
	ctx := meta.WithAttributes(t.Context(),
		meta.NewPair("testId", meta.String("1")),
		meta.NewPair("see", meta.Ignored("secret")),
		meta.NewPair("redacted", meta.Redacted("2")),
	)

	assertStrings(t, ctx, "snake", meta.Map{"test_id": "1", "redacted": "*"}, func(ctx context.Context) meta.Map {
		return meta.SnakeStrings(ctx, meta.NoPrefix)
	})
	assertStrings(t, ctx, "camel", meta.Map{"testId": "1", "redacted": "*"}, func(ctx context.Context) meta.Map {
		return meta.CamelStrings(ctx, meta.NoPrefix)
	})
	assertStrings(t, ctx, "none", meta.Map{"testId": "1", "redacted": "*"}, func(ctx context.Context) meta.Map {
		return meta.Strings(ctx, meta.NoPrefix)
	})
	assertStrings(t, ctx, "prefix", meta.Map{"test.testId": "1", "test.redacted": "*"}, func(ctx context.Context) meta.Map {
		return meta.Strings(ctx, "test.")
	})
}

func TestWithAttributesReturnsSameContextWithoutPairs(t *testing.T) {
	ctx := t.Context()

	require.Same(t, ctx, meta.WithAttributes(ctx))
}

func TestWithAttributesDoesNotCollideWithStringContextKey(t *testing.T) {
	ctx := context.WithValue(t.Context(), context.Key("meta"), "bad")

	require.NotPanics(t, func() {
		require.True(t, meta.Attribute(ctx, meta.RequestIDKey).IsEmpty())
	})

	ctx = meta.WithAttributes(ctx, meta.WithRequestID(meta.String("request-id")))

	require.Equal(t, "bad", ctx.Value(context.Key("meta")))
	require.Equal(t, meta.String("request-id"), meta.Attribute(ctx, meta.RequestIDKey))
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
		{name: "transport", key: meta.TransportKey, value: "transport", pair: meta.WithTransport},
		{name: "service method", key: meta.ServiceMethodKey, value: "service-method", pair: meta.WithServiceMethod},
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

func TestWithAttributesKeepsParentContextIsolatedWithSinglePair(t *testing.T) {
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

func TestAccessors(t *testing.T) {
	ctx := meta.WithAttributes(t.Context(),
		meta.WithRequestID(meta.String("request-id")),
		meta.WithTransport(meta.String("transport")),
		meta.WithServiceMethod(meta.String("service-method")),
		meta.WithUserAgent(meta.String("user-agent")),
		meta.WithUserID(meta.String("user-id")),
		meta.WithIPAddr(meta.String("ip-addr")),
		meta.WithAuthorization(meta.Ignored("authorization")),
		meta.WithGeolocation(meta.String("geo:47,11")),
	)

	tests := []struct {
		name string
		got  meta.Value
		want meta.Value
	}{
		{name: "request id", got: meta.RequestID(ctx), want: meta.String("request-id")},
		{name: "transport", got: meta.Transport(ctx), want: meta.String("transport")},
		{name: "service method", got: meta.ServiceMethod(ctx), want: meta.String("service-method")},
		{name: "transport service method", got: meta.TransportServiceMethod(ctx), want: meta.Ignored("transport:service-method")},
		{name: "user agent", got: meta.UserAgent(ctx), want: meta.String("user-agent")},
		{name: "user id", got: meta.UserID(ctx), want: meta.String("user-id")},
		{name: "ip addr", got: meta.IPAddr(ctx), want: meta.String("ip-addr")},
		{name: "authorization", got: meta.Authorization(ctx), want: meta.Ignored("authorization")},
		{name: "geolocation", got: meta.Geolocation(ctx), want: meta.String("geo:47,11")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.want, test.got)
		})
	}
}

func TestTransportServiceMethod(t *testing.T) {
	ctx := meta.WithAttributes(t.Context(),
		meta.WithTransport(meta.Ignored("http")),
		meta.WithServiceMethod(meta.Ignored("GET /users/{id}")),
	)

	require.Equal(t, meta.Ignored("http:GET /users/{id}"), meta.TransportServiceMethod(ctx))
}

func assertStrings(t *testing.T, ctx context.Context, name string, want meta.Map, export func(context.Context) meta.Map) {
	t.Helper()

	t.Run(name, func(t *testing.T) {
		require.Equal(t, want, export(ctx))
	})
}
