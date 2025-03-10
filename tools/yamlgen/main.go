package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	ntStatusRegex   = regexp.MustCompile(`#define (\w+)\s+\(\(NTSTATUS\)(0x[0-9A-F]{1,8})L\)`)
	hResultRegex    = regexp.MustCompile(`#define (\w+)\s+_HRESULT_TYPEDEF_\((0x[0-9A-F]{1,8})L\)`)
	win32ErrorRegex = regexp.MustCompile(`#define (\w+)\s+(\d+)L`)
	messageRegex    = regexp.MustCompile(`(0x[0-9A-F]{0,8}),\s*"(.*)"`)
)

var (
	headerPath   = flag.String("h", "", "")
	messagesPath = flag.String("mt", "", "")
	outputPath   = flag.String("o", "-", "")
	mode         = flag.String("m", "", "")
)

var codeFormat string = "0x%08X"

type uint32Hex uint32

func (U uint32Hex) MarshalYAML() (interface{}, error) {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: fmt.Sprintf(codeFormat, U),
	}, nil
}

type errorInfo struct {
	Code        uint32Hex
	Name        string
	Description string
}

func main() {
	flag.Parse()

	if *headerPath == "" || *messagesPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	var (
		err      error
		header   []byte
		messages []byte
		// outputFile       *os.File
	)

	if header, err = os.ReadFile(*headerPath); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if messages, err = os.ReadFile(*messagesPath); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	messageMap := map[uint32]string{}
	messageMatches := messageRegex.FindAllStringSubmatch(string(messages), -1)

	for _, match := range messageMatches {
		code, err := strconv.ParseUint(match[1], 0, 32)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		message := match[2]

		message = strings.ReplaceAll(message, "\\r", "\r")
		message = strings.ReplaceAll(message, "\\n", "\n")
		message = strings.TrimSpace(message)

		messageMap[uint32(code)] = message
	}

	var headerRegex *regexp.Regexp

	switch *mode {
	case "ntstatus":
		headerRegex = ntStatusRegex
	case "hresult":
		headerRegex = hResultRegex
	case "win32error":
		codeFormat = "%d"
		headerRegex = win32ErrorRegex
	default:
		fmt.Fprintf(os.Stderr, "invalid mode %s\n", *mode)
		os.Exit(1)
	}

	errors := []errorInfo{}

	headerMatches := headerRegex.FindAllStringSubmatch(string(header), -1)

	for _, match := range headerMatches {
		code, err := strconv.ParseUint(match[2], 0, 32)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		errors = append(errors, errorInfo{
			Name:        match[1],
			Code:        uint32Hex(code),
			Description: messageMap[uint32(code)],
		})
	}

	yaml, err := yaml.Marshal(errors)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var output *os.File

	if *outputPath == "-" {
		output = os.Stdout
	} else {
		output, err = os.Create(*outputPath)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	fmt.Fprint(output, string(yaml))
}
