package youtubeAPI

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

func GetVideoIDsByPlaylistID(apikey string, playlistID string, pageToken string) ([]string, error) {
	service, err := getService(apikey)
	if err != nil {
		return nil, err
	}
	list, err := service.PlaylistItems.List([]string{"snippet"}).PlaylistId(playlistID).MaxResults(50).PageToken(pageToken).Do()
	if err != nil {
		return nil, err
	}
	res := make([]string, len(list.Items))
	for i, item := range list.Items {
		res[i] = item.Snippet.ResourceId.VideoId
	}
	if list.NextPageToken != "" {
		next, err := GetVideoIDsByPlaylistID(apikey, playlistID, list.NextPageToken)
		if err != nil {
			return nil, err
		}
		res = append(res, next...)
	}
	return res, nil
}

func GetVideoIDsByChannelID(apikey string, channelID string, pageToken string) ([]string, error) {
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
		next, err := GetVideoIDsByChannelID(apikey, channelID, list.NextPageToken)
		if err != nil {
			return nil, err
		}
		res = append(res, next...)
	}
	return res, nil
}

func GetChannelIDByUsername(apikey string, username string) (string, error) {
	service, err := getService(apikey)
	if err != nil {
		return "", err
	}
	res, err := service.Channels.List([]string{"snippet"}).ForUsername(username).Do()
	if err != nil {
		return "", err
	}
	if len(res.Items) == 0 {
		return "", errors.New("cannot find channel by username: " + username)
	}
	return res.Items[0].Id, nil
}
