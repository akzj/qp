package qp

import (
	"fmt"
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

func (s *StringObject) invoke() (Expression, error) {
	if s.init {
		return s, nil
	}
	s.objects = StringObjectBuiltInFunctionMap
	return s, nil
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

var StringObjectBuiltInFunctionMap = map[string]*Object{
	"to_lower": &Object{
		inner: &stringLowCase{},
		label: "to_lower",
	},
	"clone": &Object{
		inner: &StringObjectClone{},
		label: "clone",
	},
}

type StringObjectClone struct {
}

func (s StringObjectClone) invoke(arguments ...Expression) (Expression, error) {
	fmt.Println("StringObjectClone")
	if len(arguments) > 1 {
		panic("only one arguments")
	}
	expression, err := arguments[0].invoke()
	if err != nil {
		return nil, err
	}
	for {
		fmt.Println(reflect.TypeOf(expression).String())
		switch inner := expression.(type) {
		case *StringObject:
			return inner.clone(), nil
		case *Object:
			expression, err = inner.invoke()
			if err != nil {
				log.Println(err)
				return nil, err
			}
		case *ReturnStatement:
			expression = inner
		default:
			log.Println("type error", reflect.TypeOf(expression).String())
			return nil, fmt.Errorf("no string type error")
		}
	}
}

type stringLowCase struct{}

func (stringLowCase) invoke(arguments ...Expression) (Expression, error) {
	fmt.Println("stringLowCase")
	if len(arguments) > 1 {
		panic("only one arguments")
	}
	expression, err := arguments[0].invoke()
	if err != nil {
		return nil, err
	}
	for {
		fmt.Println(reflect.TypeOf(expression).String())
		switch inner := expression.(type) {
		case *StringObject:
			inner.data = strings.ToLower(inner.data)
			return inner, nil
		case *Object:
			expression, err = inner.invoke()
			if err != nil {
				log.Println(err)
				return nil, err
			}
		case *ReturnStatement:
			expression = inner
		default:
			log.Println("type error", reflect.TypeOf(expression).String())
			return nil, fmt.Errorf("no string type error")
		}
	}
}
