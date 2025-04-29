package calculation

import (
	"ferex/backend/models"
	"fmt"
)

// Simplified 2025 federal tax brackets (for demonstration)
var federalBrackets2025 = map[string][]struct {
	Rate   float64
	Amount float64
}{
	"single": {
		{0.10, 11000},
		{0.12, 44725 - 11000},
		{0.22, 95375 - 44725},
		{0.24, 182100 - 95375},
		{0.32, 231250 - 182100},
		{0.35, 578125 - 231250},
		{0.37, 1e12}, // top bracket
	},
	"married": {
		{0.10, 22000},
		{0.12, 89450 - 22000},
		{0.22, 190750 - 89450},
		{0.24, 364200 - 190750},
		{0.32, 462500 - 364200},
		{0.35, 693750 - 462500},
		{0.37, 1e12},
	},
}

// Standard deduction for 2025
var standardDeduction2025 = map[string]float64{
	"single":  14600,
	"married": 29200,
}

// CalculateTax estimates federal and state tax on retirement income.
func CalculateTax(input models.TaxCalculationInput) models.TaxCalculationResult {
	var notes string
	brackets, ok := federalBrackets2025[input.FilingStatus]
	if !ok {
		return models.TaxCalculationResult{
			FederalTaxOwed:    0,
			StateTaxOwed:      0,
			NetAfterTaxIncome: 0,
			EffectiveTaxRate:  0,
			Notes:             fmt.Sprintf("Unknown filing status: %s", input.FilingStatus),
		}
	}

	deduction := input.Deductions
	if deduction == 0 {
		if d, ok := standardDeduction2025[input.FilingStatus]; ok {
			deduction = d
		}
	}

	// Social Security taxability (simplified): up to 85% taxable
	// Provisional income = AGI + 0.5*SS + tax-exempt interest (ignored here)
	agi := input.TaxablePension + input.TSPWithdrawal + input.OtherTaxableIncome
	provisional := agi + 0.5*input.SocialSecurity
	ssTaxable := 0.0
	if provisional > 34000 && input.FilingStatus == "married" {
		ssTaxable = 0.85 * input.SocialSecurity
	} else if provisional > 25000 && input.FilingStatus == "single" {
		ssTaxable = 0.85 * input.SocialSecurity
	} else if provisional > 32000 && input.FilingStatus == "married" {
		ssTaxable = 0.5 * input.SocialSecurity
	} else if provisional > 25000 && input.FilingStatus == "single" {
		ssTaxable = 0.5 * input.SocialSecurity
	}

	taxableIncome := agi + ssTaxable - deduction
	if taxableIncome < 0 {
		taxableIncome = 0
	}

	// Calculate federal tax
	fedTax := 0.0
	incomeLeft := taxableIncome
	for _, b := range brackets {
		amount := b.Amount
		if incomeLeft < amount {
			fedTax += incomeLeft * b.Rate
			break
		}
		fedTax += amount * b.Rate
		incomeLeft -= amount
	}
	fedTax -= input.TaxCredits
	if fedTax < 0 {
		fedTax = 0
	}

	// State tax: placeholder flat 5% if state is specified
	stateTax := 0.0
	if input.StateOfResidence != "" && input.StateTaxableIncome > 0 {
		stateTax = input.StateTaxableIncome * 0.05
		notes += fmt.Sprintf("State tax estimated at 5%% for %s.\n", input.StateOfResidence)
	}

	netIncome := input.GrossPension + input.TSPWithdrawal + input.TSPRothWithdrawal + input.SocialSecurity + input.OtherTaxableIncome - fedTax - stateTax
	effectiveRate := 0.0
	totalIncome := input.GrossPension + input.TSPWithdrawal + input.TSPRothWithdrawal + input.SocialSecurity + input.OtherTaxableIncome
	if totalIncome > 0 {
		effectiveRate = (fedTax + stateTax) / totalIncome
	}

	return models.TaxCalculationResult{
		FederalTaxOwed:    fedTax,
		StateTaxOwed:      stateTax,
		NetAfterTaxIncome: netIncome,
		EffectiveTaxRate:  effectiveRate,
		Notes:             notes,
	}
}
