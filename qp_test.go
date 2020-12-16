package qp

import (
	"bytes"
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	fmt.Println("hello qp")
}

func TestBuffer(t *testing.T) {
	var reader = bytes.NewReader([]byte("1+1 if else"))
	for {
		c, err := reader.ReadByte()
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		fmt.Println(string(c))
		fmt.Println(isLetter(c))
		fmt.Println(c, 'a')
		fmt.Println(c, 'i')
	}
}

func TestLexer(t *testing.T) {
	lexer := newLexer(bytes.NewReader([]byte(`
if 2 > 1{
	3+3 
return 1
var a = 
a ++
for
a//hello world
"hello"
` + "` multi-line hello\n\nworld`" + `
var a = nil
==
`)))
	if lexer == nil {
		t.Fatal("Lexer nil")
	}
	var count = 1000
	for lexer.Finish() == false && count > 0 {
		fmt.Println(lexer.Peek().String())
		lexer.next()
		count--
	}
}

func TestParseBoolExpression(t *testing.T) {
	data := `
if (1 > (1+2)) == false{
	for {
		if 0 == 0{
			 var a = func(){return 0}() == 0
		}
		println("hello")
		break
	}
}
`
	Parse(data).Invoke()
}

func TestNumAddParse(t *testing.T) {
	parser := NewParse2("1*(5+5+5)*2")
	parser.initTokens()
	expression := parser.parseFactor(0)
	if expression == nil {
		t.Fatal("Parse failed")
	}
	if val := expression.Invoke(); val != nil {
		fmt.Println(val)
	}
}

func TestReturnStatement(t *testing.T) {
	cases := []struct {
		exp string
		val int64
	}{{
		exp: `
var global = 100
var main = func(left) {
	println(left,1)
	var f = func(){
		println(left,2)
		var b = func(){
			var d = func(){
					println(left,4)
					return global
				}
			return d()
		}
		return b
	}
	var d = f
	return d()
}
var b = main("closure ")()
println(b)
`,
		val: int64(100),
	},
	}

	for _, Case := range cases {
		expression := Parse(Case.exp)
		if expression == nil {
			t.Fatal("Parse failed")
		}
		fmt.Println("------------------------------------------------")
		if val := expression.Invoke(); val != nil {
			fmt.Println(val)
		}
	}
}

func TestIfStatement(t *testing.T) {
	cases := []struct {
		exp string
		val int64
	}{{
		exp: `
if 2 > 1{
	return 3+3
}
`, val: int64(6),
	}}

	for _, Case := range cases {
		statements := Parse(Case.exp)
		if statements == nil {
			t.Fatal("Parse failed")
		}
		if val := statements.Invoke(); val != nil {
			if int64(val.(ReturnStatement).returnVal.(Int)) != Case.val {
				t.Fatalf("no match %+v %+v", int64(val.(Int)), Case.val)
			}
		}
	}
}

func TestVarAssign(t *testing.T) {
	cases := []struct {
		exp string
		val interface{}
	}{{
		exp: `
var a= 1
println(a)
return a+1
`, val: int64(3),
	}}

	for _, Case := range cases {
		expression := Parse(Case.exp)
		if expression == nil {
			t.Fatal("Parse failed")
		}
		if val := expression.Invoke(); val != nil {
			fmt.Println("result", val.String())
		}
	}
}

func TestFunctionCall(t *testing.T) {
	cases := []struct {
		exp string
		val interface{}
	}{{
		exp: `
var a = 8
if a > 10{
	println(a,1)
}else if a > 9{
	println(a,2)
}else if a >8{
	println(a,3)
}else{
	println(a,0)
}

`, val: int64(3),
	}}

	for _, Case := range cases {
		statements := Parse(Case.exp)
		if statements == nil {
			t.Fatal("Parse failed")
		}
		fmt.Println("-----------------------")
		if val := statements.Invoke(); val != nil {
			fmt.Println("result", val)
		}
	}
}

func TestInc(t *testing.T) {
	cases := []struct {
		exp string
		val interface{}
	}{{
		exp: `
var a = 1
println(a)
a++
println(a)

var b = func(){
	println(a+1)
}
b()
`, val: int64(3),
	}}

	for _, Case := range cases {
		expression := Parse(Case.exp)
		if expression == nil {
			t.Fatal("Parse failed")
		}
		fmt.Println("---------------------------")
		expression.Invoke()
	}
}

func TestAssignStatement(t *testing.T) {
	cases := []struct {
		exp string
		val interface{}
	}{{
		exp: `
var a = 1
a = 100
a++
println(a)
println(a+1)
`, val: int64(3),
	}}
	for _, Case := range cases {
		expression := Parse(Case.exp)
		if expression == nil {
			t.Fatal("Parse failed")
		}
		fmt.Println("-------------")
		expression.Invoke()
	}
}

func TestFor(t *testing.T) {
	cases := []struct {
		exp string
		val interface{}
	}{{
		exp: `

var b = 1
for var a = 1 ;a < 10; a++{
	println(a,b)
	if a > 3{
		return 1
	}
	//local var 
	var a = 100
	println(a)
	a = 101
	println(a)
}

`, val: int64(3),
	}}

	for _, Case := range cases {
		expression := Parse(Case.exp)
		if expression == nil {
			t.Fatal("Parse failed")
		}
		fmt.Println("---------------------------")
		expression.Invoke()
	}
}

func TestFunction(t *testing.T) {
	data := `
func add(a,b){
	var c = a+b
	println(c) //3
}

var a = 1
var b = 2
var c = 100

println(c) //100
add(a,b)
println(c) //100

`
	expression := Parse(data)
	if expression == nil {
		t.Fatal("Parse failed")
	}
	fmt.Println("---------------------------")
	expression.Invoke()
}

func TestStructObjectDefaultInit(t *testing.T) {
	data := `
type User {
	//define user member with default Int 1
	id:  66666
	//define user member with default nil
}

var user = User{}
user.id = 1
println(user) //

`

	expression := Parse(data)
	if expression == nil {
		t.Fatal("Parse failed")
	}
	fmt.Println("---------------------------")
	expression.Invoke()
}

func TestStructObject(t *testing.T) {
	data := `
type User {
	id: 66666
}

var user = User{
	a:1+1
}
//println struct member field
println(user.a) //2
user.a = 100
println(user.a) //100


func User.print(){
	println(this.id)
}

// -----------------------------------------------------------
user.print()

`

	expression := Parse(data)
	if expression == nil {
		t.Fatal("Parse failed")
	}
	fmt.Println("---------------------------")
	expression.Invoke()
}

func TestLambdaFunction(t *testing.T) {
	data := `
//assign function objects to var

type User{
}

var user = User{}

user.a = func(){
	println(1)
}

user.a()

//call lambda function objects
`
	statements := Parse(data)
	if statements == nil {
		t.Fatal("Parse failed")
	}
	fmt.Println("---------------------------")
	statements.Invoke()
}

func TestObject(t *testing.T) {
	data := `
type user {
	//define user member with default Int 1
	id: 1
}

//define function for user
func user.print(){
	println(this.id) //i
}

// alloc field u

var ccc = 1
var u = user{
	c:ccc
}
// Get field
println(u.id) // 1 

//call user function
u.id = 100
u.print() //print(u)

// assign function abject to user objects
u.hello = func(){
	println(u.id+1)
}

//alloc int field
var b = 1


//closure
u.incB = func(b){
	println(b,"bbbbbbbbbbbbbbbbb")
	b++
}

// call function objects
println(b) //1
u.incB(b)
println(b) //2
u.hello()

`
	statements := Parse(data)
	if statements == nil {
		t.Fatal("Parse failed")
	}
	fmt.Println("-------------------------------------------------------------")
	statements.Invoke()
}

func TestStackFrame(t *testing.T) {
	cases := []struct {
		data string
		err  bool
	}{
		{
			data: `
var a = 1
func testA(){
	var a = 100
	println(a)
}
println(a)
testA()
println(a)
`,
			err: false,
		},
		{
			data: `
var a = 1
func testA(){
	var a = 100
	println(a)
	if a > 10{
		var b = 100
		println(b)
	}
}
println(a)
testA()
println(a)
`,
			err: true,
		},
	}

	for index, Case := range cases {
		if index == 1 {

			statements := Parse(Case.data)
			if statements == nil {
				t.Fatal("Parse failed")
			}
			fmt.Println("-------------------------------------------------------------")
			statements.Invoke()
		}
	}
}

func TestClosure(t *testing.T) {
	data := `

var a =1

var f = func(){
	println(a)
}

a++

f()

`
	statements := Parse(data)
	if statements == nil {
		t.Fatal("Parse failed")
	}
	fmt.Println("-------------------------------------------------------------")
	statements.Invoke()
}

func TestString(t *testing.T) {
	data := `
var a ="Hello World"
var b = a.clone() // 
a.to_lower()
println(a,b)
`
	statements := Parse(data)
	if statements == nil {
		t.Fatal("Parse failed")
	}
	fmt.Println("-------------------------------------------------------------")
	statements.Invoke()
}

func TestNil(t *testing.T) {
	data := `
var a
if a == nil{
	println("is nil",1000,a)
}else{
	println(a)
}
`
	Parse(data).Invoke()
}

func TestArray(t *testing.T) {
	data := `

var a = [1,func(){
	println("function object")
}]
a.append(2)
println(a.Get(0))
var a = a.Get(1)
a()

`
	Parse(data).Invoke()
}

func TestFunctionCallC(t *testing.T) {
	data := `

func(){
	return func(){
	return func(){
	return func(){	
}
}
}
}()()()()
`
	Parse(data).Invoke()
}

func TestBuiltInFunctionNow(t *testing.T) {
	data := `
var res = now()
println(now()-res)
`
	Parse(data).Invoke()
}

func TestList(t *testing.T) {

	data := `
type Item {
}

type List {
}

func List.insert(left){
    var item =Item{}
	item.value = left
    if this.head == nil {
        this.head = item
    }else{
        item.next = this.head
        this.head = item
    }
}


var list = List{}

for var count = 0;count < 10;count++{
	list.insert(count)
}


for var head =list.head ;head != nil; head = head.next{
	println(head.value)
}

`
	Parse(data).Invoke()
}

/*
832040
--- PASS: TestFib (0.01s)
*/
func fib(val int) int {
	if val < 2 {
		return val
	}
	var a = fib(val - 2)
	var b = fib(val - 1)
	return a + b
}

func TestFib(t *testing.T) {
	fmt.Println(fib(25))
}

//todo fix
/*
5702887
--- PASS: TestFibonacci (12.52s)
PASS
*/
func TestFibonacci(t *testing.T) {
	data := `

func fib(left){
	if left < 2 {
		return left
	}
var l = fib(left-2) 
var r = fib(left-1)
	 var res = l+r
	return res
}

println("num|result|take time")
println("---|------|---------")
for var num = 20; num < 30; num++ {
	var begin = now()
	println(num,"|",fib(num),"|",now() - begin)
}
println("")
`

	if statement := NewParse2(data).Parse(); statement == nil {
		panic("parse failed")
	} else {
		statement.Invoke()
	}

}

func TestPrintlnArray(t *testing.T) {
	data := `

var arr = []
for var i = 0; i < 100;i ++{
	arr.append(i)
}

println(arr.size())
println(arr.Get(10))
`
	Parse(data).Invoke()
}
