# qp

qp is small script program language

```
func fib(a){
	if a < 2 {
		return a
	}
	return fib(a-1) + fib(a-2)
}
var begin = now()
var a = fib(35)
println("35",a,now()-begin)
```
output
```
35 9227465 3.87979642s
```
