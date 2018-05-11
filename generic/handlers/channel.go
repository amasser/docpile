package handlers

import "sync"

type Channel struct {
	inner   Handler
	channel chan *envelope
}

func NewChannel(inner Handler, options ...ChannelOption) *Channel {
	this := &Channel{inner: inner, channel: make(chan *envelope, 1024)}

	for _, option := range options {
		option(this)
	}

	return this
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

////////////////////////////////////////////////////////

type ChannelOption func(*Channel)

func Start() ChannelOption { return func(this *Channel) { go this.Listen() } }

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
