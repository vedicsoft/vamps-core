package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/vedicsoft/vamps-core/commons"
)

const (
	ROLE_LOCATION_MANAGER string = "location_manager"
	ROLE_CAPTIVE_MANAGER  string = "captive_manager"
	ROLE_ADVERT_MANAGER   string = "advert_manager"
	ROLE_POLICY_MANAGER   string = "policy_manager"
)

type Permission struct {
	permission string
}

/**
* get scope from jwt and check for permission
* "scopes": {
*    "wifi_location": [
*      "read",
*     "write",
*      "execute"
*    ]
*  }
 */
func IsAuthorized(resourceId string, permission string, r *http.Request) bool {
	m1 := make(map[string][]string)
	json.Unmarshal([]byte(r.Header.Get("scopes")), &m1)
	m2 := m1[resourceId]
	if m2 != nil {
		for _, element := range m2 {
			if element == permission {
				return true
			}
		}
	}
	return false
}

func IsUserAuthorized(username string, resourceId string, permission string, r *http.Request) bool {
	m1 := make(map[string][]string)
	json.Unmarshal([]byte(r.Header.Get("scopes")), &m1)

	m2 := m1[resourceId]
	if m2 != nil && username == r.Header.Get("username") {
		for _, element := range m2 {
			if element == permission {
				return true
			}
		}
	}
	return false
}

func HasRole(userID int, role string) (bool, error) {
	const GET_USER_ROLE string = `SELECT vs_roles.roleid from vs_roles INNER JOIN vs_user_roles ON
									 vs_roles.roleid IN (SELECT vs_user_roles.roleid FROM vs_user_roles WHERE
									 vs_user_roles.userid=?) AND vs_roles.name =?`
	var roles []int
	dbMap := commons.GetDBConnection(commons.PLATFORM_DB)
	var err error
	_, err = dbMap.Select(&roles, GET_USER_ROLE, userID, role)
	if err != nil {
		return false, err
	}
	if len(roles) > 0 {
		return true, nil
	}
	return false, nil
}

func RequireResourceAuthorization(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := r.Header.Get("username")
		tenantID, _ := strconv.Atoi(r.Header.Get("tenanid"))
		if !isAuthorized2(tenantID, username, r) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		inner.ServeHTTP(w, r)
	})
}
