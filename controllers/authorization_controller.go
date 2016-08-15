package controllers

import (
	"encoding/json"
	"net/http"
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
