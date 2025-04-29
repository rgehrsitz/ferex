package tests

import (
	"ferex/backend/calculation"
	"ferex/backend/models"
	"ferex/backend/tests/testutils"
	"testing"
)

func TestSRSCalculation(t *testing.T) {
	cases := []struct {
		name              string
		input             models.SRSCalculationInput
		expect            float64
		earningsReduction float64
		isEligible        bool
	}{
		{
			name: "Eligible, no earnings test",
			input: models.SRSCalculationInput{
				EstimatedSocialSecurityAt62: 20000,
				YearsOfFERSService:          30,
				RetirementAge:               60,
				MRA:                         57,
				IsImmediateUnreducedAnnuity: true,
				ProjectedEarnedIncome:       0,
				RetirementYear:              2025,
			},
			expect:            20000 / 40 * 30, // $15,000
			earningsReduction: 0,
			isEligible:        true,
		},
		{
			name: "Eligible, with earnings test reduction",
			input: models.SRSCalculationInput{
				EstimatedSocialSecurityAt62: 24000,
				YearsOfFERSService:          25,
				RetirementAge:               60,
				MRA:                         57,
				IsImmediateUnreducedAnnuity: true,
				ProjectedEarnedIncome:       25400, // $2,000 over limit
				RetirementYear:              2025,
			},
			expect:            (24000 / 40 * 25) - 1000, // $15,000 - $1,000
			earningsReduction: 1000,
			isEligible:        true,
		},
		{
			name: "Not eligible: MRA+10 (not immediate, unreduced)",
			input: models.SRSCalculationInput{
				EstimatedSocialSecurityAt62: 18000,
				YearsOfFERSService:          15,
				RetirementAge:               57,
				MRA:                         57,
				IsImmediateUnreducedAnnuity: false,
				ProjectedEarnedIncome:       0,
				RetirementYear:              2025,
			},
			expect:            0,
			earningsReduction: 0,
			isEligible:        false,
		},
		{
			name: "Not eligible: retire at 62",
			input: models.SRSCalculationInput{
				EstimatedSocialSecurityAt62: 22000,
				YearsOfFERSService:          20,
				RetirementAge:               62,
				MRA:                         57,
				IsImmediateUnreducedAnnuity: true,
				ProjectedEarnedIncome:       0,
				RetirementYear:              2025,
			},
			expect:            0,
			earningsReduction: 0,
			isEligible:        false,
		},
		{
			name: "Earnings test wipes out SRS",
			input: models.SRSCalculationInput{
				EstimatedSocialSecurityAt62: 16000,
				YearsOfFERSService:          20,
				RetirementAge:               60,
				MRA:                         57,
				IsImmediateUnreducedAnnuity: true,
				ProjectedEarnedIncome:       100000, // Well over limit
				RetirementYear:              2025,
			},
			expect:            0,
			earningsReduction: 8000, // SRS would be 8,000, but all wiped out
			isEligible:        true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := calculation.CalculateSRS(tc.input)
			if testutils.Abs(got.AnnualSRSAmount-tc.expect) > 0.01 {
				t.Errorf("%s: got %.2f, want %.2f", tc.name, got.AnnualSRSAmount, tc.expect)
			}
			if testutils.Abs(got.EarningsTestReduction-tc.earningsReduction) > 0.01 {
				t.Errorf("%s: earnings test reduction got %.2f, want %.2f", tc.name, got.EarningsTestReduction, tc.earningsReduction)
			}
			if got.IsEligible != tc.isEligible {
				t.Errorf("%s: eligibility got %v, want %v", tc.name, got.IsEligible, tc.isEligible)
			}
		})
	}
}
