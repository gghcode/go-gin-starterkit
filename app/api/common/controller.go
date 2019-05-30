package common

import (
	"github.com/gin-gonic/gin"
	"github.com/gyuhwankim/go-gin-starterkit/app/api"
)

// Controller is common controller.
type Controller struct{}

// NewController return new controller instance.
func NewController() *Controller {
	return &Controller{}
}

// RegisterRoutes is method that register api routes.
func (c Controller) RegisterRoutes(handle api.HandleFunc) {
	handle("GET", "/healthy", c.getHealthy)
}

func (c *Controller) getHealthy(ctx *gin.Context) {
	ctx.Status(200)
}
