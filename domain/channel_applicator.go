package domain

type ChannelApplicator struct {
	channel chan []interface{}
	inner   Applicator
}

func NewChannelApplicator(inner Applicator) *ChannelApplicator {
	return &ChannelApplicator{
		inner:   inner,
		channel: make(chan []interface{}, 1024),
	}
}

func (this *ChannelApplicator) Start() *ChannelApplicator {
	go this.Listen()
	return this
}

func (this *ChannelApplicator) Listen() {
	for messages := range this.channel {
		this.inner.Apply(messages...)
	}
}

func (this *ChannelApplicator) Apply(messages ...interface{}) {
	this.channel <- messages
}
