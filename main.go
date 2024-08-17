package main

import (
	"log"
	"net/http"
	"os"

	tempest "github.com/Amatsagu/Tempest"
)

func main() {
	client := tempest.NewClient(tempest.ClientOptions{
		PublicKey: os.Getenv("DISCORD_PUBLIC_KEY"),
		Rest:      tempest.NewRestClient(os.Getenv("DISCORD_BOT_TOKEN")),
	})

	client.RegisterCommand(HelloCommand)

	err := client.SyncCommands(nil, nil, false)

	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/interactions", client.HandleDiscordRequest)

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
