package main

import "log"

const (
	listenAddress = "127.0.0.1:8888"
	workspacePath = "./workspace"
)

func init() {
	log.SetFlags(log.Llongfile | log.Lmicroseconds)
}

func main() {
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
	wireup.BuildHTTPServer(listenAddress, httpHandler)

	listener := wireup.BuildListener()
	listener.Listen()
}
