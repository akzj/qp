package qp

type BaseObject interface {
	Expression
	getObject(label string) *Object
	allocObject(label string) *Object
	addObject(k string, v *Object)
	clone() BaseObject
}

type TypeObject struct {
	vm    *VMContext
	label string
	//init statement when create objects
	init          bool
	initStatement Statements
	//user define function
	objects map[string]*Object
}

func (sObj *TypeObject) invoke() Expression {
	if sObj.init {
		return sObj
	}
	for _, statement := range sObj.initStatement {
		statement.invoke()
	}
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
		object = &Object{label: label}
		sObj.objects[label] = object
	}
	return object
}

func (sObj *TypeObject) clone() BaseObject {
	clone := *sObj
	for k, v := range sObj.objects {
		clone.addObject(k, v)
	}
	if len(sObj.initStatement) != 0 {
		clone.initStatement = make(Statements, len(sObj.initStatement))
		copy(clone.initStatement, sObj.initStatement)
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
