package ast

import (
	"fmt"
	"gitlab.com/akzj/qp"
	"gitlab.com/akzj/qp/lexer"
	"log"
)

type Array struct {
	data   []qp.Expression
	object map[string]*Object
}

func (a *Array) String() string {
	return fmt.Sprintf("%+v", a.data)
}

func (a *Array) Invoke() qp.Expression {
	return a
}

func (a *Array) GetType() lexer.Type {
	return lexer.ArrayObjectType
}

func (a *Array) getObject(label string) *Object {
	return qp.ArrayBuiltInFunctions[label]
}

func (a *Array) allocObject(label string) *Object {
	if obj := a.getObject(label); obj != nil {
		return obj
	}
	if a.object == nil {
		a.object = map[string]*Object{}
	}
	obj := &Object{
		inner: qp.nilObject,
		label: label,
		typ:   lexer.NilType,
	}
	a.object[label] = obj
	return obj
}

func (a *Array) clone() qp.BaseObject {
	var data = make([]qp.Expression, len(a.data))
	copy(data, a.data)
	return &Array{data: data}
}

func RegisterArrayFunction() {
	qp.registerBuiltInFunc(qp.ArrayBuiltInFunctions, "append", func(arguments ...qp.Expression) qp.Expression {
		array := arguments[0].Invoke().(*Array)
		for _, exp := range arguments[1:] {
			array.data = append(array.data, exp.Invoke())
		}
		return array
	})("size", func(arguments ...qp.Expression) qp.Expression {
		return qp.Int(len(arguments[0].Invoke().(*Array).data))
	})("Get", func(arguments ...qp.Expression) qp.Expression {
		if len(arguments) != 2 {
			log.Panic("array Get() arguments error")
		}
		array, ok := arguments[0].Invoke().(*Array)
		if ok == false {
			log.Panic("object not array type")
		}
		i, ok := arguments[1].(qp.Int)
		if ok == false {
			log.Panic("is not array arguments error")
		}
		if len(array.data) <= int(i) {
			log.Panic("index out of range")
		}
		return array.data[i]
	})
}
