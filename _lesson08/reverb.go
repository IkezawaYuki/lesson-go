package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func echo(c net.Conn, shout string, delay time.Duration) {
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
}

func handleConnection(c net.Conn) {
	input := bufio.NewScanner(c)
	alive := make(chan string, 1)
	go watchDog(c, alive)

	for input.Scan() {
		text := input.Text()
		alive <- text
		go echo(c, text, 1*time.Second)
	}

	close(alive)
	c.Close()
	fmt.Println("handleConn: closed")
}

func watchDog(c net.Conn, alive chan string) {
	for {
		select {
		case <-time.After(10 * time.Second):
			c.Close()
		case _, ok := <-alive:
			if !ok {
				return
			}
		}
	}
}

func main() {
	l, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConnection(conn)
	}
}
