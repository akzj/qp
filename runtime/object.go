package runtime

import (
	"gitlab.com/akzj/qp/lexer"
	"reflect"
)

type Object struct {
	Pointer Invokable
	Label   string
}

func (obj *Object) Invoke() Invokable {
	switch obj.Pointer.(type) {
	case Invokable:
		return obj.Pointer.(Invokable).Invoke()
	default:
		panic(reflect.TypeOf(obj.Pointer).String())
	}
}
func (obj *Object) isNil() bool {
	return obj.Pointer == nil
}
func (obj *Object) GetType() lexer.Type {
	return lexer.ObjectType
}

func (obj *Object) String() string {
	return obj.Pointer.String()
}
