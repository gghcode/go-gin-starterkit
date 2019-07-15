package todo

import (
	"net/http"

	"github.com/gghcode/go-gin-starterkit/app/api/common"
	"github.com/gghcode/go-gin-starterkit/middleware"
	"github.com/gin-gonic/gin"
)

// Controller handles http request.
type Controller struct {
	repo Repository
}

// NewController return new bindTodo controller instance.
func NewController(repo Repository) *Controller {
	return &Controller{
		repo: repo,
	}
}

// RegisterRoutes register handler routes.
func (controller Controller) RegisterRoutes(router gin.IRouter) {
	router.Handle("GET", "/", controller.getAllTodos)
	router.Handle("POST", "/", controller.createTodo)

	authorized := router.Use(middleware.AuthRequired())
	{
		authorized.Handle("GET", "/:id", controller.getTodoByTodoID)
		authorized.Handle("PUT", "/:id", controller.updateTodoByTodoID)
		authorized.Handle("DELETE", "/:id", controller.removeTodoByTodoID)
	}
}

// @Description Create new todo
// @Accept json
// @Produce json
// @Param payload body todo.CreateTodoRequest true "todo payload"
// @Success 201 {object} todo.TodoResponse "ok"
// @Failure 400 {object} common.ErrorResponse "Invalid todo payload"
// @Tags Todo API
// @Router /api/todos [post]
func (controller *Controller) createTodo(ctx *gin.Context) {
	var dtoReq CreateTodoRequest

	if err := ctx.ShouldBindJSON(&dtoReq); err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewErrResp(err))
		return
	}

	todoEntity := Todo{
		Title:    dtoReq.Title,
		Contents: dtoReq.Contents,
	}

	createdTodo, err := controller.repo.CreateTodo(todoEntity)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.NewErrResp(err))
		return
	}

	ctx.JSON(http.StatusCreated, createdTodo.TodoResponse())
}

// @Description Get all todos
// @Accept json
// @Produce json
// @Success 200 {array} todo.TodoResponse "ok"
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
// @Security ApiKeyAuth
// @Produce json
// @Param id path string true "Todo ID"
// @Success 200 {object} todo.TodoResponse "ok"
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

	ctx.JSON(http.StatusOK, bindTodo.TodoResponse())
}

// @Description Update todo by todo id
// @Security ApiKeyAuth
// @Produce json
// @Param id path string true "Todo ID"
// @Param payload body todo.CreateTodoRequest true "todo payload"
// @Success 200 {object} todo.TodoResponse "ok"
// @Failure 400 {object} common.ErrorResponse "Invalid todo payload"
// @Failure 404 {object} common.ErrorResponse "Not found entity"
// @Tags Todo API
// @Router /api/todos/{id} [put]
func (controller *Controller) updateTodoByTodoID(ctx *gin.Context) {
	todoID := ctx.Param("id")

	var dtoReq CreateTodoRequest
	if err := ctx.ShouldBindJSON(&dtoReq); err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewErrResp(err))
		return
	}

	todoEntity := Todo{
		Title:    dtoReq.Title,
		Contents: dtoReq.Contents,
	}

	todo, err := controller.repo.UpdateTodoByTodoID(todoID, todoEntity)
	if err == common.ErrEntityNotFound {
		ctx.JSON(http.StatusNotFound, common.NewErrResp(err))
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.NewErrResp(err))
		return
	}

	ctx.JSON(http.StatusOK, todo.TodoResponse())
}

// @Description Remove todo by todo id
// @Security ApiKeyAuth
// @Produce json
// @Param id path string true "Todo ID"
// @Success 200 {object} todo.TodoResponse "ok"
// @Failure 404 {object} common.ErrorResponse "Not found entity"
// @Tags Todo API
// @Router /api/todos/{id} [delete]
func (controller *Controller) removeTodoByTodoID(ctx *gin.Context) {
	todoID := ctx.Param("id")

	removedTodo, err := controller.repo.RemoveTodoByTodoID(todoID)
	if err == common.ErrEntityNotFound {
		ctx.JSON(http.StatusNotFound, common.NewErrResp(err))
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.NewErrResp(err))
		return
	}

	ctx.JSON(http.StatusOK, removedTodo.TodoResponse())
}
