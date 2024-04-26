package main

import "fmt"

type D struct {
	a int
	b int
}
type B struct {
	a int
	b int
}

type m interface{}

func main() {

	var k m

	k = 1

	_, e := k.(D)

	fmt.Println(e)

}
