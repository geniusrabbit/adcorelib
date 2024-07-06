package formataccessor

import (
	"context"

	"github.com/geniusrabbit/adcorelib/admodels"
	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/models"
	"github.com/geniusrabbit/adcorelib/storage/loader"
)

var ctxFormatAccessor = struct{ s string }{s: "formataccessor"}

// WithContext returns new context with format accessor
func WithContext(ctx context.Context, accessor types.FormatsAccessor) context.Context {
	return context.WithValue(ctx, ctxFormatAccessor, accessor)
}

// FromContext returns format accessor
func FromContext(ctx context.Context) types.FormatsAccessor {
	return ctx.Value(ctxFormatAccessor).(types.FormatsAccessor)
}

// NewFormatAccessor from format dataAccessor
func NewFormatAccessor(formatDataAccessor loader.DataAccessor) types.FormatsAccessor {
	return types.NewSimpleFormatAccessorWithLoader(func() ([]*types.Format, error) {
		data, err := formatDataAccessor.Data()
		if err != nil {
			return nil, err
		}
		formats := make([]*types.Format, 0, len(data))
		for _, it := range data {
			formats = append(formats, admodels.FormatFromModel(it.(*models.Format)))
		}
		return formats, nil
	})
}
