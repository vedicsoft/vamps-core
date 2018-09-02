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
	"fmt"
	"strings"
	"net/http"
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

func (backend *JWTAuthenticationBackend) GenerateToken(user *models.SystemUser, expInHours int) (string, error) {
	exp := time.Now().Add(time.Hour * time.Duration(expInHours)).Unix()
	roles, err := getUserSystemRoles(user)
	if err != nil {
		return "", errors.New("could not load user scopes stack trace: " + err.Error())
	}
	groups, err := getUserGroups(user)
	if err != nil {
		return "", errors.New("could not load user scopes stack trace: " + err.Error())
	}
	t := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iat":    time.Now().Unix(),
		"exp":    exp,
		"sub":    user.Username,
		"tenantid": user.TenantId,
		"userid":     getUserId(user),
		"roles":     roles,
		"groups":    groups,
	})
	tokenString, err := t.SignedString(backend.privateKey)
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

func getUserSystemRoles(user *models.SystemUser) ([]string, error) {
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

func getUserGroups(user *models.SystemUser) ([]string, error) {
	const GET_USER_GROUPS string = `SELECT vs_groups.name  from vs_groups WHERE vs_groups.id IN (SELECT vs_group_users.groupid FROM vs_group_users WHERE vs_group_users.userid= ?)`
	var groups []string
	dbMap := commons.GetDBConnection(commons.PLATFORM_DB)
	var err error
	_, err = dbMap.Select(&groups, GET_USER_GROUPS, user.UserId)
	if err != nil {
		return groups, err
	}
	return groups, err
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

func (backend *JWTAuthenticationBackend) InvalidateJWT(r *http.Request) error {
	t, err := backend.ProcessToken(r)
	if err != nil || !t.Valid {
		return errors.New("invalid token")
	}
	// get remaining token validity
	if claims, ok := t.Claims.(jwt.MapClaims); ok {
		if validity, ok := claims["exp"].(float64); ok {
			tm := time.Unix(int64(validity), 0)
			remainer := tm.Sub(time.Now())
			if remainer > 0 {
				return redis.SetValue(t.Raw, t.Raw, remainer.Seconds() + expireOffset)
			} else {
				return nil
			}
		}
	}
	return nil
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

func (backend *JWTAuthenticationBackend) ProcessToken(r *http.Request) (*jwt.Token, error) {
	// 1. extract token from request
	tokenString, err := extractToken(r)
	if err != nil {
		return nil, err
	}
	// 2. verify
	if backend.privateKey != nil {
		t, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return backend.privateKey, nil
		})
		if err != nil || !t.Valid {
			return nil, err
		}
		// 3. Check the black list
		if !backend.IsInBlacklist(tokenString) {
			t.Valid = false
			return nil, errors.New("black listed token")
		}
	}
	return nil, errors.New("public key is not initialized")
}

func extractToken(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")
	ts, err := getAuthToken("Bearer", header)
	if err != nil {
		// unable to find the token from the header
		tc, err := r.Cookie("access_token")
		if err == nil {
			return ts, errors.New("access token is not present in the request")
		} else {
			return tc.Raw, nil
		}
	} else {
		return ts, err
	}
}

//returns the token string from the header value
//Input header format should be "<tokenType> <token>"
func getAuthToken(tokenType, header string) (string, error) {
	tmp := strings.Split(header, " ")
	if len(tmp) == 2 && tmp[0] == tokenType {
		//jwt token authentication
		return tmp[1], nil
	} else {
		return "", errors.New("unable to extract the token")
	}
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