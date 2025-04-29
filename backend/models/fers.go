package models

// FERSCalculationInput holds all user-provided data needed for FERS annuity calculation.
type FERSCalculationInput struct {
	High3Salary               float64 // Highest average basic pay over 3 consecutive years
	YearsOfService            float64 // Total years (including partial years) of creditable service
	AgeAtRetirement           int     // Age at retirement
	UnusedSickLeaveMonths     int     // Unused sick leave in months (converted to service credit)
	IsPartTime                bool    // True if any part-time service
	PartTimeProrationFactor   float64 // Proration factor (1.0 = full time)
	IsAge62With20Years        bool    // True if age 62+ with 20+ years (for 1.1% multiplier)
	SurvivorBenefitElection   string  // e.g., "max", "partial", "none" (optional)
	EmployeeContributions     float64 // For tax-free portion calculation (optional)
}

// FERSCalculationResult holds the calculated pension and related details.
type FERSCalculationResult struct {
	AnnualPension             float64 // Gross annual pension
	MonthlyPension            float64 // Gross monthly pension
	EarlyRetirementReduction  float64 // Total reduction for early retirement (if any)
	SickLeaveServiceCredit    float64 // Years added from unused sick leave
	ProrationApplied          bool    // True if part-time proration applied
	ProratedPension           float64 // Pension after proration (if applicable)
	SurvivorBenefitReduction  float64 // Reduction for survivor benefit election
	Notes                     string  // Any warnings, special conditions, or info
}
