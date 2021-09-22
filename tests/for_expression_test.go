package tests

import (
	"fmt"
	"testing"

	"gitlab.com/akzj/qp/parser"
	stackmachine "gitlab.com/akzj/qp/stack-machine"
)

func runScript(script string, printIns bool) {
	parser := parser.New(script)

	statements := parser.Parse()
	vm := parser.GetVMContext()
	objects := vm.Objects()
	for _, it := range objects {
		statements = append(statements, it)
	}
	GC := stackmachine.NewCodeGenerator()
	fmt.Println(GC.Gen(statements))

	m := stackmachine.NewMachine(GC, stackmachine.Options{Debug: printIns})
	m.Run()
}

func TestForExpressionTest(t *testing.T) {
	cases := []struct {
		exp string
		val interface{}
	}{{
		exp: `
		
		for k := 0; k < 10;k++{
			for i := 1;i < 4;i++{
				println(k,i)
			}
			println("---------------")
		}
`, val: int64(3),
	}}

	for _, Case := range cases {
		expression := parser.Parse(Case.exp)
		if expression == nil {
			t.Fatal("Parse failed")
		}
		fmt.Println("---------------------------")
		fmt.Println(expression.String())
		fmt.Println()
		fmt.Println("---------------------------")
		runScript(Case.exp, false)
	}
}
