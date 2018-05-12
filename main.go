package main

import (
	"fmt"
	"net/http"
)

/*
TODOs
- projections
- better/safer text encoding of event store OR disable advanced characters in input text fields, e.g. tags, description etc.
  or event per file?

- apply/remove tags from documents
- delete document
- update document description/published date/effective dates

- document search
- tag search (finding tags that match partial text)
- during document search, as combinations of tags are specified,
    only the remaining intersection of tags is suggested
  e.g. /search/tags?text=asdf&tag=123&tag=456&tag=789 (list of available tags)
*/

func main() {
	const workspacePath = "/Users/jonathan/Downloads/docpile/workspace"
	wireup := NewWireup(workspacePath, workspacePath)

	aggregate := wireup.BuildDomain()
	store := wireup.BuildEventStore(aggregate)
	projector := wireup.BuildProjector()

	for message := range store.Load() {
		aggregate.Apply([]interface{}{message})
		projector.Apply([]interface{}{message})
	}

	application := wireup.BuildMessageHandler(aggregate, store, projector)
	httpHandler := wireup.BuildHTTPHandler(application, projector)

	fmt.Println("Listening...")
	http.ListenAndServe("127.0.0.1:8080", httpHandler)
}
