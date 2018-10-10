package app

type ApiController interface {
	RouteInfos() RouteInfos
}

func RegisterControllers(engine ServerEngine, controllers []ApiController) {
	for _, controller := range controllers {
		RegisterController(engine, controller)
	}
}

func RegisterController(engine ServerEngine, controller ApiController) {
	routeInfos := controller.RouteInfos()

	for _, route := range routeInfos {
		engine.RegisterRoute(route)
	}
}
