package discord

import (
	"encoding/json"
	"errors"
	"github/ekusonarot/discord-bot-ChatAloud/textToSpeech"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func (myContext *MyContext) getcommand() CMD {
	cmd := make(CMD)
	cmd[".cmon"] = Command{
		`
		".cmon"
			Enter the voice channel and execute the command
		`,
		func(s *discordgo.Session, m *discordgo.MessageCreate, cmds []string) {
			vcID, err := findVoiceChannel(s, m)
			if err != nil {
				if _, err := s.ChannelMessageSend(m.ChannelID, "Voice Channel was not found"); err != nil {
					log.Println(err)
				}
			}

			_, err = s.ChannelVoiceJoin(m.GuildID, vcID, false, true)
			if err != nil {
				if _, ok := s.VoiceConnections[m.GuildID]; ok {
					_ = s.VoiceConnections[m.GuildID]
				} else {
					log.Fatal(err)
				}
			}
			myContext.DeadMan = map[string]map[string]bool{
				m.GuildID: map[string]bool{},
			}
		},
	}

	cmd[".bye"] = Command{
		`
		".bye"
			This bot leaves
		`,
		func(s *discordgo.Session, m *discordgo.MessageCreate, cmds []string) {
			vc, ok := s.VoiceConnections[m.GuildID]
			if ok {
				if err := vc.Disconnect(); err != nil {
					log.Fatal(err)
				}
			}
		},
	}

	cmd[".vcset"] = Command{
		`
		".vcset"
			Change audio settings
			-speaker: Change the speaker
				-speaker {1...15}
			-style: Change the tone
				-style {1...15}
			-rate: Change the Speech rate
				-rate {0.5...10.0}
			-vctype: Change the Voice type
				-vctype {0.5...2.0}
			exp. .vcset -speaker 3 -style 4 -rate 1.0
		`,
		func(s *discordgo.Session, m *discordgo.MessageCreate, cmds []string) {
			if err := voiceSettingChange(cmds[1:], s, m, myContext.VoiceSetting); err != nil {
				if _, err := s.ChannelMessageSend(m.ChannelID, "command \""+cmds[0]+"\" : please enter the argument"); err != nil {
					log.Println(err)
				}
			}
		},
	}

	cmd[".help"] = Command{
		`
		".help"
			Show help
		`,
		func(s *discordgo.Session, m *discordgo.MessageCreate, cmds []string) {
			if _, err := s.ChannelMessageSend(m.ChannelID, myContext.helpString()); err != nil {
				log.Println(err)
			}
		},
	}

	cmd[".m"] = Command{
		`
		".m"
			Mute everyone in the Voice Channel
		`,
		func(s *discordgo.Session, m *discordgo.MessageCreate, cmds []string) {
			vc, ok := s.VoiceConnections[m.GuildID]
			if !ok {
				if _, err := s.ChannelMessageSend(m.ChannelID, "This bot is not in any Voice Channel"); err != nil {
					log.Println(err)
				}
				return
			}
			usrIDs, err := findUsersIDInVoiceChannel(s, m, vc.ChannelID)
			if err != nil {
				log.Fatal(err)
			}
			for _, usrID := range usrIDs {
				if err := s.GuildMemberMute(m.GuildID, usrID, true); err != nil {
					log.Println(err)
				}
			}
		},
	}

	cmd[".nm"] = Command{
		`
		".nm"
			Unmute everyone in the Voice Channel
		`,
		func(s *discordgo.Session, m *discordgo.MessageCreate, cmds []string) {
			vc, ok := s.VoiceConnections[m.GuildID]
			if !ok {
				if _, err := s.ChannelMessageSend(m.ChannelID, "This bot is not in any Voice Channel"); err != nil {
					log.Println(err)
				}
				return
			}
			usrIDs, err := findUsersIDInVoiceChannel(s, m, vc.ChannelID)
			if err != nil {
				log.Fatal(err)
			}
			for _, usrID := range usrIDs {
				if err := s.GuildMemberMute(m.GuildID, usrID, false); err != nil {
					log.Println(err)
				}
			}
		},
	}

	cmd[".d"] = Command{
		`
		".d"
			Deaf everyone in the Voice Channel
		`,
		func(s *discordgo.Session, m *discordgo.MessageCreate, cmds []string) {
			vc, ok := s.VoiceConnections[m.GuildID]
			if !ok {
				if _, err := s.ChannelMessageSend(m.ChannelID, "This bot is not in any Voice Channel"); err != nil {
					log.Println(err)
				}
				return
			}
			usrIDs, err := findUsersIDInVoiceChannel(s, m, vc.ChannelID)
			if err != nil {
				log.Fatal(err)
			}
			for _, usrID := range usrIDs {
				if err := s.GuildMemberDeafen(m.GuildID, usrID, true); err != nil {
					log.Println(err)
				}
			}
		},
	}

	cmd[".nd"] = Command{
		`
		".nd"
			Undeaf everyone in the Voice Channel
		`,
		func(s *discordgo.Session, m *discordgo.MessageCreate, cmds []string) {
			vc, ok := s.VoiceConnections[m.GuildID]
			if !ok {
				if _, err := s.ChannelMessageSend(m.ChannelID, "This bot is not in any Voice Channel"); err != nil {
					log.Println(err)
				}
				return
			}
			usrIDs, err := findUsersIDInVoiceChannel(s, m, vc.ChannelID)
			if err != nil {
				log.Fatal(err)
			}
			for _, usrID := range usrIDs {
				if err := s.GuildMemberDeafen(m.GuildID, usrID, false); err != nil {
					log.Println(err)
				}
			}
		},
	}

	cmd[".k"] = Command{
		`
		".k"
			Set status killed
			exp> .k
			exp> .k {usr name}
		`,
		func(s *discordgo.Session, m *discordgo.MessageCreate, cmds []string) {
			_, ok := s.VoiceConnections[m.GuildID]
			if !ok {
				if _, err := s.ChannelMessageSend(m.ChannelID, "This bot is not in any Voice Channel"); err != nil {
					log.Println(err)
				}
				return
			}
			var usrID string
			if len(cmds) > 1 {
				var err error
				usrID, err = findUsersIDbyUserName(s, m, cmds[1])
				if err != nil {
					if _, err := s.ChannelMessageSend(m.ChannelID, "UserID was not found"); err != nil {
						log.Println(err)
					}
				}
			} else {
				usrID = m.Author.ID
			}
			myContext.DeadMan[m.GuildID][usrID] = true
			if myContext.IsMeeting[m.GuildID] {
				if err := s.GuildMemberDeafen(m.GuildID, usrID, false); err != nil {
					log.Println(err)
				}
				if err := s.GuildMemberMute(m.GuildID, usrID, true); err != nil {
					log.Println(err)
				}
			} else {
				if err := s.GuildMemberDeafen(m.GuildID, usrID, false); err != nil {
					log.Println(err)
				}
				if err := s.GuildMemberMute(m.GuildID, usrID, false); err != nil {
					log.Println(err)
				}
			}
		},
	}

	cmd[".nk"] = Command{
		`
		".nk"
			Set status unkilled
			exp> .nk
			exp> .nk {usr name}
		`,
		func(s *discordgo.Session, m *discordgo.MessageCreate, cmds []string) {
			_, ok := s.VoiceConnections[m.GuildID]
			if !ok {
				if _, err := s.ChannelMessageSend(m.ChannelID, "This bot is not in any Voice Channel"); err != nil {
					log.Println(err)
				}
				return
			}
			var usrID string
			if len(cmds) > 1 {
				var err error
				usrID, err = findUsersIDbyUserName(s, m, cmds[1])
				if err != nil {
					if _, err := s.ChannelMessageSend(m.ChannelID, "UserID was not found"); err != nil {
						log.Println(err)
					}
				}
			} else {
				usrID = m.Author.ID
			}
			myContext.DeadMan[m.GuildID][usrID] = false
			if myContext.IsMeeting[m.GuildID] {
				if err := s.GuildMemberDeafen(m.GuildID, usrID, false); err != nil {
					log.Println(err)
				}
				if err := s.GuildMemberMute(m.GuildID, usrID, false); err != nil {
					log.Println(err)
				}
			} else {
				if err := s.GuildMemberDeafen(m.GuildID, usrID, true); err != nil {
					log.Println(err)
				}
				if err := s.GuildMemberMute(m.GuildID, usrID, true); err != nil {
					log.Println(err)
				}
			}
		},
	}

	cmd[".ck"] = Command{
		`
		".ck"
			kill status all clear
		`,
		func(s *discordgo.Session, m *discordgo.MessageCreate, cmds []string) {
			_, ok := s.VoiceConnections[m.GuildID]
			if !ok {
				if _, err := s.ChannelMessageSend(m.ChannelID, "This bot is not in any Voice Channel"); err != nil {
					log.Println(err)
				}
				return
			}

			delete(myContext.DeadMan, m.GuildID)
			myContext.DeadMan = map[string]map[string]bool{
				m.GuildID: map[string]bool{},
			}
		},
	}

	cmd[".mtg"] = Command{
		`
		".mtg"
			Meeting start
		`,
		func(s *discordgo.Session, m *discordgo.MessageCreate, cmds []string) {
			vc, ok := s.VoiceConnections[m.GuildID]
			if !ok {
				if _, err := s.ChannelMessageSend(m.ChannelID, "This bot is not in any Voice Channel"); err != nil {
					log.Println(err)
				}
				return
			}
			myContext.IsMeeting[m.GuildID] = true

			UsrIDs, err := findUsersIDInVoiceChannel(s, m, vc.ChannelID)
			if err != nil {
				if _, err := s.ChannelMessageSend(m.ChannelID, "UserIDs was not found"); err != nil {
					log.Println(err)
				}
				return
			}
			for _, usrID := range UsrIDs {
				if myContext.DeadMan[m.GuildID][usrID] {
					if err := s.GuildMemberDeafen(m.GuildID, usrID, false); err != nil {
						log.Println(err)
					}
					if err := s.GuildMemberMute(m.GuildID, usrID, true); err != nil {
						log.Println(err)
					}
				} else {
					if err := s.GuildMemberDeafen(m.GuildID, usrID, false); err != nil {
						log.Println(err)
					}
					if err := s.GuildMemberMute(m.GuildID, usrID, false); err != nil {
						log.Println(err)
					}
				}
			}

		},
	}

	cmd[".fin"] = Command{
		`
		".fin"
			Meeting finish
		`,
		func(s *discordgo.Session, m *discordgo.MessageCreate, cmds []string) {
			vc, ok := s.VoiceConnections[m.GuildID]
			if !ok {
				if _, err := s.ChannelMessageSend(m.ChannelID, "This bot is not in any Voice Channel"); err != nil {
					log.Println(err)
				}
				return
			}
			myContext.IsMeeting[m.GuildID] = false

			UsrIDs, err := findUsersIDInVoiceChannel(s, m, vc.ChannelID)
			if err != nil {
				if _, err := s.ChannelMessageSend(m.ChannelID, "UserIDs was not found"); err != nil {
					log.Println(err)
				}
				return
			}
			for _, usrID := range UsrIDs {
				if myContext.DeadMan[m.GuildID][usrID] {
					if err := s.GuildMemberDeafen(m.GuildID, usrID, false); err != nil {
						log.Println(err)
					}
					if err := s.GuildMemberMute(m.GuildID, usrID, false); err != nil {
						log.Println(err)
					}
				} else {
					if err := s.GuildMemberDeafen(m.GuildID, usrID, true); err != nil {
						log.Println(err)
					}
					if err := s.GuildMemberMute(m.GuildID, usrID, true); err != nil {
						log.Println(err)
					}
				}
			}
		},
	}
	return cmd
}

func (myContext *MyContext) helpString() string {
	cmds := myContext.getcommand()
	var help string
	for _, cmd := range cmds {
		help += cmd.Description
	}
	return help
}

func findUsersIDbyUserName(s *discordgo.Session, m *discordgo.MessageCreate, aUsername string) (string, error) {
	for _, guild := range s.State.Guilds {
		for _, voiceState := range guild.VoiceStates {
			User, err := s.User(voiceState.UserID)
			if err != nil {
				return "", errors.New("UserID was not found")
			}
			if User.Username == aUsername {
				return User.ID, nil
			}
		}
	}
	return "", errors.New("UserID was not found")
}

func findUsersIDInVoiceChannel(s *discordgo.Session, m *discordgo.MessageCreate, VoiceChannelID string) ([]string, error) {
	usrs := make([]string, 0)
	for _, guild := range s.State.Guilds {
		for _, voiceState := range guild.VoiceStates {
			if voiceState.ChannelID == VoiceChannelID {
				usrs = append(usrs, voiceState.UserID)
			}
		}
	}
	return usrs, nil
}

func findVoiceChannel(s *discordgo.Session, m *discordgo.MessageCreate) (string, error) {
	for _, guild := range s.State.Guilds {
		for _, voiceState := range guild.VoiceStates {
			if voiceState.UserID == m.Author.ID {
				return voiceState.ChannelID, nil
			}
		}
	}

	return "", errors.New("channel id was not found")
}
func voiceSettingChange(slice []string, s *discordgo.Session, m *discordgo.MessageCreate, VoiceSetting map[string]*textToSpeech.VoiceSetting) error {
	if len(slice) == 0 {
		return errors.New("no argument")
	}

	var apiSetting_json []byte
	var err error
	apiSetting, ok := VoiceSetting[m.Author.ID]
	if !ok {
		apiSetting_json, err = ioutil.ReadFile("defaultVoice.json")
		if err != nil {
			log.Fatal(err)
		}
		if err := json.Unmarshal(apiSetting_json, &apiSetting); err != nil {
			log.Fatal(err)
		}
		VoiceSetting[m.Author.ID] = apiSetting
	}

	for i := 0; i < len(slice)/2; i += 2 {
		switch slice[i] {
		case "-speaker":
			t, _ := strconv.Atoi(slice[i+1])

			if (t <= 15) && (t >= 1) {
				VoiceSetting[m.Author.ID].SpeakerID = int(t)
				break
			}
			if _, err := s.ChannelMessageSend(m.ChannelID, "argument \""+slice[i]+"\" : 1 <= arg <= 15"); err != nil {
				log.Println(err)
			}
		case "-style":
			t, _ := strconv.Atoi(slice[i+1])

			if (t <= 15) && (t >= 1) {
				VoiceSetting[m.Author.ID].StyleID = int(t)
				break
			}
			if _, err := s.ChannelMessageSend(m.ChannelID, "argument \""+slice[i]+"\" : 1 <= arg <= 15"); err != nil {
				log.Println(err)
			}
		case "-rate":
			t, _ := strconv.ParseFloat(slice[i+1], 32)

			if (t <= 10.0) && (t >= 0.5) {
				VoiceSetting[m.Author.ID].SpeechRate = float32(t)
				break
			}
			if _, err := s.ChannelMessageSend(m.ChannelID, "argument \""+slice[i]+"\" : 0.5 <= arg <= 10.0"); err != nil {
				log.Println(err)
			}
		case "-vctype":
			t, _ := strconv.ParseFloat(slice[i+1], 32)

			if (t <= 2.0) && (t >= 0.5) {
				VoiceSetting[m.Author.ID].VoiceType = float32(t)
				break
			}
			if _, err := s.ChannelMessageSend(m.ChannelID, "argument \""+slice[i]+"\" : 0.5 <= arg <= 2.0"); err != nil {
				log.Println(err)
			}
		default:
			if _, err := s.ChannelMessageSend(m.ChannelID, "argument \""+slice[i]+"\" was not found"); err != nil {
				log.Println(err)
			}
		}
	}
	return nil
}
