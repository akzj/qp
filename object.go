package qp

import (
	"fmt"
	"reflect"
)

type Object struct {
	inner   Expression
	pointer int
	label   string
	typ     Type
}

func (obj *Object) Invoke() Expression {
	switch obj.inner.(type) {
	case Expression:
		return obj.inner.(Expression).Invoke()
	default:
		panic(reflect.TypeOf(obj.inner).String())
	}
}
func (obj *Object) isNil() bool {
	return obj.inner == nil
}
func (obj *Object) getType() Type {
	return ObjectType
}

func (obj *Object) String() string {
	return fmt.Sprintf("pointer %d label %s", obj.pointer, obj.label)
}

func (obj *Object) initType() {
	switch obj.inner.(type) {
	case *Int:
		obj.typ = IntType
	}
}

func (obj *Object) unwrapFunction() Function {
	var object = obj
	if object.inner == nil {
		panic(object.inner)
	}
	if function, ok := object.inner.(Function); ok {
		return function
	}
	panic("unknown type" + reflect.TypeOf(object.inner).String())
}
