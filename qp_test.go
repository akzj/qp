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
			if val.(bool) != Case.expect {
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

func TestIfStatement(t *testing.T) {
	cases := []struct {
		exp string
		val interface{}
	}{{
		exp: `
if 2 > 1{
	return 3+3
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
