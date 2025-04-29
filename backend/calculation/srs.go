package calculation

import (
	"ferex/backend/models"
	"fmt"
	"math"
)

// Earnings test threshold for 2025; adjust as needed for future years
const earningsTestLimit2025 = 23400.0

// CalculateSRS computes the FERS Annuity Supplement (SRS) based on user input.
func CalculateSRS(input models.SRSCalculationInput) models.SRSCalculationResult {
	var notes string
	isEligible := false

	// Eligibility: Must be immediate, unreduced annuity (not MRA+10, deferred, or disability), and retire before 62
	if input.IsImmediateUnreducedAnnuity && input.RetirementAge < 62 {
		isEligible = true
	} else {
		notes += "Not eligible for SRS (must be immediate, unreduced FERS annuity and retire before age 62).\n"
	}

	// Only proceed with calculation if eligible
	if !isEligible {
		return models.SRSCalculationResult{
			AnnualSRSAmount:       0,
			MonthlySRSAmount:      0,
			EarningsTestReduction: 0,
			IsEligible:            false,
			Notes:                 notes,
		}
	}

	// SRS formula: (SS at 62 / 40) * Years of FERS service (rounded up)
	yearsOfFERSService := math.Ceil(input.YearsOfFERSService)
	srsAmount := (input.EstimatedSocialSecurityAt62 / 40.0) * yearsOfFERSService

	// Earnings test: $1 reduction for every $2 earned above the limit
	// Use 2025 limit for now; in production, this should be dynamic
	earningsTestLimit := earningsTestLimit2025
	overLimit := input.ProjectedEarnedIncome - earningsTestLimit
	earningsReduction := 0.0
	if overLimit > 0 {
		earningsReduction = overLimit / 2.0
		if earningsReduction > srsAmount {
			earningsReduction = srsAmount
		}
		notes += fmt.Sprintf("Earnings test reduction applied: $%.2f\n", earningsReduction)
	}

	finalSRS := srsAmount - earningsReduction
	if finalSRS < 0 {
		finalSRS = 0
	}

	return models.SRSCalculationResult{
		AnnualSRSAmount:       finalSRS,
		MonthlySRSAmount:      finalSRS / 12.0,
		EarningsTestReduction: earningsReduction,
		IsEligible:            true,
		Notes:                 notes,
	}
}
