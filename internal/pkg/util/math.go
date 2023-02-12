// Package util is the package that contains all of the utility functions and types.
package util

import "golang.org/x/exp/constraints"

// OrderedMin returns the minimum of x and y.
func OrderedMin[T constraints.Ordered](x T, y T) T {
	if x < y {
		return x
	}

	return y
}
