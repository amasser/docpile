package projections

import (
	"errors"
	"time"

	"bitbucket.org/jonathanoliver/docpile/app/events"
)

type AllDocuments struct {
	index map[uint64]int
	items []Document
}

func NewAllDocuments() *AllDocuments {
	return &AllDocuments{index: map[uint64]int{}}
}

func (this *AllDocuments) Transform(message interface{}) {
	switch message := message.(type) {
	case events.DocumentDefined:
		this.documentDefined(message)
	}
}
func (this *AllDocuments) documentDefined(message events.DocumentDefined) {
	if _, contains := this.index[message.DocumentID]; !contains {
		this.index[message.DocumentID] = len(this.items)
		this.items = append(this.items, newDocument(message))
	}
}

func (this *AllDocuments) List() []Document { return this.items }
func (this *AllDocuments) Load(id uint64) (Document, error) {
	if index, contains := this.index[id]; contains {
		return this.items[index], nil
	} else {
		return Document{}, DocumentNotFoundError
	}
}

var (
	DocumentNotFoundError = errors.New("document not found")
)

//////////////////////////////////////////////////////////////

type Document struct {
	DocumentID  uint64     `json:"document_id"`
	Timestamp   time.Time  `json:"timestamp"`
	AssetID     uint64     `json:"asset_id"`
	AssetOffset uint64     `json:"asset_offset,omitempty"`
	Published   *time.Time `json:"published,omitempty"`
	PeriodBegin *time.Time `json:"effective_begin,omitempty"`
	PeriodEnd   *time.Time `json:"effective_end,omitempty"`
	Tags        []uint64   `json:"tags,omitempty"`
	Documents   []uint64   `json:"documents,omitempty"`
	Description string     `json:"description,omitempty"`
}

func newDocument(message events.DocumentDefined) Document {
	return Document{
		DocumentID:  message.DocumentID,
		Timestamp:   message.Timestamp,
		AssetID:     message.AssetID,
		AssetOffset: message.AssetOffset,
		Published:   message.Published,
		PeriodBegin: message.PeriodBegin,
		PeriodEnd:   message.PeriodEnd,
		Tags:        message.Tags,
		Documents:   message.Documents,
		Description: message.Description,
	}
}
