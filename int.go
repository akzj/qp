package qp

import (
	"strconv"
)

type Int int64

func (i Int) Invoke() Expression {
	return i
}

func (Int) GetType() Type {
	return IntType
}

func (i Int) String() string {
	return strconv.FormatInt(int64(i), 10)
}
