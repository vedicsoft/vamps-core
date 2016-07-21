package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/handlers"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"strconv"
	"time"
	"github.com/vamps-core/routes"
	"github.com/vamps-core/commons"
	_ "github.com/go-sql-driver/mysql"
	"path/filepath"
)

const SERVER_HOME string = "SERVER_HOME"
const SERVER_LOG_FILE_NAME string = "server.log"
const ACCESS_LOG_FILE_NAME string = "http-access.log"
const DEFUALT_CONFIG_FILE_NAME string = "config.default.yaml"
const CONFIG_FILE_NAME string = "config.yaml"
const SERVER_CONFIGS_DIRECTORY = "configs"

type ServerConfigs struct {
	Hostname           string
	HttpPort           int
	HttpsPort          int
	ReadTimeOut        int
	WriteTimeOut       int
	SSLCertificateFile string
	SSLKeyFile         string
	TraceLogFile       string
	EnableTrace        bool
	EnableAccessLogs   bool
	LogsDirectory      string
}

var ServerHome string
var serverConfigs ServerConfigs
var logHandler http.Handler
var serverLogFile os.File
var httpAccessLogFile os.File

func init() {
	ServerHome = os.Getenv(SERVER_HOME)
	if ( len(ServerHome) <= 0 ) {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Fatal("Error while determining the server home. Please set the SERVER_HOME varaible and restart.")
		}
		ServerHome = dir
		log.Info(ServerHome)
	}

	viper.New()
	viper.AddConfigPath(ServerHome + "/" + SERVER_CONFIGS_DIRECTORY)
	viper.SetConfigName("config")
	if _, err := os.Stat(ServerHome + "/" + SERVER_CONFIGS_DIRECTORY + "/" + CONFIG_FILE_NAME); os.IsNotExist(err) {
		viper.SetConfigName("config.default")
	}

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		log.Error("Fatal error config file: %s \n", err)
	}

	configsMap := viper.GetStringMap("serverConfigs")
	serverConfigs.HttpPort = configsMap["httpPort"].(int)
	serverConfigs.ReadTimeOut = configsMap["readTimeOut"].(int)
	serverConfigs.WriteTimeOut = configsMap["writeTimeOut"].(int)
	serverConfigs.LogsDirectory = configsMap["logsDirectory"].(string)
	serverConfigs.EnableAccessLogs = configsMap["enableAccessLogs"].(bool)
	serverConfigs.SSLCertificateFile = configsMap["certificateFile"].(string)
	serverConfigs.SSLKeyFile = configsMap["keyFile"].(string)

	serverLogFile, err := os.OpenFile(serverConfigs.LogsDirectory + "/" + SERVER_LOG_FILE_NAME, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error while opening server log file: %v", err)
	}

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(serverLogFile)
	log.SetLevel(log.DebugLevel)

	httpAccessLogFile, err := os.OpenFile(serverConfigs.LogsDirectory + "/" + ACCESS_LOG_FILE_NAME, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error while trying to open the access log file: %v", err)
	}

	if serverConfigs.EnableAccessLogs {
		logHandler = handlers.LoggingHandler(httpAccessLogFile, http.DefaultServeMux)
	}

	dbConfigs := make(map[string] commons.DBConfigs)
	databases := viper.Get("dbConfigs").([]interface{})
	for i, _ := range databases {
		database := databases[i].(map[interface{}]interface{})
		dbConfigs[database["name"].(string)] = commons.DBConfigs{
			Dialect: database["dialect"].(string),
			DBName: database["dbname"].(string),
			Address: database["address"].(string),
			Parameters: database["parameters"].(string),
			Username: database["username"].(string),
			Password: database["password"].(string),
		}
	}
	commons.ConstructConnectionPool(dbConfigs)
}

func main() {
	defer serverLogFile.Close()
	defer httpAccessLogFile.Close()

	router := routes.NewRouter()
	http.Handle("/", router)

	httpsServer := &http.Server{
		Addr:           ":" + strconv.Itoa(serverConfigs.HttpPort),
		Handler:        logHandler,
		ReadTimeout:    time.Duration(serverConfigs.ReadTimeOut) * time.Second,
		WriteTimeout:   time.Duration(serverConfigs.WriteTimeOut) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Info("Starting server on port : " + strconv.Itoa(serverConfigs.HttpPort))
	log.Fatal("HTTP Server error: ", httpsServer.ListenAndServeTLS(serverConfigs.SSLCertificateFile, serverConfigs.SSLKeyFile))
}

