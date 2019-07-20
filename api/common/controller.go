package common

import (
	"github.com/gin-gonic/gin"
)

// Controller is common controller.
type Controller struct{}

// NewController return new controller instance.
func NewController() *Controller {
	return &Controller{}
}

// RegisterRoutes is method that register api routes.
func (c Controller) RegisterRoutes(router gin.IRouter) {
	router.Handle("GET", "/healthy", c.getHealthy)
}

// @Description Get server healthy
// @Success 200
// @Tags App API
// @Router /healthy [get]
func (c *Controller) getHealthy(ctx *gin.Context) {
	ctx.Status(200)
}
