package commons

import (
	log "github.com/Sirupsen/logrus"
	"net/http"
	"strconv"
)

func GetTenantId(r *http.Request) int {
	tenantId, err := strconv.Atoi(r.Header.Get("tenantid"))
	if err != nil {
		log.Error("Error while reading tenantid from request header", err)
	}
	return tenantId
}
