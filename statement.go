package qp

import (
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

type FuncCallQueueStatement struct {
	statement []*FuncCallStatement
}

type FuncCallStatement struct {
	vm        *VMContext
	label     string
	getObject *getObjectPropStatement
	function  Function
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

var nopStatement = NopStatement{}

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

func (f *FuncCallQueueStatement) Invoke() Expression {
	var function Function
	for index, call := range f.statement {

		//prepare closure for function object
		if index == 0 {
			if call.function != nil {
				if statement, ok := call.function.(*FuncStatement); ok {
					statement.doClosureInit()
				}
			}
		} else {
			call.function = function
		}
		expression := call.Invoke()
		if index != len(f.statement)-1 {
			var ok bool
			function, ok = expression.(Function)
			if ok == false {
				log.Panic("statement no callable object",
					reflect.TypeOf(expression).String())
			}
			continue
		}
		return expression
	}
	return nil
}

func (f *FuncCallQueueStatement) getType() Type {
	return funcCallQueueStatementType
}

func (f *FuncStatement) Invoke() Expression {
	f.doClosureInit()
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

type makeArrayStatement struct {
	vm             *VMContext
	initStatements Statements
}

func (m *makeArrayStatement) Invoke() Expression {
	var array = &Array{
		TypeObject: TypeObject{
			vm:      m.vm,
			label:   "array",
			init:    true,
			objects: arrayBuiltInFunctions,
		},
	}
	for _, statement := range m.initStatements {
		array.data = append(array.data, statement.Invoke())
	}
	return array
}

func (m *makeArrayStatement) getType() Type {
	return arrayObjectType
}

func (g *getObjectPropStatement) Invoke() Expression {
	obj := g.getObject.Invoke()
	if obj == nilObject {
		return obj
	}
	return obj.(*Object).inner
}

func (g *getObjectObjectStatement) Invoke() Expression {
	object := g.vmContext.getObject(g.labels[0])
	if object == nil {
		log.Panicf("getObject failed `%s`", g.labels[0])
	}
	structObj, ok := object.inner.(BaseObject)
	if ok == false {
		log.Panic("objects type no struct objects,error",
			g.labels, reflect.TypeOf(object.inner).String())
	}
	/*
	 user.id = 1 // bind 1 to user.id
	 printlnFunc(user.id)// visit user.id
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

func (statement *StructObjectInitStatement) Invoke() Expression {
	object := statement.vm.cloneTypeObject(statement.label)
	if object == nil {
		log.Panicf("cloneTypeObject with label `%s` failed", statement.label)
	}
	for _, initStatement := range statement.initStatements {
		object.initStatement = append(object.initStatement, initStatement)
	}
	for _, init := range object.initStatement {
		switch s := init.(type) {
		case VarAssignStatement:
			s.object = object
		case VarStatement:
			s.object = object
		case NopStatement:
			continue
		default:
			panic("unknown statement " + reflect.TypeOf(init).String())
		}
		init.Invoke()
	}
	return object
}

func (statement *StructObjectInitStatement) getType() Type {
	return typeObjectInitStatementType
}

func (f *FuncStatement) prepareArgumentBind(inArguments Expressions) {
	if len(f.parameters) != len(inArguments) {
		log.Panic("argument size no match ", len(f.parameters), len(inArguments))
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
	for index, result := range inArguments {
		f.vm.allocObject(f.parameters[index]).inner = result
	}
}

func (f *FuncStatement) call(arguments ...Expression) Expression {
	defer f.vm.popStackFrame()
	f.prepareArgumentBind(arguments)
	for _, statement := range f.statements {
		result := statement.Invoke()
		if ret, ok := result.(ReturnStatement); ok {
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
	f.closureInit = true
	var closureObjs []Expression
	for _, label := range f.closureLabel {
		obj := f.vm.getObject(label)
		if obj == nil {
			log.Panicf("no find obj with label `%s`", label)
		}
		closureObjs = append(closureObjs, obj.inner)
	}
	f.closureObjs = closureObjs
}

func (expression *AssignStatement) Invoke() Expression {
	val := expression.expression.Invoke()
	var inner = val
	switch obj := val.(type) {
	case *Object:
		inner = obj.inner
	}
	if expression.getObject != nil {
		object := expression.getObject.Invoke()
		object.(*Object).inner = inner
	} else {
		object := expression.ctx.getObject(expression.label)
		if object == nil {
			log.Panicf("AssignStatement getObject failed `%s`", expression.label)
		}
		object.inner = inner
	}
	return nil
}

func (expression *AssignStatement) getType() Type {
	return assignStatementType
}

func (NopStatement) Invoke() Expression {
	return nopStatement
}

func (n NopStatement) getType() Type {
	return nopStatementType
}

func (f *ForStatement) Invoke() Expression {
	f.vm.pushStackFrame(false) //make stack frame

	//make for brock stack
	f.preStatement.Invoke()

	for ; ; {
		val := f.checkStatement.Invoke()
		bObj, ok := val.(Bool)
		if ok == false {
			log.Panic("for checkStatement expect Bool")
		}
		if bObj == false {
			f.vm.popStackFrame() //end of for
			return nil
		}
		f.vm.pushStackFrame(false) //make stack frame for `{` brock
		for _, statement := range f.statements {
			val := statement.Invoke()
			if val == breakObject {
				return nil
			}
			if _, ok := val.(ReturnStatement); ok {
				return val
			}
		}
		f.vm.popStackFrame()
		f.postStatement.Invoke()
	}
}

func (f *ForStatement) getType() Type {
	return forTokenType
}

func (statement IncFieldStatement) Invoke() Expression {
	object := statement.ctx.getObject(statement.label)
	if object == nil {
		log.Panic("no find Object with label ", statement.label)
	}
	object.inner = Int(int(object.inner.(Int)) + 1)
	return nil
}

func (statement *IncFieldStatement) getType() Type {
	return incOperatorTokenType
}

func (Statements) getType() Type {
	return statementsType
}

func (f *FuncCallStatement) Invoke() Expression {
	var function Function
	var arguments Expressions
	if f.function != nil {
		function = f.function
	} else if f.getObject != nil {
		var object = f.getObject.Invoke()
		if object == nil {
			log.Panic("no find function", f.label)
		}
		//lambda function can't bind this to first argument
		if funcStatement, ok := object.(*FuncStatement); ok {
			if funcStatement.closure && len(f.arguments) != 0 {
				statement, ok := f.arguments[0].(*getObjectPropStatement)
				if ok && statement.this {
					f.arguments = f.arguments[1:]
				}
			}
		}
		function = object.(Function)
	} else {
		var err error
		function, err = f.vm.getFunction(f.label)
		if err != nil {
			log.Panicf("no find function with label`%s`", f.label)
		}
	}
	for _, argument := range f.arguments {
		arguments = append(arguments, argument.Invoke())
	}
	result := function.call(arguments...)
	if ret, ok := result.(ReturnStatement); ok {
		return ret.returnVal
	}
	return result
}

func (f *FuncCallStatement) getType() Type {
	return funcTokenType
}

func (f *getVarStatement) Invoke() Expression {
	object := f.ctx.getObject(f.label)
	if object == nil {
		log.Panicf("getObject faild `%s`", f.label)
	}
	if object.inner == nil {
		object.inner = nilObject
	}
	return object.Invoke()
}

func (f *getVarStatement) getType() Type {
	return labelType
}

func (v VarStatement) Invoke() Expression {
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

func (expression VarAssignStatement) Invoke() Expression {
	obj := expression.expression.Invoke()
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

func (expression VarAssignStatement) getType() Type {
	return varAssignTokenType
}

func (r ReturnStatement) Invoke() Expression {
	if r.returnVal != nil {
		return r
	}
	return ReturnStatement{returnVal: r.express.Invoke()}
}

func (ReturnStatement) getType() Type {
	return returnTokenType
}

func (IfStatement) getType() Type {
	return ifTokenType
}

func (statements Statements) Invoke() Expression {
	var val Expression
	for _, statement := range statements {
		val = statement.Invoke()
		if _, ok := val.(ReturnStatement); ok {
			return val
		} else if _, ok := val.(*BreakObject); ok {
			return breakObject
		}
	}
	return val
}

func (ifStm *IfStatement) Invoke() Expression {
	check := ifStm.check.Invoke()
	if _, ok := check.(Bool); ok == false {
		log.Panic("if statement check require boolObject")
	}
	if check.(Bool) {
		ifStm.vm.pushStackFrame(false) //make  if brock stack
		val := ifStm.statement.Invoke()
		ifStm.vm.popStackFrame() //release  if brock stack
		return val
	} else {
		for _, stm := range ifStm.elseIfStatements {
			elseIf := stm.check.Invoke()
			if _, ok := elseIf.(Bool); ok == false {
				log.Panicln("else if require bool result")
			}
			if elseIf.(Bool) {
				ifStm.vm.pushStackFrame(false) //make  if brock stack
				val := stm.statement.Invoke()
				ifStm.vm.popStackFrame() //release  if brock stack
				return val
			}
		}
		if ifStm.elseStatement != nil {
			ifStm.vm.pushStackFrame(false) //make  brock stack
			val := ifStm.elseStatement.Invoke()
			ifStm.vm.popStackFrame() //release  if brock stack
			return val
		}
	}
	return nil
}
