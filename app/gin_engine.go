package app

import "github.com/gin-gonic/gin"

type GinEngine struct {
	gin *gin.Engine
}

func (engine *GinEngine) RegisterRoute(route RouteInfo) {
	engine.gin.Handle(route.Method, route.Path, route.Handle)
}

func (engine *GinEngine) Run(addr string) error {
	return engine.gin.Run(addr)
}

func NewGinEngine() *GinEngine {
	result := GinEngine{
		gin: gin.New(),
	}

	return &result
}
