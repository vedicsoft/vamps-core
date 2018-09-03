package commons

import (
	"database/sql"
	"encoding/json"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/gorp.v1"
	"gopkg.in/mgo.v2"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/mysql"
	_ "github.com/mattes/migrate/source/file"
	"errors"
)

const (
	DIALECT_MYSQL string = "mysql"
	DIALECT_POSTGRES string = "postgres"
	DIALECT_SQLITE3 string = "sqlite3"
	DIALECT_MONGO string = "mongodb"

	USER_STORE string = "userstore"
	DATA_STORE string = "datastore"
	ANALYTICS_STORE string = "analyticsstore"
)

var (
	mongoConnectionUrl string
	mgoSession *mgo.Session
	createdStores = make(map[string]*gorp.DbMap)
)

type Store struct {
	Type          string
	Dialect       string
	Host          string
	Port          string
	Username      string
	Password      string
	DBName        string
	ShouldMigrate bool
}

func (s *Store) RegisterDB() error {
	switch s.Dialect {
	case DIALECT_MYSQL:
		if s.ShouldMigrate {
			err := s.Migrate()
			if err != nil {
				return err
			}
		}
		connURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&" +
		"multiStatements=true", s.Username, s.Password, s.Host, s.Port, s.DBName)
		db, err := sql.Open(s.Dialect, connURL)
		if err != nil {
			return err
		}
		createdStores[s.Type] = &gorp.DbMap{Db:db, Dialect: s.Dialect}
	case DIALECT_POSTGRES:
	//TO DO
	}
	return nil
}

func (s *Store) Migrate() error {
	switch s.Dialect {
	case DIALECT_MYSQL:
		cURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&" +
		"multiStatements=true", s.Username, s.Password, s.Host, s.Port, s.DBName)
		db, err := sql.Open("mysql", cURL)
		if err != nil {
			return err
		}
		defer db.Close()
		driver, err := mysql.WithInstance(db, &mysql.Config{
			MigrationsTable: s.Type + "_migrations",
		})
		if err != nil {
			return errors.New(fmt.Sprintf("failed to initialize DB migration err: %s", err.Error()))

		}
		defer driver.Close()
		m, err := migrate.NewWithDatabaseInstance(
			"file://resources/db_scripts/" + s.Dialect + "/" + s.Type, s.Dialect, driver)
		if err != nil {
			return err
		}
		defer m.Close()
		err = m.Up()
		if err != nil && migrate.ErrNoChange != err {
			return err
		}
	case DIALECT_POSTGRES:
	// TODO: implement
	}
	return nil
}

func GetDBConnection(storeType string) (*gorp.DbMap, error) {
	if store, ok := createdStores[storeType]; ok {
		return store.Db, nil
	} else {
		return nil, errors.New(fmt.Sprintf("store not initialized"))
	}
}

func GetMongoSession() (*mgo.Session, error) {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial(mongoConnectionUrl)
		if err != nil {
			return nil, err
		}
	}
	return mgoSession.Clone(), nil
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
