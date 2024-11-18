package adresponse

type responseDataAccessor interface {
	Get(key string) any
}

type defaultData map[string]any

func (d defaultData) Get(key string) any {
	return d[key]
}
