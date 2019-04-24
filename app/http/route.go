package http

import (
	"github.com/gin-gonic/gin"
)

// Route include infomation about api route.
type Route struct {
	Method  string
	Path    string
	Handler func(*gin.Context)
}
