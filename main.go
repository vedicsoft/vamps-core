package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/handlers"
	"github.com/vedicsoft/vamps-core/commons"
	"github.com/vedicsoft/vamps-core/routes"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"
)

var logHandler http.Handler

func main() {
	configFile := flag.String("flagname", "", "serverconfiguration file")
	commons.InitConfigurations(*configFile)
	os.Chdir(commons.ServerConfigurations.Home)
	serverLogFile, err := os.OpenFile(commons.ServerConfigurations.LogsDirectory+"/"+commons.SERVER_LOG_FILE_NAME, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error while opening server log file: %v", err)
	}

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(serverLogFile)
	log.SetLevel(log.DebugLevel)

	httpAccessLogFile, err := os.OpenFile(commons.ServerConfigurations.LogsDirectory+"/"+commons.ACCESS_LOG_FILE_NAME, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error while trying to open the access log file: %v", err)
	}

	if commons.ServerConfigurations.EnableAccessLogs {
		logHandler = handlers.LoggingHandler(httpAccessLogFile, http.DefaultServeMux)
	}

	defer serverLogFile.Close()
	defer httpAccessLogFile.Close()

	commons.ConstructConnectionPool(commons.ServerConfigurations.DBConfigMap)

	args := []string{"bin/caddy", "--conf=" + commons.ServerConfigurations.CaddyFile, "-pidfile=bin/caddy.pid"}

	if err := exec.Command("nohup", args...).Start(); err != nil {
		log.Fatalln("Error occourred while starting caddy server : ", err.Error())
		os.Exit(1)
	}

	//Starting the API server
	router := routes.NewRouter()
	http.Handle("/", router)

	httpsServer := &http.Server{
		Addr:           ":" + strconv.Itoa(commons.ServerConfigurations.HttpsPort+commons.ServerConfigurations.PortOffset),
		Handler:        logHandler,
		ReadTimeout:    time.Duration(commons.ServerConfigurations.ReadTimeOut) * time.Second,
		WriteTimeout:   time.Duration(commons.ServerConfigurations.WriteTimeOut) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Info("Starting server on port : " + strconv.Itoa(commons.ServerConfigurations.HttpsPort+commons.ServerConfigurations.PortOffset))
	log.Fatal("HTTP Server error: ", httpsServer.ListenAndServeTLS(commons.ServerConfigurations.SSLCertificateFile, commons.ServerConfigurations.SSLKeyFile))
}
