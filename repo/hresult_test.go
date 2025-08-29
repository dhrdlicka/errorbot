package repo

import (
	"reflect"
	"testing"

	"github.com/dhrdlicka/errorbot/winerror"
)

// createTestHResultRepo creates a test repository with sample HRESULT data
func createTestHResultRepo() HResultRepo {
	return HResultRepo{
		Facilities: map[uint16]string{
			0: "FACILITY_NULL",
			1: "FACILITY_RPC",
			7: "FACILITY_WIN32",
		},
		Codes: []ErrorInfo{
			{Code: 0x00000000, Name: "S_OK", Description: "Operation successful."},
			{Code: 0x80004001, Name: "E_NOTIMPL", Description: "Not implemented."},
			{Code: 0x80004005, Name: "E_FAIL", Description: "Unspecified failure."},
			{Code: 0x80070005, Name: "E_ACCESSDENIED", Description: "General access denied error."},
		},
	}
}

// createTestRepo creates a complete test repository for integration testing
func createTestRepo() Repo {
	return Repo{
		HResult: createTestHResultRepo(),
		Win32Error: Win32ErrorRepo{
			{Code: 5, Name: "ERROR_ACCESS_DENIED", Description: "Access is denied."},
			{Code: 87, Name: "ERROR_INVALID_PARAMETER", Description: "The parameter is incorrect."},
		},
		NTStatus: NTStatusRepo{
			Facilities: map[uint16]string{
				0: "FACILITY_NTWIN32",
			},
			Codes: []ErrorInfo{
				{Code: 0xC0000001, Name: "STATUS_UNSUCCESSFUL", Description: "Unsuccessful."},
				{Code: 0xC0000022, Name: "STATUS_ACCESS_DENIED", Description: "Access denied."},
			},
		},
	}
}

func TestHResultRepo_FindCode(t *testing.T) {
	repo := createTestHResultRepo()

	tests := []struct {
		name     string
		code     uint32
		expected []ErrorInfo
	}{
		{
			name: "find S_OK",
			code: 0x00000000,
			expected: []ErrorInfo{
				{Code: 0x00000000, Name: "S_OK", Description: "Operation successful."},
			},
		},
		{
			name: "find E_ACCESSDENIED",
			code: 0x80070005,
			expected: []ErrorInfo{
				{Code: 0x80070005, Name: "E_ACCESSDENIED", Description: "General access denied error."},
			},
		},
		{
			name:     "find non-existent code",
			code:     0x80000999,
			expected: []ErrorInfo{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := repo.FindCode(tt.code)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("HResultRepo.FindCode(0x%08X) = %v, expected %v", tt.code, result, tt.expected)
			}
		})
	}
}

func TestRepo_FindHResult_DirectLookup(t *testing.T) {
	repo := createTestRepo()

	tests := []struct {
		name     string
		code     uint32
		expected []ErrorInfo
	}{
		{
			name: "direct HRESULT lookup - S_OK",
			code: 0x00000000,
			expected: []ErrorInfo{
				{Code: 0x00000000, Name: "S_OK", Description: "Operation successful."},
			},
		},
		{
			name: "direct HRESULT lookup - E_ACCESSDENIED (maps to Win32)",
			code: 0x80070005,
			expected: []ErrorInfo{
				{Code: 0x80070005, Name: "HRESULT_FROM_WIN32(ERROR_ACCESS_DENIED)", Description: "Access is denied."},
			},
		},
		{
			name:     "direct HRESULT lookup - not found",
			code:     0x80000999,
			expected: []ErrorInfo{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := repo.FindHResult(tt.code)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Repo.FindHResult(0x%08X) = %v, expected %v", tt.code, result, tt.expected)
			}
		})
	}
}

func TestRepo_FindHResult_NTStatusMapping(t *testing.T) {
	repo := createTestRepo()

	tests := []struct {
		name        string
		code        uint32
		description string
	}{
		{
			name:        "NTSTATUS mapped to HRESULT",
			code:        0xC0000001 | winerror.FACILITY_NT_BIT, // STATUS_UNSUCCESSFUL with N bit
			description: "Should find NTSTATUS and wrap with HRESULT_FROM_NT",
		},
		{
			name:        "NTSTATUS ACCESS_DENIED mapped to HRESULT",
			code:        0xC0000022 | winerror.FACILITY_NT_BIT, // STATUS_ACCESS_DENIED with N bit
			description: "Should find NTSTATUS ACCESS_DENIED and wrap with HRESULT_FROM_NT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := repo.FindHResult(tt.code)

			if len(result) == 0 {
				t.Errorf("Repo.FindHResult(0x%08X) returned no results, expected NTSTATUS mapping", tt.code)
				return
			}

			// Check that the result has HRESULT_FROM_NT prefix
			if len(result[0].Name) < 16 || result[0].Name[:16] != "HRESULT_FROM_NT(" {
				t.Errorf("Expected HRESULT_FROM_NT prefix, got: %s", result[0].Name)
			}

			// Check that the code was updated to the HRESULT code
			if result[0].Code != tt.code {
				t.Errorf("Expected code 0x%08X, got 0x%08X", tt.code, result[0].Code)
			}
		})
	}
}

func TestRepo_FindHResult_Win32Mapping(t *testing.T) {
	repo := createTestRepo()

	// Create HRESULT from Win32 error: 0x80070000 | win32_code
	// This should have S=1, R=0, Facility=7 (FACILITY_WIN32)
	win32ToHResult := uint32(0x80070000) // S=1, R=0, Facility=7, Code will be added

	tests := []struct {
		name        string
		code        uint32
		win32Code   uint32
		description string
	}{
		{
			name:        "Win32 ERROR_ACCESS_DENIED mapped to HRESULT",
			code:        win32ToHResult | 5, // ERROR_ACCESS_DENIED
			win32Code:   5,
			description: "Should find Win32 error and wrap with HRESULT_FROM_WIN32",
		},
		{
			name:        "Win32 ERROR_INVALID_PARAMETER mapped to HRESULT",
			code:        win32ToHResult | 87, // ERROR_INVALID_PARAMETER
			win32Code:   87,
			description: "Should find Win32 error and wrap with HRESULT_FROM_WIN32",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := repo.FindHResult(tt.code)

			if len(result) == 0 {
				t.Errorf("Repo.FindHResult(0x%08X) returned no results, expected Win32 mapping", tt.code)
				return
			}

			// Check that the result has HRESULT_FROM_WIN32 prefix
			if len(result[0].Name) < 19 || result[0].Name[:19] != "HRESULT_FROM_WIN32(" {
				t.Errorf("Expected HRESULT_FROM_WIN32 prefix, got: %s", result[0].Name)
			}

			// Check that the code was updated to the HRESULT code
			if result[0].Code != tt.code {
				t.Errorf("Expected code 0x%08X, got 0x%08X", tt.code, result[0].Code)
			}
		})
	}
}

func TestRepo_FindHResult_ConversionLogic(t *testing.T) {
	repo := createTestRepo()

	tests := []struct {
		name        string
		code        uint32
		expectType  string
		description string
	}{
		{
			name:        "N bit set - should map to NTSTATUS",
			code:        0x90000001, // S=1, N=1
			expectType:  "NTSTATUS",
			description: "N bit set should trigger NTSTATUS lookup",
		},
		{
			name:        "S=1, R=0, Facility=7 - should map to Win32",
			code:        0x80070005, // S=1, R=0, Facility=7, Code=5
			expectType:  "WIN32",
			description: "Should trigger Win32 error lookup",
		},
		{
			name:        "S=1, R=1, Facility=7 - should NOT map to Win32",
			code:        0xC0070005, // S=1, R=1, Facility=7, Code=5
			expectType:  "DIRECT",
			description: "R bit set should prevent Win32 mapping",
		},
		{
			name:        "S=0 - should be direct lookup",
			code:        0x00000000, // S=0 (success)
			expectType:  "DIRECT",
			description: "Success codes should be direct lookup",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := repo.FindHResult(tt.code)

			switch tt.expectType {
			case "NTSTATUS":
				if len(result) > 0 && len(result[0].Name) >= 16 {
					if result[0].Name[:16] != "HRESULT_FROM_NT(" {
						t.Errorf("Expected NTSTATUS mapping, got: %s", result[0].Name)
					}
				}
			case "WIN32":
				if len(result) > 0 && len(result[0].Name) >= 19 {
					if result[0].Name[:19] != "HRESULT_FROM_WIN32(" {
						t.Errorf("Expected Win32 mapping, got: %s", result[0].Name)
					}
				}
			case "DIRECT":
				if len(result) > 0 {
					// Should not have mapping prefixes
					name := result[0].Name
					if len(name) >= 16 && name[:16] == "HRESULT_FROM_NT(" {
						t.Errorf("Unexpected NTSTATUS mapping for direct lookup: %s", name)
					}
					if len(name) >= 19 && name[:19] == "HRESULT_FROM_WIN32(" {
						t.Errorf("Unexpected Win32 mapping for direct lookup: %s", name)
					}
				}
			}
		})
	}
}

// Test edge cases
func TestRepo_FindHResult_EdgeCases(t *testing.T) {
	// Empty repository
	emptyRepo := Repo{}

	result := emptyRepo.FindHResult(0x80070005)
	if len(result) != 0 {
		t.Errorf("Empty repository should return no results, got %d", len(result))
	}

	// Repository with empty sub-repositories
	partialRepo := Repo{
		HResult:    HResultRepo{Codes: []ErrorInfo{}},
		Win32Error: Win32ErrorRepo{},
		NTStatus:   NTStatusRepo{Codes: []ErrorInfo{}},
	}

	result = partialRepo.FindHResult(0x80070005)
	if len(result) != 0 {
		t.Errorf("Repository with empty sub-repos should return no results, got %d", len(result))
	}
}

// Benchmark tests
func BenchmarkHResultRepo_FindCode(b *testing.B) {
	repo := createTestHResultRepo()

	for i := 0; i < b.N; i++ {
		_ = repo.FindCode(0x80070005)
	}
}

func BenchmarkRepo_FindHResult_Direct(b *testing.B) {
	repo := createTestRepo()

	for i := 0; i < b.N; i++ {
		_ = repo.FindHResult(0x80070005)
	}
}

func BenchmarkRepo_FindHResult_NTStatusMapping(b *testing.B) {
	repo := createTestRepo()
	code := uint32(0xC0000001 | winerror.FACILITY_NT_BIT)

	for i := 0; i < b.N; i++ {
		_ = repo.FindHResult(code)
	}
}

func BenchmarkRepo_FindHResult_Win32Mapping(b *testing.B) {
	repo := createTestRepo()
	code := uint32(0x80070005) // Win32 mapping

	for i := 0; i < b.N; i++ {
		_ = repo.FindHResult(code)
	}
}

// Example tests
func ExampleHResultRepo_FindCode() {
	repo := HResultRepo{
		Codes: []ErrorInfo{
			{Code: 0x80070005, Name: "E_ACCESSDENIED", Description: "Access denied."},
		},
	}

	matches := repo.FindCode(0x80070005)
	if len(matches) > 0 {
		println(matches[0].Name) // E_ACCESSDENIED
	}
}

func ExampleRepo_FindHResult() {
	repo := Repo{
		Win32Error: Win32ErrorRepo{
			{Code: 5, Name: "ERROR_ACCESS_DENIED", Description: "Access denied."},
		},
	}

	// This will map Win32 error to HRESULT
	matches := repo.FindHResult(0x80070005)
	if len(matches) > 0 {
		println(matches[0].Name) // HRESULT_FROM_WIN32(ERROR_ACCESS_DENIED)
	}
}
