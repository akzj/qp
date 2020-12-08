package qp

import (
	"fmt"
	"log"
	"testing"
)

func init() {
	log.SetFlags(log.Lshortfile)
}
func TestVMMemoryStackTest(t *testing.T) {
	vm := newVMContext()

	t.Run("allocObject", func(t *testing.T) {
		objectA := vm.allocObject("a")
		objectA1 := vm.getObject("a")
		if objectA != objectA1 {
			t.Fatal("allocObject error")
		}
	})

	t.Run("stackFrame", func(t *testing.T) {
		objectB := vm.allocObject("b")

		//
		vm.pushStackFrame(false)
		objectB1 := vm.allocObject("b")
		if object := vm.getObject("b"); object != objectB1 {
			if object == nil {
				t.Fatal("no find objects")
			}
			fmt.Println(object)
			fmt.Println(objectB1)
			t.Fatalf("pushStackFrame failed %s %s", object.String(), objectB1.String())
		}

		//
		vm.popStackFrame()
		if vm.getObject("b") == objectB1 {
			t.Fatal("popStackFrame failed")
		}
		if vm.getObject("b") != objectB {
			t.Fatal("popStackFrame failed")
		}
	})

	t.Run("stackFrameIsolate", func(t *testing.T) {
		a := vm.allocObject("a")

		//isolate
		vm.pushStackFrame(true)

		if vm.getTypeObject("a") != nil {
			t.Fatal("isolate failed")
		}

		vm.popStackFrame()
		if vm.getObject("a") != a {
			t.Fatal("isolate failed")
		}
	})
}
