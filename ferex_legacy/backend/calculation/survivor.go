package calculation

import (
	"ferex/backend/models" // Assuming 'ferex' is module name
)

// getSurvivorReduction calculates the financial impact of electing a survivor benefit.
// It returns:
// - retireeAnnuityActualReductionAmount: The actual monetary amount the retiree's annuity is reduced by.
// - survivorBenefitFraction: The percentage (as a decimal) of an annuity amount the survivor will receive (e.g., 0.55 for 55%).
// - survivorBenefitBasedOnUnreduced: A boolean. True if survivorBenefitFraction applies to the retiree's original, unreduced annuity (CSRS). False if it applies to the retiree's reduced annuity (FERS).
// - notes: Descriptive information about the calculation.
func getSurvivorReduction(pensionType, election string, initialAnnuity float64) (retireeAnnuityActualReductionAmount, survivorBenefitFraction float64, survivorBenefitBasedOnUnreduced bool, notes string) {
	switch pensionType {
	case "FERS":
		switch election {
		case "max": // 50% survivor benefit, 10% cost to retiree
			actualReduction := initialAnnuity * 0.10
			return actualReduction, 0.50, false, "FERS max: Survivor receives 50% of retiree's reduced annuity. Retiree's annuity reduced by 10%."
		case "partial": // 25% survivor benefit, 5% cost to retiree
			actualReduction := initialAnnuity * 0.05
			return actualReduction, 0.25, false, "FERS partial: Survivor receives 25% of retiree's reduced annuity. Retiree's annuity reduced by 5%."
		default:
			return 0.0, 0.0, false, "FERS: No survivor benefit elected"
		}
	case "CSRS", "CSRSOffset":
		switch election {
		case "max":
			// Max survivor benefit is 55% of unreduced annuity.
			// Cost: 2.5% of the first $3600 of unreduced annuity + 10% of unreduced annuity over $3600.
			cost := 0.0
			if initialAnnuity <= 0 {
			    cost = 0.0
			} else if initialAnnuity <= 3600 {
				cost = initialAnnuity * 0.025
			} else {
				cost = (3600 * 0.025) + ((initialAnnuity - 3600) * 0.10)
			}
			// Ensure cost does not exceed the benefit or the annuity itself (edge case for very low annuities)
			if cost > initialAnnuity { cost = initialAnnuity } // Cost cannot be more than the annuity
			return cost, 0.55, true, "CSRS max: Survivor receives 55% of retiree's unreduced annuity. Retiree cost based on formula (2.5% of first $3600 + 10% over $3600)."
		case "partial":
			// Survivor to receive 25% of the retiree's unreduced annuity.
			// The cost for this benefit is calculated based on a 'deemed' base amount.
			// If survivor gets 25% of unreduced_X, and this 25%_unreduced_X must be 55% of some Base_Y,
			// then Base_Y = (0.25 * unreduced_X) / 0.55.
			targetSurvivorBenefitFractionOfUnreduced := 0.25
			if initialAnnuity <= 0 {
			    return 0.0, 0.0, true, "CSRS partial: No benefit or cost for zero/negative initial annuity."
			}
			baseForCostFormula := (targetSurvivorBenefitFractionOfUnreduced * initialAnnuity) / 0.55

			cost := 0.0
			if baseForCostFormula <= 0 {
			    cost = 0.0
			} else if baseForCostFormula <= 3600 {
				cost = baseForCostFormula * 0.025
			} else {
				cost = (3600 * 0.025) + ((baseForCostFormula - 3600) * 0.10)
			}
			// Ensure cost does not exceed the annuity itself.
			if cost > initialAnnuity { cost = initialAnnuity }
			if cost < 0 { cost = 0 } // Cost cannot be negative

			return cost, targetSurvivorBenefitFractionOfUnreduced, true, "CSRS partial: Survivor receives 25% of retiree's unreduced annuity. Cost calculated based on derived base for CSRS formula."
		default:
			return 0.0, 0.0, false, "CSRS: No survivor benefit elected"
		}
	default:
		return 0.0, 0.0, false, "Unknown pension type"
	}
}

// CalculateSurvivorBenefit projects survivor annuity/income
func CalculateSurvivorBenefit(input models.SurvivorBenefitCalculationInput) models.SurvivorBenefitCalculationResult {
	retireeAnnuityActualReductionAmount, survivorBenefitFraction, survivorBenefitBasedOnUnreduced, notes := getSurvivorReduction(input.PensionType, input.SurvivorElection, input.InitialAnnuity)

	retireeReducedAnnuity := input.InitialAnnuity - retireeAnnuityActualReductionAmount
	var initialSurvivorAnnuity float64

	if survivorBenefitBasedOnUnreduced {
		// Typically CSRS: survivor benefit is a fraction of the retiree's UNREDUCED annuity.
		initialSurvivorAnnuity = input.InitialAnnuity * survivorBenefitFraction
	} else {
		// Typically FERS (and legacy simplified CSRS partial): survivor benefit is a fraction of the retiree's REDUCED annuity.
		initialSurvivorAnnuity = retireeReducedAnnuity * survivorBenefitFraction
	}

	// Ensure retiree's reduced annuity isn't negative (if reduction somehow exceeds annuity)
	if retireeReducedAnnuity < 0 {
	    retireeReducedAnnuity = 0
	    // If retiree annuity is zero, survivor annuity from pension should also be zero unless based on unreduced amount that was positive
	    if !survivorBenefitBasedOnUnreduced { // if based on reduced, and reduced is now 0, then this part is 0
	        initialSurvivorAnnuity = 0
	    }
	}
	// Ensure initial survivor annuity isn't negative
	if initialSurvivorAnnuity < 0 {
	    initialSurvivorAnnuity = 0
	}

	projectedSurvivorIncome := make([]float64, input.YearsToProject)
	currentAnnualSurvivorIncome := initialSurvivorAnnuity
	totalProjectedSurvivorIncome := 0.0

	for i := 0; i < input.YearsToProject; i++ {
		if i > 0 {
			currentAnnualSurvivorIncome *= (1 + input.COLARate)
		}
		annualIncome := currentAnnualSurvivorIncome

		if input.IncludeSSSurvivor {
			annualIncome += input.SSSurvivorAmount // Assuming SSSurvivorAmount is annual
		}
		if input.IncludeTSP && input.TSPBalanceAtDeath > 0 && input.YearsToProject > 0 {
			annualIncome += input.TSPBalanceAtDeath / float64(input.YearsToProject) // Simple TSP spread
		}
		if input.OtherSurvivorIncome > 0 {
			annualIncome += input.OtherSurvivorIncome // Assuming OtherSurvivorIncome is annual
		}
		projectedSurvivorIncome[i] = annualIncome
		totalProjectedSurvivorIncome += annualIncome
	}

	return models.SurvivorBenefitCalculationResult{
		InitialSurvivorAnnuity: initialSurvivorAnnuity,
		ProjectedAnnuities:     projectedSurvivorIncome,
		TotalSurvivorIncome:    totalProjectedSurvivorIncome,
		Notes:                  notes,
	}
}

// CalculateSurvivorBenefits is a stub function.
// It currently calls the existing CalculateSurvivorBenefit function,
// ignoring fersResult and csrsResult as the existing logic in CalculateSurvivorBenefit
// primarily uses fields already present in survivorInput (like InitialAnnuity and PensionType).
// TODO: Review if fersResult and csrsResult are needed for more detailed survivor calculations.
func CalculateSurvivorBenefits(survivorInput models.SurvivorBenefitCalculationInput, fersResult models.FERSCalculationResult, csrsResult models.CSRSCalculationResult) models.SurvivorBenefitCalculationResult {
	// Call the existing function that matches the core logic for now
	return CalculateSurvivorBenefit(survivorInput)
}
