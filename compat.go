package replica

// CopyAny creates a deep copy of v. Drop-in replacement for
// copystructure.Copy. Returns the copied value as any.
func CopyAny(v any) (any, error) {
	return Copy(v)
}

// MustAny is like CopyAny but panics on error. Drop-in replacement for
// copystructure.Must.
func MustAny(v any) any {
	return Must(v)
}
