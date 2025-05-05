package repo

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type BugCheck struct {
	Code        uint32   `yaml:"code"`
	Name        string   `yaml:"name"`
	URL         string   `yaml:"url"`
	Description string   `yaml:"description"`
	Parameters  []string `yaml:"parameters"`
}

type BugCheckRepo []BugCheck

func LoadBugChecks(name string) (BugCheckRepo, error) {
	file, err := os.ReadFile(name)

	if err != nil {
		return nil, err
	}

	var bugChecks []BugCheck
	err = yaml.Unmarshal(file, &bugChecks)

	if err != nil {
		return nil, err
	}

	return bugChecks, nil
}

func (repo Repo) FindBugCheck(code uint32) []ErrorInfo {
	return repo.BugCheck.FindCode(code)
}

func (bugChecks BugCheckRepo) FindCode(code uint32) []ErrorInfo {
	matches := []ErrorInfo{}

	for _, bugCheck := range bugChecks.FindBugCheckCode(code) {
		matches = append(matches, bugCheck.ErrorInfo())
	}

	return matches
}

func (bugChecks BugCheckRepo) FindBugCheckCode(code uint32) []BugCheck {
	matches := []BugCheck{}

	for _, item := range bugChecks {
		if item.Code == code {
			matches = append(matches, item)
		}
	}

	return matches
}

func (bugCheck BugCheck) String() string {
	description := fmt.Sprintf("`0x%08X`\n\n", bugCheck.Code)

	if bugCheck.Description != "" {
		description = fmt.Sprintf("%s#### Description\n%s\n\n", description, bugCheck.Description)
	}

	if len(bugCheck.Parameters) > 0 {
		description = fmt.Sprintf("%s#### Parameters\n\n", description)
		for i, item := range bugCheck.Parameters {
			description = fmt.Sprintf("%s%d. ", description, i+1)
			for j, line := range strings.Split(item, "\n") {
				if j > 0 {
					description = fmt.Sprintf("%s   ", description)
				}

				description = fmt.Sprintf("%s%s\n", description, line)
			}
		}
	}

	return description
}

func (bugCheck BugCheck) ErrorInfo() ErrorInfo {
	return ErrorInfo{
		Code:        bugCheck.Code,
		Name:        bugCheck.Name,
		Description: bugCheck.Description,
	}
}
