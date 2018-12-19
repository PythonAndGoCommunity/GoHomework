package server

import (
	"testing"
)

func TestNew(t *testing.T) {
	server, err := New(8080, nil)
	if err != nil {
		t.Fatal(err)
	}
	if server.l == nil {
		t.Fatal("listener is nil")
	}
	if server.storage != nil {
		t.Fatal("storage should be nil")
	}
}
