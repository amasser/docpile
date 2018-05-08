package domain

import (
	"errors"
	"time"
)

type IdentityGenerator interface {
	Next() uint64
}

var (
	TagAlreadyExistsError       = errors.New("tag already exists")
	AssetAlreadyExistsError     = errors.New("asset already exists")
	AssetNotFoundError          = errors.New("asset not found")
	TagNotFoundError            = errors.New("tag not found")
	LinkedDocumentNotFoundError = errors.New("linked document not found")
)

type DocumentDefinition struct {
	AssetID     uint64
	AssetOffset uint64
	Published   *time.Time
	PeriodBegin *time.Time
	PeriodEnd   *time.Time
	Tags        []uint64
	Documents   []uint64
	Description string
}

type TagAdder interface {
	AddTag(string) (uint64, error)
}
