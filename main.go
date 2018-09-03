package main

import (
	"net/http"
	"os"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/vedicsoft/vamps-core/commons"
	"github.com/vedicsoft/vamps-core/routes"
	"github.com/astaxie/beego/config"
	_ "github.com/astaxie/beego/config/xml"
	"strings"
	"fmt"
	"errors"
)


func initStore(storeType string, conf config.Configer) error {
	s := commons.Store{
		Type:     storeType,
		Dialect:  conf.String(storeType + "::dialect"),
		Host:     conf.String(storeType + "::url"),
		Port:     conf.String(storeType + "::port"),
		Username: conf.String(storeType + "::username"),
		Password: conf.String(storeType + "::password"),
		DBName:   conf.String(storeType + "::db"),
	}
	if strings.ToLower(conf.DefaultString("runMode", "dev")) == "dev" {
		s.ShouldMigrate = true
	}
	err := s.RegisterDB()
	if err != nil {
		msg := fmt.Sprintf("store: %s %s:%s@tcp(%s:%s)/%s?charset=utf8&multiStatements=true err: %s",
			storeType, s.Username, "***", s.Host, s.Port, s.DBName, err.Error())
		return errors.New(msg)
	}
	return err
}

func main() {
	conf, err := config.NewConfig("ini", "conf/app.conf")
	if err != nil {
		log.Fatalf("failed to parse config file err: %s", err.Error())
	}
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	switch conf.String("logLevel") {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}

	// DB initialization
	err = initStore("userstore", conf)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatal("failed to initialize userstore")
	}

	err = initStore("datastore", conf)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatal("failed to initialize datastore")
	}
	router := routes.NewRouter()
	httpsServer := &http.Server{
		Addr: ":" + strconv.Itoa(conf.DefaultInt("httpsPort", 443)),
		Handler:        router,
		ReadTimeout:    time.Duration(conf.DefaultInt("readTimeout", 20)) * time.Second,
		WriteTimeout:   time.Duration(conf.DefaultInt("writeTimeout", 20)) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Info("Starting server on port  " + httpsServer.Addr)
	log.Fatal("HTTP Server error: ", httpsServer.ListenAndServeTLS(
		conf.DefaultString("certFile", "resources/security/server.pem"),
		conf.DefaultString("keyFile", "resources/security/server.key")))
}
