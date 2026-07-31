package encoding_test

import (
	"testing"

	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/encoding/bytes"
	"github.com/alexfalkowski/go-service/v2/encoding/gob"
	"github.com/alexfalkowski/go-service/v2/encoding/hjson"
	"github.com/alexfalkowski/go-service/v2/encoding/json"
	"github.com/alexfalkowski/go-service/v2/encoding/msgpack"
	"github.com/alexfalkowski/go-service/v2/encoding/proto"
	"github.com/alexfalkowski/go-service/v2/encoding/toml"
	"github.com/alexfalkowski/go-service/v2/encoding/yaml"
	"github.com/alexfalkowski/go-service/v2/internal/test"
	"github.com/stretchr/testify/require"
)

func TestNewMapRegistersDefaultEncoders(t *testing.T) {
	encoders := encoding.NewMap()

	expected := map[string]encoding.Encoder{
		"json":      json.NewEncoder(),
		"hjson":     hjson.NewEncoder(),
		"yaml":      yaml.NewEncoder(),
		"toml":      toml.NewEncoder(),
		"msgpack":   msgpack.NewEncoder(),
		"protobuf":  proto.NewBinary(),
		"prototext": proto.NewText(),
		"protojson": proto.NewJSON(),
		"gob":       gob.NewEncoder(),
		"bytes":     bytes.NewEncoder(),
	}

	for kind, expectedEncoder := range expected {
		t.Run(kind, func(t *testing.T) {
			require.IsType(t, expectedEncoder, encoders.Get(kind))
		})
	}
}

func TestMapRegister(t *testing.T) {
	encoders := encoding.NewMap()
	custom := test.NewEncoder(test.ErrFailed)
	replacement := bytes.NewEncoder()

	encoders.Register("custom", custom)
	require.Same(t, custom, encoders.Get("custom"))

	encoders.Register("custom", replacement)
	require.Same(t, replacement, encoders.Get("custom"))
}
