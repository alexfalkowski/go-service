package metrics

import (
	"unique"
)

var (
	unaryKind  = unique.Make("unary")
	streamKind = unique.Make("stream")
)
