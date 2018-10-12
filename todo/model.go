package todo

import (
	"github.com/globalsign/mgo/bson"
	"time"
)

type Todo struct {
	Id          bson.ObjectId `bson: "_id" json:"id, omitempty"`
	Title       string        `bson:"title" json:"title"`
	Contents    string        `bson:"contents" json:"contents"`
	CreatedTime time.Time     `bson:"created_timestamp" json:"created_timestamp"`
}
