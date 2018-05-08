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
	tags          map[string]struct{}
	managedHashes map[events.SHA256Hash]struct{}
	cloudAssets   map[string]struct{}
}

func NewAggregate(identity IdentityGenerator) *Aggregate {
	return &Aggregate{
		identity:      identity,
		tags:          make(map[string]struct{}),
		managedHashes: make(map[events.SHA256Hash]struct{}),
		cloudAssets:   make(map[string]struct{}),
	}
}

func (this *Aggregate) AddTag(name string) {
	if _, contains := this.tags[strings.ToLower(name)]; contains {
		return
	}

	this.raise(events.TagAdded{
		TagID:     this.identity.Next(),
		Timestamp: this.clock.UTCNow(),
		TagName:   name,
	})
}
func (this *Aggregate) ImportManagedAsset(name, mime string, hash events.SHA256Hash) {
	if _, contains := this.managedHashes[hash]; contains {
		return
	}

	this.raise(events.ManagedAssetImported{
		AssetID:   this.identity.Next(),
		Timestamp: this.clock.UTCNow(),
		Hash:      hash,
		MIMEType:  mime,
		Name:      name,
	})
}
func (this *Aggregate) ImportCloudAsset(name, provider, resource string) {
	if _, contains := this.cloudAssets[composeCloudKey(provider, resource)]; contains {
		return
	}

	this.raise(events.CloudAssetImported{
		AssetID:   this.identity.Next(),
		Timestamp: this.clock.UTCNow(),
		Name:      name,
		Provider:  provider,
		Resource:  resource,
	})
}
func (this *Aggregate) DefineDocument() {
}
func composeCloudKey(provider, resource string) string {
	return fmt.Sprintf("%s.%s", strings.ToLower(provider), resource)
}

func (this *Aggregate) raise(event interface{}) {
	this.Apply(event)
	this.events = append(this.events, event)
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
	this.tags[strings.ToLower(event.TagName)] = struct{}{}
}
func (this *Aggregate) applyManagedAssetImported(event events.ManagedAssetImported) {
	this.managedHashes[event.Hash] = struct{}{}
}
func (this *Aggregate) applyCloudAssetImported(event events.CloudAssetImported) {
	this.cloudAssets[composeCloudKey(event.Provider, event.Resource)] = struct{}{}
}
func (this *Aggregate) applyDocumentDefined(event events.DocumentDefined) {
}

func (this *Aggregate) Consume() []interface{} {
	consumed := this.events
	this.events = nil // don't re-use the buffer
	return consumed
}
