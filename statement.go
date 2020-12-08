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
type propObjectStatement struct {
	this      bool
	vmContext *VMContext
	labels    []string
}

type FuncCallStatement struct {
	vm        *VMContext
	label     string
	getObject *propObjectStatement
	arguments Expressions
}

type AssignStatement struct {
	ctx        *VMContext
	label      string
	expression Expression
	getObject  *propObjectStatement
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

func (g *propObjectStatement) invoke() (Expression, error) {
	object := g.vmContext.getObject(g.labels[0])
	if object == nil {
		log.Panic("getObject failed", g.labels[0])
	}
	if object.inner == nil {
		log.Panic("objects nil", g.labels[0])
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
	return obj, nil
}

func (g *propObjectStatement) getType() Type {
	return propObjectStatementType
}

func (statement *StructObjectInitStatement) invoke() (Expression, error) {
	object := statement.vm.cloneTypeObject(statement.label)
	if object == nil {
		return nil, fmt.Errorf("cloneTypeObject with label `%s` failed", statement.label)
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
		if _, err := init.invoke(); err != nil {
			log.Panic("struct objects init failed")
		}
	}
	return object, nil
}

func (statement *StructObjectInitStatement) getType() Type {
	return typeObjectInitStatementType
}

func (f *FuncStatement) prepareArgumentBind(inArguments Expressions) error {
	//lambda function no bind this to objects
	if f.closure && len(inArguments) != 0 {
		statement, ok := inArguments[0].(*propObjectStatement)
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
		result, err := expression.invoke()
		if err != nil {
			log.Println("argument invoke() failed", err)
			return err
		}
		results = append(results, result)
	}

	f.vm.pushStackFrame(true)

	// put closure objects to stack
	for index := range f.closureLabel {
		label := f.closureLabel[index]
		object := f.vm.allocObject(label)
		switch closureObj := f.closureObjs[index].(type) {
		case *Object:
			object.inner = closureObj.inner
		default:
			object.inner = closureObj
		}
	}

	//make stack for this function
	for index, result := range results {
		label := f.parameters[index]
		object := f.vm.allocObject(label)
		if object == nil {
			panic("allocObject nil")
		}
		switch obj := result.(type) {
		case *Object:
			if obj == nil {
				panic("obj nil")
			}
			object.inner = obj.inner
		default:
			object.inner = result
		}
	}

	return nil
}

func (f *FuncStatement) invoke(arguments ...Expression) (Expression, error) {
	defer f.vm.popStackFrame()
	if err := f.prepareArgumentBind(arguments); err != nil {
		return nil, err
	}
	for _, statement := range f.statements {
		val, err := statement.invoke()
		if err != nil {
			log.Panic("statement.invoke() failed", err)
		}
		if val != nil {
			//function return
			if ret, ok := val.(*ReturnStatement); ok {
				return ret.returnVal, nil
			}
		}
	}
	return nil, nil
}

func (f *FuncStatement) getType() Type {
	return FuncStatementType
}

func (f *FuncStatement) doClosureInit() error {
	if f.closureInit {
		return nil
	}
	var closureObjs []Expression
	for _, label := range f.closureLabel {
		obj := f.vm.getObject(label)
		if obj == nil {
			log.Panic("no find obj with label", label)
		}
		closureObjs = append(closureObjs, obj)
	}
	f.closureObjs = closureObjs
	return nil
}

func (expression *AssignStatement) invoke() (Expression, error) {
	val, err := expression.expression.invoke()
	if err != nil {
		log.Panic("AssignStatement .expression.invoke() failed", err.Error())
	}
	var inner interface{} = val
	switch obj := val.(type) {
	case *Object:
		inner = obj.inner
	}
	if expression.getObject != nil {
		object, err := expression.getObject.invoke()
		if err != nil {
			log.Panic("on.getObject.invoke() failed", err)
			return nil, err
		}
		object.(*Object).inner = inner
	} else {
		object := expression.ctx.getObject(expression.label)
		if object == nil {
			log.Println("AssignStatement .expression.getObject failed", object.label)
			return nil, err
		}
		object.inner = inner
	}
	return nil, nil
}

func (expression *AssignStatement) getType() Type {
	return assignStatementType
}

func (n *NopStatement) invoke() (Expression, error) {
	return nil, nil
}

func (n *NopStatement) getType() Type {
	return nopStatementType
}

func (f *ForStatement) invoke() (Expression, error) {
	f.vm.pushStackFrame(false) //make stack frame
	val, err := f.preStatement.invoke()
	if err != nil {
		log.Println("for preStatement.invoke() error", err)
		return nil, err
	}
	if val != nil {
		log.Println("for preStatement.invoke() must nil")
		return nil, fmt.Errorf("for preStatement.invoke() must nil")
	}

	for ; ; {
		val, err := f.checkStatement.invoke()
		if err != nil {
			log.Println("for checkStatement.invoke() error", err)
			return nil, err
		}
		bObj, ok := val.(*BoolObject)
		if ok == false {
			fmt.Errorf("for checkStatement expect BoolObject")
			return nil, fmt.Errorf("for checkStatement expect BoolObject")
		}
		if bObj.val == false {
			f.vm.popStackFrame() //end of for
			return nil, nil
		}
		f.vm.pushStackFrame(false) //make stack frame for `{` brock
		for _, statement := range f.statements {
			val, err := statement.invoke()
			if err != nil {
				log.Println("for checkStatement.invoke() error", err)
				return nil, err
			}
			if val == breakObject {
				log.Println("break from for")
				return nil, nil
			}
			if _, ok := val.(*ReturnStatement); ok {
				return val, nil
			}
		}
		f.vm.popStackFrame()
		if _, err = f.postStatement.invoke(); err != nil {
			log.Println("for postStatement.invoke() error", err)
			return nil, err
		}
	}
}

func (f *ForStatement) getType() Type {
	return forTokenType
}

func (statement *IncFieldStatement) invoke() (Expression, error) {
	object := statement.ctx.getObject(statement.label)
	if object == nil {
		return nil, fmt.Errorf("no find Object with label `%s`", statement.label)
	}
	innerObject, err := object.invoke()
	if err != nil {
		panic(err)
	}
	switch obj := innerObject.(type) {
	case *IntObject:
		obj.val++
		return object, nil
	default:
		panic("unknown type " + reflect.TypeOf(innerObject).String())
	}
	return nil, nil
}

func (statement *IncFieldStatement) getType() Type {
	return incOperatorTokenType
}

func (Statements) getType() Type {
	return statementsType
}

func (f *FuncCallStatement) invoke() (Expression, error) {
	if f.getObject != nil {
		object, err := f.getObject.invoke()
		if err != nil {
			log.Println("getObject.invoke() failed", err)
			return nil, err
		}
		function := object.(*Object).unwrapFunction()
		if function == nil {
			err = fmt.Errorf("no finction objects")
			log.Println(err)
			return nil, err
		}
		return function.invoke(f.arguments...)
	} else {
		function, err := f.vm.getFunction(f.label)
		if err == nil {
			return function.invoke(f.arguments...)
		}
		log.Panic("getFunction failed", f.label, err)
		return nil, err
	}
}

func (f *FuncCallStatement) getType() Type {
	return funcTokenType
}

func (f *getVarStatement) invoke() (Expression, error) {
	object := f.ctx.getObject(f.label)
	if object == nil {
		return nil, fmt.Errorf("no find Object with label `%s`", f.label)
	}
	return object.invoke()
}

func (f *getVarStatement) getType() Type {
	return labelType
}

func (v *VarStatement) invoke() (Expression, error) {
	if v.object != nil {
		v.object.allocObject(v.label)
	} else {
		v.ctx.allocObject(v.label)
	}
	return nil, nil
}

func (v VarStatement) getType() Type {
	return varTokenType
}

func (expression *VarAssignStatement) invoke() (Expression, error) {
	obj, err := expression.expression.invoke()
	if err != nil {
		return nil, err
	}
	var object *Object
	if expression.object != nil {
		object = expression.object.allocObject(expression.label)
	} else {
		object = expression.ctx.allocObject(expression.label)
	}
	switch obj := obj.(type) {
	case *Object:
		object.inner = obj.inner
	default:
		object.inner = obj
	}
	object.initType()
	return nil, nil
}

func (expression *VarAssignStatement) getType() Type {
	return varAssignTokenType
}

func (r *ReturnStatement) invoke() (Expression, error) {
	//log.Println("ReturnStatement invoke")
	if r.returnVal != nil {
		return r, nil
	}
	val, err := r.express.invoke()
	if err != nil {
		log.Println("invoke return statement failed")
		return nil, err
	}
	if val == nil {
		log.Println("return nil error")
		return nil, fmt.Errorf("return expression nil")
	}
	switch inner := val.(type) {
	case *ReturnStatement:
		return inner, nil
	default:
		r.returnVal = val
	}
	return r, nil
}

func (ReturnStatement) getType() Type {
	return returnTokenType
}

func (IfStatement) getType() Type {
	return ifTokenType
}

func (statements Statements) invoke() (Expression, error) {
	var val Expression
	var err error
	for _, statement := range statements {
		if val, err = statement.invoke(); err != nil {
			return nil, err
		} else if val != nil {
			if _, ok := val.(*ReturnStatement); ok {
				return val, nil
			}
		}
	}
	return nil, nil
}

func (ifStm *IfStatement) invoke() (Expression, error) {
	check, err := ifStm.check.invoke()
	if err != nil {
		log.Println("IfStatement check error", err.Error())
		return nil, err
	}
	if _, ok := check.(*BoolObject); ok == false {
		log.Panic("if statement check require boolObject")
	}
	if check.(*BoolObject).val {
		ifStm.vm.pushStackFrame(false) //make  if brock stack
		val, err := ifStm.statement.invoke()
		ifStm.vm.popStackFrame() //release  if brock stack
		if err != nil {
			log.Panic("IfStatement statement error", err.Error())
			return nil, err
		}
		return val, nil
	} else {
		for _, stm := range ifStm.elseIfStatements {
			check2, err := stm.check.invoke()
			if err != nil {
				log.Panic("IfStatement check error", err.Error())
				return nil, err
			}
			if check2.(*BoolObject).val {
				ifStm.vm.pushStackFrame(false) //make  if brock stack
				val, err := stm.statement.invoke()
				ifStm.vm.popStackFrame() //release  if brock stack
				if err != nil {
					log.Panic("IfStatement statement error", err.Error())
					return nil, err
				}
				return val, nil
			} else {
				log.Println("false")
			}
		}
		if ifStm.elseStatement != nil {
			ifStm.vm.pushStackFrame(false) //make  brock stack
			val, err := ifStm.elseStatement.invoke()
			ifStm.vm.popStackFrame() //release  if brock stack
			if err != nil {
				log.Panic(err.Error())
				return nil, err
			}
			return val, nil
		}
	}
	return nil, err
}
