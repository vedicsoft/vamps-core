package models

import (
    "net/textproto"
)

type NameValue struct {
    Name  string  `db:"name"json:"name"`
    Value float64  `db:"value"json:"value"`
}

type Tenant struct {
    TenantId  int       `db:"tenantid"json:"tenantid"`
    Domain    string    `db:"domain"json:"domain"`
    Status    string    `db:"status"json:"status"`
    CreatedOn string    `db:"createdon"json:"createdon"`
}

type DashboardUser struct {
    UserId      int64     `db:"userid"json:"userid"`
    TenantId    int       `db:"tenantid"json:"tenantid"`
    Username    string    `db:"username"json:"username"`
    Password    string    `db:"password"json:"password"`
    Email       string    `db:"email"json:"email"`
    Status      string    `db:"status"json:"status"`
    Roles       []string  `json:"roles"`
    Permissions []Permission   `json:"permissions"`
    ApGroups    []string  `json:"apgroups"`
    SSIDs       []string   `json:"ssids"`
}

type UserInfo struct {
    TenantId    int       `db:"tenantid"json:"tenantid"`
    Username    string    `db:"username"json:"username"`
    Email       string    `db:"email"json:"email"`
    Status      string    `db:"status"json:"status"`
    Permissions []Permission   `json:"permissions"`
    ApGroups    []string  `json:"apgroups"`
}

type DashboardUserDetails struct {
    TenantId  int       `json:"tenantid"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    ContactNo string  `json:"contactno"`
}

type DashboardUserResetPassword struct {
    Username    string  `json:"username"`
    OldPassword string  `json:"oldpassword"`
    NewPassword string  `json:"newpassword"`
}


type Role struct {
    Name     string `json:"name"`
    TenantId string `json:"tenantId"`
}

type Permission struct {
    PermissionId int64    `json:"permissionid"`
    TenantId     int64        `json:"tenantid"`
    Name         string        `json:"name"`
    Action       string        `json:"action"`
}

type AuthUser struct {
    Username string `json:"username"`
    Role     Role   `json:"role"`
}

type Constrains struct {
    TenantId   int               `json:"tenantid"`
    From       string            `json:"from"`
    To         string            `json:"to"`
    PreFrom    string            `json:"prefrom"`
    PreTo      string            `json:"preto"`
    ACL        string            `json:"acl"`
    Criteria   string            `json:"criteria"`
    GroupNames []string          `json:"groupnames"`
	Parameters []string 		 `json:"parameters"`
}

type ApGroup struct {
    TenantId    int                  `db:"tenantid"json:"tenantid"`
    GroupName   string              `db:"groupname"json:"groupname"`
    GroupSymbol string              `db:"groupsymbol"json:"groupsymbol"`
}

type DashboardMetric struct {
    TenantId int                  `db:"tenantid"json:"tenantid"`
    MetricId int                  `db:"metricid"json:"metricid"`
    Name     string                      `db:"name"json:"name"`
}

type DashboardAppInfo struct {
	AppId          int64                   `db:"appid"json:"appid"`
	TenantId       int                   `db:"tenantid"json:"tenantid"`
	Aggregate      string                `db:"aggregate"json:"aggregate"`
	Name           string                `db:"name"json:"name"`
	FilterCriteria string           	 `db:"filtercriteria"json:"filtercriteria"`
	Parameters     []string				 `json:"parameters"`
	Users          []DashboardAppUser    `db:"users"json:"users"`
	Metrics        []DashboardAppMetric  `db:"metrics"json:"metrics"`
	Acls           string                `db:"acl"json:"acls"`
}

type DashboardGroups struct {
    TenantId int                  `db:"tenantid"json:"tenantid"`
    Groups   []DashboardAppGroup      `db:"groups"json:"groups"`
}

type DashboardApp struct {
    AppId     int                 `db:"appid"json:"appid"`
    TenantId  int                 `db:"tenantid"json:"tenantid"`
    Name      string              `db:"name"json:"name"`
    Aggregate string              `db:"aggregate"json:"aggregate"`
    FilterCriteria string         `db:"filtercriteria"json:"filtercriteria"`
}

type DashboardAppUser struct {
    TenantId int                  `db:"tenantid"json:"tenantid"`
    AppId    int                      `db:"appid"json:"appid"`
    UserName string                  `db:"username"json:"username"`
}

type DashboardAppMetric struct {
    MetricId int                  `db:"metricid"json:"metricid"`
    Name     string                      `db:"name"json:"name"`
}

type DashboardAppGroup struct {
    AppId     int                      `db:"appid"json:"appid"`
    GroupName string              `db:"groupname"json:"groupname"`
}

type DashboardAppAcls struct {
    AppId int                      `db:"appid"json:"appid"`
    Acls  string                     `db:"acl"json:"acls"`
}

type Response struct {
    Status  string `json:"status"`
    Message string `json:"message"`
}

type FileHeader struct {
	Filename string
	Header   textproto.MIMEHeader
	// contains filtered or unexported fields
}

