package main

import (
	"bufio"
	"fmt"
	"net"
)

// Client connects and communicates with a server.
type Client struct {
	conn net.Conn
	w    *bufio.Writer
	ch   chan string
}

// NewClient initializes a new client.
func NewClient() *Client {
	return &Client{}
}

// Connect connects the client to a server on the specified address.
func (c *Client) Connect(addr string) (net.Addr, error) {
	conn, err := net.Dial(protocol, addr)
	if err != nil {
		return nil, err
	}
	c.conn = conn
	c.w = bufio.NewWriter(conn)
	return c.conn.RemoteAddr(), nil
}

// Close closes connection with a server.
func (c *Client) Close() {
	c.conn.Close()
}

// Send sends a message to a server.
func (c *Client) Send(msg string) error {
	_, err := fmt.Fprintf(c.conn, msg)
	if err != nil {
		return err
	}
	return nil
}
