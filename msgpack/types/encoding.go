package types

import "io"

type Encoder interface {
	Encode(val interface{}) error
}

type Decoder interface {
	Decode(val interface{}) error
}

type EncodeGenerator interface {
	NewEncoder(w io.Writer) Encoder
}

type DecodeGenerator interface {
	NewDecoder(reader io.Reader, buf []byte) Decoder
}
