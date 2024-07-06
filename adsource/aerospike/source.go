package aerospike

import (
	"log"

	as "github.com/aerospike/aerospike-client-go"
	"github.com/geniusrabbit/adcorelib/adtype"
)

type source struct {
	client    *as.Client
	udfPath   string
	namespace string
}

// NewSource accessor
func NewSource(opts ...OptionFnk) (*source, error) {
	var (
		opt = option{hostname: "127.0.0.1", port: 3000, namespace: "default", udfPath: "/udf/"}
		err error
	)

	for _, o := range opts {
		o(&opt)
	}

	src := &source{udfPath: opt.udfPath, namespace: opt.namespace}
	if src.client, err = as.NewClient(opt.hostname, opt.port); err != nil {
		return nil, err
	}
	if err = src.initUDF(); err != nil {
		return nil, err
	}
	return src, nil
}

// Bid request for standart system filter
func (src *source) Bid(request *adtype.BidRequest) adtype.Responser {
	args := []as.Value{}
	stm := as.NewStatement(src.namespace, "campaigns")
	res, err := src.client.QueryAggregate(nil, stm, "search_ads", "search_ads", args...)

	if err != nil {
		return adtype.NewErrorResponse(request, err)
	}

	for rec := range res.Results() {
		res := rec.Record.Bins["SUCCESS"].(map[any]any)
		log.Printf("Result from Map/Reduce: %v\n", res)
		log.Printf("Result %f\n", res["sum"].(float64)/res["count"].(float64))
	}

	return adtype.NewErrorResponse(request, err)
}

// ProcessResponseItem result or error
func (src *source) ProcessResponseItem(adtype.Responser, adtype.ResponserItem) {

}

func (src *source) initUDF() error {
	as.SetLuaPath(src.udfPath)

	files := []string{"search_ads", "search_strict_ads"}
	for _, filename := range files {
		regTask, err := src.client.RegisterUDFFromFile(nil, src.udfPath+filename+".lua", filename+".lua", as.LUA)
		if err == nil {
			err = <-regTask.OnComplete()
		}
		if err != nil {
			return err
		}
	}
	return nil
}
