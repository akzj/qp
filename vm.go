package qp

import (
	"fmt"
	"log"
	"strings"
)

type stackFrame struct {
	stackTopPointer    int
	stackBottomPointer int
	stackGCPointer     int
}

type memory struct {
	stackTopPointer    int
	stackBottomPointer int
	stackGCPointer     int
	stackSize          int
	stack              []Object
	stackFrames        []stackFrame
}

func newMemory() *memory {
	return &memory{
		stackTopPointer:    0,
		stackBottomPointer: 0,
		stackSize:          1024,
		stack:              make([]Object, 1024),
		stackFrames:        make([]stackFrame, 0, 1024),
	}
}

func (m *memory) alloc(label string) *Object {
	if m.stackSize-1 <= m.stackTopPointer {
		newStack := make([]Object, m.stackSize*2)
		copy(newStack, m.stack[:m.stackTopPointer])
		m.stack = newStack
		m.stackSize = m.stackSize * 2
	}
	object := &m.stack[m.stackTopPointer]
	object.label = label
	m.stackTopPointer++
	return object
}

func (m *memory) getObject(label string) *Object {
	for index := m.stackTopPointer - 1; index >= m.stackBottomPointer; index-- {
		if m.stack[index].label == label {
			return &m.stack[index]
		}
	}
	return nil
}

func (m *memory) pushStackFrame(isolate bool) {
	m.stackFrames = append(m.stackFrames, stackFrame{
		stackTopPointer:    m.stackTopPointer,
		stackBottomPointer: m.stackBottomPointer,
		stackGCPointer:     m.stackGCPointer,
	})
	m.stackGCPointer = m.stackTopPointer
	if isolate {
		m.stackBottomPointer = m.stackTopPointer
	}
}

func (m *memory) popStackFrame() {
	if len(m.stackFrames) == 0 {
		panic("stackFrames empty")
	}
	toGc := m.stack[m.stackGCPointer:m.stackTopPointer]
	frame := m.stackFrames[len(m.stackFrames)-1]
	m.stackFrames = m.stackFrames[:len(m.stackFrames)-1]
	m.stackTopPointer = frame.stackTopPointer
	m.stackBottomPointer = frame.stackBottomPointer
	m.stackGCPointer = frame.stackGCPointer
	for i := range toGc {
		toGc[i].inner = nil
	}
}

type VMContext struct {
	mem           *memory
	functions     map[string]*Object
	structObjects map[string]*Object
}

func newVMContext() *VMContext {
	return &VMContext{
		mem:           newMemory(),
		structObjects: map[string]*Object{},
		functions:     map[string]*Object{},
	}
}

func (ctx *VMContext) allocObject(label string) *Object {
	return ctx.mem.alloc(label)
}

func (ctx *VMContext) getObject(label string) *Object {
	if obj, ok := builtInFunctions[label]; ok {
		return &Object{inner: obj}
	}
	if obj, ok := ctx.functions[label]; ok {
		return obj
	}
	if obj, ok := ctx.structObjects[label]; ok {
		return obj
	}
	return ctx.mem.getObject(label)
}

func (ctx *VMContext) pushStackFrame(isolate bool) {
	ctx.mem.pushStackFrame(isolate)
}

func (ctx *VMContext) popStackFrame() {
	ctx.mem.popStackFrame()
}

func (ctx *VMContext) addUserFunction(function *FuncStatement) {
	if function.labels != nil {
		structObject := ctx.getTypeObject(function.labels[0])
		if structObject == nil { //todo fix parse order
			log.Panic("no find structObject", function.labels[0])
		}
		structObject.addObject(function.labels[1], &Object{
			inner: function,
			label: strings.Join(function.labels, "."),
			typ:   FuncStatementType,
		})
	}

	if _, ok := builtInFunctions[function.label]; ok {
		log.Panic("function name conflict with built in function", function.label)
	}
	if _, ok := ctx.functions[function.label]; ok {
		log.Panic("function name repeated")
	}
	ctx.functions[function.label] = &Object{inner: function}
}

func (ctx *VMContext) getFunction(label string) (Function, error) {
	if function, ok := builtInFunctions[label]; ok {
		return function.inner.(Function), nil
	}
	if function, ok := ctx.functions[label]; ok {
		return function.inner.(Function), nil
	}

	if object := ctx.getObject(label); object != nil {
		if function := object.unwrapFunction(); function != nil {
			return function, nil
		}
	}
	return nil, fmt.Errorf("no find function with name `%s`", label)
}

func (ctx *VMContext) addStructObject(object *TypeObject) {
	if _, ok := ctx.structObjects[object.label]; ok {
		log.Panic("structObject repeated", object.label)
	}
	ctx.structObjects[object.label] = &Object{inner: object}
}

func (ctx *VMContext) getTypeObject(label string) *TypeObject {
	obj, _ := ctx.structObjects[label]
	if obj == nil {
		return nil
	}
	return obj.inner.(*TypeObject)
}

func (ctx *VMContext) isGlobal(label string) bool {
	if _, ok := builtInFunctions[label]; ok {
		return true
	}
	return false
}
