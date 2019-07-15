package middleware

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gghcode/go-gin-starterkit/app/api/common"
	"github.com/gghcode/go-gin-starterkit/config"
	"github.com/gin-gonic/gin"
)

// VerifyHandlerKey is key that identify inner handler.
const VerifyHandlerKey = "INNER_FUNC_AUTH_REQUIRED"

var (
	// ErrTokenExpired is occurred when token expired
	ErrTokenExpired = errors.New("Token expired")

	// ErrUnauthorizedToken is occurred when token is invalid
	ErrUnauthorizedToken = errors.New("Token is unauthorized")
)

// AddAuthHandler is
func AddAuthHandler(conf config.JwtConfig) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var innerHandler gin.HandlerFunc = func(ctx *gin.Context) {
			token := ctx.GetHeader("Authorization")

			claims, err := verifyAccessToken(conf.SecretKey, token)
			if err != nil {
				ctx.AbortWithStatusJSON(
					http.StatusUnauthorized,
					common.NewErrResp(err),
				)
				return
			}

			userID, _ := strconv.ParseInt(claims["sub"].(string), 10, 64)

			ctx.Set("user_id", userID)
			ctx.Next()
		}

		ctx.Set(VerifyHandlerKey, innerHandler)
		ctx.Next()
	}
}

// AuthRequired Middleware
func AuthRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		handler := ctx.MustGet(VerifyHandlerKey).(gin.HandlerFunc)
		handler(ctx)

		ctx.Next()
	}
}

func verifyAccessToken(secret string, accessToken string) (jwt.MapClaims, error) {
	tokenInfo := strings.Split(accessToken, " ")
	if len(tokenInfo) != 2 || tokenInfo[0] != "Bearer" {
		return nil, ErrUnauthorizedToken
	}

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(
		tokenInfo[1],
		&claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
	)

	if err != nil {
		validationErr, ok := err.(*jwt.ValidationError)
		if ok && validationErr.Errors == jwt.ValidationErrorExpired {
			return nil, ErrTokenExpired
		}

		return nil, ErrUnauthorizedToken
	}

	return claims, nil
}
