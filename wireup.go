package main

import (
	stdhttp "net/http"

	"bitbucket.org/jonathanoliver/docpile/app/domain"
	"bitbucket.org/jonathanoliver/docpile/app/events"
	"bitbucket.org/jonathanoliver/docpile/app/http"
	"bitbucket.org/jonathanoliver/docpile/generic/applicators"
	"bitbucket.org/jonathanoliver/docpile/generic/eventstore"
	"bitbucket.org/jonathanoliver/docpile/generic/handlers"
	"bitbucket.org/jonathanoliver/docpile/generic/identity"
	"bitbucket.org/jonathanoliver/docpile/generic/serialization"
	"bitbucket.org/jonathanoliver/docpile/generic/storage"
	"github.com/julienschmidt/httprouter"
	"github.com/smartystreets/detour"
)

type Wireup struct {
	indexPath string
	dataPath  string
}

func NewWireup(indexPath string, dataPath string) *Wireup {
	return &Wireup{indexPath: indexPath, dataPath: dataPath}
}

func (this *Wireup) BuildDomain() *domain.Aggregate {
	return domain.NewAggregate(identity.NewEpochGenerator())
}

func (this *Wireup) BuildEventStore(aggregate handlers.Aggregate) *eventstore.DelimitedText {
	return eventstore.NewDelimitedText(
		storage.NewFileStorage(this.indexPath, storage.Append(), storage.EnsureWorkspace()),
		this.buildTypeRegistry(),
		serialization.NewJSONSerializer())
}
func (this *Wireup) buildTypeRegistry() eventstore.TypeRegistry {
	var registry = eventstore.NewRegistry(eventstore.PanicOnUnknownType())

	registry.Add("tag-added", events.TagAdded{})
	registry.Add("tag-removed", events.TagRenamed{})
	registry.Add("tag-synonym-defined", events.TagSynonymDefined{})
	registry.Add("tag-synonym-removed", events.TagSynonymRemoved{})
	registry.Add("managed-asset-imported", events.ManagedAssetImported{})
	registry.Add("cloud-asset-imported", events.CloudAssetImported{})
	registry.Add("document-defined", events.DocumentDefined{})

	return registry
}

func (this *Wireup) BuildMessageHandler(aggregate handlers.Aggregate, store eventstore.EventStore) handlers.Handler {
	var applicator = this.buildApplicator(store)

	var application handlers.Handler = handlers.NewDomain(aggregate, applicator)
	application = handlers.NewChannel(application, handlers.StartChannel())
	application = domain.NewWriteAssetHandler(application, storage.NewFileStorage(this.dataPath))
	return application
}
func (this *Wireup) buildApplicator(store eventstore.EventStore) applicators.Applicator {
	applicator := SampleApplicator()
	applicator = eventstore.NewApplicator(applicator, store)
	return applicator
}

func (this *Wireup) BuildHTTPHandler(application handlers.Handler) stdhttp.Handler {
	tagController := http.NewTagWriteController(application)
	assetController := http.NewAssetWriteController(application)
	documentController := http.NewDocumentWriteController(application)

	router := buildRouter()
	router.Handler("PUT", "/tags", detour.New(tagController.Add))
	router.Handler("POST", "/tags/name", detour.New(tagController.Rename))
	router.Handler("PUT", "/tags/synonym", detour.New(tagController.DefineSynonym))
	router.Handler("DELETE", "/tags/synonym", detour.New(tagController.RemoveSynonym))

	router.Handler("PUT", "/assets", detour.New(assetController.ImportManaged))
	router.Handler("PUT", "/documents", detour.New(documentController.Define))

	return router

	// GET /search/documents = document search criteria
	// GET /search/tags = tag auto-complete search

	// apply/remove one or more tags to a single document
	//   PUT /documents/:id/tags
	//   DELETE /document/:id/tags/:tags
	// apply/remove one tag to one or more documents
	//   PUT /tags/:id/documents
	//   DELETE /tags/:id/documents/:documents
}
func buildRouter() *httprouter.Router {
	router := httprouter.New()
	router.HandleMethodNotAllowed = true
	router.HandleOPTIONS = true
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false
	return router
}
