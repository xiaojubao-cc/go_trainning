package main

import "fmt"

type Foo map[string]string
type Bar struct {
	thingOne string
	thingTwo int
}

func main() {
	// OK:
	y := new(Bar)
	(*y).thingOne = "hello"
	(*y).thingTwo = 1
	// not OK:
	// OK:
	x := make(Foo)
	x["x"] = "goodbye"
	x["y"] = "world"
	fmt.Printf("x: %T\n", *y)
}
