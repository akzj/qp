package qp

import (
	"log"
	"reflect"
	"strings"
)

type String string

func (s String) getObject(label string) *Object {
	return stringBuiltInFunctions[label]
}

func (s String) allocObject(label string) *Object {
	return stringBuiltInFunctions[label]
}

func (s String) addObject(k string, v *Object) {
	panic("implement me")
}

func (s String) String() string {
	return string(s)
}

func (s String) Invoke() Expression {
	return s
}

func (s String) clone() BaseObject {
	return String(string(s))
}

func (s String) getType() Type {
	return stringType
}

type StringObjectClone struct {
}

func (s StringObjectClone) String() string {
	panic("implement me")
}

func (s StringObjectClone) Invoke() Expression {
	return s
}

func (s StringObjectClone) getType() Type {
	return TypeObjectType
}

func (s StringObjectClone) call(arguments ...Expression) Expression {
	if len(arguments) > 1 {
		log.Panicln("only one arguments")
	}

	for {
		switch inner := arguments[0].(type) {
		case String:
			return inner.clone()
		default:
			log.Panicln("type error", reflect.TypeOf(arguments[0]).String())
		}
	}
}

type stringLowCase struct{}

func (c stringLowCase) String() string {
	panic("implement me")
}

func (c stringLowCase) Invoke() Expression {
	return c
}

func (c stringLowCase) getType() Type {
	return TypeObjectType
}

func (stringLowCase) call(arguments ...Expression) Expression {
	if len(arguments) > 1 {
		log.Panicln("only one arguments")
	}
	for {
		switch inner := arguments[0].(type) {
		case String:
			inner = String(strings.ToLower(string(inner)))
			return inner
		default:
			log.Panicln("type error", reflect.TypeOf(arguments[0]).String())
		}
	}
}
