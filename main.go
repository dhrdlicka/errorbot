package main

import (
	"log"
	"net/http"
	"os"

	tempest "github.com/Amatsagu/Tempest"
	"github.com/dhrdlicka/errorbot/commands"
	"github.com/dhrdlicka/errorbot/repo"
)

func main() {
	err := repo.Load()

	if err != nil {
		log.Fatal(err)
	}

	client := tempest.NewClient(tempest.ClientOptions{
		PublicKey: os.Getenv("DISCORD_PUBLIC_KEY"),
		Rest:      tempest.NewRestClient(os.Getenv("DISCORD_BOT_TOKEN")),
	})

	//client.RegisterCommand(commands.HelloCommand)
	client.RegisterCommand(commands.ErrorCommand)
	client.RegisterCommand(commands.BugCheckCommand)
	client.RegisterCommand(commands.NTStatusCommand)

	err = client.SyncCommands(nil, nil, false)

	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("POST /interactions", client.HandleDiscordRequest)

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
