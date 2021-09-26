package tests

import (
	"gitlab.com/akzj/qp/parser"
	stackmachine "gitlab.com/akzj/qp/stack-machine"
)

func run(script string, printIns bool) {
	parser := parser.New(script)
	statements := parser.Parse()
	GC := stackmachine.NewCodeGenerator()
	GC.Gen(statements)
	m := stackmachine.NewMachine(GC, stackmachine.Options{Debug: printIns})
	m.Run()
}

func ExampleForExpression_i_j() {
	run(`
	for i := 0;i < 3;i++{
		for j := 0;j < 3;j++{
			println(i,j)
		}
	}
	`, false)
	//Output:
	//0 0
	//0 1
	//0 2
	//1 0
	//1 1
	//1 2
	//2 0
	//2 1
	//2 2
}

func ExampleForExpression_if() {
	run(`
	for i := 0;i < 3;i++{
		for j := 0;j < 3;j++{
			if i + j  == 2{
				println(i,j)
			}
		}
	}
	`, false)
	//Output:
	//0 2
	//1 1
	//2 0
}

func ExampleForExpression_3() {
	run(`
	a := 0
	for i := 0;i < 3;i++{
		for j := 0;j < 3;j++{
			for k := 0;k < 3;k++{
				a++			
			}	
		}
	}
	println(a)
	`, false)
	//Output:
	//27
}

func ExampleIfExpression() {
	run(`
	if 0 < 1{
		println("hello")
	}
	`, false)
	//Output:
	//hello
}

func ExampleIfVarCheck() {
	run(`
	a := 1
	if (a +1) * 0 <= 2{
		println(a+1)
	}
	`, false)
	//Output:
	//2
}
