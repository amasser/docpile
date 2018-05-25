package handlers

import "sync"

type Channel struct {
	inner   Handler
	channel chan *envelope
}

func NewChannel(inner Handler) *Channel {
	return &Channel{inner: inner, channel: make(chan *envelope, 1024)}
}

func (this *Channel) Handle(message interface{}) Result {
	envelope := newEnvelope(message)
	this.channel <- envelope
	return envelope.Result()
}

func (this *Channel) Listen() {
	for envelope := range this.channel {
		envelope.SetResult(this.inner.Handle(envelope.Message))
	}
}

func (this *Channel) Close() {
	close(this.channel)
}

////////////////////////////////////////////////////////

type envelope struct {
	Message interface{}
	waiter  *sync.WaitGroup
	result  Result
}

func newEnvelope(message interface{}) *envelope {
	this := &envelope{waiter: &sync.WaitGroup{}, Message: message}
	this.waiter.Add(1)
	return this
}
func (this *envelope) SetResult(result Result) {
	this.result = result
	this.waiter.Done()
}
func (this *envelope) Result() Result {
	this.waiter.Wait()
	return this.result
}
