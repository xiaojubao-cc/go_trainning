package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	str := `{
            "name":"golang"
            }`
	maps := make(map[string]interface{})
	err := json.Unmarshal([]byte(str), &maps)
	if err != nil {
		panic(err)
	}
	fmt.Println(maps)
}
