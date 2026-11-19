package events

import (
	"bytes"
	"encoding/base64"
	"strings"
	"testing"
)

func TestCodeCompression(t *testing.T) {
	dataTest := []byte("Test")
	code1 := CodeObj(dataTest, nil)
	code2 := code1.Compress()
	if !bytes.Equal(code2.Decompress().Data(), dataTest) {
		t.Errorf("Invalid compression: %s", code2.Error())
	}
}

func TestURLEncodeRoundTripOmitsPadding(t *testing.T) {
	payloads := [][]byte{
		[]byte("a"),
		[]byte("ab"),
		[]byte("abc"),
		[]byte("hello world!!"),
		bytes.Repeat([]byte{1, 2, 3, 4}, 50),
	}
	for _, payload := range payloads {
		enc := CodeObj(payload, nil).URLEncode()
		if enc.ErrorObj() != nil {
			t.Fatalf("encode: %v", enc.ErrorObj())
		}
		if strings.ContainsAny(enc.String(), "=%+/") {
			t.Fatalf("encoded %q is not a raw url token", enc.String())
		}
		dec := enc.URLDecode()
		if dec.ErrorObj() != nil {
			t.Fatalf("decode %q: %v", enc.String(), dec.ErrorObj())
		}
		if !bytes.Equal(dec.Data(), payload) {
			t.Fatalf("round trip got %q want %q", dec.Data(), payload)
		}
	}
}

func TestURLDecodeAcceptsLegacyTokenShapes(t *testing.T) {
	payload := []byte("hello world!!")
	padded := base64.URLEncoding.EncodeToString(payload)
	if !strings.Contains(padded, "=") {
		t.Fatal("fixture should include padding")
	}
	stripped := strings.TrimRight(padded, "=")

	cases := []string{
		padded,
		strings.ReplaceAll(padded, "=", "%3D"),
		stripped,
		padded + `\u0026prc=1.23`,
		strings.ReplaceAll(padded, "=", "%3D") + `\u0026prc=1.23`,
	}
	for _, token := range cases {
		dec := CodeObj([]byte(token), nil).URLDecode()
		if dec.ErrorObj() != nil {
			t.Fatalf("decode %q: %v", token, dec.ErrorObj())
		}
		if !bytes.Equal(dec.Data(), payload) {
			t.Fatalf("decode %q got %q want %q", token, dec.Data(), payload)
		}
	}
}
