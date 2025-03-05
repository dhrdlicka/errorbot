package repo

import (
	"os"

	"gopkg.in/yaml.v3"
)

type HResultRepo struct {
	Facilities map[uint16]string
	Codes      []ErrorInfo
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

func (repo HResultRepo) FindHResult(code uint32) []ErrorInfo {
	matches := []ErrorInfo{}

	for _, item := range repo.Codes {
		if item.Code == code {
			matches = append(matches, item)
		}
	}

	return matches
}
