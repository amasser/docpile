package applicators

type Channel struct {
	channel chan []interface{}
	inner   Applicator
}

func NewChannel(inner Applicator) *Channel {
	return &Channel{inner: inner, channel: make(chan []interface{}, 1024)}
}

func (this *Channel) Listen() {
	for messages := range this.channel {
		this.inner.Apply(messages)
	}
}

func (this *Channel) Apply(messages []interface{}) {
	this.channel <- messages
}

func (this *Channel) Close() {
	close(this.channel)
}
