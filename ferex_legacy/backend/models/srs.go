package models

// SRSCalculationInput holds all user-provided data needed for FERS Annuity Supplement (SRS) calculation.
type SRSCalculationInput struct {
	EstimatedSocialSecurityAt62 float64 // User's estimated SS benefit at age 62
	YearsOfFERSService          float64 // Only actual FERS service (not military buyback, not non-FERS)
	RetirementAge               int     // Age at retirement
	MRA                         int     // Minimum Retirement Age for eligibility
	IsImmediateUnreducedAnnuity bool    // True if eligible for immediate, unreduced annuity
	ProjectedEarnedIncome       float64 // Earned income before age 62 (for earnings test)
	RetirementYear              int     // Year of retirement (for earnings test threshold)
}

// SRSCalculationResult holds the calculated SRS benefit and related details.
type SRSCalculationResult struct {
	AnnualSRSAmount           float64 // Annual SRS benefit after reductions
	MonthlySRSAmount          float64 // Monthly SRS benefit after reductions
	EarningsTestReduction     float64 // Reduction due to earnings test
	IsEligible                bool    // True if eligible for SRS
	Notes                     string  // Any warnings, special conditions, or info
}
