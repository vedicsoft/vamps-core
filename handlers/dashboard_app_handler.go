package handlers

import (
	"encoding/json"
	"net/http"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/Sirupsen/logrus"
	"wislabs.wifi.manager/dao"
	"wislabs.wifi.manager/controllers/dashboard"
)

/**
* POST
* @path dashboard/apps/
*
*/
func CreateDashboardApp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var dashboardApp dao.DashboardAppInfo
	err := decoder.Decode(&dashboardApp)
	if (err != nil) {
		log.Error("Error while decoding location json")
	}
	dashboard.CreateNewDashboardApp(dashboardApp)
	w.WriteHeader(http.StatusOK)
}