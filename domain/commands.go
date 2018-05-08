package domain

import "bitbucket.org/jonathanoliver/docpile/events"

type AddTag struct {
	Name string
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
