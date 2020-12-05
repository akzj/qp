package qp

type memory struct {
	heapObject  []*Object
	stackObject []*Object
}

func (m *memory) alloc(label string) *Object {
	if o := m.getObject(label); o != nil {
		return o
	}
	var obj = new(Object)
	m.heapObject = append(m.heapObject, obj)
	obj.pointer = len(m.heapObject) - 1
	obj.label = label
	return obj
}

func (m *memory) getObject(label string) *Object {
	for _, object := range m.heapObject {
		if object.label == label {
			return object
		}
	}
	return nil
}

type VMContext struct {
	mem *memory
}

func newVMContext() *VMContext {
	return &VMContext{
		mem: &memory{},
	}
}

func (ctx *VMContext) allocObject(label string) *Object {
	return ctx.mem.alloc(label)
}

func (ctx *VMContext) getObject(label string) *Object {
	return ctx.mem.getObject(label)
}
