package runtime

import (
	"fmt"
	"gitlab.com/akzj/qp/lexer"
)

type Invokable interface {
	Invoke() Invokable
	GetType() lexer.Type
	fmt.Stringer
}