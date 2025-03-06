package repo

import (
	"fmt"
	"os"

	"github.com/dhrdlicka/errorbot/winerror"
	"gopkg.in/yaml.v3"
)

type HResultRepo struct {
	Facilities map[uint16]string
	Codes      []ErrorInfo
}

func LoadHResults(name string) (HResultRepo, error) {
	file, err := os.ReadFile(name)

	if err != nil {
		return HResultRepo{}, err
	}

	var repo HResultRepo
	err = yaml.Unmarshal(file, &repo)

	if err != nil {
		return HResultRepo{}, err
	}

	return repo, nil
}

func (repo HResultRepo) FindCode(code uint32) []ErrorInfo {
	return FindCode(repo.Codes, code)
}

func (repo Repo) FindHResult(code uint32) []ErrorInfo {
	hr := winerror.HResult(code)

	if hr.N() {
		// this is a mapped NTSTATUS
		ntStatusMatches := repo.FindNTStatus(code ^ winerror.FACILITY_NT_BIT)

		for i := range ntStatusMatches {
			ntStatusMatches[i].Name = fmt.Sprintf("HRESULT_FROM_NT(%s)", ntStatusMatches[i].Name)
			ntStatusMatches[i].Code = code
		}

		return ntStatusMatches
	} else if hr.Facility() == winerror.FACILITY_WIN32 {
		// this is a mapped Win32 error
		win32ErrorMatches := repo.FindWin32Error(uint32(hr.Code()))

		for i := range win32ErrorMatches {
			win32ErrorMatches[i].Name = fmt.Sprintf("HRESULT_FROM_WIN32(%s)", win32ErrorMatches[i].Name)
			win32ErrorMatches[i].Code = code
		}

		return win32ErrorMatches

	}

	return repo.HResult.FindCode(code)
}
