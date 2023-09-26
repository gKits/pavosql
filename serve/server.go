package serve

import (
	"github.com/gKits/PavoSQL/core"
	_ "net/http"
)

type Server struct {
	DBSM *core.DBMS
	Port int
}

func (s *Server) Serve() {
}
