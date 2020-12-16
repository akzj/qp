package ast

import (
	"gitlab.com/akzj/qp/lexer"
)

type String string

func (s String) GetObject(label string) *Object {
	return StringFunctions[label]
}

func (s String) AllocObject(label string) *Object {
	return StringFunctions[label]
}

func (s String) AddObject(k string, v *Object) {
	panic("implement me")
}

func (s String) String() string {
	return string(s)
}

func (s String) Invoke() Expression {
	return s
}

func (s String) Clone() BaseObject {
	return String(string(s))
}

func (s String) GetType() lexer.Type {
	return lexer.StringType
}
