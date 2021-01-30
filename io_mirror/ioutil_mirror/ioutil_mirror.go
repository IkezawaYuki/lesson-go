package ioutil_mirror

import (
	"bytes"
	"io"
	"os"
	"sort"
)

func readAll(r io.Reader, capacity int64) (b []byte, err error) {
	var buf bytes.Buffer

	defer func() {
		e := recover()
		if panicErr, ok := e.(error); ok && panicErr == bytes.ErrTooLarge {
			err = panicErr
		} else {
			panic(e)
		}
	}()
	if int64(int(capacity)) == capacity {
		buf.Grow(int(capacity))
	}
	_, err = buf.ReadFrom(r)
	return buf.Bytes(), err
}

func ReadAll(r io.Reader) ([]byte, error) {
	return readAll(r, bytes.MinRead)
}

func ReadFile(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var n int64 = bytes.MinRead
	if fi, err := f.Stat(); err == nil {
		if size := fi.Size() + bytes.MinRead; size > n {
			n = size
		}
	}
	return readAll(f, n)
}

func WriteFile(filename string, data []byte, perm os.FileMode) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}

func ReadDir(dirname string) ([]os.FileInfo, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Name() < list[j].Name()
	})
	return list, nil
}

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error {
	return nil
}

func NopCloser(r io.Reader) io.ReadCloser {
	return nopCloser{r}
}

type devNull int
