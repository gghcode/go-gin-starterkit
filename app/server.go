package app

import (
	"github.com/gghcode/go-gin-starterkit/app/api/common"
	"github.com/gghcode/go-gin-starterkit/config"
	"github.com/gin-gonic/gin"
)

// Controller is interface about api Controller.
type Controller interface {
	RegisterRoutes(gin.IRouter)
}

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

	registerControllerPrefix(server.core, "api", common.NewController())

	return &server
}

// Run start to listen.
func (server *Server) Run() error {
	addr := server.conf.Addr
	return server.core.Run(addr)
}

func registerController(core *gin.Engine, controller Controller) {
	registerControllerPrefix(core, "", controller)
}

func registerControllerPrefix(core *gin.Engine, prefixPath string, controller Controller) {
	controller.RegisterRoutes(core.Group(prefixPath))
}
