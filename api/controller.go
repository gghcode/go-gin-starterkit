package api

import "github.com/gin-gonic/gin"

type Controller interface {
	Router() Router
}

func RegisterControllers(engine *gin.Engine, controllers []Controller) {
	for _, controller := range controllers {
		RegisterController(engine, controller)
	}
}

func RegisterController(engine *gin.Engine, controller Controller) {
	router := controller.Router()
	routeInfos := router.RouteInfos()

	for _, route := range routeInfos {
		engine.Handle(route.Method, route.Path, route.Handle)
	}
}
