package todo

import (
	"github.com/gin-gonic/gin"
)

type todoModelValidator struct {
	BindModel struct {
		Title    string `json:"title" binding:"exists,min=4,max=100"`
		Contents string `json:"contents" binding:"exists,min=2,max=2048"`
	}

	todo Todo
}

func newTodoModelValidator() todoModelValidator {
	return todoModelValidator{}
}

func (validator *todoModelValidator) Bind(c *gin.Context) error {
	if err := c.ShouldBindJSON(&validator.BindModel); err != nil {
		return err
	}

	validator.todo.Title = validator.BindModel.Title
	validator.todo.Contents = validator.BindModel.Contents

	return nil
}
