package qp

type Bool bool

var trueObject = Bool(true)
var falseObject = Bool(false)

func (b Bool) String() string {
	if b {
		return "true"
	}
	return "false"
}

func (b Bool) Invoke() Expression {
	return b
}

func (b Bool) GetType() Type {
	return BoolObjectType
}
