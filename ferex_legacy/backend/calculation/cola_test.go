package calculation

import (
	"reflect"
	"strings" // Added for strings.Contains
	"testing"

	"ferex/backend/models" // Import models package
	"github.com/google/go-cmp/cmp"
)

func TestGetCorrectFERSCOLA(t *testing.T) {
	testCases := []struct {
		name              string
		cpiW_IncreaseRate float64
		expectedCOLARate  float64
	}{
		{"CPI 1.5%", 0.015, 0.015},
		{"CPI 2.0%", 0.02, 0.02},
		{"CPI 2.5%", 0.025, 0.02}, // FERS COLA is 2% if CPI-W is 2-3%
		{"CPI 3.0%", 0.03, 0.02}, // FERS COLA is 2% if CPI-W is 2-3%
		{"CPI 4.5%", 0.045, 0.035}, // FERS COLA is CPI-W - 1% if CPI-W > 3%
		{"CPI 0%", 0.0, 0.0},
		{"CPI -1% (Deflation)", -0.01, 0.0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := getCorrectFERSCOLA(tc.cpiW_IncreaseRate)
			const epsilon = 0.0000001 // A small tolerance for float comparison
			if diff := got - tc.expectedCOLARate; diff < -epsilon || diff > epsilon {
				t.Errorf("getCorrectFERSCOLA(%f) = %f; want %f (difference: %f)", tc.cpiW_IncreaseRate, got, tc.expectedCOLARate, diff)
			}
		})
	}
}

func TestProjectPensionWithCOLA(t *testing.T) {
	testCases := []struct {
		name                 string
		input                models.PensionCOLAInput
		expectedResult       models.PensionCOLAResult // Or just expectedProjection []float64 and expectedNotesContain string
		expectedNotesContain string // Kept separate for easier partial string matching in notes
	}{
		{
			name: "CSRS, 2% COLA, 3 years, starts at retirement age 60",
			input: models.PensionCOLAInput{
				InitialPension:     50000,
				PensionType:        "CSRS",
				AssumedCPIWRate:    0.02,
				RetirementAge:      60,
				ProjectionStartAge: 60,
				YearsToProject:     3,
			},
			expectedResult: models.PensionCOLAResult{
				ProjectedPension: []float64{51000.00, 52020.00, 53060.40},
				// Notes can be checked separately if complex or for partial matches
			},
			expectedNotesContain: "CSRS: Full COLA applied each year.", // Updated to match new note phrasing
		},
		{
			name: "FERS, 4.5% CPI-W, 5 years, projection starts at age 57",
			input: models.PensionCOLAInput{
				InitialPension:     30000,
				PensionType:        "FERS",
				AssumedCPIWRate:    0.045,
				RetirementAge:      57, // Assuming retirement at projection start for simplicity here
				ProjectionStartAge: 57,
				YearsToProject:     5,
			},
			expectedResult: models.PensionCOLAResult{
				ProjectedPension: []float64{30000.00, 30000.00, 30000.00, 30000.00, 30000.00},
			},
			expectedNotesContain: "FERS: No COLA (age < 62).",
		},
		{
			name: "FERS, 4.5% CPI-W, 3 years, projection starts at age 61",
			input: models.PensionCOLAInput{
				InitialPension:     30000,
				PensionType:        "FERS",
				AssumedCPIWRate:    0.045,
				RetirementAge:      57,
				ProjectionStartAge: 61,
				YearsToProject:     3,
			},
			expectedResult: models.PensionCOLAResult{
				ProjectedPension: []float64{30000.00, 31050.00, 32136.75},
			},
			expectedNotesContain: "FERS: COLA applied (age >= 62).",
		},
		{
			name: "FERS, 2.5% CPI-W, 3 years, projection starts at age 61",
			input: models.PensionCOLAInput{
				InitialPension:     40000,
				PensionType:        "FERS",
				AssumedCPIWRate:    0.025,
				RetirementAge:      60,
				ProjectionStartAge: 61,
				YearsToProject:     3,
			},
			expectedResult: models.PensionCOLAResult{
				ProjectedPension: []float64{40000.00, 40800.00, 41616.00},
			},
			expectedNotesContain: "FERS: COLA applied (age >= 62).",
		},
		{
			name: "FERS, 1.5% CPI-W, 3 years, projection starts at age 61",
			input: models.PensionCOLAInput{
				InitialPension:     20000,
				PensionType:        "FERS",
				AssumedCPIWRate:    0.015,
				RetirementAge:      60,
				ProjectionStartAge: 61,
				YearsToProject:     3,
			},
			expectedResult: models.PensionCOLAResult{
				ProjectedPension: []float64{20000.00, 20300.00, 20604.50},
			},
			expectedNotesContain: "FERS: COLA applied (age >= 62).",
		},
		{
			name: "No Pension Type, 2% COLA, 3 years",
			input: models.PensionCOLAInput{
				InitialPension:     10000,
				PensionType:        "OTHER",
				AssumedCPIWRate:    0.02,
				RetirementAge:      60,
				ProjectionStartAge: 60,
				YearsToProject:     3,
			},
			expectedResult: models.PensionCOLAResult{
				ProjectedPension: []float64{10000.00, 10000.00, 10000.00},
			},
			expectedNotesContain: "Unknown pension type or no COLA policy",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ProjectPensionWithCOLA(tc.input)
			if !reflect.DeepEqual(result.ProjectedPension, tc.expectedResult.ProjectedPension) {
				t.Errorf("TestProjectPensionWithCOLA() %s: projection mismatch (-want +got):\n%s", tc.name, cmp.Diff(tc.expectedResult.ProjectedPension, result.ProjectedPension))
			}
			if tc.expectedNotesContain != "" && !containsSubstring(result.Notes, tc.expectedNotesContain) {
				t.Errorf("TestProjectPensionWithCOLA() %s: notes = %q, want to contain %q", tc.name, result.Notes, tc.expectedNotesContain)
			}
		})
	}
}

// Helper function for checking substrings in notes, can be moved to a testutils package later
func containsSubstring(s, substr string) bool {
    return strings.Contains(s, substr)
}
