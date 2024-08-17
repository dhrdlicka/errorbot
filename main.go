package main

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
)

var (
	applicationID        string
	applicationPublicKey ed25519.PublicKey
	botToken             string
)

func main() {
	var err error

	http.HandleFunc("/interactions", handleDiscordInteraction)

	port := os.Getenv("PORT")
	applicationID = os.Getenv("DISCORD_APPLICATION_ID")
	botToken = os.Getenv("DISCORD_BOT_TOKEN")

	applicationPublicKey, err = hex.DecodeString(os.Getenv("DISCORD_APPLICATION_PUBLIC_KEY"))

	if err != nil {
		log.Fatal(err)
	}

	err = registerDiscordCommands()

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
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !validateSignature(r) {
		w.WriteHeader(http.StatusUnauthorized)
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

func validateSignature(r *http.Request) bool {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		log.Println(err)
		return false
	}

	signature := []byte(r.Header.Get("x-signature-ed25519"))
	timestamp := r.Header.Get("x-signature-timestamp")

	payload := []byte(timestamp + string(body))

	return ed25519.Verify(applicationPublicKey, payload, signature)
}
