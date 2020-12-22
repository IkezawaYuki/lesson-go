package main

import "fmt"

type TokenType int

const (
	Int TokenType = iota
	Plus
	Minus
	Lparen
	Rparen
)

type Element interface {
	Value() int
}

type Integer struct {
	value int
}

func NewInteger(value int) *Integer {
	return &Integer{value: value}
}

func (i *Integer) Value() int {
	return i.value
}

type Operation int

const (
	Addition Operation = iota
	Subtraction
)

type BinaryOperation struct {
	Type        Operation
	Left, Right Element
}

func (b *BinaryOperation) Value() int {
	switch b.Type {
	case Addition:
		return b.Left.Value() + b.Right.Value()
	case Subtraction:
		return b.Left.Value() - b.Right.Value()
	default:
		panic("Unsupported operation")
	}
}

type Token struct {
	Type TokenType
	Text string
}

func (t *Token) String() string {
	return fmt.Sprintf("`%s`", t.Text)
}

func Lex(input string) []Token {
	// todo
}

func main() {
	input := "(13+4)-(12-1)"
	tokens := Lex(input)
	fmt.Println(tokens)
}
