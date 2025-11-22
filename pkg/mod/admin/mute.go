package admin


import (
	"fmt"
	"time"

	client "github.com/bwmarrin/discordgo"
	handler "github.com/yookudohot/naji/handler"
)
func init(){
	handler.RegisterCommand(handler.Command{
		Name: "mute",
		Description: "Dando muito trabalho? Shiii",
		Execute: executeMute,
		Options: []*client.ApplicationCommandOption{
			{
				Type:	client.ApplicationCommandOptionUser,
				Name: "usuario",
				Description: "Usuario a ser silenciado",
				Required: true,
			},
			{
				Type:	client.ApplicationCommandOptionInteger,
				Name:	"duracao",
				Description: "Duracao do mute em minutis",
				Required: true,
				MinValue: func() *float64 { v:= 1.0; return &v}(),
				MaxValue:43200 ,
			},
			{
				Type:	client.ApplicationCommandOptionString,
				Name:	"motivo",
				Description: "Motivo do mute",
				Required: false,
			},
		},
	})
}

func executeMute(s* client.Session, i*client.InteractionCreate){
	if !temPermissao(i.Member.Permissions){
		responderErro(s, i, "Você não tem permissão.")
		return
	}
	options := i.ApplicationCommandData().Options
	usuario := options[0].UserValue(s)
	duracao := options[1].IntValue()

	motivo := "Sem motivo especificado..."
	if len(options) > 2 {
		motivo = options[2].StringValue()
	}
	timeout := time.Now().Add(time.Duration(duracao)*time.Minute)

	err := s.GuildMemberTimeout(i.GuildID, usuario.ID, &timeout)
	if err != nil {
		responderErro(s, i, fmt.Sprintf("Erro ao mutar o usuario: %v", err))
		return
	}
	embed := &client.MessageEmbed {
		Title: "[🔇]Usuario Punido!",
		Description: fmt.Sprintf("%s foi silenciado por **%d** minutos!", usuario.Mention(), duracao),
		Color:	0xe74c3c,
		Fields: []*client.MessageEmbedField{
			{
				Name:	"[📝] Motivo",
				Value: motivo,
				Inline: false,
			},
			{
				 Name: "⏰ Até",
				 Value: fmt.Sprintf("<t:%d:F>", timeout.Unix()),
				 Inline: false,
			},
			{
				Name: "[👮] Moderador",
				Value: i.Member.User.Mention(),
				Inline: true,
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




func temPermissao(perms int64) bool {
	return perms&client.PermissionModerateMembers != 0 ||
		perms&client.PermissionAdministrator != 0
}

func responderErro(s *client.Session, i *client.InteractionCreate, mensagem string) {
	s.InteractionRespond(i.Interaction, &client.InteractionResponse{
		Type: client.InteractionResponseChannelMessageWithSource,
		Data: &client.InteractionResponseData{
			Content: mensagem,
			Flags:   client.MessageFlagsEphemeral, // Mensagem só visível para quem usou
		},
	})
}
