package app

import (
	"path"

	"github.com/gin-gonic/gin"
	"github.com/gyuhwankim/go-gin-starterkit/app/http"
	"github.com/gyuhwankim/go-gin-starterkit/app/resources/todo"
	"github.com/gyuhwankim/go-gin-starterkit/config"
)

// Server is api-server instance.
// it contains gin.Engine, middlewares, configuration.
type Server struct {
	core *gin.Engine
	conf config.Configuration
}

// New return new server instance.
func New(conf config.Configuration) *Server {
	server := Server{
		core: gin.New(),
		conf: conf,
	}

	registerDefaultRoutes(server.core)

	return &server
}

// Run start to listen.
func (server *Server) Run() error {
	registerResource(server.core, "/api/v1/todos", todo.NewV1Resource())

	addr := server.conf.Addr
	return server.core.Run(addr)
}

func registerResource(core *gin.Engine, basePath string, routes []http.Route) {
	for _, route := range routes {
		core.Handle(route.Method, path.Join(basePath, route.Path), route.Handler)
	}
}

func registerDefaultRoutes(core *gin.Engine) {
	core.GET("ping", func(ctx *gin.Context) {
		ctx.String(200, "pong")
	})
}
