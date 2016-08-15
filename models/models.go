package models

import "net/http"

type NameValue struct {
	Name  string  `db:"name"json:"name"`
	Value float64 `db:"value"json:"value"`
}

type Tenant struct {
	TenantId  int    `db:"tenantid"json:"tenantid"`
	Domain    string `db:"domain"json:"domain"`
	Status    string `db:"status"json:"status"`
	CreatedOn string `db:"createdon"json:"createdon"`
}

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

type WifiUser struct {
	UserId          int64  `db:"userid"json:"user_id"`
	TenantId        int    `db:"tenantid"`
	Username        string `db:"username"json:"username"`
	Password        string `db:"password"json:"password"`
	Email           string `db:"email"json:"email"`
	AccountStatus   string `db:"account_status"json:"account_status"`
	FirstName       string `db:"first_name"json:"first_name"`
	LastName        string `db:"last_name"json:"last_name"`
	Gender          string `db:"gender"json:"gender"`
	BirthDay        string `db:"birthday"json:"birthday"`
	Age             int    `db:"age"json:"age"`
	AgeUpper        int    `db:"age_upper"json:"age_upper"`
	AgeLower        int    `db:"age_lower"json:"age_lower"`
	Religion        string `db:"religion"json:"religion"`
	Occupation      string `db:"occupation"json:"occupation"`
	MaritalStatus   string `db:"marital_status"json:"marital_status"`
	ProfileImage    string `db:"profile_image"json:"profile_image"`
	MobileNUmber    string `db:"mobile_number"json:"mobile_number"`
	AdminNotes      string `db:"admin_notes"json:"admin_notes"`
	LastUpdatedTime string `db:"last_updatedtime"json:"last_updated_time"`
}

type DataTablesResponse struct {
	Draw            int        `json:"draw"`
	RecordsTotal    int64      `json:"recordsTotal"`
	RecordsFiltered int64      `json:"recordsFiltered"`
	Data            []WifiUser `json:"data"`
	Error           string
}

type Role struct {
	Name     string `json:"name"`
	TenantId string `json:"tenantId"`
}

type Permission struct {
	PermissionId int64  `json:"permissionid"`
	TenantId     int64  `json:"tenantid"`
	Name         string `json:"name"`
	Action       string `json:"action"`
}

type Route struct {
	Name        string
	Method      string
	Pattern     string
	Secured     bool
	HandlerFunc http.HandlerFunc
}

type Routes []Route
