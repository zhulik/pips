package pips

import (
	"context"
	"sync"
)

// F is a Future.
type F[T any] interface {
	// Wait waits for F to be resolved.
	// If F is not yet resolved - waits and returns the result.
	// If F is already resolved - returns the cached result.
	Wait() (T, error)
}

// NewF creates a new Future and a function to resolve it.
func NewF[T any]() (F[T], func(T, error)) {
	f := &ff[T]{done: make(chan any)}

	return f, func(value T, err error) {
		f.m.Lock()
		defer f.m.Unlock()

		if f.value != nil {
			return
		}

		f.value = NewD(value, err)

		close(f.done)
	}
}

// GoF wraps a given function in a Future and resolves it in the background.
func GoF[T any](ctx context.Context, fn func(context.Context) (T, error)) F[T] {
	f, resolve := NewF[T]()

	go resolve(fn(ctx))

	return f
}

type ff[T any] struct {
	m sync.Mutex

	value D[T]
	done  chan any
}

func (f *ff[T]) Wait() (T, error) {
	f.m.Lock()

	if f.value != nil {
		f.m.Unlock()
		return f.value.Unpack()
	}

	<-f.done

	f.m.Unlock()

	return f.value.Unpack()
}
