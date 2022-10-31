package config

import (
	"errors"
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

	dbPath := viper.GetString("db_path")
	if dbPath == "" {
		return nil, errors.New("please specify env DB_PATH")
	}

	dbUser := viper.GetString("db_user")
	if dbUser == "" {
		return nil, errors.New("please specify env DB_USER")
	}

	dbPass := viper.GetString("db_pass")
	if dbPass == "" {
		return nil, errors.New("please specify env DB_PASS")
	}

	apiServAddr := viper.GetString("api_serv_addr")
	if apiServAddr == "" {
		return nil, errors.New("please specify env API_SERV_ADDR")
	}

	cfg := Config{
		DBPath:      dbPath,
		DBUser:      dbUser,
		DBPass:      dbPass,
		ApiServAddr: apiServAddr,
	}

	return &cfg, nil
}
