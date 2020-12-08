package qp

import (
	"fmt"
	"reflect"
)

type Object struct {
	inner   interface{}
	pointer int
	label   string
	typ     Type
}

func (obj *Object) invoke() (Expression, error) {
	if obj == nil {
		panic("")
	}
	if obj.inner == nil {
		fmt.Println(reflect.TypeOf(obj).String())
		panic("obj.inner")
	}
	if obj.inner == obj{
		panic("dead loop")
	}
	fmt.Println("inner", reflect.TypeOf(obj.inner).String())
	switch inner := obj.inner.(type) {
	case *FuncStatement:
		if err := inner.doClosureInit(); err != nil {
			fmt.Println("function statement do function closure init failed")
			return nil, err
		}
		return obj, nil
	case Expression:
		return obj.inner.(Expression), nil
	default:
		panic(reflect.TypeOf(obj.inner).String())
	}
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

func (obj *Object) unwrapFunction() Function {
	fmt.Println("unwrapFunction start")
	var object = obj
	for object != nil {
		if object.inner == object {
			panic("objects inner pointer to self ,dead loop")
		}
		if object.inner == nil {
			panic("objects.inner nil")
		}
		switch inner := object.inner.(type) {
		case *Object:
			object = inner
			continue
		case *FuncStatement:
			return inner
		case Function:
			return inner
		case *ReturnStatement:
			switch inner2 := inner.returnVal.(type) {
			case *Object:
				object = inner2
				continue
			default:
				panic("unknown type" + reflect.TypeOf(inner.returnVal).String())
			}
		default:
			panic("unknown type" + reflect.TypeOf(inner).String())
		}
	}
	return nil
}
