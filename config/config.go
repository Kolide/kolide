package config

import (
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	envPrefix = "KOLIDE"
)

// MysqlConfig defines configs related to MySQL
type MysqlConfig struct {
	Address  string
	Username string
	Password string
	Database string
}

// ServerConfig defines configs related to the Kolide server
type ServerConfig struct {
	Address string
	Cert    string
	Key     string
}

// AuthConfig defines configs related to user authorization
type AuthConfig struct {
	JwtKey      string
	BcryptCost  int
	SaltKeySize int
}

// AppConfig defines configs related to HTTP
type AppConfig struct {
	WebAddress string
}

// SMTPConfig defines configs related to SMTP email
type SMTPConfig struct {
	Server          string
	Username        string
	Password        string
	PoolConnections int
	TokenKeySize    int
}

// SessionConfig defines configs related to user sessions
type SessionConfig struct {
	KeySize           int
	ExpirationSeconds int
	CookieName        string
}

// OsqueryConfig defines configs related to osquery
type OsqueryConfig struct {
	EnrollSecret  string
	NodeKeySize   int
	StatusLogFile string
	ResultLogFile string
}

// LoggingConfig defines configs related to logging
type LoggingConfig struct {
	Debug         bool
	DisableBanner bool
}

// KolideConfig stores the application configuration. Each subcategory is
// broken up into it's own struct, defined above. When editing any of these
// structs, Manager.addConfigs and Manager.LoadConfig should be
// updated to set and retrieve the configurations as appropriate.
type KolideConfig struct {
	Mysql   MysqlConfig
	Server  ServerConfig
	Auth    AuthConfig
	App     AppConfig
	SMTP    SMTPConfig
	Session SessionConfig
	Osquery OsqueryConfig
	Logging LoggingConfig
}

// addConfigs adds the configuration keys and default values that will be
// filled into the KolideConfig struct
func (man Manager) addConfigs() {
	// MySQL
	man.addConfigString("mysql.address", "localhost:3306")
	man.addConfigString("mysql.username", "kolide")
	man.addConfigString("mysql.password", "kolide")
	man.addConfigString("mysql.database", "kolide")

	// Server
	man.addConfigString("server.address", "0.0.0.0:8080")
	man.addConfigString("server.cert", "./tools/osquery/kolide.crt")
	man.addConfigString("server.key", "./tools/osquery/kolide.key")

	// Auth
	man.addConfigString("auth.jwt_key", "CHANGEME")
	man.addConfigInt("auth.bcrypt_cost", 12)
	man.addConfigInt("auth.salt_key_size", 24)

	// App
	man.addConfigString("app.web_address", "0.0.0.0:8080")

	// SMTP
	man.addConfigString("smtp.server", "")
	man.addConfigString("smtp.username", "")
	man.addConfigString("smtp.password", "")
	man.addConfigInt("smtp.pool_connections", 4)
	man.addConfigInt("smtp.token_key_size", 24)

	// Session
	man.addConfigInt("session.key_size", 64)
	man.addConfigInt("session.expiration_seconds", 60*60*24*90)
	man.addConfigString("session.cookie_name", "KolideSession")

	// Osquery
	man.addConfigString("osquery.enroll_secret", "")
	man.addConfigInt("osquery.node_key_size", 24)
	man.addConfigString("osquery.status_log_file", "/tmp/osquery_status")
	man.addConfigString("osquery.result_log_file", "/tmp/osquery_result")

	// Logging
	man.addConfigBool("logging.debug", false)
	man.addConfigBool("logging.disable_banner", false)
}

// LoadConfig will load the config variables into a fully initialized
// KolideConfig struct
func (man Manager) LoadConfig() KolideConfig {
	man.loadConfigFile()

	return KolideConfig{
		Mysql: MysqlConfig{
			Address:  man.getConfigString("mysql.address"),
			Username: man.getConfigString("mysql.username"),
			Password: man.getConfigString("mysql.password"),
			Database: man.getConfigString("mysql.database"),
		},
		Server: ServerConfig{
			Address: man.getConfigString("server.address"),
			Cert:    man.getConfigString("server.cert"),
			Key:     man.getConfigString("server.key"),
		},
		Auth: AuthConfig{
			JwtKey:      man.getConfigString("auth.jwt_key"),
			BcryptCost:  man.getConfigInt("auth.bcrypt_cost"),
			SaltKeySize: man.getConfigInt("auth.salt_key_size"),
		},
		App: AppConfig{
			WebAddress: man.getConfigString("app.web_address"),
		},
		SMTP: SMTPConfig{
			Server:          man.getConfigString("smtp.server"),
			Username:        man.getConfigString("smtp.username"),
			Password:        man.getConfigString("smtp.password"),
			PoolConnections: man.getConfigInt("smtp.pool_connections"),
			TokenKeySize:    man.getConfigInt("smtp.token_key_size"),
		},
		Session: SessionConfig{
			KeySize:           man.getConfigInt("session.key_size"),
			ExpirationSeconds: man.getConfigInt("session.expiration_seconds"),
			CookieName:        man.getConfigString("session.cookie_name"),
		},
		Osquery: OsqueryConfig{
			EnrollSecret:  man.getConfigString("osquery.enroll_secret"),
			NodeKeySize:   man.getConfigInt("osquery.node_key_size"),
			StatusLogFile: man.getConfigString("osquery.status_log_file"),
			ResultLogFile: man.getConfigString("osquery.result_log_file"),
		},
		Logging: LoggingConfig{
			Debug:         man.getConfigBool("logging.debug"),
			DisableBanner: man.getConfigBool("logging.disable_banner"),
		},
	}
}

// envNameFromConfigKey converts a config key into the corresponding
// environment variable name
func envNameFromConfigKey(key string) string {
	return envPrefix + "_" + strings.ToUpper(strings.Replace(key, ".", "_", -1))
}

// flagNameFromConfigKey converts a config key into the corresponding flag name
func flagNameFromConfigKey(key string) string {
	return strings.Replace(key, ".", "_", -1)
}

// Manager manages the addition and retrieval of config values for Kolide
// configs. It's only public API method is LoadConfig, which will return the
// populated KolideConfig struct.
type Manager struct {
	viper    *viper.Viper
	command  *cobra.Command
	defaults map[string]interface{}
}

// NewManager initializes a Manager wrapping the provided cobra
// command. All config flags will be attached to that command (and inherited by
// the subcommands). Typically this should be called just once, with the root
// command.
func NewManager(command *cobra.Command) Manager {
	man := Manager{
		viper:    viper.New(),
		command:  command,
		defaults: map[string]interface{}{},
	}
	man.addConfigs()
	return man
}

// addDefault will check for duplication, then add a default value to the
// defaults map
func (man Manager) addDefault(key string, defVal interface{}) {
	if _, exists := man.defaults[key]; exists {
		panic("Trying to add duplicate config for key " + key)
	}

	man.defaults[key] = defVal
}

// getInterfaceVal is a helper function used by the getConfig* functions to
// retrieve the config value as interface{}, which will then be cast to the
// appropriate type by the getConfig* function.
func (man Manager) getInterfaceVal(key string) interface{} {
	interfaceVal := man.viper.Get(key)
	if interfaceVal == nil {
		var ok bool
		interfaceVal, ok = man.defaults[key]
		if !ok {
			panic("Tried to look up default value for nonexistent config option: " + key)
		}
	}
	return interfaceVal
}

// addConfigString adds a string config to the config options
func (man Manager) addConfigString(key string, defVal string) {
	man.command.PersistentFlags().String(flagNameFromConfigKey(key), defVal, "Env: "+envNameFromConfigKey(key))
	man.viper.BindPFlag(key, man.command.PersistentFlags().Lookup(flagNameFromConfigKey(key)))
	man.viper.BindEnv(key, envNameFromConfigKey(key))

	// Add default
	man.addDefault(key, defVal)
}

// getConfigString retrieves a string from the loaded config
func (man Manager) getConfigString(key string) string {
	interfaceVal := man.getInterfaceVal(key)
	stringVal, err := cast.ToStringE(interfaceVal)
	if err != nil {
		panic("Unable to cast to string for key " + key + ": " + err.Error())
	}

	return stringVal
}

// addConfigInt adds a int config to the config options
func (man Manager) addConfigInt(key string, defVal int) {
	man.command.PersistentFlags().Int(flagNameFromConfigKey(key), defVal, "Env: "+envNameFromConfigKey(key))
	man.viper.BindPFlag(key, man.command.PersistentFlags().Lookup(flagNameFromConfigKey(key)))
	man.viper.BindEnv(key, envNameFromConfigKey(key))

	// Add default
	man.addDefault(key, defVal)
}

// getConfigInt retrieves a int from the loaded config
func (man Manager) getConfigInt(key string) int {
	interfaceVal := man.getInterfaceVal(key)
	intVal, err := cast.ToIntE(interfaceVal)
	if err != nil {
		panic("Unable to cast to int for key " + key + ": " + err.Error())
	}

	return intVal
}

// addConfigBool adds a bool config to the config options
func (man Manager) addConfigBool(key string, defVal bool) {
	man.command.PersistentFlags().Bool(flagNameFromConfigKey(key), defVal, "Env: "+envNameFromConfigKey(key))
	man.viper.BindPFlag(key, man.command.PersistentFlags().Lookup(flagNameFromConfigKey(key)))
	man.viper.BindEnv(key, envNameFromConfigKey(key))

	// Add default
	man.addDefault(key, defVal)
}

// getConfigBool retrieves a bool from the loaded config
func (man Manager) getConfigBool(key string) bool {
	interfaceVal := man.getInterfaceVal(key)
	boolVal, err := cast.ToBoolE(interfaceVal)
	if err != nil {
		panic("Unable to cast to bool for key " + key + ": " + err.Error())
	}

	return boolVal
}

// loadConfigFile handles the loading of the config file.
func (man Manager) loadConfigFile() {
	configFile := man.command.PersistentFlags().Lookup("config").Value.String()
	if configFile != "" {
		man.viper.SetConfigFile(configFile)
	} else {
		man.viper.SetConfigName("kolide")
		man.viper.AddConfigPath(".")
		man.viper.AddConfigPath("$HOME")
		man.viper.AddConfigPath("./tools/app")
		man.viper.AddConfigPath("/etc/kolide")
	}

	man.viper.SetConfigType("yaml")

	err := man.viper.ReadInConfig()
	if err != nil {
		logrus.Fatalf("Error reading config file: %s", man.viper.ConfigFileUsed())
	}

	logrus.Info("Using config file: ", man.viper.ConfigFileUsed())
}
