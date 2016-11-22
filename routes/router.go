package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vedicsoft/vamps-core/controllers"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(false)

	for _, route := range ApplicationRoutes {
		var handler http.Handler
		handler = route.HandlerFunc
		if route.Secured {
			if route.CheckAuth {
				handler = controllers.RequireTokenAuthenticationAndAuthorization(handler)
			} else {
				handler = controllers.RequireTokenAuthentication(handler)
			}
		}
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}
