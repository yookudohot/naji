package handler

import (
	"github.com/bwmarrin/discordgo"
	//client "github.com/bwmarrin/discordgo"
)
type Command struct {
	Name 	string 
	Description string 
	Options	[]*discordgo.ApplicationCommandOption
	Execute 	func(s *discordgo.Session, i *discordgo.InteractionCreate)
}
var Commands = make(map[string]Command)
func RegisterCommand(cmd Command) {
	Commands[cmd.Name] = cmd
}