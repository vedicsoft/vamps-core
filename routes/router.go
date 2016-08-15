package routes

import (
	"github.com/gorilla/mux"
	"github.com/vedicsoft/vamps-core/controllers"
	"net/http"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(false)

	for _, route := range ApplicationRoutes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)
		if route.Secured {
			handler = controllers.RequireTokenAuthentication(handler)
		}
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}
