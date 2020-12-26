package stackmachine

import (
	"fmt"
	"gitlab.com/akzj/qp/parser"
	"testing"
)

func TestGenStoreIns(t *testing.T) {
	script := `
var a = 1
func test(b,c){
	var a = 1000 + b
	println(b,c)
	println(a)
}
test(2,3)
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
func fib(a){
	if a < 2 {
		return a
	}
	//println(a)
	var b = fib(a-1)
	var c = fib(a-2)
	return  b+c
}
var a = fib(35)
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


func hello2(a,b,c){
	println(a,b,c)
}

func hello(a,b,c){
	println(a,b,c)
	hello2(b,c,a)
}

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
