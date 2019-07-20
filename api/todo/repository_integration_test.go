package todo_test

import (
	"testing"

	"github.com/gghcode/go-gin-starterkit/api/common"
	"github.com/gghcode/go-gin-starterkit/api/todo"
	"github.com/gghcode/go-gin-starterkit/config"
	"github.com/gghcode/go-gin-starterkit/db"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	WillFetchedTodoIdx = 0
	WillUpdatedTodoIdx = 1
	WillRemovedTodoIdx = 2
)

type repoIntegration struct {
	suite.Suite

	gormDB *gorm.DB
	dbConn *db.Conn

	repo todo.Repository

	testTodos []todo.Todo
}

func TestTodoRepoIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	suite.Run(t, new(repoIntegration))
}

func (suite *repoIntegration) SetupSuite() {
	conf, err := config.NewBuilder().
		BindEnvs("TEST").
		Build()

	dbConn, err := db.NewConn(conf)
	require.NoError(suite.T(), err)

	suite.dbConn = dbConn
	suite.repo = todo.NewRepository(suite.dbConn)

	suite.testTodos, err = pushTestDataToDB(suite.repo)
	require.NoError(suite.T(), err)
}

func pushTestDataToDB(repo todo.Repository) ([]todo.Todo, error) {
	todos := []todo.Todo{
		todo.Todo{Title: "will fetched todo", Contents: "first new contents"},
		todo.Todo{Title: "will updated todo", Contents: "second new contents"},
		todo.Todo{Title: "will removed todo", Contents: "third new contents"},
	}

	var result []todo.Todo

	for _, todo := range todos {
		insertedTodo, err := repo.CreateTodo(todo)
		if err != nil {
			return nil, err
		}

		result = append(result, insertedTodo)
	}

	return result, nil
}

func (suite *repoIntegration) TearDownSuite() {
	suite.dbConn.Close()
}

func (suite *repoIntegration) TestGetTodos() {
	testCases := []struct {
		description     string
		expectedTodosFn func([]todo.Todo) []todo.Todo
		expectedErr     error
	}{
		{
			description: "ShouldFetchTodos",
			expectedTodosFn: func(actualTodo []todo.Todo) []todo.Todo {
				return actualTodo
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			actualTodos, actualErr := suite.repo.GetTodos()

			suite.Equal(tc.expectedTodosFn(actualTodos), actualTodos)
			suite.Equal(tc.expectedErr, actualErr)
		})
	}
}

func (suite *repoIntegration) TestGetTodoByID() {
	testCases := []struct {
		description  string
		argsTodoID   string
		expectedTodo todo.Todo
		expectedErr  error
	}{
		{
			description:  "ShouldFetchTodo",
			argsTodoID:   suite.testTodos[WillFetchedTodoIdx].ID.String(),
			expectedTodo: suite.testTodos[WillFetchedTodoIdx],
			expectedErr:  nil,
		},
		{
			description:  "ShouldReturnNotFoundErr",
			argsTodoID:   todo.EmptyTodo.ID.String(),
			expectedTodo: todo.EmptyTodo,
			expectedErr:  common.ErrEntityNotFound,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			actualTodo, actualErr := suite.repo.GetTodoByTodoID(tc.argsTodoID)

			suite.Equal(tc.expectedTodo, actualTodo)
			suite.Equal(tc.expectedErr, actualErr)
		})
	}
}

func (suite *repoIntegration) TestCreateTodo() {
	testCases := []struct {
		description    string
		argsTodo       todo.Todo
		expectedTodoFn func(todo.Todo) todo.Todo
		expectedErr    error
	}{
		{
			description: "ShouldCreateTodo",
			argsTodo:    todo.Todo{Title: "new title", Contents: "new contents"},
			expectedTodoFn: func(insertedTodo todo.Todo) todo.Todo {
				todo := todo.Todo{Title: "new title", Contents: "new contents"}
				todo.ID = insertedTodo.ID
				todo.CreatedAt = insertedTodo.CreatedAt

				return todo
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			actualTodo, actualErr := suite.repo.CreateTodo(tc.argsTodo)

			suite.Equal(tc.expectedTodoFn(actualTodo), actualTodo)
			suite.Equal(tc.expectedErr, actualErr)
		})
	}
}

func (suite *repoIntegration) TestUpdateTodoByID() {
	testCases := []struct {
		description  string
		argsTodoID   string
		argsTodo     todo.Todo
		expectedTodo todo.Todo
		expectedErr  error
	}{
		{
			description: "ShouldUpdateTodo",
			argsTodoID:  suite.testTodos[WillUpdatedTodoIdx].ID.String(),
			argsTodo: todo.Todo{
				ID:        suite.testTodos[WillUpdatedTodoIdx].ID,
				Title:     "will update title",
				Contents:  "will update contents",
				CreatedAt: suite.testTodos[WillUpdatedTodoIdx].CreatedAt,
			},
			expectedTodo: todo.Todo{
				ID:        suite.testTodos[WillUpdatedTodoIdx].ID,
				Title:     "will update title",
				Contents:  "will update contents",
				CreatedAt: suite.testTodos[WillUpdatedTodoIdx].CreatedAt,
			},
			expectedErr: nil,
		},
		{
			description:  "ShouldReturnNotFoundErr",
			argsTodoID:   todo.EmptyTodo.ID.String(),
			argsTodo:     todo.EmptyTodo,
			expectedTodo: todo.EmptyTodo,
			expectedErr:  common.ErrEntityNotFound,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			actualTodo, actualErr := suite.repo.UpdateTodoByTodoID(
				tc.argsTodoID, tc.argsTodo,
			)

			suite.Equal(tc.expectedTodo, actualTodo)
			suite.Equal(tc.expectedErr, actualErr)
		})
	}
}

func (suite *repoIntegration) testRemoveTodoByID() {
	testCases := []struct {
		description  string
		argsTodoID   string
		expectedTodo todo.Todo
		expectedErr  error
	}{
		{
			description:  "ShouldRemoveTodo",
			argsTodoID:   suite.testTodos[WillRemovedTodoIdx].ID.String(),
			expectedTodo: suite.testTodos[WillRemovedTodoIdx],
			expectedErr:  nil,
		},
		{
			description:  "ShouldReturnNotFoundErr",
			argsTodoID:   todo.EmptyTodo.ID.String(),
			expectedTodo: todo.EmptyTodo,
			expectedErr:  common.ErrEntityNotFound,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			actualTodo, actualErr := suite.repo.RemoveTodoByTodoID(tc.argsTodoID)

			suite.Equal(tc.expectedTodo, actualTodo)
			suite.Equal(tc.expectedErr, actualErr)
		})
	}
}
