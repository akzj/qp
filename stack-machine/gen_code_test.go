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
	GC := NewGenCode(statements)
	fmt.Println(GC.Gen())
}

func TestGenCallCode(t *testing.T) {
	script := `
println(1+1)
`
	statements := parser.New(script).Parse()
	GC := NewGenCode(statements)
	fmt.Println(GC.Gen())
}
