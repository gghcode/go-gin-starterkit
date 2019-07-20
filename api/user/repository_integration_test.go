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
)

const (
	WillFetchedEntityIdx = 0
	WillUpdatedEntityIdx = 1
	WillRemovedEntityIdx = 2
)

type repoIntegration struct {
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

	suite.Run(t, new(repoIntegration))
}

func (suite *repoIntegration) SetupSuite() {
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

func (suite *repoIntegration) TearDownSuite() {
	suite.dbConn.Close()
}

func (suite *repoIntegration) TestCreateUser() {
	testCases := []struct {
		description    string
		argsUser       user.User
		expectedUserFn func(user.User) user.User
		expectedErr    error
	}{
		{
			description: "ShouldCreateUser",
			argsUser:    user.User{UserName: "newUser", PasswordHash: []byte("password")},
			expectedUserFn: func(actualUser user.User) user.User {
				user := user.User{UserName: "newUser", PasswordHash: []byte("password")}
				user.ID = actualUser.ID
				user.CreatedAt = actualUser.CreatedAt

				return user
			},
			expectedErr: nil,
		},
		{
			description: "ShouldReturnConflictErr_WhenAlreadyExistsUserName",
			argsUser: user.User{
				UserName:     suite.testUsers[WillFetchedEntityIdx].UserName,
				PasswordHash: []byte("password"),
			},
			expectedUserFn: func(user.User) user.User {
				return user.EmptyUser
			},
			expectedErr: common.ErrAlreadyExistsEntity,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			actualUser, actualErr := suite.repo.CreateUser(tc.argsUser)

			suite.Equal(tc.expectedUserFn(actualUser), actualUser)
			suite.Equal(tc.expectedErr, actualErr)
		})
	}
}

func (suite *repoIntegration) TestGetUserByUserName() {
	testCases := []struct {
		description  string
		argsUserName string
		expectedUser user.User
		expectedErr  error
	}{
		{
			description:  "ShouldFetchUser",
			argsUserName: suite.testUsers[WillFetchedEntityIdx].UserName,
			expectedUser: suite.testUsers[WillFetchedEntityIdx],
			expectedErr:  nil,
		},
		{
			description:  "ShouldReturnNotFoundErr",
			argsUserName: user.EmptyUser.UserName,
			expectedUser: user.EmptyUser,
			expectedErr:  common.ErrEntityNotFound,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			actualUser, actualErr := suite.repo.GetUserByUserName(tc.argsUserName)

			suite.Equal(tc.expectedUser, actualUser)
			suite.Equal(tc.expectedErr, actualErr)
		})
	}
}

func (suite *repoIntegration) TestGetUserByID() {
	testCases := []struct {
		description  string
		argsUserID   int64
		expectedUser user.User
		expectedErr  error
	}{
		{
			description:  "ShouldFetchUser",
			argsUserID:   suite.testUsers[WillFetchedEntityIdx].ID,
			expectedUser: suite.testUsers[WillFetchedEntityIdx],
			expectedErr:  nil,
		},
		{
			description:  "ShouldReturnNotFoundErr",
			argsUserID:   user.EmptyUser.ID,
			expectedUser: user.EmptyUser,
			expectedErr:  common.ErrEntityNotFound,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			actualUser, actualErr := suite.repo.GetUserByUserID(tc.argsUserID)

			suite.Equal(tc.expectedUser, actualUser)
			suite.Equal(tc.expectedErr, actualErr)
		})
	}
}

func (suite *repoIntegration) TestUpdateUserByID() {
	testCases := []struct {
		description  string
		argsUserID   int64
		argsUser     user.User
		expectedUser user.User
		expectedErr  error
	}{
		{
			description: "ShouldUpdateUser",
			argsUserID:  suite.testUsers[WillUpdatedEntityIdx].ID,
			argsUser: user.User{
				UserName: "willUpdateUserName",
			},
			expectedUser: user.User{
				ID:           suite.testUsers[WillUpdatedEntityIdx].ID,
				UserName:     "willUpdateUserName",
				PasswordHash: suite.testUsers[WillUpdatedEntityIdx].PasswordHash,
				CreatedAt:    suite.testUsers[WillUpdatedEntityIdx].CreatedAt,
			},
			expectedErr: nil,
		},
		{
			description:  "ShouldReturnNotFoundErr",
			argsUserID:   user.EmptyUser.ID,
			argsUser:     user.EmptyUser,
			expectedUser: user.EmptyUser,
			expectedErr:  common.ErrEntityNotFound,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			actualUser, actualErr := suite.repo.UpdateUserByUserID(
				tc.argsUserID, tc.argsUser,
			)

			suite.Equal(tc.expectedUser, actualUser)
			suite.Equal(tc.expectedErr, actualErr)
		})
	}
}

func (suite *repoIntegration) TestRemoveUserByID() {
	testCases := []struct {
		description  string
		argsUserID   int64
		expectedUser user.User
		expectedErr  error
	}{
		{
			description:  "ShouldRemoveUser",
			argsUserID:   suite.testUsers[WillRemovedEntityIdx].ID,
			expectedUser: suite.testUsers[WillRemovedEntityIdx],
			expectedErr:  nil,
		},
		{
			description:  "ShouldReturnNotFoundErr",
			argsUserID:   user.EmptyUser.ID,
			expectedUser: user.EmptyUser,
			expectedErr:  common.ErrEntityNotFound,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			actualUser, actualErr := suite.repo.RemoveUserByUserID(tc.argsUserID)

			suite.Equal(tc.expectedUser, actualUser)
			suite.Equal(tc.expectedErr, actualErr)
		})
	}
}
