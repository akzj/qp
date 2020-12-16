package ast

import (
	"gitlab.com/akzj/qp/lexer"
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

func (b Bool) Invoke() Expression {
	return b
}

func (b Bool) GetType() lexer.Type {
	return lexer.BoolObjectType
}
