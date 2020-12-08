package qp

type BoolObject struct {
	val bool
}

func (b *BoolObject) invoke() (Expression, error) {
	return b, nil
}

func (b *BoolObject) getType() Type {
	return BoolObjectType
}
