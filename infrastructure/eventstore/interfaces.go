package eventstore

import "reflect"

type TypeRegistry interface {
	Name(reflect.Type) (string, error)
	Type(string) (reflect.Type, error)
}

type EventStore interface {
	Store([]interface{}) error
}
