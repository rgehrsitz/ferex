package calculation

import (
	"ferex/backend/models"
	"fmt"
	"math"
)

const csrsSickLeaveHoursPerYear = 2087.0 // OPM standard

// CalculateCSRS performs the main CSRS/CSRS Offset annuity calculation.
func CalculateCSRS(retirementAge float64, yearsOfService float64, input models.CSRSCalculationInput) models.CSRSCalculationResult {
	var result models.CSRSCalculationResult

	// 1. Calculate Service Credit including Sick Leave
	sickLeaveYears := 0.0
	if input.UnusedSickLeaveHours > 0 {
		sickLeaveYears = input.UnusedSickLeaveHours / csrsSickLeaveHoursPerYear // Use float64 directly from input
	}
	totalServiceYears := yearsOfService + sickLeaveYears
	result.SickLeaveServiceCredit = sickLeaveYears // Store converted years
	result.TotalServiceYears = totalServiceYears   // Store total years used in calculation
	// --- IRS Simplified Method Exclusion ---
	result.IRSSimplifiedMethodExclusion = 0.0
	if input.EmployeeContributions > 0 && retirementAge > 0 {
		hasSurvivor := input.SurvivorBenefitType != "None"
		annualExclusion := CalculateIRSSimplifiedMethodExclusion(input.EmployeeContributions, int(retirementAge), hasSurvivor)
		result.IRSSimplifiedMethodExclusion = annualExclusion / 12.0 // Convert to monthly
		if result.IRSSimplifiedMethodExclusion > 0 {
			result.Notes += fmt.Sprintf("IRS Simplified Method monthly exclusion: $%.2f\n", result.IRSSimplifiedMethodExclusion)
		}
	}
	// 2. Calculate Monthly Gross Annuity using the CSRS formula including sick leave
	// Calculate annual first, then convert to monthly for primary storage
	annualGrossAnnuity := calculateCSRSFormula(input.High3Salary, totalServiceYears)
	result.MonthlyGrossAnnuity = annualGrossAnnuity / 12.0 // Store monthly as primary value

	// 3. Apply Part-Time Proration if applicable
	if input.IsPartTime {
		result.MonthlyGrossAnnuityBeforeProration = result.MonthlyGrossAnnuity // Store pre-proration value
		result.MonthlyGrossAnnuity *= input.PartTimeProrationFactor
		result.IsProrated = true
	}
	// 4. Apply CSRS Offset Reduction if applicable
	if input.IsCSRSOffset {
		// Placeholder for Offset Reduction logic
		// monthlyOffsetReduction := calculateCSRSOffsetReduction(input, totalServiceYears, retirementAge) / 12.0
		// result.MonthlyOffsetReduction = monthlyOffsetReduction // Store the monthly reduction amount
		// result.MonthlyNetAnnuity = result.MonthlyGrossAnnuity - monthlyOffsetReduction // Apply reduction
		result.MonthlyNetAnnuity = result.MonthlyGrossAnnuity // TEMP: No reduction yet
	} else {
		result.MonthlyNetAnnuity = result.MonthlyGrossAnnuity
	}
	// 5. Calculate Survivor Benefit Reduction
	monthlySurvivorReduction := 0.0
	if input.SurvivorBenefitType != "None" {
		// Calculate cost based on the annuity *before* survivor reduction (Net Annuity)
		if input.SurvivorBenefitType == "Full" {
			monthlySurvivorReduction = result.MonthlyNetAnnuity * 0.10 // Approx. 10% for full survivor benefit
		} else if input.SurvivorBenefitType == "Partial" {
			monthlySurvivorReduction = result.MonthlyNetAnnuity * 0.05 // Approx. 5% for partial
		}
		// TODO: Add logic for specific spousal vs insurable interest rates
	}
	result.MonthlySurvivorBenefitReduction = monthlySurvivorReduction

	// 6. Apply Early Retirement Reduction (if applicable - simplified logic)
	// This should ideally be applied after survivor cost, before final calculation
	// Note: Current implementation skipped as logic needs review.

	// 7. Set final monthly annuity and derive annual values
	result.MonthlyFinalAnnuity = result.MonthlyNetAnnuity - monthlySurvivorReduction

	// Populate legacy annual fields for backward compatibility
	result.GrossAnnuity = result.MonthlyGrossAnnuity * 12.0
	result.GrossAnnuityBeforeProration = result.MonthlyGrossAnnuityBeforeProration * 12.0
	result.OffsetReduction = result.MonthlyOffsetReduction * 12.0
	result.NetAnnuity = result.MonthlyNetAnnuity * 12.0
	result.SurvivorBenefitReduction = result.MonthlySurvivorBenefitReduction * 12.0
	result.FinalAnnuity = result.MonthlyFinalAnnuity * 12.0

	// --- Set new fields ---
	// RetirementType: Immediate if age >= 55 and service >= 30, Deferred if age < 55 and service >= 5, else Other
	if retirementAge >= 55 && yearsOfService >= 30 {
		result.RetirementType = "Immediate"
	} else if retirementAge < 55 && yearsOfService >= 5 {
		result.RetirementType = "Deferred"
	} else {
		result.RetirementType = "Other"
	}
	if input.IsPartTime {
		result.ProrationFactor = input.PartTimeProrationFactor
	} else {
		result.ProrationFactor = 1.0
	}

	return result
}

// CalculateCSRSComponentAnnuity calculates the basic CSRS annuity component for a given High-3 and CSRS service years.
// This is typically used for the CSRS portion of a FERS transferee's annuity.
// It does not include sick leave, as sick leave for FERS transferees is credited under FERS.
func CalculateCSRSComponentAnnuity(high3Salary float64, csrsServiceYears float64) float64 {
	return calculateCSRSFormula(high3Salary, csrsServiceYears)
}

// calculateCSRSFormula applies the core CSRS tiered percentage formula.
func calculateCSRSFormula(high3 float64, years float64) float64 {
	if years <= 0 || high3 <= 0 {
		return 0.0
	}

	annuity := 0.0
	// 1.5% for the first 5 years
	years1 := math.Min(years, 5.0)
	annuity += 0.015 * high3 * years1

	// 1.75% for the next 5 years (years 6-10)
	if years > 5.0 {
		years2 := math.Min(years-5.0, 5.0)
		annuity += 0.0175 * high3 * years2
	}

	// 2.0% for years over 10
	if years > 10.0 {
		years3 := years - 10.0
		annuity += 0.02 * high3 * years3
	}

	return annuity
}
