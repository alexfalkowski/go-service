package grpc

import (
	"github.com/alexfalkowski/go-service/v2/os"
)

var fs *os.FS

// Register for grpc.
func Register(f *os.FS) {
	fs = f
}
