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
	"github.com/astaxie/beego"
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
	if strings.ToLower(beego.BConfig.RunMode) == "dev" {
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

	//Starting the API server
	router := routes.NewRouter()
	httpsServer := &http.Server{
		Addr: ":" + strconv.Itoa(commons.ServerConfigurations.HttpsPort+
			commons.ServerConfigurations.PortOffset),
		Handler:        router,
		ReadTimeout:    time.Duration(commons.ServerConfigurations.ReadTimeOut) * time.Second,
		WriteTimeout:   time.Duration(commons.ServerConfigurations.WriteTimeOut) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Info("Starting server on port : " + strconv.Itoa(commons.ServerConfigurations.HttpsPort+
		commons.ServerConfigurations.PortOffset))
	log.Fatal("HTTP Server error: ", httpsServer.ListenAndServeTLS(commons.ServerConfigurations.SSLCertificateFile,
		commons.ServerConfigurations.SSLKeyFile))
}
