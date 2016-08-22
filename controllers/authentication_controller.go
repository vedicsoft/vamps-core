package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	log "github.com/Sirupsen/logrus"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/vedicsoft/vamps-core/commons"
	"github.com/vedicsoft/vamps-core/models"
)

type TokenAuthentication struct {
	Token    string `json:"token" form:"token"`
	TenantId int64  `json:"tenantid" form:"tenantid"`
}

func Login(requestUser *models.SystemUser) (int, []byte) {
	authEngine := InitJWTAuthenticationEngine()
	requestUser.TenantId = getTenantId(requestUser)
	if authEngine.Authenticate(requestUser) {
		token, err := authEngine.GenerateToken(requestUser)
		if err != nil {
			return http.StatusInternalServerError, []byte("")
		} else {
			response, _ := json.Marshal(TokenAuthentication{token, requestUser.TenantId})
			return http.StatusOK, response
		}
	}
	return http.StatusUnauthorized, []byte("")
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
		authBackend := InitJWTAuthenticationEngine()
		token, err := jwt.ParseFromRequest(
			r,
			func(token *jwt.Token) (interface{}, error) {
				return authBackend.PublicKey, nil
			})
		if err != nil || !token.Valid || authBackend.IsInBlacklist(r.Header.Get("Authorization")) {
			w.WriteHeader(http.StatusForbidden)
			return
		} else {
			sClaims, _ := json.Marshal(token.Claims["scopes"])
			r.Header.Set("scopes", string(sClaims))
			r.Header.Set("username", token.Claims["sub"].(string))
			r.Header.Set("tenantid", strconv.FormatFloat((token.Claims["tenantid"]).(float64), 'f', 0, 64))
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
		log.Fatalln(msg, err)
	}
}
