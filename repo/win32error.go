package repo

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Win32ErrorRepo Errors

func LoadWin32Errors(name string) (Win32ErrorRepo, error) {
	file, err := os.ReadFile(name)

	if err != nil {
		return nil, err
	}

	var repo Win32ErrorRepo
	err = yaml.Unmarshal(file, &repo)

	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (repo Win32ErrorRepo) FindCode(code uint32) []ErrorInfo {
	return Errors(repo).FindCode(code)
}

func (repo Repo) FindWin32Error(code uint32) []ErrorInfo {
	return repo.Win32Error.FindCode(code)
}
