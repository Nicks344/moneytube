package serverAPI

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/Nicks344/moneytube/client/core/src/config"

	"github.com/Nicks344/moneytube/moneytubemodel"

	"github.com/imroc/req"
	"github.com/Nicks344/moneytube/server/backend/src/modules/audio/speechpro_api"
)

func GenerateAudio(data moneytubemodel.AudioGenerateInput, resultFile string) (err error) {

	if data.Api == 40 {
		ctx := context.TODO()
		credentials := speechpro_api.AuthRequestDto{data.SpeechProCredentials.Login, int64(data.SpeechProCredentials.ID), data.SpeechProCredentials.Password}
		client := speechpro_api.NewAPIClient(speechpro_api.NewConfiguration())
		sessionApi := client.SessionApi

		loginResponse, _, err := sessionApi.Login(ctx, credentials)
		if err != nil {
			return err
		}

		sessionId := loginResponse.SessionId
		synthesis := client.SynthesizeApi
		voiceSplit := strings.Split(data.Voice, "(")
		voice := voiceSplit[0]

		synthesisText := &speechpro_api.SynthesizeText{"text/plain", data.Text}
		synthesisRequest := speechpro_api.SynthesizeRequest{synthesisText, voice, "audio/wav"}
		synthesisResponse, _, err := synthesis.Synthesize(context.Background(), sessionId, synthesisRequest, nil)
		if err != nil {
			return err
		}
		sound, err := base64.StdEncoding.DecodeString(synthesisResponse.Data)
		if err != nil {
			return err
		}

		os.Remove(resultFile)
		os.Remove(resultFile + ".wav")
		if err := ioutil.WriteFile(resultFile+".wav", sound, 0644); err != nil {
			return err
		}
		defer os.Remove(resultFile + ".wav")

		return ConvertToMp3(resultFile+".wav", resultFile)
	}
	var resp *req.Resp
	resp, err = req.Post(host+"/api/user/v1/audio/generate/", getAuthHeaders(), req.BodyJSON(data))
	if err != nil {
		return
	}

	if err = checkError(resp); err != nil {
		return
	}

	return resp.ToFile(resultFile)
}
func ConvertToMp3(filename string, resultFile string) error {
	cmd := exec.Command(config.GetFFmpegBin(), "-i", filename, "-ar", "44100", "-acodec", "libmp3lame", resultFile)
	return cmd.Run()
}
