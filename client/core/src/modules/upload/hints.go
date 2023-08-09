package upload

import (
	"errors"
	"fmt"
	"github.com/meandrewdev/logger"
	"time"

	"github.com/imroc/req"
	"github.com/tidwall/gjson"
)

type HintType int

const (
	HT_Video    HintType = 10
	HT_Playlist HintType = 20
	HT_Channel  HintType = 30
)

type VideoHintsParams struct {
	Context interface{}
	Query   string
	Headers req.Header
	Hints   []HintInfo
	VideoID string
}

type HintInfo struct {
	Type    HintType
	Data    string
	Time    time.Duration
	Message string
	Teaser  string
}

func editVideoHints(params VideoHintsParams, proxy string) (*req.Resp, error) {
	postBody := EditVideoData{
		Context:         params.Context,
		ExternalVideoID: params.VideoID,
	}

	infoCards := []InfoCards{}

	for i, hintInfo := range params.Hints {
		infoCard := InfoCards{
			VideoID:          params.VideoID,
			TeaserStartMs:    fmt.Sprintf("%d", hintInfo.Time.Milliseconds()),
			InfoCardEntityID: fmt.Sprintf("new-addition-%d", i+1),
			CustomMessage:    hintInfo.Message,
			TeaserText:       hintInfo.Teaser,
		}

		switch hintInfo.Type {
		case HT_Video:
			infoCard.VideoInfoCard = &VideoInfoCard{
				VideoID: hintInfo.Data,
			}

		case HT_Playlist:
			infoCard.PlaylistInfoCard = &PlaylistInfoCard{
				FullPlaylistID: hintInfo.Data,
			}

		case HT_Channel:
			infoCard.CollaboratorInfoCard = &CollaboratorInfoCard{
				ChannelID: hintInfo.Data,
			}
		}

		infoCards = append(infoCards, infoCard)
	}

	postBody.InfoCardEdit = InfoCardEdit{
		InfoCards: infoCards,
	}
	request := req.New()
	if proxy != "" {
		request.SetProxyUrl(proxy)
	}
	request.EnableCookie(false)

	resp, err := request.Post("https://studio.youtube.com/youtubei/v1/video_editor/edit_video?"+params.Query, params.Headers, req.BodyJSON(postBody))
	if err != nil {
		return nil, err
	}

	body := resp.String()
	status := gjson.Get(body, "executionStatus").String()

	if status != "EDIT_EXECUTION_STATUS_DONE" {
		logger.WarningF("error on edit video hints, body:\r\n%s", body)
		return nil, errors.New("edit video hints unsuccessfull")
	}

	return resp, nil
}

type EditVideoData struct {
	Context         interface{}  `json:"context,omitempty"`
	InfoCardEdit    InfoCardEdit `json:"infoCardEdit"`
	ExternalVideoID string       `json:"externalVideoId"`
}
type VideoInfoCard struct {
	VideoID string `json:"videoId"`
}
type PlaylistInfoCard struct {
	FullPlaylistID string `json:"fullPlaylistId"`
}
type CollaboratorInfoCard struct {
	ChannelID string `json:"channelId"`
}
type InfoCards struct {
	VideoID              string                `json:"videoId"`
	TeaserStartMs        string                `json:"teaserStartMs"`
	VideoInfoCard        *VideoInfoCard        `json:"videoInfoCard,omitempty"`
	InfoCardEntityID     string                `json:"infoCardEntityId"`
	CustomMessage        string                `json:"customMessage"`
	TeaserText           string                `json:"teaserText"`
	PlaylistInfoCard     *PlaylistInfoCard     `json:"playlistInfoCard,omitempty"`
	CollaboratorInfoCard *CollaboratorInfoCard `json:"collaboratorInfoCard,omitempty"`
}
type InfoCardEdit struct {
	InfoCards []InfoCards `json:"infoCards"`
}
