package main

import tempest "github.com/Amatsagu/Tempest"

var HelloCommand tempest.Command = tempest.Command{
	Name:                "hello",
	Description:         "Hello World!",
	SlashCommandHandler: handleHello,
}

func handleHello(itx *tempest.CommandInteraction) {
	itx.SendReply(tempest.ResponseMessageData{
		Content: "Hello World!",
	}, false, nil)
}
