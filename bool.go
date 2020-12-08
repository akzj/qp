package qp

type BoolObject struct {
	val bool
}

func (b *BoolObject) invoke() Expression {
	return b
}

func (b *BoolObject) getType() Type {
	return BoolObjectType
}
