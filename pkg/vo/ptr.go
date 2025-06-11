package vo

func ToPointer[T any](val T) *T {
	return &val
}

func FromPointer[T any](ptr *T) T {
	if ptr == nil {
		var zeroValue T
		return zeroValue
	}
	return *ptr
}
