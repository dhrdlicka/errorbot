package commands

import (
	"fmt"
	"log/slog"

	tempest "github.com/Amatsagu/Tempest"
	"github.com/dhrdlicka/errorbot/repo"
	"github.com/dhrdlicka/errorbot/winerror"
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
	codes, err := parseCode(value)

	if err != nil {
		slog.Error("failed to parse command option", err)
		return
	}

	matches := []repo.ErrorInfo{}

	for _, code := range codes {
		matches = append(matches, repoInstance.FindNTStatus(code)...)
	}

	var response tempest.ResponseMessageData

	if len(matches) > 0 {
		for _, match := range matches {
			response.Embeds = append(response.Embeds, createNTStatusEmbed(match))
		}
	} else {
		// only break down the hexadecimal code if possible
		response.Embeds = append(response.Embeds, createUnknownNTStatusEmbed(codes[0]))
	}

	itx.SendReply(response, false, nil)
}

func ntStatusSeverityToString(severity uint8) string {
	switch severity {
	case winerror.STATUS_SEVERITY_SUCCESS:
		return "Success"
	case winerror.STATUS_SEVERITY_INFORMATIONAL:
		return "Informational"
	case winerror.STATUS_SEVERITY_WARNING:
		return "Warning"
	case winerror.STATUS_SEVERITY_ERROR:
		return "Error"
	}

	return ""
}

func createNTStatusEmbed(ntStatus repo.ErrorInfo) *tempest.Embed {
	return &tempest.Embed{
		Title:       ntStatus.Name,
		Description: ntStatus.Description,
		Fields: append(
			[]*tempest.EmbedField{
				{
					Name:  "NTSTATUS code",
					Value: fmt.Sprintf("`0x%08X` (%d)", ntStatus.Code, ntStatus.Code),
				},
			}, createNTStatusEmbedFields(winerror.NTStatus(ntStatus.Code))...),
	}
}

func createUnknownNTStatusEmbed(code uint32) *tempest.Embed {
	return &tempest.Embed{
		Fields: append(
			[]*tempest.EmbedField{
				{
					Name:  "NTSTATUS code",
					Value: fmt.Sprintf("`0x%08X` (%d)", code, code),
				},
			}, createNTStatusEmbedFields(winerror.NTStatus(code))...),
	}
}

func createNTStatusEmbedFields(status winerror.NTStatus) []*tempest.EmbedField {
	facility := fmt.Sprintf("%d", status.Facility())

	if facility_name, ok := repoInstance.NTStatus.Facilities[status.Facility()]; ok {
		facility = fmt.Sprintf("%s (%s)", facility_name, facility)
	}

	return []*tempest.EmbedField{
		{
			Name:   "Severity",
			Value:  fmt.Sprintf("%s (%d)", ntStatusSeverityToString(status.Sev()), status.Sev()),
			Inline: true,
		},
		{
			Name:   "Customer",
			Value:  fmt.Sprintf("%t", status.C()),
			Inline: true,
		},
		{
			Name:   "Reserved (N)",
			Value:  fmt.Sprintf("%d", boolToInt(status.N())),
			Inline: true,
		},
		{
			Name:   "Facility",
			Value:  facility,
			Inline: true,
		},
		{
			Name:   "Code",
			Value:  fmt.Sprintf("%d", status.Code()),
			Inline: true,
		},
	}
}
