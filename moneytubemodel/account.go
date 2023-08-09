package moneytubemodel

import (
	"errors"
	"fmt"
	"strings"
)

const (
	ASNew    = 0
	ASAuth   = 10 //Авторизуюсь
	ASUpdate = 20 //Обновление данных
	ASReady  = 30 //Готов к работе
	ASError  = 40 //Ошибка
)

type Account struct {
	ID               int `bson:"_id"`
	Login            string
	Password         string
	ApiKey           string
	ChannelName      string
	ChannelID        string
	VideoCount       int
	PlaylistsCount   int
	SubscribersCount int
	ViewsCount       int
	Status           int
	ErrorMessage     string
	Proxy            string
	Group            string
	Cookies          string
	UserAgent        string
}

func (acc *Account) GetChannelLink() string {
	return fmt.Sprintf("https://www.youtube.com/channel/%s", acc.ChannelID)
}

func (acc *Account) GetFormattedProxy() (string, error) {
	if acc.Proxy == "" {
		return "", nil
	}
	result := ""
	protAndProxy := strings.Split(acc.Proxy, "://")
	protocol := protAndProxy[0]
	if len(protAndProxy) != 2 || (protocol != "socks5" && protocol != "http" && protocol != "https") {
		return "", errors.New("invalid proxy protocol")
	}

	result += protocol + "://"

	proxySplit := strings.Split(protAndProxy[1], ":")
	if len(proxySplit) < 2 {
		return "", errors.New("invalid proxy")
	}

	if len(proxySplit) == 4 {
		result += proxySplit[2] + ":" + proxySplit[3] + "@"
	}

	result += proxySplit[0] + ":" + proxySplit[1]

	return result, nil
}
