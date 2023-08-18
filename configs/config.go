package configs

import (
	"github.com/spf13/viper"
)

var cfg *config

type config struct {
	EnableDebug bool
	Server      struct {
		Port string
	}
	Postgres struct {
		Host     string
		Port     string
		User     string
		Password string
		Database string
	}
}

func init() {
	viper.SetDefault("ENABLE_DEBUG", false)
	viper.SetDefault("PORT", "80")
	viper.SetDefault("POSTGRES_DB", nil)
	viper.SetDefault("POSTGRES_HOST", nil)
	viper.SetDefault("POSTGRES_PORT", nil)
	viper.SetDefault("POSTGRES_USER", nil)
	viper.SetDefault("POSTGRES_PASS", nil)

	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()
	viper.ReadInConfig()

	cfg = new(config)

	cfg.EnableDebug = viper.GetBool("ENABLE_DEBUG")
	cfg.Server = struct {
		Port string
	}{
		Port: viper.GetString("PORT"),
	}
	cfg.Postgres = struct {
		Host     string
		Port     string
		User     string
		Password string
		Database string
	}{
		Host:     viper.GetString("POSTGRES_HOST"),
		Port:     viper.GetString("POSTGRES_PORT"),
		User:     viper.GetString("POSTGRES_USER"),
		Password: viper.GetString("POSTGRES_PASS"),
		Database: viper.GetString("POSTGRES_DB"),
	}

}

func GetConfig() config {
	return *cfg
}
