package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	
	"github.com/bwmarrin/discordgo"
	logrus "github.com/sirupsen/logrus" // pacote para logs mais limpos
	config "github.com/yookudohot/naji/config/setting" // config
	"github.com/yookudohot/naji/handler" //handler
	_ "github.com/yookudohot/naji/pkg/mod/cmd" // pacote de comandos
	_ "github.com/yookudohot/naji/pkg/mod/admin"
)

func main() {
	sess, err := discordgo.New("Bot " + config.Token_return())
	if err != nil {
		log.Fatal(err)
	}
	
	sess.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages
	
	sess.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Bot logado como: %s", s.State.User.Username)
		
		
		for _, cmd := range handler.Commands {
			_, err := s.ApplicationCommandCreate(s.State.User.ID, "", &discordgo.ApplicationCommand{
				Name:        cmd.Name,
				Description: cmd.Description,
				Options:	cmd.Options,
			})
			if err != nil {
				log.Printf("Erro ao criar o comando %s :: %s", cmd.Name, err)
			} else {
				log.Printf("Comando %s registrado.", cmd.Name)
			}
		}
	})
	
	sess.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if cmd, ok := handler.Commands[i.ApplicationCommandData().Name]; ok {
			cmd.Execute(s, i)
		}
	})
	
	err = sess.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer sess.Close()
	
	logger := logrus.New()
	logger.Info("[INFO] Bot Online")
	
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	logger.Info("[INFO] Bot desligando...")
}