package identity

type Generator interface {
	Next() uint64
}
