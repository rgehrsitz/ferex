package models

// TSPCalculationInput holds all user-provided data for TSP projection and withdrawal modeling.
type TSPCalculationInput struct {
	CurrentBalance           float64 // Current TSP account balance
	AnnualEmployeeContribution float64 // Annual employee contribution (dollars)
	AnnualAgencyMatch        float64 // Annual agency match (dollars)
	YearsUntilRetirement     int     // Number of years until retirement
	ExpectedAnnualReturnRate float64 // Expected annual return rate (as decimal, e.g., 0.06 for 6%)
	WithdrawalStartAge       int     // Age at which withdrawals begin
	WithdrawalMethod         string  // "fixed", "percent", "RMD", etc.
	WithdrawalAmountOrPercent float64 // Amount (if fixed) or percent (if percent method)
	ProjectedWithdrawalYears int     // Number of years for withdrawals (0 for "lifetime")
	RothBalance              float64 // Optional: Roth TSP balance
	TraditionalBalance       float64 // Optional: Traditional TSP balance
}

// TSPCalculationResult holds the output of the TSP projection and withdrawal modeling.
type TSPCalculationResult struct {
	ProjectedBalanceAtRetirement float64 // Projected TSP balance at retirement
	AnnualWithdrawalIncome       float64 // Projected annual withdrawal income
	MonthlyWithdrawalIncome      float64 // Projected monthly withdrawal income
	YearsBalanceLasts            int     // Number of years balance lasts (if applicable)
	Notes                        string  // Any warnings, special conditions, or info
}
