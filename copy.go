package replica

import (
	"reflect"
)

// Copy creates a deep copy of v. Returns an error if a circular reference is
// detected or a custom copier fails.
func Copy[T any](v T) (T, error) {
	return CopyWith(v)
}

// Must is like Copy but panics on error. Useful for variable initialization.
func Must[T any](v T) T {
	result, err := Copy(v)
	if err != nil {
		panic(err)
	}
	return result
}

// CopyWith creates a deep copy of v with the given options.
func CopyWith[T any](v T, opts ...Option) (T, error) {
	cfg := newConfig(opts)
	vis := newVisited()

	original := reflect.ValueOf(&v).Elem()
	copied, err := deepCopy(original, cfg, vis)
	if err != nil {
		var zero T
		return zero, err
	}

	return copied.Interface().(T), nil
}
