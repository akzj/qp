package stackmachine

import (
	"fmt"
	"gitlab.com/akzj/qp/parser"
	"testing"
)

func TestGenStoreIns(t *testing.T) {
	srcipt := `
var a = 1
if a < 1{
	
}
`
	statements := parser.New(srcipt).Parse()
	GC := NewGenCode(statements)
	fmt.Println(GC.Gen())
}
