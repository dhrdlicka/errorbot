package repo

import (
	"reflect"
	"testing"
)

func TestFindCode(t *testing.T) {
	testErrors := []ErrorInfo{
		{Code: 1, Name: "ERROR_INVALID_FUNCTION", Description: "Incorrect function."},
		{Code: 2, Name: "ERROR_FILE_NOT_FOUND", Description: "The system cannot find the file specified."},
		{Code: 5, Name: "ERROR_ACCESS_DENIED", Description: "Access is denied."},
		{Code: 5, Name: "ERROR_ACCESS_DENIED_DUPLICATE", Description: "Another access denied error."}, // Duplicate code
		{Code: 0, Name: "ERROR_SUCCESS", Description: "The operation completed successfully."},
	}

	tests := []struct {
		name     string
		errors   []ErrorInfo
		code     uint32
		expected []ErrorInfo
	}{
		{
			name:   "find existing single code",
			errors: testErrors,
			code:   1,
			expected: []ErrorInfo{
				{Code: 1, Name: "ERROR_INVALID_FUNCTION", Description: "Incorrect function."},
			},
		},
		{
			name:   "find existing duplicate codes",
			errors: testErrors,
			code:   5,
			expected: []ErrorInfo{
				{Code: 5, Name: "ERROR_ACCESS_DENIED", Description: "Access is denied."},
				{Code: 5, Name: "ERROR_ACCESS_DENIED_DUPLICATE", Description: "Another access denied error."},
			},
		},
		{
			name:     "find non-existent code",
			errors:   testErrors,
			code:     999,
			expected: []ErrorInfo{},
		},
		{
			name:   "find zero code",
			errors: testErrors,
			code:   0,
			expected: []ErrorInfo{
				{Code: 0, Name: "ERROR_SUCCESS", Description: "The operation completed successfully."},
			},
		},
		{
			name:     "empty error list",
			errors:   []ErrorInfo{},
			code:     1,
			expected: []ErrorInfo{},
		},
		{
			name:     "nil error list",
			errors:   nil,
			code:     1,
			expected: []ErrorInfo{},
		},
		{
			name:   "find max uint32 code",
			errors: []ErrorInfo{{Code: 0xFFFFFFFF, Name: "MAX_CODE", Description: "Max code test"}},
			code:   0xFFFFFFFF,
			expected: []ErrorInfo{
				{Code: 0xFFFFFFFF, Name: "MAX_CODE", Description: "Max code test"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindCode(tt.errors, tt.code)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("FindCode() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestErrorInfo_ErrorInfo(t *testing.T) {
	tests := []struct {
		name      string
		errorInfo ErrorInfo
		expected  ErrorInfo
	}{
		{
			name: "basic error info",
			errorInfo: ErrorInfo{
				Code:        5,
				Name:        "ERROR_ACCESS_DENIED",
				Description: "Access is denied.",
			},
			expected: ErrorInfo{
				Code:        5,
				Name:        "ERROR_ACCESS_DENIED",
				Description: "Access is denied.",
			},
		},
		{
			name:      "zero value error info",
			errorInfo: ErrorInfo{},
			expected:  ErrorInfo{},
		},
		{
			name: "error info with empty strings",
			errorInfo: ErrorInfo{
				Code:        0,
				Name:        "",
				Description: "",
			},
			expected: ErrorInfo{
				Code:        0,
				Name:        "",
				Description: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.errorInfo.ErrorInfo()
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ErrorInfo.ErrorInfo() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// Test edge cases and boundary conditions
func TestFindCode_EdgeCases(t *testing.T) {
	largeErrorList := make([]ErrorInfo, 1000)
	for i := 0; i < 1000; i++ {
		largeErrorList[i] = ErrorInfo{
			Code:        uint32(i),
			Name:        "ERROR_" + string(rune(i)),
			Description: "Description for error " + string(rune(i)),
		}
	}

	tests := []struct {
		name     string
		errors   []ErrorInfo
		code     uint32
		expected int // Expected number of matches
	}{
		{
			name:     "large error list - find first",
			errors:   largeErrorList,
			code:     0,
			expected: 1,
		},
		{
			name:     "large error list - find last",
			errors:   largeErrorList,
			code:     999,
			expected: 1,
		},
		{
			name:     "large error list - find middle",
			errors:   largeErrorList,
			code:     500,
			expected: 1,
		},
		{
			name:     "large error list - not found",
			errors:   largeErrorList,
			code:     1000,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindCode(tt.errors, tt.code)
			if len(result) != tt.expected {
				t.Errorf("FindCode() returned %d matches, expected %d", len(result), tt.expected)
			}
			if tt.expected > 0 && len(result) > 0 {
				if result[0].Code != tt.code {
					t.Errorf("FindCode() returned wrong code %d, expected %d", result[0].Code, tt.code)
				}
			}
		})
	}
}

// Benchmark tests
func BenchmarkFindCode_SingleMatch(b *testing.B) {
	errors := []ErrorInfo{
		{Code: 1, Name: "ERROR_1", Description: "Error 1"},
		{Code: 2, Name: "ERROR_2", Description: "Error 2"},
		{Code: 3, Name: "ERROR_3", Description: "Error 3"},
	}
	
	for i := 0; i < b.N; i++ {
		_ = FindCode(errors, 2)
	}
}

func BenchmarkFindCode_NoMatch(b *testing.B) {
	errors := []ErrorInfo{
		{Code: 1, Name: "ERROR_1", Description: "Error 1"},
		{Code: 2, Name: "ERROR_2", Description: "Error 2"},
		{Code: 3, Name: "ERROR_3", Description: "Error 3"},
	}
	
	for i := 0; i < b.N; i++ {
		_ = FindCode(errors, 999)
	}
}

func BenchmarkFindCode_LargeList(b *testing.B) {
	errors := make([]ErrorInfo, 1000)
	for i := 0; i < 1000; i++ {
		errors[i] = ErrorInfo{Code: uint32(i), Name: "ERROR", Description: "Desc"}
	}
	
	for i := 0; i < b.N; i++ {
		_ = FindCode(errors, 500)
	}
}

// Example tests for documentation
func ExampleFindCode() {
	errors := []ErrorInfo{
		{Code: 1, Name: "ERROR_INVALID_FUNCTION", Description: "Incorrect function."},
		{Code: 2, Name: "ERROR_FILE_NOT_FOUND", Description: "File not found."},
	}
	
	matches := FindCode(errors, 1)
	if len(matches) > 0 {
		println(matches[0].Name) // ERROR_INVALID_FUNCTION
	}
}

func ExampleErrorInfo_ErrorInfo() {
	err := ErrorInfo{Code: 5, Name: "ERROR_ACCESS_DENIED", Description: "Access denied."}
	info := err.ErrorInfo()
	println(info.Name) // ERROR_ACCESS_DENIED
}
