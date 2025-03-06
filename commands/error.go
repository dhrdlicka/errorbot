package commands

import (
	"fmt"
	"log/slog"

	tempest "github.com/Amatsagu/Tempest"
	"github.com/dhrdlicka/errorbot/repo"
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
	value := itx.Data.Options[0].Value.(string)
	codes, err := parseCode(value)

	if err != nil {
		slog.Error("failed to parse command option", err)
		return
	}

	hResultMatches := []repo.ErrorInfo{}
	win32ErrorMatches := []repo.ErrorInfo{}
	ntStatusMatches := []repo.ErrorInfo{}

	for _, code := range codes {
		hResultMatches = append(hResultMatches, repoInstance.FindHResult(code)...)
		win32ErrorMatches = append(win32ErrorMatches, repoInstance.FindWin32Error(code)...)
		ntStatusMatches = append(ntStatusMatches, repoInstance.FindNTStatus(code)...)
	}

	var response tempest.ResponseMessageData

	var hasAnyMatch = false

	if len(hResultMatches) > 0 {
		response.Embeds = append(response.Embeds, &tempest.Embed{
			Title:       "Possible HRESULT error codes",
			Description: formatResults(hResultMatches),
		})

		hasAnyMatch = true
	}

	if len(win32ErrorMatches) > 0 {
		response.Embeds = append(response.Embeds, &tempest.Embed{
			Title:       "Possible Win32 error codes",
			Description: formatResults(win32ErrorMatches),
		})

		hasAnyMatch = true
	}

	if len(ntStatusMatches) > 0 {
		response.Embeds = append(response.Embeds, &tempest.Embed{
			Title:       "Possible NTSTATUS error codes",
			Description: formatResults(ntStatusMatches),
		})

		hasAnyMatch = true
	}

	if !hasAnyMatch {
		response.Content = fmt.Sprintf("Could not find error code %s (`0x%08X`)", value, codes[0])
	}

	itx.SendReply(response, false, nil)
}

func formatResults(errors []repo.ErrorInfo) string {
	var result []byte

	for _, item := range errors {
		result = fmt.Appendf(result, "%v\n", item)
	}

	return string(result)
}
