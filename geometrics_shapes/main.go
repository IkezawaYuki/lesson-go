package main

import "fmt"

type GraphicObject struct {
	Name, Color string
	Children    []GraphicObject
}

func main() {
	fmt.Println("hello world")
}
