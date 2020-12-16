package ast

import (
	"gitlab.com/akzj/qp/lexer"
)

type BaseObject interface {
	Expression
	GetObject(label string) *Object
	AllocObject(label string) *Object
	Clone() BaseObject
}

type TypeObjectPropTemplate struct {
	Name string
	Exp  Expression
}

func (t TypeObjectPropTemplate) String() string {
	return t.Name + ":" + t.Exp.String()
}

type TypeObject struct {
	VM    *VMContext
	Label string
	//Init Statement when create objects
	Init                    bool
	TypeObjectPropTemplates []TypeObjectPropTemplate
	//user define Function
	objects map[string]*Object
}

func (sObj *TypeObject) String() string {
	str := "type " + sObj.Label + "{"
	for index, exp := range sObj.TypeObjectPropTemplates {
		if index == 0 {
			str += "\n"
		}
		str += "\t" + exp.String() + ";\n"

	}
	return str + "}"
}

func (sObj *TypeObject) Invoke() Expression {
	return sObj
}

func (sObj *TypeObject) GetType() lexer.Type {
	return lexer.TypeObjectType
}

func (sObj *TypeObject) GetObject(label string) *Object {
	object, ok := sObj.objects[label]
	if ok {
		return object
	}
	return nil
}

func (sObj *TypeObject) AllocObject(label string) *Object {
	if sObj.objects == nil {
		sObj.objects = map[string]*Object{}
	}
	object, ok := sObj.objects[label]
	if ok {
		return object
	} else {
		object = &Object{
			Inner: NilObj,
			Label: label,
		}
		sObj.objects[label] = object
	}
	return object
}

func (sObj *TypeObject) Clone() BaseObject {
	clone := *sObj
	for k, v := range sObj.objects {
		clone.AddObject(k, v)
	}
	if len(sObj.TypeObjectPropTemplates) != 0 {
		clone.TypeObjectPropTemplates = make([]TypeObjectPropTemplate, len(sObj.TypeObjectPropTemplates))
		copy(clone.TypeObjectPropTemplates, sObj.TypeObjectPropTemplates)
	}
	return &clone
}

func (sObj *TypeObject) AddObject(k string, v *Object) {
	if v == nil {
		panic(v)
	}
	if sObj.objects == nil {
		sObj.objects = map[string]*Object{}
	}
	sObj.objects[k] = v
}
