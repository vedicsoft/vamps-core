package models

import (
	"github.com/vedicsoft/vamps-core/commons"
)

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

type Subscriber struct {
	UserId        int64              `db:"userid"json:"user_id"`
	TenantId      int                `db:"tenantid"`
	Username      string             `db:"username"json:"username"`
	Password      string             `db:"password"json:"password"`
	Email         string             `db:"email"json:"email"`
	AccountStatus string             `db:"account_status"json:"accountStatus"`
	FirstName     commons.NullString `db:"first_name"json:"firstName"`
	LastName      commons.NullString `db:"last_name"json:"lastName"`
	Gender        commons.NullString `db:"gender"json:"gender"`
	BirthDay      commons.NullString `db:"birthday"json:"birthday"`
	Age           commons.NullInt64  `db:"age"json:"age"`
	AgeUpper      commons.NullInt64  `db:"age_upper"json:"ageUpper"`
	AgeLower      commons.NullInt64  `db:"age_lower"json:"ageLower"`
	Religion      commons.NullString `db:"religion"json:"religion"`
	Occupation    commons.NullString `db:"occupation"json:"occupation"`
	MaritalStatus commons.NullString `db:"marital_status"json:"maritalStatus"`
	ProfileImage  commons.NullString `db:"profile_image"json:"profileImage"`
	MobileNUmber  commons.NullString `db:"mobile_number"json:"mobileNumber"`
	AdminNotes    commons.NullString `db:"admin_notes"json:"adminNotes"`
	Created       commons.NullString `db:"created"json:"created"`
	Updated       commons.NullString `db:"updated"json:"updated"`
}

type DataTablesResponse struct {
	Draw            int          `json:"draw"`
	RecordsTotal    int64        `json:"recordsTotal"`
	RecordsFiltered int64        `json:"recordsFiltered"`
	Data            []Subscriber `json:"data"`
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
	HandlerFunc commons.AppHandler
}

type Routes []Route
