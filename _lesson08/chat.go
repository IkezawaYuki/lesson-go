package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"regexp"
	"time"
)

type client struct {
	who string
	out chan<- string
}

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string)
)

const clientTimeOut = time.Minute * 5

func broadcaster() {
	clients := make(map[client]bool)
	for {
		select {
		case msg := <-messages:
			for cli := range clients {
				select {
				case cli.out <- msg:
				default:
				}
			}
		case cli := <-entering:
			for client := range clients {
				cli.out <- client.who + " is in now"
			}
			clients[cli] = true
		case cli := <-leaving:
			delete(clients, cli)
			close(cli.out)
		}
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string)
	go clientWriter(conn, ch)

	input := bufio.NewScanner(conn)

	var firstMsg string
	who := conn.RemoteAddr().String()
	if input.Scan() {
		firstMsg = input.Text()
		if user, ok := extractUser(firstMsg); ok {
			who = user
			firstMsg = ""
		}
	}

	ch <- "You are " + who
	messages <- who + " has arrived"
	entering <- client{who, ch}

	if len(firstMsg) > 0 {
		messages <- who + ": " + input.Text()
	}

	timer := time.AfterFunc(clientTimeOut, func() {
		conn.Close()
	})
	for input.Scan() {
		timer.Stop()
		messages <- who + ": " + input.Text()
		timer = time.AfterFunc(clientTimeOut, func() {
			conn.Close()
		})
	}
	timer.Stop()

	leaving <- client{who, ch}
	messages <- who + " has left"
	conn.Close()
}

var userPattern = regexp.MustCompile("<user>(.+)</user>")

func extractUser(msg string) (string, bool) {
	matches := userPattern.FindAllStringSubmatch(msg, -1)
	if matches == nil {
		return "", false
	}
	return matches[0][1], true
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConn(conn)
	}
}
