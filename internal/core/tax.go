package core

import "strconv"

// TaxBracket describes a single marginal rate slice for progressive tax
// calculations. Amounts are in dollars and should be inflation-adjusted before
// use.
type TaxBracket struct {
	Threshold float64 `json:"threshold"`
	Rate      float64 `json:"rate"`
}

// TaxSettings stores the filing status, deductions, and bracket tables required
// to compute annual liabilities. Future iterations will allow state-specific
// details; the MVP focuses on federal taxation.
type TaxSettings struct {
	FilingStatus      TaxFilingStatus `json:"filing_status"`
	StandardDeduction float64         `json:"standard_deduction"`
	ItemizedDeduction float64         `json:"itemized_deduction"`
	UseItemized       bool            `json:"use_itemized"`
	FederalBrackets   []TaxBracket    `json:"federal_brackets"`
	TaxFreeBasis      float64         `json:"tax_free_basis"`
	StateName         string          `json:"state_name"`
	StateRate         float64         `json:"state_rate"`
}

// Validate ensures basic correctness of the tax configuration.
func (t TaxSettings) Validate() *ValidationError {
	var verr ValidationError
	if t.FilingStatus == "" {
		verr.Add("filing_status", "filing status is required")
	}
	if len(t.FederalBrackets) == 0 {
		verr.Add("federal_brackets", "at least one bracket is required")
	}
	lastThreshold := 0.0
	for i, bracket := range t.FederalBrackets {
		if bracket.Threshold < lastThreshold {
			verr.Add("federal_brackets", "brackets must be sorted by threshold")
			break
		}
		if bracket.Rate < 0 {
			verr.Add("federal_brackets", "bracket "+strconv.Itoa(i)+" rate cannot be negative")
		}
		lastThreshold = bracket.Threshold
	}
	if verr.Empty() {
		return nil
	}
	return &verr
}
