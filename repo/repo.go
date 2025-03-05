package repo

var NTStatus *NTStatusRepo
var HResult *HResultRepo
var Win32Error *Win32ErrorRepo

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

	Win32Error, err = LoadWin32Errors("yaml/win32error.yml")

	if err != nil {
		return err
	}

	return nil
}
