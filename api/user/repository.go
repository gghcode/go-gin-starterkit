package user

import (
	"time"

	"github.com/gghcode/go-gin-starterkit/api/common"
	"github.com/gghcode/go-gin-starterkit/db"
	"github.com/jinzhu/gorm"
	pg "github.com/lib/pq"
)

// Repository communications with db connection.
type Repository interface {
	CreateUser(User) (User, error)

	GetUserByUserName(userName string) (User, error)

	GetUserByUserID(userID int64) (User, error)

	UpdateUserByUserID(userID int64, user User) (User, error)

	RemoveUserByUserID(userID int64) (User, error)
}

type repository struct {
	dbConn *db.Conn
}

// NewRepository return new instance.
func NewRepository(dbConn *db.Conn) Repository {
	if dbConn != nil {
		dbConn.GetDB().AutoMigrate(User{})
	}

	return &repository{
		dbConn: dbConn,
	}
}

func (repo *repository) CreateUser(user User) (User, error) {
	user.CreatedAt = time.Now().Unix()

	err := repo.dbConn.GetDB().
		Create(&user).
		Error

	if pgErr, ok := err.(*pg.Error); ok {
		if pgErr.Code == "23505" {
			// handle duplicate insert
			return EmptyUser, common.ErrAlreadyExistsEntity
		}
	} else if err != nil {
		return EmptyUser, err
	}

	return user, nil
}

func (repo *repository) GetUserByUserName(userName string) (User, error) {
	var result User

	err := repo.dbConn.GetDB().
		Where("user_name=?", userName).
		First(&result).
		Error

	if err == gorm.ErrRecordNotFound {
		return EmptyUser, common.ErrEntityNotFound
	} else if err != nil {
		return EmptyUser, err
	}

	return result, nil
}

func (repo *repository) GetUserByUserID(userID int64) (User, error) {
	var result User

	err := repo.dbConn.GetDB().
		Where("id=?", userID).
		First(&result).
		Error

	if err == gorm.ErrRecordNotFound {
		return EmptyUser, common.ErrEntityNotFound
	} else if err != nil {
		return EmptyUser, err
	}

	return result, nil
}

func (repo *repository) UpdateUserByUserID(userID int64, user User) (User, error) {
	entity, err := repo.GetUserByUserID(userID)
	if err != nil {
		return EmptyUser, err
	}

	err = repo.dbConn.GetDB().
		Model(&entity).
		Updates(&user).
		Error

	if err != nil {
		return EmptyUser, err
	}

	return entity, nil
}

func (repo *repository) RemoveUserByUserID(userID int64) (User, error) {
	entity, err := repo.GetUserByUserID(userID)
	if err != nil {
		return EmptyUser, err
	}

	err = repo.dbConn.GetDB().
		Delete(&entity).
		Error

	if err != nil {
		return EmptyUser, err
	}

	return entity, nil
}
