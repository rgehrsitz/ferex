package models

// RetirementCalculationInput aggregates all user data for a full retirement projection.
type RetirementCalculationInput struct {
	// Pension Inputs
	FERSInput          FERSCalculationInput
	CSRSInput          CSRSCalculationInput
	SRSInput           SRSCalculationInput
	TSPInput           TSPCalculationInput
	TaxInput           TaxCalculationInput
	SocialSecurityInput SocialSecurityCalculationInput
	COLAInput          COLACalculationInput
	SurvivorInput      SurvivorBenefitCalculationInput
	HealthInput        HealthPremiumCalculationInput

	// Monte Carlo simulation (optional)
	MonteCarloInput    MonteCarloInput
}

// RetirementCalculationResult aggregates all outputs for the full retirement projection.
type RetirementCalculationResult struct {
	// Individual Results
	FERSResult          FERSCalculationResult
	CSRSResult          CSRSCalculationResult
	SRSResult           SRSCalculationResult
	TSPResult           TSPCalculationResult
	TaxResult           TaxCalculationResult
	SocialSecurityResult SocialSecurityCalculationResult
	COLAResult          COLACalculationResult
	SurvivorResult      SurvivorBenefitCalculationResult
	HealthResult        HealthPremiumCalculationResult

	// Monte Carlo simulation (optional)
	MonteCarloResult    MonteCarloResult

	// Aggregated/summary fields
	NetAfterTaxIncome   float64 // Final net income after all deductions
	EffectiveTaxRate    float64 // Overall effective tax rate
	TotalRetirementIncome float64 // Sum of all projected retirement income streams
	Notes               string  // Any warnings or summary notes
}
