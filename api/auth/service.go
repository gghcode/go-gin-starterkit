package auth

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gghcode/go-gin-starterkit/api/user"
	"github.com/gghcode/go-gin-starterkit/config"
	"github.com/gghcode/go-gin-starterkit/db"
	"github.com/gghcode/go-gin-starterkit/service"
)

const prefixRefreshToken = "refresh_token"

// Service is auth authService.
type Service interface {
	VerifyAuthentication(username, password string) (user.User, error)
	GenerateAccessToken(userID int64) (string, error)
	IssueRefreshToken(userID int64) (string, error)
	VerifyRefreshToken(userID int64, refreshToken string) bool
	ExtractTokenClaims(token string) (jwt.MapClaims, error)
}

// NewService return new auth authService instance.
func NewService(
	conf config.Configuration,
	userRepo user.Repository,
	passport service.Passport,
	redisConn db.RedisConn) Service {

	return &authService{
		secretKeyBytes:      []byte(conf.Jwt.SecretKey),
		accessExpiresInSec:  time.Duration(conf.Jwt.AccessExpiresInSec),
		refreshExpiresInSec: time.Duration(conf.Jwt.RefreshExpiresInSec),
		userRepo:            userRepo,
		passport:            passport,
		redis:               redisConn,
	}
}

type authService struct {
	secretKeyBytes      []byte
	accessExpiresInSec  time.Duration
	refreshExpiresInSec time.Duration

	userRepo user.Repository
	passport service.Passport
	redis    db.RedisConn
}

func (authService *authService) VerifyAuthentication(username, password string) (user.User, error) {
	loginUser, err := authService.userRepo.GetUserByUserName(username)
	if err != nil {
		return user.EmptyUser, err
	}

	if !authService.passport.IsValidPassword(password, loginUser.PasswordHash) {
		return user.EmptyUser, ErrInvalidPassword
	}

	return loginUser, nil
}

func (authService *authService) GenerateAccessToken(userID int64) (string, error) {
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(authService.accessExpiresInSec * time.Second).Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   strconv.FormatInt(userID, 10),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(authService.secretKeyBytes)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (authService *authService) IssueRefreshToken(userID int64) (string, error) {
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(authService.refreshExpiresInSec * time.Second).Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   strconv.FormatInt(userID, 10),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(authService.secretKeyBytes)
	if err != nil {
		return "", err
	}

	if err := authService.redis.Client().Set(
		RefreshTokenRedisStorageKey(userID),
		tokenString,
		authService.refreshExpiresInSec*time.Second,
	).Err(); err != nil {
		return "", err
	}

	return tokenString, nil
}

func (authService *authService) VerifyRefreshToken(userID int64, refreshToken string) bool {
	storageKey := RefreshTokenRedisStorageKey(userID)

	token, err := authService.redis.Client().Get(storageKey).Result()
	if err != nil {
		return false
	}

	return token == refreshToken
}

func (authService *authService) ExtractTokenClaims(token string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(
		token,
		&claims,
		func(token *jwt.Token) (interface{}, error) {
			return authService.secretKeyBytes, nil
		},
	)

	if err != nil {
		validationErr, ok := err.(*jwt.ValidationError)
		if ok && validationErr.Errors == jwt.ValidationErrorExpired {
			return nil, ErrInvalidRefreshToken
		}

		return nil, ErrInvalidRefreshToken
	}

	return claims, nil
}

// RefreshTokenRedisStorageKey return key that stored on redis
func RefreshTokenRedisStorageKey(userID int64) string {
	return fmt.Sprintf("%s_%d", prefixRefreshToken, userID)
}
