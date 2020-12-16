package qp

import (
	"gitlab.com/akzj/qp/ast"
	"gitlab.com/akzj/qp/lexer"
)

type NilObject struct {
}

var nilObject = NilObject{}

func (NilObject) String() string {
	return "nil"
}

func (n NilObject) Invoke() ast.Expression {
	return n
}

func (NilObject) GetType() lexer.Type {
	return lexer.NilType
}
