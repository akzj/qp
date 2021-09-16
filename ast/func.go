package ast

import (
	"log"
	"strings"

	"gitlab.com/akzj/qp/lexer"
	"gitlab.com/akzj/qp/runtime"
)

type FuncExpression struct {
	Closure      bool
	Label        string
	Labels       []string // struct objects Function eg:user.add(){}
	Parameters   []string // parameter Label
	ClosureLabel []string // Closure Label
	ClosureInit  bool
	Statements   Expressions        // Function body
	VM           *runtime.VMRuntime // VM context
	ClosureObjs  []runtime.Invokable
}

func (f *FuncExpression) String() string {
	var str = "func " + f.Label + "("
	for index, argument := range f.Parameters {
		if index != 0 {
			str += ","
		}
		str += argument
	}
	str += "){\n"
	for _, statement := range f.Statements {
		for _, line := range strings.Split(statement.String(), "\n") {
			str += "\t" + line + "\n"
		}
	}
	str += "\n}"
	return str
}

func (f *FuncExpression) Invoke() runtime.Invokable {
	f.doClosureInit()
	return f
}

func (f *FuncExpression) prepareArgumentBind(inArguments []runtime.Invokable) {
	if len(f.Parameters) != len(inArguments) {
		log.Panicf("call Function %s argument count %d %d no match ", f.Label, len(f.Parameters), len(inArguments))
	}

	for index := range f.ClosureLabel {
		// put Closure objects to stack
		f.VM.AllocObject(f.ClosureLabel[index]).Pointer = f.ClosureObjs[index]
	}
	//make stack for this Function
	for index, result := range inArguments {
		f.VM.AllocObject(f.Parameters[index]).Pointer = result
	}
}

func (f *FuncExpression) Call(arguments ...runtime.Invokable) runtime.Invokable {
	f.VM.PushStackFrame(true)
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

func (f *FuncExpression) GetType() lexer.Type {
	return lexer.FuncStatementType
}

func (f *FuncExpression) doClosureInit() {
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
