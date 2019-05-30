package api

import "github.com/gin-gonic/gin"

// HandleFunc is function that register routes of controller.
type HandleFunc func(httpMethod, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes

// Controller is interface about api Controller.
type Controller interface {
	RegisterRoutes(HandleFunc)
}
