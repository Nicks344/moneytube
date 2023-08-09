package config

import "github.com/spf13/viper"

func init() {
	viper.AddConfigPath(".")
	viper.ReadInConfig()
	if GetAELang() == "" {
		viper.Set("ae_lang", "EN")
		viper.WriteConfig()
	}
}
