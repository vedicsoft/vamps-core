package routes

import (
	dashboard_handlers "github.com/vamps-core/handlers"
)

var AuthenticationRoutes = Routes{
	Route{
		"Login",
		"POST",
		"/dashboard/login",
		false,
		dashboard_handlers.Login,
	},
	Route{
		"Logout",
		"POST",
		"/dashboard/logout",
		true,
		dashboard_handlers.Logout,
	},
}
