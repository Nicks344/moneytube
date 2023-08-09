package config

import (
	"strconv"

	"github.com/Nicks344/moneytube/moneytubemodel"
	"github.com/spf13/viper"
)

func GetAEMaxMemoryPercent() int {
	return viper.GetInt("ae_memory_persent")
}

func GetYandexSpeechApiKey() string {
	return viper.GetString("ys_api_key")
}

func GetGoogleSpeechApiKey() string {
	return viper.GetString("gs_api_key")
}

func GetVoiceRSSApiKey() string {
	return viper.GetString("vrs_api_key")
}

func GetFFmpegBin() string {
	return viper.GetString("ffmpeg_bin")
}

func GetFFprobeBin() string {
	return viper.GetString("ffprobe_bin")
}

func GetYouTubeApiKey() string {
	return viper.GetString("yt_api_key")
}

func GetAerenderPath() string {
	return viper.GetString("ae_exe")
}

func GetApiKey() string {
	return viper.GetString("api_key")
}

func GetVersion() string {
	return viper.GetString("version")
}

func GetLicenseKey() string {
	return viper.GetString("license_key")
}

func GetShowBrowser() bool {
	return viper.GetBool("show_browser")
}

func GetAELang() string {
	return viper.GetString("ae_lang")
}

func GetSpeechProCreds() moneytubemodel.SpeechProCredentials {
	id, _ := strconv.Atoi(viper.GetString("speech_pro_id"))
	return moneytubemodel.SpeechProCredentials{
		ID:       id,
		Login:    viper.GetString("speech_pro_login"),
		Password: viper.GetString("speech_pro_password"),
	}
}
