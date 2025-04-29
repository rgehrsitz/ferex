package tests

import (
	"ferex/backend/calculation"
	"ferex/backend/models"
	"ferex/backend/tests/testutils"
	"math"
	"testing"
)

func TestTSPCalculation(t *testing.T) {
	cases := []struct {
		name          string
		input         models.TSPCalculationInput
		expectBalance float64
		expectAnnual  float64
		expectYears   int
		notesContains string
	}{
		{
			name: "Standard projection, fixed withdrawal, balance lasts",
			input: models.TSPCalculationInput{
				CurrentBalance:             200000,
				AnnualEmployeeContribution: 10000,
				AnnualAgencyMatch:          5000,
				YearsUntilRetirement:       10,
				ExpectedAnnualReturnRate:   0.06,
				WithdrawalMethod:           "fixed",
				WithdrawalAmountOrPercent:  40000,
				ProjectedWithdrawalYears:   20,
			},
			expectBalance: 200000*math.Pow(1.06, 10) + sumContributions(10000+5000, 0.06, 10),
			expectAnnual:  40000,
			expectYears:   20,
			notesContains: "",
		},
		{
			name: "Fixed withdrawal, balance exhausted early",
			input: models.TSPCalculationInput{
				CurrentBalance:             100000,
				AnnualEmployeeContribution: 0,
				AnnualAgencyMatch:          0,
				YearsUntilRetirement:       5,
				ExpectedAnnualReturnRate:   0.03,
				WithdrawalMethod:           "fixed",
				WithdrawalAmountOrPercent:  60000,
				ProjectedWithdrawalYears:   20,
			},
			expectBalance: 100000 * math.Pow(1.03, 5),
			expectAnnual:  60000,
			expectYears:   3, // Simulation result: balance lasts 3 years
			notesContains: "Balance exhausted",
		},
		{
			name: "Percent withdrawal",
			input: models.TSPCalculationInput{
				CurrentBalance:             500000,
				AnnualEmployeeContribution: 0,
				AnnualAgencyMatch:          0,
				YearsUntilRetirement:       0,
				ExpectedAnnualReturnRate:   0.05,
				WithdrawalMethod:           "percent",
				WithdrawalAmountOrPercent:  6,
				ProjectedWithdrawalYears:   30,
			},
			expectBalance: 500000,
			expectAnnual:  500000 * 0.06,
			expectYears:   30,
			notesContains: "",
		},
		{
			name: "RMD withdrawal",
			input: models.TSPCalculationInput{
				CurrentBalance:             400000,
				AnnualEmployeeContribution: 0,
				AnnualAgencyMatch:          0,
				YearsUntilRetirement:       0,
				ExpectedAnnualReturnRate:   0.05,
				WithdrawalMethod:           "RMD",
				WithdrawalAmountOrPercent:  0,
				ProjectedWithdrawalYears:   20,
			},
			expectBalance: 400000,
			expectAnnual:  400000 / 20,
			expectYears:   20,
			notesContains: "",
		},
		{
			name: "Default 4% rule",
			input: models.TSPCalculationInput{
				CurrentBalance:             300000,
				AnnualEmployeeContribution: 0,
				AnnualAgencyMatch:          0,
				YearsUntilRetirement:       0,
				ExpectedAnnualReturnRate:   0.05,
				WithdrawalMethod:           "",
				WithdrawalAmountOrPercent:  0,
				ProjectedWithdrawalYears:   0,
			},
			expectBalance: 300000,
			expectAnnual:  300000 * 0.04,
			expectYears:   25,
			notesContains: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := calculation.CalculateTSP(tc.input)
			if testutils.Abs(got.ProjectedBalanceAtRetirement-tc.expectBalance) > 1.0 {
				t.Errorf("%s: projected balance got %.2f, want %.2f", tc.name, got.ProjectedBalanceAtRetirement, tc.expectBalance)
			}
			if testutils.Abs(got.AnnualWithdrawalIncome-tc.expectAnnual) > 1.0 {
				t.Errorf("%s: annual withdrawal got %.2f, want %.2f", tc.name, got.AnnualWithdrawalIncome, tc.expectAnnual)
			}
			if got.YearsBalanceLasts != tc.expectYears {
				t.Errorf("%s: years balance lasts got %d, want %d", tc.name, got.YearsBalanceLasts, tc.expectYears)
			}
			if tc.notesContains != "" && !testutils.Contains(got.Notes, tc.notesContains) {
				t.Errorf("%s: notes missing expected string: %q", tc.name, tc.notesContains)
			}
		})
	}
}

// sumContributions calculates the future value of a series of annual contributions with growth
func sumContributions(annual float64, rate float64, years int) float64 {
	fv := 0.0
	for i := 0; i < years; i++ {
		fv = (fv + annual) * (1 + rate)
	}
	return fv
}
