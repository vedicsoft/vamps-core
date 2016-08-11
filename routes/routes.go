package routes

import (
	"github.com/vedicsoft/vamps-core/commons"
)

var ApplicationRoutes commons.Routes

func init() {
	routes := []commons.Routes{
		ConsoleRoutes,
		AuthenticationRoutes,
	}

	for _, r := range routes {
		ApplicationRoutes = append(ApplicationRoutes, r...)
	}
}

