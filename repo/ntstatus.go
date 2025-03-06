package repo

import (
	"fmt"
	"os"

	"github.com/dhrdlicka/errorbot/winerror"
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
	s := winerror.NTStatus(code)

	if s.N() {
		// this is an NTSTATUS mapped into an HRESULT
		return []ErrorInfo{}
	} else if s.Sev() == winerror.STATUS_SEVERITY_ERROR && s.Facility() == winerror.FACILITY_NTWIN32 {
		// this is a mapped Win32 error
		win32ErrorMatches := repo.Win32Error.FindCode(uint32(s.Code()))

		for i := range win32ErrorMatches {
			win32ErrorMatches[i].Name = fmt.Sprintf("NTSTATUS_FROM_WIN32(%s)", win32ErrorMatches[i].Name)
			win32ErrorMatches[i].Code = code
		}

		return win32ErrorMatches
	}

	return repo.NTStatus.FindCode(code)
}
