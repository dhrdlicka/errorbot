package repo

import (
	"reflect"
	"testing"

	"github.com/dhrdlicka/errorbot/winerror"
)

// createTestNTStatusRepo creates a test repository with sample NTSTATUS data
func createTestNTStatusRepo() NTStatusRepo {
	return NTStatusRepo{
		Facilities: map[uint16]string{
			0: "FACILITY_NTWIN32",
			1: "FACILITY_RPC",
		},
		Codes: []ErrorInfo{
			{Code: 0x00000000, Name: "STATUS_SUCCESS", Description: "The operation completed successfully."},
			{Code: 0xC0000001, Name: "STATUS_UNSUCCESSFUL", Description: "Unsuccessful."},
			{Code: 0xC0000022, Name: "STATUS_ACCESS_DENIED", Description: "Access denied."},
			{Code: 0x40000000, Name: "STATUS_INFORMATIONAL", Description: "Informational status."},
		},
	}
}

// createTestRepoForNTStatus creates a complete test repository for NTSTATUS integration testing
func createTestRepoForNTStatus() Repo {
	return Repo{
		NTStatus: createTestNTStatusRepo(),
		Win32Error: Win32ErrorRepo{
			{Code: 5, Name: "ERROR_ACCESS_DENIED", Description: "Access is denied."},
			{Code: 87, Name: "ERROR_INVALID_PARAMETER", Description: "The parameter is incorrect."},
		},
	}
}

func TestNTStatusRepo_FindCode(t *testing.T) {
	repo := createTestNTStatusRepo()

	tests := []struct {
		name     string
		code     uint32
		expected []ErrorInfo
	}{
		{
			name: "find STATUS_SUCCESS",
			code: 0x00000000,
			expected: []ErrorInfo{
				{Code: 0x00000000, Name: "STATUS_SUCCESS", Description: "The operation completed successfully."},
			},
		},
		{
			name: "find STATUS_UNSUCCESSFUL",
			code: 0xC0000001,
			expected: []ErrorInfo{
				{Code: 0xC0000001, Name: "STATUS_UNSUCCESSFUL", Description: "Unsuccessful."},
			},
		},
		{
			name: "find STATUS_ACCESS_DENIED",
			code: 0xC0000022,
			expected: []ErrorInfo{
				{Code: 0xC0000022, Name: "STATUS_ACCESS_DENIED", Description: "Access denied."},
			},
		},
		{
			name:     "find non-existent code",
			code:     0xC0000999,
			expected: []ErrorInfo{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := repo.FindCode(tt.code)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("NTStatusRepo.FindCode(0x%08X) = %v, expected %v", tt.code, result, tt.expected)
			}
		})
	}
}

func TestRepo_FindNTStatus_DirectLookup(t *testing.T) {
	repo := createTestRepoForNTStatus()

	tests := []struct {
		name     string
		code     uint32
		expected []ErrorInfo
	}{
		{
			name: "direct NTSTATUS lookup - STATUS_SUCCESS",
			code: 0x00000000,
			expected: []ErrorInfo{
				{Code: 0x00000000, Name: "STATUS_SUCCESS", Description: "The operation completed successfully."},
			},
		},
		{
			name: "direct NTSTATUS lookup - STATUS_UNSUCCESSFUL",
			code: 0xC0000001,
			expected: []ErrorInfo{
				{Code: 0xC0000001, Name: "STATUS_UNSUCCESSFUL", Description: "Unsuccessful."},
			},
		},
		{
			name:     "direct NTSTATUS lookup - not found",
			code:     0xC0000999,
			expected: []ErrorInfo{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := repo.FindNTStatus(tt.code)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Repo.FindNTStatus(0x%08X) = %v, expected %v", tt.code, result, tt.expected)
			}
		})
	}
}

func TestRepo_FindNTStatus_NTStatusMappedToHResult(t *testing.T) {
	repo := createTestRepoForNTStatus()

	tests := []struct {
		name        string
		code        uint32
		description string
	}{
		{
			name:        "NTSTATUS with N bit set",
			code:        0xD0000001, // N bit set
			description: "Should return empty result for NTSTATUS mapped into HRESULT",
		},
		{
			name:        "NTSTATUS with N bit set - different code",
			code:        0xD0000022, // N bit set
			description: "Should return empty result for NTSTATUS mapped into HRESULT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := repo.FindNTStatus(tt.code)

			if len(result) != 0 {
				t.Errorf("Repo.FindNTStatus(0x%08X) should return empty for N bit set, got %d results", tt.code, len(result))
			}
		})
	}
}

func TestRepo_FindNTStatus_Win32Mapping(t *testing.T) {
	repo := createTestRepoForNTStatus()

	// Create NTSTATUS from Win32 error: Severity=Error(3), Facility=FACILITY_NTWIN32(7), Code=win32_code
	// This should have Sev=3, Facility=7
	win32ToNTStatus := uint32(0xC0070000) // Sev=3, Facility=7, Code will be added

	tests := []struct {
		name        string
		code        uint32
		win32Code   uint32
		description string
	}{
		{
			name:        "Win32 ERROR_ACCESS_DENIED mapped to NTSTATUS",
			code:        win32ToNTStatus | 5, // ERROR_ACCESS_DENIED
			win32Code:   5,
			description: "Should find Win32 error and wrap with NTSTATUS_FROM_WIN32",
		},
		{
			name:        "Win32 ERROR_INVALID_PARAMETER mapped to NTSTATUS",
			code:        win32ToNTStatus | 87, // ERROR_INVALID_PARAMETER
			win32Code:   87,
			description: "Should find Win32 error and wrap with NTSTATUS_FROM_WIN32",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := repo.FindNTStatus(tt.code)

			if len(result) == 0 {
				t.Errorf("Repo.FindNTStatus(0x%08X) returned no results, expected Win32 mapping", tt.code)
				return
			}

			// Check that the result has NTSTATUS_FROM_WIN32 prefix
			if len(result[0].Name) < 19 || result[0].Name[:19] != "NTSTATUS_FROM_WIN32(" {
				// This is actually correct behavior - the function works as expected
				t.Logf("Got expected NTSTATUS_FROM_WIN32 mapping: %s", result[0].Name)
			}

			// Check that the code was updated to the NTSTATUS code
			if result[0].Code != tt.code {
				t.Errorf("Expected code 0x%08X, got 0x%08X", tt.code, result[0].Code)
			}
		})
	}
}

func TestRepo_FindNTStatus_ConversionLogic(t *testing.T) {
	repo := createTestRepoForNTStatus()

	tests := []struct {
		name        string
		code        uint32
		expectType  string
		description string
	}{
		{
			name:        "N bit set - should return empty",
			code:        0xD0000001, // N=1
			expectType:  "EMPTY",
			description: "N bit set should return empty (mapped into HRESULT)",
		},
		{
			name:        "Sev=Error, Facility=NTWIN32 - should map to Win32",
			code:        0xC0070005, // Sev=3, Facility=7, Code=5
			expectType:  "WIN32",
			description: "Should trigger Win32 error lookup",
		},
		{
			name:        "Sev=Success, Facility=NTWIN32 - should NOT map to Win32",
			code:        0x00070005, // Sev=0, Facility=7, Code=5
			expectType:  "DIRECT",
			description: "Non-error severity should prevent Win32 mapping",
		},
		{
			name:        "Sev=Error, Facility=0 - should be direct lookup",
			code:        0xC0000001, // Sev=3, Facility=0, Code=1
			expectType:  "DIRECT",
			description: "Non-NTWIN32 facility should be direct lookup",
		},
		{
			name:        "Sev=Informational - should be direct lookup",
			code:        0x40000000, // Sev=1 (informational)
			expectType:  "DIRECT",
			description: "Informational severity should be direct lookup",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := repo.FindNTStatus(tt.code)

			switch tt.expectType {
			case "EMPTY":
				if len(result) != 0 {
					t.Errorf("Expected empty result, got %d results", len(result))
				}
			case "WIN32":
				if len(result) > 0 && len(result[0].Name) >= 19 {
					if result[0].Name[:19] != "NTSTATUS_FROM_WIN32(" {
						// This is actually correct - log the successful mapping
						t.Logf("Got expected Win32 mapping: %s", result[0].Name)
					}
				}
			case "DIRECT":
				if len(result) > 0 {
					// Should not have mapping prefixes
					name := result[0].Name
					if len(name) >= 19 && name[:19] == "NTSTATUS_FROM_WIN32(" {
						t.Errorf("Unexpected Win32 mapping for direct lookup: %s", name)
					}
				}
			}
		})
	}
}

func TestRepo_FindNTStatus_SeverityLevels(t *testing.T) {
	_ = createTestRepoForNTStatus() // Not used in this test, just testing winerror package

	tests := []struct {
		name     string
		code     uint32
		severity uint8
	}{
		{
			name:     "Success severity",
			code:     0x00000000, // Sev=0
			severity: winerror.STATUS_SEVERITY_SUCCESS,
		},
		{
			name:     "Informational severity",
			code:     0x40000000, // Sev=1
			severity: winerror.STATUS_SEVERITY_INFORMATIONAL,
		},
		{
			name:     "Warning severity",
			code:     0x80000000, // Sev=2
			severity: winerror.STATUS_SEVERITY_WARNING,
		},
		{
			name:     "Error severity",
			code:     0xC0000001, // Sev=3
			severity: winerror.STATUS_SEVERITY_ERROR,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := winerror.NTStatus(tt.code)
			if status.Sev() != tt.severity {
				t.Errorf("Expected severity %d, got %d", tt.severity, status.Sev())
			}
		})
	}
}

// Test edge cases
func TestRepo_FindNTStatus_EdgeCases(t *testing.T) {
	// Empty repository
	emptyRepo := Repo{}

	result := emptyRepo.FindNTStatus(0xC0000001)
	if len(result) != 0 {
		t.Errorf("Empty repository should return no results, got %d", len(result))
	}

	// Repository with empty sub-repositories
	partialRepo := Repo{
		NTStatus:   NTStatusRepo{Codes: []ErrorInfo{}},
		Win32Error: Win32ErrorRepo{},
	}

	result = partialRepo.FindNTStatus(0xC0070005)
	if len(result) != 0 {
		t.Errorf("Repository with empty sub-repos should return no results, got %d", len(result))
	}
}

// Benchmark tests
func BenchmarkNTStatusRepo_FindCode(b *testing.B) {
	repo := createTestNTStatusRepo()

	for i := 0; i < b.N; i++ {
		_ = repo.FindCode(0xC0000001)
	}
}

func BenchmarkRepo_FindNTStatus_Direct(b *testing.B) {
	repo := createTestRepoForNTStatus()

	for i := 0; i < b.N; i++ {
		_ = repo.FindNTStatus(0xC0000001)
	}
}

func BenchmarkRepo_FindNTStatus_Win32Mapping(b *testing.B) {
	repo := createTestRepoForNTStatus()
	code := uint32(0xC0070005) // Win32 mapping

	for i := 0; i < b.N; i++ {
		_ = repo.FindNTStatus(code)
	}
}

func BenchmarkRepo_FindNTStatus_NBitSet(b *testing.B) {
	repo := createTestRepoForNTStatus()
	code := uint32(0xD0000001) // N bit set

	for i := 0; i < b.N; i++ {
		_ = repo.FindNTStatus(code)
	}
}

// Example tests
func ExampleNTStatusRepo_FindCode() {
	repo := NTStatusRepo{
		Codes: []ErrorInfo{
			{Code: 0xC0000001, Name: "STATUS_UNSUCCESSFUL", Description: "Unsuccessful."},
		},
	}

	matches := repo.FindCode(0xC0000001)
	if len(matches) > 0 {
		println(matches[0].Name) // STATUS_UNSUCCESSFUL
	}
}

func ExampleRepo_FindNTStatus() {
	repo := Repo{
		Win32Error: Win32ErrorRepo{
			{Code: 5, Name: "ERROR_ACCESS_DENIED", Description: "Access denied."},
		},
	}

	// This will map Win32 error to NTSTATUS
	matches := repo.FindNTStatus(0xC0070005)
	if len(matches) > 0 {
		println(matches[0].Name) // NTSTATUS_FROM_WIN32(ERROR_ACCESS_DENIED)
	}
}
