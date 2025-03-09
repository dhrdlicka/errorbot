package commands

import (
	"fmt"
	"log/slog"

	tempest "github.com/Amatsagu/Tempest"
	"github.com/dhrdlicka/errorbot/repo"
)

type codeFinder interface {
	FindCode(uint32) []repo.ErrorInfo
}

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
		slog.Error("failed to parse command option", "error", err)
		return
	}

	repos := []struct {
		codeFinder
		string
	}{
		{repoInstance.BugCheck, "bug check"},
		{repoInstance.HResult, "HRESULT"},
		{repoInstance.Win32Error, "Win32 error"},
		{repoInstance.NTStatus, "NTSTATUS"},
	}

	var response tempest.ResponseMessageData

	found := false

	for _, errorRepo := range repos {
		matches := []repo.ErrorInfo{}

		for _, code := range codes {
			matches = append(matches, errorRepo.FindCode(code)...)
		}

		if len(matches) > 0 {
			found = true

			response.Embeds = append(response.Embeds, &tempest.Embed{
				Title:       fmt.Sprintf("Possible %s codes", errorRepo.string),
				Description: formatResults(matches),
			})
		}
	}

	if !found {
		response.Content = fmt.Sprintf("Could not find error code %s (`0x%08X`)", value, codes[0])
	}

	itx.SendReply(response, false, nil)
}

func formatResults(errors []repo.ErrorInfo) string {
	var result []byte

	for _, item := range errors {
		result = fmt.Appendf(result, "`%s` (`0x%08X`)\n> %s\n\n", item.Name, item.Code, item.Description)
	}

	return string(result)
}
