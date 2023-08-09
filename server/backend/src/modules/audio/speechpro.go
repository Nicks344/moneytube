package audio

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"os"
	"strings"

	"github.com/Nicks344/moneytube/moneytubemodel"
	"github.com/Nicks344/moneytube/server/backend/src/modules/audio/speechpro_api"
)

func SpeechProGenerate(data moneytubemodel.AudioGenerateInput, resultFile string) error {
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
