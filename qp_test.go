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
		val interface{}
	}{{
		exp: `
if 2 < 3{
	return 1+1
}
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
var a = 1
var b = 1
for ;a < 10; a++{
	println(a,b)
	if a > 3{
		return 1
	}
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
	c:1,
}
// get field
println(u.id) //i

//call user function
u.print() //print(u)

// assign function abject to user object
u.hello = func(){
	println(222)
}

//alloc int field
var b = 1

//closure
u.incB= func(){
	b++
}

// call function object
println(b) //1
u.getB()
println(b) //2

`
	fmt.Println(data)
}
