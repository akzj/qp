package qp

import (
	"gitlab.com/akzj/qp/ast"
	"gitlab.com/akzj/qp/lexer"
)

type Bool bool

var trueObject = Bool(true)
var falseObject = Bool(false)

func (b Bool) String() string {
	if b {
		return "true"
	}
	return "false"
}

func (b Bool) Invoke() ast.Expression {
	return b
}

func (b Bool) GetType() lexer.Type {
	return lexer.BoolObjectType
}
