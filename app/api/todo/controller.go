package todo

import (
	"net/http"

	"github.com/gghcode/go-gin-starterkit/app/api/common"
	"github.com/gin-gonic/gin"
)

// Controller handles http request.
type Controller struct {
	repo                 Repository
	todoValidatorFactory func() TodoModelValidator
}

// NewController return new bindTodo controller instance.
func NewController(repo Repository) *Controller {
	return &Controller{
		repo: repo,
		todoValidatorFactory: func() TodoModelValidator {
			return newTodoModelValidator()
		},
	}
}

// RegisterRoutes register handler routes.
func (controller Controller) RegisterRoutes(router gin.IRouter) {
	router.Handle("POST", "/", controller.createTodo)
	router.Handle("GET", "/", controller.getAllTodos)
	router.Handle("GET", "/:id", controller.getTodoByTodoID)
	router.Handle("PUT", "/:id", controller.updateTodoByTodoID)
	router.Handle("DELETE", "/:id", controller.removeTodoByTodoID)
}

func (controller *Controller) createTodo(ctx *gin.Context) {
	todoModelValidator := controller.todoValidatorFactory()
	if err := todoModelValidator.Bind(ctx); err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewError("error", err))
		return
	}

	bindTodo := todoModelValidator.Todo()
	createdTodo, err := controller.repo.createTodo(bindTodo)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, createdTodo)
}

func (controller *Controller) getAllTodos(ctx *gin.Context) {
	todos, err := controller.repo.getTodos()
	if err != nil {
		ctx.String(500, err.Error())
		return
	}

	ctx.JSON(200, todos)
}

func (controller *Controller) getTodoByTodoID(ctx *gin.Context) {
	todoID := ctx.Param("id")

	bindTodo, err := controller.repo.getTodoByTodoID(todoID)
	if err == common.ErrEntityNotFound {
		ctx.Status(404)
		return
	} else if err != nil {
		ctx.AbortWithStatusJSON(500, map[string]string{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(200, bindTodo)
}

func (controller *Controller) updateTodoByTodoID(ctx *gin.Context) {
	todoID := ctx.Param("id")
	bindTodo := Todo{}

	if err := ctx.BindJSON(&bindTodo); err != nil {
		ctx.AbortWithError(400, err)
	}

	bindTodo, err := controller.repo.updateTodoByTodoID(todoID, bindTodo)
	if err == common.ErrEntityNotFound {
		ctx.AbortWithError(404, err)
	} else if err != nil {
		ctx.AbortWithError(500, err)
	}

	ctx.JSON(200, bindTodo)
}

func (controller *Controller) removeTodoByTodoID(ctx *gin.Context) {
	todoID := ctx.Param("id")

	removedTodoID, err := controller.repo.removeTodoByTodoID(todoID)
	if err == common.ErrEntityNotFound {
		ctx.AbortWithError(404, err)
		return
	} else if err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	ctx.JSON(204, removedTodoID)
}
