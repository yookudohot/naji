package admin

import (
	fmt "fmt"
	client "github.com/bwmarrin/discordgo"
	handler "github.com/yookudohot/naji/handler"
)

func init(){
	handler.RegisterCommand(handler.Command{
		Name: "unmute",
		Description: "desmuta alguém mutado",
		Execute: executeUnmute,
		Options: []*client.ApplicationCommandOption{
			{
				Type: client.ApplicationCommandOptionUser,
				Name: "usuario",
				Description: "usuario a ser desmutado",
				Required: true,
			},
			{
				Type: client.ApplicationCommandOptionString,
				Name: "motivo",
				Description: "Motivo do unmute",
				Required: false,
			},
		},
	})
}

func executeUnmute(s* client.Session, i * client.InteractionCreate){
	if i.Member == nil {
		responderErro(s, i, "O usuario não existe.")
		return 
	}
	if !temPermissao(i.Member.Permissions){
		responderErro(s, i, "Você não tem permissão.")
		return 
	}
	options := i.ApplicationCommandData().Options 
	usuario := options[0].UserValue(s)
	motivo := "Sem motivo especificado"

	if len(options) > 1 {
		motivo = options[1].StringValue()
	}
	err := s.GuildMemberTimeout(i.GuildID, usuario.ID, nil)
	if err != nil {
		responderErro(s, i, fmt.Sprintf("Erro! Impossivel desmutar o usuario %v", usuario.Mention()))
		return 
	}
	embed := &client.MessageEmbed{
		Title: "[🔊] Usuário Desmutado!",
		Description: fmt.Sprintf("%s teve o mute removido", usuario.Mention()),
		Color: 0x00000,
		Fields: []*client.MessageEmbedField{
			{
				Name: "[📝] Motivo",
				Value: motivo,
				Inline: false,
			},
			{
				Name: "[👮] Moderador",
				Value: i.Member.User.Mention(),
				Inline: true,
			},
		},
	}
	s.InteractionRespond(i.Interaction, &client.InteractionResponse{
		Type: client.InteractionResponseChannelMessageWithSource,
		Data: &client.InteractionResponseData{
			Embeds: []*client.MessageEmbed{embed},
		},
	})
}