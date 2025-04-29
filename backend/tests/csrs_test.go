package tests

import (
	"ferex/backend/calculation"
	"ferex/backend/models"
	"ferex/backend/tests/testutils"
	"testing"
)

func TestCSRSCalculation(t *testing.T) {
	cases := []struct {
		name              string
		input             models.CSRSCalculationInput
		expect            float64
		offsetReduction   float64
		survivorReduction float64
		prorationApplied  bool
		notesContains     string
	}{
		{
			name: "Standard CSRS (30 years, no proration, no survivor)",
			input: models.CSRSCalculationInput{
				High3Salary:             100000,
				YearsOfService:          30,
				UnusedSickLeaveMonths:   0,
				IsPartTime:              false,
				PartTimeProrationFactor: 1.0,
				SurvivorBenefitElection: "none",
				IsCSRSOffset:            false,
			},
			expect:            100000 * (5*0.015 + 5*0.0175 + 20*0.02), // $56,250
			offsetReduction:   0,
			survivorReduction: 0,
			prorationApplied:  false,
			notesContains:     "",
		},
		{
			name: "CSRS 80% High-3 maximum",
			input: models.CSRSCalculationInput{
				High3Salary:             120000,
				YearsOfService:          45,
				UnusedSickLeaveMonths:   0,
				IsPartTime:              false,
				PartTimeProrationFactor: 1.0,
				SurvivorBenefitElection: "none",
				IsCSRSOffset:            false,
			},
			expect:            120000 * 0.8, // $96,000
			offsetReduction:   0,
			survivorReduction: 0,
			prorationApplied:  false,
			notesContains:     "80% High-3 maximum",
		},
		{
			name: "CSRS with part-time proration (0.75)",
			input: models.CSRSCalculationInput{
				High3Salary:             80000,
				YearsOfService:          20,
				UnusedSickLeaveMonths:   0,
				IsPartTime:              true,
				PartTimeProrationFactor: 0.75,
				SurvivorBenefitElection: "none",
				IsCSRSOffset:            false,
			},
			expect:            80000 * (5*0.015 + 5*0.0175 + 10*0.02) * 0.75, // $28,500 * 0.75 = $21,375
			offsetReduction:   0,
			survivorReduction: 0,
			prorationApplied:  true,
			notesContains:     "Part-time proration factor applied",
		},
		{
			name: "CSRS with sick leave credit (6 months)",
			input: models.CSRSCalculationInput{
				High3Salary:             90000,
				YearsOfService:          25,
				UnusedSickLeaveMonths:   6,
				IsPartTime:              false,
				PartTimeProrationFactor: 1.0,
				SurvivorBenefitElection: "none",
				IsCSRSOffset:            false,
			},
			expect:            90000 * (5*0.015 + 5*0.0175 + 15.5*0.02), // 25.5 years
			offsetReduction:   0,
			survivorReduction: 0,
			prorationApplied:  false,
			notesContains:     "",
		},
		{
			name: "CSRS with max survivor benefit",
			input: models.CSRSCalculationInput{
				High3Salary:             100000,
				YearsOfService:          30,
				UnusedSickLeaveMonths:   0,
				IsPartTime:              false,
				PartTimeProrationFactor: 1.0,
				SurvivorBenefitElection: "max",
				IsCSRSOffset:            false,
			},
			expect:            100000 * (5*0.015 + 5*0.0175 + 20*0.02) * 0.9, // 10% reduction
			offsetReduction:   0,
			survivorReduction: 100000 * (5*0.015 + 5*0.0175 + 20*0.02) * 0.1,
			prorationApplied:  false,
			notesContains:     "Maximum survivor benefit reduction",
		},
		{
			name: "CSRS Offset (option 1 lower)",
			input: models.CSRSCalculationInput{
				High3Salary:             110000,
				YearsOfService:          30,
				UnusedSickLeaveMonths:   0,
				IsPartTime:              false,
				PartTimeProrationFactor: 1.0,
				SurvivorBenefitElection: "none",
				IsCSRSOffset:            true,
				YearsOfOffsetService:    10,
				SSAt62WithOffset:        12000,
				SSAt62WithoutOffset:     10000,
				AgeAtRetirement:         62,
			},
			expect:            110000*(5*0.015+5*0.0175+20*0.02) - 2000, // Offset reduction = 2000 (option 1 lower)
			offsetReduction:   2000,
			survivorReduction: 0,
			prorationApplied:  false,
			notesContains:     "CSRS Offset reduction applied",
		},
		{
			name: "CSRS Offset (option 2 lower)",
			input: models.CSRSCalculationInput{
				High3Salary:             110000,
				YearsOfService:          30,
				UnusedSickLeaveMonths:   0,
				IsPartTime:              false,
				PartTimeProrationFactor: 1.0,
				SurvivorBenefitElection: "none",
				IsCSRSOffset:            true,
				YearsOfOffsetService:    5,
				SSAt62WithOffset:        12000,
				SSAt62WithoutOffset:     10000,
				AgeAtRetirement:         62,
			},
			expect:            110000*(5*0.015+5*0.0175+20*0.02) - 1500, // Offset reduction = 1500 (option 2 lower)
			offsetReduction:   1500,
			survivorReduction: 0,
			prorationApplied:  false,
			notesContains:     "CSRS Offset reduction applied",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := calculation.CalculateCSRS(tc.input)
			if testutils.Abs(got.AnnualPension-tc.expect) > 0.01 {
				t.Errorf("%s: got %.2f, want %.2f", tc.name, got.AnnualPension, tc.expect)
			}
			if testutils.Abs(got.OffsetReduction-tc.offsetReduction) > 0.01 {
				t.Errorf("%s: offset reduction got %.2f, want %.2f", tc.name, got.OffsetReduction, tc.offsetReduction)
			}
			if testutils.Abs(got.SurvivorBenefitReduction-tc.survivorReduction) > 0.01 {
				t.Errorf("%s: survivor reduction got %.2f, want %.2f", tc.name, got.SurvivorBenefitReduction, tc.survivorReduction)
			}
			if got.ProrationApplied != tc.prorationApplied {
				t.Errorf("%s: proration applied got %v, want %v", tc.name, got.ProrationApplied, tc.prorationApplied)
			}
			if tc.notesContains != "" && !testutils.Contains(got.Notes, tc.notesContains) {
				t.Errorf("%s: notes missing expected string: %q", tc.name, tc.notesContains)
			}
		})
	}
}
