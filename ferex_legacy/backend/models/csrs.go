package models

// CSRSCalculationInput holds all user-provided data for CSRS calculations.
type CSRSCalculationInput struct {
	High3Salary             float64 // Highest average basic pay over 3 consecutive years
	UnusedSickLeaveHours    float64 // Unused sick leave in HOURS (converted to service credit)
	SurvivorBenefitType     string  // e.g., "None", "Full" (55%), "Partial" (custom, often 25% or 50% of max)
	IsPartTime              bool    // True if any part-time service
	PartTimeProrationFactor float64 // Proration factor (1.0 = full time)
	EmployeeContributions   float64 // For tax-free portion calculation (optional)
	IsCSRSOffset            bool    // True if CSRS Offset
	YearsOfOffsetService    float64 // Only for CSRS Offset
	SSAt62WithOffset        float64 // Only for CSRS Offset: SS benefit at 62 with Offset earnings
	SSAt62WithoutOffset     float64 // Only for CSRS Offset: SS benefit at 62 without Offset earnings

	// Date fields for derivation
	DateOfBirth            string `json:"DateOfBirth"`            // YYYY-MM-DD
	ServiceComputationDate string `json:"ServiceComputationDate"` // YYYY-MM-DD
	PlannedRetirementDate  string `json:"PlannedRetirementDate"`  // YYYY-MM-DD
}

// CSRSCalculationResult holds the calculated pension and related details.
// NOTE: All monetary calculations are now performed monthly-first for better user relatability.
type CSRSCalculationResult struct {
	IRSSimplifiedMethodExclusion       float64 // Monthly tax-free portion of annuity per IRS Simplified Method
	MonthlyGrossAnnuity                float64 `json:"monthlyGrossAnnuity"`                // Primary calculation - monthly annuity before deductions/offset
	SickLeaveServiceCredit             float64 `json:"sickLeaveServiceCredit"`             // Years added from sick leave
	TotalServiceYears                  float64 `json:"totalServiceYears"`                  // Total years used in calc (service + sick leave)
	IsProrated                         bool    `json:"isProrated"`                         // Flag if part-time proration applied
	MonthlyGrossAnnuityBeforeProration float64 `json:"monthlyGrossAnnuityBeforeProration"` // Monthly gross annuity before part-time proration
	MonthlyOffsetReduction             float64 `json:"monthlyOffsetReduction"`             // Monthly offset reduction for CSRS Offset
	MonthlyNetAnnuity                  float64 `json:"monthlyNetAnnuity"`                  // Monthly annuity after offset (if any)
	MonthlySurvivorBenefitReduction    float64 `json:"monthlySurvivorBenefitReduction"`    // Monthly amount reduced for survivor benefit
	MonthlyFinalAnnuity                float64 `json:"monthlyFinalAnnuity"`                // Primary result - monthly annuity after all deductions
	RetirementType                     string  `json:"retirementType"`                     // Computed retirement type (Immediate, Deferred, etc.)
	ProrationFactor                    float64 `json:"prorationFactor"`                    // Actual proration factor used
	Notes                              string  `json:"notes"`

	// Legacy annual fields - computed from monthly values for backward compatibility
	GrossAnnuity                float64 `json:"grossAnnuity"`                // Derived: MonthlyGrossAnnuity * 12
	GrossAnnuityBeforeProration float64 `json:"grossAnnuityBeforeProration"` // Derived: MonthlyGrossAnnuityBeforeProration * 12
	OffsetReduction             float64 `json:"offsetReduction"`             // Derived: MonthlyOffsetReduction * 12
	NetAnnuity                  float64 `json:"netAnnuity"`                  // Derived: MonthlyNetAnnuity * 12
	SurvivorBenefitReduction    float64 `json:"survivorBenefitReduction"`    // Derived: MonthlySurvivorBenefitReduction * 12
	FinalAnnuity                float64 `json:"finalAnnuity"`                // Derived: MonthlyFinalAnnuity * 12
}
