package repo

import (
	"os"

	"gopkg.in/yaml.v3"
)

type NTStatusRepo struct {
	Facilities map[uint16]string
	Codes      []ErrorInfo
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

func (repo NTStatusRepo) FindNTStatus(code uint32) []ErrorInfo {
	matches := []ErrorInfo{}

	for _, item := range repo.Codes {
		if item.Code == code {
			matches = append(matches, item)
		}
	}

	return matches
}
