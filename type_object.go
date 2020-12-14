package qp

type BaseObject interface {
	Expression
	getObject(label string) *Object
	allocObject(label string) *Object
	clone() BaseObject
}

type TypeObjectPropTemplate struct {
	name string
	exp  Expression
}

func (t TypeObjectPropTemplate) String() string {
	return t.name + ":" + t.exp.String()
}

type TypeObject struct {
	vm    *VMContext
	label string
	//init statement when create objects
	init                    bool
	typeObjectPropTemplates []TypeObjectPropTemplate
	//user define function
	objects map[string]*Object
}

func (sObj *TypeObject) String() string {
	str := "type " + sObj.label + "{"
	for index, exp := range sObj.typeObjectPropTemplates {
		if index == 0 {
			str += "\n"
		}
		str += "\t" + exp.String() + ";\n"

	}
	return str + "}"
}

func (sObj *TypeObject) Invoke() Expression {
	return sObj
}

func (sObj *TypeObject) getType() Type {
	return TypeObjectType
}

func (sObj *TypeObject) getObject(label string) *Object {
	object, ok := sObj.objects[label]
	if ok {
		return object
	}
	return nil
}

func (sObj *TypeObject) allocObject(label string) *Object {
	if sObj.objects == nil {
		sObj.objects = map[string]*Object{}
	}
	object, ok := sObj.objects[label]
	if ok {
		return object
	} else {
		object = &Object{
			inner: nilObject,
			label: label,
		}
		sObj.objects[label] = object
	}
	return object
}

func (sObj *TypeObject) clone() BaseObject {
	clone := *sObj
	for k, v := range sObj.objects {
		clone.addObject(k, v)
	}
	if len(sObj.typeObjectPropTemplates) != 0 {
		clone.typeObjectPropTemplates = make([]TypeObjectPropTemplate, len(sObj.typeObjectPropTemplates))
		copy(clone.typeObjectPropTemplates, sObj.typeObjectPropTemplates)
	}
	return &clone
}

func (sObj *TypeObject) addObject(k string, v *Object) {
	if v == nil {
		panic(v)
	}
	if sObj.objects == nil {
		sObj.objects = map[string]*Object{}
	}
	sObj.objects[k] = v
}
