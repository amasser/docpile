package storage

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"reflect"
)

type TextEventStore struct {
	filename   string
	store      ReadWriter
	registry   Registry
	serializer Serializer
}

func NewTextEventStore(store ReadWriter, registry Registry, serializer Serializer) *TextEventStore {
	return &TextEventStore{
		filename:   defaultFilename,
		store:      store,
		registry:   registry,
		serializer: serializer,
	}
}

func (this *TextEventStore) Store(messages []interface{}) error {
	buffer := bytes.NewBuffer([]byte{})
	this.writeMessagesToBuffer(messages, buffer)
	return this.store.Write(this.filename, ioutil.NopCloser(buffer))
}
func (this *TextEventStore) writeMessagesToBuffer(messages []interface{}, destination *bytes.Buffer) {
	for _, message := range messages {
		this.writeMessageToBuffer(message, destination)
	}
}
func (this *TextEventStore) writeMessageToBuffer(message interface{}, destination *bytes.Buffer) {
	destination.WriteString(this.typeName(message))
	destination.WriteString(fieldDelimiter)
	this.serializer.Serialize(message, destination)
	destination.WriteString(lineBreak)
}
func (this *TextEventStore) typeName(message interface{}) string {
	if typeName, err := this.registry.Name(reflect.TypeOf(message)); err == nil {
		return typeName
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
	reader, err := this.store.Read(this.filename)
	if err != nil && err == NotFoundError {
		close(channel)
		return
	} else if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			panic(err)
		}
		channel <- this.parseLine(scanner.Bytes())
	}

	close(channel)
}
func (this *TextEventStore) parseLine(line []byte) interface{} {
	index := bytes.Index(line, []byte(fieldDelimiterBytes))
	if index < 0 {
		log.Panic(missingDelimiterError)
	}

	messageType := string(line[0:index])
	return this.deserialize(messageType, line[index:])
}
func (this *TextEventStore) deserialize(messageType string, body []byte) interface{} {
	instance := this.createInstance(messageType)
	if err := json.Unmarshal(body, instance.Interface()); err != nil {
		panic(err)
	}

	return instance.Elem().Interface()
}
func (this *TextEventStore) createInstance(name string) reflect.Value {
	if messageType, err := this.registry.Type(name); err == nil {
		return reflect.New(messageType)
	} else {
		panic(err)
	}
}

const (
	fieldDelimiter  = "\t"
	lineBreak       = "\n"
	defaultFilename = "events.txt"
)

var (
	fieldDelimiterBytes   = []byte(fieldDelimiter)
	missingDelimiterError = errors.New("missing field delimiter")
)
