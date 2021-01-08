package main

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"text/scanner"
)

func Unmarshal(data []byte, out interface{}) (err error) {
	lex := &lexer{scan: scanner.Scanner{Mode: scanner.GoTokens}}
	lex.scan.Init(bytes.NewReader(data))
	lex.next()
	defer func() {
		if x := recover(); x != nil {
			err = fmt.Errorf("error at %s: %v", lex.scan.Position, x)
		}
	}()

}

type lexer struct {
	scan  scanner.Scanner
	token rune
}

func (lex *lexer) next() {
	lex.token = lex.scan.Scan()
}

func (lex *lexer) text() string {
	return lex.scan.TokenText()
}

func (lex *lexer) consume(want rune) {
	if lex.token != want {
		panic(fmt.Sprintf("got %q, want %q", lex.text(), want))
	}
	lex.next()
}

func read(lex *lexer, v reflect.Value) {
	switch lex.token {
	case scanner.Ident:
		switch lex.text() {
		case "nil":
			v.Set(reflect.Zero(v.Type()))
			lex.next()
			return
		case "t":
			v.SetBool(true)
			lex.next()
			return
		}
	case scanner.String:
		s, _ := strconv.Unquote(lex.text())
		v.SetString(s)
		lex.next()
		return
	case scanner.Int:
		i, _ := strconv.Atoi(lex.text())
		v.SetInt(int64(i))
		lex.next()
		return
	case scanner.Float:
		switch v.Kind() {
		case reflect.Float32:
			f, _ := strconv.ParseFloat(lex.text(), 32)
			v.SetFloat(f)
		case reflect.Float64:
			f, _ := strconv.ParseFloat(lex.text(), 64)
			v.SetFloat(f)
		default:
			panic(fmt.Sprintf("unexpected type: %d", v.Kind()))
		}
		lex.next()
		return
	case '#':
		lex.next()
		lex.next()
		lex.next()
		r := lex.text()
		lex.next()
		i := lex.text()
		lex.next()
		lex.consume(')')

		var bitSize int
		switch v.Kind() {
		case reflect.Complex64:
			bitSize = 32
		case reflect.Complex128:
			bitSize = 64
		default:
			panic(fmt.Sprintf("unexpected type: %d", v.Kind()))
		}
		fr, _ := strconv.ParseFloat(r, bitSize)
		fi, _ := strconv.ParseFloat(i, bitSize)
		v.SetComplex(complex(fr, fi))
		return
	case '(':

	}
}
