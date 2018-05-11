package infrastructure

type Handler interface {
	Handle(interface{}) Result
}

type Applicator interface {
	Apply([]interface{})
}

type Result struct {
	ID    uint64
	Error error
}