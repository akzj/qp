package ast

import (
	"fmt"
	"gitlab.com/akzj/qp/lexer"
	"gitlab.com/akzj/qp/runtime"
)

type Array struct {
	Data   []runtime.Invokable
	Object map[string]*runtime.Object
}

func (a *Array) String() string {
	return fmt.Sprintf("%+v", a.Data)
}

func (a *Array) Invoke() runtime.Invokable {
	return a
}

func (a *Array) GetType() lexer.Type {
	return lexer.ArrayObjectType
}

func (a *Array) GetObject(label string) *runtime.Object {
	return runtime.ArrayFunctions[label]
}

func (a *Array) AllocObject(label string) *runtime.Object {
	if obj := a.GetObject(label); obj != nil {
		return obj
	}
	if a.Object == nil {
		a.Object = map[string]*runtime.Object{}
	}
	obj := &runtime.Object{
		Pointer: NilObj,
		Label:   label,
	}
	a.Object[label] = obj
	return obj
}

func (a *Array) Clone() BaseObject {
	var data = make([]runtime.Invokable, len(a.Data))
	copy(data, a.Data)
	return &Array{Data: data}
}