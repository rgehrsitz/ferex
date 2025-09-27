package models

// SurvivorBenefitCalculationInput holds data for projecting survivor annuity/income.
type SurvivorBenefitCalculationInput struct {
	PensionType         string  // "FERS", "CSRS", "CSRSOffset"
	InitialAnnuity      float64 // Retiree's annual annuity before survivor reduction
	SurvivorElection    string  // "max", "partial", "none"
	SpouseAge           int     // Age of spouse/beneficiary
	RetireeAgeAtDeath   int     // Age of retiree at death
	COLARate            float64 // Annual COLA applied to survivor annuity
	YearsToProject      int     // Years to project survivor income
	IncludeSSSurvivor   bool    // Whether to include SS survivor benefit
	SSSurvivorAmount    float64 // Annual SS survivor benefit (if included)
	IncludeTSP          bool    // Whether to include TSP continuation
	TSPBalanceAtDeath   float64 // TSP balance at death (if included)
	OtherSurvivorIncome float64 // Other survivor income (optional)
}

// SurvivorBenefitCalculationResult holds projected survivor income details.
type SurvivorBenefitCalculationResult struct {
	InitialSurvivorAnnuity float64   // Survivor annuity in year 1
	ProjectedAnnuities     []float64 // Survivor annuity for each projected year
	TotalSurvivorIncome    float64   // Cumulative survivor income over projection
	Notes                  string    // Election details, warnings, etc.
}
