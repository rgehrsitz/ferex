package models

// TaxCalculationInput holds all user-provided data for tax estimation on retirement income.
type TaxCalculationInput struct {
	FilingStatus         string  // Filing status: "single", "married", etc.
	TaxYear              int     // Tax year for bracket selection
	Age                  int     // Age of taxpayer
	GrossPension         float64 // Total FERS/CSRS pension
	TaxablePension       float64 // Taxable portion of pension (after exclusions)
	TSPWithdrawal        float64 // Taxable TSP withdrawals (Traditional)
	TSPRothWithdrawal    float64 // Roth TSP withdrawals (not federally taxable)
	SocialSecurity       float64 // Social Security benefit (taxable portion computed in logic)
	OtherTaxableIncome   float64 // Other taxable income (interest, dividends, etc.)
	StateOfResidence     string  // State for state tax calculation (optional)
	StateTaxableIncome   float64 // State taxable income (optional)
	Deductions           float64 // Standard or itemized deduction
	TaxCredits           float64 // Tax credits (optional)
}

// TaxCalculationResult holds the output of the tax estimation.
type TaxCalculationResult struct {
	FederalTaxOwed       float64 // Total federal tax liability
	StateTaxOwed         float64 // Total state tax liability (optional)
	NetAfterTaxIncome    float64 // Net income after federal and state taxes
	EffectiveTaxRate     float64 // Overall effective tax rate
	Notes                string  // Any warnings, special conditions, or info
}
