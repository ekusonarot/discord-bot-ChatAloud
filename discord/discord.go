package discord

import (
	"github/ekusonarot/discord-bot-ChatAloud/textToSpeech"
	"time"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
)

type DiscordAPI struct {
	Token   string `json:"DISCORD_TOKEN"`
	BotName string `json:"CLIENT_ID"`
}

type AudioFormat struct {
	headchunkId    string
	headchunkSize  int64
	formType       string
	fmtchunkID     string
	fmtchunkSize   int64
	waveFormatType int
	channel        int
	samplePerSec   int64
	bytePerSec     int
	blockSize      int
	bitsPerSample  int
	datachunkID    string
	datachunkSize  int64
	data           []byte
}

type MyContext struct {
	VoiceSetting map[string]*textToSpeech.VoiceSetting
	DocomoAPI    textToSpeech.DocomoAPI
}

func BotRespounse(VoiceSetting map[string]*textToSpeech.VoiceSetting, docomoAPI textToSpeech.DocomoAPI) func(*discordgo.Session, *discordgo.MessageCreate) {
	myContext := MyContext{VoiceSetting, docomoAPI}
	mySpeaking := make(map[string]bool)
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.Bot {
			return
		}

		cmd := []byte(m.Content)
		if cmd[0] == '.' {
			myContext.command(s, m)
			return
		}

		vc, ok := s.VoiceConnections[m.GuildID]
		if !ok {
			return
		}

		buffer, err := myContext.docomoAPIrequest(m)
		if err != nil {
			return
		}

		int16_slice := byteslice2int16slice(buffer)

		send := make(chan []int16, 2)
		go dgvoice.SendPCM(vc, send)

		for mySpeaking[vc.ChannelID] {
			time.Sleep(100 * time.Millisecond)
		}
		mySpeaking[vc.ChannelID] = true
		vc.Speaking(true)
		time.Sleep(300 * time.Millisecond)
		for i := 0; i < len(int16_slice); i += 960 {
			if i+960 > len(int16_slice) {
				send <- int16_slice[i:]
			} else {
				send <- int16_slice[i : i+960]
			}
		}
		time.Sleep(300 * time.Millisecond)
		vc.Speaking(false)
		mySpeaking[vc.ChannelID] = false
	}
}
