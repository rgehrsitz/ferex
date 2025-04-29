package calculation

import (
	"ferex/backend/models"
	"math"
)

// CalculateTSP projects the TSP balance at retirement and models withdrawals.
func CalculateTSP(input models.TSPCalculationInput) models.TSPCalculationResult {
	var notes string

	// Project balance at retirement with annual contributions and compound growth
	balance := input.CurrentBalance
	for i := 0; i < input.YearsUntilRetirement; i++ {
		balance += input.AnnualEmployeeContribution + input.AnnualAgencyMatch
		balance *= 1 + input.ExpectedAnnualReturnRate
	}
	projectedBalance := balance

	// Withdrawal modeling
	withdrawalYears := input.ProjectedWithdrawalYears
	if withdrawalYears == 0 {
		withdrawalYears = 30 // default to 30 years for "lifetime" if not specified
	}

	annualWithdrawal := 0.0
	monthlyWithdrawal := 0.0
	yearsBalanceLasts := withdrawalYears

	switch input.WithdrawalMethod {
	case "fixed":
		annualWithdrawal = input.WithdrawalAmountOrPercent
		if annualWithdrawal > 0 {
			// Simulate withdrawals until balance runs out or withdrawal years reached
			b := projectedBalance
			years := 0
			for years < withdrawalYears && b > 0 {
				b = b*(1+input.ExpectedAnnualReturnRate) - annualWithdrawal
				years++
			}
			yearsBalanceLasts = years
			if b <= 0 {
				notes += "Balance exhausted before end of withdrawal period.\n"
			}
		}
	case "percent":
		annualWithdrawal = projectedBalance * (input.WithdrawalAmountOrPercent / 100.0)
		// No exhaustion simulation; just a percent of initial balance
		monthlyWithdrawal = annualWithdrawal / 12.0
	case "RMD":
		// Placeholder: RMD logic could be more complex
		annualWithdrawal = projectedBalance / float64(withdrawalYears)
		monthlyWithdrawal = annualWithdrawal / 12.0
		yearsBalanceLasts = withdrawalYears
	default:
		// Default to fixed, 4% rule
		annualWithdrawal = projectedBalance * 0.04
		monthlyWithdrawal = annualWithdrawal / 12.0
		yearsBalanceLasts = int(math.Min(float64(withdrawalYears), projectedBalance/annualWithdrawal))
	}

	if monthlyWithdrawal == 0 && annualWithdrawal > 0 {
		monthlyWithdrawal = annualWithdrawal / 12.0
	}

	return models.TSPCalculationResult{
		ProjectedBalanceAtRetirement: projectedBalance,
		AnnualWithdrawalIncome:       annualWithdrawal,
		MonthlyWithdrawalIncome:      monthlyWithdrawal,
		YearsBalanceLasts:            yearsBalanceLasts,
		Notes:                        notes,
	}
}
