package personification

import (
	"context"

	"github.com/geniusrabbit/adcorelib/personification/dummy"
)

var (
	dummyClient         dummy.DummyClient
	ContextKeySignature = struct{ s string }{s: "personification.signature"}
)

func WithSignature(ctx context.Context, signature *Signature) context.Context {
	return context.WithValue(ctx, ContextKeySignature, signature)
}

func SignatureFromContext(ctx context.Context) *Signature {
	if sign, ok := ctx.Value(ContextKeySignature).(*Signature); ok {
		return sign
	}
	return nil
}

func DetectorFromContext(ctx context.Context) Client {
	if sign := SignatureFromContext(ctx); sign != nil {
		return sign.Detector
	}
	return &dummyClient
}
