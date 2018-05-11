package main

import (
	"fmt"
	"log"
	stdhttp "net/http"

	"bitbucket.org/jonathanoliver/docpile/app/domain"
	"bitbucket.org/jonathanoliver/docpile/app/events"
	"bitbucket.org/jonathanoliver/docpile/app/http"
	"bitbucket.org/jonathanoliver/docpile/generic"
	"bitbucket.org/jonathanoliver/docpile/generic/eventstore"
	"bitbucket.org/jonathanoliver/docpile/generic/handlers"
	"bitbucket.org/jonathanoliver/docpile/generic/identity"
	"bitbucket.org/jonathanoliver/docpile/generic/serialization"
	"bitbucket.org/jonathanoliver/docpile/generic/storage"
	"github.com/julienschmidt/httprouter"
	"github.com/smartystreets/detour"
)

/*
TODOs
- projections
- better/safer text encoding of event store OR disable advanced characters in input text fields, e.g. tags, description etc.

- tag "synonym" vs "alias" (which word is better)
- apply/remove tags from documents

- document search
- tag search (finding tags that match partial text)
- during document search, as combinations of tags are specified,
    only the remaining intersection of tags is suggested
*/

const workspacePath = "/Users/jonathan/Downloads/docpile/workspace"

func main() {
	var registry = eventstore.NewRegistry(eventstore.PanicOnUnknownType())
	registry.Add("tag-added", events.TagAdded{})
	registry.Add("tag-removed", events.TagRenamed{})
	registry.Add("tag-synonym-defined", events.TagSynonymDefined{})
	registry.Add("tag-synonym-removed", events.TagSynonymRemoved{})
	registry.Add("managed-asset-imported", events.ManagedAssetImported{})
	registry.Add("cloud-asset-imported", events.CloudAssetImported{})
	registry.Add("document-defined", events.DocumentDefined{})

	aggregate := domain.NewAggregate(identity.NewEpochGenerator())
	store := eventstore.NewDelimitedText(
		storage.NewFileStorage(workspacePath, storage.Append(), storage.EnsureWorkspace()),
		registry,
		serialization.NewJSONSerializer())

	for message := range store.Load() {
		aggregate.Apply(message)
	}

	var applicator generic.Applicator = &Applicator{}
	applicator = eventstore.NewApplicator(applicator, store)

	var handler generic.Handler = handlers.NewDomainHandler(aggregate, applicator)
	handler = domain.NewWriteAssetHandler(handler, storage.NewFileStorage(workspacePath))

	tagController := http.NewTagWriteController(handler)
	assetController := http.NewAssetWriteController(handler)
	documentController := http.NewDocumentWriteController(handler)

	router := buildRouter()
	router.Handler("PUT", "/tags", detour.New(tagController.Add))
	router.Handler("POST", "/tags/name", detour.New(tagController.Rename))
	router.Handler("PUT", "/tags/synonym", detour.New(tagController.DefineSynonym))
	router.Handler("DELETE", "/tags/synonym", detour.New(tagController.RemoveSynonym))

	router.Handler("PUT", "/assets", detour.New(assetController.ImportManaged))
	router.Handler("PUT", "/documents", detour.New(documentController.Define))

	// GET /search/documents = document search criteria
	// GET /search/tags = tag auto-complete search

	// apply/remove one or more tags to a single document
	//   PUT /documents/:id/tags
	//   DELETE /document/:id/tags/:tags
	// apply/remove one tag to one or more documents
	//   PUT /tags/:id/documents
	//   DELETE /tags/:id/documents/:documents

	fmt.Println("Listening...")
	stdhttp.ListenAndServe("127.0.0.1:8080", router)
}

func buildRouter() *httprouter.Router {
	router := httprouter.New()
	router.HandleMethodNotAllowed = true
	router.HandleOPTIONS = true
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false
	return router
}

type Applicator struct{}

func (this *Applicator) Apply(messages []interface{}) {
	for _, message := range messages {
		log.Println("Applying:", message)
	}
}
