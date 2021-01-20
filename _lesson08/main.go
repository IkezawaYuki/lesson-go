package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

var done = make(chan struct{})
var wg sync.WaitGroup

func cancelled(done <-chan struct{}) bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}

func main() {
	start := time.Now()
	ch := make(chan string)
	for _, url := range os.Args[1:] {
		wg.Add(1)
		go fetch(url, ch)
	}
}

func fetch(url string, ch chan<- string) {
	defer wg.Done()
	start := time.Now()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	if cancelled(done) {
		return
	}
	cancelChan := make(chan struct{})
	req.Cancel = cancelChan
	go func() {
		select {
		case <-done:
			close(cancelChan)
		}
	}()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		select {
		case ch <- fmt.Sprint(err):
		case <-done:
		}
		return
	}

}
