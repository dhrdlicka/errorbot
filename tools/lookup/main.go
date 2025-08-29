package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dhrdlicka/errorbot/repo"
	"github.com/dhrdlicka/errorbot/util"
)

var value = flag.String("c", "", "`error code` in decimal or hexadecimal format [e.g. 1, -2147024894, 0x7B, C0000005]")

func main() {
	log.SetFlags(0)
	log.SetPrefix("lookup: ")

	flag.Parse()

	if *value == "" {
		flag.Usage()
		os.Exit(1)
	}

	codes, err := util.ParseCode(*value)

	if err != nil {
		log.Fatal(err)
	}

	repoInstance, err := repo.Load()

	if err != nil {
		log.Fatal(err)
	}

	repos := []struct {
		findCode func(uint32) []repo.ErrorInfo
		name     string
	}{
		{repoInstance.FindBugCheck, "bug check"},
		{repoInstance.FindHResult, "HRESULT"},
		{repoInstance.FindWin32Error, "Win32 error"},
		{repoInstance.FindNTStatus, "NTSTATUS"},
	}

	found := false

	for _, errorRepo := range repos {
		matches := []repo.ErrorInfo{}

		for _, code := range codes {
			matches = append(matches, errorRepo.findCode(code)...)
		}

		if len(matches) > 0 {
			found = true

			fmt.Printf("# Possible %s codes:\n\n%s\n", errorRepo.name, formatResults(matches))
		}
	}

	if !found {
		log.Fatalf("could not find error code %s (`0x%08X`)\n", *value, codes[0])
	}
}

func formatResults(errors []repo.ErrorInfo) string {
	var result []byte

	for _, item := range errors {
		result = fmt.Appendf(result, "`%s` (`0x%08X`)\n", item.Name, item.Code)

		for _, line := range strings.Split(item.Description, "\n") {
			result = fmt.Appendf(result, "> %s\n", strings.TrimSpace(line))
		}
	}

	return string(result)
}
