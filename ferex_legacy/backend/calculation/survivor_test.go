package calculation // Changed from 'tests'

import (
	// "ferex/backend/calculation" // Removed, test in same package
	"ferex/backend/models"      // Assuming 'ferex' is module name
	"ferex/backend/testutils" // Assuming 'ferex' is module name
	"testing"
	"math"
)

// Using testutils.Abs, so no specific floatEquals needed here unless different tolerance.

func TestSurvivorBenefitCalculation(t *testing.T) {
	cases := []struct {
		name                       string
		input                      models.SurvivorBenefitCalculationInput
		expectInitialSurvivorAnnuity float64
		// expectTotalProjectedIncome float64 // Legacy tests didn't check exact total, rather relative
		notesContains              string
		// Used for more complex total check if needed later
		expectProjectedYear1IfApplicable float64 
	}{
		{
			name: "FERS max survivor, 3 years projection, 2% COLA, no extras",
			input: models.SurvivorBenefitCalculationInput{
				PensionType:      "FERS",
				InitialAnnuity:   40000,
				SurvivorElection: "max", // 10% retiree reduction, 50% of reduced to survivor
				COLARate:         0.02,
				YearsToProject:   3,
			},
			// Retiree reduced: 40000 * (1 - 0.10) = 36000
			// Survivor initial: 36000 * 0.50 = 18000
			expectInitialSurvivorAnnuity: 18000.00,
			notesContains:              "FERS max",
			// Year 1: 18000
			// Year 2: 18000 * 1.02 = 18360
			// Year 3: 18360 * 1.02 = 18727.20
			// Total: 18000 + 18360 + 18727.20 = 55087.20
			expectProjectedYear1IfApplicable: 18000.00, // For checking first year of projection
		},
		{
			name: "CSRS partial survivor (25% of unreduced), 2 years, 3% COLA, with SS & TSP",
			input: models.SurvivorBenefitCalculationInput{
				PensionType:       "CSRS",
				InitialAnnuity:    60000,
				SurvivorElection:  "partial",
				COLARate:          0.03,
				YearsToProject:    2,
				IncludeSSSurvivor: true,
				SSSurvivorAmount:  12000, // annual
				IncludeTSP:        true,
				TSPBalanceAtDeath: 10000, // total, spread over YearsToProject
			},
			// CSRS Partial (Survivor gets 25% of unreduced):
			// Target Survivor Benefit: 0.25 * 60000 = 15000.
			// Base for Cost Formula: 15000 / 0.55 = 27272.7272...
			// Cost: (0.025 * 3600) + (0.10 * (27272.7272 - 3600)) = 90 + (0.10 * 23672.7272) = 90 + 2367.2727 = 2457.2727
			// Retiree's Reduced Annuity: 60000 - 2457.2727 = 57542.7273
			// Initial Survivor Annuity (from pension): 0.25 * 60000 = 15000.
			expectInitialSurvivorAnnuity: 15000.00,
			notesContains:              "CSRS partial: Survivor receives 25% of retiree's unreduced annuity",
			// Year 1 income: InitialSurvivor (15000) + SSSurvivor (12000) + TSP (10000/2 = 5000) = 32000
			expectProjectedYear1IfApplicable: 32000.00,
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
			expectInitialSurvivorAnnuity: 0,
			notesContains:              "No survivor benefit elected",
			expectProjectedYear1IfApplicable: 0.00,
		},
		{
			name: "CSRS max survivor (accurate calc), 1 year",
			input: models.SurvivorBenefitCalculationInput{
				PensionType:      "CSRS",
				InitialAnnuity:   70000,
				SurvivorElection: "max",
				COLARate:         0.01,
				YearsToProject:   1,
			},
			// Accurate CSRS Max:
			// Retiree Cost: (0.025 * 3600) + (0.10 * (70000 - 3600)) = 90 + (0.10 * 66400) = 90 + 6640 = 6730.
			// Retiree Reduced Annuity: 70000 - 6730 = 63270.
			// Survivor Initial Annuity: 70000 * 0.55 = 38500.
			expectInitialSurvivorAnnuity: 38500.00,
			notesContains:              "CSRS max: Survivor receives 55% of retiree's unreduced annuity", // Note updated
			expectProjectedYear1IfApplicable: 38500.00, // Year 1, no COLA, no extras
		},
	}

	tolerance := 0.01 // For monetary comparisons

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := CalculateSurvivorBenefit(tc.input)
			if math.Abs(got.InitialSurvivorAnnuity-tc.expectInitialSurvivorAnnuity) > tolerance {
				t.Errorf("%s: InitialSurvivorAnnuity got %.2f, want %.2f", tc.name, got.InitialSurvivorAnnuity, tc.expectInitialSurvivorAnnuity)
			}

			if tc.notesContains != "" && !testutils.Contains(got.Notes, tc.notesContains) {
				t.Errorf("%s: notes ('%s') missing expected string: %q", tc.name, got.Notes, tc.notesContains)
			}

			// Check first year of projected annuities if applicable
			if tc.input.YearsToProject > 0 && len(got.ProjectedAnnuities) > 0 {
				if math.Abs(got.ProjectedAnnuities[0]-tc.expectProjectedYear1IfApplicable) > tolerance {
					t.Errorf("%s: ProjectedAnnuities[0] got %.2f, want %.2f", tc.name, got.ProjectedAnnuities[0], tc.expectProjectedYear1IfApplicable)
				}
			} else if tc.input.YearsToProject > 0 && len(got.ProjectedAnnuities) == 0 {
				t.Errorf("%s: Expected projected annuities, but got none", tc.name)
			} else if tc.input.YearsToProject == 0 && len(got.ProjectedAnnuities) > 0 {
				t.Errorf("%s: Did not expect projected annuities, but got some", tc.name)
			}

			// Optionally, a more robust check for TotalSurvivorIncome could be added if exact values were determined.
			// For now, the legacy test's approach of not strictly checking total is maintained.
		})
	}
}
