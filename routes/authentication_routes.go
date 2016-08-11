package routes

import (
	dashboard_handlers "github.com/vedicsoft/vamps-core/api"
	"github.com/vedicsoft/vamps-core/commons"
)

var AuthenticationRoutes = commons.Routes{
	commons.Route{
		"Login",
		"POST",
		"/api/login",
		false,
		dashboard_handlers.Login,
	},
	commons.Route{
		"Logout",
		"POST",
		"/api/logout",
		true,
		dashboard_handlers.Logout,
	},
}
