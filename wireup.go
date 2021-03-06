package main

import (
	stdhttp "net/http"
	"sync"

	"github.com/joliver/docpile/app/domain"
	"github.com/joliver/docpile/app/events"
	"github.com/joliver/docpile/app/http"
	"github.com/joliver/docpile/app/projections"
	"github.com/joliver/docpile/generic/applicators"
	"github.com/joliver/docpile/generic/eventstore"
	"github.com/joliver/docpile/generic/handlers"
	"github.com/joliver/docpile/generic/identity"
	"github.com/joliver/docpile/generic/serialization"
	"github.com/joliver/docpile/generic/storage"
	"github.com/julienschmidt/httprouter"
	"github.com/smartystreets/detour"
	"github.com/smartystreets/httpx"
	"github.com/smartystreets/httpx/middleware"
	"github.com/smartystreets/listeners"
)

type Wireup struct {
	indexPath string
	dataPath  string
	mutex     *sync.RWMutex
	listeners []listeners.Listener
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
	mutexApplicator := applicators.NewMutex(projector, this.mutex)
	channelApplicator := applicators.NewChannel(mutexApplicator)
	eventstoreApplicator := eventstore.NewApplicator(channelApplicator, store)

	domainHandler := handlers.NewDomain(aggregate, eventstoreApplicator)
	channelHandler := handlers.NewChannel(domainHandler)

	this.listeners = append(this.listeners, channelHandler, channelApplicator)

	return domain.NewWriteAssetHandler(channelHandler, storage.NewFileStorage(this.dataPath))
}

func (this *Wireup) BuildHTTPHandler(application handlers.Handler, projector *projections.Projector) stdhttp.Handler {
	tagWriter := http.NewTagWriter(application)
	assetWriter := http.NewAssetWriter(application)
	documentWriter := http.NewDocumentWriter(application)
	reader := http.NewReader(projector.AllTags, projector.AllDocuments)
	search := http.NewSearch(projector.AllDocuments, projector.MatchingTags)

	router := buildRouter()
	router.Handler("OPTIONS", "/*wildcard", middleware.OriginCORSHeadersHandler("http://localhost:3000"))

	router.Handler("PUT", "/tags", this.writerAction(tagWriter.Add))
	router.Handler("DELETE", "/tags/:id", this.writerAction(tagWriter.Remove))
	router.Handler("POST", "/tags/:id/name", this.writerAction(tagWriter.Rename))
	router.Handler("PUT", "/tags/:id/name", this.writerAction(tagWriter.DefineSynonym))
	router.Handler("DELETE", "/tags/:id/name", this.writerAction(tagWriter.RemoveSynonym))

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
	handler := middleware.BrowserHeadersHandler(map[string]string{"Access-Control-Allow-Origin": "*"})
	handler.Install(detour.New(action))
	return handler
}
func (this *Wireup) readerAction(action interface{}) stdhttp.Handler {
	handler := middleware.BrowserHeadersHandler(map[string]string{"Access-Control-Allow-Origin": "*"})
	handler.Install(detour.New(action))
	return middleware.NewLockHandler(this.mutex.RLocker(), handler)
}

func buildRouter() *httprouter.Router {
	router := httprouter.New()
	router.HandleMethodNotAllowed = true
	router.HandleOPTIONS = true
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false
	return router
}

func (this *Wireup) BuildHTTPServer(listenAddress string, handler stdhttp.Handler) {
	server := httpx.NewHTTPServer(listenAddress, handler)
	this.listeners = append([]listeners.Listener{server}, this.listeners...) // prepend HTTP server
}

func (this *Wireup) BuildListener() listeners.ListenCloser {
	listener := listeners.NewCascadingWaitListener(this.listeners...)
	return listeners.NewCompositeWaitShutdownListener(listener)
}
