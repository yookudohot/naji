package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	logrus "github.com/sirupsen/logrus"
	config "github.com/yookudohot/naji/config/setting"
	"github.com/yookudohot/naji/handler"
	_ "github.com/yookudohot/naji/pkg/mod/admin"
	_ "github.com/yookudohot/naji/pkg/mod/cmd"

)

const GUILD_ID = "" 

func main() {
	sess, err := discordgo.New("Bot " + config.Token_return())
	if err != nil {
		log.Fatalf("Erro ao criar sessão do Discord: %v", err)
	}
	sess.Identify.Intents = discordgo.IntentsGuildMembers |
		discordgo.IntentsGuilds |
		discordgo.IntentsGuildMessages |
		discordgo.IntentsGuildIntegrations 

	sess.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Bot logado como: %s#%s", s.State.User.Username, s.State.User.Discriminator)
		log.Printf("Iniciando sincronização de comandos em: %s (ID da Guilda)", GUILD_ID)

		synchronizeCommands(s, s.State.User.ID, GUILD_ID)
	})

	sess.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			cmdName := i.ApplicationCommandData().Name
			if cmd, ok := handler.Commands[cmdName]; ok {
				cmd.Execute(s, i)
			} else {
				log.Printf("Comando não encontrado no handler: %s", cmdName)
			}
		case discordgo.InteractionMessageComponent:
			handler.HandleButton(s, i)
		}
	})

	err = sess.Open()
	if err != nil {
		log.Fatalf("Erro ao abrir a conexão: %v", err)
	}
	defer sess.Close()

	
	logger := logrus.New()
	logger.Info("[INFO] Bot Online. Pressione CTRL+C para desligar.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	logger.Info("[INFO] Bot desligando...")
}

func synchronizeCommands(s *discordgo.Session, appID string, guildID string) {

	log.Println("--- 🧹 FASE 1: Limpando comandos antigos... ---")
	oldCommands, err := s.ApplicationCommands(appID, guildID)
	if err != nil {
		log.Printf("❌ Erro ao buscar comandos antigos: %v", err)
		return
	}

	deletedCount := 0
	for _, cmd := range oldCommands {
		err := s.ApplicationCommandDelete(appID, guildID, cmd.ID)
		if err != nil {
			log.Printf("❌ Erro ao deletar '%s' (ID: %s): %v", cmd.Name, cmd.ID, err)
		} else {
			log.Printf("✅ Deletado: %s", cmd.Name)
			deletedCount++
		}
	
		time.Sleep(200 * time.Millisecond)
	}
	log.Printf("✅ Limpeza concluída. Total de comandos deletados: %d.", deletedCount)

	
	log.Println("--- 🚀 FASE 2: Registrando novos comandos... ---")
	registeredCount := 0
	for _, cmd := range handler.Commands {
		log.Printf("Registrando: %s", cmd.Name)

		_, err := s.ApplicationCommandCreate(appID, guildID, &discordgo.ApplicationCommand{
			Name:        cmd.Name,
			Description: cmd.Description,
			Options:     cmd.Options,
		})
		if err != nil {
			log.Printf("❌ Erro ao registrar '%s': %v", cmd.Name, err)
		} else {
			log.Printf("✅ Registrado: %s", cmd.Name)
			registeredCount++
		}

		time.Sleep(500 * time.Millisecond)
	}

	log.Printf("🎉 Sincronização COMPLETA. Total registrado: %d/%d comandos.", registeredCount, len(handler.Commands))
	log.Println("--------------------------------------------------")
}