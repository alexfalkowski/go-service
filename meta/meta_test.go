package meta_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/stretchr/testify/require"
)

func TestSnakeCase(t *testing.T) {
	ctx := t.Context()
	ctx = meta.WithAttribute(ctx, "testId", meta.String("1"))
	ctx = meta.WithAttribute(ctx, "see", meta.Ignored("secret"))
	ctx = meta.WithAttribute(ctx, "redacted", meta.Redacted("2"))

	require.Equal(t, meta.Map{"test_id": "1", "redacted": "*"}, meta.SnakeStrings(ctx, meta.NoPrefix))
}

func TestCamelCase(t *testing.T) {
	ctx := t.Context()
	ctx = meta.WithAttribute(ctx, "testId", meta.String("1"))
	ctx = meta.WithAttribute(ctx, "see", meta.Ignored("secret"))
	ctx = meta.WithAttribute(ctx, "redacted", meta.Redacted("2"))

	require.Equal(t, meta.Map{"testId": "1", "redacted": "*"}, meta.CamelStrings(ctx, meta.NoPrefix))
}

func TestNoneCase(t *testing.T) {
	ctx := t.Context()
	ctx = meta.WithAttribute(ctx, "testId", meta.String("1"))
	ctx = meta.WithAttribute(ctx, "see", meta.Ignored("secret"))
	ctx = meta.WithAttribute(ctx, "redacted", meta.Redacted("2"))

	require.Equal(t, meta.Map{"testId": "1", "redacted": "*"}, meta.Strings(ctx, meta.NoPrefix))
}

func TestPrefix(t *testing.T) {
	ctx := t.Context()
	ctx = meta.WithAttribute(ctx, "testId", meta.String("1"))
	ctx = meta.WithAttribute(ctx, "see", meta.Ignored("secret"))
	ctx = meta.WithAttribute(ctx, "redacted", meta.Redacted("2"))

	require.Equal(t, meta.Map{"test.testId": "1", "test.redacted": "*"}, meta.Strings(ctx, "test."))
}

func TestUserID(t *testing.T) {
	ctx := meta.WithUserID(t.Context(), meta.String("user-id"))
	require.Equal(t, meta.String("user-id"), meta.UserID(ctx))
}

func TestGeolocation(t *testing.T) {
	ctx := meta.WithGeolocation(t.Context(), meta.String("geo:47,11"))
	require.Equal(t, meta.String("geo:47,11"), meta.Geolocation(ctx))
}
