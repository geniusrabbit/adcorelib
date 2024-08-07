// Package msgpack provides premetive methods and object for message pack and unpacking
package msgpack

import (
	"bytes"

	"github.com/pierrec/lz4"

	jsongen "github.com/geniusrabbit/adcorelib/msgpack/json"
)

var (
	DefaultEncodeGenerator = jsongen.EncodeGenerator{}
	DefaultDecodeGenerator = jsongen.DecodeGenerator{}
)

// StdPack message
func StdPack(msg any) ([]byte, error) {
	var (
		buff   bytes.Buffer
		writer = lz4.NewWriter(&buff)
		err    = DefaultEncodeGenerator.NewEncoder(writer).Encode(msg)
	)
	if err != nil {
		return nil, err
	}
	if err = writer.Flush(); err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

// StdUnpack message
func StdUnpack(data []byte, msg any) error {
	return DefaultDecodeGenerator.NewDecoder(lz4.NewReader(bytes.NewReader(data)), nil).Decode(msg)
}
