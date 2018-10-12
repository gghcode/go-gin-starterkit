package todo

import "golang-webapi-starterkit/api"

type Router struct {
	routeInfos api.RouteInfos
}

func (router *Router) RouteInfos() api.RouteInfos {
	return router.routeInfos
}

func NewRouter(controller *Controller) *Router {
	router := Router{
		routeInfos: api.RouteInfos{
			api.Route("GET", "api/v1/todos/:id", controller.todo),
			api.Route("POST", "api/v1/todos", controller.addTodo),
			api.Route("PUT", "api/v1/todos/:id", controller.updateTodo),
			api.Route("DELETE", "api/v1/todos/:id", controller.removeTodo),
		},
	}

	return &router
}
