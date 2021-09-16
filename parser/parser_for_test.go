package parser

import (
	"fmt"
	"testing"
)

func TestFor(t *testing.T) {
	cases := []struct {
		exp string
		val interface{}
	}{{
		exp: `

for var a = 1 ;a < 10; a++{
	for var b = 1 ;b < 10; b++{
		for var c = 1 ;c < 10; c++{
			println(a,b,c)
		}
	}
}

`, val: int64(3),
	}}

	for _, Case := range cases {
		expression := Parse(Case.exp)
		if expression == nil {
			t.Fatal("Parse failed")
		}
		fmt.Println("---------------------------")
		expression.Invoke()
	}
}
