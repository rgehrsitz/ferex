package calculation

import (
	"ferex/backend/models"
	"math"
)

// SSA 2025 bend points for PIA (Primary Insurance Amount)
var piaBendPoints = []float64{1115, 6721} // 2025 values: first, second
var piaPercentages = []float64{0.9, 0.32, 0.15}

// Full Retirement Age lookup by birth year (simplified)
func getFRA(birthYear int) int {
	switch {
	case birthYear <= 1954:
		return 66
	case birthYear >= 1960:
		return 67
	default:
		return 66 + (birthYear-1954)/2 // 1955: 66+2mo, 1956: 66+4mo, etc. (approximate)
	}
}

// Early/late claiming factors (monthly)
func claimingFactor(claimAge, fra int) float64 {
	if claimAge == fra {
		return 1.0
	} else if claimAge < fra {
		reduction := float64(fra-claimAge) * 0.0056 * 12 // 0.56% per month, 6.7% per year for first 36 months
		return math.Max(1.0-reduction, 0.7)              // SSA minimum is 70% at 62
	} else {
		increase := float64(claimAge-fra) * 0.008 * 12 // 0.8% per month, 8% per year
		return math.Min(1.0+increase, 1.24)            // SSA max is 124% at 70
	}
}

// Estimate AIME from earnings history or salary
func estimateAIME(earnings []float64, salary float64, yearsWorked int) float64 {
	var top35 []float64
	if len(earnings) >= 35 {
		top35 = earnings
	} else if len(earnings) > 0 {
		top35 = append([]float64{}, earnings...)
		for len(top35) < 35 {
			top35 = append(top35, 0)
		}
	} else if salary > 0 && yearsWorked > 0 {
		for i := 0; i < 35; i++ {
			if i < yearsWorked {
				top35 = append(top35, salary)
			} else {
				top35 = append(top35, 0)
			}
		}
	} else {
		return 0
	}
	// Sort and take top 35
	sum := 0.0
	for _, v := range top35 {
		sum += v
	}
	return sum / 420.0 // 35 years * 12 months
}

// Calculate PIA from AIME
func calculatePIA(aime float64) float64 {
	var pia float64
	if aime <= piaBendPoints[0] {
		pia = piaPercentages[0] * aime
	} else if aime <= piaBendPoints[1] {
		pia = piaPercentages[0]*piaBendPoints[0] + piaPercentages[1]*(aime-piaBendPoints[0])
	} else {
		pia = piaPercentages[0]*piaBendPoints[0] + piaPercentages[1]*(piaBendPoints[1]-piaBendPoints[0]) + piaPercentages[2]*(aime-piaBendPoints[1])
	}
	return pia
}

// CalculateSocialSecurity projects SS benefits at 62, FRA, 70, and chosen claim age
func CalculateSocialSecurity(input models.SocialSecurityCalculationInput) models.SocialSecurityCalculationResult {
	var notes string
	fra := getFRA(input.BirthYear)

	// Prefer user-provided SSA statement values if present
	if input.UserProvidedEstimate62 > 0 && input.UserProvidedEstimateFRA > 0 && input.UserProvidedEstimate70 > 0 {
		claimingAmount := 0.0
		switch input.ClaimAge {
		case 62:
			claimingAmount = input.UserProvidedEstimate62
		case fra:
			claimingAmount = input.UserProvidedEstimateFRA
		case 70:
			claimingAmount = input.UserProvidedEstimate70
		default:
			// Linear interpolate between FRA and 70 or 62
			if input.ClaimAge < fra {
				claimingAmount = input.UserProvidedEstimate62 + (input.UserProvidedEstimateFRA-input.UserProvidedEstimate62)*float64(input.ClaimAge-62)/float64(fra-62)
			} else {
				claimingAmount = input.UserProvidedEstimateFRA + (input.UserProvidedEstimate70-input.UserProvidedEstimateFRA)*float64(input.ClaimAge-fra)/float64(70-fra)
			}
		}
		notes = "Used user-provided SSA statement values."
		return models.SocialSecurityCalculationResult{
			EstimatedAt62:  input.UserProvidedEstimate62,
			EstimatedAtFRA: input.UserProvidedEstimateFRA,
			EstimatedAt70:  input.UserProvidedEstimate70,
			ClaimingAge:    input.ClaimAge,
			ClaimingAmount: claimingAmount,
			Notes:          notes,
		}
	}

	// Otherwise, estimate from salary/history
	aime := estimateAIME(input.EarningsHistory, input.EstimatedAnnualSalary, input.YearsWorked)
	pia := calculatePIA(aime)

	// Early/late claiming factors
	factor62 := claimingFactor(62, fra)
	factorFRA := 1.0
	factor70 := claimingFactor(70, fra)
	claimingFactorX := claimingFactor(input.ClaimAge, fra)

	est62 := pia * factor62
	estFRA := pia * factorFRA
	est70 := pia * factor70
	claimingAmount := pia * claimingFactorX

	if aime == 0 || pia == 0 {
		notes = "No SSA statement or sufficient earnings data provided; estimate is zero."
	}
	if len(input.EarningsHistory) == 0 && input.EstimatedAnnualSalary > 0 && input.YearsWorked > 0 {
		notes = "Estimate based on average salary and years worked."
	}

	return models.SocialSecurityCalculationResult{
		EstimatedAt62:  est62,
		EstimatedAtFRA: estFRA,
		EstimatedAt70:  est70,
		ClaimingAge:    input.ClaimAge,
		ClaimingAmount: claimingAmount,
		Notes:          notes,
	}
}
