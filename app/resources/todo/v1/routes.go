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
			Handler: controller.getAllHandler,
		},
		http.Route{
			Method:  "GET",
			Path:    "/:id",
			Handler: controller.getByTodoIDHandler,
		},
		http.Route{
			Method:  "POST",
			Path:    "/",
			Handler: controller.createHandler,
		},
		http.Route{
			Method:  "PUT",
			Path:    "/:id",
			Handler: controller.updateByTodoIDHandler,
		},
		http.Route{
			Method:  "DELETE",
			Path:    "/:id",
			Handler: controller.removeByTodoIDHandler,
		},
	}
}
