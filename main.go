package main

import (
	"log"
	"net/http"
	"os"

	tempest "github.com/amatsagu/tempest"
	"github.com/dhrdlicka/errorbot/commands"
)

func main() {
	err := commands.LoadRepo()

	if err != nil {
		log.Fatal(err)
	}

	client := tempest.NewClient(tempest.ClientOptions{
		PublicKey: os.Getenv("DISCORD_PUBLIC_KEY"),
		Token:     os.Getenv("DISCORD_BOT_TOKEN"),
	})

	//client.RegisterCommand(commands.HelloCommand)
	client.RegisterCommand(commands.ErrorCommand)
	client.RegisterCommand(commands.BugCheckCommand)
	client.RegisterCommand(commands.NTStatusCommand)
	client.RegisterCommand(commands.HResultCommand)

	err = client.SyncCommandsWithDiscord(nil, nil, false)

	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("POST /interactions", client.DiscordRequestHandler)

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
