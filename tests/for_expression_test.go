package tests

import (
	"fmt"
	"testing"

	"gitlab.com/akzj/qp/parser"
	stackmachine "gitlab.com/akzj/qp/stack-machine"
)

func runScript(script string) {
	parser := parser.New(script)

	statements := parser.Parse()
	vm := parser.GetVMContext()
	objects := vm.Objects()
	for _, it := range objects {
		statements = append(statements, it)
	}
	GC := stackmachine.NewCodeGenerator()
	fmt.Println(GC.Gen(statements))

	m := stackmachine.NewMachine(GC)
	m.Run()
}

func TestForExpressionTest(t *testing.T) {
	cases := []struct {
		exp string
		val interface{}
	}{{
		exp: `
var a = 10
var a1 = a
if a < 2{
	a = 3
}

println(a)

`, val: int64(3),
	}}

	for _, Case := range cases {
		expression := parser.Parse(Case.exp)
		if expression == nil {
			t.Fatal("Parse failed")
		}
		fmt.Println("---------------------------")
		expression.Invoke()
		fmt.Println()
		fmt.Println("---------------------------")
		runScript(Case.exp)
	}
}
