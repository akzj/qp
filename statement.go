package qp

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

type Statements []Statement

type Statement interface {
	Expression
}

type IfStatement struct {
	vm               *VMContext
	check            Expression
	statement        Statements
	elseIfStatements []*IfStatement
	elseStatement    Statements
}

type ReturnStatement struct {
	express   Expression
	returnVal Expression
}

//just new Object
type VarStatement struct {
	ctx    *VMContext
	label  string
	object *TypeObject
}

type getVarStatement struct {
	ctx   *VMContext
	label string
}

//a.b.c.d
type getObjectPropStatement struct {
	this      bool
	getObject *getObjectObjectStatement
}

type getObjectObjectStatement struct {
	vmContext *VMContext
	labels    []string
}

type FuncCallStatement struct {
	vm        *VMContext
	label     string
	getObject *getObjectPropStatement
	arguments Expressions
}

type AssignStatement struct {
	ctx        *VMContext
	label      string
	expression Expression
	getObject  *getObjectObjectStatement
}

type VarAssignStatement struct {
	object     *TypeObject //belong to struct object member field
	ctx        *VMContext  //global or stack var
	label      string      //var name : var a,`a` is the label
	expression Expression  // init expression : var a = 1+1
}

type IncFieldStatement struct {
	ctx   *VMContext
	label string
}

type BreakStatement struct {
}

type NopStatement struct {
}

type FuncStatement struct {
	closure      bool
	label        string
	labels       []string // struct objects function eg:user.add(){}
	parameters   []string // parameter label
	closureLabel []string // closure label
	closureInit  bool
	statements   Statements // function body
	vm           *VMContext // vm context
	closureObjs  []Expression
}

func (f *FuncStatement) invoke() Expression {
	return f
}

type ForStatement struct {
	vm             *VMContext
	preStatement   Expression
	checkStatement Expression
	postStatement  Expression
	statements     Statements
}

type StructObjectInitStatement struct {
	label          string // TypeObject label
	vm             *VMContext
	initStatements Statements
}

func (g *getObjectPropStatement) invoke() Expression {
	return g.getObject.invoke().(*Object).inner
}

func (g *getObjectObjectStatement) invoke() Expression {
	object := g.vmContext.getObject(g.labels[0])
	if object == nil {
		log.Panic("getObject failed", g.labels[0])
	}
	structObj, ok := object.inner.(BaseObject)
	if ok == false {
		log.Panic("objects type no struct objects,error",
			g.labels, reflect.TypeOf(object.inner).String())
	}
	/*
	 user.id = 1 // bind 1 to user.id
	 println(user.id)// visit user.id
	*/
	var obj = object
	for i := 1; i < len(g.labels); i++ {
		obj = structObj.allocObject(g.labels[i])
		//last label
		if i != len(g.labels)-1 {
			var ok bool
			structObj, ok = obj.inner.(*TypeObject)
			if ok == false {
				label := strings.Join(g.labels[:i+1], ".")
				log.Panic("objects is no struct objects type", label)
			}
		}
	}
	return obj
}

func (g *getObjectObjectStatement) getType() Type {
	return getObjectObjectStatementType
}

func (g *getObjectPropStatement) getType() Type {
	return propObjectStatementType
}

func (statement *StructObjectInitStatement) invoke() Expression {
	object := statement.vm.cloneTypeObject(statement.label)
	if object == nil {
		log.Panicf("cloneTypeObject with label `%s` failed", statement.label)
	}
	for _, initStatement := range statement.initStatements {
		object.initStatement = append(object.initStatement, initStatement)
	}
	for _, init := range object.initStatement {
		switch s := init.(type) {
		case *VarAssignStatement:
			s.object = object
		case *VarStatement:
			s.object = object
		case *NopStatement:
			continue
		default:
			panic("unknown statement " + reflect.TypeOf(init).String())
		}
		init.invoke()
	}
	return object
}

func (statement *StructObjectInitStatement) getType() Type {
	return typeObjectInitStatementType
}

func (f *FuncStatement) prepareArgumentBind(inArguments Expressions) error {
	//lambda function no bind this to objects
	if f.closure && len(inArguments) != 0 {
		statement, ok := inArguments[0].(*getObjectPropStatement)
		if ok && statement.this {
			inArguments = inArguments[1:]
		}
	}
	if len(f.parameters) != len(inArguments) {
		log.Println("argument size no match", len(f.parameters), len(inArguments))
		return fmt.Errorf("argument size no match")
	}

	var results []Expression
	for _, expression := range inArguments {
		results = append(results, expression.invoke())
	}

	f.vm.pushStackFrame(true)
	// put closure objects to stack
	for index := range f.closureLabel {
		//object :=
			f.vm.allocObject(f.closureLabel[index]).inner = f.closureObjs[index]
		/*switch closureObj := f.closureObjs[index].(type) {
		case *Object:
			object.inner = closureObj.inner
		default:
			object.inner = closureObj
		}*/
	}

	//make stack for this function
	for index, result := range results {
		f.vm.allocObject(f.parameters[index]).inner = result
	}
	return nil
}

func (f *FuncStatement) call(arguments ...Expression) Expression {
	defer f.vm.popStackFrame()
	if err := f.prepareArgumentBind(arguments); err != nil {
		return nil
	}
	for _, statement := range f.statements {
		if ret, ok := statement.invoke().(*ReturnStatement); ok {
			return ret.returnVal
		}
	}
	return nil
}

func (f *FuncStatement) getType() Type {
	return FuncStatementType
}

func (f *FuncStatement) doClosureInit() {
	if f.closureInit {
		return
	}
	var closureObjs []Expression
	for _, label := range f.closureLabel {
		obj := f.vm.getObject(label)
		if obj == nil {
			log.Panic("no find obj with label", label)
		}
		closureObjs = append(closureObjs, obj.inner)
	}
	f.closureObjs = closureObjs
}

func (expression *AssignStatement) invoke() Expression {
	val := expression.expression.invoke()
	var inner = val
	switch obj := val.(type) {
	case *Object:
		inner = obj.inner
	}
	if expression.getObject != nil {
		object := expression.getObject.invoke()
		object.(*Object).inner = inner
	} else {
		object := expression.ctx.getObject(expression.label)
		if object == nil {
			log.Panic("AssignStatement getObject failed", expression.label)
		}
		object.inner = inner
	}
	return nil
}

func (expression *AssignStatement) getType() Type {
	return assignStatementType
}

func (n *NopStatement) invoke() Expression {
	return n
}

func (n *NopStatement) getType() Type {
	return nopStatementType
}

func (f *ForStatement) invoke() Expression {
	f.vm.pushStackFrame(false) //make stack frame

	//make for brock stack
	f.preStatement.invoke()

	for ; ; {
		val := f.checkStatement.invoke()
		bObj, ok := val.(*BoolObject)
		if ok == false {
			log.Panic("for checkStatement expect BoolObject")
		}
		if bObj.val == false {
			f.vm.popStackFrame() //end of for
			return nil
		}
		f.vm.pushStackFrame(false) //make stack frame for `{` brock
		for _, statement := range f.statements {
			val := statement.invoke()
			if val == breakObject {
				return nil
			}
			if _, ok := val.(*ReturnStatement); ok {
				return val
			}
		}
		f.vm.popStackFrame()
		f.postStatement.invoke()
	}
}

func (f *ForStatement) getType() Type {
	return forTokenType
}

func (statement *IncFieldStatement) invoke() Expression {
	object := statement.ctx.getObject(statement.label)
	if object == nil {
		log.Panic("no find Object with label ", statement.label)
	}
	innerObject := object.invoke()
	switch obj := innerObject.(type) {
	case *IntObject:
		obj.val++
	default:
		panic("unknown type " + reflect.TypeOf(innerObject).String())
	}
	return nil
}

func (statement *IncFieldStatement) getType() Type {
	return incOperatorTokenType
}

func (Statements) getType() Type {
	return statementsType
}

func (f *FuncCallStatement) invoke() Expression {
	if f.getObject != nil {
		var object = f.getObject.invoke()
		return object.(Function).call(f.arguments...)
	} else {
		function, err := f.vm.getFunction(f.label)
		if err == nil {
			return function.call(f.arguments...)
		}
		log.Panic("getFunction failed", f.label, err)
		return nil
	}
}

func (f *FuncCallStatement) getType() Type {
	return funcTokenType
}

func (f *getVarStatement) invoke() Expression {
	object := f.ctx.getObject(f.label)
	if object == nil {
		log.Panicln("no find Object with label", f.label)
	}
	return object.invoke()
}

func (f *getVarStatement) getType() Type {
	return labelType
}

func (v *VarStatement) invoke() Expression {
	if v.object != nil {
		v.object.allocObject(v.label).inner = nilObject
	} else {
		v.ctx.allocObject(v.label).inner = nilObject
	}
	return nil
}

func (v VarStatement) getType() Type {
	return varTokenType
}

func (expression *VarAssignStatement) invoke() Expression {
	obj := expression.expression.invoke()
	var object *Object
	if expression.object != nil {
		object = expression.object.allocObject(expression.label)
	} else {
		object = expression.ctx.allocObject(expression.label)
	}
	if obj == nil {
		panic(obj)
	}
	switch obj := obj.(type) {
	case *Object:
		object.inner = obj.inner
	default:
		object.inner = obj
	}
	object.initType()
	return nil
}

func (expression *VarAssignStatement) getType() Type {
	return varAssignTokenType
}

func (r *ReturnStatement) invoke() Expression {
	//log.Println("ReturnStatement invoke")
	if r.returnVal != nil {
		return r
	}
	val := r.express.invoke()
	switch inner := val.(type) {
	case *ReturnStatement:
		return inner
	default:
		r.returnVal = val
	}
	return r
}

func (ReturnStatement) getType() Type {
	return returnTokenType
}

func (IfStatement) getType() Type {
	return ifTokenType
}

func (statements Statements) invoke() Expression {
	var val Expression
	for _, statement := range statements {
		val = statement.invoke()
		if _, ok := val.(*ReturnStatement); ok {
			return val
		}
	}
	return val
}

func (ifStm *IfStatement) invoke() Expression {
	check := ifStm.check.invoke()
	if _, ok := check.(*BoolObject); ok == false {
		log.Panic("if statement check require boolObject")
	}
	if check.(*BoolObject).val {
		ifStm.vm.pushStackFrame(false) //make  if brock stack
		val := ifStm.statement.invoke()
		ifStm.vm.popStackFrame() //release  if brock stack
		return val
	} else {
		for _, stm := range ifStm.elseIfStatements {
			elseIf := stm.check.invoke()
			if _, ok := elseIf.(*BoolObject); ok == false {
				log.Panicln("else if require bool result")
			}
			if elseIf.(*BoolObject).val {
				ifStm.vm.pushStackFrame(false) //make  if brock stack
				val := stm.statement.invoke()
				ifStm.vm.popStackFrame() //release  if brock stack
				return val
			}
		}
		if ifStm.elseStatement != nil {
			ifStm.vm.pushStackFrame(false) //make  brock stack
			val := ifStm.elseStatement.invoke()
			ifStm.vm.popStackFrame() //release  if brock stack
			return val
		}
	}
	return nil
}
