package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/vedicsoft/vamps-core/commons"
	"github.com/vedicsoft/vamps-core/models"
	"github.com/dgrijalva/jwt-go"
)

type TokenContext struct {
	Token          string `json:"token" form:"token"`
	TenantId       int64  `json:"tenantid" form:"tenantid"`
	GrantType      string   `json:"grantType,omitempty"`
	UserID         int64    `json:"userID,omitempty"`
	ExpirationTime int64    `json:"expirationTime,omitempty"`
	Scopes         []string `json:"scopes,omitempty"`
}

func Login(requestUser *models.SystemUser) (int, []byte, error) {
	authEngine := InitJWTAuthenticationEngine()
	requestUser.TenantId = getTenantId(requestUser)
	authenticaed, err := authEngine.Authenticate(requestUser)
	if err != nil {
		return http.StatusInternalServerError, []byte(""), err
	}
	if authenticaed {
		token, err := authEngine.GenerateToken(requestUser, commons.ServerConfigurations.JWTExpirationDelta)
		if err != nil {
			return http.StatusInternalServerError, []byte(""), err
		} else {
			response, _ := json.Marshal(TokenContext{Token:token, TenantId:requestUser.TenantId})
			return http.StatusOK, response, nil
		}
	}
	return http.StatusUnauthorized, []byte(""), nil
}

func CustomToken(requestUser *models.SystemUser, expirationHours int) []byte {
	authEngine := InitJWTAuthenticationEngine()
	token, err := authEngine.GenerateToken(requestUser, expirationHours)
	if err != nil {
		panic(err)
	}
	requestUser.TenantId = getTenantId(requestUser)
	response, err := json.Marshal(TokenContext{Token:token, TenantId:requestUser.TenantId})
	if err != nil {
		panic(err)
	}
	return response
}

func RefreshToken(requestUser *models.SystemUser) []byte {
	authEngine := InitJWTAuthenticationEngine()
	token, err := authEngine.GenerateToken(requestUser, commons.ServerConfigurations.JWTExpirationDelta)
	if err != nil {
		panic(err)
	}
	requestUser.TenantId = getTenantId(requestUser)
	response, err := json.Marshal(TokenContext{Token:token, TenantId:requestUser.TenantId})
	if err != nil {
		panic(err)
	}
	return response
}

func Logout(req *http.Request) error {
	authEngine := InitJWTAuthenticationEngine()
	return authEngine.InvalidateJWT(req)
}

func RequireTokenAuthentication(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authBackend := InitJWTAuthenticationEngine()
		token, err := authBackend.ProcessToken(r)
		if err != nil || !token.Valid || authBackend.IsInBlacklist(r.Header.Get("Authorization")) {
			log.Debug("Authentication failed " + err.Error())
			w.WriteHeader(http.StatusForbidden)
			return
		} else {
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				sClaims, _ := json.Marshal(claims["roles"])
				r.Header.Set("roles", string(sClaims))
				r.Header.Set("username", claims["sub"].(string))
				r.Header.Set("userid", strconv.FormatFloat((claims["userid"]).(float64), 'f', 0, 64))
				r.Header.Set("tenantid", strconv.FormatFloat((claims["tenantid"]).(float64), 'f', 0, 64))
			}else{
				w.WriteHeader(http.StatusForbidden)
				return
			}
		}
		inner.ServeHTTP(w, r)
	})
}

// Check valid token or not and extract request header
func RequireTokenAuthenticationAndAuthorization(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		// This for superadmin purpose
		// To create tenant users for each tenants
		// To create tenant user groups for each users
		// To create tenant user roles
		// To create user polices
		authBackend := InitJWTAuthenticationEngine() //
		token, err := authBackend.ProcessToken(r)
		if err != nil || !token.Valid {
			w.WriteHeader(http.StatusForbidden) // 403
			return
		} else {
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				sClaims, _ := json.Marshal(claims["roles"])
				userID := claims["userid"].(float64)
				tenantID := claims["tenantid"].(float64)
				r.Header.Set("roles", string(sClaims))
				r.Header.Set("username", claims["sub"].(string))
				r.Header.Set("userid", strconv.FormatFloat(userID, 'f', 0, 64))
				r.Header.Set("tenantid", strconv.FormatFloat(tenantID, 'f', 0, 64))
				a, err := isAuthorized2(int(tenantID), int(userID), r) // check authorization policy for the user
				if err != nil  || !a{
					w.WriteHeader(http.StatusForbidden) // 403
					return
				}
			}else{
				w.WriteHeader(http.StatusUnauthorized) // 401
				return
			}
		}
		inner.ServeHTTP(w, r)
	})
}

func getTenantId(user *models.SystemUser) int64 {
	dbMap, err := commons.GetDBConnection(commons.USER_STORE)
	checkErr(err, "failed to get userstore")
	tenantId, err := dbMap.SelectInt("SELECT tenantid FROM vs_tenants WHERE domain=?", user.TenantDomain)
	checkErr(err, "Select failed")
	return tenantId
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Error(msg, err)
	}
}
