package utils

import "github.com/go-errors/errors"

// Returns a pointer to the given value
func Ptr[V any](v V) *V {
	return &v
}

// Wraps an error with a stack trace
func WrapErr(err error) error {
	return errors.Wrap(err,1)
}

// Concatenates two slices into a new slice
func Concat[T any](item []T, item2 []T) []T {
	itemLength := len(item)

	new := make([]T, itemLength, itemLength + len(item2))

	copy(new, item)
	new = append(new, item2...)

	return new
}