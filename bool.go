package qp

type BoolObject bool

var trueObject = BoolObject(true)
var falseObject = BoolObject(false)

func (b *BoolObject) String() string {
	if *b {
		return "true"
	}
	return "false"
}

func (b *BoolObject) Invoke() Expression {
	return b
}

func (b *BoolObject) getType() Type {
	return BoolObjectType
}
