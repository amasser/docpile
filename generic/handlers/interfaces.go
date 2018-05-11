package handlers

type Aggregate interface {
	Handler
	Consume() []interface{}
}

type Handler interface {
	Handle(interface{}) Result
}

type Result struct {
	ID    uint64
	Error error
}
