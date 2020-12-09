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
		t.Fatal("lexer nil")
	}
	var count = 1000
	for lexer.finish() == false && count > 0 {
		fmt.Println(lexer.peek().String())
		lexer.next()
		count--
	}
}

func TestLessExpression(t *testing.T) {

	cases := []struct {
		expStr string
		expect bool
	}{
		{
			"1 < 2",
			true,
		},
		{
			"1 <= 2",
			true,
		},
		{
			"2 <= 1",
			false,
		},
		{
			"2 < 1",
			false,
		},
	}

	for _, Case := range cases {
		parser := newParser(bytes.NewReader([]byte(Case.expStr)))
		expression := parser.parseExpression()
		if expression == nil {
			t.Fatal("Parse failed")
		}
		if val := expression.Invoke(); val != nil {
			if bool(*(val.(*BoolObject))) != Case.expect {
				t.Fatalf("expression parse failed,`%s` "+
					"result `%+v` expect `%+v`", Case.expStr, val, Case.expect)
			}
		}
	}
}

func TestNumAddParse(t *testing.T) {
	parser := newParser(bytes.NewReader([]byte("1*(5+5+5)*2")))
	expression := parser.parseExpression()
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
var global = 1
var main = func(val) {
	println(val)
	var out = 10 
	var f = func(){
		var a = 100
		var b = func(){
			var c = 1000
			var d = func(){
					return a + out + global +c
				}
			return d()
		}
		return b
	}
	var b = f()
	return b()
}
println(main("hello"))
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
			if int64(*val.(*ReturnStatement).returnVal.(*IntObject)) != Case.val {
				t.Fatalf("no match %+v %+v", int64(*val.(*IntObject)), Case.val)
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
			fmt.Println("result", val)
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
	//define user member with default IntObject 1
	var id = 66666
	//define user member with default nil
	var id2
}

var user = User{}
println(user.id) //

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
	//define user member with default IntObject 1
	var id = 66666
	//define user member with default nil
	var id2
}

var user = User{
a:1+1
}
//println struct member field
println(user.a) //2
user.a = 100
println(user.a) //100


func User.print(){
	println(.id)
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
	//define user member with default IntObject 1
	var id = 1
	//define user member with default nil
	var id2
}

//define function for user
func user.print(){
	println(.id) //i
}

// alloc field u
var u = user{
//init field
	c:1
}
// get field
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
println(a.get(0))
var a = a.get(1)
a()

`
	Parse(data).Invoke()
}

func TestFunctionCallC(t *testing.T) {
	data := `
var a0 = 1
var a1 = 2
var a2 = 3
var a3 = 4
var a4 = 5

println(a0)
func(){
	println(a1)
	return func(){
			println(a2)
			return func(){
				println(a3)
				return func(){
					println(a4)
			}
		}
	}
}()()()()
`
	Parse(data).Invoke()
}
