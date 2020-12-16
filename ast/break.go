package ast

import (
	"gitlab.com/akzj/qp/lexer"
)

var BreakObj = &BreakObject{}

type BreakObject struct {
}

func (b *BreakObject) String() string {
	return "break"
}

func (b *BreakObject) Invoke() Expression {
	return b
}

func (b *BreakObject) GetType() lexer.Type {
	return lexer.BreakType
}
