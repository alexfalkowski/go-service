package logger_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/transport/grpc/telemetry/logger"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
)

func TestLogger(t *testing.T) {
	require.Equal(t, logger.LevelError, logger.CodeToLevel(codes.DeadlineExceeded))
}
