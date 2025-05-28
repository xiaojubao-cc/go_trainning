package main

import "errors"

func main() {
	_, err := byZero(-1)
	if err != nil {
		panic(err)
	}
}

func byZero(n int) (r int, err error) {
	if n > 0 {
		return n, nil
	} else {
		return 0, errors.New("by zero")
	}
}
