package commands

import (
	"fmt"
	"log/slog"

	tempest "github.com/Amatsagu/Tempest"
	"github.com/dhrdlicka/errorbot/repo"
	"github.com/dhrdlicka/errorbot/winerror"
)

var HResultCommand = tempest.Command{
	Type:        tempest.CHAT_INPUT_COMMAND_TYPE,
	Name:        "hresult",
	Description: "Look up a HRESULT error code",
	Options: []tempest.CommandOption{
		{
			Type:        tempest.STRING_OPTION_TYPE,
			Name:        "code",
			Description: "HRESULT code",
			Required:    true,
		},
	},
	SlashCommandHandler: handleHResult,
}

func handleHResult(itx *tempest.CommandInteraction) {
	value := itx.Data.Options[0].Value.(string)
	codes, err := ParseCode(value)

	if err != nil {
		slog.Error("failed to parse command option", err)
		return
	}

	matches := []repo.ErrorInfo{}

	for _, code := range codes {
		hr := winerror.HResult(code)

		if hr.N() {
			// this is a mapped NTSTATUS
			ntStatusMatches := repo.NTStatus.FindNTStatus(uint32(hr) ^ winerror.FACILITY_NT_BIT)

			for _, ntStatus := range ntStatusMatches {
				ntStatus.Name = fmt.Sprintf("HRESULT_FROM_NT(%s)", ntStatus.Name)
			}

			matches = append(matches, ntStatusMatches...)
		} else {
			matches = append(matches, repo.HResult.FindHResult(code)...)
		}
	}

	var response tempest.ResponseMessageData

	if len(matches) > 0 {
		for _, match := range matches {
			response.Embeds = append(response.Embeds, createHResultEmbed(match))
		}
	} else {
		response.Content = fmt.Sprintf("Could not find HRESULT code %s (`0x%08X`)", value, codes[0])
	}

	itx.SendReply(response, false, nil)
}

func hResultSeverityToString(severity bool) string {
	if severity {
		return "Failure"
	} else {
		return "Success"
	}
}

func createHResultEmbed(hResult repo.ErrorInfo) *tempest.Embed {
	return &tempest.Embed{
		Title:       hResult.Name,
		Description: hResult.Description,
		Fields: append(
			[]*tempest.EmbedField{
				{
					Name:  "HRESULT code",
					Value: fmt.Sprintf("`0x%08X` (%d)", hResult.Code, hResult.Code),
				},
			}, createHResultEmbedFields(winerror.HResult(hResult.Code))...),
	}
}

func createHResultEmbedFields(hResult winerror.HResult) []*tempest.EmbedField {
	if hResult.N() {
		// mapped NTSTATUS
		return createNTStatusEmbedFields(winerror.NTStatus(hResult))
	}

	facility := fmt.Sprintf("%d", hResult.Facility())

	if facility_name, ok := repo.HResult.Facilities[hResult.Facility()]; ok {
		facility = fmt.Sprintf("%s (%s)", facility_name, facility)
	}

	return []*tempest.EmbedField{
		{
			Name:   "Severity",
			Value:  fmt.Sprintf("%s (%d)", hResultSeverityToString(hResult.S()), boolToInt(hResult.S())),
			Inline: true,
		},
		{
			Name:   "Reserved (R)",
			Value:  fmt.Sprintf("%d", boolToInt(hResult.R())),
			Inline: true,
		},
		{
			Name:   "Customer",
			Value:  fmt.Sprintf("%t", hResult.C()),
			Inline: true,
		},
		{
			Name:   "Reserved (N)",
			Value:  fmt.Sprintf("%d", boolToInt(hResult.N())),
			Inline: true,
		},
		{
			Name:   "Reserved (X)",
			Value:  fmt.Sprintf("%d", boolToInt(hResult.X())),
			Inline: true,
		},
		{
			Name:   "Facility",
			Value:  facility,
			Inline: true,
		},
		{
			Name:   "Code",
			Value:  fmt.Sprintf("%d", hResult.Code()),
			Inline: true,
		},
	}
}
