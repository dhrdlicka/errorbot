package repo

import "fmt"

type ErrorInfo struct {
	Code        uint32 `yaml:"code"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

func FindCode(errors []ErrorInfo, code uint32) []ErrorInfo {
	matches := []ErrorInfo{}

	for _, item := range errors {
		if item.Code == code {
			matches = append(matches, item)
		}
	}

	return matches
}

func (errorInfo ErrorInfo) String() string {
	return fmt.Sprintf("`%s` (`0x%08X`)\n> %s\n", errorInfo.Name, errorInfo.Code, errorInfo.Description)
}
