package utils

// prepend places the element at the beginning of the slice and moves the potentially replaced element to the end.
func Prepend[T any](slice []T, element T) []T {

	if len(slice) == 0 {
		slice = append(slice, element)
		return slice
	}

	slice = append(slice, slice[0])

	slice[0] = element

	return slice

}

func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func Values[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}
