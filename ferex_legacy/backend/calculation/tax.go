package calculation

import (
	"ferex/backend/config"
	"ferex/backend/models"
	"fmt"
	"math"
)

// CalculateTax estimates federal and state tax on retirement income using configuration data.
func CalculateTax(input models.TaxCalculationInput) models.TaxCalculationResult {
	// Load configuration
	cfg, err := config.GetFerexConfig()
	if err != nil {
		return models.TaxCalculationResult{
			GrossIncome:       0,
			TaxableIncome:     0,
			FederalTaxOwed:    0,
			StateTaxOwed:      0,
			NetAfterTaxIncome: 0,
			EffectiveTaxRate:  0,
			Notes:             fmt.Sprintf("Error loading configuration: %v", err),
		}
	}

	var notes string

	// Subtract IRS Simplified Method exclusion from taxable pension if provided
	taxablePension := input.GrossPension
	if input.TaxablePension > 0 {
		taxablePension = input.TaxablePension
	}
	if input.IRSSimplifiedMethodExclusion > 0 && taxablePension > 0 {
		taxablePension = taxablePension - input.IRSSimplifiedMethodExclusion
		if taxablePension < 0 {
			taxablePension = 0
		}
		notes += fmt.Sprintf("IRS Simplified Method exclusion applied: $%.2f deducted from taxable pension. ", input.IRSSimplifiedMethodExclusion)
	}

	// Get tax brackets for filing status
	brackets, ok := cfg.FederalTax.Brackets[input.FilingStatus]
	if !ok {
		return models.TaxCalculationResult{
			GrossIncome:       0,
			TaxableIncome:     0,
			FederalTaxOwed:    0,
			StateTaxOwed:      0,
			NetAfterTaxIncome: 0,
			EffectiveTaxRate:  0,
			Notes:             fmt.Sprintf("Unknown filing status: %s", input.FilingStatus),
		}
	}

	// Calculate deductions
	deduction := input.Deductions
	if deduction == 0 {
		// Use standard deduction
		if stdDeduction, ok := cfg.FederalTax.StandardDeduction[input.FilingStatus]; ok {
			deduction = stdDeduction

			// Add additional standard deduction for age 65+ or blind
			if input.Age >= 65 {
				if ageDeduction, ok := cfg.FederalTax.AdditionalStandardDeduction.Age65OrOlder[input.FilingStatus]; ok {
					deduction += ageDeduction
					notes += fmt.Sprintf("Additional standard deduction for age 65+: $%.0f. ", ageDeduction)
				}
			}
			if input.IsBlind {
				if blindDeduction, ok := cfg.FederalTax.AdditionalStandardDeduction.Blind[input.FilingStatus]; ok {
					deduction += blindDeduction
					notes += fmt.Sprintf("Additional standard deduction for blind: $%.0f. ", blindDeduction)
				}
			}
		}
	}

	// Calculate Social Security taxability using configuration thresholds
	agi := taxablePension + input.TSPWithdrawal + input.OtherTaxableIncome
	provisional := agi + 0.5*input.SocialSecurity
	ssTaxable := 0.0

	if ssThresholds, ok := cfg.FederalTax.SocialSecurityTaxabilityThresholds[input.FilingStatus]; ok {
		if provisional > ssThresholds.Tier2 {
			ssTaxable = math.Min(0.85*input.SocialSecurity,
				0.85*(provisional-ssThresholds.Tier2)+math.Min(0.5*input.SocialSecurity,
					0.5*(ssThresholds.Tier2-ssThresholds.Tier1)))
		} else if provisional > ssThresholds.Tier1 {
			ssTaxable = math.Min(0.5*input.SocialSecurity, 0.5*(provisional-ssThresholds.Tier1))
		}
	}

	taxableIncome := agi + ssTaxable - deduction
	if taxableIncome < 0 {
		taxableIncome = 0
	}

	// Calculate federal tax using progressive brackets
	fedTax := 0.0
	for _, bracket := range brackets {
		if taxableIncome <= bracket.Min {
			break
		}

		taxableAtThisBracket := math.Min(taxableIncome, bracket.Max) - bracket.Min
		if taxableAtThisBracket > 0 {
			fedTax += taxableAtThisBracket * bracket.Rate
		}
	}

	// Apply tax credits
	fedTax -= input.TaxCredits
	if fedTax < 0 {
		fedTax = 0
	}

	// Calculate state tax using configuration
	stateTax := 0.0
	if input.StateOfResidence != "" {
		if stateInfo, ok := cfg.StateTax.States[input.StateOfResidence]; ok {
			if stateInfo.HasStateTax {
				// Calculate state taxable income based on state rules
				stateTaxableIncome := 0.0

				// Federal pension taxation by state
				if stateInfo.FederalPensionTax == "taxed" {
					stateTaxableIncome += taxablePension
				} else if stateInfo.FederalPensionTax == "partial" {
					stateTaxableIncome += taxablePension * 0.5 // Simplified partial taxation
				}

				// TSP taxation by state
				if stateInfo.TSPTax == "taxed" {
					stateTaxableIncome += input.TSPWithdrawal
				}

				// Social Security taxation by state
				if stateInfo.SocialSecurityTax == "taxed" {
					stateTaxableIncome += ssTaxable
				} else if stateInfo.SocialSecurityTax == "partial" {
					stateTaxableIncome += ssTaxable * 0.5
				}

				// Other taxable income
				stateTaxableIncome += input.OtherTaxableIncome

				// Apply state effective rate
				stateTax = stateTaxableIncome * (stateInfo.EffectiveRate / 100)
				notes += fmt.Sprintf("State tax calculated for %s at %.2f%% effective rate. %s ",
					input.StateOfResidence, stateInfo.EffectiveRate, stateInfo.Notes)
			} else {
				notes += fmt.Sprintf("No state income tax for %s. ", input.StateOfResidence)
			}
		} else {
			// Fallback to user-provided state taxable income if state not in config
			if input.StateTaxableIncome > 0 {
				stateTax = input.StateTaxableIncome * 0.05
				notes += fmt.Sprintf("State tax estimated at 5%% for %s (not in configuration). ", input.StateOfResidence)
			}
		}
	}

	netIncome := input.GrossPension + input.TSPWithdrawal + input.TSPRothWithdrawal + input.SocialSecurity + input.OtherTaxableIncome - fedTax - stateTax
	effectiveRate := 0.0
	totalIncome := input.GrossPension + input.TSPWithdrawal + input.TSPRothWithdrawal + input.SocialSecurity + input.OtherTaxableIncome

	// Debug: Log tax income components
	// log.Printf("DEBUG TAX CALC - Pension: $%.0f, TSP: $%.0f, TSP Roth: $%.0f, SS: $%.0f, Other: $%.0f, Total: $%.0f",
	// 	input.GrossPension, input.TSPWithdrawal, input.TSPRothWithdrawal, input.SocialSecurity, input.OtherTaxableIncome, totalIncome)
	if totalIncome > 0 {
		effectiveRate = (fedTax + stateTax) / totalIncome
	}

	return models.TaxCalculationResult{
		GrossIncome:       totalIncome,
		TaxableIncome:     taxableIncome,
		FederalTaxOwed:    fedTax,
		StateTaxOwed:      stateTax,
		NetAfterTaxIncome: netIncome,
		EffectiveTaxRate:  effectiveRate,
		Notes:             notes,
	}
}

// CalculateFederalTax is a stub wrapper around CalculateTax.
// It's intended to satisfy calls expecting federal tax calculation separately.
// The actual calculation is done in CalculateTax.
func CalculateFederalTax(input models.TaxCalculationInput) models.TaxCalculationResult {
	// For now, this stub simply calls the existing comprehensive CalculateTax function.
	// The caller is expected to use the FederalTaxOwed field from the result.
	// TODO: Consider refactoring to avoid redundant calculations if both federal and state are needed.
	// log.Printf("Stub: CalculateFederalTax called with year %d, GrossPension: %.2f, SS: %.2f, TSP: %.2f", input.TaxYear, input.GrossPension, input.SocialSecurity, input.TSPWithdrawal)
	return CalculateTax(input)
}

// CalculateStateTax is a stub wrapper around CalculateTax.
// It's intended to satisfy calls expecting state tax calculation separately.
// The actual calculation is done in CalculateTax.
func CalculateStateTax(input models.TaxCalculationInput) models.TaxCalculationResult {
	// For now, this stub simply calls the existing comprehensive CalculateTax function.
	// The caller is expected to use the StateTaxOwed field from the result.
	// TODO: Consider refactoring to avoid redundant calculations if both federal and state are needed.
	// log.Printf("Stub: CalculateStateTax called with year %d, GrossPension: %.2f, SS: %.2f, TSP: %.2f", input.TaxYear, input.GrossPension, input.SocialSecurity, input.TSPWithdrawal)
	return CalculateTax(input)
}
