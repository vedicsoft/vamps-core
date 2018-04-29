package controllers

import (
	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"os"
	"time"

	"errors"

	log "github.com/Sirupsen/logrus"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/vedicsoft/vamps-core/commons"
	"github.com/vedicsoft/vamps-core/models"
	"github.com/vedicsoft/vamps-core/redis"
	"golang.org/x/crypto/bcrypt"
)

type JWTAuthenticationBackend struct {
	privateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

const (
	expireOffset = 3600
)

var authBackendInstance *JWTAuthenticationBackend = nil

func InitJWTAuthenticationEngine() *JWTAuthenticationBackend {
	if authBackendInstance == nil {
		authBackendInstance = &JWTAuthenticationBackend{
			privateKey: getPrivateKey(),
			PublicKey:  getPublicKey(),
		}
	}
	return authBackendInstance
}

func (backend *JWTAuthenticationBackend) GenerateToken(user *models.SystemUser) (string, error) {
	token := jwt.New(jwt.SigningMethodRS512)
	i := commons.ServerConfigurations.JWTExpirationDelta
	token.Claims["exp"] = time.Now().Add(time.Hour * time.Duration(i)).Unix()
	token.Claims["iat"] = time.Now().Unix()
	token.Claims["sub"] = user.Username
	token.Claims["tenantid"] = user.TenantId
	token.Claims["userid"] = getUserId(user)
	scopes, err := getUserScopes(user)
	if err != nil {
		return "", errors.New("could not load user scopes stack trace: " + err.Error())
	}
	token.Claims["scopes"] = scopes
	types := GetTypesOfPoliciesByUserID(user.UserId)
	token.Claims["types"] = types
	tokenString, err := token.SignedString(backend.privateKey)
	if err != nil {
		return "", errors.New("unable to sign the jwt stack trace: " + err.Error())
	}
	return tokenString, nil
}

func (backend *JWTAuthenticationBackend) GenerateCustomToken(user *models.SystemUser, expirationHours int) (string, error) {
	token := jwt.New(jwt.SigningMethodRS512)
	token.Claims["exp"] = time.Now().Add(time.Hour * time.Duration(expirationHours)).Unix()
	token.Claims["iat"] = time.Now().Unix()
	token.Claims["sub"] = user.Username
	token.Claims["tenantid"] = user.TenantId
	token.Claims["userid"] = getUserId(user)
	scopes, err := getUserScopes(user)
	if err != nil {
		return "", errors.New("could not load user scopes stack trace: " + err.Error())
	}
	token.Claims["scopes"] = scopes
	types := GetTypesOfPoliciesByUserID(user.UserId)
	token.Claims["types"] = types
	tokenString, err := token.SignedString(backend.privateKey)
	if err != nil {
		return "", errors.New("unable to sign the jwt stack trace: " + err.Error())
	}
	return tokenString, nil
}

func getUserId(user *models.SystemUser) int64 {
	dbMap := commons.GetDBConnection(commons.PLATFORM_DB)
	var userId sql.NullInt64
	smtOut, err := dbMap.Db.Prepare("SELECT userid FROM vs_users WHERE username=? AND tenantid=?")
	defer smtOut.Close()
	err = smtOut.QueryRow(user.Username, user.TenantId).Scan(&userId)
	if err != nil {
		log.Debug("User authentication failed " + user.Username)
		return -1
	} else {
		user.UserId = userId.Int64
		return userId.Int64
	}
}

func getUserScopes(user *models.SystemUser) ([]string, error) {
	const GET_USER_ROLES string = `SELECT vs_roles.name from vs_roles WHERE vs_roles.roleid IN (SELECT
								   vs_user_roles.roleid FROM vs_user_roles WHERE
								   vs_user_roles.userid=?) AND vs_roles.type='system'`
	var roles []string
	dbMap := commons.GetDBConnection(commons.PLATFORM_DB)
	var err error
	_, err = dbMap.Select(&roles, GET_USER_ROLES, user.UserId)
	if err != nil {
		return roles, err
	}
	return roles, err
}

func (backend *JWTAuthenticationBackend) Authenticate(user *models.SystemUser) bool {
	dbMap := commons.GetDBConnection(commons.PLATFORM_DB)
	var hashedPassword sql.NullString
	smtOut, err := dbMap.Db.Prepare("SELECT password FROM vs_users where username=? AND tenantid=? AND status='active'")
	defer smtOut.Close()

	err = smtOut.QueryRow(user.Username, user.TenantId).Scan(&hashedPassword)
	if err == nil && hashedPassword.Valid {
		if len(hashedPassword.String) > 0 {
			err = bcrypt.CompareHashAndPassword([]byte(hashedPassword.String), []byte(user.Password))
			if err == nil {
				log.Debug("User authenticated successfully " + user.Username)
				return true
			}
		}
	} else {
		log.Debug("User authentication failed for user " + user.Username)
		return false
	}
	return false
}

func (backend *JWTAuthenticationBackend) getTokenRemainingValidity(timestamp interface{}) int {
	if validity, ok := timestamp.(float64); ok {
		tm := time.Unix(int64(validity), 0)
		remainer := tm.Sub(time.Now())
		if remainer > 0 {
			return int(remainer.Seconds() + expireOffset)
		}
	}
	return expireOffset
}

func (backend *JWTAuthenticationBackend) Logout(tokenString string, token *jwt.Token) error {
	return redis.SetValue(tokenString, tokenString, backend.getTokenRemainingValidity(token.Claims["exp"]))
}

func (backend *JWTAuthenticationBackend) IsInBlacklist(token string) bool {
	redisToken, err := redis.GetValue(token)
	if err != nil {
		log.Error("Error occourred while checking for black listed jwt :" + token + " stack :" + err.Error())
	}
	if redisToken == nil {
		log.Debug("Token is not in the black list")
		return false
	}
	log.Debug("Found a blacklisted token")
	return true
}

func getPrivateKey() *rsa.PrivateKey {
	privateKeyFile, err := os.Open(commons.ServerConfigurations.JWTPrivateKeyFile)
	if err != nil {
		panic(err)
	}

	pemfileinfo, _ := privateKeyFile.Stat()
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(privateKeyFile)
	_, err = buffer.Read(pembytes)

	data, _ := pem.Decode([]byte(pembytes))

	defer privateKeyFile.Close()

	privateKeyImported, err := x509.ParsePKCS1PrivateKey(data.Bytes)

	if err != nil {
		panic(err)
	}
	return privateKeyImported
}

func getPublicKey() *rsa.PublicKey {
	publicKeyFile, err := os.Open(commons.ServerConfigurations.JWTPublicKeyFile)
	if err != nil {
		panic(err)
	}

	pemfileinfo, _ := publicKeyFile.Stat()
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(publicKeyFile)
	_, err = buffer.Read(pembytes)

	data, _ := pem.Decode([]byte(pembytes))

	defer publicKeyFile.Close()

	publicKeyImported, err := x509.ParsePKIXPublicKey(data.Bytes)

	if err != nil {
		panic(err)
	}

	rsaPub, ok := publicKeyImported.(*rsa.PublicKey)

	if !ok {
		panic(err)
	}

	return rsaPub
}

const GET_ROLES_BY_USERID = `SELECT roleid FROM vs_user_roles where userid = ?`

func GetTypesOfPoliciesByUserID(userid int64) []string {
	println("GetTypesOfPoliciesByUserID")
	var roles []int
	dbMap := commons.GetDBConnection(commons.PLATFORM_DB)
	_, err := dbMap.Select(&roles, GET_ROLES_BY_USERID, userid)
	if err != nil {
		println(err.Error())
	}
	return GetPoliciesByRoleID(roles)
}

const GET_POLICIES_BY_ROLEID = `SELECT policyid FROM vs_role_policies where roleid = ?`

func GetPoliciesByRoleID(roleids []int) []string {
	println("GetPoliciesByRoleID")
	var policies []int
	dbMap := commons.GetDBConnection(commons.PLATFORM_DB)
	for _, roleid := range roleids {
		_, err := dbMap.Select(&policies, GET_POLICIES_BY_ROLEID, roleid)
		if err != nil {
			println(err.Error())
		}
	}
	return GetTypesOfPolicies(unique(policies))
}

func unique(intSlice []int) []int {
	keys := make(map[int]bool)
	list := []int{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

const GET_TYPES_OF_POLICIES = `SELECT type FROM vamps.vs_policies where policyid = ?`

func GetTypesOfPolicies(policyIDs []int) []string {
	println("GetTypesOfPolicies")
	var types []string
	dbMap := commons.GetDBConnection(commons.PLATFORM_DB)
	for _, policyID := range policyIDs {
		_, err := dbMap.Select(&types, GET_TYPES_OF_POLICIES, policyID)
		if err != nil {
			println(err.Error())
		}
	}
	return types
}
