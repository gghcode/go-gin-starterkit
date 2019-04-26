package v1

import (
	"github.com/gyuhwankim/go-gin-starterkit/app/http"
)

// Routes returns slice that contain http route.
func Routes(controller *Controller) []http.Route {
	return []http.Route{
		http.Route{
			Method:  "GET",
			Path:    "/",
			Handler: controller.getHandler,
		},
	}
}
