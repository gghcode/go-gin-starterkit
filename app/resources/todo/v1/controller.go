package v1

import (
	"github.com/gin-gonic/gin"
)

// Controller handles http request.
type Controller struct {
}

// NewController return new todo controller instance.
func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) getHandler(ctx *gin.Context) {

}
