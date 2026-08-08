package preprocessors

import "github.com/geniusrabbit/adcorelib/adtype"

// CustomFunc is a function that preprocesses a response
type CustomFunc func(response adtype.Response) (adtype.Response, error)

// PreprocessResponse implements adsource.ResponsePreprocessor
func (p CustomFunc) PreprocessResponse(response adtype.Response) (adtype.Response, error) {
	return p(response)
}
