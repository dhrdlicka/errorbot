package commands

import (
	"fmt"
	"strings"

	tempest "github.com/Amatsagu/Tempest"
	"github.com/dhrdlicka/errorbot/repo"
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
	AutoCompleteHandler: handleBugCheckAutoComplete,
}

func handleBugCheck(itx *tempest.CommandInteraction) {
	value := itx.Data.Options[0].Value.(string)
	codes, err := parseCode(value)

	var response tempest.ResponseMessageData
	var matches []repo.BugCheck

	if err != nil {
		matches = repoInstance.BugCheck.FindBugCheckString(value)
	} else {
		matches = repoInstance.BugCheck.FindBugCheckCode(codes[0])
	}

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

func handleBugCheckAutoComplete(itx tempest.CommandInteraction) []tempest.Choice {
	value := itx.Data.Options[0].Value.(string)

	matches := repoInstance.BugCheck.FindBugCheckString(value)

	var choices []tempest.Choice

	for _, match := range matches {
		choices = append(choices, tempest.Choice{
			Name:  match.Name,
			Value: fmt.Sprintf("0x%08X", match.Code),
		})
	}

	return choices
}
