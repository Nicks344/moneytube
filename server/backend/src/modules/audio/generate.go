package audio

import (
	"errors"
	"os/exec"

	"github.com/Nicks344/moneytube/client/core/src/config"
	"github.com/Nicks344/moneytube/client/core/src/utils"
	"github.com/Nicks344/moneytube/moneytubemodel"
)

const (
	GoogleSpeech = 10
	YandexSpeech = 20
	VoiceRSS     = 30
	SpeechPro    = 40
)

func Generate(data moneytubemodel.AudioGenerateInput, resultFile string) error {
	var err error
	switch data.Api {
	case YandexSpeech:
		err = YandexSpeechGenerate(data, resultFile)

	case VoiceRSS:
		err = VoiceRSSGenerate(data, resultFile)

	case GoogleSpeech:
		err = GSpeechGenerate(data, resultFile)

	case SpeechPro:
		err = SpeechProGenerate(data, resultFile)

	default:
		err = errors.New("unknown service")
	}

	if err != nil {
		return err
	}

	if data.CropAudio {
		maxTime := float64(data.MaxTime)
		if err := utils.AdaptAudio(resultFile, maxTime, data.CropAudioVariant); err != nil {
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func ConvertToMp3(filename string, resultFile string) error {
	cmd := exec.Command(config.GetFFmpegBin(), "-i", filename, "-ar", "44100", "-acodec", "libmp3lame", resultFile)
	return cmd.Run()
}
