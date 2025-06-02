package apply

type Configer[C any] interface {
	Config() C
}

type Stage[C any] struct {
	config C
	source string
}

func (s Stage[C]) Config() C {
	return s.config
}

func (s Stage[C]) Source() string {
	return s.source
}
