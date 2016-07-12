package routes

import (
	dashboard_handlers "github.com/vamps-core/handlers"
)

var DashoardAppRoutes = Routes{
	Route{
		"Add Dashboard User App",
		"POST",
		"/dashboard/apps",
		true,
		dashboard_handlers.CreateDashboardApp,
	},
}
