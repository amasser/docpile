package storage

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"reflect"
)

type CSVEventStore struct {
	filename string
	store    ReadWriter
}

func NewCSVReader(store ReadWriter) *CSVEventStore {
	return &CSVEventStore{
		filename: "events.csv",
		store:    store,
	}
}

func (this *CSVEventStore) Store(messages ...interface{}) error {
	buffer := bytes.NewBuffer([]byte{})
	writeToBuffer(buffer, messages)
	return this.store.Write(this.filename, ioutil.NopCloser(buffer))
}
func writeToBuffer(buffer *bytes.Buffer, messages []interface{}) {
	writer := csv.NewWriter(buffer)
	writer.Comma = '\t'
	writer.UseCRLF = false
	for _, message := range messages {
		writer.Write([]string{
			reflect.TypeOf(message).String(),
			serialize(message),
		})
	}
	writer.Flush()
}
func serialize(message interface{}) string {
	if serialized, err := json.Marshal(message); err == nil {
		return string(serialized)
	} else {
		panic(err)
	}
}

func (this *CSVEventStore) Load() <-chan interface{} {
	output := make(chan interface{}, 1024)
	go this.load(output)
	return output
}
func (this *CSVEventStore) load(channel chan<- interface{}) {
	// load the file
	// iterate/deserialize
	// close the channel when done
}
