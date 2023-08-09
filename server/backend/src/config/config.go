package config

import "github.com/spf13/viper"

func GetDomain() string {
	return viper.GetString("domain")
}

func GetEnigmaProject() string {
	return viper.GetString("enigma_project_name")
}

func GetFfmpegBin() string {
	return viper.GetString("ffmpeg_bin")
}

func GetPort() string {
	return viper.GetString("port")
}
