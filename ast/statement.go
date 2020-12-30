package ast

import (
	"gitlab.com/akzj/qp/lexer"
	"gitlab.com/akzj/qp/runtime"
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
	runtime.Invokable
}

type IfStatement struct {
	VM         *runtime.VMContext
	Check      runtime.Invokable
	Statements Statements
	ElseIf     []IfStatement
	Else       Statements
}

func (ifStm IfStatement) String() string {
	return "if " + ifStm.Check.String() + "{}"
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

//just new Object
type VarStatement struct {
	VM    *runtime.VMContext
	Label string
	Exp   runtime.Invokable
}

func (v VarStatement) String() string {
	return "var " + v.Label + "=" + v.Exp.String()
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
	VM    *runtime.VMContext
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
	vmContext *runtime.VMContext
	labels    []string
}

type CallStatement struct {
	ParentExp runtime.Invokable
	Function  runtime.Invokable
	Arguments Statements
}

func (f *CallStatement) String() string {
	var str = f.Function.String() + "("
	for index, statement := range f.Arguments {
		if index != 0 {
			str += ","
		}
		str += statement.String()
	}
	return str + ")"
}

type AssignStatement struct {
	Exp  runtime.Invokable
	Left runtime.Invokable
}

func (expression AssignStatement) String() string {
	return expression.Left.String() + "=" + expression.Exp.String()
}

type VarAssignStatement struct {
	Ctx  *runtime.VMContext //global or stack var
	Name string             //var Name : var a,`a` is the Name
	Exp  runtime.Invokable  // Init Exp : var a = 1+1
}

func (expression VarAssignStatement) String() string {
	return "var " + expression.Name + "=" + expression.Exp.String()
}

type IncFieldStatement struct {
	Exp runtime.Invokable
}

func (statement IncFieldStatement) String() string {
	return statement.Exp.String()+"++"
}

type BreakStatement struct {
}

type NopStatement struct {
}

func (n NopStatement) String() string {
	return "nop"
}

type FuncStatement struct {
	Closure      bool
	Label        string
	Labels       []string // struct objects Function eg:user.add(){}
	Parameters   []string // parameter Label
	ClosureLabel []string // Closure Label
	ClosureInit  bool
	Statements   Statements         // Function body
	VM           *runtime.VMContext // VM context
	ClosureObjs  []runtime.Invokable
}

func (f *FuncStatement) String() string {
	var str = "func " + f.Label + "("
	for index, argument := range f.Parameters {
		if index != 0 {
			str += ","
		}
		str += argument
	}
	str += "){\n"
	for _, statement := range f.Statements {
		str += statement.String() + "\n"
	}
	str += "}"
	return str
}

func (f *FuncStatement) Invoke() runtime.Invokable {
	f.doClosureInit()
	return f
}

type ForStatement struct {
	VM         *runtime.VMContext
	Pre        runtime.Invokable
	Check      runtime.Invokable
	Post       runtime.Invokable
	Statements Statements
}

func (f ForStatement) String() string {
	return "for"
}

type ObjectInitStatement struct {
	VM            *runtime.VMContext
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
	vm    *runtime.VMContext
	Inits Statements
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

func (f *FuncStatement) prepareArgumentBind(inArguments []runtime.Invokable) {
	if len(f.Parameters) != len(inArguments) {
		if f.Closure {
		}
		log.Panicf("call Function %s argument count %d %d no match ", f.Label, len(f.Parameters), len(inArguments))
	}

	f.VM.PushStackFrame(true)
	for index := range f.ClosureLabel {
		// put Closure objects to stack
		f.VM.AllocObject(f.ClosureLabel[index]).Pointer = f.ClosureObjs[index]
	}

	//make stack for this Function
	for index, result := range inArguments {
		f.VM.AllocObject(f.Parameters[index]).Pointer = result
	}
}

func (f *FuncStatement) Call(arguments ...runtime.Invokable) runtime.Invokable {
	defer f.VM.PopStackFrame()
	f.prepareArgumentBind(arguments)
	for _, statement := range f.Statements {
		result := statement.Invoke()
		if ret, ok := result.(ReturnStatement); ok {
			return ret.Val
		}
	}
	return nil
}

func (f *FuncStatement) GetType() lexer.Type {
	return lexer.FuncStatementType
}

func (f *FuncStatement) doClosureInit() {
	if f.ClosureInit {
		return
	}
	f.ClosureInit = true
	var closureObjs []runtime.Invokable
	var closureLabel []string
	for _, label := range f.ClosureLabel {
		if f.VM.IsGlobal(label) {
			continue
		}
		obj := f.VM.GetObject(label)
		if obj == nil {
			log.Panicf("no find obj with Name `%s`", label)
		}
		closureObjs = append(closureObjs, obj.Pointer)
		closureLabel = append(closureLabel, label)
	}
	f.ClosureObjs = closureObjs
	f.ClosureLabel = closureLabel
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

func (f ForStatement) Invoke() runtime.Invokable {
	f.VM.PushStackFrame(false) //make stack frame

	//make for brock stack
	f.Pre.Invoke()

	for ; ; {
		val := f.Check.Invoke()
		bObj, ok := val.(Bool)
		if ok == false {
			log.Panic("for Check expect Bool")
		}
		if bObj == false {
			f.VM.PopStackFrame() //end of for
			return nil
		}
		f.VM.PushStackFrame(false) //make stack frame for `{` brock
		for _, statement := range f.Statements {
			val := statement.Invoke()
			if val == BreakObj {
				return nil
			}
			if _, ok := val.(ReturnStatement); ok {
				return val
			}
		}
		f.VM.PopStackFrame()
		f.Post.Invoke()
	}
}

func (f ForStatement) GetType() lexer.Type {
	return lexer.ForType
}

func (statement IncFieldStatement) Invoke() runtime.Invokable {
	object := statement.Exp.Invoke().(*runtime.Object)
	object.Pointer = object.Pointer.(Int) + 1
	return nil
}

func (statement IncFieldStatement) GetType() lexer.Type {
	return lexer.IncType
}

func (Statements) GetType() lexer.Type {
	return lexer.StatementsType
}

func (f *CallStatement) Invoke() runtime.Invokable {
	exp := f.Function.Invoke()
	switch obj := exp.(type) {
	case *runtime.Object:
		exp = obj.Invoke()
	case ReturnStatement:
		exp = obj.Val
	}
	if exp == nil {
		log.Panic("Function nil")
	}
	var arguments []runtime.Invokable
	if Func, ok := exp.(*FuncStatement);
		f.ParentExp != nil && (ok == false || Func.Closure == false) {
		switch argument := f.ParentExp.Invoke().(type) {
		case *runtime.Object:
			if argument.Pointer == nil {
				panic(argument.Label)
			}
			arguments = append(arguments, argument.Pointer)
		default:
			if argument == nil {
				panic("argument nil")
			}
			arguments = append(arguments, argument)
		}
	}

	if function, ok := exp.(Function); ok {
		for _, argument := range f.Arguments {
			switch job := argument.Invoke().(type) {
			case *runtime.Object:
				if job.Pointer == nil {
					panic(job.Label + " " + f.Function.String())
				}
				arguments = append(arguments, job.Pointer)
			default:
				if job == nil {
					panic("argument nil")
				}
				arguments = append(arguments, job)
			}
		}
		return function.Call(arguments...)
	}
	log.Panicf("Exp`%s` `%s` is no callable", exp.String(), reflect.TypeOf(exp).String())
	return nil
}

func (f *CallStatement) GetType() lexer.Type {
	return lexer.CallType
}

func (f GetVarStatement) Invoke() runtime.Invokable {
	return f.VM.GetObject(f.Label)
}

func (f GetVarStatement) GetType() lexer.Type {
	return lexer.IDType
}

func (v VarStatement) Invoke() runtime.Invokable {
	if v.Exp != nil {
		v.VM.AllocObject(v.Label).Pointer = v.Exp.Invoke()
	} else {
		v.VM.AllocObject(v.Label).Pointer = NilObj
	}
	return nil
}

func (v VarStatement) GetType() lexer.Type {
	return lexer.VarType
}

func (expression VarAssignStatement) Invoke() runtime.Invokable {
	obj := expression.Exp.Invoke()
	var object = expression.Ctx.AllocObject(expression.Name)
	if obj == nil {
		panic(obj)
	}
	switch obj := obj.(type) {
	case *runtime.Object:
		object.Pointer = obj.Pointer
	default:
		object.Pointer = obj
	}
	return nil
}

func (expression VarAssignStatement) GetType() lexer.Type {
	return lexer.VarAssignType
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

func (IfStatement) GetType() lexer.Type {
	return lexer.IfType
}

func (statements Statements) Invoke() runtime.Invokable {
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

func (ifStm IfStatement) Invoke() runtime.Invokable {
	check := ifStm.Check.Invoke()
	if _, ok := check.(Bool); ok == false {
		log.Panic("if Statements Check require boolObject", reflect.TypeOf(check).String())
	}
	if check.(Bool) {
		ifStm.VM.PushStackFrame(false) //make  if brock stack
		val := ifStm.Statements.Invoke()
		ifStm.VM.PopStackFrame() //release  if brock stack
		return val
	} else {
		for _, stm := range ifStm.ElseIf {
			elseIf := stm.Check.Invoke()
			if _, ok := elseIf.(Bool); ok == false {
				log.Panicln("else if require bool result")
			}
			if elseIf.(Bool) {
				ifStm.VM.PushStackFrame(false) //make  if brock stack
				val := stm.Statements.Invoke()
				ifStm.VM.PopStackFrame() //release  if brock stack
				return val
			}
		}
		if ifStm.Else != nil {
			ifStm.VM.PushStackFrame(false) //make  brock stack
			val := ifStm.Else.Invoke()
			ifStm.VM.PopStackFrame() //release  if brock stack
			return val
		}
	}
	return nil
}
