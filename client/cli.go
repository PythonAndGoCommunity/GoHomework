package main

import (
	"bufio"
	"fmt"
	"os"
)

// Cli represents command line interface of a client.
type Cli struct {
	r *bufio.Reader
}

// NewCli returns an initialized Cli struct.
func NewCli() *Cli {
	return &Cli{
		r: bufio.NewReader(os.Stdin),
	}
}

// Write writes text to stdout.
func (c *Cli) Write(text string) {
	fmt.Print(text)
}

// Read returns user's input text.
func (c *Cli) Read() (string, error) {
	return c.r.ReadString('\n')
}

// Prompt formats text and writes it to stdout.
func (c *Cli) Prompt(text string) {
	fmt.Printf("%s > ", text)
}
