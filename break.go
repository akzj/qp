package qp

var breakObject = &BreakObject{}

type BreakObject struct {
}

func (b *BreakObject) String() string {
	return "break"
}

func (b *BreakObject) Invoke() Expression {
	return b
}

func (b *BreakObject) GetType() Type {
	return BreakType
}
