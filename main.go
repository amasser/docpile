package main

import (
	"github.com/smartystreets/httpx"
	"github.com/smartystreets/listeners"
)

const listenAddress = "127.0.0.1:8080"

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
	httpServer := httpx.NewHTTPServer(listenAddress, httpHandler)
	listener := listeners.NewCompositeWaitShutdownListener(httpServer)
	listener.Listen()
}
