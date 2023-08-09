package youtubeAPI

import (
	"github.com/Nicks344/moneytube/client/core/src/serverAPI"
)

func GetVideoIDsByPlaylistID(playlistID string) ([]string, error) {
	return serverAPI.GetVideoIDsByPlaylistID(playlistID)
}

func GetVideoIDsByChannelID(channelID string) ([]string, error) {
	return serverAPI.GetVideoIDsByChannelID(channelID)
}

func GetChannelIDByUsername(username string) (string, error) {
	return serverAPI.GetChannelIDByUsername(username)
}
