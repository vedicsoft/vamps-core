package routes

import (
	dashboard_handlers "github.com/vedicsoft/vamps-core/api"
	"github.com/vedicsoft/vamps-core/commons"
)

var ConsoleRoutes = commons.Routes{
	commons.Route{
		"Add Dashboard User App",
		"POST",
		"/apps",
		true,
		dashboard_handlers.CreateDashboardApp,
	},
}
