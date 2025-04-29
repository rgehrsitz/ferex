package calculation

import (
	"ferex/backend/models"
	"fmt"
)

// CalculateHealthPremiums projects FEHB/Medicare/other premiums over time
func CalculateHealthPremiums(input models.HealthPremiumCalculationInput) models.HealthPremiumCalculationResult {
	premiums := make([]float64, input.YearsToProject)
	notes := ""
	currentFEHB := input.FEHBPremium
	currentMedicare := input.MedicarePremium
	currentOther := input.OtherHealthPremium
	total := 0.0

	for i := 0; i < input.YearsToProject; i++ {
		annual := 0.0
		if input.IncludeFEHB {
			annual += currentFEHB
		}
		if input.IncludeMedicare {
			annual += currentMedicare
		}
		if currentOther > 0 {
			annual += currentOther
		}
		premiums[i] = annual
		total += annual
		// Apply COLA for next year
		currentFEHB *= (1 + input.COLARate)
		currentMedicare *= (1 + input.COLARate)
		currentOther *= (1 + input.COLARate)
	}

	if input.IncludeFEHB {
		notes += "FEHB included. "
	} else {
		notes += "FEHB not included. "
	}
	if input.IncludeMedicare {
		notes += "Medicare included. "
	} else {
		notes += "Medicare not included. "
	}
	if input.COLARate > 0 {
		notes += fmt.Sprintf("COLA/inflation rate applied: %.2f%%. ", input.COLARate*100)
	}

	return models.HealthPremiumCalculationResult{
		ProjectedPremiums: premiums,
		TotalPremiums:     total,
		Notes:             notes,
	}
}
