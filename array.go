package qp

import "log"

type Array struct {
	TypeObject
	data []Expression
}

type appendArray struct {
}

func (a *appendArray) Invoke() Expression {
	return a
}

func (a *appendArray) getType() Type {
	return FuncStatementType
}

func (a *appendArray) call(arguments ...Expression) Expression {
	array := arguments[0].(*Array)
	for _, exp := range arguments[1:] {
		array.data = append(array.data, exp.Invoke())
	}
	return array
}

type getArray struct {
}

func (g *getArray) Invoke() Expression {
	return g
}

func (g *getArray) getType() Type {
	return FuncStatementType
}

func (g *getArray) call(arguments ...Expression) Expression {
	if len(arguments) != 2 {
		log.Panic("array get() arguments error")
	}
	array, ok := arguments[0].(*Array)
	if ok == false {
		log.Panic("object not array type")
	}
	i, ok := arguments[1].(*IntObject)
	if ok == false {
		log.Panic("is not array arguments error")
	}
	if len(array.data) <= int(*i) {
		log.Panic("index out of range")
	}
	return array.data[*i]
}
