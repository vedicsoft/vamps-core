package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/handlers"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"strconv"
	"time"
	"github.com/vamps-core/routes"
	"github.com/vamps-core/commons"
)

var ServerHome string
var serverConfigs ServerConfigs
var logHandler http.Handler
var serverLogFile os.File
var httpAccessLogFile os.File

func init() {
	ServerHome = os.Getenv(commons.SERVER_HOME)
	if ( len(ServerHome) <= 0 ) {
		ServerHome = os.Args[1]
	}

	viper.New()
	viper.AddConfigPath(ServerHome + "/configs")
	viper.SetConfigName("config")
	if _, err := os.Stat("../configs/config.yaml"); os.IsNotExist(err) {
		viper.SetConfigName("config.default")
	}
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		// Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	configsMap := viper.GetStringMap("serverConfigs")
	serverConfigs.HttpPort = configsMap["httpPort"].(int)
	serverConfigs.ReadTimeOut = configsMap["readTimeOut"].(int)
	serverConfigs.WriteTimeOut = configsMap["writeTimeOut"].(int)
	serverConfigs.LogsDirectory = configsMap["logsdirectory"].(string)

	serverLogFile, err := os.OpenFile(serverConfigs.LogsDirectory + "/" + SERVER_LOG_FILE_NAME, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error while opening server log file: %v", err)
	}

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(serverLogFile)

	httpAccessLogFile, err := os.OpenFile(serverConfigs.LogsDirectory + "/http-access.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	enableHttpAccessLogs := viper.GetBool("httpAccessLogs")

	if enableHttpAccessLogs {
		logHandler = handlers.LoggingHandler(httpAccessLogFile, http.DefaultServeMux)
	}
	//
	//utils.Init_(serverHome)

}

func main() {
	initConfigurations(ServerHome)
	InitConfigs(ServerHome)
	commons.ServerHome = ServerHome
	defer serverLogFile.Close()
	defer httpAccessLogFile.Close()

	router := routes.NewRouter()
	router.PathPrefix("/dashboard/").Handler(http.StripPrefix("/dashboard/", http.FileServer(http.Dir(ServerHome + "/webapps/dashboard/"))))
	http.Handle("/", router)

	httpsServer := &http.Server{
		Addr:           ":" + strconv.Itoa(serverConfigs.HttpPort),
		Handler:        logHandler,
		ReadTimeout:    time.Duration(serverConfigs.ReadTimeOut) * time.Second,
		WriteTimeout:   time.Duration(serverConfigs.WriteTimeOut) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Info("Starting server on port : " + strconv.Itoa(serverConfigs.HttpPort))
	log.Fatal("HTTP Server error: ", httpsServer.ListenAndServeTLS(ServerHome + "/resources/security/server.pem", ServerHome + "/resources/security/server.key"))
}

func initConfigurations(serverHome string) {

}
