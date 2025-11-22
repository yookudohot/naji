package config 
import (
	"log"
	"os"
	get_env "github.com/joho/godotenv"

)

func Token_return() string {
	err := get_env.Load("config/setting/.env")
	if err != nil {
		log.Println("Não foi possivel carregar o .env")
	}
	botToken := os.Getenv("BOT_TOKEN")
	return botToken
}
func Prefix_return() string {
	err := get_env.Load("config/setting/.env")
	if err != nil {
		log.Println("Não foi possivel carregar o .env")
	}
	botPrefix := os.Getenv("PREFIX")
	return botPrefix
}
func InviteLink_return() string {err := get_env.Load("config/setting/.env");if err != nil {
		log.Println("Nao foi possivel carregar o .env")
};botInviteLink := os.Getenv("INVITE_LINK"); return botInviteLink}

func GithubReturn() string {err := get_env.Load("config/setting/.env");if err != nil {
		log.Println("Nao foi possivel carregar o .env")
};github := os.Getenv("GITHUB_REPO"); return github}