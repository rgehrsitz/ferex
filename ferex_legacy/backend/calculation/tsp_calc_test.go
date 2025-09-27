package calculation

import (
	"testing"
	"ferex/backend/models"
	"math"
)

func TestCalculateTSP_RMDs(t *testing.T) {
	testCases := []struct {
		name string
		birthYear int
		testAge            int     // Specific age for which to check RMD and withdrawal
		initialBalance     float64 // Balance at the END of testAge-1 (start of testAge before RMD)
		expectedRmdAmount  float64 // Approximate, based on factor and balance
		expectRmdApplicable bool
	}{
		// --- Test Case Group 1: Born 1950 (RMD starts at 72) ---
		{
			name: "Born 1950, Age 71, RMD not applicable", birthYear: 1950, testAge: 71,
			initialBalance: 100000, expectedRmdAmount: 0, expectRmdApplicable: false,
		},
		{
			name: "Born 1950, Age 72, RMD applicable", birthYear: 1950, testAge: 72,
			initialBalance: 100000, expectedRmdAmount: 100000 / 27.4, expectRmdApplicable: true,
		},
		{
			name: "Born 1950, Age 73, RMD applicable (post-start)", birthYear: 1950, testAge: 73,
			initialBalance: 95000, expectedRmdAmount: 95000 / 26.5, expectRmdApplicable: true, // Assuming some growth/withdrawal happened
		},

		// --- Test Case Group 2: Born 1955 (RMD starts at 73) ---
		{
			name: "Born 1955, Age 72, RMD not applicable", birthYear: 1955, testAge: 72,
			initialBalance: 100000, expectedRmdAmount: 0, expectRmdApplicable: false,
		},
		{
			name: "Born 1955, Age 73, RMD applicable", birthYear: 1955, testAge: 73,
			initialBalance: 100000, expectedRmdAmount: 100000 / 26.5, expectRmdApplicable: true,
		},

		// --- Test Case Group 3: Born 1960 (RMD starts at 75) ---
		{
			name: "Born 1960, Age 74, RMD not applicable", birthYear: 1960, testAge: 74,
			initialBalance: 100000, expectedRmdAmount: 0, expectRmdApplicable: false,
		},
		{
			name: "Born 1960, Age 75, RMD applicable", birthYear: 1960, testAge: 75,
			initialBalance: 100000, expectedRmdAmount: 100000 / 24.6, expectRmdApplicable: true,
		},
		{
			name: "Born 1960, Age 72, RMD not applicable (well before start)", birthYear: 1960, testAge: 72,
			initialBalance: 100000, expectedRmdAmount: 0, expectRmdApplicable: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tspInput := models.TSPCalculationInput{
				BirthYear:                 tc.birthYear,
				CurrentAgeYears:           tc.testAge,       // Test RMD AT this age
				RetirementAgeYears:        tc.testAge,       // Retired AT this age (no growth/contrib phase before withdrawal calc)
				CurrentTraditionalBalance: tc.initialBalance,
				CurrentRothBalance:        0,
				UserReturnAssumptions:     models.TSPReturnAssumptions{Overall: float64Ptr(0.0), UseOverallForPre: true, UseOverallForPost: true}, // No growth during the single year of withdrawal projection
				CurrentAllocation:         models.TSPFundAllocationPercentages{G: float64Ptr(100.0)}, // Default, growth is 0 anyway
				ContributionFundAllocation:models.TSPFundAllocationPercentages{G: float64Ptr(100.0)}, // Default
				ContributeUntilRetirement: false,             // No contributions during the withdrawal year projection
				EmployeeContributionPercentage: 0,
				EmployeeContributionAmount:    0,
				IsContributionPercentage:      true,
				BaseSalaryForContributions:    0,             // Not relevant if no contributions
				TraditionalContributionAllocationPct: 100,
				RothContributionAllocationPct: 0,
				CatchUpContributionsEligible: false,
				WithdrawalStrategy:        "IRSMinimumRequiredDistribution",
				WithdrawalStartDate:       "SpecificAge",
				WithdrawalStartAge:        tc.testAge,
				YearsToProjectWithdrawals: 1,                // Only project for the RMD year
				ExpenseRatio:              0.0,              // No expenses for this precise RMD test
				WithdrawalOrder:           "ProRata",        // Explicitly set
			}

			result := CalculateTSP(tspInput)

			if len(result.WithdrawalSchedule) > 0 {
				detail := result.WithdrawalSchedule[0]

				if detail.Age != tc.testAge {
					t.Errorf("Withdrawal detail age %d does not match testAge %d. Notes: %s", detail.Age, tc.testAge, result.Notes)
				}

				rmdCalculatedForTestYear := detail.RmdAmount
				actualWithdrawalForTestYear := detail.TotalWithdrawn

				if tc.expectRmdApplicable {
					if rmdCalculatedForTestYear <= 0 && tc.initialBalance > 0 {
						t.Errorf("Expected RMD to be applicable and > 0 (initialBal: %.2f), got %.2f. Notes: %s", tc.initialBalance, rmdCalculatedForTestYear, result.Notes)
					}
					if math.Abs(rmdCalculatedForTestYear-tc.expectedRmdAmount) > 0.01 { // Allow for float precision
						t.Errorf("Expected RMD amount %.2f, got %.2f. Notes: %s", tc.expectedRmdAmount, rmdCalculatedForTestYear, result.Notes)
					}
					// If RMD is applicable, total withdrawn should be at least the RMD.
					// For IRSMinimumRequiredDistribution strategy, it should be exactly the RMD unless balance is insufficient.
					expectedWithdrawal := math.Min(tc.expectedRmdAmount, tc.initialBalance) // Cannot withdraw more than available
					if math.Abs(actualWithdrawalForTestYear-expectedWithdrawal) > 0.01 {
						t.Errorf("Expected total withdrawn %.2f (RMD or balance), got %.2f. Notes: %s", expectedWithdrawal, actualWithdrawalForTestYear, result.Notes)
					}
				} else { // RMD not applicable
					if rmdCalculatedForTestYear != 0 {
						t.Errorf("Expected RMD to be not applicable (amount 0), got %.2f. Notes: %s", rmdCalculatedForTestYear, result.Notes)
					}
					// If RMD is not applicable, and strategy is *only* RMD, then withdrawal should be 0.
					if actualWithdrawalForTestYear != 0 && tspInput.WithdrawalStrategy == "IRSMinimumRequiredDistribution" {
					    t.Errorf("Expected 0 withdrawal when RMD not applicable for IRSMinimumRequiredDistribution strategy, got %.2f. Notes: %s", actualWithdrawalForTestYear, result.Notes)
					}
				}
			} else if tc.expectRmdApplicable && tc.initialBalance > 0 {
				t.Errorf("Expected RMD to be applicable (initialBal: %.2f), but no withdrawal schedule generated. Notes: %s", tc.initialBalance, result.Notes)
			} else if !tc.expectRmdApplicable && len(result.WithdrawalSchedule) > 0 && result.WithdrawalSchedule[0].TotalWithdrawn > 0 {
			    // If RMD not applicable, and schedule IS generated, withdrawal should be 0 if strategy is only RMD.
			    if tspInput.WithdrawalStrategy == "IRSMinimumRequiredDistribution" && result.WithdrawalSchedule[0].TotalWithdrawn != 0 {
			        t.Errorf("RMD not applicable, but withdrawal of %.2f occurred with IRSMinimumRequiredDistribution strategy. Notes: %s", result.WithdrawalSchedule[0].TotalWithdrawn, result.Notes)
			    }
			}
		})
	}
}

// TODO: Add more test functions for other aspects:
// - TestCalculateTSP_ContributionLimitsAndGrowth()
// - TestCalculateTSP_WithdrawalOrders()
// - TestCalculateTSP_FixedVsPercentageWithdrawals()
// - TestCalculateTSP_DeferredWithdrawalsAndGrowth()
