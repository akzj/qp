package qp

import "fmt"

type memory struct {
	objects      []*Object
	stackPointer []int
}

func (m *memory) alloc(label string) *Object {
	var obj = new(Object)
	m.objects = append(m.objects, obj)
	obj.pointer = len(m.objects) - 1
	obj.label = label
	return obj
}

func (m *memory) getObject(label string) *Object {
	for index := len(m.objects) - 1; index >= 0; index-- {
		object := m.objects[index]
		if object.label == label {
			return object
		}
	}
	return nil
}

func (m *memory) pushStackFrame() {
	m.stackPointer = append(m.stackPointer, len(m.objects))
}
func (m *memory) popStackFrame() {
	if len(m.stackPointer) == 0 {
		panic("stack empty")
	}
	pointer := m.stackPointer[len(m.stackPointer)-1]
	m.stackPointer = m.stackPointer[:len(m.stackPointer)-1]
	m.objects = m.objects[:pointer]
}

type VMContext struct {
	mem       *memory
	functions map[string]Function
}

func newVMContext() *VMContext {
	return &VMContext{
		mem:       &memory{},
		functions: map[string]Function{},
	}
}

func (ctx *VMContext) allocObject(label string) *Object {
	return ctx.mem.alloc(label)
}

func (ctx *VMContext) getObject(label string) *Object {
	return ctx.mem.getObject(label)
}

func (ctx *VMContext) pushStackFrame() {
	ctx.mem.pushStackFrame()
}

func (ctx *VMContext) popStackFrame() {
	ctx.mem.popStackFrame()
}

func (ctx *VMContext) addUserFunction(function *FuncStatement) error {
	fmt.Println("addUserFunction with label", function.label)
	if _, ok := builtInFunctionMap[function.label]; ok {
		fmt.Println("function name conflict with built in function", function.label)
		return fmt.Errorf("conflict")
	}
	if _, ok := ctx.functions[function.label]; ok {
		fmt.Println("function name repeated")
		return fmt.Errorf("function name repeated")
	}
	ctx.functions[function.label] = func(arguments ...Expression) (Expression, error) {
		return function.invoke(arguments...)
	}
	return nil
}

func (ctx *VMContext) getFunction(label string) (Function, error) {
	function, ok := builtInFunctionMap[label]
	if ok {
		return function, nil
	}
	function, ok = ctx.functions[label]
	if ok {
		return function, nil
	}
	fmt.Println("no function ", label)
	return nil, fmt.Errorf("no find function with label `%s`", label)
}
