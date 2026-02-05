package config

import (
	"errors"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL   string        `mapstructure:"DATABASE_URL"`
	ServerPort    string        `mapstructure:"SERVER_PORT"`
	JWTSecret     string        `mapstructure:"JWT_SECRET"`
	TokenDuration time.Duration `mapstructure:"TOKEN_DURATION"`
	Environment   string        `mapstructure:"ENVIRONMENT"`
}

func LoadConfig() (config Config, err error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("TOKEN_DURATION", "24h")
	viper.SetDefault("ENVIRONMENT", "development")

	err = viper.ReadInConfig()
	if err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			err = nil
		} else {
			return
		}
	}

	err = viper.Unmarshal(&config)
	return
}
