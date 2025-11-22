package handler

import (
	"strings"
	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Name        string
	Description string
	Options     []*discordgo.ApplicationCommandOption
	Execute     func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

var Commands = make(map[string]Command)

func RegisterCommand(cmd Command) {
	Commands[cmd.Name] = cmd
}
var ButtonHandlers = make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))

func RegisterButtonHandler(prefix string, handler func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	ButtonHandlers[prefix] = handler
}

func HandleButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
	customID := i.MessageComponentData().CustomID
	
	for prefix, handler := range ButtonHandlers {
		if strings.HasPrefix(customID, prefix) {
			handler(s, i)
			return
		}
	}
	
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "❌ Botão desconhecido",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}