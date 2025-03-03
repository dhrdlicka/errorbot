package commands

import (
	"strconv"
	"strings"
)

func parseCode(code string) ([]uint32, error) {
	if strings.HasPrefix(code, "0x") || strings.HasPrefix(code, "0X") {
		// hex prefix, we are almost there
		longCode, err := strconv.ParseUint(code, 16, 32)

		if err != nil {
			return nil, err
		}

		return []uint32{uint32(longCode)}, nil
	}

	if code[0] == '-' {
		// negative number, probably signed decimal
		intCode, err := strconv.ParseInt(code, 10, 32)

		if err != nil {
			return nil, err
		}

		return []uint32{uint32(intCode)}, nil
	}

	// no prefix, now the real fun begins
	var codes []uint32
	var lastErr error

	if hexCode, err := strconv.ParseUint(code, 16, 32); err == nil {
		codes = append(codes, uint32(hexCode))
	} else {
		lastErr = err
	}

	if decCode, err := strconv.ParseUint(code, 10, 32); err == nil {
		codes = append(codes, uint32(decCode))
	} else {
		lastErr = err
	}

	if len(codes) == 0 {
		return nil, lastErr
	}

	return codes, nil
}
