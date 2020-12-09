package qp

import (
	"log"
	"reflect"
	"strings"
)

type StringObject struct {
	TypeObject
	data string
}

func (s *StringObject) String() string {
	return s.data
}

func (s *StringObject) Invoke() Expression {
	if s.init {
		return s
	}
	s.objects = stringBuiltInFunctions
	return s
}

func (s *StringObject) clone() BaseObject {
	return &StringObject{
		TypeObject: s.TypeObject,
		data:       s.data,
	}
}

func (s *StringObject) getType() Type {
	return stringTokenType
}

type StringObjectClone struct {
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
		case *StringObject:
			return inner.clone()
		default:
			log.Panicln("type error", reflect.TypeOf(arguments[0]).String())
		}
	}
}

type stringLowCase struct{}

func (c *stringLowCase) Invoke() Expression {
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
		case *StringObject:
			inner.data = strings.ToLower(inner.data)
			return inner
		default:
			log.Panicln("type error", reflect.TypeOf(arguments[0]).String())
		}
	}
}
