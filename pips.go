package pips

import "fmt"

// TryCast tries to cast the value to the given type, return ErrWrongType if fails.
func TryCast[O any](i any) (O, error) {
	o, ok := i.(O)
	if ok {
		return o, nil
	}

	return o, fmt.Errorf("%w: expected: %T, given: %T", ErrWrongType, o, i)
}
