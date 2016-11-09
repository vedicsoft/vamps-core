package main

import (
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/vedicsoft/vamps-core/commons"
	"github.com/vedicsoft/vamps-core/routes"
)

var logHandler http.Handler

func main() {
	os.Chdir(commons.ServerConfigurations.Home)
	serverLogFile, err := os.OpenFile(commons.ServerConfigurations.LogsDirectory+"/"+commons.SERVER_LOG_FILE_NAME,
		os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error while opening server log file: %v", err)
	}

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(serverLogFile)

	switch commons.ServerConfigurations.LogLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}

	defer serverLogFile.Close()

	commons.ConstructConnectionPool(commons.ServerConfigurations.DBConfigMap)

	// Starting caddy server to server static files
	args := []string{"bin/caddy", "--conf=" + commons.ServerConfigurations.CaddyFile, "-pidfile=bin/caddy.pid"}

	if err := exec.Command("nohup", args...).Start(); err != nil {
		log.Fatalln("Error occourred while starting caddy server : ", err.Error())
		os.Exit(1)
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
