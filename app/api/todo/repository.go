package todo

import "github.com/gyuhwankim/go-gin-starterkit/db"

// Repository communications with db connection.
type Repository struct {
	dbConn *db.Conn
}

// NewRepository return new instance.
func NewRepository(dbConn *db.Conn) *Repository {
	return &Repository{
		dbConn: dbConn,
	}
}

func (repo Repository) getTodos() ([]TodoModel, error) {
	var todos []TodoModel

	db := repo.dbConn.GetDB()
	if err := db.Find(&todos).Error; err != nil {
		return nil, err
	}

	return todos, nil
}
