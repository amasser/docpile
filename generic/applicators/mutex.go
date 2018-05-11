package applicators

type Mutex struct {
	inner Applicator
	mutex Locker
}

func NewMutex(inner Applicator, mutex Locker) *Mutex {
	return &Mutex{inner: inner, mutex: mutex}
}

func (this *Mutex) Apply(messages []interface{}) {
	defer this.mutex.Unlock()
	this.mutex.Lock()
	this.inner.Apply(messages)
}
