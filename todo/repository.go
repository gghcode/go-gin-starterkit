package todo

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"golang-webapi-starterkit/config"
	"time"
)

type Repository struct {
	session        *mgo.Session
	dbName         string
	collectionName string
}

func NewRepository(configuration config.Configuration) *Repository {
	mongoSession, err := mgo.Dial(configuration.MongoConnectionString)
	if err != nil {
		return nil
	}

	result := Repository{
		session:        mongoSession,
		dbName:         configuration.MongoDbName,
		collectionName: "todo",
	}

	return &result
}

func (repository *Repository) Todo(id string) (Todo, error) {
	session := repository.session.Clone()
	defer session.Close()

	var result Todo

	collection := session.DB(repository.dbName).C(repository.collectionName)
	if err := collection.FindId(id).One(result); err != nil {
		return Todo{}, err
	}

	return result, nil
}

func (repository *Repository) AddTodo(todo Todo) (Todo, error) {
	sessionCopy := repository.session.Copy()
	defer sessionCopy.Close()

	todo.Id = bson.NewObjectId()
	todo.CreatedTime = time.Now()

	collection := sessionCopy.DB(repository.dbName).C(repository.collectionName)
	if err := collection.Insert(todo); err != nil {
		return Todo{}, err
	}

	return todo, nil
}

func (repository *Repository) UpdateTodo(todo Todo) (Todo, error) {
	sessionCopy := repository.session.Copy()
	defer sessionCopy.Close()

	collection := sessionCopy.DB(repository.dbName).C(repository.collectionName)
	if err := collection.UpdateId(todo.Id, todo); err != nil {
		return Todo{}, err
	}

	return todo, nil
}

func (repository *Repository) RemoveTodo(id string) error {
	sessionCopy := repository.session.Copy()
	defer sessionCopy.Close()

	collection := sessionCopy.DB(repository.dbName).C(repository.collectionName)
	if err := collection.RemoveId(id); err != nil {
		return err
	}

	return nil
}
