package runtime

import (
	"log"
)

var (
	Functions       = map[string]*Object{}
	ArrayFunctions  = map[string]*Object{}
	StringFunctions = map[string]*Object{}
)

type StackFrame struct {
	stackTopPointer    int
	stackBottomPointer int
	stackGCPointer     int
}

type Memory struct {
	stackTopPointer    int
	stackBottomPointer int
	stackGCPointer     int
	stackSize          int
	stack              []*Object
	stackFrames        []StackFrame
}

func NewMemory() *Memory {
	return &Memory{
		stackTopPointer:    0,
		stackBottomPointer: 0,
		stackSize:          1024,
		stack:              make([]*Object, 1024),
		stackFrames:        make([]StackFrame, 0, 1024),
	}
}

func (m *Memory) Alloc(label string) *Object {
	if m.stackSize-1 <= m.stackTopPointer {
		newStack := make([]*Object, m.stackSize*2)
		copy(newStack, m.stack[:m.stackTopPointer])
		m.stack = newStack
		m.stackSize = m.stackSize * 2
	}
	object := &Object{}
	m.stack[m.stackTopPointer] = object
	object.Label = label
	m.stackTopPointer++
	return object
}

func (m *Memory) GetObject(label string) *Object {
	for index := m.stackTopPointer - 1; index >= m.stackBottomPointer; index-- {
		if m.stack[index].Label == label {
			return m.stack[index]
		}
	}
	return nil
}

func (m *Memory) pushStackFrame(isolate bool) {
	m.stackFrames = append(m.stackFrames, StackFrame{
		stackTopPointer:    m.stackTopPointer,
		stackBottomPointer: m.stackBottomPointer,
		stackGCPointer:     m.stackGCPointer,
	})
	m.stackGCPointer = m.stackTopPointer
	if isolate {
		m.stackBottomPointer = m.stackTopPointer
	}
}

func (m *Memory) popStackFrame() {
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
		toGc[i] = nil
	}
}

type VMRuntime struct {
	mem             *Memory
	GlobalFunctions map[string]*Object
	structObjects   map[string]*Object
}

func New() *VMRuntime {
	return &VMRuntime{
		mem:             NewMemory(),
		structObjects:   map[string]*Object{},
		GlobalFunctions: map[string]*Object{},
	}
}

func (ctx *VMRuntime) AllocObject(label string) *Object {
	return ctx.mem.Alloc(label)
}

func (ctx *VMRuntime) GetObject(label string) *Object {
	if obj, ok := Functions[label]; ok {
		return &Object{Pointer: obj}
	}
	if obj, ok := ctx.GlobalFunctions[label]; ok {
		return obj
	}
	if obj, ok := ctx.structObjects[label]; ok {
		return obj
	}
	return ctx.mem.GetObject(label)
}

func (ctx *VMRuntime) PushStackFrame(isolate bool) {
	ctx.mem.pushStackFrame(isolate)
}

func (ctx *VMRuntime) PopStackFrame() {
	ctx.mem.popStackFrame()
}

func (ctx *VMRuntime) AddGlobalFunction(object *Object) {
	if _, ok := ctx.GlobalFunctions[object.Label]; ok {
		log.Panic("Object name repeated")
	}
	ctx.GlobalFunctions[object.Label] = object
}

func (ctx *VMRuntime) AddStructObject(object *Object) {
	if _, ok := ctx.structObjects[object.Label]; ok {
		log.Panic("structObject repeated", object.Label)
	}
	ctx.structObjects[object.Label] = object
}

func (ctx *VMRuntime) GetTypeObject(label string) *Object {
	return ctx.structObjects[label]
}

func (ctx *VMRuntime) IsGlobal(label string) bool {
	if _, ok := Functions[label]; ok {
		return true
	}
	return false
}

func (ctx *VMRuntime) Objects() []*Object {
	var objs []*Object
	for _, object := range ctx.structObjects {
		objs = append(objs, object)
	}
	for _, object := range ctx.GlobalFunctions {
		objs = append(objs, object)
	}
	return objs
}
