package eventstore

import (
	"errors"
	"reflect"
	"strings"
)

type Registry struct {
	nameToType map[string]reflect.Type
	typeToName map[reflect.Type]string
	panic      bool
}

func NewRegistry() *Registry {
	return &Registry{
		nameToType: map[string]reflect.Type{},
		typeToName: map[reflect.Type]string{},
	}
}
func (this *Registry) PanicWhenNotFound() *Registry {
	this.panic = true
	return this
}

func (this *Registry) Add(typeName string, instance interface{}) {
	typeName = strings.TrimSpace(typeName)
	instanceType := reflect.TypeOf(instance)
	this.nameToType[typeName] = instanceType
	this.typeToName[instanceType] = typeName
}

func (this *Registry) Name(registeredType reflect.Type) (string, error) {
	if typeName, contains := this.typeToName[registeredType]; contains {
		return typeName, nil
	} else if this.panic {
		panic(typeNotFound)
	} else {
		return "", typeNotFound
	}
}
func (this *Registry) Type(typeName string) (reflect.Type, error) {
	if registeredType, contains := this.nameToType[typeName]; contains {
		return registeredType, nil
	} else if this.panic {
		panic(typeNotFound)
	} else {
		return nil, typeNotFound
	}
}

var typeNotFound = errors.New("requested type not found")
