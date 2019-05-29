package app

import (
	"github.com/gin-gonic/gin"
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

	registerDefaultRoutes(server.core)

	return &server
}

// Run start to listen.
func (server *Server) Run() error {
	addr := server.conf.Addr
	return server.core.Run(addr)
}

func registerDefaultRoutes(core *gin.Engine) {
	core.GET("healthy", func(ctx *gin.Context) {
		ctx.Status(200)
	})
}
