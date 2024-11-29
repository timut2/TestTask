package config

import (
	"log"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Port           uint32 `mapstructure:"PORT"`
	DBUser         string `mapstructure:"DB_USER"`
	DBPassword     string `mapstructure:"DB_PASSWORD"`
	DBHost         string `mapstructure:"DB_HOST"`
	DBPort         uint32 `mapstructure:"DB_PORT"`
	DBName         string `mapstructure:"DB_NAME"`
	ExternalApiUrl string `mapstrucutre:"API_URL"`
}

func Load() (*Config, error) {
	config := &Config{}

	viper.SetConfigName("config")
	viper.SetConfigType("env")

	configDir, err := filepath.Abs("./config") // Adjust this path if necessary
	if err != nil {
		log.Fatal("Error getting config directory: ", err)
	}
	viper.AddConfigPath(configDir)

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the .env file: ", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatal("Can't load the .env file: ", err)
	}

	return config, nil
}
