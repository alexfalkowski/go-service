package config_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
)

type configSize struct {
	Size bytes.Size `validate:"config_size"`
}

type secondPrecisionDuration struct {
	Duration time.Duration `validate:"duration_second_precision"`
}

func TestValidatorConfigSize(t *testing.T) {
	tests := []struct {
		name  string
		size  bytes.Size
		valid bool
	}{
		{name: "negative", size: -1},
		{name: "zero", valid: true},
		{name: "default", size: bytes.DefaultSize, valid: true},
		{name: "max", size: bytes.MaxConfigSize, valid: true},
		{name: "oversized", size: bytes.MaxConfigSize + 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := test.Validator.Struct(&configSize{Size: tt.size})
			if tt.valid {
				require.NoError(t, err)
				return
			}

			require.Error(t, err)
		})
	}
}

func TestValidatorDurationSecondPrecision(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		valid    bool
	}{
		{name: "negative", duration: -time.Second},
		{name: "zero"},
		{name: "sub second", duration: 500 * time.Millisecond},
		{name: "second", duration: time.Second, valid: true},
		{name: "minute", duration: time.Minute, valid: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := test.Validator.Struct(&secondPrecisionDuration{Duration: tt.duration})
			if tt.valid {
				require.NoError(t, err)
				return
			}

			require.Error(t, err)
		})
	}
}
