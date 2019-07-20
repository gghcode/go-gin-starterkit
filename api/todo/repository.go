package todo

import (
	"time"

	"github.com/gghcode/go-gin-starterkit/api/common"
	"github.com/gghcode/go-gin-starterkit/db"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// Repository communications with db connection.
type Repository interface {
	CreateTodo(todo Todo) (Todo, error)

	GetTodos() ([]Todo, error)

	GetTodoByTodoID(todoID string) (Todo, error)

	UpdateTodoByTodoID(todoID string, todo Todo) (Todo, error)

	RemoveTodoByTodoID(todoID string) (Todo, error)
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

func (repo *repository) GetTodos() ([]Todo, error) {
	var todos []Todo

	db := repo.dbConn.GetDB()
	if err := db.Find(&todos).Error; err != nil {
		return nil, err
	}

	return todos, nil
}

func (repo *repository) GetTodoByTodoID(todoID string) (Todo, error) {
	var todo Todo

	err := repo.dbConn.GetDB().
		Where("id=?", todoID).
		First(&todo).
		Error

	if err == gorm.ErrRecordNotFound {
		return EmptyTodo, common.ErrEntityNotFound
	} else if err != nil {
		return EmptyTodo, err
	}

	return todo, nil
}

func (repo *repository) CreateTodo(todo Todo) (Todo, error) {
	todo.ID = uuid.NewV4()
	todo.CreatedAt = time.Now().Unix()

	err := repo.dbConn.GetDB().
		Create(&todo).
		Error

	if err != nil {
		return EmptyTodo, err
	}

	return todo, nil
}

func (repo *repository) UpdateTodoByTodoID(todoID string, todo Todo) (Todo, error) {
	fetchedTodo, err := repo.GetTodoByTodoID(todoID)
	if err != nil {
		return EmptyTodo, err
	}

	err = repo.dbConn.GetDB().
		Model(&fetchedTodo).
		Updates(&todo).
		Error

	if err != nil {
		return EmptyTodo, err
	}

	return fetchedTodo, nil
}

func (repo *repository) RemoveTodoByTodoID(todoID string) (Todo, error) {
	todo, err := repo.GetTodoByTodoID(todoID)
	if err != nil {
		return EmptyTodo, err
	}

	err = repo.dbConn.GetDB().
		Delete(&todo).
		Error

	if err != nil {
		return EmptyTodo, err
	}

	return todo, nil
}
