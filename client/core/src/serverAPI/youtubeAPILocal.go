package serverAPI

import (
	"errors"
	"net/http"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

func getService(apikey string) (*youtube.Service, error) {
	return youtube.New(&http.Client{
		Transport: &transport.APIKey{Key: apikey},
	})
}

func GetVideoIDsByChannelIDA(apikey string, channelID string, pageToken string) ([]string, error) {
	service, err := getService(apikey)
	if err != nil {
		return nil, err
	}
	list, err := service.Search.List([]string{"snippet"}).Type("video").ChannelId(channelID).MaxResults(50).PageToken(pageToken).Do()
	if err != nil {
		return nil, err
	}
	res := make([]string, len(list.Items))
	for i, item := range list.Items {
		res[i] = item.Id.VideoId
	}
	if list.NextPageToken != "" {
		next, err := GetVideoIDsByChannelIDA(apikey, channelID, list.NextPageToken)
		if err != nil {
			return nil, err
		}
		res = append(res, next...)
	}
	return res, nil
}

func GetChannelIDByUsernameA(apikey string, username string) (string, error) {
	service, err := getService(apikey)
	if err != nil {
		return "", err
	}
	res, err := service.Search.List([]string{"id"}).Type("channel").Q(username).Do()
	if err != nil {
		return "", err
	}
	if len(res.Items) == 0 {
		return "", errors.New("cannot find channel by username: " + username)
	}
	return res.Items[0].Id.ChannelId, nil
}
