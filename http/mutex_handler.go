package http

import (
	"net/http"

	"bitbucket.org/jonathanoliver/docpile/domain"
)

type MutexHandler struct {
	mutex domain.Mutex
	inner http.Handler
}

func NewMutexHandler(mutex domain.Mutex, inner http.Handler) *MutexHandler {
	return &MutexHandler{mutex: mutex, inner: inner}
}

func (this *MutexHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	defer this.mutex.Unlock()
	this.mutex.Lock()
	this.inner.ServeHTTP(response, request)
}
