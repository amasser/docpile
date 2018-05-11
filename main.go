package main

import (
	"fmt"
	"net/http"
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

func main() {
	const workspacePath = "/Users/jonathan/Downloads/docpile/workspace"
	wireup := NewWireup(workspacePath, workspacePath)

	aggregate := wireup.BuildDomain()
	store := wireup.BuildEventStore(aggregate)

	for message := range store.Load() {
		aggregate.Apply(message)
	}

	application := wireup.BuildMessageHandler(aggregate, store)
	httpHandler := wireup.BuildHTTPHandler(application)

	fmt.Println("Listening...")
	http.ListenAndServe("127.0.0.1:8080", httpHandler)
}
