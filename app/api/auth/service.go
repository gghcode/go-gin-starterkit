package auth

import (
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gghcode/go-gin-starterkit/app/api/user"
	"github.com/gghcode/go-gin-starterkit/config"
	"github.com/gghcode/go-gin-starterkit/service"
)

// Service is auth authService.
type Service interface {
	VerifyAuthentication(username, password string) (user.User, error)
	GenerateAccessToken(userID int64) (string, error)
	IssueRefreshToken(userID int64) (string, error)
}

// NewService return new auth authService instance.
func NewService(
	conf config.Configuration,
	userRepo user.Repository,
	passport service.Passport) Service {

	return &authService{
		secretKeyBytes:      []byte(conf.Jwt.SecretKey),
		accessExpiresInSec:  time.Duration(conf.Jwt.AccessExpiresInSec),
		refreshExpiresInSec: time.Duration(conf.Jwt.RefreshExpiresInSec),
		userRepo:            userRepo,
		passport:            passport,
	}
}

type authService struct {
	secretKeyBytes      []byte
	accessExpiresInSec  time.Duration
	refreshExpiresInSec time.Duration

	userRepo user.Repository
	passport service.Passport
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

	return tokenString, nil
}
