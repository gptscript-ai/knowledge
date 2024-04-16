// Package z exposes a curated set of utility functions.
//
// Its primary goals are two-fold:
//
//  1. Make the best commonly used utilities discoverable
//  2. Reduce hand strain; "I can finally stop writing `func must(...)` everywhere!"
package z

import (
	"errors"
)

// Must panics IFF any error in errs is non-nil.
func Must(errs ...error) {
	if err := errors.Join(errs...); err != nil {
		panic(err)
	}
}

// MustBe panics IFF err is non-nil, otherwise it returns t.
func MustBe[T any](t T, err error) T {
	Must(err)
	return t
}

// Pointer returns a pointer to v.
func Pointer[T any](v T) *T {
	return &v
}

// Dereference returns the dereferenced value of p.
// If p is nil, the zero value of T is returned instead.
// This function is intended to be used to dereference native types (e.g. *int, *string, etc.);
// for structs, use of Dereference may decrease readability and obscure intent, so prefer a conditional
// (e.g. `if p != nil { ... }`) instead.
func Dereference[T any](p *T) (v T) {
	if p != nil {
		v = *p
	}
	return
}

// AddToMap will add the key value pair to the map, ensuring that the map is not nil.
func AddToMap[K comparable, V any](m map[K]V, k K, v V) map[K]V {
	return ConcatMaps(m, map[K]V{k: v})
}

// ConcatMaps will iteratively add all the key/value pairs from each map, overwriting existing keys.
// That is, if every map in ms has the same key, then the value of that key in the resulting map will be
// the value of the key in the last map in ms.
// If no maps are provided, then a nil map is returned.
// If at least one map is provided (nil or not), then a non-nil map is returned.
func ConcatMaps[K comparable, V any](ms ...map[K]V) map[K]V {
	if len(ms) == 0 {
		return nil
	}

	m := ms[0]
	if m == nil {
		m = make(map[K]V)
	}
	for i := 1; i < len(ms); i++ {
		for k, v := range ms[i] {
			m[k] = v
		}
	}

	return m
}
