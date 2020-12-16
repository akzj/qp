package ast

import (
	"gitlab.com/akzj/qp/lexer"
	"reflect"
)

type Object struct {
	Inner   Expression
	Pointer int
	Label   string
	Typ     lexer.Type
}

func (obj *Object) Invoke() Expression {
	switch obj.Inner.(type) {
	case Expression:
		return obj.Inner.(Expression).Invoke()
	default:
		panic(reflect.TypeOf(obj.Inner).String())
	}
}
func (obj *Object) isNil() bool {
	return obj.Inner == nil
}
func (obj *Object) GetType() lexer.Type {
	return lexer.ObjectType
}

func (obj *Object) String() string {
	return obj.Inner.String()
}

func (obj *Object) InitType() {
	switch obj.Inner.(type) {
	case *Int:
		obj.Typ = lexer.IntType
	}
}

func (obj *Object) UnwrapFunction() Function {
	var object = obj
	if object.Inner == nil {
		panic(object.Inner)
	}
	if function, ok := object.Inner.(Function); ok {
		return function
	}
	panic("unknown type" + reflect.TypeOf(object.Inner).String())
}
