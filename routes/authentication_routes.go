package routes

import (
	dashboard_handlers "github.com/vedicsoft/vamps-core/api"
)

var AuthenticationRoutes = Routes{
	Route{
		"Login",
		"POST",
		"/api/login",
		false,
		dashboard_handlers.Login,
	},
	Route{
		"Logout",
		"POST",
		"/api/logout",
		true,
		dashboard_handlers.Logout,
	},
}
