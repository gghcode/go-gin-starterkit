package todo

import (
	"testing"

	"github.com/gyuhwankim/go-gin-starterkit/app/api/common"
	"github.com/gyuhwankim/go-gin-starterkit/db"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	mocket "github.com/selvatico/go-mocket"
)

type repoTestSuite struct {
	suite.Suite

	mockGormDB *gorm.DB
	mockDbConn *db.Conn

	catcher *mocket.MockCatcher

	repo Repository
}

func TestRepoTestSuite(t *testing.T) {
	suite.Run(t, new(repoTestSuite))
}

func (suite *repoTestSuite) SetupTest() {
	mocket.Catcher.Register()

	mockGormDB, err := gorm.Open(mocket.DriverName, "connectionString")

	require.NoError(suite.T(), err)

	suite.catcher = mocket.Catcher
	suite.mockGormDB = mockGormDB
	suite.mockDbConn = db.NewConn(mockGormDB)
	suite.repo = NewRepository(suite.mockDbConn)

	require.NotNil(suite.T(), suite.repo)
}

func (suite *repoTestSuite) TearDownTest() {
	suite.mockDbConn.GetDB().Close()
}

func (suite *repoTestSuite) TestShouldGetTodos() {
	expectedTodos := []Todo{
		Todo{
			ID:       uuid.NewV4(),
			Title:    "FIRST TITLE",
			Contents: "FIRST CONTENTS",
		},
		Todo{
			ID:       uuid.NewV4(),
			Title:    "SECOND TITLE",
			Contents: "SECOND CONTENTS",
		},
	}

	reply := []map[string]interface{}{}
	for _, todo := range expectedTodos {
		reply = append(reply, map[string]interface{}{
			"id":       todo.ID.String(),
			"title":    todo.Title,
			"contents": todo.Contents,
		})
	}

	suite.catcher.Reset().NewMock().
		WithQuery(`SELECT * FROM "todos"`).
		WithReply(reply)

	actualTodos, err := suite.repo.getTodos()

	require.NoError(suite.T(), err)
	require.Equal(suite.T(), len(expectedTodos), len(actualTodos))
	require.Equal(suite.T(), expectedTodos, actualTodos)
}

func (suite *repoTestSuite) TestShouldGetTodo() {
	expected := Todo{
		ID:       uuid.NewV4(),
		Title:    "todo title",
		Contents: "todo contents",
	}

	reply := []map[string]interface{}{{
		"id":       expected.ID.String(),
		"title":    expected.Title,
		"contents": expected.Contents,
	}}

	suite.catcher.Reset().NewMock().
		WithQuery(`SELECT * FROM "todos"`).
		WithArgs(expected.ID.String()).
		WithReply(reply)

	actual, err := suite.repo.getTodoByTodoID(expected.ID.String())

	require.NoError(suite.T(), err)
	require.Equal(suite.T(), expected, actual)
}

func (suite *repoTestSuite) TestShouldBeNotFound() {
	expectedError := common.ErrEntityNotFound
	notExistsTodoID := uuid.NewV4()

	suite.catcher.Reset().NewMock().
		WithQuery(`SELECT * FROM "todos"`).
		WithArgs(notExistsTodoID.String()).
		WithError(gorm.ErrRecordNotFound)

	_, actualError := suite.repo.getTodoByTodoID(notExistsTodoID.String())

	require.Equal(suite.T(), expectedError, actualError)
}

func (suite *repoTestSuite) TestShouldBeCreated() {
	willCreateTodo := Todo{
		Title:    "new title",
		Contents: "new contents",
	}

	reply := []map[string]interface{}{{
		"id":       uuid.NewV4(),
		"title":    willCreateTodo.Title,
		"contents": willCreateTodo.Contents,
	}}

	suite.catcher.Reset().NewMock().
		WithQuery(`INSERT INTO "todos"`).
		WithReply(reply)

	actual, err := suite.repo.createTodo(willCreateTodo)
	expected := Todo{
		ID:       actual.ID,
		Title:    willCreateTodo.Title,
		Contents: willCreateTodo.Contents,
	}

	require.Nil(suite.T(), err)
	require.Equal(suite.T(), expected, actual)
}

func (suite *repoTestSuite) TestShouldBeUpdated() {
	expectedTodo := Todo{
		ID:       uuid.NewV4(),
		Title:    "updated title",
		Contents: "updated contents",
	}

	suite.catcher.Reset().NewMock().
		WithArgs(expectedTodo.ID.String())

	actualTodo, err := suite.repo.updateTodoByTodoID(expectedTodo.ID.String(), expectedTodo)

	require.Nil(suite.T(), err)
	require.Equal(suite.T(), expectedTodo, actualTodo)
}

func (suite *repoTestSuite) TestShouldBeNotFoundWhenUpdateTodo() {
	expectedError := common.ErrEntityNotFound
	notExistsTodoID := uuid.NewV4()
	notExistsTodo := Todo{
		ID:       notExistsTodoID,
		Title:    "NOT EXISTS TODO",
		Contents: "NOT EXISTS CONTENTS",
	}

	suite.catcher.Reset().NewMock().
		WithError(gorm.ErrRecordNotFound)

	_, actualError := suite.repo.updateTodoByTodoID(notExistsTodoID.String(), notExistsTodo)

	require.Equal(suite.T(), expectedError, actualError)
}

func (suite *repoTestSuite) TestShouldBeDeleted() {
	expectedTodoID := uuid.NewV4().String()

	reply := []map[string]interface{}{{
		"id":       expectedTodoID,
		"title":    "title",
		"contents": "contents",
	}}

	suite.catcher.Reset().NewMock().
		WithQuery(`DELETE * FROM "todos"`).
		WithArgs(expectedTodoID).
		WithReply(reply)

	actualTodoID, err := suite.repo.removeTodoByTodoID(expectedTodoID)

	require.Nil(suite.T(), err)
	require.Equal(suite.T(), expectedTodoID, actualTodoID)
}

func (suite *repoTestSuite) TestShouldBeNotFoundWhenRemoveTodo() {
	expectedError := common.ErrEntityNotFound
	notExistsTodoID := uuid.NewV4().String()

	suite.catcher.Reset().NewMock().
		WithArgs(notExistsTodoID).
		WithError(gorm.ErrRecordNotFound)

	_, actualError := suite.repo.removeTodoByTodoID(notExistsTodoID)

	require.Equal(suite.T(), expectedError, actualError)
}
