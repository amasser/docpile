package web

import "net/http"

type LockHandler struct {
	mutex locker
	inner http.Handler
}

func NewLockHandler(mutex locker, inner http.Handler) *LockHandler {
	return &LockHandler{mutex: mutex, inner: inner}
}

func (this *LockHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	defer this.mutex.Unlock()
	this.mutex.Lock()
	this.inner.ServeHTTP(response, request)
}
