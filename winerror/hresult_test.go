package winerror

import (
	"testing"
)

func TestHResult_S(t *testing.T) {
	tests := []struct {
		name     string
		hresult  HResult
		expected bool
	}{
		{
			name:     "severity bit set (failure)",
			hresult:  HResult(0x80000000),
			expected: true,
		},
		{
			name:     "severity bit not set (success)",
			hresult:  HResult(0x00000000),
			expected: false,
		},
		{
			name:     "severity bit set with other bits",
			hresult:  HResult(0x80070005), // E_ACCESSDENIED
			expected: true,
		},
		{
			name:     "severity bit not set with other bits",
			hresult:  HResult(0x00000001),
			expected: false,
		},
		{
			name:     "max uint32 value",
			hresult:  HResult(0xFFFFFFFF),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.hresult.S()
			if result != tt.expected {
				t.Errorf("HResult(0x%08X).S() = %v, expected %v", uint32(tt.hresult), result, tt.expected)
			}
		})
	}
}

func TestHResult_R(t *testing.T) {
	tests := []struct {
		name     string
		hresult  HResult
		expected bool
	}{
		{
			name:     "reserved R bit set",
			hresult:  HResult(0x40000000),
			expected: true,
		},
		{
			name:     "reserved R bit not set",
			hresult:  HResult(0x00000000),
			expected: false,
		},
		{
			name:     "reserved R bit set with other bits",
			hresult:  HResult(0xC0000000),
			expected: true,
		},
		{
			name:     "only other bits set",
			hresult:  HResult(0x80000000),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.hresult.R()
			if result != tt.expected {
				t.Errorf("HResult(0x%08X).R() = %v, expected %v", uint32(tt.hresult), result, tt.expected)
			}
		})
	}
}

func TestHResult_C(t *testing.T) {
	tests := []struct {
		name     string
		hresult  HResult
		expected bool
	}{
		{
			name:     "customer bit set",
			hresult:  HResult(0x20000000),
			expected: true,
		},
		{
			name:     "customer bit not set",
			hresult:  HResult(0x00000000),
			expected: false,
		},
		{
			name:     "customer bit set with other bits",
			hresult:  HResult(0xA0000000),
			expected: true,
		},
		{
			name:     "only other bits set",
			hresult:  HResult(0x80000000),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.hresult.C()
			if result != tt.expected {
				t.Errorf("HResult(0x%08X).C() = %v, expected %v", uint32(tt.hresult), result, tt.expected)
			}
		})
	}
}

func TestHResult_N(t *testing.T) {
	tests := []struct {
		name     string
		hresult  HResult
		expected bool
	}{
		{
			name:     "NT bit set",
			hresult:  HResult(0x10000000),
			expected: true,
		},
		{
			name:     "NT bit not set",
			hresult:  HResult(0x00000000),
			expected: false,
		},
		{
			name:     "NT bit set with other bits",
			hresult:  HResult(0x90000000),
			expected: true,
		},
		{
			name:     "only other bits set",
			hresult:  HResult(0x80000000),
			expected: false,
		},
		{
			name:     "FACILITY_NT_BIT constant value",
			hresult:  HResult(FACILITY_NT_BIT),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.hresult.N()
			if result != tt.expected {
				t.Errorf("HResult(0x%08X).N() = %v, expected %v", uint32(tt.hresult), result, tt.expected)
			}
		})
	}
}

func TestHResult_X(t *testing.T) {
	tests := []struct {
		name     string
		hresult  HResult
		expected bool
	}{
		{
			name:     "X bit set",
			hresult:  HResult(0x08000000),
			expected: true,
		},
		{
			name:     "X bit not set",
			hresult:  HResult(0x00000000),
			expected: false,
		},
		{
			name:     "X bit set with other bits",
			hresult:  HResult(0x88000000),
			expected: true,
		},
		{
			name:     "only other bits set",
			hresult:  HResult(0x80000000),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.hresult.X()
			if result != tt.expected {
				t.Errorf("HResult(0x%08X).X() = %v, expected %v", uint32(tt.hresult), result, tt.expected)
			}
		})
	}
}

func TestHResult_Facility(t *testing.T) {
	tests := []struct {
		name     string
		hresult  HResult
		expected uint16
	}{
		{
			name:     "facility 0",
			hresult:  HResult(0x00000000),
			expected: 0,
		},
		{
			name:     "facility WIN32 (7)",
			hresult:  HResult(0x80070005), // E_ACCESSDENIED
			expected: FACILITY_WIN32,
		},
		{
			name:     "facility max value (0x7FF)",
			hresult:  HResult(0x07FF0000),
			expected: 0x7FF,
		},
		{
			name:     "facility with other bits set",
			hresult:  HResult(0x80010001), // Facility 1
			expected: 1,
		},
		{
			name:     "facility extraction ignores other bits",
			hresult:  HResult(0xFFFF0000), // All bits set except code
			expected: 0x7FF,               // Only facility bits matter
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.hresult.Facility()
			if result != tt.expected {
				t.Errorf("HResult(0x%08X).Facility() = %d, expected %d", uint32(tt.hresult), result, tt.expected)
			}
		})
	}
}

func TestHResult_Code(t *testing.T) {
	tests := []struct {
		name     string
		hresult  HResult
		expected uint16
	}{
		{
			name:     "code 0",
			hresult:  HResult(0x00000000),
			expected: 0,
		},
		{
			name:     "code 5 (ACCESS_DENIED)",
			hresult:  HResult(0x80070005),
			expected: 5,
		},
		{
			name:     "code max value (0xFFFF)",
			hresult:  HResult(0x0000FFFF),
			expected: 0xFFFF,
		},
		{
			name:     "code extraction ignores other bits",
			hresult:  HResult(0xFFFFFFFF),
			expected: 0xFFFF, // Only code bits matter
		},
		{
			name:     "code 1 with other bits",
			hresult:  HResult(0x80070001),
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.hresult.Code()
			if result != tt.expected {
				t.Errorf("HResult(0x%08X).Code() = %d, expected %d", uint32(tt.hresult), result, tt.expected)
			}
		})
	}
}

// Test constants
func TestHResult_Constants(t *testing.T) {
	t.Run("FACILITY_WIN32 value", func(t *testing.T) {
		if FACILITY_WIN32 != 7 {
			t.Errorf("FACILITY_WIN32 = %d, expected 7", FACILITY_WIN32)
		}
	})

	t.Run("FACILITY_NT_BIT value", func(t *testing.T) {
		if FACILITY_NT_BIT != 0x10000000 {
			t.Errorf("FACILITY_NT_BIT = 0x%08X, expected 0x10000000", FACILITY_NT_BIT)
		}
	})
}

// Test real-world HRESULT values
func TestHResult_RealWorldValues(t *testing.T) {
	tests := []struct {
		name    string
		hresult HResult
		expS    bool
		expR    bool
		expC    bool
		expN    bool
		expX    bool
		expFac  uint16
		expCode uint16
	}{
		{
			name:    "S_OK",
			hresult: HResult(0x00000000),
			expS:    false, expR: false, expC: false, expN: false, expX: false,
			expFac: 0, expCode: 0,
		},
		{
			name:    "E_ACCESSDENIED",
			hresult: HResult(0x80070005),
			expS:    true, expR: false, expC: false, expN: false, expX: false,
			expFac: 7, expCode: 5,
		},
		{
			name:    "E_INVALIDARG",
			hresult: HResult(0x80070057),
			expS:    true, expR: false, expC: false, expN: false, expX: false,
			expFac: 7, expCode: 87,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.hresult.S() != tt.expS {
				t.Errorf("%s.S() = %v, expected %v", tt.name, tt.hresult.S(), tt.expS)
			}
			if tt.hresult.R() != tt.expR {
				t.Errorf("%s.R() = %v, expected %v", tt.name, tt.hresult.R(), tt.expR)
			}
			if tt.hresult.C() != tt.expC {
				t.Errorf("%s.C() = %v, expected %v", tt.name, tt.hresult.C(), tt.expC)
			}
			if tt.hresult.N() != tt.expN {
				t.Errorf("%s.N() = %v, expected %v", tt.name, tt.hresult.N(), tt.expN)
			}
			if tt.hresult.X() != tt.expX {
				t.Errorf("%s.X() = %v, expected %v", tt.name, tt.hresult.X(), tt.expX)
			}
			if tt.hresult.Facility() != tt.expFac {
				t.Errorf("%s.Facility() = %d, expected %d", tt.name, tt.hresult.Facility(), tt.expFac)
			}
			if tt.hresult.Code() != tt.expCode {
				t.Errorf("%s.Code() = %d, expected %d", tt.name, tt.hresult.Code(), tt.expCode)
			}
		})
	}
}

// Benchmark tests
func BenchmarkHResult_S(b *testing.B) {
	hr := HResult(0x80070005)
	for i := 0; i < b.N; i++ {
		_ = hr.S()
	}
}

func BenchmarkHResult_Facility(b *testing.B) {
	hr := HResult(0x80070005)
	for i := 0; i < b.N; i++ {
		_ = hr.Facility()
	}
}

func BenchmarkHResult_Code(b *testing.B) {
	hr := HResult(0x80070005)
	for i := 0; i < b.N; i++ {
		_ = hr.Code()
	}
}

// Edge cases and boundary conditions
func TestHResult_EdgeCases(t *testing.T) {
	tests := []struct {
		name string
		hr   HResult
		desc string
	}{
		{
			name: "zero value",
			hr:   HResult(0),
			desc: "All bits zero should work correctly",
		},
		{
			name: "max uint32",
			hr:   HResult(0xFFFFFFFF),
			desc: "All bits set should work correctly",
		},
		{
			name: "only severity bit",
			hr:   HResult(0x80000000),
			desc: "Only severity bit set",
		},
		{
			name: "only facility bits",
			hr:   HResult(0x07FF0000),
			desc: "Only facility bits set",
		},
		{
			name: "only code bits",
			hr:   HResult(0x0000FFFF),
			desc: "Only code bits set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just ensure no panics occur and methods return reasonable values
			s := tt.hr.S()
			r := tt.hr.R()
			c := tt.hr.C()
			n := tt.hr.N()
			x := tt.hr.X()
			fac := tt.hr.Facility()
			code := tt.hr.Code()

			// Basic sanity checks
			if fac > 0x7FF {
				t.Errorf("Facility() returned invalid facility: %d", fac)
			}
			if code > 0xFFFF {
				t.Errorf("Code() returned invalid code: %d", code)
			}

			t.Logf("%s: S=%v, R=%v, C=%v, N=%v, X=%v, Fac=%d, Code=%d", tt.desc, s, r, c, n, x, fac, code)
		})
	}
}

// Test type conversion and casting
func TestHResult_TypeConversion(t *testing.T) {
	tests := []struct {
		name     string
		input    uint32
		expected HResult
	}{
		{
			name:     "convert from uint32 zero",
			input:    0,
			expected: HResult(0),
		},
		{
			name:     "convert from uint32 max",
			input:    0xFFFFFFFF,
			expected: HResult(0xFFFFFFFF),
		},
		{
			name:     "convert from E_ACCESSDENIED",
			input:    0x80070005,
			expected: HResult(0x80070005),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HResult(tt.input)
			if result != tt.expected {
				t.Errorf("HResult(%d) = %v, expected %v", tt.input, result, tt.expected)
			}
			// Verify round-trip conversion
			if uint32(result) != tt.input {
				t.Errorf("uint32(HResult(%d)) = %d, expected %d", tt.input, uint32(result), tt.input)
			}
		})
	}
}

// Test bit manipulation patterns
func TestHResult_BitPatterns(t *testing.T) {
	tests := []struct {
		name        string
		hr          HResult
		description string
	}{
		{
			name:        "all reserved bits set",
			hr:          HResult(0x78000000), // R, C, N, X bits
			description: "Reserved bits combination",
		},
		{
			name:        "severity and facility only",
			hr:          HResult(0x80070000), // S bit + FACILITY_WIN32
			description: "Common pattern for Win32 errors",
		},
		{
			name:        "NT status mapping",
			hr:          HResult(0x90000000), // S + N bits
			description: "NTSTATUS mapped to HRESULT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that all methods work without panic
			_ = tt.hr.S()
			_ = tt.hr.R()
			_ = tt.hr.C()
			_ = tt.hr.N()
			_ = tt.hr.X()
			_ = tt.hr.Facility()
			_ = tt.hr.Code()
			t.Logf("%s: 0x%08X", tt.description, uint32(tt.hr))
		})
	}
}

// Example tests for documentation
func ExampleHResult_S() {
	hr := HResult(0x80070005) // E_ACCESSDENIED
	println(hr.S())           // true (failure)
}

func ExampleHResult_Facility() {
	hr := HResult(0x80070005) // E_ACCESSDENIED
	println(hr.Facility())    // 7 (FACILITY_WIN32)
}

func ExampleHResult_Code() {
	hr := HResult(0x80070005) // E_ACCESSDENIED
	println(hr.Code())        // 5 (ERROR_ACCESS_DENIED)
}
