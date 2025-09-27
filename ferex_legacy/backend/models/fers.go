package models

// FERSCalculationInput holds all user-provided data needed for FERS annuity calculation.
type FERSCalculationInput struct {
	High3Salary float64 // Highest average basic pay over 3 consecutive years
	// YearsOfService is derived
	// AgeAtRetirement is derived
	UnusedSickLeaveHours  int             // Unused sick leave in hours (converted to service credit)
	EmployeeContributions float64         // Optional: For tax-free portion calculation
	ServicePeriods        []ServicePeriod // Detailed service periods for proration calculation
	// IsAge62With20Years is derived
	SurvivorBenefitElection string  // e.g., "None", "Partial", "Full"
	ExpectedSSBenefitAt62   float64 `json:"expectedSSBenefitAt62"` // Estimated Social Security benefit at age 62 (for supplement calc)
	// Date fields for derivation
	DateOfBirth            string // YYYY-MM-DD
	ServiceComputationDate string // YYYY-MM-DD
	PlannedRetirementDate  string // YYYY-MM-DD
	SwitchedToFERSDate     string // YYYY-MM-DD, Optional: If employee switched from CSRS
}

// AuditStep represents a single step in the calculation process
type AuditStep struct {
	StepNumber  int                    `json:"stepNumber"`
	StepName    string                 `json:"stepName"`
	Description string                 `json:"description"`
	Formula     string                 `json:"formula"`
	Inputs      map[string]interface{} `json:"inputs"`
	Calculation string                 `json:"calculation"`
	Result      float64                `json:"result"`
	Notes       string                 `json:"notes,omitempty"`
}

// CalculationAuditTrail provides detailed step-by-step calculation breakdown
type CalculationAuditTrail struct {
	CalculationType string      `json:"calculationType"` // "FERS", "CSRS", "TSP", etc.
	InputSummary    string      `json:"inputSummary"`
	Steps           []AuditStep `json:"steps"`
	FinalResult     float64     `json:"finalResult"`
	Warnings        []string    `json:"warnings,omitempty"`
	OPMReferences   []string    `json:"ompReferences,omitempty"` // Links to official OPM documentation
}

// FERSCalculationResult holds the calculated pension and related details.
// NOTE: All monetary calculations are now performed monthly-first for better user relatability.
type FERSCalculationResult struct {
	IRSSimplifiedMethodExclusion    float64 // Monthly tax-free portion of annuity per IRS Simplified Method
	MonthlyPension                  float64 `json:"monthlyPension"`                  // Primary calculation - monthly pension amount
	AnnualPension                   float64 `json:"annualPension"`                   // Derived: MonthlyPension * 12
	MonthlyEarlyRetirementReduction float64 `json:"monthlyEarlyRetirementReduction"` // Monthly reduction for early retirement (if any)
	SickLeaveServiceCredit          float64 // Years added from unused sick leave
	ProrationApplied                bool    // True if part-time proration applied
	MonthlyProratedPension          float64 `json:"monthlyProratedPension"`          // Monthly pension after proration (if applicable)
	MonthlySurvivorBenefitReduction float64 `json:"monthlySurvivorBenefitReduction"` // Monthly reduction for survivor benefit election
	TotalServiceYears               float64 `json:"totalServiceYears"`               // Added
	MonthlyBasicAnnuity             float64 `json:"monthlyBasicAnnuity"`             // Monthly gross annuity before reductions
	IsEligibleForSupplement         bool    `json:"isEligibleForSupplement"`         // Added
	FersSupplement                  float64 `json:"fersSupplement"`                  // Already monthly - FERS supplement amount
	RetirementType                  string  `json:"retirementType"`                  // Added: Computed retirement type (Immediate, MRA+10, etc.)
	ProrationFactor                 float64 `json:"prorationFactor"`                 // Added: Actual proration factor used
	Notes                           string  `json:"notes"`

	// NEW: Detailed audit trail
	AuditTrail *CalculationAuditTrail `json:"auditTrail,omitempty"`

	// Legacy annual fields - computed from monthly values for backward compatibility
	BasicAnnuity             float64 `json:"basicAnnuity"` // Derived: MonthlyBasicAnnuity * 12
	EarlyRetirementReduction float64 // Derived: MonthlyEarlyRetirementReduction * 12
	ProratedPension          float64 // Derived: MonthlyProratedPension * 12
	SurvivorBenefitReduction float64 // Derived: MonthlySurvivorBenefitReduction * 12
}
