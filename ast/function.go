package ast

import "gitlab.com/akzj/qp/runtime"

type Function interface {
	runtime.Invokable
	Call(arguments ...runtime.Invokable) runtime.Invokable
}
