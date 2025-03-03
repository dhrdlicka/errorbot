package repo

import (
	"os"

	"github.com/dhrdlicka/errorbot/winerror"
	"gopkg.in/yaml.v3"
)

type NTStatusDetails struct {
	Code        winerror.NTStatus `yaml:"code"`
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
}

type NTStatusRepo struct {
	Facilities map[uint16]string
	Severities map[uint8]string
	Codes      []NTStatusDetails
}

func LoadNTStatuses(name string) (*NTStatusRepo, error) {
	file, err := os.ReadFile(name)

	if err != nil {
		return nil, err
	}

	var repo NTStatusRepo
	err = yaml.Unmarshal(file, &repo)

	if err != nil {
		return nil, err
	}

	return &repo, nil
}

func (repo NTStatusRepo) FindNTStatus(code uint32) []NTStatusDetails {
	matches := []NTStatusDetails{}

	for _, item := range repo.Codes {
		if item.Code == winerror.NTStatus(code) {
			matches = append(matches, item)
		}
	}

	return matches
}
