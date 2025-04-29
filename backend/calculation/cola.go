package calculation

import (
	"ferex/backend/models"
	"fmt"
)

// getFERSCOLA applies the FERS COLA cap rules per year
func getFERSCOLA(cpi float64) float64 {
	if cpi <= 2.0 {
		return cpi
	} else if cpi <= 3.0 {
		return 2.0
	} else {
		fersCola := cpi - 1.0
		if fersCola < 2.0 {
			return 2.0
		}
		return fersCola
	}
}

// CalculateCOLA projects COLA-adjusted values for pensions, SS, etc.
func CalculateCOLA(input models.COLACalculationInput) models.COLACalculationResult {
	amounts := make([]float64, input.Years)
	current := input.InitialAmount
	cumulative := 0.0
	notes := ""
	policy := input.COLAPolicy
	if policy == "" {
		policy = "Generic"
	}

	for i := 0; i < input.Years; i++ {
		adjRate := input.COLARate
		switch policy {
		case "FERS":
			adjRate = getFERSCOLA(input.COLARate*100) / 100.0
			if i == 0 {
				notes += fmt.Sprintf("FERS COLA cap logic applied: %.2f%%\n", adjRate*100)
			}
		case "CSRS", "SocialSecurity":
			// Full COLA, no cap
		case "None":
			adjRate = 0
		case "Generic":
			// Use as provided
		default:
			// Custom logic if needed
		}
		current = current * (1 + adjRate)
		amounts[i] = current
	}

	if input.Years > 0 {
		cumulative = (amounts[input.Years-1] - input.InitialAmount) / input.InitialAmount
	}

	return models.COLACalculationResult{
		ProjectedAmounts: amounts,
		FinalAmount:      amounts[input.Years-1],
		CumulativeCOLA:   cumulative,
		Notes:            notes,
	}
}
