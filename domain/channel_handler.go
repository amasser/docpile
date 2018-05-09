package domain

type ChannelHandler struct {
	inner   Handler
	channel chan *channelEnvelope
}

func NewChannelHandler(inner Handler) *ChannelHandler {
	return &ChannelHandler{
		inner:   inner,
		channel: make(chan *channelEnvelope, 1024),
	}
}

func (this *ChannelHandler) Start() *ChannelHandler {
	go this.Listen()
	return this
}

func (this *ChannelHandler) Handle(message interface{}) (uint64, error) {
	envelope := newEnvelope(message)
	this.channel <- envelope
	return envelope.Result()
}

func (this *ChannelHandler) Listen() {
	for envelope := range this.channel {
		envelope.SetResult(this.inner.Handle(envelope.Message))
	}
}
