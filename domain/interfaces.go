package domain

import (
	"errors"
	"io"
	"time"
)

type IdentityGenerator interface {
	Next() uint64
}

var (
	TagAlreadyExistsError   = errors.New("tag already exists")
	AssetAlreadyExistsError = errors.New("asset already exists")
	AssetNotFoundError      = errors.New("asset not found")
	TagNotFoundError        = errors.New("tag not found")
	DocumentNotFoundError   = errors.New("document not found")
	StoreAssetError         = errors.New("unable to storage asset")
)

type (
	TagAdder interface {
		AddTag(string) (uint64, error)
	}
	DocumentDefiner interface {
		DefineDocument(DocumentDefinition) (uint64, error)
	}
	ManagedAssetImporter interface {
		ImportManagedAsset(string, string, io.ReadCloser) (uint64, error)
	}
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
