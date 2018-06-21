package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"fmt"

	log "github.com/Sirupsen/logrus"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/vedicsoft/vamps-core/commons"
	"github.com/vedicsoft/vamps-core/models"
)

type TokenAuthentication struct {
	Token    string `json:"token" form:"token"`
	TenantId int64  `json:"tenantid" form:"tenantid"`
}

func Login(requestUser *models.SystemUser) (int, []byte, error) {
	authEngine := InitJWTAuthenticationEngine()
	requestUser.TenantId = getTenantId(requestUser)
	if authEngine.Authenticate(requestUser) {
		token, err := authEngine.GenerateToken(requestUser)
		if err != nil {
			return http.StatusInternalServerError, []byte(""), err
		} else {
			response, _ := json.Marshal(TokenAuthentication{token, requestUser.TenantId})
			return http.StatusOK, response, nil
		}
	}
	return http.StatusUnauthorized, []byte(""), nil
}

func CustomToken(requestUser *models.SystemUser, expirationHours int) []byte {
	authEngine := InitJWTAuthenticationEngine()
	token, err := authEngine.GenerateCustomToken(requestUser, expirationHours)
	if err != nil {
		panic(err)
	}
	requestUser.TenantId = getTenantId(requestUser)
	response, err := json.Marshal(TokenAuthentication{token, requestUser.TenantId})
	if err != nil {
		panic(err)
	}
	return response
}

func RefreshToken(requestUser *models.SystemUser) []byte {
	authEngine := InitJWTAuthenticationEngine()
	token, err := authEngine.GenerateToken(requestUser)
	if err != nil {
		panic(err)
	}
	requestUser.TenantId = getTenantId(requestUser)
	response, err := json.Marshal(TokenAuthentication{token, requestUser.TenantId})
	if err != nil {
		panic(err)
	}
	return response
}

func Logout(req *http.Request) error {
	authEngine := InitJWTAuthenticationEngine()
	tokenRequest, err := jwt.ParseFromRequest(req, func(token *jwt.Token) (interface{}, error) {
		return authEngine.PublicKey, nil
	})
	if err != nil {
		return err
	}
	tokenString := req.Header.Get("Authorization")
	return authEngine.Logout(tokenString, tokenRequest)
}

func RequireTokenAuthentication(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var domainTenant string
		var err error
		if t := r.Header.Get("tenantid"); t != "" {
			domainTenant = t
		}
		fmt.Println(domainTenant)
		authBackend := InitJWTAuthenticationEngine()
		token, err := jwt.ParseFromRequest(
			r,
			func(token *jwt.Token) (interface{}, error) {
				return authBackend.PublicKey, nil
			})
		if err != nil || !token.Valid || authBackend.IsInBlacklist(r.Header.Get("Authorization")) {
			log.Debug("Authentication failed " + err.Error())
			w.WriteHeader(http.StatusForbidden)
			return
		} else {
			sClaims, _ := json.Marshal(token.Claims["roles"])
			r.Header.Set("roles", string(sClaims))
			r.Header.Set("username", token.Claims["sub"].(string))
			r.Header.Set("userid", strconv.FormatFloat((token.Claims["userid"]).(float64), 'f', 0, 64))
			if domainTenant != "" {
				r.Header.Set("tenantid", domainTenant)
			}else {
				r.Header.Set("tenantid", strconv.FormatFloat((token.Claims["tenantid"]).(float64), 'f', 0, 64))
			}

		}
		inner.ServeHTTP(w, r)
	})
}

// Check valid token or not and extract request header
func RequireTokenAuthenticationAndAuthorization(inner http.Handler) http.Handler {
	fmt.Println("test")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var domainTenant string
		var err error
		// This for superadmin purpose
		// To create tenant users for each tenants
		// To create tenant user groups for each users
		// To create tenant user roles
		// To create user polices
		if t := r.Header.Get("tenantid"); t != "" {
			domainTenant = t
		}
		authBackend := InitJWTAuthenticationEngine() //
		token, err := jwt.ParseFromRequest(          // parse token from the request with checking private key and public key
			r,
			func(token *jwt.Token) (interface{}, error) {
				return authBackend.PublicKey, nil
			})
		if err != nil || !token.Valid || authBackend.IsInBlacklist(r.Header.Get("Authorization")) {
			log.Debug("Authentication failed " + err.Error())
			w.WriteHeader(http.StatusForbidden) // 403
			return
		} else {
			sClaims, _ := json.Marshal(token.Claims["roles"])
			userID := token.Claims["userid"]
			tenantID := token.Claims["tenantid"]
			r.Header.Set("roles", string(sClaims))
			r.Header.Set("username", token.Claims["sub"].(string))
			r.Header.Set("userid", strconv.FormatFloat(userID.(float64), 'f', 0, 64))
			if domainTenant != "" {
				r.Header.Set("tenantid", domainTenant)
			}else {
				r.Header.Set("tenantid", strconv.FormatFloat((token.Claims["tenantid"]).(float64), 'f', 0, 64))
			}

			a, err := isAuthorized2(int(tenantID.(float64)), int(userID.(float64)), r) // check authorization policy for the user
			if err != nil {
				log.Debug("authorization failed due to error " + err.Error())
				w.WriteHeader(http.StatusUnauthorized) // 401
				return
			}
			if !a {
				log.Debug("authorization failed for user " + strconv.Itoa(int(userID.(float64))))
				w.WriteHeader(http.StatusUnauthorized) //401
				return
			}
		}
		inner.ServeHTTP(w, r)
	})
}

func getTenantId(user *models.SystemUser) int64 {
	dbMap := commons.GetDBConnection(commons.USER_STORE_DB)
	tenantId, err := dbMap.SelectInt("SELECT tenantid FROM vs_tenants WHERE domain=?", user.TenantDomain)
	checkErr(err, "Select failed")
	return tenantId
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Error(msg, err)
	}
}
