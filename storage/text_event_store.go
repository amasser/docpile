package storage

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"reflect"
)

type TextEventStore struct {
	filename string
	store    ReadWriter
}

func NewTextEventStore(store ReadWriter) *TextEventStore {
	return &TextEventStore{
		filename: "events.txt",
		store:    store,
	}
}

func (this *TextEventStore) Store(messages []interface{}) error {
	buffer := bytes.NewBuffer([]byte{})
	writeToBuffer(buffer, messages)
	return this.store.Write(this.filename, ioutil.NopCloser(buffer))
}
func writeToBuffer(buffer *bytes.Buffer, messages []interface{}) {
	for _, message := range messages {
		buffer.WriteString(reflect.TypeOf(message).Name())
		buffer.WriteString("\t")
		buffer.WriteString(serialize(message))
		buffer.WriteString("\n")
	}
}
func serialize(message interface{}) string {
	if serialized, err := json.Marshal(message); err == nil {
		return string(serialized)
	} else {
		panic(err)
	}
}

func (this *TextEventStore) Load() <-chan interface{} {
	output := make(chan interface{}, 1024)
	go this.load(output)
	return output
}
func (this *TextEventStore) load(channel chan<- interface{}) {
	// load the file
	// iterate/deserialize
	// close the channel when done
}
