package json

import (
	"bytes"
	"encoding/json"
	"io"

	"geniusrabbit.dev/adcorelib/msgpack/types"
)

type EncodeGenerator struct{}

func (g *EncodeGenerator) NewEncoder(w io.Writer) types.Encoder {
	return json.NewEncoder(w)
}

type DecodeGenerator struct{}

func (g *DecodeGenerator) NewDecoder(reader io.Reader, buf []byte) types.Decoder {
	if reader == nil {
		reader = bytes.NewReader(buf)
	}
	return json.NewDecoder(reader)
}
