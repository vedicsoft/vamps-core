package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
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

func RequireResourceAuthorization(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//authBackend := InitJWTAuthenticationEngine()
		//token, err := jwt.ParseFromRequest(
		//	r,
		//	func(token *jwt.Token) (interface{}, error) {
		//		return authBackend.PublicKey, nil
		//	})
		//if err != nil || !token.Valid || authBackend.IsInBlacklist(r.Header.Get("Authorization")) {
		//	w.WriteHeader(http.StatusForbidden)
		//	return
		//} else {
		//	sClaims, _ := json.Marshal(token.Claims["scopes"])
		//	r.Header.Set("scopes", string(sClaims))
		//	r.Header.Set("username", token.Claims["sub"].(string))
		//	r.Header.Set("tenantid", strconv.FormatFloat((token.Claims["tenantid"]).(float64), 'f', 0, 64))
		//}

		username := r.Header.Get("username")
		tenantID, _ := strconv.Atoi(r.Header.Get("tenanid"))
		if !isAuthorized2(tenantID, username, r) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		inner.ServeHTTP(w, r)
	})
}
