package events

import (
	"errors"
	"reflect"
	"strings"
)

type Registry struct {
	nameToType map[string]reflect.Type
	typeToName map[reflect.Type]string
}

func NewRegistry() *Registry {
	return &Registry{
		nameToType: map[string]reflect.Type{},
		typeToName: map[reflect.Type]string{},
	}
}

func (this *Registry) Register(typeName string, instance interface{}) {
	typeName = strings.TrimSpace(typeName)
	instanceType := reflect.TypeOf(instance)
	this.nameToType[typeName] = instanceType
	this.typeToName[instanceType] = typeName
}

func (this *Registry) Name(registeredType reflect.Type) (string, error) {
	if typeName, contains := this.typeToName[registeredType]; contains {
		return typeName, nil
	} else {
		return "", typeNotFound
	}
}
func (this *Registry) Type(typeName string) (reflect.Type, error) {
	if registeredType, contains := this.nameToType[typeName]; contains {
		return registeredType, nil
	} else {
		return nil, typeNotFound
	}
}

var typeNotFound = errors.New("requested type not found")
