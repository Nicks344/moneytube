package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var currEnv = os.Getenv("MONEYTUBE_ENV")

func Init(configPath, envPath string) error {
	if currEnv == "" {
		panic("MONEYTUBE_ENV variable is empty!")
	}

	consfigName := currEnv + ".config"

	viper.AddConfigPath(configPath)
	viper.SetConfigName(consfigName)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	viper.WatchConfig()

	return godotenv.Load(filepath.Join(envPath, ".env"))
}
