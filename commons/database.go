package commons

import (
	"database/sql"
	"encoding/json"

	log "github.com/Sirupsen/logrus"
	//_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/gorp.v1"

	_ "github.com/mattes/migrate/driver/mysql"
	"github.com/mattes/migrate/migrate"
	"gopkg.in/mgo.v2"
)

const DIALECT_MYSQL string = "mysql"
const DIALECT_SQLITE3 string = "sqlite3"
const DIALECT_MONGO string = "mongodb"

type DBConnection struct {
	connectionURL string
	dbMap         *gorp.DbMap
}

var dbConnections map[string]DBConnection
var mongoConnectionUrl string
var mgoSession *mgo.Session

func GetDBConnection(dbIdentifier string) *gorp.DbMap {
	return dbConnections[dbIdentifier].dbMap
}

func GetDBDetailedConnection(dbName string) DBConnection {
	return dbConnections[dbName]
}

func GetMongoSession() *mgo.Session {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial(mongoConnectionUrl)
		if err != nil {
			log.Fatal("Failed to start the Mongo session")
		}
	}
	return mgoSession.Clone()
}

func ConstructConnectionPool(dbConfigs map[string]DBConfigs) {
	dbConnections = make(map[string]DBConnection)
	var connectionURL string
	var dialect gorp.Dialect
	for dbName, dbConfig := range dbConfigs {
		switch dbConfig.Dialect {
		case DIALECT_MYSQL:
			connectionURL = dbConfig.Username + ":" + dbConfig.Password + "@tcp(" + dbConfig.Address + ")/" +
				dbConfig.DBName + dbConfig.Parameters
			dialect = gorp.MySQLDialect{"InnoDB", "UTF8"}
			break
		case DIALECT_SQLITE3:
			connectionURL = dbConfig.Address
			dialect = gorp.SqliteDialect{}
			break
		case DIALECT_MONGO:
			//mongoConnectionUrl = "mongodb://"+ dbConfig.Username+":"+ dbConfig.Password+"@"+dbConfig.Address
			mongoConnectionUrl = "mongodb://" + dbConfig.Address
			continue
		}
		db, err := sql.Open(dbConfig.Dialect, connectionURL)
		if err != nil {
			log.Error("Error occurred while constructing a the DB connection to : " + connectionURL +
				" with dialect:" + dbConfig.Dialect + " stack:" + err.Error())
		}
		dbConnections[dbName] = DBConnection{connectionURL, &gorp.DbMap{Db: db, Dialect: dialect}}
	}
}

func UpgradeDBSchema(dbIdentifier string) {
	connection := GetDBDetailedConnection(dbIdentifier)
	allErrors, ok := migrate.UpSync("mysql://"+connection.connectionURL, GetServerHome()+"/resources/sql/"+
		dbIdentifier)
	if !ok {
		log.Fatal("Error occurred while migrating the database :"+dbIdentifier, allErrors)
	}
}

func DowngradeDBSchema(dbIdentifier string) {
	connection := GetDBDetailedConnection(dbIdentifier)
	allErrors, ok := migrate.DownSync("mysql://"+connection.connectionURL, GetServerHome()+"/resources/sql/"+
		dbIdentifier)
	if !ok {
		log.Fatal("Error occurred while migrating the database :"+dbIdentifier, allErrors)
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
