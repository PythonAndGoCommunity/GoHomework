package server

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"

	redislight "github.com/ITandElectronics/GoHomework"
	"github.com/ITandElectronics/GoHomework/protocol"
)

// Server tcp server. Server will handle each client into separate gorutine
type Server struct {
	l       net.Listener
	storage redislight.Storage
	port    int
}

// New contruct tcp server from provided port and storage
func New(port int, storage redislight.Storage) (*Server, error) {
	l, err := net.Listen("tcp4", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, fmt.Errorf("coundn't create tcp listener: %v", err)
	}
	return &Server{l: l, port: port, storage: storage}, nil
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				log.Println("connection has been closed by client")
				return
			}
			log.Printf("couldn't read from connection: %v\n", err)
			continue
		}
		msg, err := protocol.DecodeMessage(line)
		if err != nil {
			fmt.Fprintf(conn, "%s\r\n", err.Error())
			continue
		}
		if err := protocol.ValidateMessage(msg); err != nil {
			fmt.Fprintf(conn, "%v\r\n", err)
			log.Printf("invalid message: %v\n", err)
			continue
		}

		log.Printf("client request: %s %s %s\n", msg.Cmd, msg.Key, msg.Val)
		switch msg.Cmd {
		case "GET":
			val, err := s.storage.Get(msg.Key)
			if err != nil {
				fmt.Fprintf(conn, "(%s, absent)\r\n", err.Error())
				continue
			}
			fmt.Fprintf(conn, "(%s, present)\r\n", val)
			continue
		case "SET":
			if err := s.storage.Set(msg.Key, msg.Val); err != nil {
				fmt.Fprintf(conn, "(err, %v)\r\n", err)
				continue
			}
			fmt.Fprintf(conn, "(ok)\r\n")
			continue
		case "DEL":
			if err := s.storage.Del(msg.Key); err != nil {
				fmt.Fprintf(conn, "(ignored)\r\n")
				continue
			}
			fmt.Fprintf(conn, "(absent)\r\n")
			continue
		default:
			fmt.Fprintf(conn, "command '%s' is not supported\r\n", msg.Cmd)
			continue
		}
	}
}

func (Server) welcome(conn net.Conn) {
	fmt.Fprintf(conn, `Welcome to redis-server!
Commands should be in following format: [cmd name] [param1] [param2]
Supported commands: GET [key name], SET [key name] [value], DEL [key name]
`)
}

// Run start accepting connections from clients and handle their requests. Blocking operation
func (s *Server) Run() error {
	if s.l == nil {
		return fmt.Errorf("server is not initialized. please call server.New fist")
	}
	log.Printf("server is accepting connections on %d\n", s.port)
	for {
		conn, err := s.l.Accept()
		if err != nil {
			log.Printf("coudln't accept connection: %v\n", err)
			continue
		}
		go s.handleConnection(conn)
	}
}
