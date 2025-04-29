package tests

import (
	"ferex/backend/calculation"
	"ferex/backend/models"
	"ferex/backend/tests/testutils"
	"testing"
)

func TestFERSCalculation(t *testing.T) {
	cases := []struct {
		name   string
		input  models.FERSCalculationInput
		expect float64
	}{
		{
			name: "Standard FERS (1.0%)",
			input: models.FERSCalculationInput{
				High3Salary:             100000,
				YearsOfService:          30,
				AgeAtRetirement:         60,
				IsAge62With20Years:      false,
				SurvivorBenefitElection: "none",
				IsPartTime:              false,
				PartTimeProrationFactor: 1.0,
			},
			expect: 100000 * 30 * 0.01 * 0.9, // $30,000 * 0.9 = $27,000 (10% reduction for age 60)
		},
		{
			name: "FERS Age 62+ with 20+ years (1.1%)",
			input: models.FERSCalculationInput{
				High3Salary:             100000,
				YearsOfService:          25,
				AgeAtRetirement:         65,
				IsAge62With20Years:      true,
				SurvivorBenefitElection: "none",
				IsPartTime:              false,
				PartTimeProrationFactor: 1.0,
			},
			expect: 100000 * 25 * 0.011, // $27,500
		},
		{
			name: "Early Retirement Reduction (MRA+10, 5 years under 62)",
			input: models.FERSCalculationInput{
				High3Salary:             90000,
				YearsOfService:          20,
				AgeAtRetirement:         57,
				IsAge62With20Years:      false,
				SurvivorBenefitElection: "none",
				IsPartTime:              false,
				PartTimeProrationFactor: 1.0,
			},
			expect: (90000 * 20 * 0.01) * (1 - 0.05*5), // $18,000 - 25% = $13,500
		},
		{
			name: "Part-Time Proration (0.8)",
			input: models.FERSCalculationInput{
				High3Salary:             80000,
				YearsOfService:          25,
				AgeAtRetirement:         60,
				IsPartTime:              true,
				PartTimeProrationFactor: 0.8,
				SurvivorBenefitElection: "none",
			},
			expect: (80000 * 25 * 0.01 * 0.9) * 0.8, // $20,000 * 0.9 * 0.8 = $14,400 (10% reduction for age 60)
		},
		{
			name: "Sick Leave Credit (6 months)",
			input: models.FERSCalculationInput{
				High3Salary:             75000,
				YearsOfService:          20,
				UnusedSickLeaveMonths:   6,
				AgeAtRetirement:         60,
				SurvivorBenefitElection: "none",
				IsPartTime:              false,
				PartTimeProrationFactor: 1.0,
			},
			expect: 75000 * (20.5) * 0.01 * 0.9, // 10% reduction for age 60
		},
		{
			name: "Max Survivor Benefit (10% reduction)",
			input: models.FERSCalculationInput{
				High3Salary:             100000,
				YearsOfService:          30,
				AgeAtRetirement:         62,
				IsAge62With20Years:      true,
				SurvivorBenefitElection: "max",
			},
			expect: (100000 * 30 * 0.011) * 0.9, // $33,000 * 0.9 = $29,700
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := calculation.CalculateFERSPension(tc.input)
			if testutils.Abs(got.AnnualPension-tc.expect) > 0.01 {
				t.Errorf("%s: got %.2f, want %.2f", tc.name, got.AnnualPension, tc.expect)
			}
		})
	}
}
