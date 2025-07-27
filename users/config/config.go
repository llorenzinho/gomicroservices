package config

import (
	"strings"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port int    `json:"port"`
	Host string `json:"host"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
	SSLMode  string `json:"sslmode"`
}
type JWTConfig struct {
	Secret            string `json:"secret"`
	AccessExpiration  int    `json:"accessExpiration"`
	RefreshExpiration int    `json:"refreshExpiration"`
}

type RabbitMQConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`

	UserCreateQueue string `json:"userCreateQueue"`
	UserUpdateQueue string `json:"userUpdateQueue"`
}

type Config struct {
	Server         ServerConfig   `json:"server"`
	Database       DatabaseConfig `json:"database"`
	RabbitMQConfig RabbitMQConfig `json:"rabbitmq" mapstructure:"rabbitmq"`
	JWT            JWTConfig      `json:"jwt"`
}

func setDeafult() {
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("jwt.expiration", 3600)
}

func LoadConfig() (*Config, error) {
	var config Config
	setDeafult()

	viper.SetConfigFile("config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		return &config, err
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.Unmarshal(&config)
	return &config, nil
}
