package main

import (
	"testing"
	"time"
	"net"
	"bufio"
	"fmt"
)

type fakeConn struct{
	lenBufer int
	bufer [1024*10]byte
}

func (f *fakeConn) Read(b []byte) (int, error) {
	copy(b, f.bufer[0:f.lenBufer])
	n := f.lenBufer
	f.lenBufer = 0
	fmt.Println("Прочитанно ", n, " байт,", string(b[:n]))
	var err error
	err = nil
	return n, err
}

func (f *fakeConn) Write(b []byte) (int, error) {
	n := len(b)
	copy(f.bufer[f.lenBufer:n+f.lenBufer], b)
	fmt.Println("Записанно ", n, " байт,", string(f.bufer[f.lenBufer:n+f.lenBufer]))
	f.lenBufer += n
	var err error
	err = nil
	return n, err
}

func (f *fakeConn) Close() error {

	return nil

}

func (f *fakeConn) LocalAddr() net.Addr {

	return nil

}

func (f *fakeConn) RemoteAddr() net.Addr {

	return nil

}

func (f *fakeConn) SetDeadline(t time.Time) error {

	return nil

}

func (f *fakeConn) SetReadDeadline(t time.Time) error {

	return nil

}

func (f *fakeConn) SetWriteDeadline(t time.Time) error {

	return nil

}

func (f *fakeConn) SetReadBuffer(bytes int) error {

	return nil

}

func (f *fakeConn) SetWriteBuffer(bytes int) error {

	return nil

}


func TestHandleConnection(t *testing.T){
	fConn := fakeConn{}
	go handleConnection(&fConn)
	serverWriter := bufio.NewWriter(&fConn)
	serverReader := bufio.NewReader(&fConn)

	//set test
	n, err := serverWriter.WriteString("set a b\n")
	if n == 0 || err != nil {
			fmt.Println(err)
			return
	}
	serverWriter.Flush()
	reply, err := serverReader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
	}
	if reply != "Ok\n" {
		t.Error("Set test. Excepted Ok, got ", reply)
	}
}