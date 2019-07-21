package auth

import (
	"strconv"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/gghcode/go-gin-starterkit/api/common"
	"github.com/gghcode/go-gin-starterkit/api/user"
	"github.com/gghcode/go-gin-starterkit/config"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type fakeUserRepo struct {
	mock.Mock
}

func (r *fakeUserRepo) CreateUser(usr user.User) (user.User, error) {
	args := r.Called(usr)
	return args.Get(0).(user.User), args.Error(1)
}

func (r *fakeUserRepo) GetUserByUserName(userName string) (user.User, error) {
	args := r.Called(userName)
	return args.Get(0).(user.User), args.Error(1)
}

func (r *fakeUserRepo) GetUserByUserID(userID int64) (user.User, error) {
	args := r.Called(userID)
	return args.Get(0).(user.User), args.Error(1)
}

func (r *fakeUserRepo) UpdateUserByUserID(userID int64, usr user.User) (user.User, error) {
	args := r.Called(userID, usr)
	return args.Get(0).(user.User), args.Error(1)
}

func (r *fakeUserRepo) RemoveUserByUserID(userID int64) (user.User, error) {
	args := r.Called(userID)
	return args.Get(0).(user.User), args.Error(1)
}

type fakePassport struct {
	mock.Mock
}

func (p *fakePassport) HashPassword(password string) ([]byte, error) {
	args := p.Called(password)
	return args.Get(0).([]byte), args.Error(1)
}

func (p *fakePassport) IsValidPassword(password string, hash []byte) bool {
	args := p.Called(password, hash)
	return args.Bool(0)
}

type fakeRedisConn struct {
}

func (r *fakeRedisConn) Client() *redis.Client {
	return nil
}

type serviceUnit struct {
	suite.Suite

	configuration config.Configuration
	authService   Service
	userRepo      fakeUserRepo
	passport      fakePassport
}

func (suite *serviceUnit) SetupTest() {
	suite.configuration = config.Configuration{
		Jwt: config.JwtConfig{
			SecretKey:           "testkey",
			AccessExpiresInSec:  300,
			RefreshExpiresInSec: 3000,
		},
	}

	suite.userRepo = fakeUserRepo{}
	suite.passport = fakePassport{}
	suite.authService = NewService(
		suite.configuration,
		&suite.userRepo,
		&suite.passport,
		&fakeRedisConn{},
	)
}

func TestAuthServiceUnit(t *testing.T) {
	suite.Run(t, new(serviceUnit))
}

func (suite *serviceUnit) TestVerifyAuthentication() {
	testCases := []struct {
		description       string
		inputUserName     string
		inputPassword     string
		stubUser          user.User
		stubErr           error
		stubPasswordValid bool
		expectedUser      user.User
		expectedErr       error
	}{
		{
			description:       "ShouldReturnSuccess",
			inputUserName:     "username",
			inputPassword:     "password",
			stubUser:          user.User{ID: 10},
			stubErr:           nil,
			stubPasswordValid: true,
			expectedUser:      user.User{ID: 10},
			expectedErr:       nil,
		},
		{
			description:   "ShouldReturnNotFoundErr",
			inputUserName: "NOT_EXISTS_USER",
			inputPassword: "password",
			stubUser:      user.EmptyUser,
			stubErr:       common.ErrEntityNotFound,
			expectedUser:  user.EmptyUser,
			expectedErr:   common.ErrEntityNotFound,
		},
		{
			description:       "ShouldReturnInvalidPasswordErr",
			inputUserName:     "username",
			inputPassword:     "invalidPassword",
			stubUser:          user.User{ID: 10},
			stubErr:           nil,
			stubPasswordValid: false,
			expectedUser:      user.EmptyUser,
			expectedErr:       ErrInvalidPassword,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			suite.userRepo.
				On("GetUserByUserName", tc.inputUserName).
				Return(tc.stubUser, tc.stubErr)

			suite.passport.
				On("IsValidPassword", tc.inputPassword, tc.stubUser.PasswordHash).
				Return(tc.stubPasswordValid)

			actualUser, actualErr := suite.authService.VerifyAuthentication(
				tc.inputUserName,
				tc.inputPassword,
			)

			suite.Equal(tc.expectedUser, actualUser)
			suite.Equal(tc.expectedErr, actualErr)
		})
	}

}

func (suite *serviceUnit) TestGenerateAccessToken() {
	userID := int64(1)

	accessToken, err := suite.authService.GenerateAccessToken(userID)
	suite.NoError(err)

	suite.T().Log(accessToken)
	suite.NotNil(accessToken)
	suite.NotEqual(accessToken, "")

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(suite.configuration.Jwt.SecretKey), nil
	})

	suite.NoError(err)
	suite.True(token.Valid)

	claims, ok := token.Claims.(jwt.MapClaims)

	suite.True(ok)
	suite.Equal(strconv.FormatInt(userID, 10), claims["sub"])

	expectedExpiresInSec := suite.configuration.Jwt.AccessExpiresInSec
	actualExpiresInSec := ActualExpiresInSec(suite.T(), claims)

	suite.Equal(expectedExpiresInSec, actualExpiresInSec)
}

func ActualExpiresInSec(t *testing.T, claims jwt.MapClaims) int64 {
	expiresAt := int64(claims["exp"].(float64))
	issuedAt := int64(claims["iat"].(float64))

	return expiresAt - issuedAt
}
