package parseLinks

import (
	"errors"
	"net/url"
	"path"
	"strings"

	"github.com/Nicks344/moneytube/client/core/src/modules/youtubeAPI"
)

func GetPlaylistIDByLink(link string) (string, error) {
	urlObj, err := url.ParseRequestURI(link)
	if err != nil {
		return "", errors.New("invalid playlist link")
	}
	if !strings.Contains(urlObj.Path, "playlist") {
		return "", errors.New("invalid playlist link")
	}
	return urlObj.Query().Get("list"), nil
}

func GetChannelIDByLink(link string) (string, error) {
	urlObj, err := url.ParseRequestURI(link)
	if err != nil {
		return "", errors.New("invalid channel link")
	}
	if strings.Contains(urlObj.Path, "channel") {
		return path.Base(urlObj.Path), nil
	}
	if strings.Contains(urlObj.Path, "user") || strings.Contains(urlObj.Path, "c") {
		split := strings.Split(urlObj.Path, "/")
		if len(split) < 2 {
			return "", errors.New("invalid channel link")
		}
		for i, el := range split {
			if el == "user" || el == "c" {
				channelID, err := youtubeAPI.GetChannelIDByUsername(split[i+1])
				if err != nil {
					return "", err
				}
				return channelID, nil
			}
		}
	}
	if strings.HasPrefix(urlObj.Path, "/@") {
		channelID, err := youtubeAPI.GetChannelIDByUsername(urlObj.Path[2:])
		if err != nil {
			return "", err
		}
		return channelID, nil
	}
	return "", errors.New("invalid channel link")
}

func GetVideoIDsByLinks(links []string) []string {
	videoIDs := make([]string, len(links))
	for i, link := range links {
		urlObj, err := url.ParseRequestURI(link)
		if err != nil {
			continue
		}

		var id string

		switch urlObj.Host {
		case "youtu.be":
			id = path.Base(urlObj.Path)

		case "youtube.com", "www.youtube.com":
			id = urlObj.Query().Get("v")

		default:
			continue
		}

		videoIDs[i] = strings.Trim(id, "\r\n ")
	}
	return videoIDs
}
