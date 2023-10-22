package util

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type EnvVar struct {
	DbName               string        `mapstructure:"DB_NAME"`
	DbUser               string        `mapstructure:"DB_USER"`
	DbPassword           string        `mapstructure:"DB_PASSWORD"`
	DbHost               string        `mapstructure:"DB_HOST"`
	DbPort               string        `mapstructure:"DB_PORT"`
	ServerHost           string        `mapstructure:"SERVER_HOST"`
	ServerPort           string        `mapstructure:"SERVER_PORT"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}
type Config struct {
	DBDriver             string
	DBSource             string
	ServerAddress        string
	ServerPort           string
	TokenSymmetricKey    string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

func getEnvVar(path string) (env EnvVar, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&env)
	return
}

func LoadConfig(path string) (Config, error) {
	var config Config
	env, err := getEnvVar(path)
	if err != nil {
		return config, err
	}

	config.DBDriver = "postgres" //hard-coded driver as postgres
	config.DBSource = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", env.DbUser, env.DbPassword, env.DbHost, env.DbPort, env.DbName)
	config.ServerAddress = fmt.Sprintf("%s:%s", env.ServerHost, env.ServerPort)
	config.ServerPort = env.ServerPort
	config.TokenSymmetricKey = env.TokenSymmetricKey
	config.AccessTokenDuration = env.AccessTokenDuration
	config.RefreshTokenDuration = env.RefreshTokenDuration
	return config, nil
}
