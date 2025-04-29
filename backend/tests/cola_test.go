package tests

import (
	"ferex/backend/calculation"
	"ferex/backend/models"
	"ferex/backend/tests/testutils"
	"testing"
)

func TestCOLACalculation(t *testing.T) {
	cases := []struct {
		name          string
		input         models.COLACalculationInput
		expectFinal   float64
		notesContains string
	}{
		{
			name: "FERS COLA cap, CPI 4%",
			input: models.COLACalculationInput{
				InitialAmount: 10000,
				COLARate:      0.04, // 4% CPI-W
				Years:         3,
				COLAPolicy:    "FERS",
			},
			expectFinal:   10000 * 1.03 * 1.03 * 1.03, // FERS cap: 3% per year
			notesContains: "FERS COLA cap logic applied",
		},
		{
			name: "CSRS full COLA, 2%",
			input: models.COLACalculationInput{
				InitialAmount: 20000,
				COLARate:      0.02,
				Years:         2,
				COLAPolicy:    "CSRS",
			},
			expectFinal:   20000 * 1.02 * 1.02,
			notesContains: "",
		},
		{
			name: "No COLA",
			input: models.COLACalculationInput{
				InitialAmount: 15000,
				COLARate:      0.03,
				Years:         5,
				COLAPolicy:    "None",
			},
			expectFinal:   15000,
			notesContains: "",
		},
		{
			name: "Generic COLA 1.5%",
			input: models.COLACalculationInput{
				InitialAmount: 12000,
				COLARate:      0.015,
				Years:         4,
				COLAPolicy:    "Generic",
			},
			expectFinal:   12000 * 1.015 * 1.015 * 1.015 * 1.015,
			notesContains: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := calculation.CalculateCOLA(tc.input)
			if testutils.Abs(got.FinalAmount-tc.expectFinal) > 0.5 {
				t.Errorf("%s: final amount got %.2f, want %.2f", tc.name, got.FinalAmount, tc.expectFinal)
			}
			if tc.notesContains != "" && !testutils.Contains(got.Notes, tc.notesContains) {
				t.Errorf("%s: notes missing expected: %q", tc.name, tc.notesContains)
			}
		})
	}
}
