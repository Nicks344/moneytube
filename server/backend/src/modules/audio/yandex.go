package audio

import (
	"fmt"
	"os"
	"strings"

	"github.com/Nicks344/moneytube/moneytubemodel"

	"github.com/imroc/req"
)

const ysURL = "https://tts.api.cloud.yandex.net/speech/v1/tts:synthesize"

func YandexSpeechGenerate(data moneytubemodel.AudioGenerateInput, resultFile string) error {
	lang := data.Lang
	voice := ""

	if lang == "" {
		voiceSplit := strings.Split(data.Voice, "_")
		lang = voiceSplit[0]
		voice = strings.Split(voiceSplit[1], "-")[0]
	} else {
		voiceSplit := strings.Split(data.Voice, "(")
		voice = strings.ToLower(voiceSplit[0])
	}

	params := req.Param{
		"text":  data.Text,
		"lang":  lang,
		"voice": voice,
		"speed": data.Speed / 10,
	}

	headers := req.Header{
		"Authorization": "Api-Key " + data.ApiKey,
	}

	res, err := req.Post(ysURL, params, headers)
	if err != nil {
		return err
	}

	if res.Response().StatusCode != 200 {
		var resObj map[string]interface{}
		err = res.ToJSON(&resObj)
		if err != nil {
			return err
		}
		return fmt.Errorf("%s: %s", resObj["error_code"], resObj["error_message"])
	}

	os.Remove(resultFile)
	os.Remove(resultFile + ".ogg")
	if err := res.ToFile(resultFile + ".ogg"); err != nil {
		return err
	}
	defer os.Remove(resultFile + ".ogg")

	return ConvertToMp3(resultFile+".ogg", resultFile)
}
