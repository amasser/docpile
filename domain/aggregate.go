package domain

import (
	"fmt"
	"log"
	"path"
	"reflect"
	"strings"

	"bitbucket.org/jonathanoliver/docpile/events"
	"github.com/smartystreets/clock"
)

type Aggregate struct {
	events               []interface{}
	clock                *clock.Clock
	identity             IdentityGenerator
	tagsByID             map[uint64]string
	tagsByNormalizedName map[string]uint64
	assetsByID           map[uint64]struct{}
	managedAssets        map[managedKey]uint64
	cloudAssets          map[cloudKey]struct{}
	documentsByID        map[uint64]struct{}
}

func NewAggregate(identity IdentityGenerator) *Aggregate {
	return &Aggregate{
		identity:             identity,
		tagsByID:             make(map[uint64]string),
		tagsByNormalizedName: make(map[string]uint64),
		assetsByID:           make(map[uint64]struct{}),
		managedAssets:        make(map[managedKey]uint64),
		cloudAssets:          make(map[cloudKey]struct{}),
		documentsByID:        make(map[uint64]struct{}),
	}
}

func (this *Aggregate) AddTag(name string) (uint64, error) {
	if id, contains := this.tagsByNormalizedName[normalizeTag(name)]; contains {
		return id, TagAlreadyExistsError
	}

	id := this.identity.Next()
	return id, this.raise(events.TagAdded{
		TagID:     id,
		Timestamp: this.clock.UTCNow(),
		TagName:   name,
	})
}
func (this *Aggregate) RenameTag(id uint64, name string) (uint64, error) {
	if id, err := this.validTagInput(id, name); err != nil {
		return id, err
	}

	return id, this.raise(events.TagRenamed{
		TagID:     id,
		Timestamp: this.clock.UTCNow(),
		OldName:   this.tagsByID[id],
		NewName:   name,
	})
}
func (this *Aggregate) DefineTagSynonym(id uint64, name string) (uint64, error) {
	if id, err := this.validTagInput(id, name); err != nil {
		return id, err
	}

	return id, this.raise(events.TagSynonymDefined{
		TagID:     id,
		Timestamp: this.clock.UTCNow(),
		TagName:   name,
	})
}
func (this *Aggregate) RemoveTagSynonym(id uint64, name string) (uint64, error) {
	if id, err := this.validTagInput(id, name); err != nil {
		return id, err
	}

	return id, this.raise(events.TagSynonymRemoved{
		TagID:     id,
		Timestamp: this.clock.UTCNow(),
		TagName:   name,
	})
}
func (this *Aggregate) validTagInput(id uint64, name string) (uint64, error) {
	if _, contains := this.tagsByID[id]; !contains {
		return 0, TagNotFoundError
	} else if id, contains = this.tagsByNormalizedName[normalizeTag(name)]; contains {
		return id, TagAlreadyExistsError
	} else {
		return id, nil
	}
}
func normalizeTag(value string) string {
	return strings.ToLower(value)
}

func (this *Aggregate) ImportManagedAsset(name, mime string, hash events.SHA256Hash) (uint64, error) {
	if id, contains := this.managedAssets[managedKey(hash)]; contains {
		return id, AssetAlreadyExistsError
	}

	id := this.identity.Next()
	return id, this.raise(events.ManagedAssetImported{
		AssetID:   id,
		Timestamp: this.clock.UTCNow(),
		Hash:      events.SHA256Hash(hash),
		MIMEType:  mime,
		Name:      name,
		Key:       fmt.Sprintf("%d%s", id, path.Ext(name)),
	})
}
func (this *Aggregate) ImportCloudAsset(name, provider, resource string) (uint64, error) {
	if _, contains := this.cloudAssets[newCloudAssetKey(provider, resource)]; contains {
		return 0, AssetAlreadyExistsError
	}

	id := this.identity.Next()
	return id, this.raise(events.CloudAssetImported{
		AssetID:   id,
		Timestamp: this.clock.UTCNow(),
		Name:      name,
		Provider:  provider,
		Resource:  resource,
	})
}

func (this *Aggregate) DefineDocument(doc DocumentDefinition) (uint64, error) {
	if err := this.validDefinition(doc); err != nil {
		return 0, err
	}

	id := this.identity.Next()
	return id, this.raise(events.DocumentDefined{
		DocumentID:  id,
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
			return DocumentNotFoundError
		}
	}

	// FUTURE: other invariants, e.g. begin after end, end without begin
	return nil
}

func (this *Aggregate) raise(event interface{}) error {
	this.apply(event)
	this.events = append(this.events, event)
	return nil
}

func (this *Aggregate) Apply(messages ...interface{}) {
	for _, message := range messages {
		this.apply(message)
	}
}
func (this *Aggregate) apply(message interface{}) {
	switch message := message.(type) {

	case events.TagAdded:
		this.applyTagAdded(message)
	case events.TagRenamed:
		this.applyTagRenamed(message)
	case events.TagSynonymDefined:
		this.applyTagSynonymDefined(message)
	case events.TagSynonymRemoved:
		this.applyTagSynonymRemoved(message)

	case events.ManagedAssetImported:
		this.applyManagedAssetImported(message)
	case events.CloudAssetImported:
		this.applyCloudAssetImported(message)

	case events.DocumentDefined:
		this.applyDocumentDefined(message)
	default:
		log.Panicf(fmt.Sprintf("Aggregate cannot apply '%s'", reflect.TypeOf(message)))
	}
}

func (this *Aggregate) applyTagAdded(message events.TagAdded) {
	this.tagsByID[message.TagID] = message.TagName // full, not-normalized value
	this.tagsByNormalizedName[normalizeTag(message.TagName)] = message.TagID
}
func (this *Aggregate) applyTagRenamed(message events.TagRenamed) {
	this.tagsByID[message.TagID] = message.NewName // full, not-normalized value
	delete(this.tagsByNormalizedName, normalizeTag(message.OldName))
	this.tagsByNormalizedName[normalizeTag(message.NewName)] = message.TagID
}
func (this *Aggregate) applyTagSynonymDefined(message events.TagSynonymDefined) {
	this.tagsByNormalizedName[normalizeTag(message.TagName)] = message.TagID
}
func (this *Aggregate) applyTagSynonymRemoved(message events.TagSynonymRemoved) {
	delete(this.tagsByNormalizedName, normalizeTag(message.TagName))
}

func (this *Aggregate) applyManagedAssetImported(message events.ManagedAssetImported) {
	this.assetsByID[message.AssetID] = struct{}{}
	this.managedAssets[managedKey(message.Hash)] = message.AssetID
}
func (this *Aggregate) applyCloudAssetImported(message events.CloudAssetImported) {
	this.assetsByID[message.AssetID] = struct{}{}
	this.cloudAssets[newCloudAssetKey(message.Provider, message.Resource)] = struct{}{}
}

func (this *Aggregate) applyDocumentDefined(message events.DocumentDefined) {
	this.documentsByID[message.DocumentID] = struct{}{}
}

func (this *Aggregate) Consume() []interface{} {
	consumed := this.events
	this.events = nil // don't re-use the buffer
	return consumed
}
