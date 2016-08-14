package api

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/vedicsoft/vamps-core/controllers"
	"github.com/vedicsoft/vamps-core/models"
	"net/http"
)

/**
* POST
* @path dashboard/apps/
*
 */
func CreateDashboardApp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var dashboardApp models.DashboardAppInfo
	err := decoder.Decode(&dashboardApp)
	if err != nil {
		log.Error("Error while decoding location json")
	}
	controllers.CreateNewDashboardApp(dashboardApp)
	w.WriteHeader(http.StatusOK)
}
