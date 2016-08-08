package routes

import (
	dashboard_handlers "github.com/vedicsoft/vamps-core/api"
)

var AuthenticationRoutes = Routes{
	Route{
		"Login",
		"POST",
		"/login",
		false,
		dashboard_handlers.Login,
	},
	Route{
		"Logout",
		"POST",
		"/logout",
		true,
		dashboard_handlers.Logout,
	},
}
