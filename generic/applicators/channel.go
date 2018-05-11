package applicators

type Channel struct {
	channel chan []interface{}
	inner   Applicator
}

func NewChannel(inner Applicator, options ...ChannelOption) *Channel {
	this := &Channel{inner: inner, channel: make(chan []interface{}, 1024)}

	for _, option := range options {
		option(this)
	}

	return this
}

func (this *Channel) Listen() {
	for messages := range this.channel {
		this.inner.Apply(messages)
	}
}

func (this *Channel) Apply(messages []interface{}) {
	this.channel <- messages
}

////////////////////////////////////////////////////////

type ChannelOption func(*Channel)

func StartChannel() ChannelOption { return func(this *Channel) { go this.Listen() } }
