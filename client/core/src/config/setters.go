package config

import (
	"github.com/spf13/viper"
)

func SetApiKey(key string) error {
	viper.Set("api_key", key)
	return viper.WriteConfig()
}

func SetLicenseKey(key string) error {
	viper.Set("license_key", key)
	return viper.WriteConfig()
}
