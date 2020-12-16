package ast

import (
	"gitlab.com/akzj/qp/lexer"
)

type NilObject struct {
}

var NilObj = NilObject{}

func (NilObject) String() string {
	return "nil"
}

func (n NilObject) Invoke() Expression {
	return n
}

func (NilObject) GetType() lexer.Type {
	return lexer.NilType
}
