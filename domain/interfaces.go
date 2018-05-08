package domain

type IdentityGenerator interface {
	Next() uint64
}
