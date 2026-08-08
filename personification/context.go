package personification

import "context"

var (
	ContextKeySignature = struct{ s string }{s: "personification.signature"}
)

func WithSignature(ctx context.Context, signature *Signature) context.Context {
	return context.WithValue(ctx, ContextKeySignature, signature)
}

func SignatureFromContext(ctx context.Context) *Signature {
	return ctx.Value(ContextKeySignature).(*Signature)
}
