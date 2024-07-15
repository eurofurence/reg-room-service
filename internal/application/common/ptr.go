package common

func Deref[T any](ptr *T) T {
	var def T
	if ptr == nil {
		return def
	}

	return *ptr
}

func To[T any](val T) *T {
	return &val
}
