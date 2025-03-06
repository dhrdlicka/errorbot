package commands

import "github.com/dhrdlicka/errorbot/repo"

var repoInstance repo.Repo

func LoadRepo() error {
	var err error
	repoInstance, err = repo.Load()

	return err
}
