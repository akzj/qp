package ast

import (
	"fmt"
	"gitlab.com/akzj/qp/lexer"
)

type Array struct {
	Data   []Expression
	Object map[string]*Object
}

func (a *Array) String() string {
	return fmt.Sprintf("%+v", a.Data)
}

func (a *Array) Invoke() Expression {
	return a
}

func (a *Array) GetType() lexer.Type {
	return lexer.ArrayObjectType
}

func (a *Array) GetObject(label string) *Object {
	return ArrayFunctions[label]
}

func (a *Array) AllocObject(label string) *Object {
	if obj := a.GetObject(label); obj != nil {
		return obj
	}
	if a.Object == nil {
		a.Object = map[string]*Object{}
	}
	obj := &Object{
		Inner: NilObj,
		Label: label,
		Typ:   lexer.NilType,
	}
	a.Object[label] = obj
	return obj
}

func (a *Array) Clone() BaseObject {
	var data = make([]Expression, len(a.Data))
	copy(data, a.Data)
	return &Array{Data: data}
}