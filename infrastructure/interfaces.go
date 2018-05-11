package infrastructure

type Handler interface {
	Handle(interface{}) (uint64, error)
}

type Applicator interface {
	Apply([]interface{})
}
