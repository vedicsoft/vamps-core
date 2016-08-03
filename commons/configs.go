package commons

import (
	"os"
	"path/filepath"
	"github.com/spf13/viper"
	log "github.com/Sirupsen/logrus"
	"strconv"
)

type serverConfigs struct {
	Home               string
	Prefix             string
	IsMaster           bool
	PortOffset         int
	Hostname           string
	HttpPort           int
	HttpsPort          int
	CaddyPort          int
	ReadTimeOut        int
	WriteTimeOut       int
	CaddyPath          string
	CaddyFile          string
	SSLCertificateFile string
	SSLKeyFile         string
	JWTPrivateKeyFile  string
	JWTPublicKeyFile   string
	JWTExpirationDelta int
	TraceLogFile       string
	EnableTrace        bool
	EnableAccessLogs   bool
	LogsDirectory      string
	DBConfigMap        map[string]DBConfigs
}

type DBConfigs struct {
	Username   string
	Password   string
	Dialect    string
	DBName     string
	Address    string
	Parameters string
}

var ServerConfigurations serverConfigs

func init() {
	ServerConfigurations.Home = os.Getenv(SERVER_HOME)
	if ( len(ServerConfigurations.Home) <= 0 ) {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Fatal("Error while determining the server home. Please set the SERVER_HOME varaible and restart.")
		}
		ServerConfigurations.Home = dir
		os.Setenv(SERVER_HOME, dir)
	}

	viper.New()
	viper.AddConfigPath(ServerConfigurations.Home + "/" + SERVER_CONFIGS_DIRECTORY)
	viper.SetConfigName("config")
	if _, err := os.Stat(ServerConfigurations.Home + "/" + SERVER_CONFIGS_DIRECTORY + "/" + CONFIG_FILE_NAME); os.IsNotExist(err) {
		viper.SetConfigName("config.default")
	}

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		log.Error("Fatal error config file: %s \n", err)
	}

	configsMap := viper.GetStringMap("serverConfigs")
	ServerConfigurations.Prefix = configsMap["prefix"].(string)
	SERVER_PREFIX := ServerConfigurations.Prefix
	ServerConfigurations.IsMaster = configsMap["isMaster"].(bool)
	ServerConfigurations.PortOffset = configsMap["portOffset"].(int)
	ServerConfigurations.HttpPort = configsMap["httpPort"].(int)
	ServerConfigurations.HttpsPort = configsMap["httpsPort"].(int)
	ServerConfigurations.CaddyPort = configsMap["caddyPort"].(int)
	ServerConfigurations.ReadTimeOut = configsMap["readTimeOut"].(int)
	ServerConfigurations.WriteTimeOut = configsMap["writeTimeOut"].(int)
	ServerConfigurations.LogsDirectory = configsMap["logsDirectory"].(string)
	ServerConfigurations.EnableAccessLogs = configsMap["enableAccessLogs"].(bool)
	ServerConfigurations.CaddyPath = configsMap["caddyPath"].(string)
	ServerConfigurations.CaddyFile = configsMap["caddyFile"].(string)
	ServerConfigurations.JWTPrivateKeyFile = configsMap["JWTPrivateKeyFile"].(string)
	ServerConfigurations.JWTPublicKeyFile = configsMap["JWTPublicKeyFile"].(string)
	ServerConfigurations.JWTExpirationDelta = configsMap["JWTExpirationDelta"].(int)
	ServerConfigurations.SSLCertificateFile = configsMap["certificateFile"].(string)
	ServerConfigurations.SSLKeyFile = configsMap["keyFile"].(string)

	//Exporting variables for other services (Caddy)
	os.Setenv("PATH", os.Getenv("PATH") + ":" + ServerConfigurations.Home + "/bin")
	os.Setenv("CADDYPATH", ServerConfigurations.CaddyPath)
	os.Setenv(SERVER_PREFIX + "CADDY_PORT", strconv.Itoa(ServerConfigurations.CaddyPort + ServerConfigurations.PortOffset))
	os.Setenv(SERVER_PREFIX + "HTTPS_PORT", strconv.Itoa(ServerConfigurations.HttpsPort + ServerConfigurations.PortOffset))
	os.Setenv(SERVER_PREFIX + "CERTIFICATE_FILE", ServerConfigurations.SSLCertificateFile)
	os.Setenv(SERVER_PREFIX + "KEY_FILE", ServerConfigurations.SSLKeyFile)
	os.Setenv(SERVER_PREFIX + JWT_PRIVATE_KEY_FILE, ServerConfigurations.JWTPrivateKeyFile)
	os.Setenv(SERVER_PREFIX + JWT_PUBLIC_KEY_FILE, ServerConfigurations.JWTPublicKeyFile)
	os.Setenv(SERVER_PREFIX + JWT_EXPIRATION_DELTA, strconv.Itoa(ServerConfigurations.JWTExpirationDelta))

	ServerConfigurations.DBConfigMap = make(map[string]DBConfigs)
	databases := viper.Get("dbConfigs").([]interface{})
	for i, _ := range databases {
		database := databases[i].(map[interface{}]interface{})
		ServerConfigurations.DBConfigMap[database["name"].(string)] = DBConfigs{
			Dialect: database["dialect"].(string),
			DBName: database["dbname"].(string),
			Address: database["address"].(string),
			Parameters: database["parameters"].(string),
			Username: database["username"].(string),
			Password: database["password"].(string),
		}
	}
}