package domain

import (
	"bitbucket.org/jonathanoliver/docpile/events"
	"io"
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
