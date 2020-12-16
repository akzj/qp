package qp

import (
	"gitlab.com/akzj/qp/ast"
	"gitlab.com/akzj/qp/lexer"
	vm2 "gitlab.com/akzj/qp/vm"
)

type BaseObject interface {
	ast.Expression
	GetObject(label string) *ast.Object
	AllocObject(label string) *ast.Object
	Clone() BaseObject
}

type TypeObjectPropTemplate struct {
	name string
	exp  ast.Expression
}

func (t TypeObjectPropTemplate) String() string {
	return t.name + ":" + t.exp.String()
}

type TypeObject struct {
	vm    *vm2.VMContext
	label string
	//init statement when create objects
	init                    bool
	typeObjectPropTemplates []TypeObjectPropTemplate
	//user define function
	objects map[string]*ast.Object
}

func (sObj *TypeObject) String() string {
	str := "type " + sObj.label + "{"
	for index, exp := range sObj.typeObjectPropTemplates {
		if index == 0 {
			str += "\n"
		}
		str += "\t" + exp.String() + ";\n"

	}
	return str + "}"
}

func (sObj *TypeObject) Invoke() ast.Expression {
	return sObj
}

func (sObj *TypeObject) GetType() lexer.Type {
	return lexer.TypeObjectType
}

func (sObj *TypeObject) GetObject(label string) *ast.Object {
	object, ok := sObj.objects[label]
	if ok {
		return object
	}
	return nil
}

func (sObj *TypeObject) AllocObject(label string) *ast.Object {
	if sObj.objects == nil {
		sObj.objects = map[string]*ast.Object{}
	}
	object, ok := sObj.objects[label]
	if ok {
		return object
	} else {
		object = &ast.Object{
			inner: ast.nilObject,
			label: label,
		}
		sObj.objects[label] = object
	}
	return object
}

func (sObj *TypeObject) Clone() BaseObject {
	clone := *sObj
	for k, v := range sObj.objects {
		clone.addObject(k, v)
	}
	if len(sObj.typeObjectPropTemplates) != 0 {
		clone.typeObjectPropTemplates = make([]TypeObjectPropTemplate, len(sObj.typeObjectPropTemplates))
		copy(clone.typeObjectPropTemplates, sObj.typeObjectPropTemplates)
	}
	return &clone
}

func (sObj *TypeObject) addObject(k string, v *ast.Object) {
	if v == nil {
		panic(v)
	}
	if sObj.objects == nil {
		sObj.objects = map[string]*ast.Object{}
	}
	sObj.objects[k] = v
}
