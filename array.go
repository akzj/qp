package qp

import (
	"fmt"
	"log"
)

type Array struct {
	data   []Expression
	object map[string]*Object
}

func (a *Array) String() string {
	return fmt.Sprintf("%+v", a.data)
}

func (a *Array) Invoke() Expression {
	return a
}

func (a *Array) getType() Type {
	return arrayObjectType
}

func (a *Array) getObject(label string) *Object {
	return arrayBuiltInFunctions[label]
}

func (a *Array) allocObject(label string) *Object {
	if obj := a.getObject(label); obj != nil {
		return obj
	}
	if a.object == nil {
		a.object = map[string]*Object{}
	}
	obj := &Object{
		inner: nilObject,
		label: label,
		typ:   nilType,
	}
	a.object[label] = obj
	return obj
}

func (a *Array) clone() BaseObject {
	var data = make([]Expression, len(a.data))
	copy(data, a.data)
	return &Array{data: data}
}

func registerArrayFunction() {
	registerBuiltInFunc(arrayBuiltInFunctions, "append", func(arguments ...Expression) Expression {
		array := arguments[0].Invoke().(*Array)
		for _, exp := range arguments[1:] {
			array.data = append(array.data, exp.Invoke())
		}
		return array
	})("size", func(arguments ...Expression) Expression {
		return Int(len(arguments[0].(*Array).data))
	})("get", func(arguments ...Expression) Expression {
		if len(arguments) != 2 {
			log.Panic("array get() arguments error")
		}
		array, ok := arguments[0].(*Array)
		if ok == false {
			log.Panic("object not array type")
		}
		i, ok := arguments[1].(Int)
		if ok == false {
			log.Panic("is not array arguments error")
		}
		if len(array.data) <= int(i) {
			log.Panic("index out of range")
		}
		return array.data[i]
	})
}
