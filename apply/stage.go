package apply

type Configer[C any] interface {
	Config() C
}

type Stage[C any] struct {
	config C
}

func (s Stage[C]) Config() C {
	return s.config
}
