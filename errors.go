package pips

import "errors"

var (
	ErrPanicInStage = errors.New("panic in stage")
	ErrWrongType    = errors.New("wrong stage input or output type")
)
