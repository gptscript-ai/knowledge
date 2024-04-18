package internal

// Number represents a numeric type that can be used for arithmetic operations.
// It includes int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64, float32, and float64.
type Number interface {
	int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64 | float32 | float64
}

// Sum calculates the sum of all numbers in the given slice.
func Sum[T Number](slice []T) T {
	total := T(0)
	for _, num := range slice {
		total += num
	}

	return total
}

// Mean calculates the mean (average) of the numeric values in the given slice.
func Mean[T Number](data []T) float64 {
	if len(data) == 0 {
		return 0
	}

	var sum float64
	for _, d := range data {
		sum += float64(d)
	}

	return sum / float64(len(data))
}

// Values extracts the values from a map and returns them as a slice.
func Values[M ~map[K]V, K comparable, V any](m M) []V {
	r := make([]V, 0, len(m))
	for _, v := range m {
		r = append(r, v)
	}

	return r
}

// Concatenate concatenates multiple slices of type T into a single slice.
// The function takes a variadic number of slices and returns a new slice containing
// all elements of the input slices, preserving the order.
func Concatenate[T any](slices ...[]T) []T {
	var result []T

	resultCap := 0

	for _, s := range slices {
		resultCap += len(s)
	}

	result = make([]T, 0, resultCap)
	for _, s := range slices {
		result = append(result, s...)
	}

	return result
}
