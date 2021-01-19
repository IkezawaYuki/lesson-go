package main

import (
	"log"
	"net"
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

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
}
