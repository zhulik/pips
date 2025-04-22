package pips

// D is a pair of value and error, represents a piece of data in the pipeline.
type D[T any] interface {
	Unpack() (T, error)
	Value() T
	Error() error
}

func AnyD[T any](value T) D[any] {
	return NewD[any](value)
}

func NewD[T any](value T, err ...error) D[T] {
	pdt := pd[T]{pp[T, error]{value, nil}}
	if len(err) > 0 {
		pdt.p.b = err[0]
	}

	return pdt
}

func ErrD[T any](err error) D[T] {
	var t T
	return pd[T]{pp[T, error]{t, err}}
}

func CastD[I any, O any](d D[I]) D[O] {
	if d.Error() != nil {
		return ErrD[O](d.Error())
	}
	v := any(d.Value())
	if v == nil {
		var t O
		return NewD(t)
	}
	return NewD(v.(O))
}

type pd[T any] struct {
	p pp[T, error]
}

func (r pd[T]) Value() T {
	return r.p.a
}

func (r pd[T]) Error() error {
	return r.p.b
}

func (r pd[T]) Unpack() (T, error) {
	return r.p.Unpack()
}
