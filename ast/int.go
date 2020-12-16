package qp

import (
	"gitlab.com/akzj/qp/ast"
	"gitlab.com/akzj/qp/lexer"
	"strconv"
)

type Int int64

func (i Int) Invoke() ast.Expression {
	return i
}

func (Int) GetType() lexer.Type {
	return lexer.IntType
}

func (i Int) String() string {
	return strconv.FormatInt(int64(i), 10)
}
