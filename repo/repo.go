package repo

var NTStatus *NTStatusRepo
var HResult *HResultRepo

func Load() error {
	var err error

	NTStatus, err = LoadNTStatuses("yaml/ntstatus.yml")

	if err != nil {
		return err
	}

	HResult, err = LoadHResults("yaml/hresult.yml")

	if err != nil {
		return err
	}

	return nil
}
