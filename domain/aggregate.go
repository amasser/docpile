package domain

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"bitbucket.org/jonathanoliver/docpile/events"
	"github.com/smartystreets/clock"
)

type Aggregate struct {
	events        []interface{}
	clock         *clock.Clock
	identity      IdentityGenerator
	tagsByID      map[uint64]struct{}
	tagsByName    map[string]struct{}
	assetsByID    map[uint64]struct{}
	managedAssets map[managedAssetKey]struct{}
	cloudAssets   map[cloudAssetKey]struct{}
	documentsByID map[uint64]struct{}
}
type (
	cloudAssetKey   string
	managedAssetKey events.SHA256Hash
)

func NewAggregate(identity IdentityGenerator) *Aggregate {
	return &Aggregate{
		identity:      identity,
		tagsByID:      make(map[uint64]struct{}),
		tagsByName:    make(map[string]struct{}),
		assetsByID:    make(map[uint64]struct{}),
		managedAssets: make(map[managedAssetKey]struct{}),
		cloudAssets:   make(map[cloudAssetKey]struct{}),
		documentsByID: make(map[uint64]struct{}),
	}
}

func (this *Aggregate) AddTag(name string) error {
	if _, contains := this.tagsByName[strings.ToLower(name)]; contains {
		return TagAlreadyExistsError
	}

	return this.raise(events.TagAdded{
		TagID:     this.identity.Next(),
		Timestamp: this.clock.UTCNow(),
		TagName:   name,
	})
}
func (this *Aggregate) ImportManagedAsset(name, mime string, hash managedAssetKey) error {
	if _, contains := this.managedAssets[hash]; contains {
		return AssetAlreadyExistsError
	}

	return this.raise(events.ManagedAssetImported{
		AssetID:   this.identity.Next(),
		Timestamp: this.clock.UTCNow(),
		Hash:      events.SHA256Hash(hash),
		MIMEType:  mime,
		Name:      name,
	})
}
func (this *Aggregate) ImportCloudAsset(name, provider, resource string) error {
	if _, contains := this.cloudAssets[newCloudAssetKey(provider, resource)]; contains {
		return AssetAlreadyExistsError
	}

	return this.raise(events.CloudAssetImported{
		AssetID:   this.identity.Next(),
		Timestamp: this.clock.UTCNow(),
		Name:      name,
		Provider:  provider,
		Resource:  resource,
	})
}
func (this *Aggregate) DefineDocument(doc DocumentDefinition) error {
	if err := this.validDefinition(doc); err != nil {
		return err
	}

	this.raise(events.DocumentDefined{
		DocumentID:  this.identity.Next(),
		Timestamp:   this.clock.UTCNow(),
		AssetID:     doc.AssetID,
		AssetOffset: doc.AssetOffset,
		Published:   doc.Published,
		PeriodBegin: doc.PeriodBegin,
		PeriodEnd:   doc.PeriodEnd,
		Tags:        doc.Tags,
		Documents:   doc.Documents,
		Description: doc.Description,
	})

	return nil
}
func (this *Aggregate) validDefinition(doc DocumentDefinition) error {
	if _, contains := this.assetsByID[doc.AssetID]; !contains {
		return AssetNotFoundError
	}

	for _, tagID := range doc.Tags {
		if _, contains := this.tagsByID[tagID]; !contains {
			return TagNotFoundError
		}
	}

	for _, documentID := range doc.Documents {
		if _, contains := this.documentsByID[documentID]; !contains {
			return LinkedDocumentNotFoundError
		}
	}

	// FUTURE: other invariants, e.g. begin after end, end without begin
	return nil
}

func (this *Aggregate) raise(event interface{}) error {
	this.Apply(event)
	this.events = append(this.events, event)
	return nil
}

func (this *Aggregate) Apply(event interface{}) {
	switch event := event.(type) {
	case events.TagAdded:
		this.applyTagAdded(event)
	case events.ManagedAssetImported:
		this.applyManagedAssetImported(event)
	case events.CloudAssetImported:
		this.applyCloudAssetImported(event)
	case events.DocumentDefined:
		this.applyDocumentDefined(event)
	default:
		log.Panicf(fmt.Sprintf("Aggregate cannot apply '%s'", reflect.TypeOf(event)))
	}
}
func (this *Aggregate) applyTagAdded(event events.TagAdded) {
	this.tagsByID[event.TagID] = struct{}{}
	this.tagsByName[strings.ToLower(event.TagName)] = struct{}{}
}
func (this *Aggregate) applyManagedAssetImported(event events.ManagedAssetImported) {
	this.assetsByID[event.AssetID] = struct{}{}
	this.managedAssets[managedAssetKey(event.Hash)] = struct{}{}
}
func (this *Aggregate) applyCloudAssetImported(event events.CloudAssetImported) {
	this.assetsByID[event.AssetID] = struct{}{}
	this.cloudAssets[newCloudAssetKey(event.Provider, event.Resource)] = struct{}{}
}
func (this *Aggregate) applyDocumentDefined(event events.DocumentDefined) {
}

func (this *Aggregate) Consume() []interface{} {
	consumed := this.events
	this.events = nil // don't re-use the buffer
	return consumed
}

func newCloudAssetKey(provider, resource string) cloudAssetKey {
	return cloudAssetKey(fmt.Sprintf("%s.%s", strings.ToLower(provider), resource))
}
