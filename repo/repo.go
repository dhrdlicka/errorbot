package repo

var NTStatus *NTStatusRepo

func Load() error {
	repo, err := LoadNTStatuses("yaml/ntstatus.yml")

	NTStatus = repo

	return err
}
