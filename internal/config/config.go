package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	DBPath      string
	DBUser      string
	DBPass      string
	ApiServAddr string
}

func Init() (*Config, error) {

	err := viper.BindEnv("db_path")
	if err != nil {
		return nil, fmt.Errorf("can't read env: %w", err)
	}

	err = viper.BindEnv("db_user")
	if err != nil {
		return nil, fmt.Errorf("can't read env: %w", err)
	}

	err = viper.BindEnv("db_pass")
	if err != nil {
		return nil, fmt.Errorf("can't read env: %w", err)
	}

	err = viper.BindEnv("api_serv_addr")
	if err != nil {
		return nil, fmt.Errorf("can't read env: %w", err)
	}

	cfg := Config{
		DBPath:      viper.GetString("db_path"),
		DBUser:      viper.GetString("db_user"),
		DBPass:      viper.GetString("db_pass"),
		ApiServAddr: viper.GetString("api_serv_addr"),
	}

	return &cfg, nil
}
