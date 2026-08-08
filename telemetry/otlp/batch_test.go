package otlp_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/telemetry/otlp"
	"github.com/alexfalkowski/go-service/v2/time"
	"github.com/stretchr/testify/require"
)

func TestValidateBatchConfig(t *testing.T) {
	tests := []struct {
		config  otlp.BatchConfig
		name    string
		wantErr bool
	}{
		{name: "defaults"},
		{name: "limits", config: otlp.BatchConfig{MaxQueueSize: otlp.MaxQueueSize, MaxExportBatchSize: otlp.MaxExportBatchSize}},
		{name: "queue above limit", config: otlp.BatchConfig{MaxQueueSize: otlp.MaxQueueSize + 1}, wantErr: true},
		{name: "batch above limit", config: otlp.BatchConfig{MaxExportBatchSize: otlp.MaxExportBatchSize + 1}, wantErr: true},
		{name: "configured batch above queue", config: otlp.BatchConfig{MaxQueueSize: 1, MaxExportBatchSize: 2}, wantErr: true},
		{name: "default batch above queue", config: otlp.BatchConfig{MaxQueueSize: 1}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := otlp.ValidateBatchConfig(tt.config)
			if tt.wantErr {
				require.ErrorIs(t, err, otlp.ErrInvalidBatchConfig)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestValidateCadence(t *testing.T) {
	tests := []struct {
		cadence time.Duration
		name    string
		wantErr bool
	}{
		{name: "default"},
		{name: "whole second", cadence: time.Second},
		{name: "fractional second", cadence: 1500 * time.Millisecond, wantErr: true},
		{name: "negative", cadence: -time.Second, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := otlp.ValidateCadence(tt.cadence)
			if tt.wantErr {
				require.ErrorIs(t, err, otlp.ErrInvalidCadence)
				return
			}

			require.NoError(t, err)
		})
	}
}
