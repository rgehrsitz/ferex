package calculation

import (
	"ferex/backend/models"
	"fmt"
)

// CalculateFERSPension computes the FERS annuity based on user input.
func CalculateFERSPension(input models.FERSCalculationInput) models.FERSCalculationResult {
	var notes string
	serviceYears := input.YearsOfService

	// Add sick leave (convert months to years)
	sickLeaveYears := float64(input.UnusedSickLeaveMonths) / 12.0
	serviceYearsWithSick := serviceYears + sickLeaveYears

	// Determine multiplier
	multiplier := 0.01
	if input.IsAge62With20Years {
		multiplier = 0.011
	}

	// Calculate base annuity
	baseAnnuity := input.High3Salary * serviceYearsWithSick * multiplier
	fmt.Printf("DEBUG: baseAnnuity=%.2f, IsPartTime=%v, PartTimeProrationFactor=%.2f\n", baseAnnuity, input.IsPartTime, input.PartTimeProrationFactor)

	// Early retirement reduction (FERS MRA+10: 5% per year under 62)
	earlyReduction := 0.0
	if input.AgeAtRetirement < 62 {
		yearsUnder := 62 - input.AgeAtRetirement
		earlyReduction = baseAnnuity * 0.05 * float64(yearsUnder)
		baseAnnuity -= earlyReduction
		if earlyReduction > 0 {
			notes += fmt.Sprintf("Early retirement reduction applied: %.2f\n", earlyReduction)
		}
	}

	// Part-time proration
	proratedPension := baseAnnuity
	prorationApplied := false
	if input.IsPartTime && input.PartTimeProrationFactor > 0 && input.PartTimeProrationFactor < 1.0 {
		proratedPension = baseAnnuity * input.PartTimeProrationFactor
		prorationApplied = true
		notes += fmt.Sprintf("Part-time proration factor applied: %.2f\n", input.PartTimeProrationFactor)
	}
	fmt.Printf("DEBUG: after proration, proratedPension=%.2f\n", proratedPension)

	// Survivor benefit reduction (if any)
	survivorReduction := 0.0
	switch input.SurvivorBenefitElection {
	case "max":
		survivorReduction = proratedPension * 0.10 // 10% reduction for max survivor
		proratedPension -= survivorReduction
		notes += "Maximum survivor benefit elected (10% reduction).\n"
	case "partial":
		survivorReduction = proratedPension * 0.05 // 5% reduction for partial survivor
		proratedPension -= survivorReduction
		notes += "Partial survivor benefit elected (5% reduction).\n"
	default:
		// No reduction for "none", empty, or any other value
	}

	return models.FERSCalculationResult{
		AnnualPension:            proratedPension,
		MonthlyPension:           proratedPension / 12.0,
		EarlyRetirementReduction: earlyReduction,
		SickLeaveServiceCredit:   sickLeaveYears,
		ProrationApplied:         prorationApplied,
		ProratedPension:          proratedPension,
		SurvivorBenefitReduction: survivorReduction,
		Notes:                    notes,
	}
}
