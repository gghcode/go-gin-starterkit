package app

import (
	"github.com/gghcode/go-gin-starterkit/app/api/common"
	"github.com/gghcode/go-gin-starterkit/app/api/todo"
	"github.com/gghcode/go-gin-starterkit/config"
	"github.com/gghcode/go-gin-starterkit/db"
	_ "github.com/gghcode/go-gin-starterkit/docs"
	"github.com/gin-gonic/gin"

	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/gin-swagger/swaggerFiles"
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

	return &server
}

// Run start to listen.
func (server *Server) Run() error {
	dbConn, err := db.NewConn(server.conf)
	if err != nil {
		return err
	}

	attachSwaggerUI(server.core)

	registerControllerPrefix(server.core, "api", common.NewController())
	registerControllerPrefix(server.core, "api/todos", todo.NewController(todo.NewRepository(dbConn)))

	addr := server.conf.Addr
	return server.core.Run(addr)
}

func attachSwaggerUI(core *gin.Engine) {
	core.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func registerController(core *gin.Engine, controller Controller) {
	registerControllerPrefix(core, "", controller)
}

func registerControllerPrefix(core *gin.Engine, prefixPath string, controller Controller) {
	controller.RegisterRoutes(core.Group(prefixPath))
}
