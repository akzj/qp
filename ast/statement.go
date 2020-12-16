package ast

import (
	"gitlab.com/akzj/qp"
	"gitlab.com/akzj/qp/lexer"
	"log"
	"reflect"
	"strings"
)

type Statements []Statement

func (statements Statements) String() string {
	var str string
	for index, state := range statements {
		str += state.String()
		if index != len(statements)-1 {
			str += "\n"
		}
	}
	return str
}

type Statement interface {
	qp.Expression
}

type IfStatement struct {
	vm               *qp.VMContext
	check            qp.Expression
	statement        Statements
	elseIfStatements []*IfStatement
	elseStatement    Statements
}

func (ifStm IfStatement) String() string {
	return "if " + ifStm.check.String() + "{}"
}

type ReturnStatement struct {
	express   qp.Expression
	returnVal qp.Expression
}

func (r ReturnStatement) String() string {
	if r.returnVal != nil {
		return "return " + r.returnVal.String()
	} else {
		return "return " + r.express.String()
	}
}

//just new Object
type VarStatement struct {
	ctx    *qp.VMContext
	label  string
	object qp.Expression
}

func (v VarStatement) String() string {
	return "var " + v.label + "=" + v.object.String()
}

type periodStatement struct {
	val string
	exp qp.Expression
}

func (p periodStatement) Invoke() qp.Expression {
	object := p.exp.Invoke().(*qp.Object)
	switch obj := object.inner.(type) {
	case qp.BaseObject:
		return obj.allocObject(p.val)
	default:
		log.Panicf("left `%s` `%s` is no object type", p.val, reflect.TypeOf(obj).String())
	}
	return nil
}

func (p periodStatement) GetType() lexer.Type {
	return lexer.PeriodType
}

func (p periodStatement) String() string {
	return p.exp.String() + "." + p.val
}

type getVarStatement struct {
	ctx   *qp.VMContext
	label string
}

func (f getVarStatement) String() string {
	return f.label
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
	vmContext *qp.VMContext
	labels    []string
}

type FuncCallStatement struct {
	parentExp qp.Expression
	function  qp.Expression
	arguments qp.Expressions
}

func (f *FuncCallStatement) String() string {
	var str = f.function.String() + "("
	for index, statement := range f.arguments {
		if index != 0 {
			str += ","
		}
		str += statement.String()
	}
	return str + ")"
}

type AssignStatement struct {
	exp  qp.Expression
	left qp.Expression
}

func (expression AssignStatement) String() string {
	return expression.left.String() + "=" + expression.exp.String()
}

type VarAssignStatement struct {
	ctx        *qp.VMContext //global or stack var
	name       string        //var name : var a,`a` is the name
	expression qp.Expression // init expression : var a = 1+1
}

func (expression VarAssignStatement) String() string {
	return "var " + expression.name + "=" + expression.expression.String()
}

type IncFieldStatement struct {
	exp qp.Expression
}

func (statement IncFieldStatement) String() string {
	panic("implement me")
}

type BreakStatement struct {
}

var nopStatement = NopStatement{}

type NopStatement struct {
}

func (n NopStatement) String() string {
	return "nop"
}

type FuncStatement struct {
	closure      bool
	label        string
	labels       []string // struct objects function eg:user.add(){}
	parameters   []string // parameter label
	closureLabel []string // closure label
	closureInit  bool
	statements   Statements    // function body
	vm           *qp.VMContext // vm context
	closureObjs  []qp.Expression
}

func (f *FuncStatement) String() string {
	var str = "func " + f.label + "("
	for index, argument := range f.parameters {
		if index != 0 {
			str += ","
		}
		str += argument
	}
	str += "){\n"
	for _, statement := range f.statements {
		str += statement.String() + "\n"
	}
	str += "}"
	return str
}

func (f *FuncStatement) Invoke() qp.Expression {
	f.doClosureInit()
	return f
}

type ForStatement struct {
	vm             *qp.VMContext
	preStatement   qp.Expression
	checkStatement qp.Expression
	postStatement  qp.Expression
	statements     Statements
}

func (f *ForStatement) String() string {
	return "for"
}

type objectInitStatement struct {
	exp           qp.Expression
	vm            *qp.VMContext
	propTemplates []qp.TypeObjectPropTemplate
}

func (statement *objectInitStatement) String() string {
	var str string
	for _, statement := range statement.propTemplates {
		str += statement.String() + "\n"
	}
	return "{" + str + "}"
}

type getArrayElement struct {
	arrayExp qp.Expression
	indexExp qp.Expression
}

func (g getArrayElement) Invoke() qp.Expression {
	panic("implement me")
}

func (g getArrayElement) GetType() lexer.Type {
	panic("implement me")
}

func (g getArrayElement) String() string {
	panic("implement me")
}

type makeArrayStatement struct {
	vm             *qp.VMContext
	initStatements Statements
}

func (m *makeArrayStatement) String() string {
	var str = "["
	for index, statement := range m.initStatements {
		if index != 0 {
			str += ","
		}
		str += statement.String()
	}
	return str + "]"
}

func (m *makeArrayStatement) Invoke() qp.Expression {
	var array = &qp.Array{}
	for _, statement := range m.initStatements {
		array.data = append(array.data, statement.Invoke())
	}
	return array
}

func (m *makeArrayStatement) GetType() lexer.Type {
	return lexer.ArrayObjectType
}

func (g *getObjectPropStatement) Invoke() qp.Expression {
	obj := g.getObject.Invoke()
	if obj == qp.nilObject {
		return obj
	}
	return obj.(*qp.Object).inner
}

func (g *getObjectObjectStatement) Invoke() qp.Expression {
	object := g.vmContext.getObject(g.labels[0])
	if object == nil {
		log.Panicf("left failed `%s`", g.labels[0])
	}
	structObj, ok := object.inner.(qp.BaseObject)
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
		//last name
		if i != len(g.labels)-1 {
			var ok bool
			structObj, ok = obj.inner.(*qp.TypeObject)
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

func (statement *objectInitStatement) Invoke() qp.Expression {
	object := statement.exp.Invoke().(*qp.Object).inner.(qp.BaseObject).clone().(*qp.TypeObject)

Loop:
	for _, init := range object.typeObjectPropTemplates {
		for _, prod := range statement.propTemplates {
			if init.name == prod.name {
				continue Loop
			}
		}
		propObject := object.allocObject(init.name)
		propObject.inner = init.exp.Invoke()
	}

	for _, init := range statement.propTemplates {
		propObject := object.allocObject(init.name)
		propObject.inner = init.exp.Invoke()
	}
	return object
}

func (statement *objectInitStatement) GetType() lexer.Type {
	return lexer.TypeObjectInitStatementType
}

func (f *FuncStatement) prepareArgumentBind(inArguments qp.Expressions) {
	if len(f.parameters) != len(inArguments) {
		if f.closure {
		}
		log.Panicf("call function %s argument count %d %d no match ", f.label, len(f.parameters), len(inArguments))
	}

	f.vm.pushStackFrame(true)
	for index := range f.closureLabel {
		// put closure objects to stack
		f.vm.allocObject(f.closureLabel[index]).inner = f.closureObjs[index]
	}

	//make stack for this function
	for index, result := range inArguments {
		f.vm.allocObject(f.parameters[index]).inner = result
	}
}

func (f *FuncStatement) call(arguments ...qp.Expression) qp.Expression {
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

func (f *FuncStatement) GetType() lexer.Type {
	return lexer.FuncStatementType
}

func (f *FuncStatement) doClosureInit() {
	if f.closureInit {
		return
	}
	f.closureInit = true
	var closureObjs []qp.Expression
	var closureLabel []string
	for _, label := range f.closureLabel {
		if f.vm.isGlobal(label) {
			continue
		}
		obj := f.vm.getObject(label)
		if obj == nil {
			log.Panicf("no find obj with name `%s`", label)
		}
		closureObjs = append(closureObjs, obj.inner)
		closureLabel = append(closureLabel, label)
	}
	f.closureObjs = closureObjs
	f.closureLabel = closureLabel
}

func (expression AssignStatement) Invoke() qp.Expression {
	left := expression.left.Invoke()
	switch right := expression.exp.Invoke().(type) {
	case *qp.Object:
		left.(*qp.Object).inner = right.inner
	default:
		left.(*qp.Object).inner = right
	}
	return nil
}

func (expression AssignStatement) GetType() lexer.Type {
	return lexer.AssignStatementType
}

func (NopStatement) Invoke() qp.Expression {
	return nopStatement
}

func (n NopStatement) GetType() lexer.Type {
	return lexer.NopStatementType
}

func (f *ForStatement) Invoke() qp.Expression {
	f.vm.pushStackFrame(false) //make stack frame

	//make for brock stack
	f.preStatement.Invoke()

	for ; ; {
		val := f.checkStatement.Invoke()
		bObj, ok := val.(qp.Bool)
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
			if val == qp.breakObject {
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

func (f *ForStatement) GetType() lexer.Type {
	return lexer.ForType
}

func (statement IncFieldStatement) Invoke() qp.Expression {
	object := statement.exp.Invoke().(*qp.Object)
	object.inner = object.inner.(qp.Int) + 1
	return nil
}

func (statement IncFieldStatement) GetType() lexer.Type {
	return lexer.IncType
}

func (Statements) GetType() lexer.Type {
	return lexer.StatementsType
}

func (f *FuncCallStatement) Invoke() qp.Expression {
	exp := f.function.Invoke()
	switch obj := exp.(type) {
	case *qp.Object:
		exp = obj.Invoke()
	case ReturnStatement:
		exp = obj.returnVal
	}
	if exp == nil {
		log.Panic("function nil")
	}
	var arguments []qp.Expression
	if Func, ok := exp.(*FuncStatement);
		f.parentExp != nil && (ok == false || Func.closure == false) {
		switch argument := f.parentExp.Invoke().(type) {
		case *qp.Object:
			arguments = append(arguments, argument.inner)
		default:
			arguments = append(arguments, argument)
		}
	}

	if function, ok := exp.(qp.Function); ok {
		for _, argument := range f.arguments {
			switch job := argument.Invoke().(type) {
			case *qp.Object:
				arguments = append(arguments, job.inner)
			default:
				arguments = append(arguments, job)
			}
		}
		return function.call(arguments...)
	}
	log.Panicf("object`%s` `%s` is no callable", exp.String(), reflect.TypeOf(exp).String())
	return nil
}

func (f *FuncCallStatement) GetType() lexer.Type {
	return lexer.FuncType
}

func (f getVarStatement) Invoke() qp.Expression {
	return f.ctx.getObject(f.label)
}

func (f getVarStatement) GetType() lexer.Type {
	return lexer.IDType
}

func (v VarStatement) Invoke() qp.Expression {
	if v.object != nil {
		v.ctx.allocObject(v.label).inner = v.object.Invoke()
	} else {
		v.ctx.allocObject(v.label).inner = qp.nilObject
	}
	return nil
}

func (v VarStatement) GetType() lexer.Type {
	return lexer.VarType
}

func (expression VarAssignStatement) Invoke() qp.Expression {
	obj := expression.expression.Invoke()
	var object = expression.ctx.allocObject(expression.name)
	if obj == nil {
		panic(obj)
	}
	switch obj := obj.(type) {
	case *qp.Object:
		object.inner = obj.inner
	default:
		object.inner = obj
	}
	object.initType()
	return nil
}

func (expression VarAssignStatement) GetType() lexer.Type {
	return lexer.VarAssignType
}

func (r ReturnStatement) Invoke() qp.Expression {
	if r.returnVal != nil {
		return r
	}
	exp := r.express.Invoke()
	switch obj := exp.(type) {
	case *qp.Object:
		exp = obj.inner
	case ReturnStatement:
		return obj
	}
	return ReturnStatement{returnVal: exp}
}

func (ReturnStatement) GetType() lexer.Type {
	return lexer.ReturnType
}

func (IfStatement) GetType() lexer.Type {
	return lexer.IfType
}

func (statements Statements) Invoke() qp.Expression {
	var val qp.Expression
	for _, statement := range statements {
		val = statement.Invoke()
		if _, ok := val.(ReturnStatement); ok {
			return val
		} else if _, ok := val.(*qp.BreakObject); ok {
			return qp.breakObject
		}
	}
	return val
}

func (ifStm *IfStatement) Invoke() qp.Expression {
	check := ifStm.check.Invoke()
	if _, ok := check.(qp.Bool); ok == false {
		log.Panic("if statement check require boolObject", reflect.TypeOf(check).String())
	}
	if check.(qp.Bool) {
		ifStm.vm.pushStackFrame(false) //make  if brock stack
		val := ifStm.statement.Invoke()
		ifStm.vm.popStackFrame() //release  if brock stack
		return val
	} else {
		for _, stm := range ifStm.elseIfStatements {
			elseIf := stm.check.Invoke()
			if _, ok := elseIf.(qp.Bool); ok == false {
				log.Panicln("else if require bool result")
			}
			if elseIf.(qp.Bool) {
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
