package main

import (
	"compressor/src/compressor"
	"fmt"
)

func main() {
	var c compressor.Compressor
	c.Compress("helllo")
	fmt.Println(c.Yield)
	fmt.Println(string([]byte("hello")))
}
