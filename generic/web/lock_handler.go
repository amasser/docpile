package web

import (
	"net/http"
	"sync"
)

type LockHandler struct {
	mutex sync.Locker
	inner http.Handler
}

func NewLockHandler(mutex sync.Locker, inner http.Handler) *LockHandler {
	return &LockHandler{mutex: mutex, inner: inner}
}

func (this *LockHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	defer this.mutex.Unlock()
	this.mutex.Lock()
	this.inner.ServeHTTP(response, request)
}
