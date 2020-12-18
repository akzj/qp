package ast

import (
	"gitlab.com/akzj/qp/lexer"
	"gitlab.com/akzj/qp/runtime"
	"time"
)

type TimeObject time.Time

func (t TimeObject) String() string {
	return time.Time(t).String()
}

func (t TimeObject) Invoke() runtime.Invokable {
	return t
}

func (t TimeObject) GetType() lexer.Type {
	return lexer.TimeObjectType
}

type DurationObject time.Duration

func (d DurationObject) Invoke() runtime.Invokable {
	return d
}

func (d DurationObject) GetType() lexer.Type {
	return lexer.DurationObjectType
}

func (d DurationObject) String() string {
	return time.Duration(d).String()
}