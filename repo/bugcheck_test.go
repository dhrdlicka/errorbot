package repo

import (
	"reflect"
	"testing"
)

// createTestBugCheckRepo creates a test repository with sample BugCheck data
func createTestBugCheckRepo() BugCheckRepo {
	return BugCheckRepo{
		{
			Code:        0x0000000A,
			Name:        "IRQL_NOT_LESS_OR_EQUAL",
			URL:         "https://docs.microsoft.com/en-us/windows-hardware/drivers/debugger/bug-check-0xa--irql-not-less-or-equal",
			Description: "This indicates that Microsoft Windows or a kernel-mode driver accessed paged memory at DISPATCH_LEVEL or above.",
			Parameters: []string{
				"Memory address that was referenced",
				"IRQL that was required",
				"Type of access (0 = read, 1 = write)",
				"Address of instruction which referenced memory",
			},
		},
		{
			Code:        0x0000001E,
			Name:        "KMODE_EXCEPTION_NOT_HANDLED",
			URL:         "https://docs.microsoft.com/en-us/windows-hardware/drivers/debugger/bug-check-0x1e--kmode-exception-not-handled",
			Description: "This indicates that a kernel-mode program generated an exception which the error handler did not catch.",
			Parameters: []string{
				"The exception code that was not handled",
				"The address at which the exception occurred",
				"Parameter 0 of the exception",
				"Parameter 1 of the exception",
			},
		},
		{
			Code:        0x00000050,
			Name:        "PAGE_FAULT_IN_NONPAGED_AREA",
			URL:         "https://docs.microsoft.com/en-us/windows-hardware/drivers/debugger/bug-check-0x50--page-fault-in-nonpaged-area",
			Description: "This indicates that invalid system memory has been referenced.",
			Parameters:  []string{},
		},
		{
			Code:        0x000000D1,
			Name:        "DRIVER_IRQL_NOT_LESS_OR_EQUAL",
			URL:         "https://docs.microsoft.com/en-us/windows-hardware/drivers/debugger/bug-check-0xd1--driver-irql-not-less-or-equal",
			Description: "This indicates that a kernel-mode driver attempted to access pageable memory at a process IRQL that was too high.",
			Parameters: []string{
				"Memory address referenced",
				"IRQL at time of reference",
				"Type of access (0 = read, 1 = write)",
				"Address of instruction which referenced memory",
			},
		},
	}
}

func TestBugCheckRepo_FindBugCheckCode(t *testing.T) {
	repo := createTestBugCheckRepo()

	tests := []struct {
		name     string
		code     uint32
		expected []BugCheck
	}{
		{
			name: "find IRQL_NOT_LESS_OR_EQUAL",
			code: 0x0000000A,
			expected: []BugCheck{
				{
					Code:        0x0000000A,
					Name:        "IRQL_NOT_LESS_OR_EQUAL",
					URL:         "https://docs.microsoft.com/en-us/windows-hardware/drivers/debugger/bug-check-0xa--irql-not-less-or-equal",
					Description: "This indicates that Microsoft Windows or a kernel-mode driver accessed paged memory at DISPATCH_LEVEL or above.",
					Parameters: []string{
						"Memory address that was referenced",
						"IRQL that was required",
						"Type of access (0 = read, 1 = write)",
						"Address of instruction which referenced memory",
					},
				},
			},
		},
		{
			name: "find PAGE_FAULT_IN_NONPAGED_AREA",
			code: 0x00000050,
			expected: []BugCheck{
				{
					Code:        0x00000050,
					Name:        "PAGE_FAULT_IN_NONPAGED_AREA",
					URL:         "https://docs.microsoft.com/en-us/windows-hardware/drivers/debugger/bug-check-0x50--page-fault-in-nonpaged-area",
					Description: "This indicates that invalid system memory has been referenced.",
					Parameters:  []string{},
				},
			},
		},
		{
			name:     "find non-existent code",
			code:     0x00000999,
			expected: []BugCheck{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := repo.FindBugCheckCode(tt.code)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("BugCheckRepo.FindBugCheckCode(0x%08X) = %v, expected %v", tt.code, result, tt.expected)
			}
		})
	}
}

func TestBugCheckRepo_FindBugCheckString(t *testing.T) {
	repo := createTestBugCheckRepo()

	tests := []struct {
		name     string
		search   string
		expected int // Expected number of matches
	}{
		{
			name:     "find by exact name",
			search:   "IRQL_NOT_LESS_OR_EQUAL",
			expected: 1,
		},
		{
			name:     "find by partial name - case insensitive",
			search:   "irql",
			expected: 1, // Only matches IRQL_NOT_LESS_OR_EQUAL (case sensitive search)
		},
		{
			name:     "find by partial name - uppercase",
			search:   "PAGE_FAULT",
			expected: 1,
		},
		{
			name:     "find by partial name - lowercase",
			search:   "page_fault",
			expected: 1,
		},
		{
			name:     "find by single character",
			search:   "K",
			expected: 1, // Should match KMODE_EXCEPTION_NOT_HANDLED
		},
		{
			name:     "find non-existent string",
			search:   "NONEXISTENT",
			expected: 0,
		},
		{
			name:     "empty search string",
			search:   "",
			expected: 4, // Should match all (empty string is prefix of all strings)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := repo.FindBugCheckString(tt.search)
			if len(result) != tt.expected {
				t.Errorf("BugCheckRepo.FindBugCheckString(%q) returned %d matches, expected %d", tt.search, len(result), tt.expected)
			}
		})
	}
}

func TestBugCheckRepo_FindCode(t *testing.T) {
	repo := createTestBugCheckRepo()

	tests := []struct {
		name     string
		code     uint32
		expected []ErrorInfo
	}{
		{
			name: "find IRQL_NOT_LESS_OR_EQUAL as ErrorInfo",
			code: 0x0000000A,
			expected: []ErrorInfo{
				{
					Code:        0x0000000A,
					Name:        "IRQL_NOT_LESS_OR_EQUAL",
					Description: "This indicates that Microsoft Windows or a kernel-mode driver accessed paged memory at DISPATCH_LEVEL or above.",
				},
			},
		},
		{
			name: "find PAGE_FAULT_IN_NONPAGED_AREA as ErrorInfo",
			code: 0x00000050,
			expected: []ErrorInfo{
				{
					Code:        0x00000050,
					Name:        "PAGE_FAULT_IN_NONPAGED_AREA",
					Description: "This indicates that invalid system memory has been referenced.",
				},
			},
		},
		{
			name:     "find non-existent code as ErrorInfo",
			code:     0x00000999,
			expected: []ErrorInfo{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := repo.FindCode(tt.code)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("BugCheckRepo.FindCode(0x%08X) = %v, expected %v", tt.code, result, tt.expected)
			}
		})
	}
}

func TestRepo_FindBugCheck(t *testing.T) {
	repo := Repo{
		BugCheck: createTestBugCheckRepo(),
	}

	tests := []struct {
		name     string
		code     uint32
		expected []ErrorInfo
	}{
		{
			name: "find through Repo.FindBugCheck",
			code: 0x0000000A,
			expected: []ErrorInfo{
				{
					Code:        0x0000000A,
					Name:        "IRQL_NOT_LESS_OR_EQUAL",
					Description: "This indicates that Microsoft Windows or a kernel-mode driver accessed paged memory at DISPATCH_LEVEL or above.",
				},
			},
		},
		{
			name:     "find non-existent through Repo.FindBugCheck",
			code:     0x00000999,
			expected: []ErrorInfo{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := repo.FindBugCheck(tt.code)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Repo.FindBugCheck(0x%08X) = %v, expected %v", tt.code, result, tt.expected)
			}
		})
	}
}

func TestBugCheck_ErrorInfo(t *testing.T) {
	bugCheck := BugCheck{
		Code:        0x0000000A,
		Name:        "IRQL_NOT_LESS_OR_EQUAL",
		URL:         "https://example.com",
		Description: "Test description",
		Parameters:  []string{"param1", "param2"},
	}

	expected := ErrorInfo{
		Code:        0x0000000A,
		Name:        "IRQL_NOT_LESS_OR_EQUAL",
		Description: "Test description",
	}

	result := bugCheck.ErrorInfo()
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("BugCheck.ErrorInfo() = %v, expected %v", result, expected)
	}
}

// Test edge cases
func TestBugCheckRepo_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		repo     BugCheckRepo
		code     uint32
		expected int
	}{
		{
			name:     "empty repository",
			repo:     BugCheckRepo{},
			code:     0x0000000A,
			expected: 0,
		},
		{
			name:     "nil repository",
			repo:     nil,
			code:     0x0000000A,
			expected: 0,
		},
		{
			name: "repository with duplicate codes",
			repo: BugCheckRepo{
				{Code: 0x0000000A, Name: "DUPLICATE_1", Description: "First"},
				{Code: 0x0000000A, Name: "DUPLICATE_2", Description: "Second"},
			},
			code:     0x0000000A,
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.repo.FindBugCheckCode(tt.code)
			if len(result) != tt.expected {
				t.Errorf("BugCheckRepo.FindBugCheckCode() returned %d matches, expected %d", len(result), tt.expected)
			}
		})
	}
}

// Benchmark tests
func BenchmarkBugCheckRepo_FindBugCheckCode(b *testing.B) {
	repo := createTestBugCheckRepo()

	for i := 0; i < b.N; i++ {
		_ = repo.FindBugCheckCode(0x0000000A)
	}
}

func BenchmarkBugCheckRepo_FindBugCheckString(b *testing.B) {
	repo := createTestBugCheckRepo()

	for i := 0; i < b.N; i++ {
		_ = repo.FindBugCheckString("IRQL")
	}
}

func BenchmarkBugCheckRepo_FindCode(b *testing.B) {
	repo := createTestBugCheckRepo()

	for i := 0; i < b.N; i++ {
		_ = repo.FindCode(0x0000000A)
	}
}

func BenchmarkBugCheck_String(b *testing.B) {
	bugCheck := BugCheck{
		Code:        0x0000000A,
		Name:        "IRQL_NOT_LESS_OR_EQUAL",
		Description: "Test description",
		Parameters:  []string{"param1", "param2", "param3"},
	}

	for i := 0; i < b.N; i++ {
		_ = bugCheck.String()
	}
}

// Example tests
func ExampleBugCheckRepo_FindBugCheckCode() {
	repo := BugCheckRepo{
		{Code: 0x0000000A, Name: "IRQL_NOT_LESS_OR_EQUAL", Description: "IRQL error"},
	}

	matches := repo.FindBugCheckCode(0x0000000A)
	if len(matches) > 0 {
		println(matches[0].Name) // IRQL_NOT_LESS_OR_EQUAL
	}
}

func ExampleBugCheckRepo_FindBugCheckString() {
	repo := BugCheckRepo{
		{Code: 0x0000000A, Name: "IRQL_NOT_LESS_OR_EQUAL", Description: "IRQL error"},
		{Code: 0x000000D1, Name: "DRIVER_IRQL_NOT_LESS_OR_EQUAL", Description: "Driver IRQL error"},
	}

	matches := repo.FindBugCheckString("IRQL")
	println(len(matches)) // 2
}
