package commons

import (
	"net/http"
	"strconv"
)

func GetTenantId(r *http.Request) (int, error) {
	tenantId, err := strconv.Atoi(r.Header.Get("tenantid"))
	return tenantId, err
}
