//
// @project GeniusRabbit corelib 2018 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018 - 2019
//

package events

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"net/url"
	"strings"

	lz4 "github.com/bkaradzic/go-lz4"

	"github.com/geniusrabbit/adcorelib/msgpack"
	"github.com/geniusrabbit/adcorelib/msgpack/types"
)

// EventUnpacFunc function
type EventUnpacFunc func(code Code) Code

// LZ4BlockMaxSize constant
const LZ4BlockMaxSize = 64 << 10

// ErrEmptyData error
var ErrEmptyData = errors.New(`data is empty`)

// Code structure can contains and conver data for URL
type Code struct {
	data []byte
	err  error
}

// CodeObj by values
func CodeObj(data []byte, err error) Code {
	return Code{data: data, err: err}
}

// ObjectCode converts object to code object with defined encoder/decoder
func ObjectCode(obj any, gen ...types.EncodeGenerator) Code {
	var (
		buff bytes.Buffer
		enc  types.Encoder
	)

	if len(gen) > 0 && gen[0] != nil {
		enc = gen[0].NewEncoder(&buff)
	} else {
		enc = msgpack.DefaultEncodeGenerator.NewEncoder(&buff)
	}

	if err := enc.Encode(obj); err != nil {
		return CodeObj(nil, err)
	}
	return CodeObj(buff.Bytes(), nil)
}

func (c Code) String() string {
	return string(c.data)
}

// Data value
func (c Code) Data() []byte {
	return c.data
}

// Compress code.data
func (c Code) Compress() Code {
	if c.err != nil {
		return c
	}
	return CodeObj(lz4.Encode(nil, c.data))
}

// Decompress current code.data
func (c Code) Decompress() Code {
	if c.err != nil {
		return c
	}
	return CodeObj(lz4.Decode(nil, c.data))
}

func (c Code) Error() string {
	if c.err == nil {
		return ""
	}
	return c.err.Error()
}

// ErrorObj response
func (c Code) ErrorObj() error {
	return c.err
}

// URLEncode data for URL using. RawURLEncoding has no '=' padding, so the
// token stays inside [A-Za-z0-9-_] and needs no query escaping. A padded
// token becomes %3D, and an SSP that encodes the nurl again turns that into
// %253D, which URLEncoding.DecodeString rejects.
func (c Code) URLEncode() Code {
	if len(c.data) < 1 {
		return CodeObj(nil, ErrEmptyData)
	}
	return CodeObj([]byte(base64.RawURLEncoding.EncodeToString(c.data)), nil)
}

// URLDecode value. Accepts the raw token, a legacy padded token, a token
// whose '=' survived as %3D, and a token with trailing garbage (for example
// a JSON \u0026 that an SSP glued onto the query value).
func (c Code) URLDecode() Code {
	s := string(c.data)
	if strings.Contains(s, "%") {
		if unescaped, err := url.PathUnescape(s); err == nil {
			s = unescaped
		}
	}
	end := 0
	for end < len(s) && isRawURLBase64Byte(s[end]) {
		end++
	}
	s = strings.TrimRight(s[:end], "=")
	data, err := base64.RawURLEncoding.DecodeString(s)
	return CodeObj(data, err)
}

func isRawURLBase64Byte(c byte) bool {
	switch {
	case c >= 'A' && c <= 'Z', c >= 'a' && c <= 'z', c >= '0' && c <= '9':
		return true
	case c == '-', c == '_', c == '=':
		return true
	default:
		return false
	}
}

// DecodeObject converts current object data to target
func (c Code) DecodeObject(target any, gen ...types.DecodeGenerator) error {
	if c.err != nil {
		return c.err
	}

	var dec types.Decoder
	if len(gen) > 0 && gen[0] != nil {
		dec = gen[0].NewDecoder(nil, c.data)
	} else {
		dec = msgpack.DefaultDecodeGenerator.NewDecoder(nil, c.data)
	}
	return dec.Decode(target)
}

// ResetData object
func (c *Code) ResetData() {
	if len(c.data) > 0 {
		c.data = c.data[:]
	}
}

// Write method is implementation of io.Writer any
func (c *Code) Write(p []byte) (n int, err error) {
	c.data = append(c.data, p...)
	return len(p), nil
}

var _ io.Writer = (*Code)(nil)
