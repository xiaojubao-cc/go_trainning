package main

import (
	"fmt"
	"time"
)

const (
	DATEFORMAT = "2006-01-02 15:04:05"
	LAYOUT     = "20060102150405"
)

func main() {
	fmt.Println(time.Now().Format(DATEFORMAT))
	fmt.Println(time.Now().Add(1 * time.Hour).Format(DATEFORMAT))
	start, _ := time.Parse(LAYOUT, "20250603165341")
	end, _ := time.Parse(LAYOUT, "20250603185341")
	fmt.Println(int(end.Sub(start).Abs() / time.Hour))
	duration, _ := time.ParseDuration("2" + "ms")
	fmt.Println(duration)
}
