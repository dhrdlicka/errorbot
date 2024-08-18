package commands

import (
	"fmt"
	"log/slog"
	"strconv"

	tempest "github.com/Amatsagu/Tempest"
	"github.com/dhrdlicka/errorbot/winerror"
)

var ErrorCommand = tempest.Command{
	Type:        tempest.CHAT_INPUT_COMMAND_TYPE,
	Name:        "error",
	Description: "Look up a Windows error code",
	Options: []tempest.CommandOption{
		{
			Type:        tempest.STRING_OPTION_TYPE,
			Name:        "code",
			Description: "Error code",
			Required:    true,
		},
	},
	SlashCommandHandler: handleError,
}

func handleError(itx *tempest.CommandInteraction) {
	ntStatusList, err := winerror.LoadErrorInfo("json/ntstatus.json")

	if err != nil {
		slog.Error("failed to load ntstatus.json", err)
		return
	}

	win32ErrorList, err := winerror.LoadErrorInfo("json/win32error.json")

	if err != nil {
		slog.Error("failed to load win32error.json", err)
		return
	}

	hResultList, err := winerror.LoadErrorInfo("json/hresult.json")

	if err != nil {
		slog.Error("failed to load hresult.json", err)
		return
	}

	value := itx.Data.Options[0].Value.(string)
	code, err := strconv.ParseUint(value, 0, 32)

	if err != nil {
		intCode, err := strconv.ParseInt(value, 0, 32)

		if err != nil {
			slog.Error("failed to parse command option", err)
			return
		}

		code = uint64(intCode)
	}

	var response tempest.ResponseMessageData

	matches := winerror.FindErrorCode(uint32(code), ntStatusList)

	if len(matches) > 0 {
		embed := tempest.Embed{
			Title: "Possible NTSTATUS codes",
		}

		var description []byte

		for _, item := range matches {
			description = fmt.Appendf(description, "%v\n", item)
		}

		embed.Description = string(description)

		response.Embeds = append(response.Embeds, &embed)
	}

	matches = winerror.FindErrorCode(uint32(code), win32ErrorList)

	if len(matches) > 0 {
		embed := tempest.Embed{
			Title: "Possible Win32 error codes",
		}

		var description []byte

		for _, item := range matches {
			description = fmt.Appendf(description, "%v\n", item)
		}

		embed.Description = string(description)

		response.Embeds = append(response.Embeds, &embed)
	}

	matches = winerror.FindErrorCode(uint32(code), hResultList)

	if len(matches) > 0 {
		embed := tempest.Embed{
			Title: "Possible HRESULT codes",
		}

		var description []byte

		for _, item := range matches {
			description = fmt.Appendf(description, "%v\n", item)
		}

		embed.Description = string(description)

		response.Embeds = append(response.Embeds, &embed)
	}

	itx.SendReply(response, false, nil)
}
