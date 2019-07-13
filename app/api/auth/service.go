package auth

import (
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gghcode/go-gin-starterkit/app/api/user"
	"github.com/gghcode/go-gin-starterkit/config"
	"github.com/gghcode/go-gin-starterkit/services"
)

// Service is auth service.
type Service interface {
	VerifyAuthentication(username, password string) (user.User, error)
	GenerateAccessToken(userID int64) (string, error)
	IssueRefreshToken(userID int64) (string, error)
}

// NewService return new auth service instance.
func NewService(
	conf config.Configuration,
	userRepo user.Repository,
	passport services.Passport) Service {

	return &service{
		secretKeyBytes:      []byte(conf.Jwt.SecretKey),
		accessExpiresInSec:  time.Duration(conf.Jwt.AccessExpiresInSec),
		refreshExpiresInSec: time.Duration(conf.Jwt.RefreshExpiresInSec),
		userRepo:            userRepo,
		passport:            passport,
	}
}

type service struct {
	secretKeyBytes      []byte
	accessExpiresInSec  time.Duration
	refreshExpiresInSec time.Duration

	userRepo user.Repository
	passport services.Passport
}

func (service *service) VerifyAuthentication(username, password string) (user.User, error) {
	loginUser, err := service.userRepo.GetUserByUserName(username)
	if err != nil {
		return user.EmptyUser, err
	}

	if !service.passport.IsValidPassword(password, loginUser.PasswordHash) {
		return user.EmptyUser, ErrInvalidPassword
	}

	return loginUser, nil
}

func (service *service) GenerateAccessToken(userID int64) (string, error) {
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(service.accessExpiresInSec * time.Second).Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   strconv.FormatInt(userID, 10),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(service.secretKeyBytes)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (service *service) IssueRefreshToken(userID int64) (string, error) {
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(service.refreshExpiresInSec * time.Second).Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   strconv.FormatInt(userID, 10),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(service.secretKeyBytes)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
