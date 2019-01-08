package client

import (
	"bufio"
	"io/ioutil"
	"net"
	"os"
	"testing"
)

func FakeHandler(conn net.Conn) error {
	scnnr := bufio.NewScanner(conn)
	for scnnr.Scan() {
		line := scnnr.Text()
		conn.Write([]byte(line))
	}
	return nil
}

func TestCommands(t *testing.T) {
	var (
		protocolTCP = "tcp4"
		gPort       = ":1281"
		buff        = 128
		toggle      = true
	)
	grounds := []struct {
		gIP, input, output string
	}{
		{"127.0.0.1", "da", "net"},
		{"127.0.0.1", "--connect 127.0.0.2:2011", "--connect 127.0.0.2:2011"},
		//{"127.0.0.1", "--exit", "Exit. See ya."},
		{"127.0.0.1", "hello", "hello"},
	}
	fakeserver, err := net.Listen(protocolTCP, gPort)
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		for _, ground := range grounds {
			content := []byte(ground.input)
			tmpfile, err := ioutil.TempFile("./", "test.*.txt")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tmpfile.Name()) //clean up
			if _, werr := tmpfile.Write(content); werr != nil {
				t.Fatal(werr)
			}
			if _, serr := tmpfile.Seek(0, 0); serr != nil {
				t.Fatal(serr)
			}
			oldStdin := os.Stdin
			defer func() { os.Stdin = oldStdin }() //Restore original Stdin
			total, errC := Commands(protocolTCP, gPort, ground.gIP, buff)
			if errC != nil {
				t.Fatalf("ErrorC: %v ", errC)
			}
			if cerr := tmpfile.Close(); cerr != nil {
				t.Fatal(cerr)
			}

			if total != ground.output {
				toggle = false
				t.Errorf("Wanted:%v , got:%v .", ground.output, total)
			}
		}
	}()
	for {
		conn, aeer := fakeserver.Accept()
		if aeer != nil {
			t.Fatal(aeer)
		}
		defer conn.Close()
		go FakeHandler(conn)
		if toggle == false {
			break
		}
	}
}
