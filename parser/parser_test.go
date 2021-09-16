package parser

import (
	"fmt"
	"testing"

	"gitlab.com/akzj/qp/lexer"
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
		if index != len(testCases)-1 {
			continue
		}
		p := New(testcase.data)
		p.initTokens()
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
		p := New(testcase.data)
		p.initTokens()
		p.expectType(p.nextToken(), lexer.TypeType)
		statement := p.parseTypeStatement()
		if str := statement.String(); str != testcase.expect {
			t.Logf("parse %s failed,result \n%s\n expect\n%s", testcase.data, str, testcase.expect)
		}

	}
}

func TestIfElseIF(t *testing.T) {
	for _, token := range New(`if 1==1{}else if 1==2{}else{}`).initTokens().tokens {
		fmt.Println(token)
	}
}

func TestPeriod(t *testing.T) {
	data := `

type Item {
}

type User{
}

func User.insert(val){
	var item = Item{value:val}
	if this.head==nil{
		this.head = item
	}else{
		item.next = this.head
		this.head = item
	}
	println(this.head.value)
}


func User.getFirst(){
	return this.head.value
}

var user = User{}

user.insert(1)
user.insert(2)
user.insert(3)
user.insert(4)
user.insert(5)

var value = user.getFirst()
if value == nil{
	println("value is nil")
} else if value == 1{
	println("value is 1")
} else{
	println(value)	
}
println(value)

for var first=user.head;first != nil;first = first.next{
	println(first.value)
}

`

	fmt.Println(New(data).Parse().String())
	New(data).Parse().Invoke()
}

func TestTypeStruct(t *testing.T) {
	data := `

type Info {}
type User{}

func User.hello(){
	var info = Info{}
	info.id = 1
	this.info = info
}

var user =User{}
user.hello()
println(user.info.id)

`
	New(data).Parse().Invoke()
}

func TestNewParse2Invoke(t *testing.T) {
	testcases := []struct {
		data string
	}{
		{
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
		},
		{
			`
func(){
	var c  =func(){
		var a = func(){
			println("func c Call")
		}
		if a == nil{
			println("a nil????")
		}
		a()
		return a
	}
	return c()
}()
`,
		},
	}
	for _, testcase := range testcases {
		statements := New(testcase.data).Parse()
		fmt.Println("--------")
		statements.Invoke()
	}
}
