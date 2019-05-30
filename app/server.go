package app

import (
	"github.com/gin-gonic/gin"
	"github.com/gyuhwankim/go-gin-starterkit/app/api"
	"github.com/gyuhwankim/go-gin-starterkit/config"
)

// Server is api-server instance.
// it contains gin.Engine, middlewares, configuration.
type Server struct {
	core *gin.Engine
	conf config.Configuration
}

// NewServer return new server instance.
func NewServer(conf config.Configuration) *Server {
	server := Server{
		core: gin.New(),
		conf: conf,
	}

	return &server
}

// Run start to listen.
func (server *Server) Run() error {
	addr := server.conf.Addr
	return server.core.Run(addr)
}

func registerController(core *gin.Engine, c []api.Controller) {
	for _, item := range c {
		item.RegisterRoutes(core.Handle)
	}
}
