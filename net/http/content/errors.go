package content

import "github.com/alexfalkowski/go-service/v2/errors"

// ErrUnsupportedRequestMedia is returned when a request body uses a media type
// that is intentionally not decoded from public HTTP requests.
var ErrUnsupportedRequestMedia = errors.New("content: unsupported request media")
