package controllers

import (
	"crypto/rsa"
	"database/sql"
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
	"io/ioutil"
)

type JWTBackend struct {
	privateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

const (
	expireOffset = 3600
)

var (
	authEngine *JWTBackend = nil
)

func InitJWTAuthenticationEngine() (*JWTBackend, error) {
	if authEngine == nil {
		privateKey, err := ioutil.ReadFile(os.Getenv("JWT_PRIVATE_KEY_FILE"))
		if err != nil {
			return nil, err
		}
		parsedPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
		if err != nil {
			return nil, err
		}

		publicKey, err := ioutil.ReadFile(os.Getenv("JWT_PUBLIC_KEY_FILE"))
		if err != nil {
			return nil, err
		}
		parsedPublicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
		if err != nil {
			return nil, err
		}
		authEngine = &JWTBackend{
			privateKey: parsedPrivateKey,
			PublicKey:  parsedPublicKey,
		}
		return authEngine, err
	}
	return authEngine, nil
}

func (backend *JWTBackend) GenerateToken(user *models.SystemUser, expInHours int) (string, error) {
	exp := time.Now().Add(time.Hour * time.Duration(expInHours)).Unix()
	roles, err := getUserSystemRoles(user)
	if err != nil {
		return "", errors.New("could not load user roles err: " + err.Error())
	}
	groups, err := getUserGroups(user)
	if err != nil {
		return "", errors.New("could not load usergroups err: " + err.Error())
	}
	uID, err := getUserId(user)
	if err != nil {
		return "", errors.New("could not load userId err: " + err.Error())
	}
	t := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iat":    time.Now().Unix(),
		"exp":    exp,
		"sub":    user.Username,
		"tenantid": user.TenantId,
		"userid":   uID,
		"roles":     roles,
		"groups":    groups,
	})
	tokenString, err := t.SignedString(backend.privateKey)
	if err != nil {
		return "", errors.New("unable to sign the jwt stack trace: " + err.Error())
	}
	return tokenString, nil
}

func getUserId(user *models.SystemUser) (int64, error) {
	dbMap, err := commons.GetDBConnection(commons.USER_STORE)
	if err != nil {
		return 0, err
	}
	var userId sql.NullInt64
	smtOut, err := dbMap.Db.Prepare("SELECT userid FROM vs_users WHERE username=? AND tenantid=?")
	if err != nil {
		return 0, err
	}
	defer smtOut.Close()
	err = smtOut.QueryRow(user.Username, user.TenantId).Scan(&userId)
	if err != nil {
		return 0, err
	} else {
		user.UserId = userId.Int64
		return userId.Int64, err
	}
}

func getUserSystemRoles(user *models.SystemUser) ([]string, error) {
	const GET_USER_ROLES string = `SELECT vs_roles.name from vs_roles WHERE vs_roles.roleid IN (SELECT
								   vs_user_roles.roleid FROM vs_user_roles WHERE
								   vs_user_roles.userid=?) AND vs_roles.type='system'`
	var roles []string
	dbMap, err := commons.GetDBConnection(commons.USER_STORE)
	if err != nil {
		return roles, err
	}
	_, err = dbMap.Select(&roles, GET_USER_ROLES, user.UserId)
	if err != nil {
		return roles, err
	}
	return roles, err
}

func getUserGroups(user *models.SystemUser) ([]string, error) {
	const GET_USER_GROUPS string = `SELECT vs_groups.name  from vs_groups WHERE vs_groups.groupid IN
	(SELECT vs_group_users.groupid FROM vs_group_users WHERE vs_group_users.userid= ?)`
	var groups []string
	dbMap, err := commons.GetDBConnection(commons.USER_STORE)
	if err != nil {
		return groups, err
	}
	_, err = dbMap.Select(&groups, GET_USER_GROUPS, user.UserId)
	if err != nil {
		return groups, err
	}
	return groups, err
}

func (backend *JWTBackend) Authenticate(user *models.SystemUser) (bool, error) {
	dbMap, err := commons.GetDBConnection(commons.USER_STORE)
	if err != nil {
		return false, err
	}
	var hashedPassword sql.NullString
	smtOut, err := dbMap.Db.Prepare("SELECT password FROM vs_users where username=? AND tenantid=? AND status='active'")
	if err != nil {
		return false, err
	}
	defer smtOut.Close()
	err = smtOut.QueryRow(user.Username, user.TenantId).Scan(&hashedPassword)
	if err == nil && hashedPassword.Valid {
		if len(hashedPassword.String) > 0 {
			err = bcrypt.CompareHashAndPassword([]byte(hashedPassword.String), []byte(user.Password))
			if err == nil {
				return true, nil
			}
		}
	} else {
		return false, err
	}
	return false, err
}

func (backend *JWTBackend) getTokenRemainingValidity(timestamp interface{}) int {
	if validity, ok := timestamp.(float64); ok {
		tm := time.Unix(int64(validity), 0)
		remainer := tm.Sub(time.Now())
		if remainer > 0 {
			return int(remainer.Seconds() + expireOffset)
		}
	}
	return expireOffset
}

func (backend *JWTBackend) InvalidateJWT(r *http.Request) error {
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

func (backend *JWTBackend) IsInBlacklist(token string) bool {
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

func (backend *JWTBackend) ProcessToken(r *http.Request) (*jwt.Token, error) {
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