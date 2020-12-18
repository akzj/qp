package ast

import (
	"gitlab.com/akzj/qp/lexer"
	"gitlab.com/akzj/qp/runtime"
)

var BreakObj = &BreakObject{}

type BreakObject struct {
}

func (b *BreakObject) String() string {
	return "break"
}

func (b *BreakObject) Invoke() runtime.Invokable {
	return b
}

func (b *BreakObject) GetType() lexer.Type {
	return lexer.BreakType
}
