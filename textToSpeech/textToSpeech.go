package textToSpeech

type DocomoAPI struct {
	Token string `json:"DOCOMO_API_TOKEN"`
	Url   string `json:"URL"`
}

type VoiceSetting struct {
	Command         string  `json:"Command"`
	SpeakerID       int     `json:"SpeakerID"`
	StyleID         int     `json:"StyleID"`
	SpeechRate      float32 `json:"SpeechRate"`
	PowerRate       float32 `json:"PowerRate"`
	VoiceType       float32 `json:"VoiceType"`
	AudioFileFormat int     `json:"AudioFileFormat"`
	TextData        string  `json:"TextData"`
}
