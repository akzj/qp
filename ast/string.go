package ast

import (
	"gitlab.com/akzj/qp/lexer"
	"gitlab.com/akzj/qp/runtime"
)

type String string

func (s String) GetObject(label string) *runtime.Object {
	return runtime.StringFunctions[label]
}

func (s String) AllocObject(label string) *runtime.Object {
	return runtime.StringFunctions[label]
}

func (s String) AddObject(k string, v *runtime.Object) {
	panic("implement me")
}

func (s String) String() string {
	return "\"" + string(s) + "\""
}

func (s String) Invoke() runtime.Invokable {
	return s
}

func (s String) Clone() BaseObject {
	return String(string(s))
}

func (s String) GetType() lexer.Type {
	return lexer.StringType
}
