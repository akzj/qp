package qp

import (
	"fmt"
	"log"
	"strings"
	"testing"
)

func TestParser2_Parse(t *testing.T) {
	testCases := []struct {
		data   string
		expect string
	}{
		{
			`
if 1==2 {
}
`, `if (1==2){}`,
		},
		{
			`
if (1 == 2)==false{
	
}
`, `if (((1==2))==false){}`,
		},
		{
			`if (1==2) || (2==4) || func(){}() {}`,
			`if ((((1==2))||((2==4)))||func()){}`,
		},
		{
			`
if 1 == 2 && 2==4 || 3== 4 && true{
}
`, `if (((1==2)&&(2==4))||((3==4)&&true)){}`,
		},
		{
			`if a().a ==false || 1==2 && a == b{}`, `if ((a().a==false)||((1==2)&&(a==b))){}`,
		},
		{
			`if (a()).a()() ==false || 1==2 && a == b{}`, `if (((a()).a()()==false)||((1==2)&&(a==b))){}`,
		},
		{
			`var a = 1`, `var a=1`,
		},
	}
	for index, testcase := range testCases {
		if index != len(testCases)-1{
			continue
		}
		log.Println(strings.Repeat("-",100))
		log.Println(testcase.data)
		p := NewParse2(testcase.data)
		p.initTokens()
		for _,token := range p.tokens{
			log.Println(token)
		}
		statement := p.Parse()
		if statement == nil {
			t.Fatal("parse failed")
		}
		if str := statement.String(); str != testcase.expect {
			t.Logf("parse %s failed,result \n%s\n expect\n%s", testcase.data, str, testcase.expect)

		}

	}

	testTypeObjects := []struct {
		data   string
		expect string
	}{
		{
			`type user{}`, `type user{}`,
		},
		{
			`type user{a:1}`, `type user{
	a:1;
}`,
		},
		{
			`type user{a:func(){}}`, `type user{
	a:func;
}`,
		},
	}

	for _, testcase := range testTypeObjects {
		p := NewParse2(testcase.data)
		p.initTokens()
		p.expectType(p.nextToken(), typeType)
		statement := p.parseTypeStatement()
		if str := statement.String(); str != testcase.expect {
			t.Logf("parse %s failed,result \n%s\n expect\n%s", testcase.data, str, testcase.expect)
		}

	}
}

func TestIfElseIF(t *testing.T) {
	for _, token := range NewParse2(`if 1==1{}else if 1==2{}else{}`).initTokens().tokens {
		fmt.Println(token)
	}
}

func TestNewParse2Invoke(t *testing.T) {
	testcases := []struct {
		data string
	}{
		/*{
		`
type User{
	name:"hello"
}
var u = User{
	id :1
	print:func(){
		println("hello world")
	}
}
println(u.name,u.id)
u.print()
`,
		},*/
		{
			`
var f  = func(){
	var c  =func(){
		var a = func(){
			println("func c call")
		}

		if a == nil{
			println("a nil")
		}
		return a
	}
	return c()
}
f()
`,
		},
	}
	for _, testcase := range testcases {
		statements := NewParse2(testcase.data).Parse()
		fmt.Println("--------")
		statements.Invoke()
	}
}
