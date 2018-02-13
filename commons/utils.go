package commons

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

const GET_TENANT_ID_FROM_DOMAIN string = "SELECT tenantid from vs_tenants WHERE domain = ?"
const REGEX_MAC string = "^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$"

func GetTenantId(r *http.Request) (int, error) {
	tenantId, err := strconv.Atoi(r.Header.Get("tenantid"))
	return tenantId, err
}

func GetUserName(r *http.Request) string {
	username := r.Header.Get("username")
	return username
}

func GetUserID(r *http.Request) (int, error) {
	userID, err := strconv.Atoi(r.Header.Get("userid"))
	return userID, err
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

// format XX-XX-XX-XX-XX-XX
func NormalizeMAC(mac string) (string, error) {
	if len(mac) > 0 {
		if len(mac) == 12 {
			return strings.ToUpper(mac[:2] + "-" + mac[2:4] + "-" + mac[4:6] + "-" + mac[6:8] + "-" +
				mac[8:10] + "-" + mac[10:12]), nil
		} else {
			return strings.ToUpper(strings.Replace(mac, ":", "-", -1)), nil
		}
	}
	return "", errors.New("empty MAC address")
}

func IsValidMAC(mac string) (bool, error) {
	match, err := regexp.MatchString("^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$", mac)
	return match, err
}
