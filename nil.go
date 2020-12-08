package qp

type NilObject struct {
}

func (n NilObject) String() string {
	return "nil"
}

var nilObject = &NilObject{}

func (n *NilObject) invoke() Expression{
	return n
}

func (NilObject) getType() Type {
	return nilTokenType
}
