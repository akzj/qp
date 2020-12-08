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
		if val, err := expression.invoke(); err != nil {
			t.Errorf(err.Error())
		} else {
			if val.(*BoolObject).val != Case.expect {
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
	if val, err := expression.invoke(); err != nil {
		t.Errorf(err.Error())
	} else {
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
var main = func() {
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
println(main())
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
		if _, err := expression.invoke(); err != nil {
			t.Errorf(err.Error())
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
		if val, err := statements.invoke(); err != nil {
			t.Errorf(err.Error())
		} else {
			if val.(*ReturnStatement).returnVal.(*IntObject).val != Case.val {
				t.Fatalf("no match %+v %+v", val.(*IntObject).val, Case.val)
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
		if val, err := expression.invoke(); err != nil {
			t.Errorf(err.Error())
		} else {
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
		if val, err := statements.invoke(); err != nil {
			t.Errorf(err.Error())
		} else {
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
		if _, err := expression.invoke(); err != nil {
			t.Errorf(err.Error())
		}
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
println(a++)
println(a+1)
`, val: int64(3),
	}}
	for _, Case := range cases {
		expression := Parse(Case.exp)
		if expression == nil {
			t.Fatal("Parse failed")
		}
		fmt.Println("-------------")
		if _, err := expression.invoke(); err != nil {
			t.Errorf(err.Error())
		}
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
		if val, err := expression.invoke(); err != nil {
			t.Errorf(err.Error())
		} else {
			fmt.Println("result", val)
		}
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
	if _, err := expression.invoke(); err != nil {
		t.Errorf(err.Error())
	}
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
	if _, err := expression.invoke(); err != nil {
		t.Errorf(err.Error())
	}
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
	if _, err := expression.invoke(); err != nil {
		t.Errorf(err.Error())
	}
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
	if _, err := statements.invoke(); err != nil {
		t.Errorf(err.Error())
	}
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
u.print() //print(u)

// assign function abject to user objects
u.hello = func(){
	println(222)
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

`
	statements := Parse(data)
	if statements == nil {
		t.Fatal("Parse failed")
	}
	fmt.Println("-------------------------------------------------------------")
	if _, err := statements.invoke(); err != nil {
		t.Errorf(err.Error())
	}
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
			if _, err := statements.invoke(); err != nil {
				t.Fatal(err)
			}
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
	if _, err := statements.invoke(); err != nil {
		t.Fatal("test failed", err)
	}
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
	if _, err := statements.invoke(); err != nil {
		t.Fatal("test failed", err)
	}
}

func TestNil(t *testing.T) {
	data := `
var a
if a == nil{
	println("a is nil",1000)
}else{
	println(a)
}
`
	Parse(data).invoke()
}
