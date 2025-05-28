package main

import "fmt"

func main() {
	var i int = 10
	if i > 5 {
		fmt.Println("i > 5")
	} else if i < 5 {
		fmt.Println("i < 5")
	} else {
		fmt.Println("i == 5")
	}
	switch i {
	case 1:
		fmt.Println("i == 1")
		break
	case 2:
		fmt.Println("i == 2")
		break
	case 3:
		fmt.Println("i == 3")
		break
	case 4:
		fmt.Println("i == 4")
		break
	case 5:
		fmt.Println("i == 5")
		fallthrough
	default:
		fmt.Println("i > 5")
	}

	var str string = "h"
	switch {
	case str == "h":
		fmt.Println("str == h")
		break
	case str == "e":
		fmt.Println("str == e")
		break
	case str == "l":
		fmt.Println("str == l")
		break
	case str == "o":
		fmt.Println("str == o")
		break
	default:
		fmt.Println("str != h,e,l,o")
	}
}
