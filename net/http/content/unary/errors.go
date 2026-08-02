package unary

import "github.com/alexfalkowski/go-service/v2/errors"

// ErrUnsupportedRequestMedia is returned when a request Content-Type cannot be decoded from a public
// HTTP request body. This covers three causes: the media type is intentionally denied (see the
// decoder-bounds rule in the package documentation), the media type does not parse, or the media type
// parses but resolves to no registered codec. See [Content.NewFromRequestBody].
var ErrUnsupportedRequestMedia = errors.New("unary: unsupported request media")
