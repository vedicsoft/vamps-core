package commons

import (
	"errors"
	"net/http"
	"strconv"
)

const GET_TENANT_ID_FROM_DOMAIN string = "SELECT tenantid from vs_tenants WHERE domain = ?"

func GetTenantId(r *http.Request) (int, error) {
	tenantId, err := strconv.Atoi(r.Header.Get("tenantid"))
	return tenantId, err
}

func GetTenantIDFromDomain(tenantDomain string) (int, error) {
	dbMap := GetDBConnection(PLATFORM_DB)
	tenantID, err := dbMap.SelectNullInt(GET_TENANT_ID_FROM_DOMAIN, tenantDomain)
	if err != nil {
		errMsg := "error occurred while getting tenant id for domain: " + tenantDomain + "  stack trace:" + err.Error()
		return 0, errors.New(errMsg)
	} else {
		if tenantID.Valid {
			return int(tenantID.Int64), nil
		} else {
			return 0, errors.New("error occurred while getting tenant id for domain :" + tenantDomain)
		}

	}
}
