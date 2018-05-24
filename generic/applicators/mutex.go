package applicators

import "sync"

type Mutex struct {
	inner Applicator
	mutex sync.Locker
}

func NewMutex(inner Applicator, mutex sync.Locker) *Mutex {
	return &Mutex{inner: inner, mutex: mutex}
}

func (this *Mutex) Apply(messages []interface{}) {
	defer this.mutex.Unlock()
	this.mutex.Lock()
	this.inner.Apply(messages)
}
