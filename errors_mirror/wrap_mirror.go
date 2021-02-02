package errors_mirror

import "unsafe"

func Unwrap(err error) error {
	u, ok := err.(interface {
		Unwrap() error
	})
	if !ok {
		return nil
	}
	return u.Unwrap()
}

func Is(err, target error) bool {
	if target == nil {
		return err == target
	}
	// todo
	return false
}

type emptyinterface struct {
	typ  *rtype
	word unsafe.Pointer
}
