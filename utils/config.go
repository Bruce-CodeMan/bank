package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variables.
type Config struct {
	DBSource             string        `mapstructure:"DB_SOURCE"`
	ServerAddress        string        `mapstructure:"SERVER_ADDRESS"`
	TokenSynmmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	RedisAddress         string        `mapstructure:"REDIS_ADDRESS"`
}

func getProjectRootPath() string {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)          // /.../bank/utils
	return filepath.Join(basePath, "..") // back to /.../bank
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig() (config Config, err error) {
	rootPath := getProjectRootPath()

	// Get the Env name
	env := os.Getenv("APP_ENV")
	var fileName string
	if env == "" {
		fileName = "app.env"
	} else {
		fileName = fmt.Sprintf("app.%s.env", env)
	}

	configPath := filepath.Join(rootPath, fileName)
	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
