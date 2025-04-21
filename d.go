package pips

type RawD interface {
	RawUnpack() (any, error)
	RawValue() any
}

type D[T any] interface {
	RawD

	Unpack() (T, error)
	Value() T
	Error() error
}

func AnyD[T any](value T) D[any] {
	return NewD[any](value)
}

func NewD[T any](value T, err ...error) D[T] {
	pdt := pd[T]{value, nil}
	if len(err) > 0 {
		pdt.error = err[0]
	}

	return pdt
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
	value any
	error error
}

func (r pd[T]) RawUnpack() (any, error) {
	return r.value, r.error
}

func (r pd[T]) RawValue() any {
	return r.value
}

func (r pd[T]) Value() T {
	return r.value.(T)
}

func (r pd[T]) Error() error {
	return r.error
}

func (r pd[T]) Unpack() (T, error) {
	return r.Value(), r.error
}
