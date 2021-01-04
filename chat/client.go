package main

import "golang.org/x/net/websocket"

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}