package adsource

// SourceError contains only errors from some source drivers
type SourceError struct {
	Source  any
	Message string
}

func (e *SourceError) Error() string {
	return e.Message
}
