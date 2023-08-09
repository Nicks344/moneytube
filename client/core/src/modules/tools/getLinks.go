package tools

import (
	"encoding/json"
	"fmt"

	"github.com/Nicks344/moneytube/client/core/src/serverAPI"
	"github.com/Nicks344/moneytube/client/core/src/utils/parseLinks"

	"google.golang.org/api/googleapi"
)

type GetLinksWorker struct {
	*ToolService

	Channels []string
}

func (w *GetLinksWorker) Start() {
	w.maxProgress = len(w.Channels)
	w.handleProgressChanged("")

	for _, channel := range w.Channels {
		if w.isCancelled(true) {
			return
		}

		res, err := getLinksByChannel(channel)
		if err != nil {
			w.handleError(err)
			return
		}
		resStr, err := json.Marshal(res)
		if err != nil {
			w.handleError(err)
			return
		}

		w.progress++
		w.handleProgressChanged(string(resStr))
	}
}

func getLinksByChannel(channel string) ([]string, error) {
	channelID, err := parseLinks.GetChannelIDByLink(channel)
	if err != nil {
		return nil, err
	}
	list, err := serverAPI.GetVideoIDsByChannelID(channelID)
	if err != nil {
		e, ok := err.(*googleapi.Error)
		if !ok {
			return nil, err
		}
		if e.Message == "Invalid channel." {
			return []string{fmt.Sprintf("%s: Неверно указана ссылка на канал", channel)}, nil
		}
		return nil, err
	}
	for i, item := range list {
		list[i] = "https://www.youtube.com/watch?v=" + item
	}
	return list, nil
}
