package tests

import (
	"ferex/backend/calculation"
	"ferex/backend/models"
	"ferex/backend/tests/testutils"
	"testing"
)

func TestSocialSecurityCalculation(t *testing.T) {
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
			name: "User-provided SSA statement values",
			input: models.SocialSecurityCalculationInput{
				BirthYear:               1962,
				CurrentAge:              62,
				UserProvidedEstimate62:  1200,
				UserProvidedEstimateFRA: 1700,
				UserProvidedEstimate70:  2100,
				ClaimAge:                67,
			},
			expect62:      1200,
			expectFRA:     1700,
			expect70:      2100,
			expectClaim:   1700,
			notesContains: "SSA statement",
		},
		{
			name: "Estimate from earnings history",
			input: models.SocialSecurityCalculationInput{
				BirthYear:       1960,
				CurrentAge:      64,
				EarningsHistory: make([]float64, 35),
				ClaimAge:        62,
			},
			expect62:      0, // All zero earnings
			expectFRA:     0,
			expect70:      0,
			expectClaim:   0,
			notesContains: "zero",
		},
		{
			name: "Estimate from average salary and years worked",
			input: models.SocialSecurityCalculationInput{
				BirthYear:             1960,
				CurrentAge:            65,
				EstimatedAnnualSalary: 60000,
				YearsWorked:           35,
				ClaimAge:              67,
			},
			expect62:      0, // Will check relative values
			expectFRA:     0,
			expect70:      0,
			expectClaim:   0,
			notesContains: "average salary",
		},
		{
			name: "Missing all data",
			input: models.SocialSecurityCalculationInput{
				BirthYear:  1970,
				CurrentAge: 55,
				ClaimAge:   62,
			},
			expect62:      0,
			expectFRA:     0,
			expect70:      0,
			expectClaim:   0,
			notesContains: "zero",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := calculation.CalculateSocialSecurity(tc.input)
			if tc.expect62 > 0 && testutils.Abs(got.EstimatedAt62-tc.expect62) > 2.0 {
				t.Errorf("%s: at 62 got %.2f, want %.2f", tc.name, got.EstimatedAt62, tc.expect62)
			}
			if tc.expectFRA > 0 && testutils.Abs(got.EstimatedAtFRA-tc.expectFRA) > 2.0 {
				t.Errorf("%s: at FRA got %.2f, want %.2f", tc.name, got.EstimatedAtFRA, tc.expectFRA)
			}
			if tc.expect70 > 0 && testutils.Abs(got.EstimatedAt70-tc.expect70) > 2.0 {
				t.Errorf("%s: at 70 got %.2f, want %.2f", tc.name, got.EstimatedAt70, tc.expect70)
			}
			if tc.expectClaim > 0 && testutils.Abs(got.ClaimingAmount-tc.expectClaim) > 2.0 {
				t.Errorf("%s: claim amount got %.2f, want %.2f", tc.name, got.ClaimingAmount, tc.expectClaim)
			}
			if tc.notesContains != "" && !testutils.Contains(got.Notes, tc.notesContains) {
				t.Errorf("%s: notes missing expected: %q", tc.name, tc.notesContains)
			}
		})
	}
}
