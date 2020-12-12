package qp

import (
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
	for _, testcase := range testCases {
		statement := NewParse2(testcase.data).Parse()
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
			`type user{var a = 1}`, `type user{
	var a=1;
}`,
		},
		{
			`type user{var a = func(){}}`, `type user{
	var a=func;
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

func TestNewParse2Invoke(t *testing.T) {
	testcases := []struct {
		data string
	}{
		{
			`
var a = 1
println(a,"hello")

func hello(){
	println("hello")
}
hello()

var f = func(){
var b = 3
	println(a,"lambda")
}
f()

`,
		},
	}
	for _, testcase := range testcases {
		NewParse2(testcase.data).Parse().Invoke()
	}
}
