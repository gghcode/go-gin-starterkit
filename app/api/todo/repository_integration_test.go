package todo

import (
	"testing"

	"github.com/gghcode/go-gin-starterkit/app/api/common"
	"github.com/gghcode/go-gin-starterkit/config"
	"github.com/gghcode/go-gin-starterkit/db"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	WillFetchedTodoIdx = 0
	WillUpdatedTodoIdx = 1
	WillRemovedTodoIdx = 2
)

type repoIntegrationSuite struct {
	suite.Suite

	gormDB *gorm.DB
	dbConn *db.Conn

	repo Repository

	testTodos []Todo
}

func TestTodoRepoIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	suite.Run(t, new(repoIntegrationSuite))
}

func (suite *repoIntegrationSuite) SetupSuite() {
	conf, err := config.NewBuilder().
		BindEnvs("TEST").
		Build()

	dbConn, err := db.NewConn(conf)
	require.NoError(suite.T(), err)

	suite.dbConn = dbConn
	suite.repo = NewRepository(suite.dbConn)

	suite.testTodos, err = pushTestDataToDB(suite.repo)
	require.NoError(suite.T(), err)
}

func pushTestDataToDB(repo Repository) ([]Todo, error) {
	todos := []Todo{
		Todo{Title: "will fetched todo", Contents: "first new contents"},
		Todo{Title: "will updated todo", Contents: "second new contents"},
		Todo{Title: "will removed todo", Contents: "third new contents"},
	}

	var result []Todo

	for _, todo := range todos {
		insertedTodo, err := repo.CreateTodo(todo)
		if err != nil {
			return nil, err
		}

		result = append(result, insertedTodo)
	}

	return result, nil
}

func (suite *repoIntegrationSuite) TearDownSuite() {
	suite.dbConn.Close()
}

func (suite *repoIntegrationSuite) TestGetTodosExpectTodosFetched() {
	actualTodos, err := suite.repo.GetTodos()
	require.NoError(suite.T(), err)

	assert.NotNil(suite.T(), actualTodos)
}

func (suite *repoIntegrationSuite) TestGetTodoByIDExpectTodoFetched() {
	expectedTodo := suite.testTodos[WillFetchedTodoIdx]

	actualTodo, err := suite.repo.GetTodoByTodoID(expectedTodo.ID.String())
	require.NoError(suite.T(), err)

	expectedTodo.CreatedAt = actualTodo.CreatedAt

	assert.Equal(suite.T(), expectedTodo, actualTodo)
}

func (suite *repoIntegrationSuite) TestGetTodoByIDExpectNotFoundErrReturn() {
	expectedError := common.ErrEntityNotFound
	notExistsTodoID := uuid.Nil

	_, actualError := suite.repo.GetTodoByTodoID(notExistsTodoID.String())

	assert.Equal(suite.T(), expectedError, actualError)
}

func (suite *repoIntegrationSuite) TestCreateTodoExpectTodoCreated() {
	expectedTodo := Todo{
		Title:    "new title",
		Contents: "new contents",
	}

	actualTodo, err := suite.repo.CreateTodo(expectedTodo)
	require.NoError(suite.T(), err)

	expectedTodo.ID = actualTodo.ID
	expectedTodo.CreatedAt = actualTodo.CreatedAt

	assert.Equal(suite.T(), expectedTodo, actualTodo)
}

func (suite *repoIntegrationSuite) TestUpdateTodoByIDExpectTodoUpdated() {
	expectedTodo := suite.testTodos[WillUpdatedTodoIdx]
	expectedTodo.Title = "updated title"
	expectedTodo.Contents = "updated contents"

	actualTodo, err := suite.repo.UpdateTodoByTodoID(
		expectedTodo.ID.String(), expectedTodo)
	require.NoError(suite.T(), err)

	assert.Equal(suite.T(), expectedTodo, actualTodo)
}

func (suite *repoIntegrationSuite) TestUpdateTodoByIDExpectNotFoundErrReturn() {
	expectedError := common.ErrEntityNotFound
	notExistsTodoID := uuid.Nil

	_, actualError := suite.repo.UpdateTodoByTodoID(
		notExistsTodoID.String(), Todo{})

	assert.Equal(suite.T(), expectedError, actualError)
}

func (suite *repoIntegrationSuite) TestRemoveTodoByIDTodoExpectTodoRemoved() {
	expectedTodo := suite.testTodos[WillRemovedTodoIdx]

	actualTodo, err := suite.repo.RemoveTodoByTodoID(expectedTodo.ID.String())
	require.NoError(suite.T(), err)

	assert.Equal(suite.T(), expectedTodo, actualTodo)
}

func (suite *repoIntegrationSuite) TestRemoveTodoByIDExpectNotFoundErrReturn() {
	expectedError := common.ErrEntityNotFound
	notExistsTodoID := uuid.Nil

	_, actualError := suite.repo.RemoveTodoByTodoID(notExistsTodoID.String())

	assert.Equal(suite.T(), expectedError, actualError)
}
