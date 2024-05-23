package main

import "fmt"

//在 Go 语言中，数组是一种值类型，而且不同长度的数组属于不同的类型。例如 [2]int 和 [20]int 属于不同的类型。
//当值类型作为参数传递时，参数是该值的一个拷贝，因此更改拷贝的值并不会影响原值。

func main() {
	a := [2]int{1, 2}
	foo(a)
	fmt.Println(a)
}

func foo(a [2]int) {
	a[0] = 200
}
