package utils

func Find[T any](array []T, fn func(T) bool) int {
	for i := range array {
		elem := array[i]
		if fn(elem) {
			return i
		}
	}

	return -1
}
