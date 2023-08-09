package serverAPI

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Nicks344/moneytube/client/core/src/config"

	"github.com/imroc/req"
)

func GetVideoIDsByPlaylistID(playlistID string) (result []string, err error) {
	var resp *req.Resp
	resp, err = req.Get(fmt.Sprintf("%s/api/user/v1/youtube/videoIDs/by/playlistID?apikey=%s&playlistID=%s", host, config.GetYouTubeApiKey(), playlistID), getAuthHeaders())
	if err != nil {
		return
	}

	if err = checkError(resp); err != nil {
		if strings.Contains(err.Error(), "quotaExceeded") {
			err = errors.New("Закончилась квота youtube api")
		}
		return
	}

	err = parseAnswer(resp, &result)
	return
}

func GetVideoIDsByChannelID(channelID string) (result []string, err error) {
	var resp *req.Resp
	resp, err = req.Get(fmt.Sprintf("%s/api/user/v1/youtube/videoIDs/by/channelID?apikey=%s&channelID=%s", host, config.GetYouTubeApiKey(), channelID), getAuthHeaders())
	if err != nil {
		return
	}

	if err = checkError(resp); err != nil {
		return
	}

	err = parseAnswer(resp, &result)
	return
}

func GetChannelIDByUsername(username string) (result string, err error) {
	var resp string
	resp, err = GetChannelIDByUsernameA(config.GetYouTubeApiKey(), username)
	if err != nil {
		return "", err
	}
	return resp, nil

	//var resp *req.Resp
	//resp, err = req.Get(fmt.Sprintf("%s/api/user/v1/youtube/channelID/by/username?apikey=%s&username=%s", host, config.GetYouTubeApiKey(), username), getAuthHeaders())
	//if err != nil {
	//	return
	//}

	//if err = checkError(resp); err != nil {
	//	return
	//}

	//err = parseAnswer(resp, &result)
	//return
}
