package calculation

import (
	"ferex/backend/models"
	"fmt"
	"math"
)

// CalculateCSRS computes the CSRS or CSRS Offset annuity based on user input.
func CalculateCSRS(input models.CSRSCalculationInput) models.CSRSCalculationResult {
	var notes string

	// Add sick leave to years of service (1 month = 1/12 year)
	serviceYearsWithSick := input.YearsOfService + float64(input.UnusedSickLeaveMonths)/12.0

	// Tiered multipliers
	first5 := math.Min(serviceYearsWithSick, 5)
	next5 := math.Min(math.Max(serviceYearsWithSick-5, 0), 5)
	over10 := math.Max(serviceYearsWithSick-10, 0)

	baseAnnuity := input.High3Salary * (first5*0.015 + next5*0.0175 + over10*0.02)

	// 80% High-3 maximum rule
	maxPension := input.High3Salary * 0.8
	if baseAnnuity > maxPension {
		baseAnnuity = maxPension
		notes += "80% High-3 maximum applied.\n"
	}

	// Part-time proration
	proratedPension := baseAnnuity
	prorationApplied := false
	if input.IsPartTime && input.PartTimeProrationFactor > 0 && input.PartTimeProrationFactor < 1.0 {
		proratedPension = baseAnnuity * input.PartTimeProrationFactor
		prorationApplied = true
		notes += fmt.Sprintf("Part-time proration factor applied: %.2f\n", input.PartTimeProrationFactor)
	}

	// Survivor benefit reduction (placeholder: adjust logic as needed)
	survivorReduction := 0.0
	if input.SurvivorBenefitElection == "max" {
		survivorReduction = proratedPension * 0.10
		notes += "Maximum survivor benefit reduction (10%) applied.\n"
	} else if input.SurvivorBenefitElection == "partial" {
		survivorReduction = proratedPension * 0.05
		notes += "Partial survivor benefit reduction (5%) applied.\n"
	}

	finalPension := proratedPension - survivorReduction

	// CSRS Offset reduction (if applicable)
	offsetReduction := 0.0
	if input.IsCSRSOffset && input.AgeAtRetirement >= 62 {
		// OPM reduces by the lesser of:
		// 1. SS with vs without Offset earnings
		// 2. (SS with Offset) * (Years of Offset / 40)
		option1 := input.SSAt62WithOffset - input.SSAt62WithoutOffset
		option2 := input.SSAt62WithOffset * (input.YearsOfOffsetService / 40.0)
		offsetReduction = math.Min(option1, option2)
		if offsetReduction < 0 {
			offsetReduction = 0
		}
		notes += fmt.Sprintf("CSRS Offset reduction applied: $%.2f\n", offsetReduction)
		finalPension -= offsetReduction
		if finalPension < 0 {
			finalPension = 0
		}
	}

	return models.CSRSCalculationResult{
		AnnualPension:            finalPension,
		MonthlyPension:           finalPension / 12.0,
		EarlyRetirementReduction: 0, // Placeholder: Add logic if needed
		SickLeaveServiceCredit:   float64(input.UnusedSickLeaveMonths) / 12.0,
		ProrationApplied:         prorationApplied,
		ProratedPension:          proratedPension,
		SurvivorBenefitReduction: survivorReduction,
		OffsetReduction:          offsetReduction,
		Notes:                    notes,
	}
}
