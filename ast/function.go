package ast

type Function interface {
	Expression
	Call(arguments ...Expression) Expression
}
