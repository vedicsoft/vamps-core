package api

import (
	"encoding/json"
	"net/http"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/Sirupsen/logrus"
	"github.com/vamps-core/models"
	"github.com/vamps-core/controllers"
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
	if (err != nil) {
		log.Error("Error while decoding location json")
	}
	controllers.CreateNewDashboardApp(dashboardApp)
	w.WriteHeader(http.StatusOK)
}