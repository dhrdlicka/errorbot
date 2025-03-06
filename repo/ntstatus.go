package repo

import (
	"os"

	"gopkg.in/yaml.v3"
)

type NTStatusRepo struct {
	Facilities map[uint16]string
	Codes      []ErrorInfo
}

func LoadNTStatuses(name string) (NTStatusRepo, error) {
	file, err := os.ReadFile(name)

	if err != nil {
		return NTStatusRepo{}, err
	}

	var repo NTStatusRepo
	err = yaml.Unmarshal(file, &repo)

	if err != nil {
		return NTStatusRepo{}, err
	}

	return repo, nil
}

func (repo NTStatusRepo) FindCode(code uint32) []ErrorInfo {
	return FindCode(repo.Codes, code)
}

func (repo Repo) FindNTStatus(code uint32) []ErrorInfo {
	return repo.NTStatus.FindCode(code)
}
