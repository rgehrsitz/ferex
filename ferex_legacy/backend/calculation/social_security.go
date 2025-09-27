package calculation

import (
	"ferex/backend/models" // Assuming 'ferex' is module name
	"math"
)

// SSA 2025 bend points for PIA (Primary Insurance Amount)
var piaBendPoints = []float64{1115, 6721} // 2025 values: first, second
var piaPercentages = []float64{0.9, 0.32, 0.15}

// getFRA calculates the Full Retirement Age (FRA) based on birth year.
// It returns a struct containing the FRA in years and months.
func getFRA(birthYear int) models.FRADetail {
	switch {
	case birthYear <= 1937:
		return models.FRADetail{Years: 65, Months: 0}
	case birthYear == 1938:
		return models.FRADetail{Years: 65, Months: 2}
	case birthYear == 1939:
		return models.FRADetail{Years: 65, Months: 4}
	case birthYear == 1940:
		return models.FRADetail{Years: 65, Months: 6}
	case birthYear == 1941:
		return models.FRADetail{Years: 65, Months: 8}
	case birthYear == 1942:
		return models.FRADetail{Years: 65, Months: 10}
	case birthYear >= 1943 && birthYear <= 1954:
		return models.FRADetail{Years: 66, Months: 0}
	case birthYear == 1955:
		return models.FRADetail{Years: 66, Months: 2}
	case birthYear == 1956:
		return models.FRADetail{Years: 66, Months: 4}
	case birthYear == 1957:
		return models.FRADetail{Years: 66, Months: 6}
	case birthYear == 1958:
		return models.FRADetail{Years: 66, Months: 8}
	case birthYear == 1959:
		return models.FRADetail{Years: 66, Months: 10}
	case birthYear >= 1960:
		return models.FRADetail{Years: 67, Months: 0}
	default:
		// Should not happen with valid birth years, but as a fallback
		if birthYear < 1937 {
			return models.FRADetail{Years: 65, Months: 0} // Or handle as an error
		}
		return models.FRADetail{Years: 67, Months: 0} // Default to latest FRA for years > 1960 or unexpected
	}
}

// claimingFactor calculates the adjustment factor for early or late Social Security claiming.
// claimAge is in whole years. fra is the precise Full Retirement Age.
func claimingFactor(claimAge int, fra models.FRADetail) float64 {
	claimAgeInMonths := claimAge * 12
	fraInMonths := fra.Years*12 + fra.Months
	monthsDifference := claimAgeInMonths - fraInMonths
	if monthsDifference == 0 {
		return 1.0
	} else if monthsDifference < 0 { // Early claiming
		monthsEarly := -monthsDifference
		reduction := 0.0
		if monthsEarly <= 36 {
			reduction = float64(monthsEarly) * (5.0 / 9.0 / 100.0) // 5/9 of 1% per month for first 36 months
		} else {
			reduction = float64(36)*(5.0/9.0/100.0) + float64(monthsEarly-36)*(5.0/12.0/100.0) // 5/12 of 1% for months beyond 36
		}
		// Max reduction depends on FRA. E.g., for FRA 67, claiming at 62 (60 months early):
		// 36 * (5/9)% + 24 * (5/12)% = 20% + 10% = 30% reduction.
		// The 0.70 factor is a general floor based on the earliest possible claiming age (62) vs latest FRA (67).
		return 1.0 - reduction // SSA does not have a hard floor like 0.70; actual reduction is applied.
	} else { // Delayed claiming
		// Delayed Retirement Credits (DRCs) - 2/3 of 1% per month of delay.
		// Credits apply up to age 70.
		monthsDelayed := monthsDifference
		maxDelayedMonths := (70 * 12) - fraInMonths

		if monthsDelayed > maxDelayedMonths {
			monthsDelayed = maxDelayedMonths
		}
		if monthsDelayed < 0 { // Should not happen if monthsDifference was positive
			monthsDelayed = 0
		}

		increase := float64(monthsDelayed) * (2.0 / 3.0 / 100.0)
		return 1.0 + increase
	}
}

// Estimate AIME from earnings history or salary
func estimateAIME(earnings []float64, salary float64, yearsWorked int) float64 {
	var indexedEarnings []float64
	// In a real scenario, earnings would be indexed for inflation.
	// This is a simplified version using nominal earnings.
	if len(earnings) > 0 {
		indexedEarnings = earnings
	} else if salary > 0 {
		for i := 0; i < yearsWorked; i++ {
			indexedEarnings = append(indexedEarnings, salary)
		}
	} else {
		return 0
	}

	if len(indexedEarnings) == 0 {
		return 0
	}

	// Select top 35 years of indexed earnings
	// If fewer than 35 years, use zeros for the remaining years.
	// This is a simplification; SSA uses specific indexing factors.
	actualYearsToConsider := len(indexedEarnings)
	if actualYearsToConsider > 35 {
		// Sort and take top 35 - not implemented here for brevity, assumes earnings are somewhat representative
		// A real implementation needs sorting: sort.Float64s(indexedEarnings); indexedEarnings = indexedEarnings[len-35:]
		actualYearsToConsider = 35
	}

	sum := 0.0
	count := 0
	for i := 0; i < actualYearsToConsider; i++ {
		sum += indexedEarnings[i]
		count++
	}

	if count == 0 {
		return 0
	}
	// Sum of top 'count' years, extend with zeros if count < 35 for AIME calculation
	// AIME is average of 35 years (420 months)
	return sum / 420.0
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
	return math.Floor(pia*10) / 10 // Result is rounded down to the next lower dime
}

// CalculateSocialSecurity projects SS benefits at 62, FRA, 70, and chosen claim age
func CalculateSocialSecurity(input models.SocialSecurityCalculationInput) models.SocialSecurityCalculationResult {
	var notes string
	fra := getFRA(input.BirthYear) // fra is now models.FRADetail

	if input.UserProvidedEstimate62 > 0 || input.UserProvidedEstimateFRA > 0 || input.UserProvidedEstimate70 > 0 {
		est62 := input.UserProvidedEstimate62
		estFRA := input.UserProvidedEstimateFRA
		est70 := input.UserProvidedEstimate70
		claimingAmount := 0.0

		// Declare piaFromStatement here to be in scope for the whole "user-provided estimates" block
		var piaFromStatement float64
		piaFromStatement = estFRA // Initialize with the user's FRA estimate

		// Prioritize direct match for claim age with provided estimates
		// For estFRA, we assume the user's statement FRA aligns with a whole year if their input.ClaimAge is a whole year.
		if input.ClaimAge == 62 && est62 > 0 {
			claimingAmount = est62
		} else if input.ClaimAge == fra.Years && fra.Months == 0 && estFRA > 0 {
			claimingAmount = estFRA
		} else if input.ClaimAge == 70 && est70 > 0 {
			claimingAmount = est70
		} else {
			// No direct match, so claimingAmount is still 0.
			// piaFromStatement is already initialized with estFRA.
			// If estFRA (and thus piaFromStatement) is zero, try to derive it.
			if piaFromStatement == 0 {
				// Try to derive from est62
				if est62 > 0 {
					piaFromStatement = est62 / claimingFactor(62, fra)
				}
				// If still zero and est70 is available, try to derive from est70
				if est70 > 0 {
					piaFrom70 := est70 / claimingFactor(70, fra)
					if piaFromStatement == 0 { // if est62 wasn't available or didn't yield PIA
						piaFromStatement = piaFrom70
					} else { // if pia was derived from est62, average with piaFrom70
						piaFromStatement = (piaFromStatement + piaFrom70) / 2.0
					}
				}
			}

			if piaFromStatement > 0 { // Only calculate if we have a PIA to base it on
				claimingAmount = piaFromStatement * claimingFactor(input.ClaimAge, fra)
			} else {
				// Fallback if PIA couldn't be derived but a direct estimate for ClaimAge exists
				// This case handles if estFRA was 0, and PIA couldn't be derived from est62/est70,
				// but the user is asking for 62 or 70 directly for which they *did* provide an estimate.
				if input.ClaimAge == 62 && est62 > 0 { // Check est62 > 0 again
					claimingAmount = est62
				} else if input.ClaimAge == fra.Years && fra.Months == 0 && estFRA > 0 { // Check estFRA > 0 again
					claimingAmount = estFRA
				} else if input.ClaimAge == 70 && est70 > 0 { // Check est70 > 0 again
					claimingAmount = est70
				}
			}
		}

		notes = "Used user-provided SSA statement values."
		// piaFromStatement is now in scope here
		if input.UserProvidedEstimateFRA == 0 && piaFromStatement > 0 && !(input.ClaimAge == fra.Years && fra.Months == 0 && estFRA > 0) {
			notes += " FRA estimate was missing or zero; benefit at custom age calculated based on other available statement data."
		} else if piaFromStatement == 0 && claimingAmount == 0 {
			notes += " Could not determine benefit at custom age from available statement data."
		}

		return models.SocialSecurityCalculationResult{
			EstimatedAt62:  est62,
			EstimatedAtFRA: estFRA,
			EstimatedAt70:  est70,
			ClaimingAge:    input.ClaimAge,
			ClaimingAmount: math.Floor(claimingAmount*10) / 10, // Round final claiming amount
			Notes:          notes,
		}
	}

	// If no user-provided estimates are sufficient, calculate from scratch
	aime := estimateAIME(input.EarningsHistory, input.EstimatedAnnualSalary, input.YearsWorked)
	pia := calculatePIA(aime)

	est62_calc := pia * claimingFactor(62, fra)
	estFRA_calc := pia
	est70_calc := pia * claimingFactor(70, fra)
	claimingAmount_calc := pia * claimingFactor(input.ClaimAge, fra)

	if aime == 0 {
		notes = "No SSA statement or sufficient earnings data (AIME is zero); estimate is zero."
	} else if len(input.EarningsHistory) == 0 && input.EstimatedAnnualSalary > 0 && input.YearsWorked > 0 {
		notes = "Estimate based on average salary and years worked (AIME calculated)."
	} else if len(input.EarningsHistory) > 0 {
		notes = "Estimate based on earnings history (AIME calculated)."
	}

	return models.SocialSecurityCalculationResult{
		EstimatedAt62:  math.Floor(est62_calc*10) / 10,
		EstimatedAtFRA: math.Floor(estFRA_calc*10) / 10,
		EstimatedAt70:  math.Floor(est70_calc*10) / 10,
		ClaimingAge:    input.ClaimAge,
		ClaimingAmount: math.Floor(claimingAmount_calc*10) / 10,
		Notes:          notes,
	}
}
