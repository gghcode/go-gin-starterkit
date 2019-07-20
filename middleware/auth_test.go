package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gghcode/go-gin-starterkit/api/common"
	"github.com/gghcode/go-gin-starterkit/api/testutil"
	"github.com/gghcode/go-gin-starterkit/config"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type authUnit struct {
	suite.Suite

	conf config.JwtConfig
}

func TestAuthMiddlewareUnit(t *testing.T) {
	suite.Run(t, new(authUnit))
}

func (suite *authUnit) SetupTest() {
	suite.conf = config.JwtConfig{
		SecretKey:           "testkey",
		AccessExpiresInSec:  300,
		RefreshExpiresInSec: 3000,
	}

	gin.SetMode(gin.TestMode)
}

func (suite *authUnit) TestVerifyAccessToken() {
	testCases := []struct {
		description   string
		accessTokenFn func() string
		expectedErr   error
		expectedSub   string
	}{
		{
			description: "ShouldBeSuccess",
			accessTokenFn: func() string {
				claims := &jwt.StandardClaims{
					ExpiresAt: time.Now().Add(3000 * time.Second).Unix(),
					IssuedAt:  time.Now().Unix(),
					Subject:   "10",
				}

				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString([]byte(suite.conf.SecretKey))

				return "Bearer " + tokenString
			},
			expectedErr: nil,
			expectedSub: "10",
		},
		{
			description: "ShouldTokenExpiredErrReturn",
			accessTokenFn: func() string {
				claims := &jwt.StandardClaims{
					ExpiresAt: time.Now().Add(-300 * time.Second).Unix(),
					IssuedAt:  time.Now().Unix(),
					Subject:   "10",
				}

				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString([]byte(suite.conf.SecretKey))

				return "Bearer " + tokenString
			},
			expectedErr: ErrTokenExpired,
		},
		{
			description:   "ShouldUnauthorizedTokenErrReturn_WhenEmptyToken",
			accessTokenFn: func() string { return "" },
			expectedErr:   ErrUnauthorizedToken,
		},
		{
			description:   "ShouldUnauthorizedTokenErrReturn_WhenInvalidTokenType",
			accessTokenFn: func() string { return "Bear abcd" },
			expectedErr:   ErrUnauthorizedToken,
		},
		{
			description:   "ShouldUnauthorizedTokenErrReturn_WhenShortTokenInfo",
			accessTokenFn: func() string { return "BearerValidToken" },
			expectedErr:   ErrUnauthorizedToken,
		},
		{
			description:   "ShouldUnauthorizedTokenErrReturn_WhenInvalidAccessToken",
			accessTokenFn: func() string { return "Bearer InvalidToken" },
			expectedErr:   ErrUnauthorizedToken,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			claims, actualErr := verifyAccessToken(
				suite.conf.SecretKey,
				tc.accessTokenFn(),
			)

			suite.Equal(tc.expectedErr, actualErr)

			if tc.expectedSub != "" {
				suite.Equal(tc.expectedSub, claims["sub"])
			}
		})
	}
}

func (suite *authUnit) TestAuthMiddleware() {
	testCases := []struct {
		description    string
		accessToken    func() (accessToken string)
		expectedStatus int
		expectedJSON   string
	}{
		{
			description: "ShouldBeSuccess",
			accessToken: func() string {
				expiresIn := time.Duration(suite.conf.AccessExpiresInSec)

				claims := &jwt.StandardClaims{
					ExpiresAt: time.Now().Add(expiresIn * time.Second).Unix(),
					IssuedAt:  time.Now().Unix(),
					Subject:   "10",
				}

				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString([]byte(suite.conf.SecretKey))

				return "Bearer " + tokenString
			},
			expectedStatus: http.StatusOK,
		},
		{
			description: "ShouldTokenExpiredErrReturn",
			accessToken: func() string {
				claims := &jwt.StandardClaims{
					ExpiresAt: time.Now().Add(-300 * time.Second).Unix(),
					IssuedAt:  time.Now().Unix(),
					Subject:   "10",
				}

				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString([]byte(suite.conf.SecretKey))

				return "Bearer " + tokenString
			},
			expectedStatus: http.StatusUnauthorized,
			expectedJSON: testutil.JSONStringFromInterface(suite.T(),
				common.NewErrResp(ErrTokenExpired)),
		},
		{
			description:    "ShouldUnauthorizedTokenErrReturn_WhenShortTokenInfo",
			accessToken:    func() string { return "Bearfasdf" },
			expectedStatus: http.StatusUnauthorized,
			expectedJSON: testutil.JSONStringFromInterface(suite.T(),
				common.NewErrResp(ErrUnauthorizedToken)),
		},
		{
			description:    "ShouldUnauthorizedTokenErrReturn_WhenEmptyToken",
			accessToken:    func() string { return "" },
			expectedStatus: http.StatusUnauthorized,
			expectedJSON: testutil.JSONStringFromInterface(suite.T(),
				common.NewErrResp(ErrUnauthorizedToken)),
		},
		{
			description:    "ShouldUnauthorizedTokenErrReturn_WhenInvalidAccessToken",
			accessToken:    func() string { return "Bearer InvalidToken" },
			expectedStatus: http.StatusUnauthorized,
			expectedJSON: testutil.JSONStringFromInterface(suite.T(),
				common.NewErrResp(ErrUnauthorizedToken)),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			recorder := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/", nil)
			req.Header.Add("Authorization", tc.accessToken())

			_, engine := gin.CreateTestContext(recorder)

			engine.Use(AddAuthHandler(suite.conf))
			engine.Use(AuthRequired())
			engine.GET("/", func(ctx *gin.Context) { ctx.MustGet("user_id") })
			engine.ServeHTTP(recorder, req)

			suite.Equal(tc.expectedStatus, recorder.Code)

			actualJSON := testutil.JSONStringFromResBody(suite.T(), recorder.Body)

			suite.NotNil(actualJSON)
		})
	}
}
