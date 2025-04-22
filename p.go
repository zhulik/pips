package pips

// P is a pair.
type P[A any, B any] interface {
	Unpack() (A, B)

	A() A
	B() B
}

// NewP creates a new pair.
func NewP[A any, B any](a A, b B) P[A, B] {
	return pp[A, B]{a, b}
}

type pp[A any, B any] struct {
	a A
	b B
}

// Unpack returns the pair.
func (p pp[A, B]) Unpack() (A, B) {
	return p.a, p.b
}

// A returns the first element of the pair.
func (p pp[A, B]) A() A {
	return p.a
}

// B returns the second element of the pair.
func (p pp[A, B]) B() B {
	return p.b
}
