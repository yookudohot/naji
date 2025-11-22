package admin

import (
	"fmt"
	"strings"
	client "github.com/bwmarrin/discordgo"
	handler "github.com/yookudohot/naji/handler"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "ban",
		Description: "Bane um usuario do servidor",
		Execute:     executeBan,
		Options: []*client.ApplicationCommandOption{
			{
				Type:        client.ApplicationCommandOptionUser,
				Name:        "usuario",
				Description: "Usuario a ser banido",
				Required:    true,
			},
			{
				Type:        client.ApplicationCommandOptionString,
				Name:        "motivo",
				Description: "Motivo do ban",
				Required:    false,
			},
		},
	})
	
	handler.RegisterButtonHandler("ban_confirm", handleBanConfirm)
	handler.RegisterButtonHandler("ban_cancel", handleBanCancel)
}

func executeBan(s *client.Session, i *client.InteractionCreate) {
	if i.Member == nil {
		responderErro(s, i, "Erro: comando usado fora de um servidor")
		return
	}
	
	if !temPermissao(i.Member.Permissions) {
		responderErro(s, i, "Voce nao tem permissao")
		return
	}
	
	options := i.ApplicationCommandData().Options
	usuario := options[0].UserValue(s)
	
	motivo := "Sem motivo especificado"
	if len(options) > 1 {
		motivo = options[1].StringValue()
	}
	
	
	embed := &client.MessageEmbed{
		Title:       "⚠️ Confirmação de Ban",
		Description: fmt.Sprintf("Você tem certeza que deseja banir %s?", usuario.Mention()),
		Color:       0xe74c3c,
		Fields: []*client.MessageEmbedField{
			{
				Name:   "📝 Motivo",
				Value:  motivo,
				Inline: false,
			},
			{
				Name:   "⚠️ Atenção",
				Value:  "Esta ação é irreversível!",
				Inline: false,
			},
		},
	}
	
	components := []client.MessageComponent{
		client.ActionsRow{
			Components: []client.MessageComponent{
				client.Button{
					Label:    "✅ Confirmar Ban",
					Style:    client.DangerButton,
					CustomID: fmt.Sprintf("ban_confirm:%s:%s:%s", usuario.ID, i.Member.User.ID, motivo),
				},
				client.Button{
					Label:    "❌ Cancelar",
					Style:    client.SecondaryButton,
					CustomID: fmt.Sprintf("ban_cancel:%s:%s", usuario.ID, i.Member.User.ID),
				},
			},
		},
	}
	
	s.InteractionRespond(i.Interaction, &client.InteractionResponse{
		Type: client.InteractionResponseChannelMessageWithSource,
		Data: &client.InteractionResponseData{
			Embeds:     []*client.MessageEmbed{embed},
			Components: components,
			Flags:      client.MessageFlagsEphemeral,
		},
	})
}

func handleBanConfirm(s *client.Session, i *client.InteractionCreate) {
	customID := i.MessageComponentData().CustomID
	parts := strings.Split(customID, ":")
	
	if len(parts) < 3 {
		s.InteractionRespond(i.Interaction, &client.InteractionResponse{
			Type: client.InteractionResponseChannelMessageWithSource,
			Data: &client.InteractionResponseData{
				Content: "❌ Erro ao processar botão",
				Flags:   client.MessageFlagsEphemeral,
			},
		})
		return
	}
	
	usuarioID := parts[1]
	moderadorID := parts[2]
	motivo := "Sem motivo especificado"
	if len(parts) > 3 {
		motivo = parts[3]
	}
	
	if i.Member.User.ID != moderadorID {
		s.InteractionRespond(i.Interaction, &client.InteractionResponse{
			Type: client.InteractionResponseChannelMessageWithSource,
			Data: &client.InteractionResponseData{
				Content: "❌ Apenas quem usou o comando pode confirmar!",
				Flags:   client.MessageFlagsEphemeral,
			},
		})
		return
	}

	err := s.GuildBanCreateWithReason(i.GuildID, usuarioID, motivo, 0)
	if err != nil {
		s.InteractionRespond(i.Interaction, &client.InteractionResponse{
			Type: client.InteractionResponseUpdateMessage,
			Data: &client.InteractionResponseData{
				Content:    fmt.Sprintf("❌ Erro ao banir usuário: %v", err),
				Embeds:     []*client.MessageEmbed{},
				Components: []client.MessageComponent{},
			},
		})
		return
	}
	

	embed := &client.MessageEmbed{
		Title:       "🔨 Usuário Banido",
		Description: fmt.Sprintf("Usuário <@%s> foi banido com sucesso!", usuarioID),
		Color:       0xe74c3c,
		Fields: []*client.MessageEmbedField{
			{
				Name:   "📝 Motivo",
				Value:  motivo,
				Inline: false,
			},
			{
				Name:   "👮 Moderador",
				Value:  fmt.Sprintf("<@%s>", moderadorID),
				Inline: true,
			},
		},
	}
	
	s.InteractionRespond(i.Interaction, &client.InteractionResponse{
		Type: client.InteractionResponseUpdateMessage,
		Data: &client.InteractionResponseData{
			Embeds:     []*client.MessageEmbed{embed},
			Components: []client.MessageComponent{},
		},
	})
}

func handleBanCancel(s *client.Session, i *client.InteractionCreate) {
	customID := i.MessageComponentData().CustomID
	parts := strings.Split(customID, ":")
	
	if len(parts) < 3 {
		return
	}
	
	moderadorID := parts[2]
	
	if i.Member.User.ID != moderadorID {
		s.InteractionRespond(i.Interaction, &client.InteractionResponse{
			Type: client.InteractionResponseChannelMessageWithSource,
			Data: &client.InteractionResponseData{
				Content: "❌ Apenas quem usou o comando pode cancelar!",
				Flags:   client.MessageFlagsEphemeral,
			},
		})
		return
	}
	
	s.InteractionRespond(i.Interaction, &client.InteractionResponse{
		Type: client.InteractionResponseUpdateMessage,
		Data: &client.InteractionResponseData{
			Content:    "❌ Ban cancelado. Nenhuma ação foi tomada.",
			Embeds:     []*client.MessageEmbed{},
			Components: []client.MessageComponent{},
		},
	})
}