package commands

import tempest "github.com/amatsagu/tempest"

var HelloCommand = tempest.Command{
	Name:                "hello",
	Description:         "Hello World!",
	SlashCommandHandler: handleHello,
}

func handleHello(itx *tempest.CommandInteraction) {
	itx.SendReply(tempest.ResponseMessageData{
		Content: "Hello World!",
	}, false, nil)
}
