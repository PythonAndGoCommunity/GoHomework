package collection

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testConn struct {
	net.Conn
	id int
}

func TestConIndex__Contains__Success(t *testing.T) {
	testSlice := make([]net.Conn, 10, 10)

	neededConn := testConn{id: 123}

	for i := 0; i < len(testSlice); i++ {
		testSlice[i] = testConn{id: i}
	}

	testSlice[4] = neededConn

	assert.Equal(t, 4, ConnIndex(testSlice, neededConn))
}

func TestConIndex__DoesntContain__Fail(t *testing.T) {
	testSlice := make([]net.Conn, 10, 10)

	neededConn := testConn{id: 123}

	for i := 0; i < len(testSlice); i++ {
		testSlice[i] = testConn{id: i}
	}

	assert.Equal(t, -1, ConnIndex(testSlice, neededConn))
}
