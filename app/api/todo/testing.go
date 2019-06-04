package todo

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

type mockRepository struct {
	mock.Mock
}

func (repo *mockRepository) createTodo(todo Todo) (Todo, error) {
	args := repo.Called(todo)
	return args.Get(0).(Todo), args.Error(1)
}

func (repo *mockRepository) getTodos() ([]Todo, error) {
	args := repo.Called()
	return args.Get(0).([]Todo), args.Error(1)
}

func (repo *mockRepository) getTodoByTodoID(todoID string) (Todo, error) {
	args := repo.Called(todoID)
	return args.Get(0).(Todo), args.Error(1)
}

func (repo *mockRepository) updateTodoByTodoID(todoID string, todo Todo) (Todo, error) {
	args := repo.Called(todoID, todo)
	return args.Get(0).(Todo), args.Error(1)
}

func (repo *mockRepository) removeTodoByTodoID(todoID string) (string, error) {
	args := repo.Called(todoID)
	return args.String(0), args.Error(1)
}

type mockTodoModelValidator struct {
	mock.Mock
}

func (validator *mockTodoModelValidator) Bind(c *gin.Context) error {
	args := validator.Called(c)
	return args.Error(0)
}

func (validator *mockTodoModelValidator) Todo() Todo {
	args := validator.Called()
	return args.Get(0).(Todo)
}
