package commands

import (
	"fmt"
	"log/slog"
	"strconv"

	tempest "github.com/Amatsagu/Tempest"
	"github.com/dhrdlicka/errorbot/winerror"
)

var BugCheckCommand = tempest.Command{
	Type:        tempest.CHAT_INPUT_COMMAND_TYPE,
	Name:        "bugcheck",
	Description: "Look up a Windows NT bug check code",
	Options: []tempest.CommandOption{
		{
			Type:        tempest.STRING_OPTION_TYPE,
			Name:        "code",
			Description: "Bug check code",
			Required:    true,
		},
	},
	SlashCommandHandler: handleError,
}

func handleBugCheck(itx *tempest.CommandInteraction) {
	bugCheckList, err := winerror.LoadBugChecks("yaml/bugcheck.yml")

	if err != nil {
		slog.Error("failed to load bugcheck.yml", err)
		return
	}

	value := itx.Data.Options[0].Value.(string)
	longCode, err := strconv.ParseUint(value, 0, 32)

	code := uint32(longCode)

	if err != nil {
		intCode, err := strconv.ParseInt(value, 0, 32)

		if err != nil {
			slog.Error("failed to parse command option", err)
			return
		}

		code = uint32(intCode)
	}

	var response tempest.ResponseMessageData

	matches := winerror.FindBugCheck(uint32(code), bugCheckList)

	if len(matches) > 0 {
		response.Content = fmt.Sprintln(matches[0])
	} else {
		response.Content = fmt.Sprintf("Could not find bug check code %s (`0x%08X`)", value, code)
	}

	itx.SendReply(response, false, nil)
}
