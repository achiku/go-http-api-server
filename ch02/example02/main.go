package main

import "fmt"

type myFunc func() string

func (f myFunc) WrapFunc() {
	fmt.Print("before msg\n")
	fmt.Printf("%s\n", f())
	fmt.Print("after msg\n")
}

func main() {
	f := myFunc(func() string {
		return "hello!!"
	})
	f.WrapFunc()
}
