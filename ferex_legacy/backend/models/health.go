package models

// HealthPremiumCalculationInput holds data for projecting FEHB/Medicare premiums.
type HealthPremiumCalculationInput struct {
	FEHBPremium             float64 // Annual FEHB premium (user or default)
	MedicarePremium         float64 // Annual Medicare Part B premium (user or default)
	IncludeFEHB             bool    // Whether to include FEHB in projection
	IncludeMedicare         bool    // Whether to include Medicare in projection
	COLARate                float64 // Annual COLA/inflation rate for premiums
	YearsToProject          int     // Years to project premiums
	OtherHealthPremium      float64 // Other annual health premiums (optional)
}

// HealthPremiumCalculationResult holds projected health premium details.
type HealthPremiumCalculationResult struct {
	ProjectedPremiums   []float64 // Total annual premiums for each year
	TotalPremiums       float64   // Cumulative premiums over projection
	Notes               string    // Plan notes, warnings, etc.
}
