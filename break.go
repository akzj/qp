package qp

var breakObject = &BreakObject{}

type BreakObject struct {
}

func (b *BreakObject) Invoke() Expression {
	return b
}

func (b *BreakObject) getType() Type {
	return breakTokenType
}
