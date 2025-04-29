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
