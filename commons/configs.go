package commons

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
	LogLevel           string
	DBConfigMap        map[string]DBConfigs
	KafkaConfigs       map[string]KafkaConfig
	ConfigMap          map[string]interface{}
	RedisConfigs       RedisConfigs
	ExternalServices   map[string]ExternalServicesConfigs
	TenantConfigs      TenantConfigsInfo
}

type DBConfigs struct {
	Username   string
	Password   string
	Dialect    string
	DBName     string
	Address    string
	Parameters string
}

type ExternalServicesConfigs struct {
	Service string
	Path    string
}

type TenantConfigsInfo struct {
	DefaultRoles []string
}

type RedisConfigs struct {
	Address  string
	Password string
}

var ServerConfigurations serverConfigs

