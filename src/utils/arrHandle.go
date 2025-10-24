package utils

func Delete[T any](slice []T, index int) []T {
	if index < 0 || index >= len(slice) {
		return slice
	}
	return append(slice[:index], slice[index+1:]...)
}
func Insert[T any](slice []T, index int, value T) []T {
	if index < 0 || index > len(slice) {
		return slice
	}
	return append(slice[:index], append([]T{value}, slice[index:]...)...)
}
func Replace[T any](slice []T, index int, value T) []T {
	if index < 0 || index >= len(slice) {
		return slice
	}
	slice[index] = value
	return slice
}
