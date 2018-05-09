package main

import (
	"fmt"
	"log"

	"bitbucket.org/jonathanoliver/docpile/domain"
	"bitbucket.org/jonathanoliver/docpile/http"
	"bitbucket.org/jonathanoliver/docpile/storage"
	"github.com/julienschmidt/httprouter"
	"github.com/smartystreets/detour"

	stdhttp "net/http"
)

/*
TODOs
- loading of events from event storage at startup
- routing of events to event storage and projections
- projections

- tag "synonym" vs "alias" (which word is better)
- apply/remove tags from documents

- document search
- tag search (finding tags that match partial text)
- during document search, as combinations of tags are specified,
    only the remaining intersection of tags is suggested
*/

const workspacePath = "/Users/jonathan/Downloads/docpile"

func main() {
	identity := domain.NewEpochGenerator()
	aggregate := domain.NewAggregate(identity)

	applicator := &Applicator{} // TODO

	var handler domain.Handler = domain.NewMessageHandler(aggregate, applicator)
	handler = storage.NewLocalStorageHandler(handler, storage.NewLocalStorage(workspacePath))

	tagController := http.NewTagController(handler)
	assetController := http.NewAssetController(handler)
	documentController := http.NewDocumentController(handler)

	router := buildRouter()
	router.Handler("PUT", "/tags", detour.New(tagController.Add))
	router.Handler("POST", "/tags/:id", detour.New(tagController.Rename))
	router.Handler("PUT", "/tags/:id/synonyn", detour.New(tagController.DefineSynonym))
	router.Handler("DELETE", "/tags/:id/synonyn", detour.New(tagController.RemoveSynonym))
	router.Handler("PUT", "/assets", detour.New(assetController.ImportManaged))
	router.Handler("PUT", "/documents", detour.New(documentController.Define))

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

func (this *Applicator) Apply(messages ...interface{}) {
	for _, message := range messages {
		log.Println("Applying:", message)
	}
}
