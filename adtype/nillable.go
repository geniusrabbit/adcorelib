package adtype

// IsNillable describes a type that can report whether it is a typed nil.
// Implementations must be nil-safe: a nil receiver must return true without panicking.
type IsNillable interface {
	IsNil() bool
}

// IsNil reports whether v is an untyped nil interface or a typed nil.
func IsNil(v IsNillable) bool {
	return v == nil || v.IsNil()
}
