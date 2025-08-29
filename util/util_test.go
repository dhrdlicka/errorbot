package util

import (
	"reflect"
	"testing"
)

func TestParseCode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []uint32
		wantErr  bool
	}{
		// Hex prefix tests (0x/0X)
		{
			name:     "hex with lowercase 0x prefix",
			input:    "0x1234",
			expected: []uint32{0x1234},
			wantErr:  false,
		},
		{
			name:     "hex with uppercase 0X prefix",
			input:    "0X1234",
			expected: []uint32{0x1234},
			wantErr:  false,
		},
		{
			name:     "hex with 0x prefix - zero value",
			input:    "0x0",
			expected: []uint32{0},
			wantErr:  false,
		},
		{
			name:     "hex with 0x prefix - max uint32",
			input:    "0xFFFFFFFF",
			expected: []uint32{0xFFFFFFFF},
			wantErr:  false,
		},
		{
			name:     "hex with 0x prefix - lowercase letters",
			input:    "0xabcdef",
			expected: []uint32{0xabcdef},
			wantErr:  false,
		},
		{
			name:     "hex with 0x prefix - uppercase letters",
			input:    "0xABCDEF",
			expected: []uint32{0xABCDEF},
			wantErr:  false,
		},
		{
			name:     "hex with 0x prefix - invalid hex character",
			input:    "0x123G",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "hex with 0x prefix - overflow uint32",
			input:    "0x100000000",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "hex with 0x prefix - empty after prefix",
			input:    "0x",
			expected: nil,
			wantErr:  true,
		},

		// Negative number tests
		{
			name:     "negative decimal number",
			input:    "-1",
			expected: []uint32{0xFFFFFFFF}, // -1 as uint32
			wantErr:  false,
		},
		{
			name:     "negative decimal - large negative",
			input:    "-2147483648",
			expected: []uint32{0x80000000}, // -2147483648 as uint32
			wantErr:  false,
		},
		{
			name:     "negative decimal - small negative",
			input:    "-123",
			expected: []uint32{4294967173}, // -123 as uint32 (0xFFFFFF85)
			wantErr:  false,
		},
		{
			name:     "negative decimal - overflow int32",
			input:    "-2147483649",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "negative decimal - invalid format",
			input:    "-abc",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "negative decimal - just minus sign",
			input:    "-",
			expected: nil,
			wantErr:  true,
		},

		// Ambiguous codes (no prefix) - should try both hex and decimal
		{
			name:     "ambiguous code - valid as both hex and decimal",
			input:    "123",
			expected: []uint32{0x123, 123}, // Both hex and decimal interpretations
			wantErr:  false,
		},
		{
			name:     "ambiguous code - valid only as decimal",
			input:    "999",
			expected: []uint32{0x999, 999}, // 999 is valid hex (0x999=2457) and decimal
			wantErr:  false,
		},
		{
			name:     "ambiguous code - valid only as hex",
			input:    "ABC",
			expected: []uint32{0xABC}, // Only hex
			wantErr:  false,
		},
		{
			name:     "ambiguous code - zero",
			input:    "0",
			expected: []uint32{0}, // Same value for both hex and decimal, should be compacted
			wantErr:  false,
		},
		{
			name:     "ambiguous code - invalid for both",
			input:    "XYZ",
			expected: nil,
			wantErr:  true,
		},

		{
			name:     "ambiguous code - max uint32 as decimal",
			input:    "4294967295",
			expected: []uint32{4294967295},
			wantErr:  false,
		},
		{
			name:     "ambiguous code - overflow uint32 as decimal",
			input:    "4294967296",
			expected: nil,
			wantErr:  true,
		},

		// Edge cases
		{
			name:     "single character hex",
			input:    "F",
			expected: []uint32{0xF}, // Only valid as hex
			wantErr:  false,
		},
		{
			name:     "single character decimal",
			input:    "9",
			expected: []uint32{9}, // Same value for hex and decimal, compacted to one
			wantErr:  false,
		},
		{
			name:     "leading zeros in hex prefix",
			input:    "0x0001",
			expected: []uint32{1},
			wantErr:  false,
		},
		{
			name:     "leading zeros in ambiguous",
			input:    "0001",
			expected: []uint32{1}, // Both hex and decimal parse to 1, compacted
			wantErr:  false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseCode(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseCode(%q) expected error, but got none", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("ParseCode(%q) unexpected error: %v", tt.input, err)
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ParseCode(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestBoolToInt(t *testing.T) {
	tests := []struct {
		name     string
		input    bool
		expected int
	}{
		{
			name:     "true converts to 1",
			input:    true,
			expected: 1,
		},
		{
			name:     "false converts to 0",
			input:    false,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BoolToInt(tt.input)
			if result != tt.expected {
				t.Errorf("BoolToInt(%v) = %d, expected %d", tt.input, result, tt.expected)
			}
		})
	}
}

// Benchmark tests for performance analysis
func BenchmarkParseCode_HexPrefix(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = ParseCode("0x1234ABCD")
	}
}

func BenchmarkParseCode_Negative(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = ParseCode("-12345")
	}
}

func BenchmarkParseCode_Ambiguous(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = ParseCode("12345")
	}
}

func BenchmarkBoolToInt_True(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = BoolToInt(true)
	}
}

func BenchmarkBoolToInt_False(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = BoolToInt(false)
	}
}

// Example tests for documentation
func ExampleParseCode() {
	// Hex with prefix
	codes, _ := ParseCode("0x1234")
	println(codes[0]) // 4660

	// Negative decimal
	codes, _ = ParseCode("-1")
	println(codes[0]) // 4294967295 (0xFFFFFFFF)

	// Ambiguous code (valid as both hex and decimal)
	codes, _ = ParseCode("123")
	println(len(codes)) // 2 (both interpretations)
}

func ExampleBoolToInt() {
	println(BoolToInt(true))  // 1
	println(BoolToInt(false)) // 0
}
