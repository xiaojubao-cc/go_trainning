package main

import "os"

func main() {
	source, err := os.Open("src.txt")
	if err != nil {
		panic(err)
	}
	defer source.Close()

	dest, err := os.Create("dst.txt")
	if err != nil {
		panic(err)
	}
	defer dest.Close()

	bs := make([]byte, 5)
	source.Read(bs)
	dest.Write(bs)
}
