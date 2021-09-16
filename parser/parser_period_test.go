package parser

import (
	"fmt"
	"testing"
)

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
user.insert(6)

var value = user.getFirst()
if value == nil{
	println("value is nil")
} else if value == 5{
	println("first value is 5")
} else{
	println("first value is",value)	
}

for var first=user.head;first != nil;first = first.next{
	println(first.value)
}
`

	fmt.Println(New(data).Parse().String())
	New(data).Parse().Invoke()
}
