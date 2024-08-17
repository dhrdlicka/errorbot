package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
)

var (
	applicationID string
	botToken      string
)

func main() {
	http.HandleFunc("/interactions", handleDiscordInteraction)

	port := os.Getenv("PORT")
	applicationID = os.Getenv("DISCORD_APPLICATION_ID")
	botToken = os.Getenv("DISCORD_BOT_TOKEN")

	err := registerDiscordCommands()

	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func registerDiscordCommands() error {
	commands := []discordgo.ApplicationCommand{
		{
			Name:        "hello",
			Description: "Hello World",
			Type:        discordgo.ChatApplicationCommand,
		},
	}

	globalCommandsEndpoint := discordgo.EndpointApplicationGlobalCommands(applicationID)

	body, err := json.Marshal(commands)

	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPut, globalCommandsEndpoint, bytes.NewBuffer(body))

	if err != nil {
		return err
	}

	request.Header.Add("authorization", "Bot "+botToken)

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode >= 400 {
		return err
	}

	return nil
}

func handleDiscordInteraction(w http.ResponseWriter, r *http.Request) {
	var interaction discordgo.Interaction
	err := json.NewDecoder(r.Body).Decode(&interaction)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var response discordgo.InteractionResponse

	switch interaction.Type {
	case discordgo.InteractionPing:
		response = discordgo.InteractionResponse{
			Type: discordgo.InteractionResponsePong,
		}
	case discordgo.InteractionApplicationCommand:
		data := interaction.ApplicationCommandData()

		switch data.Name {
		case "hello":
			response = discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "hello!",
				},
			}
		}
	}

	if response.Type == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}
