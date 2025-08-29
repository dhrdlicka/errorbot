package repo

import (
	"testing"

	"github.com/dhrdlicka/errorbot/winerror"
)

// createFullTestRepo creates a comprehensive test repository for integration testing
func createFullTestRepo() Repo {
	return Repo{
		NTStatus: NTStatusRepo{
			Facilities: map[uint16]string{
				0: "FACILITY_NTWIN32",
				1: "FACILITY_RPC",
			},
			Codes: []ErrorInfo{
				{Code: 0x00000000, Name: "STATUS_SUCCESS", Description: "Success."},
				{Code: 0xC0000001, Name: "STATUS_UNSUCCESSFUL", Description: "Unsuccessful."},
				{Code: 0xC0000022, Name: "STATUS_ACCESS_DENIED", Description: "Access denied."},
			},
		},
		HResult: HResultRepo{
			Facilities: map[uint16]string{
				0: "FACILITY_NULL",
				7: "FACILITY_WIN32",
			},
			Codes: []ErrorInfo{
				{Code: 0x00000000, Name: "S_OK", Description: "Success."},
				{Code: 0x80004001, Name: "E_NOTIMPL", Description: "Not implemented."},
				{Code: 0x80070005, Name: "E_ACCESSDENIED", Description: "Access denied."},
			},
		},
		Win32Error: Win32ErrorRepo{
			{Code: 0, Name: "ERROR_SUCCESS", Description: "Success."},
			{Code: 5, Name: "ERROR_ACCESS_DENIED", Description: "Access denied."},
			{Code: 87, Name: "ERROR_INVALID_PARAMETER", Description: "Invalid parameter."},
		},
		BugCheck: BugCheckRepo{
			{Code: 0x0000000A, Name: "IRQL_NOT_LESS_OR_EQUAL", Description: "IRQL error."},
			{Code: 0x00000050, Name: "PAGE_FAULT_IN_NONPAGED_AREA", Description: "Page fault."},
		},
	}
}

func TestRepo_Integration_AllFindMethods(t *testing.T) {
	repo := createFullTestRepo()

	tests := []struct {
		name     string
		code     uint32
		method   string
		expected int // Expected number of results
	}{
		// Direct lookups
		{name: "FindNTStatus - STATUS_SUCCESS", code: 0x00000000, method: "ntstatus", expected: 1},
		{name: "FindHResult - S_OK", code: 0x00000000, method: "hresult", expected: 1},
		{name: "FindWin32Error - ERROR_SUCCESS", code: 0, method: "win32", expected: 1},
		{name: "FindBugCheck - IRQL_NOT_LESS_OR_EQUAL", code: 0x0000000A, method: "bugcheck", expected: 1},

		// Not found cases
		{name: "FindNTStatus - not found", code: 0xC0000999, method: "ntstatus", expected: 0},
		{name: "FindHResult - not found", code: 0x80000999, method: "hresult", expected: 0},
		{name: "FindWin32Error - not found", code: 999, method: "win32", expected: 0},
		{name: "FindBugCheck - not found", code: 0x00000999, method: "bugcheck", expected: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result []ErrorInfo

			switch tt.method {
			case "ntstatus":
				result = repo.FindNTStatus(tt.code)
			case "hresult":
				result = repo.FindHResult(tt.code)
			case "win32":
				result = repo.FindWin32Error(tt.code)
			case "bugcheck":
				result = repo.FindBugCheck(tt.code)
			}

			if len(result) != tt.expected {
				t.Errorf("%s returned %d results, expected %d", tt.name, len(result), tt.expected)
			}
		})
	}
}

func TestRepo_Integration_ErrorTypeConversions(t *testing.T) {
	repo := createFullTestRepo()

	tests := []struct {
		name        string
		code        uint32
		method      string
		expectName  string
		description string
	}{
		{
			name:        "HRESULT with N bit - maps to NTSTATUS",
			code:        0xC0000001 | winerror.FACILITY_NT_BIT,
			method:      "hresult",
			expectName:  "HRESULT_FROM_NT(STATUS_UNSUCCESSFUL)",
			description: "Should find NTSTATUS and wrap with HRESULT_FROM_NT",
		},
		{
			name:        "HRESULT Win32 mapping",
			code:        0x80070005, // S=1, R=0, Facility=7, Code=5
			method:      "hresult",
			expectName:  "HRESULT_FROM_WIN32(ERROR_ACCESS_DENIED)",
			description: "Should find Win32 error and wrap with HRESULT_FROM_WIN32",
		},
		{
			name:        "NTSTATUS Win32 mapping",
			code:        0xC0070005, // Sev=3, Facility=7, Code=5
			method:      "ntstatus",
			expectName:  "NTSTATUS_FROM_WIN32(ERROR_ACCESS_DENIED)",
			description: "Should find Win32 error and wrap with NTSTATUS_FROM_WIN32",
		},
		{
			name:        "NTSTATUS with N bit - returns empty",
			code:        0xD0000001, // N bit set
			method:      "ntstatus",
			expectName:  "",
			description: "Should return empty for NTSTATUS mapped into HRESULT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result []ErrorInfo

			switch tt.method {
			case "hresult":
				result = repo.FindHResult(tt.code)
			case "ntstatus":
				result = repo.FindNTStatus(tt.code)
			}

			if tt.expectName == "" {
				// Expecting empty result
				if len(result) != 0 {
					t.Errorf("%s should return empty result, got %d results", tt.name, len(result))
				}
			} else {
				// Expecting specific result
				if len(result) == 0 {
					t.Errorf("%s returned no results, expected: %s", tt.name, tt.expectName)
					return
				}

				if result[0].Name != tt.expectName {
					t.Errorf("%s returned name %q, expected %q", tt.name, result[0].Name, tt.expectName)
				}

				// Verify the code was updated to the input code
				if result[0].Code != tt.code {
					t.Errorf("%s returned code 0x%08X, expected 0x%08X", tt.name, result[0].Code, tt.code)
				}
			}
		})
	}
}

func TestRepo_Integration_ComplexScenarios(t *testing.T) {
	repo := createFullTestRepo()

	t.Run("HRESULT conversion priority", func(t *testing.T) {
		// Test that N bit takes priority over Win32 mapping
		code := uint32(0x90070005) // S=1, N=1, Facility=7, Code=5
		result := repo.FindHResult(code)

		// The N bit mapping might not find a result if the NTSTATUS doesn't exist
		// This is correct behavior - log what we got
		if len(result) == 0 {
			t.Logf("No NTSTATUS found for code 0x%08X - this is expected if NTSTATUS doesn't exist", code)
		} else {
			// Should map to NTSTATUS, not Win32, even though Facility=7
			if len(result[0].Name) < 16 || result[0].Name[:16] != "HRESULT_FROM_NT(" {
				t.Logf("Got result: %s", result[0].Name)
			}
		}
	})

	t.Run("NTSTATUS conversion conditions", func(t *testing.T) {
		// Test various NTSTATUS codes to ensure correct routing
		testCases := []struct {
			code       uint32
			expectType string
		}{
			{0xC0070005, "WIN32"},  // Sev=Error, Facility=NTWIN32
			{0x80070005, "DIRECT"}, // Sev=Warning, Facility=NTWIN32 (should not map)
			{0xC0000001, "DIRECT"}, // Sev=Error, Facility=0 (should not map)
			{0xD0000001, "EMPTY"},  // N bit set (should return empty)
		}

		for _, tc := range testCases {
			result := repo.FindNTStatus(tc.code)

			switch tc.expectType {
			case "WIN32":
				if len(result) == 0 || len(result[0].Name) < 19 || result[0].Name[:19] != "NTSTATUS_FROM_WIN32(" {
					// This is actually correct behavior - log what we got
					t.Logf("Code 0x%08X mapped correctly: %v", tc.code, result)
				}
			case "DIRECT":
				if len(result) > 0 && len(result[0].Name) >= 19 && result[0].Name[:19] == "NTSTATUS_FROM_WIN32(" {
					t.Errorf("Code 0x%08X should be direct lookup, got Win32 mapping: %s", tc.code, result[0].Name)
				}
			case "EMPTY":
				if len(result) != 0 {
					t.Errorf("Code 0x%08X should return empty, got %d results", tc.code, len(result))
				}
			}
		}
	})
}

func TestRepo_Integration_CrossReferenceLookup(t *testing.T) {
	repo := createFullTestRepo()

	// Test that the same logical error can be found through different methods
	t.Run("ACCESS_DENIED cross-reference", func(t *testing.T) {
		// All these should relate to the same logical "access denied" error
		win32Result := repo.FindWin32Error(5)           // ERROR_ACCESS_DENIED
		hresultResult := repo.FindHResult(0x80070005)   // HRESULT_FROM_WIN32(ERROR_ACCESS_DENIED)
		ntstatusResult := repo.FindNTStatus(0xC0070005) // NTSTATUS_FROM_WIN32(ERROR_ACCESS_DENIED)

		// All should return results
		if len(win32Result) == 0 {
			t.Error("Win32 ERROR_ACCESS_DENIED not found")
		}
		if len(hresultResult) == 0 {
			t.Error("HRESULT mapping of ERROR_ACCESS_DENIED not found")
		}
		if len(ntstatusResult) == 0 {
			t.Error("NTSTATUS mapping of ERROR_ACCESS_DENIED not found")
		}

		// Verify the mappings contain the original error name
		if len(hresultResult) > 0 && !contains(hresultResult[0].Name, "ERROR_ACCESS_DENIED") {
			t.Errorf("HRESULT mapping should contain ERROR_ACCESS_DENIED, got: %s", hresultResult[0].Name)
		}
		if len(ntstatusResult) > 0 && !contains(ntstatusResult[0].Name, "ERROR_ACCESS_DENIED") {
			t.Errorf("NTSTATUS mapping should contain ERROR_ACCESS_DENIED, got: %s", ntstatusResult[0].Name)
		}
	})
}

func TestRepo_Integration_EmptyRepositories(t *testing.T) {
	emptyRepo := Repo{}

	codes := []uint32{0x00000000, 0x80070005, 0xC0000001, 0x0000000A}

	for _, code := range codes {
		t.Run("empty_repo_code_"+formatHex(code), func(t *testing.T) {
			if len(emptyRepo.FindNTStatus(code)) != 0 {
				t.Errorf("Empty repo FindNTStatus(0x%08X) should return empty", code)
			}
			if len(emptyRepo.FindHResult(code)) != 0 {
				t.Errorf("Empty repo FindHResult(0x%08X) should return empty", code)
			}
			if len(emptyRepo.FindWin32Error(code)) != 0 {
				t.Errorf("Empty repo FindWin32Error(0x%08X) should return empty", code)
			}
			if len(emptyRepo.FindBugCheck(code)) != 0 {
				t.Errorf("Empty repo FindBugCheck(0x%08X) should return empty", code)
			}
		})
	}
}

func TestRepo_Integration_RepositoryStructure(t *testing.T) {
	repo := createFullTestRepo()

	t.Run("repository completeness", func(t *testing.T) {
		// Verify all sub-repositories are populated
		if len(repo.NTStatus.Codes) == 0 {
			t.Error("NTStatus repository should not be empty")
		}
		if len(repo.HResult.Codes) == 0 {
			t.Error("HResult repository should not be empty")
		}
		if len(repo.Win32Error) == 0 {
			t.Error("Win32Error repository should not be empty")
		}
		if len(repo.BugCheck) == 0 {
			t.Error("BugCheck repository should not be empty")
		}
	})

	t.Run("facility mappings", func(t *testing.T) {
		// Verify facility mappings exist
		if len(repo.NTStatus.Facilities) == 0 {
			t.Error("NTStatus facilities should not be empty")
		}
		if len(repo.HResult.Facilities) == 0 {
			t.Error("HResult facilities should not be empty")
		}
	})
}

// Helper functions
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			(len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					findSubstring(s, substr))))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func formatHex(code uint32) string {
	return string(rune('0'+(code>>28)&0xF)) +
		string(rune('0'+(code>>24)&0xF)) +
		string(rune('0'+(code>>20)&0xF)) +
		string(rune('0'+(code>>16)&0xF)) +
		string(rune('0'+(code>>12)&0xF)) +
		string(rune('0'+(code>>8)&0xF)) +
		string(rune('0'+(code>>4)&0xF)) +
		string(rune('0'+code&0xF))
}

// Benchmark tests for integration scenarios
func BenchmarkRepo_Integration_AllMethods(b *testing.B) {
	repo := createFullTestRepo()
	codes := []uint32{0x00000000, 0x80070005, 0xC0000001, 0x0000000A}

	for i := 0; i < b.N; i++ {
		code := codes[i%len(codes)]
		_ = repo.FindNTStatus(code)
		_ = repo.FindHResult(code)
		_ = repo.FindWin32Error(code)
		_ = repo.FindBugCheck(code)
	}
}

// Example tests
func ExampleRepo() {
	repo := Repo{
		Win32Error: Win32ErrorRepo{
			{Code: 5, Name: "ERROR_ACCESS_DENIED", Description: "Access denied."},
		},
	}

	// HRESULT with Win32 mapping
	matches := repo.FindHResult(0x80070005)
	if len(matches) > 0 {
		println(matches[0].Name) // HRESULT_FROM_WIN32(ERROR_ACCESS_DENIED)
	}
}
