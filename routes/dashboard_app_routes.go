package routes

import (
	dashboard_handlers "github.com/vedicsoft/vamps-core/api"
	"github.com/vedicsoft/vamps-core/commons"
)

var ConsoleRoutes = commons.Routes{
	commons.Route{
		Name:        "Add Dashboard User App",
		Method:      "POST",
		Pattern:     "/apps",
		Secured:     true,
		HandlerFunc: dashboard_handlers.CreateDashboardApp,
	},
}
