package qp

import (
	"fmt"
	"log"
	"strings"
)

type stackFrame struct {
	stack        []*Object
	stackPointer int
	isolate      bool
}

type memory struct {
	stack       []*Object
	stackFrames []stackFrame
}

func (m *memory) alloc(label string) *Object {
	var obj = new(Object)
	m.stack = append(m.stack, obj)
	obj.pointer = len(m.stack) - 1
	obj.label = label
	return obj
}

func (m *memory) getObject(label string) *Object {
	for index := len(m.stack) - 1; index >= 0; index-- {
		object := m.stack[index]
		if object.label == label {
			return object
		}
	}
	return nil
}

func (m *memory) pushStackFrame(isolate bool) {
	m.stackFrames = append(m.stackFrames, stackFrame{
		stack:        m.stack,
		stackPointer: len(m.stack),
		isolate:      isolate,
	})
	if isolate {
		m.stack = make([]*Object, 0, 32)
	}
}

func (m *memory) popStackFrame() {
	if len(m.stackFrames) == 0 {
		panic("stackFrames empty")
	}
	frame := m.stackFrames[len(m.stackFrames)-1]
	m.stackFrames = m.stackFrames[:len(m.stackFrames)-1]
	var toGc []*Object
	if frame.isolate == false {
		toGc = m.stack[frame.stackPointer:]
		m.stack = m.stack[:frame.stackPointer]
	} else {
		toGc = m.stack
		m.stack = frame.stack
	}
	for i := range toGc {
		toGc[i] = nil
	}
}

type VMContext struct {
	mem           *memory
	functions     map[string]Function
	structObjects map[string]*TypeObject
}

func newVMContext() *VMContext {
	return &VMContext{
		mem:           &memory{},
		structObjects: map[string]*TypeObject{},
		functions:     map[string]Function{},
	}
}

func (ctx *VMContext) allocObject(label string) *Object {
	return ctx.mem.alloc(label)
}

func (ctx *VMContext) getObject(label string) *Object {
	return ctx.mem.getObject(label)
}

func (ctx *VMContext) pushStackFrame(isolate bool) {
	ctx.mem.pushStackFrame(isolate)
}

func (ctx *VMContext) popStackFrame() {
	ctx.mem.popStackFrame()
}

func (ctx *VMContext) addUserFunction(function *FuncStatement) error {
	log.Println("addUserFunction")
	if function.labels != nil {
		structObject := ctx.getTypeObject(function.labels[0])
		if structObject == nil { //todo fix parse order
			fmt.Println("no find structObject", function.labels[0])
			return fmt.Errorf("no find structObject")
		}
		fmt.Println(function.labels)
		structObject.addObject(function.labels[1], &Object{
			inner: function,
			label: strings.Join(function.labels, "."),
			typ:   FuncStatementType,
		})
		return nil
	}

	fmt.Println("addUserFunction with label", function.label)
	if _, ok := builtInFunctionMap[function.label]; ok {
		fmt.Println("function name conflict with built in function", function.label)
		return fmt.Errorf("conflict")
	}
	if _, ok := ctx.functions[function.label]; ok {
		fmt.Println("function name repeated")
		return fmt.Errorf("function name repeated")
	}
	ctx.functions[function.label] = function
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

	if object := ctx.getObject(label); object != nil {
		if function := object.unwrapFuncStatement(); function != nil {
			return function, nil
		}
	}
	fmt.Println("no function ", label)
	return nil, fmt.Errorf("no find function with label `%s`", label)
}

func (ctx *VMContext) addStructObject(object *TypeObject) error {
	if _, ok := ctx.structObjects[object.label]; ok {
		fmt.Println("structObject repeated", object.label)
		return fmt.Errorf("structObject repeated")
	}
	ctx.structObjects[object.label] = object
	return nil
}

func (ctx *VMContext) getTypeObject(label string) *TypeObject {
	obj, _ := ctx.structObjects[label]
	return obj
}

func (ctx *VMContext) cloneTypeObject(label string) *TypeObject {
	obj, ok := ctx.structObjects[label]
	if ok == false {
		fmt.Println("no find structObject with label", label)
		return nil
	}
	return obj.clone()
}
