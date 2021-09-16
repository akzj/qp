package ast

import (
	"log"
	"reflect"
	"strings"

	"gitlab.com/akzj/qp/lexer"
	"gitlab.com/akzj/qp/runtime"
)

type Expressions []Expression

func (statements Expressions) String() string {
	var str string
	for index, state := range statements {
		str += state.String()
		if index != len(statements)-1 {
			str += "\n"
		}
	}
	return str
}

type Expression interface {
	runtime.Invokable
}

type ReturnStatement struct {
	Exp runtime.Invokable
	Val runtime.Invokable
}

func (r ReturnStatement) String() string {
	if r.Val != nil {
		return "return " + r.Val.String()
	} else {
		return "return " + r.Exp.String()
	}
}

type PeriodStatement struct {
	Val string
	Exp runtime.Invokable
}

func (p PeriodStatement) Invoke() runtime.Invokable {
	object := unwrapObject(p.Exp.Invoke())
	switch obj := object.(type) {
	case BaseObject:
		return obj.AllocObject(p.Val)
	default:
		log.Panicf("Left `%s` `%s` is no Exp type", p.Val, reflect.TypeOf(obj).String())
	}
	return nil
}

func (p PeriodStatement) GetType() lexer.Type {
	return lexer.PeriodType
}

func (p PeriodStatement) String() string {
	return p.Exp.String() + "." + p.Val
}

type GetVarStatement struct {
	VM    *runtime.VMRuntime
	Label string
}

func (f GetVarStatement) String() string {
	return f.Label
}

//a.b.c.d
type getObjectPropStatement struct {
	this      bool
	getObject *getObjectObjectStatement
}

func (g *getObjectPropStatement) String() string {
	panic("implement me")
}

type getObjectObjectStatement struct {
	vmContext *runtime.VMRuntime
	labels    []string
}

type AssignStatement struct {
	Exp  runtime.Invokable
	Left runtime.Invokable
}

func (expression AssignStatement) String() string {
	return expression.Left.String() + "=" + expression.Exp.String()
}

type IncFieldStatement struct {
	Exp runtime.Invokable
}

func (statement IncFieldStatement) String() string {
	return statement.Exp.String() + "++"
}

type BreakStatement struct {
}

type NopStatement struct {
}

func (n NopStatement) String() string {
	return "nop"
}

type ObjectInitStatement struct {
	VM            *runtime.VMRuntime
	Exp           runtime.Invokable
	PropTemplates []TypeObjectPropTemplate
}

func (statement ObjectInitStatement) String() string {
	var str string
	for _, statement := range statement.PropTemplates {
		str += statement.String() + "\n"
	}
	return statement.Exp.String() + "{" + str + "}"
}

type ArrayGetElement struct {
	Exp   runtime.Invokable
	Index runtime.Invokable
}

func (g ArrayGetElement) Invoke() runtime.Invokable {
	panic("implement me")
}

func (g ArrayGetElement) GetType() lexer.Type {
	panic("implement me")
}

func (g ArrayGetElement) String() string {
	panic("implement me")
}

type MakeArrayStatement struct {
	vm    *runtime.VMRuntime
	Inits Expressions
}

func (m *MakeArrayStatement) String() string {
	var str = "["
	for index, statement := range m.Inits {
		if index != 0 {
			str += ","
		}
		str += statement.String()
	}
	return str + "]"
}

func (m *MakeArrayStatement) Invoke() runtime.Invokable {
	var array = &Array{}
	for _, statement := range m.Inits {
		array.Data = append(array.Data, statement.Invoke())
	}
	return array
}

func (m *MakeArrayStatement) GetType() lexer.Type {
	return lexer.ArrayObjectType
}

func (g *getObjectPropStatement) Invoke() runtime.Invokable {
	obj := g.getObject.Invoke()
	if obj == NilObj {
		return obj
	}
	return obj.(*runtime.Object).Pointer.(runtime.Invokable)
}

func (g *getObjectObjectStatement) Invoke() runtime.Invokable {
	object := g.vmContext.GetObject(g.labels[0])
	if object == nil {
		log.Panicf("Left failed `%s`", g.labels[0])
	}
	structObj, ok := object.Pointer.(BaseObject)
	if ok == false {
		log.Panic("objects type no struct objects,error",
			g.labels, reflect.TypeOf(object.Pointer).String())
	}
	/*
	 user.id = 1 // bind 1 to user.id
	 printlnFunc(user.id)// visit user.id
	*/
	var obj = object
	for i := 1; i < len(g.labels); i++ {
		obj = structObj.AllocObject(g.labels[i])
		//last Name
		if i != len(g.labels)-1 {
			var ok bool
			structObj, ok = obj.Pointer.(*TypeObject)
			if ok == false {
				label := strings.Join(g.labels[:i+1], ".")
				log.Panic("objects is no struct objects type", label)
			}
		}
	}
	return obj
}

func (g *getObjectObjectStatement) getType() lexer.Type {
	return lexer.GetObjectObjectStatementType
}

func (g *getObjectPropStatement) GetType() lexer.Type {
	return lexer.PropObjectStatementType
}

func (statement ObjectInitStatement) Invoke() runtime.Invokable {
	object := statement.Exp.Invoke().(*runtime.Object).Pointer.(BaseObject).Clone().(*TypeObject)

Loop:
	for _, init := range object.TypeObjectPropTemplates {
		for _, prod := range statement.PropTemplates {
			if init.Name == prod.Name {
				continue Loop
			}
		}
		propObject := object.AllocObject(init.Name)
		propObject.Pointer = init.Exp.Invoke()
	}

	for _, init := range statement.PropTemplates {
		propObject := object.AllocObject(init.Name)
		propObject.Pointer = init.Exp.Invoke()
	}
	return object
}

func (statement ObjectInitStatement) GetType() lexer.Type {
	return lexer.TypeObjectInitStatementType
}

func (expression AssignStatement) Invoke() runtime.Invokable {
	left := expression.Left.Invoke()
	switch right := expression.Exp.Invoke().(type) {
	case *runtime.Object:
		left.(*runtime.Object).Pointer = right.Pointer
	default:
		left.(*runtime.Object).Pointer = right
	}
	return nil
}

func (expression AssignStatement) GetType() lexer.Type {
	return lexer.AssignStatementType
}

func (NopStatement) Invoke() runtime.Invokable {
	return NopStatement{}
}

func (n NopStatement) GetType() lexer.Type {
	return lexer.NopStatementType
}

func (statement IncFieldStatement) Invoke() runtime.Invokable {
	object := statement.Exp.Invoke().(*runtime.Object)
	object.Pointer = object.Pointer.(Int) + 1
	return nil
}

func (statement IncFieldStatement) GetType() lexer.Type {
	return lexer.IncType
}

func (Expressions) GetType() lexer.Type {
	return lexer.StatementsType
}

func (f GetVarStatement) Invoke() runtime.Invokable {
	return f.VM.GetObject(f.Label)
}

func (f GetVarStatement) GetType() lexer.Type {
	return lexer.IDType
}

func (r ReturnStatement) Invoke() runtime.Invokable {
	if r.Val != nil {
		return r
	}
	exp := r.Exp.Invoke()
	switch obj := exp.(type) {
	case *runtime.Object:
		exp = obj.Pointer
	case ReturnStatement:
		return obj
	}
	return ReturnStatement{Val: exp}
}

func (ReturnStatement) GetType() lexer.Type {
	return lexer.ReturnType
}

func (statements Expressions) Invoke() runtime.Invokable {
	var val runtime.Invokable
	for _, statement := range statements {
		val = statement.Invoke()
		if _, ok := val.(ReturnStatement); ok {
			return val
		} else if _, ok := val.(*BreakObject); ok {
			return BreakObj
		}
	}
	return val
}
