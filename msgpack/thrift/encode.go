package thrift

// import (
// 	"io"

// 	basethrift "github.com/thrift-iterator/go"

// 	"geniusrabbit.dev/adcorelib/msgpack/types"
// )

// type EncodeGenerator struct {
// 	api basethrift.API
// }

// func (g *EncodeGenerator) NewEncoder(w io.Writer) types.Encoder {
// 	if g.api != nil {
// 		return g.api.NewEncoder(w)
// 	}
// 	return NewEncoder(w)
// }

// type DecodeGenerator struct {
// 	api basethrift.API
// }

// func (g *DecodeGenerator) NewDecoder(reader io.Reader, buf []byte) types.Decoder {
// 	if g.api != nil {
// 		return g.api.NewDecoder(reader, buf)
// 	}
// 	return NewDecoder(reader, buf)
// }
