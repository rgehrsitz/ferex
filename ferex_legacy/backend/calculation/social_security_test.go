package calculation // Changed from 'tests'

import (
	// "ferex/backend/calculation" // Removed, test in same package
	"ferex/backend/models"      // Assuming 'ferex' is module name
	"ferex/backend/testutils" // Assuming 'ferex' is module name
	"testing"
	"math"
)

// Helper for comparing floats with a small tolerance
func floatEqualsSocialSecurity(a, b, tolerance float64) bool {
	return math.Abs(a-b) < tolerance
}

func TestGetFRA(t *testing.T) {
	cases := []struct {
		name      string
		birthYear int
		expectFRA models.FRADetail
	}{
		{"<=1937", 1937, models.FRADetail{Years: 65, Months: 0}},
		{"1938", 1938, models.FRADetail{Years: 65, Months: 2}},
		{"1939", 1939, models.FRADetail{Years: 65, Months: 4}},
		{"1940", 1940, models.FRADetail{Years: 65, Months: 6}},
		{"1941", 1941, models.FRADetail{Years: 65, Months: 8}},
		{"1942", 1942, models.FRADetail{Years: 65, Months: 10}},
		{"1943", 1943, models.FRADetail{Years: 66, Months: 0}},
		{"1950 (mid 1943-1954)", 1950, models.FRADetail{Years: 66, Months: 0}},
		{"1954", 1954, models.FRADetail{Years: 66, Months: 0}},
		{"1955", 1955, models.FRADetail{Years: 66, Months: 2}},
		{"1956", 1956, models.FRADetail{Years: 66, Months: 4}},
		{"1957", 1957, models.FRADetail{Years: 66, Months: 6}},
		{"1958", 1958, models.FRADetail{Years: 66, Months: 8}},
		{"1959", 1959, models.FRADetail{Years: 66, Months: 10}},
		{">=1960", 1960, models.FRADetail{Years: 67, Months: 0}},
		{"1970 (future)", 1970, models.FRADetail{Years: 67, Months: 0}},
		{"Edge case far past", 2050, models.FRADetail{Years: 67, Months: 0}},
		{"Edge case far past (default behavior test)", 1900, models.FRADetail{Years: 65, Months: 0}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := getFRA(tc.birthYear)
			if got.Years != tc.expectFRA.Years || got.Months != tc.expectFRA.Months {
				t.Errorf("getFRA(%d) = %+v, want %+v", tc.birthYear, got, tc.expectFRA)
			}
		})
	}
}

func TestSocialSecurityCalculation(t *testing.T) {
	// For 2025 PIA Bend Points: 1115, 6721. Percentages: 0.9, 0.32, 0.15
	// Example AIME: 5000
	// PIA = 0.9 * 1115 + 0.32 * (5000 - 1115) = 1003.5 + 0.32 * 3885 = 1003.5 + 1243.2 = 2246.7
	// Example AIME: 8000
	// PIA = 0.9 * 1115 + 0.32 * (6721 - 1115) + 0.15 * (8000 - 6721)
	//     = 1003.5 + 0.32 * 5606 + 0.15 * 1279
	//     = 1003.5 + 1793.92 + 191.85 = 2989.27

	cases := []struct {
		name          string
		input         models.SocialSecurityCalculationInput
		expect62      float64
		expectFRA     float64
		expect70      float64
		expectClaim   float64
		notesContains string
	}{
		{
			name: "User-provided SSA statement, claim at FRA",
			input: models.SocialSecurityCalculationInput{
				BirthYear:               1960, // FRA = 67
				UserProvidedEstimate62:  1200.50,
				UserProvidedEstimateFRA: 1700.80,
				UserProvidedEstimate70:  2100.20,
				ClaimAge:                67,
			},
			expect62:      1200.50,
			expectFRA:     1700.80,
			expect70:      2100.20,
			expectClaim:   1700.80, // Claiming at FRA
			notesContains: "SSA statement",
		},
		{
			name: "User-provided SSA statement, claim at 62",
			input: models.SocialSecurityCalculationInput{
				BirthYear:               1960, // FRA = 67
				UserProvidedEstimate62:  1200.00,
				UserProvidedEstimateFRA: 1700.00,
				UserProvidedEstimate70:  2100.00,
				ClaimAge:                62,
			},
			expect62:      1200.00,
			expectFRA:     1700.00,
			expect70:      2100.00,
			expectClaim:   1200.00, 
			notesContains: "SSA statement",
		},
		{
			name: "User-provided SSA statement, claim at 70",
			input: models.SocialSecurityCalculationInput{
				BirthYear:               1960, // FRA = 67
				UserProvidedEstimate62:  1200.00,
				UserProvidedEstimateFRA: 1700.00,
				UserProvidedEstimate70:  2100.00,
				ClaimAge:                70,
			},
			expect62:      1200.00,
			expectFRA:     1700.00,
			expect70:      2100.00,
			expectClaim:   2100.00, 
			notesContains: "SSA statement",
		},
		{
			name: "User-provided SSA statement, claim age 65 (interpolation)",
			input: models.SocialSecurityCalculationInput{
				BirthYear:               1960, // FRA = 67
				UserProvidedEstimate62:  1200.00,
				UserProvidedEstimateFRA: 1700.00, // Base for interpolation
				UserProvidedEstimate70:  2100.00,
				ClaimAge:                65, // Interpolate between 62 and 67 (FRA)
			},
			expect62:      1200.00,
			expectFRA:     1700.00,
			expect70:      2100.00,
			// PIA from FRA is 1700. Claiming factor for age 65, FRA 67 (24 months early):
			// Reduction: 24 * (5/9)/100 = 0.13333...; Benefit = 1700 * (1 - 0.13333) = 1700 * 0.86666... = 1473.33
			expectClaim:   1473.30, // PIA * factor(65,67) = 1700.00 * (1.0 - (24 * (5.0/9.0/100.0))) = 1700 * 0.86666... approx 1473.33. Rounded to 1473.3
			notesContains: "SSA statement",
		},
		{
			name: "Estimate from salary, BirthYear 1960 (FRA 67), AIME from 60k salary for 35y",
			input: models.SocialSecurityCalculationInput{
				BirthYear:             1960, // FRA = 67
				EstimatedAnnualSalary: 60000,
				YearsWorked:           35,
				ClaimAge:              67,
			},
			// AIME = 60000 / 12 = 5000
			// PIA for AIME 5000 = 0.9*1115 + 0.32*(5000-1115) = 1003.5 + 0.32*3885 = 1003.5 + 1243.2 = 2246.7
			expectFRA:     2246.7, // This is the PIA
			// At 62 (FRA 67, 5 years = 60 months early): 36m @ 5/9%, 24m @ 5/12%. Reduction: 36*(5/9/100) + 24*(5/12/100) = 0.20 + 0.10 = 0.30. Factor = 0.7
			expect62:      2246.7 * 0.7, // 1572.69 -> 1572.6
			// At 70 (FRA 67, 3 years = 36 months late): Factor = 1 + 3*0.08 = 1.24
			expect70:      2246.7 * 1.24, // 2785.908 -> 2785.9
			expectClaim:   2246.7, // Claiming at FRA
			notesContains: "average salary",
		},
		{
			name: "Estimate from salary, BirthYear 1954 (FRA 66), AIME from 80k salary for 35y, claim at 62",
			input: models.SocialSecurityCalculationInput{
				BirthYear:             1954, // FRA = 66
				EstimatedAnnualSalary: 80000,
				YearsWorked:           35,
				ClaimAge:              62,
			},
			// AIME = 80000 / 12 = 6666.666...
			// PIA for AIME 6666.67 = 0.9*1115 + 0.32*(6666.67-1115) = 1003.5 + 0.32*5551.67 = 1003.5 + 1776.53 = 2780.03
			expectFRA:     2780.0, 
			// At 62 (FRA 66, 4 years = 48 months early): 36m @ 5/9%, 12m @ 5/12%. Reduction: 36*(5/9/100) + 12*(5/12/100) = 0.20 + 0.05 = 0.25. Factor = 0.75
			expect62:      2780.0 * 0.75, // 2085.0
			// At 70 (FRA 66, 4 years = 48 months late): Factor = 1 + 4*0.08 = 1.32
			expect70:      2780.0 * 1.32, // 3669.6
			expectClaim:   2085.0, // Claiming at 62
			notesContains: "average salary",
		},
		{
			name: "Missing all data, should be zero",
			input: models.SocialSecurityCalculationInput{
				BirthYear:  1970, // FRA 67
				ClaimAge:   62,
			},
			expect62:      0,
			expectFRA:     0,
			expect70:      0,
			expectClaim:   0,
			notesContains: "AIME is zero",
		},
	}

	tolerance := 0.11 // PIA is rounded to dime, so results can differ slightly

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := CalculateSocialSecurity(tc.input) // Direct call

			if !floatEqualsSocialSecurity(got.EstimatedAt62, tc.expect62, tolerance) {
				t.Errorf("%s: EstimatedAt62 got %.2f, want %.2f", tc.name, got.EstimatedAt62, tc.expect62)
			}
			if !floatEqualsSocialSecurity(got.EstimatedAtFRA, tc.expectFRA, tolerance) {
				t.Errorf("%s: EstimatedAtFRA got %.2f, want %.2f", tc.name, got.EstimatedAtFRA, tc.expectFRA)
			}
			if !floatEqualsSocialSecurity(got.EstimatedAt70, tc.expect70, tolerance) {
				t.Errorf("%s: EstimatedAt70 got %.2f, want %.2f", tc.name, got.EstimatedAt70, tc.expect70)
			}
			if !floatEqualsSocialSecurity(got.ClaimingAmount, tc.expectClaim, tolerance) {
				t.Errorf("%s: ClaimingAmount got %.2f, want %.2f", tc.name, got.ClaimingAmount, tc.expectClaim)
			}

			if tc.notesContains != "" && !testutils.Contains(got.Notes, tc.notesContains) {
				t.Errorf("%s: notes ('%s') missing expected string: %q", tc.name, got.Notes, tc.notesContains)
			}
		})
	}
}
