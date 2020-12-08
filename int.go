package qp

import (
	"strconv"
)

type IntObject struct {
	val int64
}

func (i *IntObject) Invoke() Expression {
	return i
}

func (i *IntObject) getType() Type {
	return IntObjectType
}

func (i *IntObject) String() string {
	return strconv.FormatInt(i.val, 10)
}
