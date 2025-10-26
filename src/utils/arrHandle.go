package utils

import "errors"

func Delete[T any](slice *[]T, index int) error {
	if index < 0 || index >= len(*slice) {
		return errors.New("index out of range")
	}
	*slice = append((*slice)[:index], (*slice)[index+1:]...)
	return nil
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
