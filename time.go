package qp

import "time"

type TimeObject time.Time

func (t TimeObject) String() string {
	return time.Time(t).String()
}

func (t TimeObject) Invoke() Expression {
	return t
}

func (t TimeObject) GetType() Type {
	return TimeObjectType
}

type DurationObject time.Duration

func (d DurationObject) Invoke() Expression {
	return d
}

func (d DurationObject) GetType() Type {
	return DurationObjectType
}

func (d DurationObject) String() string {
	return time.Duration(d).String()
}