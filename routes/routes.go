package routes

import "github.com/vedicsoft/vamps-core/models"

var ApplicationRoutes models.Routes

func init() {
	routes := []models.Routes{
		AuthenticationRoutes,
	}

	for _, r := range routes {
		ApplicationRoutes = append(ApplicationRoutes, r...)
	}
}
