package qp

import (
	"strconv"
)

type IntObject int64

func (i *IntObject) Invoke() Expression {
	return i
}

func (i *IntObject) getType() Type {
	return IntObjectType
}

func (i *IntObject) String() string {
	return strconv.FormatInt(int64(*i), 10)
}
