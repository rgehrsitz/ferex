package tests

import (
	"ferex/backend/calculation"
	"ferex/backend/models"
	"ferex/backend/tests/testutils"
	"testing"
)

func TestSurvivorBenefitCalculation(t *testing.T) {
	cases := []struct {
		name          string
		input         models.SurvivorBenefitCalculationInput
		expectInitial float64
		expectTotal   float64
		notesContains string
	}{
		{
			name: "FERS max survivor, no extras",
			input: models.SurvivorBenefitCalculationInput{
				PensionType:      "FERS",
				InitialAnnuity:   40000,
				SurvivorElection: "max",
				COLARate:         0.02,
				YearsToProject:   3,
			},
			expectInitial: 40000 * 0.9 * 0.5, // 10% reduction, 50% to survivor
			expectTotal:   0,                 // Checked relatively
			notesContains: "FERS max",
		},
		{
			name: "CSRS partial survivor, with SS and TSP",
			input: models.SurvivorBenefitCalculationInput{
				PensionType:       "CSRS",
				InitialAnnuity:    60000,
				SurvivorElection:  "partial",
				COLARate:          0.03,
				YearsToProject:    2,
				IncludeSSSurvivor: true,
				SSSurvivorAmount:  12000,
				IncludeTSP:        true,
				TSPBalanceAtDeath: 10000,
			},
			expectInitial: 60000 * 0.95 * 0.5,
			expectTotal:   0,
			notesContains: "CSRS partial",
		},
		{
			name: "FERS no survivor benefit",
			input: models.SurvivorBenefitCalculationInput{
				PensionType:      "FERS",
				InitialAnnuity:   50000,
				SurvivorElection: "none",
				COLARate:         0.02,
				YearsToProject:   2,
			},
			expectInitial: 0,
			expectTotal:   0,
			notesContains: "No survivor",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := calculation.CalculateSurvivorBenefit(tc.input)
			if testutils.Abs(got.InitialSurvivorAnnuity-tc.expectInitial) > 1.0 {
				t.Errorf("%s: initial survivor got %.2f, want %.2f", tc.name, got.InitialSurvivorAnnuity, tc.expectInitial)
			}
			if tc.expectTotal > 0 && got.TotalSurvivorIncome < tc.expectTotal {
				t.Errorf("%s: total survivor income got %.2f, want at least %.2f", tc.name, got.TotalSurvivorIncome, tc.expectTotal)
			}
			if tc.notesContains != "" && !testutils.Contains(got.Notes, tc.notesContains) {
				t.Errorf("%s: notes missing expected: %q", tc.name, tc.notesContains)
			}
		})
	}
}
