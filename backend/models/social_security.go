package models

// SocialSecurityCalculationInput holds user data for SS benefit projection.
type SocialSecurityCalculationInput struct {
	BirthYear               int       // User's year of birth (for FRA)
	CurrentAge              int       // User's age (for projections)
	EarningsHistory         []float64 // Optional: indexed earnings for each year (up to 35 years)
	EstimatedAnnualSalary   float64   // Optional: if no full history, use average salary
	YearsWorked             int       // Optional: for simple estimate
	UserProvidedEstimate62  float64   // Optional: SSA statement estimate at age 62
	UserProvidedEstimateFRA float64   // Optional: SSA statement estimate at FRA
	UserProvidedEstimate70  float64   // Optional: SSA statement estimate at age 70
	ClaimAge                int       // Desired claiming age (62-70)
}

// SocialSecurityCalculationResult holds the projected SS benefits.
type SocialSecurityCalculationResult struct {
	EstimatedAt62  float64 // Monthly benefit at age 62
	EstimatedAtFRA float64 // Monthly benefit at FRA
	EstimatedAt70  float64 // Monthly benefit at age 70
	ClaimingAge    int     // Age used for benefit calculation
	ClaimingAmount float64 // Monthly benefit at chosen claiming age
	Notes          string  // Any warnings, method notes, etc.
}
