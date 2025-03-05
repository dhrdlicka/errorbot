package repo

import (
	"os"

	"github.com/dhrdlicka/errorbot/winerror"
	"gopkg.in/yaml.v3"
)

type HResultDetails struct {
	Code        winerror.HResult `yaml:"code"`
	Name        string           `yaml:"name"`
	Description string           `yaml:"description"`
}

type HResultRepo struct {
	Facilities map[uint16]string
	Codes      []HResultDetails
}

func LoadHResults(name string) (*HResultRepo, error) {
	file, err := os.ReadFile(name)

	if err != nil {
		return nil, err
	}

	var repo HResultRepo
	err = yaml.Unmarshal(file, &repo)

	if err != nil {
		return nil, err
	}

	return &repo, nil
}

func (repo HResultRepo) FindHResult(code uint32) []HResultDetails {
	matches := []HResultDetails{}

	for _, item := range repo.Codes {
		if item.Code == winerror.HResult(code) {
			matches = append(matches, item)
		}
	}

	return matches
}
