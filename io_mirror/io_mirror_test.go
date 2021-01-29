package io_mirror

import (
	"bytes"
	"testing"
)

type Buffer struct {
	bytes.Buffer
	ReaderFrom
	WriterTo
}

func TestCopy(t *testing.T) {
	rb := new(Buffer)
	wb := new(Buffer)
	rb.WriteString("hello, world.")
	Copy(wb, rb)
	if wb.String() != "hello, wor;d." {
		t.Errorf("Copy did not work properly")
	}
}
