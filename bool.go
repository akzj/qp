package qp

type BoolObject bool

var trueExpression = BoolObject(true)

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
