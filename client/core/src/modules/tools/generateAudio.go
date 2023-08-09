package tools

import (
	"io/ioutil"
	"path/filepath"

	"github.com/Nicks344/moneytube/client/core/src/config"
	"github.com/Nicks344/moneytube/client/core/src/serverAPI"
	"github.com/Nicks344/moneytube/moneytubemodel"
)

const (
	GoogleSpeech = 10
	YandexSpeech = 20
	VoiceRSS     = 30
	SpeechPro    = 40
)

type AudioGeneratorWorker struct {
	*ToolService

	Api              int
	CropAudio        bool
	CropAudioVariant int
	MaxTime          int
	Lang             string
	Voice            string
	Speed            int
	TextFiles        []string
	ResultPath       string
	GSpeechProps     struct {
		Profiles []string
		Pitch    int
	}
}

func (w *AudioGeneratorWorker) Start() {
	w.maxProgress = len(w.TextFiles)
	w.handleProgressChanged("")

	for _, file := range w.TextFiles {
		textBytes, err := ioutil.ReadFile(file)
		if err != nil {
			w.handleError(err)
			return
		}

		filename := filepath.Base(file)
		resultFile := filepath.Join(w.ResultPath, filename+".mp3")

		data := moneytubemodel.AudioGenerateInput{
			Api:              w.Api,
			CropAudio:        w.CropAudio,
			CropAudioVariant: w.CropAudioVariant,
			MaxTime:          w.MaxTime,
			Voice:            w.Voice,
			Lang:             w.Lang,
			Speed:            w.Speed,
			Text:             string(textBytes),
			GSpeechProps:     w.GSpeechProps,
		}

		switch w.Api {
		case GoogleSpeech:
			data.ApiKey = config.GetGoogleSpeechApiKey()

		case YandexSpeech:
			data.ApiKey = config.GetYandexSpeechApiKey()

		case VoiceRSS:
			data.ApiKey = config.GetVoiceRSSApiKey()

		case SpeechPro:
			data.SpeechProCredentials = config.GetSpeechProCreds()
		}

		err = serverAPI.GenerateAudio(data, resultFile)
		if err != nil {
			w.handleError(err)
			return
		}
		if w.isCancelled(true) {
			return
		}

		w.progress++
		w.handleProgressChanged("")
	}
}
