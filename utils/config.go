package utils

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variables.
type Config struct {
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

func getProjectRootPath() string {
	_, b, _, _ := runtime.Caller(0)
	log.Printf("b: %s\n", b)
	basePath := filepath.Dir(b)          // /.../bank/utils
	return filepath.Join(basePath, "..") // back to /.../bank
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig() (config Config, err error) {
	rootPath := getProjectRootPath()
	fmt.Printf("rootPath: %s\n", rootPath)

	viper.SetConfigName("app")
	viper.SetConfigType("env") // json, xml, yaml and so on
	viper.AddConfigPath(rootPath)

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
