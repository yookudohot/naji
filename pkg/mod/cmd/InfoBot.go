package cmd

import (
	time "time"
	fmt "fmt"
	client "github.com/bwmarrin/discordgo"
	handler "github.com/yookudohot/naji/handler"
	setting "github.com/yookudohot/naji/config/setting"
)

func init(){
	handler.RegisterCommand(handler.Command{
		Name: "showinfo",
		Description: "Mostra algumas informações",
		Execute: executeInfo,
	})
}

func executeInfo(s *client.Session, i *client.InteractionCreate){

	infoServers := len(s.State.Guilds)
	infoCommands := len(handler.Commands)


	embed := &client.MessageEmbed{
		Title: "[ `📄` ] Informações do bot.",
		Description: fmt.Sprintln("--------------"),
		Color: 0x808080,
		Fields: []*client.MessageEmbedField{
			{
				Name: "[`👷`] Criador",
				Value: "[Yookudohot](https://github.com/yookudohot)",
				Inline: false,
			},
			{
				Name: "[`📔`] Sobre mim",
				Value: fmt.Sprintf("* Estou em **%d** servidores\n * Tenho **%d** comandos registrados\n * Fui criado para ajudar na moderação de servidores e estou em constante evolução! xD", infoServers, infoCommands),
				Inline: false,
			},
			{
				Name: "[`🔗`] Links úteis",
				Value: fmt.Sprintf("[Convide-me, já!](%s)\n[Página no Github](%s)", setting.InviteLink_return(), setting.GithubReturn()),
				Inline: false,
			},
		},
	Timestamp: time.Now().Format(time.RFC3339),
}
s.InteractionRespond(i.Interaction, &client.InteractionResponse{
	Type: client.InteractionResponseChannelMessageWithSource,
	Data: &client.InteractionResponseData{
		Embeds: []*client.MessageEmbed{embed},
	},
})
}