package domain

import "bitbucket.org/jonathanoliver/docpile/infrastructure"

func newResult(id uint64, err error) infrastructure.Result {
	return infrastructure.Result{ID: id, Error: err}
}
