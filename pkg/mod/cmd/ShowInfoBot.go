package cmd

import (
	client "github.com/bwmarrin/discordgo"
	"github.com/yookudohot/naji/handler"
)
func init(){
	handler.RegisterCommand(handler.Command{
		Name:	"info",
		Description: "Mostra informacoes do bot",
		Execute: executeInfo,
	})
}
func executeInfo(s *client.Session, i *client.InteractionCreate){
	embed := &client.MessageEmbed{
		Title:	"[🤖] Informações do bot!",
		Description: "Olá! Sou um bot experimental escrito em Golang!",
		Color:	0x00ff00,
	}
	s.InteractionRespond(i.Interaction, &client.InteractionResponse{
		Type: client.InteractionResponseChannelMessageWithSource,
		Data:	&client.InteractionResponseData{
			Embeds: []*client.MessageEmbed{embed},
		},
	})
}