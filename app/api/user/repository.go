package user

import (
	"github.com/gghcode/go-gin-starterkit/app/api/common"
	"github.com/gghcode/go-gin-starterkit/db"
	"github.com/jinzhu/gorm"
)

// Repository communications with db connection.
type Repository interface {
	CreateUser(User) (User, error)

	GetUserByUserName(userName string) (User, error)

	GetUserByUserID(userID int64) (User, error)

	UpdateUserByUserID(userID int64, user User) (User, error)

	RemoveUserByUserID(userID int64) (int64, error)
}

type repository struct {
	dbConn *db.Conn
}

// NewRepository return new instance.
func NewRepository(dbConn *db.Conn) Repository {
	dbConn.GetDB().AutoMigrate(User{})

	return &repository{
		dbConn: dbConn,
	}
}

func (repo *repository) CreateUser(user User) (User, error) {
	err := repo.dbConn.GetDB().
		Create(&user).
		Error

	if err != nil {
		return emptyUser, err
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
		return emptyUser, common.ErrEntityNotFound
	} else if err != nil {
		return emptyUser, err
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
		return emptyUser, common.ErrEntityNotFound
	} else if err != nil {
		return emptyUser, err
	}

	return result, nil
}

func (repo *repository) UpdateUserByUserID(userID int64, user User) (User, error) {
	entity, err := repo.GetUserByUserID(userID)
	if err != nil {
		return emptyUser, err
	}

	err = repo.dbConn.GetDB().
		Model(&entity).
		Updates(&user).
		Error

	if err != nil {
		return emptyUser, err
	}

	return entity, nil
}

func (repo *repository) RemoveUserByUserID(userID int64) (int64, error) {
	entity, err := repo.GetUserByUserID(userID)
	if err != nil {
		return emptyUser.ID, err
	}

	err = repo.dbConn.GetDB().
		Delete(&entity).
		Error

	if err != nil {
		return emptyUser.ID, err
	}

	return userID, nil
}
