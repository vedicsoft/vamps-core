package api

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/vedicsoft/vamps-core/authenticator"
	"github.com/vedicsoft/vamps-core/commons"
	"net/http"
	"strings"
)

func Login(w http.ResponseWriter, r *http.Request) {
	requestUser := new(commons.SystemUser)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestUser)

	results := strings.Split(requestUser.Username, "@")
	if len(results) > 1 {
		requestUser.Username = results[0]
		requestUser.TenantDomain = results[1]
	} else {
		//setting default
		requestUser.TenantDomain = "super.com"
	}

	responseStatus, token := authenticator.Login(requestUser)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseStatus)
	w.Write(token)
}

func RefreshToken(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	requestUser := new(commons.SystemUser)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestUser)

	w.Header().Set("Content-Type", "application/json")
	w.Write(authenticator.RefreshToken(requestUser))
}

func Logout(w http.ResponseWriter, r *http.Request) {
	err := authenticator.Logout(r)
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
