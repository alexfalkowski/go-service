package meta_test

import (
	"fmt"
	"testing"

	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/meta"
	"github.com/stretchr/testify/require"
)

func TestErrorWithNil(t *testing.T) {
	require.Equal(t, meta.Blank(), meta.Error(nil))
}

func TestErrorWithTypedNil(t *testing.T) {
	var err error = (*test.MessageError)(nil)

	require.Equal(t, meta.Blank(), meta.Error(err))
}

func TestToStringWithNil(t *testing.T) {
	require.Equal(t, meta.Blank(), meta.ToString(nil))
}

func TestToStringWithTypedNil(t *testing.T) {
	var stringer fmt.Stringer = (*test.Stringer)(nil)

	require.Equal(t, meta.Blank(), meta.ToString(stringer))
}

func TestToRedactedWithTypedNil(t *testing.T) {
	var stringer fmt.Stringer = (*test.Stringer)(nil)

	require.Equal(t, meta.Blank(), meta.ToRedacted(stringer))
}

func TestToIgnoredWithTypedNil(t *testing.T) {
	var stringer fmt.Stringer = (*test.Stringer)(nil)

	require.Equal(t, meta.Blank(), meta.ToIgnored(stringer))
}

func TestRedactedWithMultiByteValue(t *testing.T) {
	require.Equal(t, "*", meta.Redacted("é").String())
	require.Equal(t, "**", meta.Redacted("éa").String())
}
