package commons

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gorp.v1"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
)

//var dbConfigs map[string]DBConfigs
var dbConnections map[string]*gorp.DbMap

func GetDBConnection(dbname string) *gorp.DbMap {
	return dbConnections[dbname]
}

func ConstructConnectionPool(dbConfigs map[string]DBConfigs) {
	dbConnections = make(map[string]*gorp.DbMap)
	for dbname, dbconfig := range dbConfigs {
		connectionUrl := dbconfig.Username + ":" + dbconfig.Password + "@tcp(" + dbconfig.Address + ")/" + dbconfig.DBName + dbconfig.Parameters
		db, err := sql.Open("mysql", connectionUrl)
		if err != nil {
			log.Error("Error occourred while constructing a the DB connection to : " + connectionUrl)
		}
		dbConnections[dbname] = &gorp.DbMap{Db: db, Dialect:gorp.MySQLDialect{"InnoDB", "UTF8"}}
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