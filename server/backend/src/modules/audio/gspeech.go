package audio

import (
	"encoding/base64"
	"errors"
	"io/ioutil"
	"strings"

	"github.com/Nicks344/moneytube/moneytubemodel"

	"github.com/imroc/req"
)

func GSpeechGenerate(data moneytubemodel.AudioGenerateInput, resultFile string) error {
	lang := data.Lang
	voice := data.Voice
	gender := "MALE"
	if lang == "" {
		voiceSplit := strings.Split(data.Voice, "-")
		lang = strings.ToLower(voiceSplit[0] + "-" + voiceSplit[1])
	} else {
		voiceSplit := strings.Split(data.Voice, "(")
		gender = strings.TrimRight(voiceSplit[1], ")")
		voice = lang + "-" + voiceSplit[0]
	}

	reqData := GSpeechReq{
		Input: GSpeechReqInput{
			Text: data.Text,
		},
		Voice: GSpeechReqVoice{
			LanguageCode: lang,
			Name:         voice,
			SsmlGender:   gender,
		},
		AudioConfig: GSpeechReqAudioConfig{
			AudioEncoding:    "MP3",
			SpeakingRate:     float64(data.Speed) / 100,
			Pitch:            data.GSpeechProps.Pitch,
			EffectsProfileId: data.GSpeechProps.Profiles,
		},
	}

	res, err := req.Post("https://texttospeech.googleapis.com/v1/text:synthesize?key="+data.ApiKey, req.BodyJSON(reqData))
	if err != nil {
		return err
	}

	var ans GSpeechAns
	err = res.ToJSON(&ans)
	if err != nil {
		return err
	}
	if ans.Error.Message != "" {
		return errors.New(ans.Error.Message)
	}

	audioBytes, err := base64.StdEncoding.DecodeString(ans.AudioContent)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(resultFile, audioBytes, 0666)
}

type GSpeechAns struct {
	AudioContent string `json:"audioContent"`
	Error        struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
		Details []struct {
			Type            string `json:"@type"`
			FieldViolations []struct {
				Field       string `json:"field"`
				Description string `json:"description"`
			} `json:"fieldViolations"`
		} `json:"details"`
	} `json:"error"`
}

type GSpeechReq struct {
	Input       GSpeechReqInput       `json:"input"`
	Voice       GSpeechReqVoice       `json:"voice"`
	AudioConfig GSpeechReqAudioConfig `json:"audioConfig"`
}

type GSpeechReqInput struct {
	Text string `json:"text"`
}

type GSpeechReqVoice struct {
	LanguageCode string `json:"languageCode"`
	Name         string `json:"name"`
	SsmlGender   string `json:"ssmlGender"`
}

type GSpeechReqAudioConfig struct {
	AudioEncoding    string   `json:"audioEncoding"`
	SpeakingRate     float64  `json:"speakingRate"`
	Pitch            int      `json:"pitch"`
	EffectsProfileId []string `json:"effectsProfileId"`
}
