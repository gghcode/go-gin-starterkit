package todo

import (
	"github.com/gyuhwankim/go-gin-starterkit/app/api/common"
	"github.com/gyuhwankim/go-gin-starterkit/db"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// Repository communications with db connection.
type Repository interface {
	createTodo(todo Todo) (Todo, error)

	getTodos() ([]Todo, error)

	getTodoByTodoID(todoID string) (Todo, error)

	updateTodoByTodoID(todoID string, todo Todo) (Todo, error)

	removeTodoByTodoID(todoID string) (string, error)
}

type repository struct {
	dbConn *db.Conn
}

// NewRepository return new instance.
func NewRepository(dbConn *db.Conn) Repository {
	dbConn.GetDB().AutoMigrate(Todo{})

	return &repository{
		dbConn: dbConn,
	}
}

func (repo *repository) getTodos() ([]Todo, error) {
	var todos []Todo

	db := repo.dbConn.GetDB()
	if err := db.Find(&todos).Error; err != nil {
		return nil, err
	}

	return todos, nil
}

func (repo *repository) getTodoByTodoID(todoID string) (Todo, error) {
	var todo Todo

	err := repo.dbConn.GetDB().
		Where("id=?", todoID).
		First(&todo).
		Error

	if err == gorm.ErrRecordNotFound {
		return Todo{}, common.ErrEntityNotFound
	} else if err != nil {
		return Todo{}, err
	}

	return todo, nil
}

func (repo *repository) createTodo(todo Todo) (Todo, error) {
	todo.ID = uuid.NewV4()

	err := repo.dbConn.GetDB().
		Create(&todo).
		Error

	if err != nil {
		return Todo{}, err
	}

	return todo, nil
}

func (repo *repository) updateTodoByTodoID(todoID string, todo Todo) (Todo, error) {
	fetchedTodo, err := repo.getTodoByTodoID(todoID)
	if err != nil {
		return Todo{}, err
	}

	err = repo.dbConn.GetDB().
		Model(&fetchedTodo).
		Updates(&todo).
		Error

	if err != nil {
		return Todo{}, err
	}

	return todo, nil
}

func (repo *repository) removeTodoByTodoID(todoID string) (string, error) {
	todo, err := repo.getTodoByTodoID(todoID)
	if err != nil {
		return "", err
	}

	err = repo.dbConn.GetDB().
		Delete(&todo).
		Error

	if err != nil {
		return "", err
	}

	return todoID, nil
}
