package calculation

import (
	"ferex/backend/models"
)

// getSurvivorReduction returns the reduction and survivor percent based on pension type and election
func getSurvivorReduction(pensionType, election string) (reduction, survivorPct float64, notes string) {
	switch pensionType {
	case "FERS":
		switch election {
		case "max":
			return 0.10, 0.50, "FERS max: 50% to survivor, 10% reduction"
		case "partial":
			return 0.05, 0.25, "FERS partial: 25% to survivor, 5% reduction"
		default:
			return 0.0, 0.0, "FERS: No survivor benefit elected"
		}
	case "CSRS", "CSRSOffset":
		switch election {
		case "max":
			return 0.10, 0.55, "CSRS max: 55% to survivor, ~10% reduction"
		case "partial":
			return 0.05, 0.50, "CSRS partial: 50% to survivor, 5% reduction"
		default:
			return 0.0, 0.0, "CSRS: No survivor benefit elected"
		}
	default:
		return 0.0, 0.0, "Unknown pension type"
	}
}

// CalculateSurvivorBenefit projects survivor annuity/income
func CalculateSurvivorBenefit(input models.SurvivorBenefitCalculationInput) models.SurvivorBenefitCalculationResult {
	reduction, survivorPct, notes := getSurvivorReduction(input.PensionType, input.SurvivorElection)
	baseAnnuity := input.InitialAnnuity * (1 - reduction)
	initialSurvivor := baseAnnuity * survivorPct
	projected := make([]float64, input.YearsToProject)
	current := initialSurvivor
	total := 0.0
	for i := 0; i < input.YearsToProject; i++ {
		if i > 0 {
			current = current * (1 + input.COLARate)
		}
		ann := current
		// Add SS survivor if included
		if input.IncludeSSSurvivor {
			ann += input.SSSurvivorAmount
		}
		// Add TSP if included (simple: spread evenly over projection)
		if input.IncludeTSP && input.TSPBalanceAtDeath > 0 && input.YearsToProject > 0 {
			ann += input.TSPBalanceAtDeath / float64(input.YearsToProject)
		}
		// Add other survivor income
		if input.OtherSurvivorIncome > 0 {
			ann += input.OtherSurvivorIncome
		}
		projected[i] = ann
		total += ann
	}
	return models.SurvivorBenefitCalculationResult{
		InitialSurvivorAnnuity: initialSurvivor,
		ProjectedAnnuities:     projected,
		TotalSurvivorIncome:    total,
		Notes:                  notes,
	}
}
