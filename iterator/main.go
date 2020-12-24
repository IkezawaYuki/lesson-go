package main

import "fmt"

type Person struct {
	FirstName, MiddleName, LastName string
}

func (p *Person) Names() []string {
	return []string{p.FirstName, p.MiddleName, p.LastName}
}

func main() {
	p := Person{
		"A", "B", "C",
	}
	for _, n := range p.Names() {
		fmt.Println(n)
	}
}
