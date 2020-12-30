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

func runScript(script string) {
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
	runScript(script)
}

func TestGenCallCode(t *testing.T) {
	script := `
var a = "HELLO"
var b = a.to_lower()
println(b,1,2,3,4)
`
	runScript(script)
}

func TestFib35(t *testing.T) {
	script := `
func fib(a){
	if a < 2 {
		return a
	}
	return fib(a-1) + fib(a-2)
}
var begin = now()
println("35",fib(35),now()-begin)
`
	runScript(script)
}

func TestFib35Cunt(t *testing.T) {
	script := `
type Count {}
var count = Count{}
func fib(a,count){
	if a < 2 {
		return a
	}
	count.i ++
	return fib(a-1,count) + fib(a-2,count)
}
var begin = now()
println("35",fib(35,count),now()-begin,count.i)
`
	runScript(script)
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

	runScript(script)
}

func TestGenTime(t *testing.T) {
	runScript(`
var a = now()
var b = now()
println(b-a)
`)
}

func TestGenFor(t *testing.T) {
	runScript(`

func fib(a){
	if a < 2 {
		return a
	}
	return fib(a-1) + fib(a-2)
}

println("num|result|time")
println("---|------|-----")

for var i = 0; i < 36; i++ {
	var s = now()
	var b = fib(i)
	var e = now()
	println(i,"|",b,"|",e-s)
}

`)
}

func TestUserObject(t *testing.T) {
	runScript(`

type User{}
func User.hello(a){
	println("hello",100, a)
	this.a = 100
}
var u = User{}
u.hello(1)
println(u.a)
var c = u
println(c.a)
`)

}

func TestUserLambda(t *testing.T) {
	runScript(`


type User{}

func User.printName(id){
	println(this.name,id)
}

var u = User{}
u.name = "jojo"
u.printName(1999)

`)

}

func TestFixFunctionCallAsArguments(t *testing.T) {
	runScript(`

func getNum(a){
	return a
}

func printlnN(a,b,c,d,e){
	println(a,b,c,d,e)
}
printlnN(getNum(1),getNum(2),getNum(3),getNum(4),getNum(5))

`)
}

func TestTestLambdaClosure(t *testing.T) {
	runScript(`

type User{}

var u = User{}

u.name = "hello"

var f = func(){
	println(u.name)
	u.id = 100
}

f()
println(u.id)

`)
}

func TestNil(t *testing.T) {
	runScript(`

var a = nil
if a == nil{
	println("a is nil ")
}

`)
}

//todo fix
func TestList(t *testing.T) {
	runScript(`

type Item {
}

type List {
}

func List.insert(left){
	println(left,"1")
    var item =Item{}
	item.value = left
    if this.head == nil {
		println(left,"2")
        this.head = item
    }else{
		println(left,"3")
        item.Next = this.head
        this.head = item
    }
}


var list = List{}

for var count = 0;count < 10;count++{
	list.insert(count)
}



`)
}