package stackmachine

import (
	"fmt"
	"gitlab.com/akzj/qp/parser"
	"log"
	"testing"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

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
var a = "HELLO"
a.to_lower()
var c = 1
println(a,c)
var d = 2
println(d)
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

func TestGenReturnVal(t *testing.T) {
	script := `
func fib(a){
	if a < 2 {
		return a
	}
	return fib(a-1) + fib(a-2)
}
var a = fib(29)
println("35",a)
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

func hello4(a,b,c){
	println(a,b,c)
}

func hello3(a,b,c){
	println(a,b,c)
	hello4(b,c,a)
}

func hello2(a,b,c){
	println(a,b,c)
	hello3(b,c,a)
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
