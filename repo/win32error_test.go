package repo

import (
	"reflect"
	"testing"
)

// createTestWin32ErrorRepo creates a test repository with sample Win32 error data
func createTestWin32ErrorRepo() Win32ErrorRepo {
	return Win32ErrorRepo{
		{Code: 0, Name: "ERROR_SUCCESS", Description: "The operation completed successfully."},
		{Code: 1, Name: "ERROR_INVALID_FUNCTION", Description: "Incorrect function."},
		{Code: 2, Name: "ERROR_FILE_NOT_FOUND", Description: "The system cannot find the file specified."},
		{Code: 5, Name: "ERROR_ACCESS_DENIED", Description: "Access is denied."},
		{Code: 87, Name: "ERROR_INVALID_PARAMETER", Description: "The parameter is incorrect."},
		{Code: 1223, Name: "ERROR_CANCELLED", Description: "The operation was canceled by the user."},
	}
}

func TestWin32ErrorRepo_FindCode(t *testing.T) {
	repo := createTestWin32ErrorRepo()

	tests := []struct {
		name     string
		code     uint32
		expected []ErrorInfo
	}{
		{
			name: "find ERROR_SUCCESS",
			code: 0,
			expected: []ErrorInfo{
				{Code: 0, Name: "ERROR_SUCCESS", Description: "The operation completed successfully."},
			},
		},
		{
			name: "find ERROR_ACCESS_DENIED",
			code: 5,
			expected: []ErrorInfo{
				{Code: 5, Name: "ERROR_ACCESS_DENIED", Description: "Access is denied."},
			},
		},
		{
			name: "find ERROR_INVALID_PARAMETER",
			code: 87,
			expected: []ErrorInfo{
				{Code: 87, Name: "ERROR_INVALID_PARAMETER", Description: "The parameter is incorrect."},
			},
		},
		{
			name:     "find non-existent code",
			code:     999,
			expected: []ErrorInfo{},
		},
		{
			name: "find high value code",
			code: 1223,
			expected: []ErrorInfo{
				{Code: 1223, Name: "ERROR_CANCELLED", Description: "The operation was canceled by the user."},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := repo.FindCode(tt.code)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Win32ErrorRepo.FindCode(%d) = %v, expected %v", tt.code, result, tt.expected)
			}
		})
	}
}

func TestRepo_FindWin32Error(t *testing.T) {
	repo := Repo{
		Win32Error: createTestWin32ErrorRepo(),
	}

	tests := []struct {
		name     string
		code     uint32
		expected []ErrorInfo
	}{
		{
			name: "find through Repo.FindWin32Error",
			code: 5,
			expected: []ErrorInfo{
				{Code: 5, Name: "ERROR_ACCESS_DENIED", Description: "Access is denied."},
			},
		},
		{
			name:     "find non-existent through Repo.FindWin32Error",
			code:     999,
			expected: []ErrorInfo{},
		},
		{
			name: "find zero code through Repo.FindWin32Error",
			code: 0,
			expected: []ErrorInfo{
				{Code: 0, Name: "ERROR_SUCCESS", Description: "The operation completed successfully."},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := repo.FindWin32Error(tt.code)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Repo.FindWin32Error(%d) = %v, expected %v", tt.code, result, tt.expected)
			}
		})
	}
}

// Test edge cases and boundary conditions
func TestWin32ErrorRepo_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		repo     Win32ErrorRepo
		code     uint32
		expected int // Expected number of matches
	}{
		{
			name:     "empty repository",
			repo:     Win32ErrorRepo{},
			code:     1,
			expected: 0,
		},
		{
			name:     "nil repository",
			repo:     nil,
			code:     1,
			expected: 0,
		},
		{
			name: "repository with duplicate codes",
			repo: Win32ErrorRepo{
				{Code: 5, Name: "ERROR_ACCESS_DENIED_1", Description: "First access denied."},
				{Code: 5, Name: "ERROR_ACCESS_DENIED_2", Description: "Second access denied."},
			},
			code:     5,
			expected: 2,
		},
		{
			name: "max uint32 code",
			repo: Win32ErrorRepo{
				{Code: 0xFFFFFFFF, Name: "MAX_ERROR", Description: "Maximum error code."},
			},
			code:     0xFFFFFFFF,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.repo.FindCode(tt.code)
			if len(result) != tt.expected {
				t.Errorf("Win32ErrorRepo.FindCode() returned %d matches, expected %d", len(result), tt.expected)
			}
		})
	}
}

// Test type conversion and interface compliance
func TestWin32ErrorRepo_TypeConversion(t *testing.T) {
	repo := createTestWin32ErrorRepo()

	// Test that Win32ErrorRepo can be used as []ErrorInfo
	var errorSlice []ErrorInfo = repo
	if len(errorSlice) != len(repo) {
		t.Errorf("Type conversion failed: expected length %d, got %d", len(repo), len(errorSlice))
	}

	// Test that individual elements are accessible
	if errorSlice[0].Code != repo[0].Code {
		t.Errorf("Element access failed: expected code %d, got %d", repo[0].Code, errorSlice[0].Code)
	}
}

// Test repository operations
func TestWin32ErrorRepo_Operations(t *testing.T) {
	repo := createTestWin32ErrorRepo()

	t.Run("repository length", func(t *testing.T) {
		expectedLength := 6
		if len(repo) != expectedLength {
			t.Errorf("Repository length = %d, expected %d", len(repo), expectedLength)
		}
	})

	t.Run("repository indexing", func(t *testing.T) {
		if repo[0].Code != 0 {
			t.Errorf("First element code = %d, expected 0", repo[0].Code)
		}
		if repo[len(repo)-1].Code != 1223 {
			t.Errorf("Last element code = %d, expected 1223", repo[len(repo)-1].Code)
		}
	})

	t.Run("repository iteration", func(t *testing.T) {
		count := 0
		for _, err := range repo {
			if err.Code == 0 || err.Code == 1 || err.Code == 2 || err.Code == 5 || err.Code == 87 || err.Code == 1223 {
				count++
			}
		}
		if count != 6 {
			t.Errorf("Iteration found %d valid codes, expected 6", count)
		}
	})
}

// Test common Win32 error codes
func TestWin32ErrorRepo_CommonCodes(t *testing.T) {
	repo := createTestWin32ErrorRepo()

	commonCodes := []struct {
		code uint32
		name string
	}{
		{0, "ERROR_SUCCESS"},
		{1, "ERROR_INVALID_FUNCTION"},
		{2, "ERROR_FILE_NOT_FOUND"},
		{5, "ERROR_ACCESS_DENIED"},
		{87, "ERROR_INVALID_PARAMETER"},
	}

	for _, cc := range commonCodes {
		t.Run("common_code_"+cc.name, func(t *testing.T) {
			result := repo.FindCode(cc.code)
			if len(result) == 0 {
				t.Errorf("Expected to find code %d (%s), but got no results", cc.code, cc.name)
				return
			}
			if result[0].Name != cc.name {
				t.Errorf("Expected name %s, got %s", cc.name, result[0].Name)
			}
		})
	}
}

// Benchmark tests
func BenchmarkWin32ErrorRepo_FindCode(b *testing.B) {
	repo := createTestWin32ErrorRepo()
	
	for i := 0; i < b.N; i++ {
		_ = repo.FindCode(5)
	}
}

func BenchmarkWin32ErrorRepo_FindCode_NotFound(b *testing.B) {
	repo := createTestWin32ErrorRepo()
	
	for i := 0; i < b.N; i++ {
		_ = repo.FindCode(999)
	}
}

func BenchmarkRepo_FindWin32Error(b *testing.B) {
	repo := Repo{Win32Error: createTestWin32ErrorRepo()}
	
	for i := 0; i < b.N; i++ {
		_ = repo.FindWin32Error(5)
	}
}

// Example tests for documentation
func ExampleWin32ErrorRepo_FindCode() {
	repo := Win32ErrorRepo{
		{Code: 5, Name: "ERROR_ACCESS_DENIED", Description: "Access is denied."},
	}
	
	matches := repo.FindCode(5)
	if len(matches) > 0 {
		println(matches[0].Name) // ERROR_ACCESS_DENIED
	}
}

func ExampleRepo_FindWin32Error() {
	repo := Repo{
		Win32Error: Win32ErrorRepo{
			{Code: 5, Name: "ERROR_ACCESS_DENIED", Description: "Access is denied."},
		},
	}
	
	matches := repo.FindWin32Error(5)
	if len(matches) > 0 {
		println(matches[0].Name) // ERROR_ACCESS_DENIED
	}
}
