package winerror

import (
	"testing"
)

func TestNTStatus_Sev(t *testing.T) {
	tests := []struct {
		name     string
		status   NTStatus
		expected uint8
	}{
		{
			name:     "severity success (0)",
			status:   NTStatus(0x00000000),
			expected: STATUS_SEVERITY_SUCCESS,
		},
		{
			name:     "severity informational (1)",
			status:   NTStatus(0x40000000),
			expected: STATUS_SEVERITY_INFORMATIONAL,
		},
		{
			name:     "severity warning (2)",
			status:   NTStatus(0x80000000),
			expected: STATUS_SEVERITY_WARNING,
		},
		{
			name:     "severity error (3)",
			status:   NTStatus(0xC0000000),
			expected: STATUS_SEVERITY_ERROR,
		},
		{
			name:     "severity with other bits set",
			status:   NTStatus(0xC0000001), // STATUS_UNSUCCESSFUL
			expected: STATUS_SEVERITY_ERROR,
		},
		{
			name:     "severity extraction ignores other bits",
			status:   NTStatus(0xFFFFFFFF),
			expected: STATUS_SEVERITY_ERROR, // Only severity bits matter
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.status.Sev()
			if result != tt.expected {
				t.Errorf("NTStatus(0x%08X).Sev() = %d, expected %d", uint32(tt.status), result, tt.expected)
			}
		})
	}
}

func TestNTStatus_C(t *testing.T) {
	tests := []struct {
		name     string
		status   NTStatus
		expected bool
	}{
		{
			name:     "customer bit set",
			status:   NTStatus(0x20000000),
			expected: true,
		},
		{
			name:     "customer bit not set",
			status:   NTStatus(0x00000000),
			expected: false,
		},
		{
			name:     "customer bit set with other bits",
			status:   NTStatus(0xE0000000),
			expected: true,
		},
		{
			name:     "only other bits set",
			status:   NTStatus(0xC0000000),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.status.C()
			if result != tt.expected {
				t.Errorf("NTStatus(0x%08X).C() = %v, expected %v", uint32(tt.status), result, tt.expected)
			}
		})
	}
}

func TestNTStatus_N(t *testing.T) {
	tests := []struct {
		name     string
		status   NTStatus
		expected bool
	}{
		{
			name:     "N bit set",
			status:   NTStatus(0x10000000),
			expected: true,
		},
		{
			name:     "N bit not set",
			status:   NTStatus(0x00000000),
			expected: false,
		},
		{
			name:     "N bit set with other bits",
			status:   NTStatus(0xD0000000),
			expected: true,
		},
		{
			name:     "only other bits set",
			status:   NTStatus(0xC0000000),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.status.N()
			if result != tt.expected {
				t.Errorf("NTStatus(0x%08X).N() = %v, expected %v", uint32(tt.status), result, tt.expected)
			}
		})
	}
}

func TestNTStatus_Facility(t *testing.T) {
	tests := []struct {
		name     string
		status   NTStatus
		expected uint16
	}{
		{
			name:     "facility 0",
			status:   NTStatus(0x00000000),
			expected: 0,
		},
		{
			name:     "facility NTWIN32 (7)",
			status:   NTStatus(0xC0070005),
			expected: FACILITY_NTWIN32,
		},
		{
			name:     "facility max value (0xFFF)",
			status:   NTStatus(0x0FFF0000),
			expected: 0xFFF,
		},
		{
			name:     "facility with other bits set",
			status:   NTStatus(0xC0010001), // Facility 1
			expected: 1,
		},
		{
			name:     "facility extraction ignores other bits",
			status:   NTStatus(0xFFFF0000), // All bits set except code
			expected: 0xFFF,                // Only facility bits matter
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.status.Facility()
			if result != tt.expected {
				t.Errorf("NTStatus(0x%08X).Facility() = %d, expected %d", uint32(tt.status), result, tt.expected)
			}
		})
	}
}

func TestNTStatus_Code(t *testing.T) {
	tests := []struct {
		name     string
		status   NTStatus
		expected uint16
	}{
		{
			name:     "code 0",
			status:   NTStatus(0x00000000),
			expected: 0,
		},
		{
			name:     "code 1 (STATUS_UNSUCCESSFUL)",
			status:   NTStatus(0xC0000001),
			expected: 1,
		},
		{
			name:     "code max value (0xFFFF)",
			status:   NTStatus(0x0000FFFF),
			expected: 0xFFFF,
		},
		{
			name:     "code extraction ignores other bits",
			status:   NTStatus(0xFFFFFFFF),
			expected: 0xFFFF, // Only code bits matter
		},
		{
			name:     "code 5 with other bits",
			status:   NTStatus(0xC0070005),
			expected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.status.Code()
			if result != tt.expected {
				t.Errorf("NTStatus(0x%08X).Code() = %d, expected %d", uint32(tt.status), result, tt.expected)
			}
		})
	}
}

// Test constants
func TestNTStatus_Constants(t *testing.T) {
	t.Run("STATUS_SEVERITY constants", func(t *testing.T) {
		if STATUS_SEVERITY_SUCCESS != 0 {
			t.Errorf("STATUS_SEVERITY_SUCCESS = %d, expected 0", STATUS_SEVERITY_SUCCESS)
		}
		if STATUS_SEVERITY_INFORMATIONAL != 1 {
			t.Errorf("STATUS_SEVERITY_INFORMATIONAL = %d, expected 1", STATUS_SEVERITY_INFORMATIONAL)
		}
		if STATUS_SEVERITY_WARNING != 2 {
			t.Errorf("STATUS_SEVERITY_WARNING = %d, expected 2", STATUS_SEVERITY_WARNING)
		}
		if STATUS_SEVERITY_ERROR != 3 {
			t.Errorf("STATUS_SEVERITY_ERROR = %d, expected 3", STATUS_SEVERITY_ERROR)
		}
	})

	t.Run("FACILITY_NTWIN32 value", func(t *testing.T) {
		if FACILITY_NTWIN32 != 7 {
			t.Errorf("FACILITY_NTWIN32 = %d, expected 7", FACILITY_NTWIN32)
		}
	})
}

// Test real-world NTSTATUS values
func TestNTStatus_RealWorldValues(t *testing.T) {
	tests := []struct {
		name    string
		status  NTStatus
		expSev  uint8
		expC    bool
		expN    bool
		expFac  uint16
		expCode uint16
	}{
		{
			name:   "STATUS_SUCCESS",
			status: NTStatus(0x00000000),
			expSev: STATUS_SEVERITY_SUCCESS, expC: false, expN: false,
			expFac: 0, expCode: 0,
		},
		{
			name:   "STATUS_UNSUCCESSFUL",
			status: NTStatus(0xC0000001),
			expSev: STATUS_SEVERITY_ERROR, expC: false, expN: false,
			expFac: 0, expCode: 1,
		},
		{
			name:   "STATUS_ACCESS_DENIED",
			status: NTStatus(0xC0000022),
			expSev: STATUS_SEVERITY_ERROR, expC: false, expN: false,
			expFac: 0, expCode: 0x22,
		},
		{
			name:   "NTSTATUS from Win32 error",
			status: NTStatus(0xC0070005), // Mapped from ERROR_ACCESS_DENIED
			expSev: STATUS_SEVERITY_ERROR, expC: false, expN: false,
			expFac: 7, expCode: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.status.Sev() != tt.expSev {
				t.Errorf("%s.Sev() = %d, expected %d", tt.name, tt.status.Sev(), tt.expSev)
			}
			if tt.status.C() != tt.expC {
				t.Errorf("%s.C() = %v, expected %v", tt.name, tt.status.C(), tt.expC)
			}
			if tt.status.N() != tt.expN {
				t.Errorf("%s.N() = %v, expected %v", tt.name, tt.status.N(), tt.expN)
			}
			if tt.status.Facility() != tt.expFac {
				t.Errorf("%s.Facility() = %d, expected %d", tt.name, tt.status.Facility(), tt.expFac)
			}
			if tt.status.Code() != tt.expCode {
				t.Errorf("%s.Code() = %d, expected %d", tt.name, tt.status.Code(), tt.expCode)
			}
		})
	}
}

// Edge cases and boundary conditions
func TestNTStatus_EdgeCases(t *testing.T) {
	tests := []struct {
		name   string
		status NTStatus
		desc   string
	}{
		{
			name:   "zero value",
			status: NTStatus(0),
			desc:   "All bits zero should work correctly",
		},
		{
			name:   "max uint32",
			status: NTStatus(0xFFFFFFFF),
			desc:   "All bits set should work correctly",
		},
		{
			name:   "only severity bits",
			status: NTStatus(0xC0000000),
			desc:   "Only severity bits set",
		},
		{
			name:   "only facility bits",
			status: NTStatus(0x0FFF0000),
			desc:   "Only facility bits set",
		},
		{
			name:   "only code bits",
			status: NTStatus(0x0000FFFF),
			desc:   "Only code bits set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just ensure no panics occur and methods return reasonable values
			sev := tt.status.Sev()
			c := tt.status.C()
			n := tt.status.N()
			fac := tt.status.Facility()
			code := tt.status.Code()

			// Basic sanity checks
			if sev > 3 {
				t.Errorf("Sev() returned invalid severity: %d", sev)
			}
			if fac > 0xFFF {
				t.Errorf("Facility() returned invalid facility: %d", fac)
			}
			if code > 0xFFFF {
				t.Errorf("Code() returned invalid code: %d", code)
			}

			t.Logf("%s: Sev=%d, C=%v, N=%v, Fac=%d, Code=%d", tt.desc, sev, c, n, fac, code)
		})
	}
}

// Benchmark tests
func BenchmarkNTStatus_Sev(b *testing.B) {
	status := NTStatus(0xC0000001)
	for i := 0; i < b.N; i++ {
		_ = status.Sev()
	}
}

func BenchmarkNTStatus_Facility(b *testing.B) {
	status := NTStatus(0xC0070005)
	for i := 0; i < b.N; i++ {
		_ = status.Facility()
	}
}

func BenchmarkNTStatus_Code(b *testing.B) {
	status := NTStatus(0xC0000001)
	for i := 0; i < b.N; i++ {
		_ = status.Code()
	}
}

// Test type conversion and casting
func TestNTStatus_TypeConversion(t *testing.T) {
	tests := []struct {
		name     string
		input    uint32
		expected NTStatus
	}{
		{
			name:     "convert from uint32 zero",
			input:    0,
			expected: NTStatus(0),
		},
		{
			name:     "convert from uint32 max",
			input:    0xFFFFFFFF,
			expected: NTStatus(0xFFFFFFFF),
		},
		{
			name:     "convert from STATUS_UNSUCCESSFUL",
			input:    0xC0000001,
			expected: NTStatus(0xC0000001),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NTStatus(tt.input)
			if result != tt.expected {
				t.Errorf("NTStatus(%d) = %v, expected %v", tt.input, result, tt.expected)
			}
			// Verify round-trip conversion
			if uint32(result) != tt.input {
				t.Errorf("uint32(NTStatus(%d)) = %d, expected %d", tt.input, uint32(result), tt.input)
			}
		})
	}
}

// Test bit manipulation patterns
func TestNTStatus_BitPatterns(t *testing.T) {
	tests := []struct {
		name        string
		status      NTStatus
		description string
	}{
		{
			name:        "all severity levels",
			status:      NTStatus(0xC0000000), // Error severity
			description: "Error severity pattern",
		},
		{
			name:        "customer and N bits",
			status:      NTStatus(0x30000000), // C + N bits
			description: "Customer and N bits combination",
		},
		{
			name:        "Win32 error mapping",
			status:      NTStatus(0xC0070000), // Error + FACILITY_NTWIN32
			description: "Win32 error mapped to NTSTATUS",
		},
		{
			name:        "informational with facility",
			status:      NTStatus(0x40010000), // Informational + facility 1
			description: "Informational status with facility",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that all methods work without panic
			_ = tt.status.Sev()
			_ = tt.status.C()
			_ = tt.status.N()
			_ = tt.status.Facility()
			_ = tt.status.Code()
			t.Logf("%s: 0x%08X", tt.description, uint32(tt.status))
		})
	}
}

// Test severity level boundaries
func TestNTStatus_SeverityBoundaries(t *testing.T) {
	tests := []struct {
		name     string
		status   NTStatus
		expected uint8
	}{
		{
			name:     "severity 0 boundary",
			status:   NTStatus(0x00000000),
			expected: STATUS_SEVERITY_SUCCESS,
		},
		{
			name:     "severity 1 boundary",
			status:   NTStatus(0x40000000),
			expected: STATUS_SEVERITY_INFORMATIONAL,
		},
		{
			name:     "severity 2 boundary",
			status:   NTStatus(0x80000000),
			expected: STATUS_SEVERITY_WARNING,
		},
		{
			name:     "severity 3 boundary",
			status:   NTStatus(0xC0000000),
			expected: STATUS_SEVERITY_ERROR,
		},
		{
			name:     "severity with other bits",
			status:   NTStatus(0xFFFFFFFF),  // All bits set
			expected: STATUS_SEVERITY_ERROR, // Should still extract severity correctly
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.status.Sev()
			if result != tt.expected {
				t.Errorf("NTStatus(0x%08X).Sev() = %d, expected %d", uint32(tt.status), result, tt.expected)
			}
		})
	}
}

// Example tests for documentation
func ExampleNTStatus_Sev() {
	status := NTStatus(0xC0000001) // STATUS_UNSUCCESSFUL
	println(status.Sev())          // 3 (STATUS_SEVERITY_ERROR)
}

func ExampleNTStatus_Facility() {
	status := NTStatus(0xC0070005) // NTSTATUS from Win32 error
	println(status.Facility())     // 7 (FACILITY_NTWIN32)
}

func ExampleNTStatus_Code() {
	status := NTStatus(0xC0000001) // STATUS_UNSUCCESSFUL
	println(status.Code())         // 1
}
