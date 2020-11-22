package discord

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type DISCORDFUNC func(*discordgo.Session, *discordgo.MessageCreate, []string)
type Command struct {
	Description string
	Exec        DISCORDFUNC
}

type CMD map[string]Command

func (myContext *MyContext) command(s *discordgo.Session, m *discordgo.MessageCreate) {
	cmds := myContext.getcommand()
	slice := strings.Split(m.Content, " ")

	cmd, ok := cmds[slice[0]]
	if !ok {
		if _, err := s.ChannelMessageSend(m.ChannelID, "command \""+slice[0]+"\" was not found"); err != nil {
			log.Println(err)
		}
		return
	}
	cmd.Exec(s, m, slice)
}
