package qp

import (
	"fmt"
	"strconv"
)

type IntObject struct {
	val int64
}

func (i *IntObject) invoke() (Expression, error) {
	fmt.Println("IntObject,invoke", i.val)
	return i, nil
}

func (i *IntObject) getType() Type {
	return IntObjectType
}

func (i *IntObject) String() string {
	return strconv.FormatInt(i.val, 10)
}
