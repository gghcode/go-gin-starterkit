package todo

import (
	"github.com/gyuhwankim/go-gin-starterkit/app/api/common"
	"github.com/gyuhwankim/go-gin-starterkit/db"
)

// Repository communications with db connection.
type Repository interface {
	getTodos() ([]TodoModel, error)
	getTodoByTodoID(todoID string) (TodoModel, error)
}

type repository struct {
	dbConn *db.Conn
}

// NewRepository return new instance.
func NewRepository(dbConn *db.Conn) Repository {
	dbConn.GetDB().AutoMigrate(TodoModel{})

	return &repository{
		dbConn: dbConn,
	}
}

func (repo *repository) getTodos() ([]TodoModel, error) {
	var todos []TodoModel

	db := repo.dbConn.GetDB()
	if err := db.Find(&todos).Error; err != nil {
		return nil, err
	}

	return todos, nil
}

func (repo *repository) getTodoByTodoID(todoID string) (TodoModel, error) {
	var todo TodoModel

	notfound := repo.dbConn.GetDB().
		Where("id=?", todoID).
		First(&todo).
		RecordNotFound()

	if notfound {
		return TodoModel{}, common.ErrEntityNotFound
	}

	return todo, nil
}
