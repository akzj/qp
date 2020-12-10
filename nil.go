package qp

type NilObject struct {
}

var nilObject = NilObject{}

func (NilObject) String() string {
	return "nil"
}

func (n NilObject) Invoke() Expression {
	return n
}

func (NilObject) getType() Type {
	return nilType
}
