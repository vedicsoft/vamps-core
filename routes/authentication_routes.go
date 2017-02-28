package routes

import (
	dashboard_handlers "github.com/vedicsoft/vamps-core/api"
	"github.com/vedicsoft/vamps-core/models"
)

var AuthenticationRoutes = models.Routes{
	models.Route{
		Name:        "Login",
		Method:      "POST",
		Pattern:     "/login",
		Secured:     false,
		HandlerFunc: dashboard_handlers.Login,
	},
	models.Route{
		Name:        "Logout",
		Method:      "POST",
		Pattern:     "/logout",
		Secured:     true,
		HandlerFunc: dashboard_handlers.Logout,
	},
	models.Route{
		Name:        "Get Customm Token",
		Method:      "GET",
		Pattern:     "/customtoken",
		Secured:     true,
		HandlerFunc: dashboard_handlers.GetCustomToken,
	},
}
