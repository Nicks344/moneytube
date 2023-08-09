package moneytubemodel

type AudioGenerateInput struct {
	Api                  int
	ApiKey               string
	CropAudio            bool
	CropAudioVariant     int
	MaxTime              int
	Voice                string
	Lang                 string
	Speed                int
	Text                 string
	GSpeechProps         GSpeechProps
	SpeechProCredentials SpeechProCredentials
}

type GSpeechProps struct {
	Profiles []string
	Pitch    int
}

type SpeechProCredentials struct {
	ID       int
	Login    string
	Password string
}
