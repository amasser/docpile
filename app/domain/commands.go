package domain

import (
	"io"
	"time"

	"bitbucket.org/jonathanoliver/docpile/app/events"
)

type AddTag struct {
	Name string
}
type RenameTag struct {
	ID   uint64
	Name string
}
type DefineTagSynonym struct {
	ID   uint64
	Name string
}
type RemoveTagSynonym struct {
	ID   uint64
	Name string
}

type ImportManagedStreamingAsset struct {
	Name     string
	MIMEType string
	Size     uint64
	Body     io.ReadCloser
}
type ImportManagedAsset struct {
	Name     string
	MIMEType string
	Hash     events.SHA256Hash
}

type ImportCloudAsset struct {
	Name     string
	Provider string
	Resource string
}

type DefineDocument struct {
	Document DocumentDefinition
}

type DocumentDefinition struct {
	AssetID     uint64
	AssetOffset uint64
	Published   *time.Time
	PeriodMin   *time.Time
	PeriodMax   *time.Time
	Tags        []uint64
	Documents   []uint64
	Description string
}

type RemoveDocument struct {
	ID uint64
}
