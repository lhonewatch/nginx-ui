package settings

import (
	"log"
	"os"
	"sync"

	"github.com/spf13/viper"
)

var (
	once sync.Once

	// App holds application-level settings
	App = &AppSettings{}

	// Nginx holds nginx-related settings
	Nginx = &NginxSettings{}

	// Server holds HTTP server settings
	Server = &ServerSettings{}

	// Auth holds authentication settings
	Auth = &AuthSettings{}
)

// AppSettings contains general application configuration
type AppSettings struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Debug   bool   `mapstructure:"debug"`
}

// NginxSettings contains nginx-specific configuration
type NginxSettings struct {
	ConfigDir  string `mapstructure:"config_dir"`
	PIDPath    string `mapstructure:"pid_path"`
	LogDir     string `mapstructure:"log_dir"`
	AccessLog  string `mapstructure:"access_log"`
	ErrorLog   string `mapstructure:"error_log"`
}

// ServerSettings contains HTTP server configuration
type ServerSettings struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	HTTPS    bool   `mapstructure:"https"`
	CertPath string `mapstructure:"cert_path"`
	KeyPath  string `mapstructure:"key_path"`
}

// AuthSettings contains authentication configuration
type AuthSettings struct {
	JWTSecret   string `mapstructure:"jwt_secret"`
	TokenExpiry int    `mapstructure:"token_expiry"` // in hours
	BCryptCost  int    `mapstructure:"bcrypt_cost"`
}

// Init loads configuration from the given path (or defaults)
func Init(cfgPath string) {
	once.Do(func() {
		v := viper.New()

		// Set defaults
		v.SetDefault("app.name", "nginx-ui")
		v.SetDefault("app.debug", false)
		v.SetDefault("server.host", "0.0.0.0")
		v.SetDefault("server.port", 9000)
		v.SetDefault("server.https", false)
		v.SetDefault("nginx.config_dir", "/etc/nginx")
		v.SetDefault("nginx.pid_path", "/var/run/nginx.pid")
		v.SetDefault("nginx.log_dir", "/var/log/nginx")
		v.SetDefault("nginx.access_log", "/var/log/nginx/access.log")
		v.SetDefault("nginx.error_log", "/var/log/nginx/error.log")
		v.SetDefault("auth.token_expiry", 24)
		v.SetDefault("auth.bcrypt_cost", 10)

		// Allow environment variable overrides with NGINX_UI_ prefix
		v.SetEnvPrefix("NGINX_UI")
		v.AutomaticEnv()

		if cfgPath != "" {
			v.SetConfigFile(cfgPath)
		} else {
			v.SetConfigName("app")
			v.SetConfigType("ini")
			v.AddConfigPath("/etc/nginx-ui/")
			v.AddConfigPath("$HOME/.nginx-ui")
			v.AddConfigPath(".")
		}

		if err := v.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				log.Printf("[settings] warning: could not read config file: %v", err)
			}
		} else {
			log.Printf("[settings] loaded config from %s", v.ConfigFileUsed())
		}

		if err := v.UnmarshalKey("app", App); err != nil {
			log.Fatalf("[settings] failed to unmarshal app settings: %v", err)
		}
		if err := v.UnmarshalKey("nginx", Nginx); err != nil {
			log.Fatalf("[settings] failed to unmarshal nginx settings: %v", err)
		}
		if err := v.UnmarshalKey("server", Server); err != nil {
			log.Fatalf("[settings] failed to unmarshal server settings: %v", err)
		}
		if err := v.UnmarshalKey("auth", Auth); err != nil {
			log.Fatalf("[settings] failed to unmarshal auth settings: %v", err)
		}

		// Warn if JWT secret is not set
		if Auth.JWTSecret == "" {
			log.Println("[settings] warning: auth.jwt_secret is not set; using insecure default")
			Auth.JWTSecret = "changeme-insecure-default"
		}

		if App.Debug {
			log.SetFlags(log.LstdFlags | log.Lshortfile)
			log.Println("[settings] debug mode enabled")
		}

		_ = os.MkdirAll(Nginx.LogDir, 0755)
	})
}
