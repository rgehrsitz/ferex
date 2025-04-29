package calculation

import (
	"ferex/backend/models"
)

// CalculateRetirementProjection orchestrates all modules for a full retirement projection
func CalculateRetirementProjection(input models.RetirementCalculationInput) models.RetirementCalculationResult {
	fersResult := CalculateFERSPension(input.FERSInput)
	csrsResult := CalculateCSRS(input.CSRSInput)
	srsResult := CalculateSRS(input.SRSInput)
	tspResult := CalculateTSP(input.TSPInput)
	taxResult := CalculateTax(input.TaxInput)
	ssResult := CalculateSocialSecurity(input.SocialSecurityInput)
	colaResult := CalculateCOLA(input.COLAInput)
	survivorResult := CalculateSurvivorBenefit(input.SurvivorInput)
	healthResult := CalculateHealthPremiums(input.HealthInput)

	// Aggregate income streams (example: sum of annuities, SS, TSP withdrawals)
	totalIncome := fersResult.AnnualPension + csrsResult.AnnualPension + srsResult.AnnualSRSAmount + tspResult.AnnualWithdrawalIncome + ssResult.ClaimingAmount

	// Subtract health premiums
	netIncome := totalIncome
	if len(healthResult.ProjectedPremiums) > 0 {
		netIncome -= healthResult.ProjectedPremiums[0]
	}
	// Subtract taxes
	if taxResult.FederalTaxOwed > 0 {
		netIncome -= taxResult.FederalTaxOwed
	}
	if taxResult.StateTaxOwed > 0 {
		netIncome -= taxResult.StateTaxOwed
	}

	notes := fersResult.Notes + "\n" + csrsResult.Notes + "\n" + srsResult.Notes + "\n" + tspResult.Notes + "\n" + taxResult.Notes + "\n" + ssResult.Notes + "\n" + colaResult.Notes + "\n" + survivorResult.Notes + "\n" + healthResult.Notes

	var monteCarloResult models.MonteCarloResult
	if input.MonteCarloInput.NumSimulations > 0 {
		monteCarloResult = RunMonteCarlo(input.MonteCarloInput)
	}

	return models.RetirementCalculationResult{
		FERSResult:            fersResult,
		CSRSResult:            csrsResult,
		SRSResult:             srsResult,
		TSPResult:             tspResult,
		TaxResult:             taxResult,
		SocialSecurityResult:  ssResult,
		COLAResult:            colaResult,
		SurvivorResult:        survivorResult,
		HealthResult:          healthResult,
		MonteCarloResult:      monteCarloResult,
		NetAfterTaxIncome:     netIncome,
		EffectiveTaxRate:      taxResult.EffectiveTaxRate,
		TotalRetirementIncome: totalIncome,
		Notes:                 notes,
	}
}
