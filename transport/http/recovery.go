package http

import (
	"github.com/alexfalkowski/go-service/v2/io"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/net/http/status"
	"github.com/alexfalkowski/go-service/v2/runtime"
	snoop "github.com/felixge/httpsnoop"
)

type recoveryHandler struct{}

func (*recoveryHandler) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	committed := false
	hooks := snoop.Hooks{
		WriteHeader: func(writeHeader snoop.WriteHeaderFunc) snoop.WriteHeaderFunc {
			return func(code int) {
				committed = true
				writeHeader(code)
			}
		},
		Write: func(write snoop.WriteFunc) snoop.WriteFunc {
			return func(bytes []byte) (int, error) {
				committed = true
				return write(bytes)
			}
		},
		ReadFrom: func(readFrom snoop.ReadFromFunc) snoop.ReadFromFunc {
			return func(reader io.Reader) (int64, error) {
				committed = true
				return readFrom(reader)
			}
		},
		WriteString: func(writeString snoop.WriteStringFunc) snoop.WriteStringFunc {
			return func(value string) (int, error) {
				committed = true
				return writeString(value)
			}
		},
		Flush: func(flush snoop.FlushFunc) snoop.FlushFunc {
			return func() {
				committed = true
				flush()
			}
		},
		FlushError: func(flushError snoop.FlushErrorFunc) snoop.FlushErrorFunc {
			return func() error {
				committed = true
				return flushError()
			}
		},
	}
	defer func() {
		if value := recover(); value != nil {
			if committed {
				// The response is already on the wire, so writing a safe 500 would mix
				// response states. Let net/http abort the broken in-flight response.
				panic(value)
			}

			err := status.SafeError(http.StatusInternalServerError, runtime.ConvertRecover(value))
			_ = status.WriteError(res, err)
		}
	}()

	next(snoop.Wrap(res, hooks), req)
}
