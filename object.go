package qp

import (
	"fmt"
	"log"
	"reflect"
)

type Object struct {
	inner   interface{}
	pointer int
	label   string
	typ     Type
}

func (obj *Object) invoke() (Expression, error) {
	switch inner := obj.inner.(type) {
	case *FuncStatement:
		if err := inner.doClosureInit(); err != nil {
			log.Println("function statement do function closure init failed")
			return nil, err
		}
		return obj, nil
	case Expression:
		return obj.inner.(Expression).invoke()
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
	case int64:
		obj.typ = IntObjectType
	}
}

func (obj *Object) unwrapFunction() Function {
	var object = obj
	for object != nil {
		switch inner := object.inner.(type) {
		case *FuncStatement:
			return inner
		case Function:
			return inner
		default:
			panic("unknown type" + reflect.TypeOf(inner).String())
		}
	}
	return nil
}
