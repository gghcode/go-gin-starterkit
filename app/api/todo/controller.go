package todo

import "github.com/gin-gonic/gin"

// Controller handles http request.
type Controller struct {
	repo Repository
}

// NewController return new todo controller instance.
func NewController(repo Repository) *Controller {
	return &Controller{
		repo: repo,
	}
}

// RegisterRoutes register handler routes.
func (controller Controller) RegisterRoutes(router gin.IRouter) {
	router.Handle("GET", "/", controller.getAllTodos)
	router.Handle("GET", "/:id", controller.getTodoByTodoID)
}

func (controller *Controller) getAllTodos(ctx *gin.Context) {
	todos, err := controller.repo.getTodos()
	if err != nil {
		ctx.String(500, err.Error())
	}

	ctx.JSON(200, todos)
}

func (controller *Controller) getTodoByTodoID(ctx *gin.Context) {
	todoID := ctx.Param("id")

	todo, err := controller.repo.getTodoByTodoID(todoID)
	if err != nil {
		ctx.String(500, err.Error())
	}

	ctx.JSON(200, todo)
}
