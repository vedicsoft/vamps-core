package api

import (
	"encoding/json"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/vedicsoft/vamps-core/commons"
	"github.com/vedicsoft/vamps-core/controllers"
	"github.com/vedicsoft/vamps-core/models"
)

func Login(w http.ResponseWriter, r *http.Request) *commons.AppError {
	requestUser := new(models.SystemUser)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestUser)

	if err != nil {
		return &commons.AppError{err, "Can't display record", 500}
	}

	results := strings.Split(requestUser.Username, "@")
	if len(results) > 1 {
		requestUser.Username = results[0]
		requestUser.TenantDomain = results[1]
	} else {
		//setting default
		requestUser.TenantDomain = "super.com"
	}

	responseStatus, token := controllers.Login(requestUser)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseStatus)
	w.Write(token)
	return &commons.AppError{err, "Can't display record", 500}
}

func RefreshToken(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	requestUser := new(models.SystemUser)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestUser)

	w.Header().Set("Content-Type", "application/json")
	w.Write(controllers.RefreshToken(requestUser))
}

func Logout(w http.ResponseWriter, r *http.Request) {
	err := controllers.Logout(r)
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
