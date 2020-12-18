package ast

import (
	"gitlab.com/akzj/qp/lexer"
	"gitlab.com/akzj/qp/runtime"
)

type Bool bool

var TrueObject = Bool(true)
var FalseObject = Bool(false)

func (b Bool) String() string {
	if b {
		return "true"
	}
	return "false"
}

func (b Bool) Invoke() runtime.Invokable {
	return b
}

func (b Bool) GetType() lexer.Type {
	return lexer.BoolObjectType
}
