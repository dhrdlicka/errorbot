package winerror

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type ErrorInfo struct {
	CodeString  string `json:"code"`
	Code        uint32
	Name        string `json:"name"`
	Description string `json:"description"`
}

func LoadErrorInfo(name string) ([]ErrorInfo, error) {
	file, err := os.ReadFile(name)

	if err != nil {
		return nil, err
	}

	var errors []ErrorInfo
	err = json.Unmarshal(file, &errors)

	if err != nil {
		return nil, err
	}

	for i, item := range errors {
		code, err := strconv.ParseUint(item.CodeString, 0, 32)

		if err != nil {
			return nil, err
		}

		errors[i].Code = uint32(code)
	}

	return errors, nil
}

func FindErrorCode(code uint32, errors []ErrorInfo) []ErrorInfo {
	matches := []ErrorInfo{}

	for _, item := range errors {
		if item.Code == code {
			matches = append(matches, item)
		}
	}

	return matches
}

func (errorInfo ErrorInfo) String() string {
	return fmt.Sprintf("`%s` (`0x%08X`)\n> %s\n", errorInfo.Name, errorInfo.Code, errorInfo.Description)
}
