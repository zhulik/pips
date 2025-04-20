package pips

type D[T any] interface {
	Unpack() (T, error)

	Value() T
	Error() error
}

func AnyD[T any](value T) D[any] {
	return NewD[any](value)
}

func NewD[T any](value T) D[T] {
	return pd[T]{value, nil}
}

func ErrD[T any](err error) D[T] {
	var t T
	return pd[T]{t, err}
}

func CastD[I any, O any](d D[I]) D[O] {
	if d.Error() != nil {
		return ErrD[O](d.Error())
	}
	return NewD(any(d.Value()).(O))
}

type pd[T any] struct {
	value T
	error error
}

func (r pd[T]) Value() T {
	return r.value
}

func (r pd[T]) Error() error {
	return r.error
}

func (r pd[T]) Unpack() (T, error) {
	return r.value, r.error
}
