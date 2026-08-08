package otlp

import "github.com/alexfalkowski/go-service/v2/errors"

const (
	// DefaultMaxQueueSize is the OpenTelemetry SDK default queue size.
	DefaultMaxQueueSize = 2048
	// DefaultMaxExportBatchSize is the OpenTelemetry SDK default export batch size.
	DefaultMaxExportBatchSize = 512
	// MaxQueueSize is the largest supported OTLP record queue.
	MaxQueueSize = 8192
	// MaxExportBatchSize is the largest supported OTLP export batch.
	MaxExportBatchSize = 2048
)

// ErrInvalidBatchConfig is returned when an OTLP batch queue or export batch
// exceeds its supported limit or the effective queue size.
var ErrInvalidBatchConfig = errors.New("otlp: invalid batch configuration")

// BatchConfig describes the OTLP queue and export batch limits.
type BatchConfig struct {
	// MaxQueueSize bounds queued records before they are dropped.
	MaxQueueSize int
	// MaxExportBatchSize bounds records exported together.
	MaxExportBatchSize int
}

// ValidateBatchConfig accepts zero SDK-default selectors, supported count
// limits, and an effective batch no larger than its effective queue.
func ValidateBatchConfig(config BatchConfig) error {
	if config.MaxQueueSize < 0 || config.MaxQueueSize > MaxQueueSize ||
		config.MaxExportBatchSize < 0 || config.MaxExportBatchSize > MaxExportBatchSize {
		return ErrInvalidBatchConfig
	}

	queueSize := config.MaxQueueSize
	if queueSize == 0 {
		queueSize = DefaultMaxQueueSize
	}

	batchSize := config.MaxExportBatchSize
	if batchSize == 0 {
		batchSize = DefaultMaxExportBatchSize
	}
	if batchSize > queueSize {
		return ErrInvalidBatchConfig
	}

	return nil
}
