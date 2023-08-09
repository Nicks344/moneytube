package audio

import (
	"errors"
	"strings"

	"github.com/Nicks344/moneytube/moneytubemodel"

	"github.com/imroc/req"
)

const voiceRSSURL = "http://api.voicerss.org"

func VoiceRSSGenerate(data moneytubemodel.AudioGenerateInput, resultFile string) error {
	lang := data.Lang
	voice := data.Voice
	if lang == "" {
		lang = data.Voice
		voice = ""
	} else {
		voice = strings.Split(voice, "(")[0]
	}

	params := req.Param{
		"src": data.Text,
		"hl":  lang,
		"v":   voice,
		"r":   data.Speed,
		"c":   "MP3",
		"f":   "48khz_16bit_stereo",
		"key": data.ApiKey,
	}
	res, err := req.Post(voiceRSSURL, params)
	if err != nil {
		return err
	}

	if res.Response().Header.Get("Content-Type") != "audio/mpeg" {
		resStr, err := res.ToString()
		if err != nil {
			return err
		}
		return errors.New(resStr)
	}

	return res.ToFile(resultFile)
}
