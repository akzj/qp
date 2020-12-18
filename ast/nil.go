package ast

import (
	"gitlab.com/akzj/qp/lexer"
	"gitlab.com/akzj/qp/runtime"
)

type NilObject struct {
}

var NilObj = NilObject{}

func (NilObject) String() string {
	return "nil"
}

func (n NilObject) Invoke() runtime.Invokable {
	return n
}

func (NilObject) GetType() lexer.Type {
	return lexer.NilType
}
