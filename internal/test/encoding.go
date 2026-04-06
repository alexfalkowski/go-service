package test

import (
	"io"

	"github.com/alexfalkowski/go-service/v2/encoding"
	"github.com/alexfalkowski/go-service/v2/encoding/bytes"
	"github.com/alexfalkowski/go-service/v2/encoding/gob"
	"github.com/alexfalkowski/go-service/v2/encoding/hjson"
	"github.com/alexfalkowski/go-service/v2/encoding/json"
	"github.com/alexfalkowski/go-service/v2/encoding/proto"
	"github.com/alexfalkowski/go-service/v2/encoding/toml"
	"github.com/alexfalkowski/go-service/v2/encoding/yaml"
	"github.com/alexfalkowski/go-service/v2/net/http/content"
)

// Encoder contains the real encoders exercised by config, cache, and transport tests.
var Encoder = encoding.NewMap(encoding.MapParams{
	JSON:        json.NewEncoder(),
	HJSON:       hjson.NewEncoder(),
	YAML:        yaml.NewEncoder(),
	TOML:        toml.NewEncoder(),
	ProtoBinary: proto.NewBinary(),
	ProtoText:   proto.NewText(),
	ProtoJSON:   proto.NewJSON(),
	GOB:         gob.NewEncoder(),
	Bytes:       bytes.NewEncoder(),
})

// Content is the shared HTTP content registry backed by Encoder.
var Content = content.NewContent(Encoder, Pool)

// NewEncoder returns an encoder test double whose Encode and Decode methods fail with the supplied error.
func NewEncoder(err error) encoding.Encoder {
	return &enc{err: err}
}

type enc struct {
	err error
}

// Encode implements encoding.Encoder and returns the configured error.
func (e *enc) Encode(_ io.Writer, _ any) error {
	return e.err
}

// Decode implements encoding.Encoder and returns the configured error.
func (e *enc) Decode(_ io.Reader, _ any) error {
	return e.err
}
