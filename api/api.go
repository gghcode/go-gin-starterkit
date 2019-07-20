package api

import "github.com/gin-gonic/gin"

// Controller is interface about api Controller.
type Controller interface {
	RegisterRoutes(router gin.IRouter)
}

// IController is instance of Container
var IController = new(Controller)
