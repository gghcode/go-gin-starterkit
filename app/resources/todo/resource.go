package todo

import (
	"github.com/gyuhwankim/go-gin-starterkit/app/http"
	v1 "github.com/gyuhwankim/go-gin-starterkit/app/resources/todo/v1"
)

// NewV1Resource returns http routes.
func NewV1Resource() []http.Route {
	controller := v1.NewController()
	return v1.Routes(controller)
}
