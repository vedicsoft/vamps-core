package routes

import (
	dashboard_handlers "wislabs.wifi.manager/handlers"
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
