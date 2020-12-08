package qp

var breakObject = &BreakObject{}

type BreakObject struct {
}

func (b BreakObject) invoke() (Expression, error) {
	return b, nil
}

func (b BreakObject) getType() Type {
	return breakTokenType
}
