package repo

type Repo struct {
	NTStatus   NTStatusRepo
	HResult    HResultRepo
	Win32Error Win32ErrorRepo
	BugCheck   BugCheckRepo
}

func Load() (Repo, error) {
	var err error

	ntStatuses, err := LoadNTStatuses("ntstatus.yml")

	if err != nil {
		return Repo{}, err
	}

	hResults, err := LoadHResults("hresult.yml")

	if err != nil {
		return Repo{}, err
	}

	win32Errors, err := LoadWin32Errors("win32error.yml")

	if err != nil {
		return Repo{}, err
	}

	bugChecks, err := LoadBugChecks("bugcheck.yml")

	if err != nil {
		return Repo{}, err
	}

	return Repo{
		NTStatus:   ntStatuses,
		HResult:    hResults,
		Win32Error: win32Errors,
		BugCheck:   bugChecks,
	}, nil
}
