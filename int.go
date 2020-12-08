package qp

import (
	"strconv"
)

type IntObject struct {
	val int64
}

func (i *IntObject) invoke() (Expression, error) {
	return i, nil
}

func (i *IntObject) getType() Type {
	return IntObjectType
}

func (i *IntObject) String() string {
	return strconv.FormatInt(i.val, 10)
}
