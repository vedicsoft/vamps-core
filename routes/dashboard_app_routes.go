package routes

import (
	dashboard_handlers "github.com/vamps-core/api"
)

var ConsoleRoutes = Routes{
	Route{
		"Add Dashboard User App",
		"POST",
		"/apps",
		true,
		dashboard_handlers.CreateDashboardApp,
	},
}
