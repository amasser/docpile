package domain

import (
	"errors"
	"fmt"
	"log"
	"path"
	"reflect"
	"strings"

	"bitbucket.org/jonathanoliver/docpile/app/events"
	"bitbucket.org/jonathanoliver/docpile/generic/handlers"
	"bitbucket.org/jonathanoliver/docpile/generic/identity"
	"github.com/smartystreets/clock"
)

type Aggregate struct {
	events               []interface{}
	clock                *clock.Clock
	identity             identity.Generator
	tagsByID             map[uint64]string
	tagsByNormalizedName map[string]uint64
	assetsByID           map[uint64]struct{}
	managedAssets        map[events.SHA256Hash]uint64
	cloudAssets          map[string]struct{}
	documentsByID        map[uint64]struct{}
}

func NewAggregate(identity identity.Generator) *Aggregate {
	return &Aggregate{
		identity:             identity,
		tagsByID:             make(map[uint64]string),
		tagsByNormalizedName: make(map[string]uint64),
		assetsByID:           make(map[uint64]struct{}),
		managedAssets:        make(map[events.SHA256Hash]uint64),
		cloudAssets:          make(map[string]struct{}),
		documentsByID:        make(map[uint64]struct{}),
	}
}

func (this *Aggregate) Handle(message interface{}) handlers.Result {
	switch message := message.(type) {

	case AddTag:
		return this.AddTag(message.Name)
	case RenameTag:
		return this.RenameTag(message.ID, message.Name)
	case DefineTagSynonym:
		return this.DefineTagSynonym(message.ID, message.Name)
	case RemoveTagSynonym:
		return this.RemoveTagSynonym(message.ID, message.Name)

	case ImportManagedAsset:
		return this.ImportManagedAsset(message.Name, message.MIMEType, message.Hash)
	case ImportCloudAsset:
		return this.ImportCloudAsset(message.Name, message.Provider, message.Resource)

	case DefineDocument:
		return this.DefineDocument(message.Document)

	default:
		log.Panicf(fmt.Sprintf("Aggregate cannot handle '%s'", reflect.TypeOf(message)))
		return newResult(0, nil)
	}
}

func (this *Aggregate) AddTag(name string) handlers.Result {
	if id, contains := this.tagsByNormalizedName[normalizeTag(name)]; contains {
		return newResult(id, TagAlreadyExistsError)
	}

	id := this.identity.Next()
	return this.raise(id, events.TagAdded{
		TagID:     id,
		Timestamp: this.clock.UTCNow(),
		TagName:   name,
	})
}
func (this *Aggregate) RenameTag(id uint64, name string) handlers.Result {
	if result := this.validTagInput(id, name); result.Error != nil {
		return result
	}

	return this.raise(id, events.TagRenamed{
		TagID:     id,
		Timestamp: this.clock.UTCNow(),
		OldName:   this.tagsByID[id],
		NewName:   name,
	})
}
func (this *Aggregate) DefineTagSynonym(id uint64, name string) handlers.Result {
	if result := this.validTagInput(id, name); result.Error != nil {
		return result
	}

	return this.raise(id, events.TagSynonymDefined{
		TagID:     id,
		Timestamp: this.clock.UTCNow(),
		TagName:   name,
	})
}
func (this *Aggregate) RemoveTagSynonym(id uint64, name string) handlers.Result {
	if result := this.validTagInput(id, name); result.Error != nil {
		return result
	}

	return this.raise(id, events.TagSynonymRemoved{
		TagID:     id,
		Timestamp: this.clock.UTCNow(),
		TagName:   name,
	})
}
func (this *Aggregate) validTagInput(id uint64, name string) handlers.Result {
	if _, contains := this.tagsByID[id]; !contains {
		return newResult(0, TagNotFoundError)
	} else if id, contains = this.tagsByNormalizedName[normalizeTag(name)]; contains {
		return newResult(id, TagAlreadyExistsError)
	} else {
		return newResult(id, nil)
	}
}
func normalizeTag(value string) string {
	return strings.ToLower(value)
}

func (this *Aggregate) ImportManagedAsset(name, mime string, hash events.SHA256Hash) handlers.Result {
	if id, contains := this.managedAssets[hash]; contains {
		return newResult(id, AssetAlreadyExistsError)
	}

	id := this.identity.Next()
	return this.raise(id, events.ManagedAssetImported{
		AssetID:   id,
		Timestamp: this.clock.UTCNow(),
		Hash:      events.SHA256Hash(hash),
		MIMEType:  mime,
		Name:      name,
		Key:       fmt.Sprintf("%d%s", id, path.Ext(name)),
	})
}
func (this *Aggregate) ImportCloudAsset(name, provider, resource string) handlers.Result {
	if _, contains := this.cloudAssets[normalizeCloudAsset(provider, resource)]; contains {
		return newResult(0, AssetAlreadyExistsError)
	}

	id := this.identity.Next()
	return this.raise(id, events.CloudAssetImported{
		AssetID:   id,
		Timestamp: this.clock.UTCNow(),
		Name:      name,
		Provider:  provider,
		Resource:  resource,
	})
}
func normalizeCloudAsset(provider, resource string) string {
	return fmt.Sprintf("%s.%s", strings.ToLower(provider), resource)
}

func (this *Aggregate) DefineDocument(doc DocumentDefinition) handlers.Result {
	if err := this.validDefinition(doc); err != nil {
		return newResult(0, err)
	}

	id := this.identity.Next()
	return this.raise(id, events.DocumentDefined{
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

func (this *Aggregate) raise(id uint64, event interface{}) handlers.Result {
	this.apply(event)
	this.events = append(this.events, event)
	return newResult(id, nil)
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
	this.managedAssets[message.Hash] = message.AssetID
}
func (this *Aggregate) applyCloudAssetImported(message events.CloudAssetImported) {
	this.assetsByID[message.AssetID] = struct{}{}
	this.cloudAssets[normalizeCloudAsset(message.Provider, message.Resource)] = struct{}{}
}

func (this *Aggregate) applyDocumentDefined(message events.DocumentDefined) {
	this.documentsByID[message.DocumentID] = struct{}{}
}

func (this *Aggregate) Consume() []interface{} {
	consumed := this.events
	this.events = nil // don't re-use the buffer
	return consumed
}

func newResult(id uint64, err error) handlers.Result {
	return handlers.Result{ID: id, Error: err}
}

var (
	TagAlreadyExistsError   = errors.New("tag already exists")
	AssetAlreadyExistsError = errors.New("asset already exists")
	AssetNotFoundError      = errors.New("asset not found")
	TagNotFoundError        = errors.New("tag not found")
	DocumentNotFoundError   = errors.New("document not found")
	StoreAssetError         = errors.New("unable to storage asset")
)