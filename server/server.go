package server

import (
	"bufio"
	"fmt"
	"net"
	_ "net/http"

	"github.com/gKits/PavoSQL/dbms"
	"github.com/gKits/PavoSQL/parse"
)

type Server struct {
	Addr     string
	Port     uint16
	listener net.Listener
	dbms     dbms.DBMS
	parse    parse.Parser
	count    int
}

func NewServer(addr string, port uint16) (Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		return Server{}, err
	}

	// TODO: implement sane and correct defaults
	return Server{
		Addr:     addr,
		Port:     port,
		listener: listener,
		dbms:     dbms.DBMS{},
		parse:    parse.Parser{},
		count:    0,
	}, nil
}

func (s *Server) Run() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			continue
		}
		go s.handleConnection(conn)
		s.count++
	}
}

func (s *Server) Stop() {
	s.listener.Close()
	// TODO: Implement graceful shutdown
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		msg := scanner.Text()

		switch msg {
		case "exit":

		case "h", "help":
		default:
			// TODO: Implement parse, operate, respond loop
		}
	}
	s.count--
}
