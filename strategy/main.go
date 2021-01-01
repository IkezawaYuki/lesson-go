package main

import "strings"

type OutputFormat int

const (
	Markdown OutputFormat = iota
	Html
)

type ListStrategy interface {
	Start(builder *strings.Builder)
	End(builder *strings.Builder)
	AddListItem(builder *strings.Builder, item string)
}

type MarkdownListStrategy struct{}

func (m MarkdownListStrategy) Start(builder *strings.Builder) {
}

func (m MarkdownListStrategy) End(builder *strings.Builder) {
}

func (m MarkdownListStrategy) AddListItem(builder *strings.Builder, item string) {
	builder.WriteString(" * " + item + "\n")
}

func main() {

}
