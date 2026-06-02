package strings_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/strings"
	"github.com/stretchr/testify/require"
)

func TestCaseConversion(t *testing.T) {
	tests := []struct {
		convert func() string
		name    string
		want    string
	}{
		{name: "delimited", convert: func() string { return strings.ToDelimited("InvalidArgument", ' ') }, want: "invalid argument"},
		{name: "lower camel", convert: func() string { return strings.ToLowerCamel("request_id") }, want: "requestId"},
		{name: "snake", convert: func() string { return strings.ToSnake("requestID") }, want: "request_id"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.want, test.convert())
		})
	}
}
