package commands

import (
	"fmt"
	"log/slog"
	"strings"

	tempest "github.com/Amatsagu/Tempest"
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
	SlashCommandHandler: handleBugCheck,
}

func handleBugCheck(itx *tempest.CommandInteraction) {
	value := itx.Data.Options[0].Value.(string)
	codes, err := parseCode(value)

	if err != nil {
		slog.Error("failed to parse command option", "error", err)
		return
	}

	var response tempest.ResponseMessageData

	matches := repoInstance.BugCheck.FindBugCheckCode(codes[0])

	if len(matches) > 0 {
		match := matches[0]

		embed := tempest.Embed{
			Title:       match.Name,
			Description: match.Description,
			Fields: []*tempest.EmbedField{
				{
					Name:  "Bugcheck code",
					Value: fmt.Sprintf("`0x%08X`", match.Code),
				},
			},
		}

		if len(match.Parameters) > 0 {
			parameters := ""

			for i, parameter := range match.Parameters {
				parameters = fmt.Sprintf("%s%d. %s\n", parameters, i, strings.ReplaceAll(parameter, "\n", "\n   "))
			}

			if len(parameters) < 1024 {
				embed.Fields = append(embed.Fields, &tempest.EmbedField{
					Name:  "Parameters",
					Value: parameters,
				})
			}
		}

		embed.Fields = append(embed.Fields, &tempest.EmbedField{
			Name:  "Documentation",
			Value: match.URL,
		})

		response.Embeds = append(response.Embeds, &embed)

	} else {
		response.Content = fmt.Sprintf("Could not find bug check code %s (`0x%08X`)", value, codes[0])
	}

	itx.SendReply(response, false, nil)
}
