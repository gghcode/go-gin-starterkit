package todo

import (
	"github.com/gin-gonic/gin"
	"golang-webapi-starterkit/api"
	"golang-webapi-starterkit/config"
	"net/http"
)

type Controller struct {
	repository *Repository
}

func NewController(configuration config.Configuration) *Controller {
	controller := Controller{
		repository: NewRepository(configuration),
	}

	return &controller
}

func (controller *Controller) Router() api.Router {
	return NewRouter(controller)
}

func (controller *Controller) todo(ctx *gin.Context) {
	todoId := ctx.Param("id")

	todo, err := controller.repository.Todo(todoId)
	if err != nil {
		ctx.JSON(500, err)
		return
	}

	ctx.JSON(200, todo)
}

func (controller *Controller) addTodo(ctx *gin.Context) {
	var todo Todo

	if err := ctx.Bind(&todo); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	todo, err := controller.repository.AddTodo(todo)
	if err != nil {
		ctx.AbortWithError(400, err)
		return
	}

	ctx.JSON(http.StatusCreated, todo)
}

func (controller *Controller) updateTodo(ctx *gin.Context) {
	var todo Todo

	if err := ctx.Bind(&todo); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	todo, err := controller.repository.UpdateTodo(todo)
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	ctx.JSON(http.StatusOK, todo)
}

func (controller *Controller) removeTodo(ctx *gin.Context) {
	todoId := ctx.Param("id")

	err := controller.repository.RemoveTodo(todoId)
	if err != nil {
		ctx.AbortWithError(400, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
