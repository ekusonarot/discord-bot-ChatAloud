package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github/ekusonarot/discord-bot-ChatAloud/discord"
	"github/ekusonarot/discord-bot-ChatAloud/textToSpeech"

	"github.com/bwmarrin/discordgo"
)

var stopBot = make(chan bool)

func main() {
	VoiceSetting := make(map[string]*textToSpeech.VoiceSetting)

	bytes, err := ioutil.ReadFile("setting.json")
	if err != nil {
		log.Fatal(err)
	}

	var setting discord.DiscordAPI
	if err := json.Unmarshal(bytes, &setting); err != nil {
		log.Fatal(err)
	}

	var docomoAPI textToSpeech.DocomoAPI
	if err := json.Unmarshal(bytes, &docomoAPI); err != nil {
		log.Fatal(err)
	}

	tDiscord, err := discordgo.New()
	if err != nil {
		log.Fatal(err)
	}
	tDiscord.Token = "Bot " + setting.Token

	tDiscord.AddHandler(discord.BotRespounse(VoiceSetting, docomoAPI))

	if err := tDiscord.Open(); err != nil {
		log.Fatal(err)
	}
	defer tDiscord.Close()

	log.Println("listening...")
	<-stopBot
}
