package ast

import (
	"gitlab.com/akzj/qp/lexer"
	"gitlab.com/akzj/qp/runtime"
)

type BaseObject interface {
	runtime.Invokable
	GetObject(label string) *runtime.Object
	AllocObject(label string) *runtime.Object
	Clone() BaseObject
}

type TypeObjectPropTemplate struct {
	Name string
	Exp  runtime.Invokable
}

func (t TypeObjectPropTemplate) String() string {
	return t.Name + ":" + t.Exp.String()
}

type TypeObject struct {
	VM    *runtime.VMContext
	Label string
	//Init Statement when create objects
	Init                    bool
	TypeObjectPropTemplates []TypeObjectPropTemplate
	//user define Function
	objects map[string]*runtime.Object
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

func (sObj *TypeObject) Invoke() runtime.Invokable {
	return sObj
}

func (sObj *TypeObject) GetType() lexer.Type {
	return lexer.TypeObjectType
}

func (sObj *TypeObject) GetObject(label string) *runtime.Object {
	object, ok := sObj.objects[label]
	if ok {
		return object
	}
	return nil
}

func (sObj *TypeObject) AllocObject(label string) *runtime.Object {
	if sObj.objects == nil {
		sObj.objects = map[string]*runtime.Object{}
	}
	object, ok := sObj.objects[label]
	if ok {
		return object
	} else {
		object = &runtime.Object{
			Pointer: NilObj,
			Label:   label,
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

func (sObj *TypeObject) AddObject(k string, v *runtime.Object) {
	if v == nil {
		panic(v)
	}
	if sObj.objects == nil {
		sObj.objects = map[string]*runtime.Object{}
	}
	sObj.objects[k] = v
}

func (sObj *TypeObject) GetObjects() map[string]*runtime.Object {
	return sObj.objects
}
