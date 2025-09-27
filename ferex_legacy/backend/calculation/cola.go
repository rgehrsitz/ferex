package calculation

import (
	"ferex/backend/models" // Added import for our models package
	"math"
	"strings" // Added for strings.Contains and strings.TrimSpace
)

// getCorrectFERSCOLA calculates the FERS COLA rate based on the CPI-W increase rate.
// Both cpiW_IncreaseRate and the return value are rates (e.g., 0.04 for 4%).
func getCorrectFERSCOLA(cpiW_IncreaseRate float64) float64 {
	if cpiW_IncreaseRate <= 0.0 {
		return 0.0 // No COLA for deflation or zero inflation
	} else if cpiW_IncreaseRate <= 0.02 {
		return cpiW_IncreaseRate // COLA = CPI-W if CPI-W is 2% or less
	} else if cpiW_IncreaseRate <= 0.03 {
		return 0.02 // COLA = 2% if CPI-W is >2% and <=3%
	} else {
		return cpiW_IncreaseRate - 0.01 // COLA = CPI-W - 1% if CPI-W is >3%
	}
}

// ProjectPensionWithCOLA projects annual pension amounts with applicable COLAs.
func ProjectPensionWithCOLA(input models.PensionCOLAInput) models.PensionCOLAResult {
	projectedAmounts := make([]float64, input.YearsToProject)
	notes := ""

	if input.YearsToProject <= 0 {
		return models.PensionCOLAResult{ProjectedPension: []float64{}, Notes: "YearsToProject must be positive."}
	}

	for i := 0; i < input.YearsToProject; i++ {
		ageInYear := input.ProjectionStartAge + i
		colaRateToApply := 0.0
		appliedCOLAThisYear := false

		switch input.PensionType {
		case "CSRS":
			colaRateToApply = input.AssumedCPIWRate
			if colaRateToApply > 0 {
				appliedCOLAThisYear = true
			}
			// Add note only once for CSRS if COLA is generally applied
			if i == 0 && appliedCOLAThisYear && notes == "" {
				notes += "CSRS: Full COLA applied each year. "
			}
		case "FERS":
			// FERS COLA is generally not paid until age 62.
			// Exception: Disability, Survivor, Special Provisions (not explicitly handled here yet)
			if ageInYear >= 62 {
				colaRateToApply = getCorrectFERSCOLA(input.AssumedCPIWRate)
				if colaRateToApply > 0 {
					appliedCOLAThisYear = true
				}
				// Add note when FERS COLA first starts applying, or if it applies from the first projection year.
				firstYearOfProjectionIsCOLAEligible := (i == 0 && appliedCOLAThisYear)
				thisYearIsFirstCOLAEligibleYear := (ageInYear == 62 && appliedCOLAThisYear && (input.ProjectionStartAge+i-1 < 62))
				if (firstYearOfProjectionIsCOLAEligible || thisYearIsFirstCOLAEligibleYear) && !strings.Contains(notes, "FERS: COLA applied") {
					notes += "FERS: COLA applied (age >= 62). "
				}
			} else {
				// No COLA if under 62
				if i == 0 && !strings.Contains(notes, "FERS: No COLA") { // Add note only once if projection starts before 62
					notes += "FERS: No COLA (age < 62). "
				}
			}
		default:
			// No COLA for unknown pension types
			colaRateToApply = 0.0
			if i == 0 && !strings.Contains(notes, "Unknown pension type") {
				notes += "Unknown pension type or no COLA policy: No COLA applied. "
			}
		}

		baseForCOLAThisYear := 0.0
		if i == 0 {
			baseForCOLAThisYear = input.InitialPension
		} else {
			baseForCOLAThisYear = projectedAmounts[i-1]
		}

		projectedValue := baseForCOLAThisYear * (1 + colaRateToApply)
		projectedAmounts[i] = math.Round(projectedValue*100) / 100
	}

	return models.PensionCOLAResult{
		ProjectedPension: projectedAmounts,
		Notes:            strings.TrimSpace(notes),
	}
}

// ProjectSocialSecurityWithCOLA projects annual Social Security benefits with COLAs.
func ProjectSocialSecurityWithCOLA(input models.SocialSecurityCOLAInput) models.SocialSecurityCOLAResult {
	projectedBenefits := make([]float64, input.YearsToProject)
	notes := ""

	if input.YearsToProject <= 0 {
		return models.SocialSecurityCOLAResult{ProjectedBenefit: []float64{}, Notes: "YearsToProject must be positive."}
	}

	for i := 0; i < input.YearsToProject; i++ {
		ageInYear := input.ProjectionStartAge + i
		colaRateToApply := 0.0

		// Social Security COLA applies if the benefit has started (ageInYear >= BenefitStartAge)
		// and the COLA rate is positive.
		if ageInYear >= input.BenefitStartAge && input.AssumedCPIWRate > 0 {
			colaRateToApply = input.AssumedCPIWRate
			if i == 0 && notes == "" { // Add note only once
				notes += "Social Security: Full COLA applied each year after benefits commence. "
			}
		} else if i == 0 && ageInYear < input.BenefitStartAge && notes == "" {
			notes += "Social Security: Benefits not yet started. No COLA applied. "
		}

		// Only apply COLA if the benefit has actually started
		currentBenefitValue := 0.0
		if ageInYear < input.BenefitStartAge {
			currentBenefitValue = 0 // No benefit yet
		} else if ageInYear == input.BenefitStartAge && i == 0 {
			// First year of benefits, use initial amount (COLA will apply next year)
			currentBenefitValue = input.InitialAnnualBenefit
		} else if ageInYear == input.BenefitStartAge && i > 0 && projectedBenefits[i-1] == 0 {
			// Benefits starting this year, but not the first year of projection loop
			currentBenefitValue = input.InitialAnnualBenefit
		} else if projectedBenefits[i-1] > 0 { // If benefits were paid last year
			currentBenefitValue = projectedBenefits[i-1] * (1 + colaRateToApply)
		} else if ageInYear > input.BenefitStartAge && projectedBenefits[i-1] == 0 && i > 0 {
			// This case handles if benefits start *during* the projection period,
			// and it's not the first year of the loop.
			// The COLA should apply to the initial benefit amount.
			currentBenefitValue = input.InitialAnnualBenefit * (1 + colaRateToApply)
		}

		if currentBenefitValue > 0 {
			projectedBenefits[i] = math.Round(currentBenefitValue*100) / 100
		} else {
			projectedBenefits[i] = 0
		}
	}

	return models.SocialSecurityCOLAResult{
		ProjectedBenefit: projectedBenefits,
		Notes:            strings.TrimSpace(notes),
	}
}
