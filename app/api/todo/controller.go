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

// @Description Create new todo
// @Accept json
// @Produce json
// @Param payload body model.CreateTodoRequest true "todo payload"
// @Success 201 {object} model.TodoResponse "ok"
// @Failure 400 {object} common.ErrorResponse "Invalid todo payload"
// @Tags Todo API
// @Router /api/todos [post]
func (controller *Controller) createTodo(ctx *gin.Context) {
	todoModelValidator := controller.todoValidatorFactory()
	if err := todoModelValidator.Bind(ctx); err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewErrResp(err))
		return
	}

	bindTodo := todoModelValidator.Todo()
	createdTodo, err := controller.repo.CreateTodo(bindTodo)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.NewErrResp(err))
		return
	}

	ctx.JSON(http.StatusCreated, createdTodo)
}

// @Description Get all todos
// @Accept json
// @Produce json
// @Success 200 {array} model.TodoResponse "ok"
// @Tags Todo API
// @Router /api/todos [get]
func (controller *Controller) getAllTodos(ctx *gin.Context) {
	todos, err := controller.repo.GetTodos()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.NewErrResp(err))
		return
	}

	ctx.JSON(http.StatusOK, todos)
}

// @Description Get todo by todo id
// @Produce json
// @Param id path string true "Todo ID"
// @Success 200 {object} model.TodoResponse "ok"
// @Failure 404 {object} common.ErrorResponse "Not found entity"
// @Tags Todo API
// @Router /api/todos/{id} [get]
func (controller *Controller) getTodoByTodoID(ctx *gin.Context) {
	todoID := ctx.Param("id")

	bindTodo, err := controller.repo.GetTodoByTodoID(todoID)
	if err == common.ErrEntityNotFound {
		ctx.JSON(http.StatusNotFound, common.NewErrResp(err))
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.NewErrResp(err))
		return
	}

	ctx.JSON(http.StatusOK, bindTodo)
}

// @Description Update todo by todo id
// @Produce json
// @Param id path string true "Todo ID"
// @Param payload body model.CreateTodoRequest true "todo payload"
// @Success 200 {object} model.TodoResponse "ok"
// @Failure 400 {object} common.ErrorResponse "Invalid todo payload"
// @Failure 404 {object} common.ErrorResponse "Not found entity"
// @Tags Todo API
// @Router /api/todos/{id} [put]
func (controller *Controller) updateTodoByTodoID(ctx *gin.Context) {
	todoID := ctx.Param("id")

	todoModelValidator := controller.todoValidatorFactory()
	if err := todoModelValidator.Bind(ctx); err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewErrResp(err))
		return
	}

	bindTodo := todoModelValidator.Todo()
	todo, err := controller.repo.UpdateTodoByTodoID(todoID, bindTodo)
	if err == common.ErrEntityNotFound {
		ctx.JSON(http.StatusNotFound, common.NewErrResp(err))
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.NewErrResp(err))
		return
	}

	ctx.JSON(http.StatusOK, todo)
}

// @Description Remove todo by todo id
// @Produce json
// @Param id path string true "Todo ID"
// @Success 200 {string} string "ok"
// @Failure 404 {object} common.ErrorResponse "Not found entity"
// @Tags Todo API
// @Router /api/todos/{id} [delete]
func (controller *Controller) removeTodoByTodoID(ctx *gin.Context) {
	todoID := ctx.Param("id")

	removedTodoID, err := controller.repo.RemoveTodoByTodoID(todoID)
	if err == common.ErrEntityNotFound {
		ctx.JSON(http.StatusNotFound, common.NewErrResp(err))
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.NewErrResp(err))
		return
	}

	ctx.JSON(http.StatusOK, removedTodoID)
}
