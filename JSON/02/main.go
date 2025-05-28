package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	str := `[100,200,300,300.04]`
	var arr []float64
	err := json.Unmarshal([]byte(str), &arr)
	if err != nil {
		panic(err)
	}
	fmt.Println(arr)
}
