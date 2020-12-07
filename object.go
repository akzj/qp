package qp

import (
	"fmt"
	"reflect"
	"strconv"
)

type Object struct {
	inner   interface{}
	pointer int
	label   string
	typ     Type
}

var breakObject = &BreakObject{}

type BreakObject struct {
}

type IntObject struct {
	val int64
}

type BoolObject struct {
	val bool
}

type StructObject struct {
	vm    *VMContext
	label string
	//init statement when create object
	init          bool
	initStatement Statements
	//user define function
	object map[string]*Object
}

func (sObj *StructObject) invoke() (Expression, error) {
	if sObj.init {
		return sObj, nil
	}
	for _, statement := range sObj.initStatement {
		if _, err := statement.invoke(); err != nil {
			return nil, err
		}
	}
	return sObj, nil
}

func (sObj *StructObject) getType() Type {
	return StructObjectType
}

func (sObj *StructObject) getObject(label string) *Object {
	object, ok := sObj.object[label]
	if ok {
		return object
	}
	return nil
}

func (sObj *StructObject) allocObject(label string) *Object {
	if sObj.object == nil {
		sObj.object = map[string]*Object{}
	}
	object, ok := sObj.object[label]
	if ok {
		return object
	} else {
		object = &Object{label: label}
		sObj.object[label] = object
	}
	return object
}

func (sObj *StructObject) clone() *StructObject {
	clone := *sObj
	for k, v := range sObj.object {
		clone.addObject(k, v)
	}
	if len(sObj.initStatement) != 0 {
		clone.initStatement = make(Statements, len(sObj.initStatement))
		copy(clone.initStatement, sObj.initStatement)
	}
	return &clone
}

func (sObj *StructObject) addObject(k string, v *Object) {
	if sObj.object == nil {
		sObj.object = map[string]*Object{}
	}
	sObj.object[k] = v
}

func (b BreakObject) invoke() (Expression, error) {
	return b, nil
}

func (b BreakObject) getType() Type {
	return breakTokenType
}

func (i *IntObject) invoke() (Expression, error) {
	fmt.Println("IntObject,invoke", i.val)
	return i, nil
}

func (i *IntObject) getType() Type {
	return IntObjectType
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

func (obj *Object) invoke() (Expression, error) {
	if obj == nil{
		panic("")
	}
	if obj.inner == nil{
		fmt.Println(reflect.TypeOf(obj).String())
		panic("obj.inner")
	}
	fmt.Println("inner",reflect.TypeOf(obj.inner).String())
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
