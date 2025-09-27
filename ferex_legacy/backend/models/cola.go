package models

// COLACalculationInput holds the data for projecting COLA-adjusted values.
type COLACalculationInput struct {
	InitialAmount float64 // Starting value (e.g. pension, SS, withdrawal)
	COLARate      float64 // Annual COLA rate (e.g. 0.02 for 2%)
	Years         int     // Number of years to project
	COLAPolicy    string  // "FERS", "CSRS", "SocialSecurity", "None", or custom (optional)
	StartYear     int     // First year of projection (optional)
}

// COLACalculationResult holds the COLA-adjusted projections.
type COLACalculationResult struct {
	ProjectedAmounts   []float64 // Annual values after COLA for each year
	FinalAmount        float64   // Value after all years
	CumulativeCOLA     float64   // Total percent increase over period
	Notes              string    // Policy notes, warnings, etc.
}

// PensionCOLAInput defines the input for COLA calculations on a pension.
type PensionCOLAInput struct {
	InitialPension     float64 // The starting annual pension amount
	PensionType        string  // "CSRS" or "FERS"
	AssumedCPIWRate    float64 // The assumed annual CPI-W rate for COLA calculations (e.g., 0.025 for 2.5%)
	RetirementAge      int     // The age at which the InitialPension starts
	ProjectionStartAge int     // The age at the beginning of the projection period (YearsToProject[0])
	YearsToProject     int     // The number of years to project the pension for
}

// PensionCOLAResult defines the result of a COLA calculation.
type PensionCOLAResult struct {
	ProjectedPension []float64 // Slice of projected annual pension amounts, including COLAs
	Notes            string    // Explanatory notes about the calculation (e.g., when FERS COLA starts)
}

// SocialSecurityCOLAInput defines the input for COLA calculations on Social Security benefits.
type SocialSecurityCOLAInput struct {
	InitialAnnualBenefit float64 // The starting annual Social Security benefit
	AssumedCPIWRate      float64 // The assumed annual CPI-W rate for COLA calculations (e.g., 0.025 for 2.5%)
	BenefitStartAge      int     // The age at which the InitialAnnualBenefit starts (claim age)
	ProjectionStartAge   int     // The age at the beginning of the projection period (YearsToProject[0])
	YearsToProject       int     // The number of years to project the benefit for
}

// SocialSecurityCOLAResult defines the result of a COLA calculation for Social Security.
type SocialSecurityCOLAResult struct {
	ProjectedBenefit []float64 // Slice of projected annual Social Security benefits, including COLAs
	Notes            string    // Explanatory notes about the calculation
}
