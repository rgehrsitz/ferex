package models

// CSRSCalculationInput holds all user-provided data for CSRS/CSRS Offset annuity calculation.
type CSRSCalculationInput struct {
	High3Salary              float64 // Highest average basic pay over 3 consecutive years
	YearsOfService           float64 // Total years (including partial years) of creditable service
	UnusedSickLeaveMonths    int     // Unused sick leave in months (converted to service credit)
	IsPartTime               bool    // True if any part-time service
	PartTimeProrationFactor  float64 // Proration factor (1.0 = full time)
	EmployeeContributions    float64 // For tax-free portion calculation (optional)
	SurvivorBenefitElection  string  // e.g., "max", "partial", "none" (optional)
	IsCSRSOffset             bool    // True if CSRS Offset
	YearsOfOffsetService     float64 // Only for CSRS Offset
	SSAt62WithOffset         float64 // Only for CSRS Offset: SS benefit at 62 with Offset earnings
	SSAt62WithoutOffset      float64 // Only for CSRS Offset: SS benefit at 62 without Offset earnings
	AgeAtRetirement          int     // Age at retirement (for reductions)
}

// CSRSCalculationResult holds the calculated pension and related details.
type CSRSCalculationResult struct {
	AnnualPension            float64 // Gross annual pension
	MonthlyPension           float64 // Gross monthly pension
	EarlyRetirementReduction float64 // Total reduction for early retirement (if any)
	SickLeaveServiceCredit   float64 // Years added from unused sick leave
	ProrationApplied         bool    // True if part-time proration applied
	ProratedPension          float64 // Pension after proration (if applicable)
	SurvivorBenefitReduction float64 // Reduction for survivor benefit election
	OffsetReduction          float64 // Reduction for CSRS Offset (if applicable)
	Notes                    string  // Any warnings, special conditions, or info
}
