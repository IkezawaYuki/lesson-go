package main

import (
	"fmt"
	"strconv"
)

func main() {
	hello := "hello"
	h := strconv.Quote(hello)
	fmt.Println(hello)
	fmt.Println(h)

	fmt.Printf("%q\n", hello)
	fmt.Printf("%s\n", hello)

	b, _ := Marshal(hello)
	fmt.Println(string(b))
}
