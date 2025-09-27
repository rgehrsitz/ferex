package models

// RetirementCalculationInput aggregates all user data for a full retirement projection.
type RetirementCalculationInput struct {
	// System Selection
	CalculationSystem string // 'FERS' | 'CSRS' | '' - user's chosen retirement system

	// Pension Inputs
	FERSInput           FERSCalculationInput
	CSRSInput           CSRSCalculationInput
	SRSInput            SRSCalculationInput
	TSPInput            TSPCalculationInput
	TaxInput            TaxCalculationInput
	SocialSecurityInput SocialSecurityCalculationInput
	COLAInput           COLACalculationInput
	SurvivorInput       SurvivorBenefitCalculationInput
	HealthInput         HealthPremiumCalculationInput

	// Monte Carlo simulation (optional)
	MonteCarloInput MonteCarloInput
}

// MonthlyFinancialProjection holds consolidated financial data for a single month of projection.
type MonthlyFinancialProjection struct {
	Year                             int     `json:"year"`
	Month                            int     `json:"month"` // 1-12 representing Jan-Dec
	AgeYears                         int     `json:"ageYears"`
	AgeMonths                        int     `json:"ageMonths"`
	PensionForMonth                  float64 `json:"pensionForMonth"`
	SocialSecurityForMonth           float64 `json:"socialSecurityForMonth"`
	TSPWithdrawalTraditionalForMonth float64 `json:"tspWithdrawalTraditionalForMonth"`
	TSPWithdrawalRothForMonth        float64 `json:"tspWithdrawalRothForMonth"`
	HealthPremiumsForMonth           float64 `json:"healthPremiumsForMonth,omitempty"`
	// SalaryIncomeForMonth          float64 `json:"salaryIncomeForMonth,omitempty"` // If pre-retirement months are included
	// OtherIncomeForMonth           float64 `json:"otherIncomeForMonth,omitempty"`  // For future flexibility
	TotalPreTaxIncomeForMonth   float64 `json:"totalPreTaxIncomeForMonth"`
	AllocatedFederalTaxForMonth float64 `json:"allocatedFederalTaxForMonth"`
	AllocatedStateTaxForMonth   float64 `json:"allocatedStateTaxForMonth,omitempty"`
	NetCashFlowForMonth         float64 `json:"netCashFlowForMonth"`
	Notes                       string  `json:"notes,omitempty"` // e.g., "Retirement Month", "SS Claim Month", "COLA Applied"
}

// RetirementCalculationResult aggregates all outputs for the full retirement projection.
type RetirementCalculationResult struct {
	// Individual Results
	FERSResult           FERSCalculationResult
	CSRSResult           CSRSCalculationResult
	SRSResult            SRSCalculationResult
	TSPResult            TSPCalculationResult
	TaxResult            TaxCalculationResult
	SocialSecurityResult SocialSecurityCalculationResult
	COLAResult           COLACalculationResult
	SurvivorResult       SurvivorBenefitCalculationResult
	HealthResult         HealthPremiumCalculationResult

	// Monte Carlo simulation (optional)
	MonteCarloResult MonteCarloResult

	// Aggregated/summary fields
	NetAfterTaxIncome     float64 // Final net income after all deductions
	EffectiveTaxRate      float64 // Overall effective tax rate
	TotalRetirementIncome float64 // Sum of all projected retirement income streams
	Notes                 string  // Any warnings or summary notes

	// Detailed time-series projection
	DetailedMonthlyProjections []MonthlyFinancialProjection `json:"detailedMonthlyProjections,omitempty"`
}
