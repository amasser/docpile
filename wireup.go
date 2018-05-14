package main

import (
	stdhttp "net/http"
	"sync"

	"bitbucket.org/jonathanoliver/docpile/app/domain"
	"bitbucket.org/jonathanoliver/docpile/app/events"
	"bitbucket.org/jonathanoliver/docpile/app/http"
	"bitbucket.org/jonathanoliver/docpile/app/projections"
	"bitbucket.org/jonathanoliver/docpile/generic/applicators"
	"bitbucket.org/jonathanoliver/docpile/generic/eventstore"
	"bitbucket.org/jonathanoliver/docpile/generic/handlers"
	"bitbucket.org/jonathanoliver/docpile/generic/identity"
	"bitbucket.org/jonathanoliver/docpile/generic/serialization"
	"bitbucket.org/jonathanoliver/docpile/generic/storage"
	"bitbucket.org/jonathanoliver/docpile/generic/web"
	"github.com/julienschmidt/httprouter"
	"github.com/smartystreets/detour"
)

type Wireup struct {
	indexPath string
	dataPath  string
	mutex     *sync.RWMutex
}

func NewWireup(indexPath string, dataPath string) *Wireup {
	return &Wireup{indexPath: indexPath, dataPath: dataPath, mutex: &sync.RWMutex{}}
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
	registry.AddMultiple(events.Types)
	return registry
}

func (this *Wireup) BuildProjector() *projections.Projector {
	return projections.NewProjector()
}

func (this *Wireup) BuildMessageHandler(aggregate handlers.Aggregate, store eventstore.EventStore, projector *projections.Projector) handlers.Handler {
	var applicator = this.buildApplicator(store, projector)

	var application handlers.Handler = handlers.NewDomain(aggregate, applicator)
	application = handlers.NewChannel(application, handlers.StartChannel())
	application = domain.NewWriteAssetHandler(application, storage.NewFileStorage(this.dataPath))
	return application
}
func (this *Wireup) buildApplicator(store eventstore.EventStore, projector *projections.Projector) applicators.Applicator {
	applicator := SampleApplicator()
	applicator = applicators.NewFanout(applicator, projector)
	applicator = applicators.NewMutex(applicator, this.mutex)
	applicator = applicators.NewChannel(applicator, applicators.StartChannel())
	applicator = eventstore.NewApplicator(applicator, store)
	return applicator
}

func (this *Wireup) BuildHTTPHandler(application handlers.Handler, projector *projections.Projector) stdhttp.Handler {
	tagWriter := http.NewTagWriter(application)
	assetWriter := http.NewAssetWriter(application)
	documentWriter := http.NewDocumentWriter(application)
	reader := http.NewReader(projector)
	search := http.NewSearch(projector)

	router := buildRouter()
	router.Handler("PUT", "/tags", this.writerAction(tagWriter.Add))
	router.Handler("DELETE", "/tags/:id", this.writerAction(tagWriter.Remove))
	router.Handler("POST", "/tags/:id/name", this.writerAction(tagWriter.Rename))
	router.Handler("PUT", "/tags/:id/synonym", this.writerAction(tagWriter.DefineSynonym))
	router.Handler("DELETE", "/tags/:id/synonym", this.writerAction(tagWriter.RemoveSynonym))

	router.Handler("PUT", "/assets", this.writerAction(assetWriter.ImportManaged))
	router.Handler("PUT", "/documents", this.writerAction(documentWriter.Define))
	router.Handler("DELETE", "/documents/:id", this.writerAction(documentWriter.Remove))

	router.Handler("GET", "/tags", this.readerAction(reader.ListTags))
	router.Handler("GET", "/tags/:id", this.readerAction(reader.LoadTag))
	router.Handler("GET", "/documents", this.readerAction(reader.ListDocuments))
	router.Handler("GET", "/documents/:id", this.readerAction(reader.LoadDocument))

	// these methods don't mutate, but binding is easier when JSON decoding the request body.
	router.Handler("POST", "/search/documents", this.readerAction(search.Documents))
	router.Handler("POST", "/search/tags", this.readerAction(search.Tags))

	return router
}
func (this *Wireup) writerAction(action interface{}) stdhttp.Handler {
	return detour.New(action)
}
func (this *Wireup) readerAction(action interface{}) stdhttp.Handler {
	return web.NewLockHandler(this.mutex.RLocker(), detour.New(action))
}

func buildRouter() *httprouter.Router {
	router := httprouter.New()
	router.HandleMethodNotAllowed = true
	router.HandleOPTIONS = true
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false
	return router
}
