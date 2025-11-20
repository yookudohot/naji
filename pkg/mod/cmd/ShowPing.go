package cmd
import (
	client "github.com/bwmarrin/discordgo"
	"github.com/yookudohot/naji/handler"
)

func init(){
	handler.RegisterCommand(handler.Command{
		Name:	"ping",
		Description: "Responde com pong!",
		Execute: executePing,
	})
}

func executePing(s *client.Session, i *client.InteractionCreate) { s.InteractionRespond(i.Interaction, &client.InteractionResponse{
	Type: client.InteractionResponseChannelMessageWithSource,
	Data: &client.InteractionResponseData{
		Content: "🏓 Pong!",
	},
})

}
