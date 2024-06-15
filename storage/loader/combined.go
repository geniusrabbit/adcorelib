package loader

type MergeFnk func(datas ...[]interface{}) []interface{}

// CombinedLoader loads and merge data into other type
type CombinedLoader struct {
	loaders []DataAccessor
	merge   MergeFnk
}

// NewCombinedLoader returns combined implementation of dataloader
func NewCombinedLoader(merge MergeFnk, loaders ...DataAccessor) *CombinedLoader {
	return &CombinedLoader{
		loaders: loaders,
		merge:   merge,
	}
}

func (l *CombinedLoader) NeedUpdate() bool {
	for _, lr := range l.loaders {
		if lr.NeedUpdate() {
			return true
		}
	}
	return false
}

// Data returns loaded data and reload if necessary
func (l *CombinedLoader) Data() ([]interface{}, error) {
	datas := make([][]interface{}, 0, len(l.loaders))
	for _, lr := range l.loaders {
		dr, err := lr.Data()
		if err != nil {
			return nil, err
		}
		datas = append(datas, dr)
	}
	return l.merge(datas...), nil
}
