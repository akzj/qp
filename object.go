package qp

import (
	"fmt"
	"strconv"
)

const (
	ObjectType     Type = 100000
	IntObjectType  Type = 10000
	BoolObjectType Type = 10001
)

var breakObject = &BreakObject{}

type BreakObject struct {
}

func (b BreakObject) invoke() (Expression, error) {
	return b, nil
}

func (b BreakObject) getType() Type {
	return breakTokenType
}

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

type BoolObject struct {
	val bool
}

func (b *BoolObject) invoke() (Expression, error) {
	return b, nil
}

func (b *BoolObject) getType() Type {
	return BoolObjectType
}

func (i *IntObject) String() string {
	return strconv.FormatInt(i.val, 10)
}

type Object struct {
	inner   interface{}
	pointer int
	label   string
	typ     Type
}

func (obj *Object) invoke() (Expression, error) {
	return obj.inner.(Expression), nil
}

func (obj *Object) getType() Type {
	return ObjectType
}

func (obj *Object) String() string {
	return fmt.Sprintf("pointer %d label %s", obj.pointer, obj.label)
}

func (obj *Object) initType() {
	switch obj.inner.(type) {
	case int64:
		obj.typ = IntObjectType
	}
}

func (obj *Object) AddObject(val *Object) (Expression, error) {
	switch obj.typ {
	case IntObjectType:
		switch val.typ {
		case IntObjectType:
			return &IntObject{val: obj.inner.(*IntObject).val + val.inner.(*IntObject).val}, nil
		}
	}
	return nil, fmt.Errorf("unknown obj")
}
