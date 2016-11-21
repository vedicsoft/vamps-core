package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/vedicsoft/vamps-core/commons"
	"github.com/vedicsoft/vamps-core/controllers"
	"github.com/vedicsoft/vamps-core/models"
)

func Login(w http.ResponseWriter, r *http.Request) *commons.AppError {
	requestUser := new(models.SystemUser)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestUser)
	if err != nil {
		return &commons.AppError{err, "Unable to decode the json request body", 500}
	}
	results := strings.Split(requestUser.Username, "@")
	if len(results) > 1 {
		requestUser.Username = results[0]
		requestUser.TenantDomain = results[1]
	} else {
		//setting default
		requestUser.TenantDomain = "super.com"
	}
	responseStatus, token, err := controllers.Login(requestUser)
	if err != nil {
		return &commons.AppError{err, "error while authenticating user", 500}
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(responseStatus)
		w.Write(token)
	}
	return nil
}

func RefreshToken(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) *commons.AppError {
	requestUser := new(models.SystemUser)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestUser)
	if err != nil {
		return &commons.AppError{err, "Unable to decode the json request body", 500}
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(controllers.RefreshToken(requestUser))
	return nil
}

func Logout(w http.ResponseWriter, r *http.Request) *commons.AppError {
	err := controllers.Logout(r)
	if err != nil {
		return &commons.AppError{err, "Unable to logout", 500}
	} else {
		w.WriteHeader(http.StatusOK)
	}
	return nil
}
