package commons

import (
	"database/sql"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/gorp.v1"
)

const DIALECT_MYSQL string = "mysql"
const DIALECT_SQLITE3 string = "sqlite3"

var dbConnections map[string]*gorp.DbMap

func GetDBConnection(dbName string) *gorp.DbMap {
	return dbConnections[dbName]
}

func ConstructConnectionPool(dbConfigs map[string]DBConfigs) {
	dbConnections = make(map[string]*gorp.DbMap)
	var connectionUrl string
	var dialect gorp.Dialect
	for dbName, dbConfig := range dbConfigs {
		switch dbConfig.Dialect {
		case DIALECT_MYSQL:
			connectionUrl = dbConfig.Username + ":" + dbConfig.Password + "@tcp(" + dbConfig.Address + ")/" + dbConfig.DBName + dbConfig.Parameters
			dialect = gorp.MySQLDialect{"InnoDB", "UTF8"}
			break
		case DIALECT_SQLITE3:
			connectionUrl = dbConfig.Address
			dialect = gorp.SqliteDialect{}
			break
		}
		db, err := sql.Open(dbConfig.Dialect, connectionUrl)
		if err != nil {
			log.Error("Error occourred while constructing a the DB connection to : " + connectionUrl + " with dialect:" + dbConfig.Dialect + " stack:" + err.Error())
		}
		dbConnections[dbName] = &gorp.DbMap{Db: db, Dialect: dialect}
	}
}

type NullString struct {
	sql.NullString
}

func (r NullString) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String)
}

func (r *NullString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	r.Valid = true
	return json.Unmarshal(data, (*string)(&r.String))
}

/* NullInt*/
type NullInt64 struct {
	sql.NullInt64
}

func (r NullInt64) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Int64)
}

func (r *NullInt64) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	r.Valid = true
	return json.Unmarshal(data, (*int64)(&r.Int64))
}

/* NullFloat64*/
type NullFloat64 struct {
	sql.NullFloat64
}

func (r NullFloat64) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Float64)
}

func (r *NullFloat64) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	r.Valid = true
	return json.Unmarshal(data, (*float64)(&r.Float64))
}
