package routes

import (
	dashboard_handlers "github.com/vedicsoft/vamps-core/api"
	"github.com/vedicsoft/vamps-core/commons"
)

var AuthenticationRoutes = commons.Routes{
	commons.Route{
		Name:        "Login",
		Method:      "POST",
		Pattern:     "/api/login",
		Secured:     false,
		HandlerFunc: dashboard_handlers.Login,
	},
	commons.Route{
		Name:        "Logout",
		Method:      "POST",
		Pattern:     "/api/logout",
		Secured:     true,
		HandlerFunc: dashboard_handlers.Logout,
	},
}
