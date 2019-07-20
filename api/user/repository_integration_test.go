package user_test

import (
	"testing"

	"github.com/gghcode/go-gin-starterkit/api/common"
	"github.com/gghcode/go-gin-starterkit/api/user"
	"github.com/gghcode/go-gin-starterkit/config"
	"github.com/gghcode/go-gin-starterkit/db"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gotest.tools/assert"
)

const (
	WillFetchedEntityIdx = 0
	WillUpdatedEntityIdx = 1
	WillRemovedEntityIdx = 2
)

type repoIntegrationSuite struct {
	suite.Suite

	gormDB *gorm.DB
	dbConn *db.Conn

	repo user.Repository

	testUsers []user.User
}

func TestUserRepoIntegration(t *testing.T) {
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
	suite.repo = user.NewRepository(suite.dbConn)

	suite.testUsers, err = pushTestDataToDB(suite.repo, "repo")
	require.NoError(suite.T(), err)
}

func pushTestDataToDB(repo user.Repository, prefix string) ([]user.User, error) {
	users := []user.User{
		user.User{UserName: prefix + "willFetchedUser", PasswordHash: []byte("passwordHash")},
		user.User{UserName: prefix + "willUpdatedUser", PasswordHash: []byte("passwordHash")},
		user.User{UserName: prefix + "willRemovedUser", PasswordHash: []byte("passwordHash")},
	}

	var result []user.User

	for _, user := range users {
		insertedUser, err := repo.CreateUser(user)
		if err != nil {
			return nil, err
		}

		result = append(result, insertedUser)
	}

	return result, nil
}

func (suite *repoIntegrationSuite) TearDownSuite() {
	suite.dbConn.Close()
}

func (suite *repoIntegrationSuite) TestCreateUserExpectUserCreated() {
	expectedUser := user.User{
		UserName:     "newUser",
		PasswordHash: []byte("passwordHash"),
	}

	actualUser, err := suite.repo.CreateUser(expectedUser)
	require.NoError(suite.T(), err)

	expectedUser.ID = actualUser.ID
	expectedUser.CreatedAt = actualUser.CreatedAt

	assertUserEqual(suite.T(), expectedUser, actualUser)
}

func (suite *repoIntegrationSuite) TestCreateUserExpectAlreadyExistsErrReturn() {
	alreadyExistsUser := user.User{
		UserName:     suite.testUsers[WillFetchedEntityIdx].UserName,
		PasswordHash: []byte("password"),
	}

	expectedErr := common.ErrAlreadyExistsEntity

	_, actualErr := suite.repo.CreateUser(alreadyExistsUser)

	assert.Equal(suite.T(), expectedErr, actualErr)
}

func (suite *repoIntegrationSuite) TestGetUserByUserNameExpectUserFetched() {
	expectedUser := suite.testUsers[WillFetchedEntityIdx]

	actualUser, err := suite.repo.GetUserByUserName(expectedUser.UserName)
	require.NoError(suite.T(), err)

	assertUserEqual(suite.T(), expectedUser, actualUser)
}

func (suite *repoIntegrationSuite) TestGetUserByUserNameExpectNotFoundErrReturn() {
	notExistsUserName := "NOT_EXISTS_USER_NAME"
	expectedError := common.ErrEntityNotFound

	_, actualError := suite.repo.GetUserByUserName(notExistsUserName)

	assert.Equal(suite.T(), expectedError, actualError)
}

func (suite *repoIntegrationSuite) TestGetUserByIDExpectUserFetched() {
	expectedUser := suite.testUsers[WillFetchedEntityIdx]

	actualUser, err := suite.repo.GetUserByUserID(expectedUser.ID)
	require.NoError(suite.T(), err)

	assertUserEqual(suite.T(), expectedUser, actualUser)
}

func (suite *repoIntegrationSuite) TestGetUserByIDExpectNotFoundErrReturn() {
	notExistsUserID := user.EmptyUser.ID
	expectedError := common.ErrEntityNotFound

	_, actualError := suite.repo.GetUserByUserID(notExistsUserID)

	assert.Equal(suite.T(), expectedError, actualError)
}

func (suite *repoIntegrationSuite) TestUpdateUserByIDExpectUserUpdated() {
	expectedUser := suite.testUsers[WillUpdatedEntityIdx]
	expectedUser.UserName = "updated name"

	actualTodo, err := suite.repo.UpdateUserByUserID(expectedUser.ID, expectedUser)
	require.NoError(suite.T(), err)

	assertUserEqual(suite.T(), expectedUser, actualTodo)
}

func (suite *repoIntegrationSuite) TestUpdateUserByIDExpectNotFoundErrReturn() {
	notExistsUserID := user.EmptyUser.ID
	expectedError := common.ErrEntityNotFound

	_, actualError := suite.repo.UpdateUserByUserID(notExistsUserID, user.User{})

	assert.Equal(suite.T(), expectedError, actualError)
}

func (suite *repoIntegrationSuite) TestRemoveUserByIDExpectUserRemoved() {
	expectedUser := suite.testUsers[WillRemovedEntityIdx]

	actualUser, err := suite.repo.RemoveUserByUserID(expectedUser.ID)
	require.NoError(suite.T(), err)

	assertUserEqual(suite.T(), expectedUser, actualUser)
}

func (suite *repoIntegrationSuite) TestRemoveUserByIDExpectNotFoundErrReturn() {
	notExistsTodoID := user.EmptyUser.ID
	expectedError := common.ErrEntityNotFound

	_, actualError := suite.repo.RemoveUserByUserID(notExistsTodoID)

	assert.Equal(suite.T(), expectedError, actualError)
}

func assertUserEqual(t *testing.T, expect user.User, actual user.User) {
	assert.Equal(t, expect.ID, actual.ID)
	assert.Equal(t, expect.UserName, actual.UserName)
	assert.Equal(t, expect.CreatedAt, actual.CreatedAt)
}
