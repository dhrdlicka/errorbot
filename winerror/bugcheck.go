package winerror

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type BugCheck struct {
	Code uint32 `yaml:"code"`
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

func LoadBugChecks(name string) ([]BugCheck, error) {
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

func FindBugCheck(code uint32, bugChecks []BugCheck) []BugCheck {
	matches := []BugCheck{}

	for _, item := range bugChecks {
		if item.Code == code {
			matches = append(matches, item)
		}
	}

	return matches
}

func (bugCheck BugCheck) String() string {
	return fmt.Sprintf("[`%s`](%s) (`0x%08X`)\n", bugCheck.Name, bugCheck.URL, bugCheck.Code)
}
