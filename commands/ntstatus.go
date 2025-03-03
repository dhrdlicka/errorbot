package commands

import (
	"fmt"
	"log/slog"

	tempest "github.com/Amatsagu/Tempest"
	"github.com/dhrdlicka/errorbot/repo"
)

var NTStatusCommand = tempest.Command{
	Type:        tempest.CHAT_INPUT_COMMAND_TYPE,
	Name:        "ntstatus",
	Description: "Look up an NTSTATUS error code",
	Options: []tempest.CommandOption{
		{
			Type:        tempest.STRING_OPTION_TYPE,
			Name:        "code",
			Description: "NTSTATUS code",
			Required:    true,
		},
	},
	SlashCommandHandler: handleNTStatus,
}

func handleNTStatus(itx *tempest.CommandInteraction) {
	value := itx.Data.Options[0].Value.(string)
	codes, err := ParseCode(value)

	if err != nil {
		slog.Error("failed to parse command option", err)
		return
	}

	matches := []repo.NTStatusDetails{}

	for _, code := range codes {
		matches = append(matches, repo.NTStatus.FindNTStatus(code)...)
	}

	var response tempest.ResponseMessageData

	if len(matches) > 0 {
		for _, match := range matches {
			response.Embeds = append(response.Embeds, createEmbed(match))
		}
	} else {
		response.Content = fmt.Sprintf("Could not find NTSTATUS code %s (`0x%08X`)", value, codes[0])
	}

	itx.SendReply(response, false, nil)
}

func createEmbed(ntStatus repo.NTStatusDetails) *tempest.Embed {
	return &tempest.Embed{
		Title:       ntStatus.Name,
		Description: ntStatus.Description,
		Fields: []*tempest.EmbedField{
			{
				Name:  "NTSTATUS code",
				Value: fmt.Sprintf("`0x%08X` (%d)", ntStatus.Code, ntStatus.Code),
			},
			{
				Name:   "Sev",
				Value:  fmt.Sprintf("%d", ntStatus.Code.Sev()),
				Inline: true,
			},
			{
				Name:   "C",
				Value:  fmt.Sprintf("%t", ntStatus.Code.C()),
				Inline: true,
			},
			{
				Name:   "N",
				Value:  fmt.Sprintf("%t", ntStatus.Code.N()),
				Inline: true,
			},
			{
				Name:   "Facility",
				Value:  fmt.Sprintf("%d", ntStatus.Code.Facility()),
				Inline: true,
			},
			{
				Name:   "Code",
				Value:  fmt.Sprintf("%d", ntStatus.Code.Code()),
				Inline: true,
			},
		},
	}
}
