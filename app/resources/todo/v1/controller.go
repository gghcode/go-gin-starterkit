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

func (c *Controller) getAllHandler(ctx *gin.Context) {
	panic("Not implement...")
}

func (c *Controller) getByTodoIdHandler(ctx *gin.Context) {
	panic("Not implement...")
}

func (c *Controller) createHandler(ctx *gin.Context) {
	panic("Not implement...")
}

func (c *Controller) updateByTodoIdHandler(ctx *gin.Context) {
	panic("Not implement...")
}

func (c *Controller) removeByTodoIdHandler(ctx *gin.Context) {
	panic("Not implement...")
}

func (c *Controller) postHandler(ctx *gin.Context) {
	var todo Todo

	if err := ctx.BindJSON(&todo); err != nil {
		ctx.AbortWithStatus(400)
	}

	ctx.JSON(201, todo)
}
