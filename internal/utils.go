package internal

// fuck you go
func Ptr[T any](val T) *T {
	return &val
}
