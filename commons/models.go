package commons

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	Secured     bool
	HandlerFunc http.HandlerFunc
}

type Routes []Route

type SystemUser struct {
	UserId       int64    `db:"userid"json:"userid"`
	TenantId     int64    `db:"tenantid"json:"tenantid"`
	Username     string   `db:"username"json:"username"`
	TenantDomain string   `db:"domain"json:"tenantdomain"`
	Password     string   `db:"password"json:"password"`
	Email        string   `db:"email"json:"email"`
	Status       string   `db:"status"json:"status"`
	Roles        []string `json:"roles"`
}

type User struct {
	UUID     string `json:"uuid" form:"-"`
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}
