package ast

import (
	"gitlab.com/akzj/qp/lexer"
	"gitlab.com/akzj/qp/runtime"
	"strconv"
)

type Int int64

func (i Int) Invoke() runtime.Invokable {
	return i
}

func (Int) GetType() lexer.Type {
	return lexer.IntType
}

func (i Int) String() string {
	return strconv.FormatInt(int64(i), 10)
}
