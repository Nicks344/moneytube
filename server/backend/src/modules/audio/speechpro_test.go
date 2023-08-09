package audio

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/Nicks344/moneytube/server/backend/src/modules/audio/speechpro_api"
	"golang.org/x/net/context"
)

type Test struct {
	Langs  map[string]string   `json:"langs"`
	Voices map[string][]string `json:"voices"`
}

func TestTTS(t *testing.T) {
	credentials := speechpro_api.AuthRequestDto{"-", 0, "-"}
	client := speechpro_api.NewAPIClient(speechpro_api.NewConfiguration())
	sessionApi := client.SessionApi
	loginResponse, _, err := sessionApi.Login(context.Background(), credentials)
	if err != nil {
		t.Fatal(err)
	}

	sessionId := loginResponse.SessionId
	synthesis := client.SynthesizeApi

	langs, _, err := synthesis.LanguageVoicesSupport(context.Background(), sessionId)
	if err != nil {
		t.Fatal(err)
	}

	test := Test{
		Langs:  map[string]string{},
		Voices: map[string][]string{},
	}

	firstVoice := ""
	for il, l := range langs {
		test.Langs[l.Name] = l.Name

		voices, _, err := synthesis.Voices(context.Background(), sessionId, l.Name)
		if err != nil {
			t.Fatal(err)
		}

		vs := []string{}
		for iv, v := range voices {
			vs = append(vs, fmt.Sprintf("%s(%s)", v.Name, v.Gender))
			if iv == 0 && il == 0 {
				firstVoice = v.Name
			}
		}

		test.Voices[l.Name] = vs
	}

	data, err := json.MarshalIndent(test, "", "\t")
	if err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile("langs.json", data, 0666); err != nil {
		t.Fatal(err)
	}

	synthesisText := &speechpro_api.SynthesizeText{"text/plain", "Тестовое сообщение"}
	synthesisRequest := speechpro_api.SynthesizeRequest{synthesisText, firstVoice, "audio/wav"}
	synthesisResponse, _, err := synthesis.Synthesize(context.Background(), sessionId, synthesisRequest, nil)
	if err != nil {
		t.Fatal(err)
	}
	sound, err := base64.StdEncoding.DecodeString(synthesisResponse.Data)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile("test.wav", sound, 0644)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGenerateAllVoices(t *testing.T) {
	credentials := speechpro_api.AuthRequestDto{"andrew.progs@gmail.com", 10437, "nKPc07nuU%"}
	client := speechpro_api.NewAPIClient(speechpro_api.NewConfiguration())
	sessionApi := client.SessionApi
	loginResponse, _, err := sessionApi.Login(context.Background(), credentials)
	if err != nil {
		t.Fatal(err)
	}

	sessionId := loginResponse.SessionId
	synthesis := client.SynthesizeApi

	langs, _, err := synthesis.LanguageVoicesSupport(context.Background(), sessionId)
	if err != nil {
		t.Fatal(err)
	}

	test := Test{
		Langs:  map[string]string{},
		Voices: map[string][]string{},
	}

	for _, l := range langs {
		test.Langs[l.Name] = l.Name

		voices, _, err := synthesis.Voices(context.Background(), sessionId, l.Name)
		if err != nil {
			t.Fatal(err)
		}

		vs := []string{}
		for _, v := range voices {
			vs = append(vs, fmt.Sprintf("%s(%s)", v.Name, v.Gender))
		}

		test.Voices[l.Name] = vs
	}

	data, err := json.MarshalIndent(test, "", "\t")
	if err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile("langs.json", data, 0666); err != nil {
		t.Fatal(err)
	}
}
