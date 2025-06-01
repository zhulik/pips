package pips

import "errors"

var (
	ErrPanicInStage = errors.New("panic in stage")
	ErrWrongType    = errors.New("value cannot be casted")
)
