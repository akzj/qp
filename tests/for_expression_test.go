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
		if 1 < 2 {
			a := 1
			b := 1
			c := 1
			d := 1
			e := 1
			e2 := 1
			e3 := 1
			e4 := 1
			e5 := 1
		}

		println(e2)
		
`, val: int64(3),
	}}

	for _, Case := range cases {
		expression := parser.Parse(Case.exp)
		if expression == nil {
			t.Fatal("Parse failed")
		}
		fmt.Println("---------------------------")
		//	expression.Invoke()
		fmt.Println()
		fmt.Println("---------------------------")
		runScript(Case.exp)
	}
}
