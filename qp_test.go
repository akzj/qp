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
`)))
	if lexer == nil {
		t.Fatal("lexer nil")
	}
	var count = 1000
	for lexer.finish() == false && count > 0 {
		fmt.Println(lexer.peek())
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
		expression := Parse(Case.expStr)
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
	expression := Parse("1*(5+5+5)*2")
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
	return 3+3
if 2 > 3{
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
		expression := Parse(Case.exp)
		if expression == nil {
			t.Fatal("Parse failed")
		}
		if val, err := expression.invoke(); err != nil {
			t.Errorf(err.Error())
		} else {
			if val.(*IntObject).val != Case.val {
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
var b=a+1
return b+1
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
var a = 1
if a > 1{
	println(a,1)
}else if a > 2{
	println(a,2)
}else if a > 3{
	println(a,3)
}else{
	println(a,0)
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
for ;; {
	println(a)
	a++
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
