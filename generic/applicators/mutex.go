package applicators

type Mutex struct {
	inner Applicator
	mutex locker
}

func NewMutex(inner Applicator, mutex locker) *Mutex {
	return &Mutex{inner: inner, mutex: mutex}
}

func (this *Mutex) Apply(messages []interface{}) {
	defer this.mutex.Unlock()
	this.mutex.Lock()
	this.inner.Apply(messages)
}
