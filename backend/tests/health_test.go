package tests

import (
	"ferex/backend/calculation"
	"ferex/backend/models"
	"ferex/backend/tests/testutils"
	"testing"
)

func TestHealthPremiumCalculation(t *testing.T) {
	cases := []struct {
		name          string
		input         models.HealthPremiumCalculationInput
		expectFirst   float64
		expectLast    float64
		expectTotal   float64
		notesContains string
	}{
		{
			name: "FEHB only, no COLA",
			input: models.HealthPremiumCalculationInput{
				FEHBPremium:     5000,
				IncludeFEHB:     true,
				IncludeMedicare: false,
				COLARate:        0.0,
				YearsToProject:  3,
			},
			expectFirst:   5000,
			expectLast:    5000,
			expectTotal:   15000,
			notesContains: "FEHB",
		},
		{
			name: "Medicare only, 2% COLA",
			input: models.HealthPremiumCalculationInput{
				MedicarePremium: 2000,
				IncludeFEHB:     false,
				IncludeMedicare: true,
				COLARate:        0.02,
				YearsToProject:  2,
			},
			expectFirst:   2000,
			expectLast:    2000 * 1.02,
			expectTotal:   2000 + 2000*1.02,
			notesContains: "Medicare",
		},
		{
			name: "FEHB + Medicare + Other, 3% COLA",
			input: models.HealthPremiumCalculationInput{
				FEHBPremium:        6000,
				MedicarePremium:    1800,
				OtherHealthPremium: 500,
				IncludeFEHB:        true,
				IncludeMedicare:    true,
				COLARate:           0.03,
				YearsToProject:     2,
			},
			expectFirst:   6000 + 1800 + 500,
			expectLast:    (6000 * 1.03) + (1800 * 1.03) + (500 * 1.03),
			expectTotal:   (6000 + 1800 + 500) + ((6000 * 1.03) + (1800 * 1.03) + (500 * 1.03)),
			notesContains: "COLA",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := calculation.CalculateHealthPremiums(tc.input)
			if testutils.Abs(got.ProjectedPremiums[0]-tc.expectFirst) > 1.0 {
				t.Errorf("%s: first year premium got %.2f, want %.2f", tc.name, got.ProjectedPremiums[0], tc.expectFirst)
			}
			if testutils.Abs(got.ProjectedPremiums[len(got.ProjectedPremiums)-1]-tc.expectLast) > 1.0 {
				t.Errorf("%s: last year premium got %.2f, want %.2f", tc.name, got.ProjectedPremiums[len(got.ProjectedPremiums)-1], tc.expectLast)
			}
			if testutils.Abs(got.TotalPremiums-tc.expectTotal) > 2.0 {
				t.Errorf("%s: total premiums got %.2f, want %.2f", tc.name, got.TotalPremiums, tc.expectTotal)
			}
			if tc.notesContains != "" && !testutils.Contains(got.Notes, tc.notesContains) {
				t.Errorf("%s: notes missing expected: %q", tc.name, tc.notesContains)
			}
		})
	}
}
