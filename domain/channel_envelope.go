package domain

import "sync"

type channelEnvelope struct {
	Message interface{}
	waiter  *sync.WaitGroup
	id      uint64
	err     error
}

func newEnvelope(message interface{}) *channelEnvelope {
	this := &channelEnvelope{waiter: &sync.WaitGroup{}, Message: message}
	this.waiter.Add(1)
	return this
}
func (this *channelEnvelope) SetResult(id uint64, err error) {
	this.id = id
	this.err = err
	this.waiter.Done()
}
func (this *channelEnvelope) Result() (uint64, error) {
	this.waiter.Wait()
	return this.id, this.err
}
