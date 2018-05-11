package eventstore

import "reflect"

type Registry interface {
	Name(reflect.Type) (string, error)
	Type(string) (reflect.Type, error)
}
