package qp

import (
	"fmt"
	"reflect"
	"strings"
)

type Statements []Statement

type Statement interface {
	Expression
}

type IfStatement struct {
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
type getStructObjectStatement struct {
	this      bool
	vmContext *VMContext
	labels    []string
}

type FuncCallStatement struct {
	vm        *VMContext
	label     string
	getObject *getStructObjectStatement
	arguments Expressions
}

type AssignStatement struct {
	ctx        *VMContext
	label      string
	expression Expression
	getObject  *getStructObjectStatement
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
	label      string
	labels     []string   // struct object function eg:user.add(){}
	parameters []string   // parameter label
	statements Statements // function body
	vm         *VMContext // vm context
}

type ForStatement struct {
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

func (g *getStructObjectStatement) invoke() (Expression, error) {
	fmt.Println("getStructObjectStatement", g.labels)
	object := g.vmContext.getObject(g.labels[0])
	if object == nil {
		fmt.Println("getObject failed", g.labels[0])
		return nil, fmt.Errorf("getObject failed")
	}
	if object.inner == nil {
		fmt.Println("object nil", g.labels[0])
		return nil, fmt.Errorf("object nil")
	}
	structObj, ok := object.inner.(*TypeObject)
	if ok == false {
		fmt.Println("object type no struct object,error",
			g.labels, reflect.TypeOf(object.inner).String())
		return nil, fmt.Errorf("object type is no TypeObject")
	}
	/*
	 user.id = 1 // bind 1 to user.id
	 println(user.id)// visit user.id
	*/
	var obj = object
	for i := 1; i < len(g.labels); i++ {
		obj = structObj.allocObject(g.labels[i])
		fmt.Println("g.labels[i]", g.labels[i])
		//last label
		if i != len(g.labels)-1 {
			var ok bool
			structObj, ok = obj.inner.(*TypeObject)
			if ok == false {
				label := strings.Join(g.labels[:i+1], ".")
				fmt.Println("object is no struct object type", label)
				return nil, fmt.Errorf("object(%s) type is no TypeObject", label)
			}
		}
	}
	fmt.Println("o------------------------0", reflect.TypeOf(obj).String())
	return obj, nil
}

func (g *getStructObjectStatement) getType() Type {
	return getTypeObjectVarStatementType
}

func (statement *StructObjectInitStatement) invoke() (Expression, error) {
	fmt.Println("StructObjectInitStatement invoke()", statement.label)
	object := statement.vm.allocStructObject(statement.label)
	if object == nil {
		return nil, fmt.Errorf("allocStructObject with label `%s` failed", statement.label)
	}
	fmt.Println("object.initStatement", len(object.initStatement))
	for _, initStatement := range statement.initStatements {
		object.initStatement = append(object.initStatement, initStatement)
	}
	for _, init := range object.initStatement {
		switch s := init.(type) {
		case *VarAssignStatement:
			s.object = object
			fmt.Println("----s----")
		case *VarStatement:
			s.object = object
			fmt.Println("----sx----")
		case *NopStatement:
			continue
		default:
			panic("unknown statement " + reflect.TypeOf(init).String())
		}
		if _, err := init.invoke(); err != nil {
			fmt.Println("struct object init failed")
			return nil, err
		}
	}
	return object, nil
}

func (statement *StructObjectInitStatement) getType() Type {
	return typeObjectInitStatementType
}

func (f *FuncStatement) prepareArgumentBind(inArguments Expressions) error {

	/*
		lambda function no bind this to object
	*/
	if f.label == "lambda" && len(inArguments) != 0 {
		statement, ok := inArguments[0].(*getStructObjectStatement)
		if ok && statement.this {
			inArguments = inArguments[1:]
		}
	}
	if len(f.parameters) != len(inArguments) {
		fmt.Println("argument size no match", len(f.parameters), len(inArguments))
		return fmt.Errorf("argument size no match")
	}

	for index, expression := range inArguments {
		label := f.parameters[index]
		fmt.Println(expression.getType(), label)
		val, err := expression.invoke()
		if err != nil {
			fmt.Println("invoke argument failed", err)
			return err
		}
		if val == nil {
			fmt.Println("invoke argument return nil error")
			return fmt.Errorf("invoke argument return nil error")
		}
		object := f.vm.allocObject(label)
		if object == nil {
			panic("allocObject nil")
		}
		switch obj := val.(type) {
		case *Object:
			if obj == nil {
				panic("obj nil")
			}
			object.inner = obj.inner
		default:
			object.inner = val
		}
		fmt.Println("---------bind argument-----",
			label, reflect.ValueOf(val).String())
	}
	return nil
}

func (f *FuncStatement) invoke(arguments ...Expression) (Expression, error) {
	//argument stack
	f.vm.pushStackFrame()

	//pop argument stack
	defer f.vm.popStackFrame()
	if err := f.prepareArgumentBind(arguments); err != nil {
		return nil, err
	}
	for _, statement := range f.statements {
		val, err := statement.invoke()
		if err != nil {
			fmt.Println("statement.invoke() failed", err)
			return nil, err
		}
		if val != nil {
			//function return
			if _, ok := val.(*ReturnStatement); ok {
				return val, nil
			}
		}
	}
	return nil, nil
}

func (f *FuncStatement) getType() Type {
	panic("implement me")
}

func (expression *AssignStatement) invoke() (Expression, error) {
	fmt.Println("AssignStatement")
	val, err := expression.expression.invoke()
	if err != nil {
		fmt.Println("AssignStatement .expression.invoke() failed", err.Error())
		return nil, err
	}
	if expression.getObject != nil {
		object, err := expression.getObject.invoke()
		if err != nil {
			fmt.Println("on.getObject.invoke() failed", err)
			return nil, err
		}
		object.(*Object).inner = val
	} else {
		fmt.Println(val.getType())
		fmt.Println(val.(*IntObject).val)
		object := expression.ctx.getObject(expression.label)
		if object == nil {
			fmt.Println("AssignStatement .expression.getObject failed", object.label)
			return nil, err
		}
		object.inner = val
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
	val, err := f.preStatement.invoke()
	if err != nil {
		fmt.Println("for preStatement.invoke() error", err)
		return nil, err
	}
	if val != nil {
		fmt.Println("for preStatement.invoke() must nil")
		return nil, fmt.Errorf("for preStatement.invoke() must nil")
	}

	for ; ; {
		val, err := f.checkStatement.invoke()
		if err != nil {
			fmt.Println("for checkStatement.invoke() error", err)
			return nil, err
		}
		bObj, ok := val.(*BoolObject)
		if ok == false {
			fmt.Errorf("for checkStatement expect BoolObject")
			return nil, fmt.Errorf("for checkStatement expect BoolObject")
		}
		if bObj.val == false {
			return nil, nil
		}
		for _, statement := range f.statements {
			val, err := statement.invoke()
			if err != nil {
				fmt.Println("for checkStatement.invoke() error", err)
				return nil, err
			}
			if val == breakObject {
				fmt.Println("break from for")
				return nil, nil
			}
			if _, ok := val.(*ReturnStatement); ok {
				return val, nil
			}
		}
		if _, err = f.postStatement.invoke(); err != nil {
			fmt.Println("for postStatement.invoke() error", err)
			return nil, err
		}
	}
}

func (f *ForStatement) getType() Type {
	return forTokenType
}

func (statement *IncFieldStatement) invoke() (Expression, error) {
	fmt.Println("IncFieldStatement")
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
	fmt.Println("FuncCallStatement invoke")
	if f.getObject != nil {
		object, err := f.getObject.invoke()
		if err != nil {
			fmt.Println("getObject.invoke() failed", err)
			return nil, err
		}
		fmt.Println("unwrapFuncStatement")
		function := object.(*Object).unwrapFuncStatement()
		if function == nil {
			err = fmt.Errorf("no finction object")
			fmt.Println(err)
			return nil, err
		}
		return function.invoke(f.arguments...)
	} else {
		function, err := f.vm.getFunction(f.label)
		if err == nil {
			return function.invoke(f.arguments...)
		}
		fmt.Println("getFunction failed", f.label, err)
		return nil, err
	}
}

func (f *FuncCallStatement) getType() Type {
	return funcTokenType
}

func (f *getVarStatement) invoke() (Expression, error) {
	fmt.Println("getVarStatement invoke")
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
	fmt.Println("VarStatement invoke")
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
	fmt.Println("VarAssignStatement invoke", expression.label)
	obj, err := expression.expression.invoke()
	if err != nil {
		return nil, err
	}
	fmt.Println("expression", expression.expression.getType())
	fmt.Println(obj.getType())
	var object *Object
	if expression.object != nil {
		fmt.Println("expression.label", expression.label)
		object = expression.object.allocObject(expression.label)
	} else {
		object = expression.ctx.allocObject(expression.label)
	}
	object.inner = obj
	object.initType()
	return nil, nil
}

func (expression *VarAssignStatement) getType() Type {
	return varAssignTokenType
}

func (r *ReturnStatement) invoke() (Expression, error) {
	//fmt.Println("ReturnStatement invoke")
	if r.returnVal != nil {
		return r, nil
	}
	val, err := r.express.invoke()
	if err != nil {
		fmt.Println("invoke return statement failed")
		return nil, err
	}
	if val == nil {
		fmt.Println("return nil error")
		return nil, fmt.Errorf("return expression nil")
	}
	fmt.Println("return val", val)
	r.returnVal = val
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
		fmt.Println("statement type", statement.getType())
		if val, err = statement.invoke(); err != nil {
			return nil, err
		} else if val != nil {
			if _, ok := val.(*ReturnStatement); ok {
				return val, nil
			} else {
				fmt.Println("---------", val.getType())
			}
		}
	}
	return nil, nil
}

func (ifStm *IfStatement) invoke() (Expression, error) {
	//fmt.Println("IfStatement invoke")
	check, err := ifStm.check.invoke()
	if err != nil {
		fmt.Println("IfStatement check error", err.Error())
		return nil, err
	}
	if check.(*BoolObject).val {
		fmt.Println("true")
		val, err := ifStm.statement.invoke()
		if err != nil {
			fmt.Println("IfStatement statement error", err.Error())
			return nil, err
		}
		fmt.Println(ifStm.statement.getType())
		return val, nil
	} else {
		fmt.Println(ifStm.check)
		for _, stm := range ifStm.elseIfStatements {
			check2, err := stm.check.invoke()
			if err != nil {
				fmt.Println("IfStatement check error", err.Error())
				return nil, err
			}
			if check2.(*BoolObject).val {
				fmt.Println("else if true")
				val, err := stm.statement.invoke()
				if err != nil {
					fmt.Println("IfStatement statement error", err.Error())
					return nil, err
				}
				return val, nil
			} else {
				fmt.Println("false")
			}
		}
		if ifStm.elseStatement != nil {
			fmt.Println("else")
			val, err := ifStm.elseStatement.invoke()
			if err != nil {
				fmt.Println(err.Error())
				return nil, err
			}
			return val, nil
		}
	}
	return nil, err
}
