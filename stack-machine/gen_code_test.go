package stackmachine

import (
	"fmt"
	"gitlab.com/akzj/qp/parser"
	"testing"
)

func TestGenStoreIns(t *testing.T) {
	script := `
var a = 2
if a > 1{
	print(a)
}
`

	/*
		push 2
		store a
		load a
		push 1
		cmp >
		jump R 3
		push 1
		jump R 4
		load a
		push 1
		call println
	*/
	statements := parser.New(script).Parse()
	GC := NewGenCode()
	fmt.Println(GC.Gen(statements))
}

func TestGenCallCode(t *testing.T) {
	script := `
println(1+1)
`
	statements := parser.New(script).Parse()
	GC := NewGenCode()
	fmt.Println(GC.Gen(statements))
}

func TestGenReturnVal(t *testing.T) {
	script := `
func fib(left){
	if left < 2 {
		return left
	}
	var a = fib(left-2)
	println(a)
	println(left)
	var b = fib(left-1)
	println(b)
	return  a + b
}
var a = fib(2)
println(a)
`
	parser := parser.New(script)

	statements := parser.Parse()
	vm := parser.GetVMContext()
	objects := vm.Objects()
	for _, it := range objects {
		statements = append(statements, it)
	}
	GC := NewGenCode()
	fmt.Println(GC.Gen(statements))

	m := New()
	m.instructions = GC.ins
	m.symbolTable = GC.symbolTable

	m.Run()
}

func TestGenFuncStatement(t *testing.T) {
	script := `

func hello3(a,b,c){
	println(a,b,c)
}

func hello2(a,b,c){
	println(a,b,c)
	hello3(b,c,a)
}

func hello(a,b,c){
	println(a,b,c)
	hello2(b,c,a)
}

hello(1,2,3)
hello(4,5,6)
`

	parser := parser.New(script)

	statements := parser.Parse()
	vm := parser.GetVMContext()
	objects := vm.Objects()
	for _, it := range objects {
		statements = append(statements, it)
	}
	GC := NewGenCode()
	fmt.Println(GC.Gen(statements))

	m := New()
	m.instructions = GC.ins
	m.symbolTable = GC.symbolTable

	m.Run()
}
