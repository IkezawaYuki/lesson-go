package context_mirror

import (
	"errors"
	"time"
)

type Context interface {
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key interface{}) interface{}
}

var Canceled = errors.New("context canceled")

var DeallineExeceeded error = deadlineExceededError{}

type deadlineExceededError struct{}

func (deadlineExceededError) Error() string { return "context deadline exceeded" }

func Background() Context {
	return bacground
}
