package main

import (
	"reflect"
	"unsafe"
)

func main() {

}

const eps = 1.0e-10

func equal(x, y reflect.Value, seen map[comparison]bool) bool {
	if !x.IsValid() || !y.IsValid() {
		return x.IsValid() == y.IsValid()
	}

	if x.Type() != y.Type() {
		return false
	}

	if x.CanAddr() && y.CanAddr() {
		xptr := unsafe.Pointer(x.UnsafeAddr())
		yptr := unsafe.Pointer(y.UnsafeAddr())
		if xptr == yptr {
			return true
		}

		c := comparison{
			x: xptr,
			y: yptr,
			t: x.Type(),
		}
		if seen[c] {
			return true
		}
		seen[c] = true
	}

	// todo
	panic("unreachable")
}

type comparison struct {
	x, y unsafe.Pointer
	t    reflect.Type
}
