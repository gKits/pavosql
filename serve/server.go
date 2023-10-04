package serve

import (
	_ "net/http"
)

type Server struct {
	Port int
}

func (s *Server) Serve() {
}
