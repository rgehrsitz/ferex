package calculation

import (
	"ferex/backend/models"
	"fmt"
	"math"
	"strings"
)

const fersSickLeaveHoursPerYear = 2087.0 // OPM standard for FERS sick leave conversion (2087 hours/year)

// CalculateFERS performs the main FERS annuity calculation.
// It takes pre-validated and derived inputs.
func CalculateFERS(
	input models.FERSCalculationInput, // Note: input.High3Salary should be pre-adjusted by prorationFactor if applicable
	ageAtRetirementYears int,
	ageAtRetirementMonths int, // For future precise calculations if needed, currently years are primary
	totalCreditableServiceYearsBeforeSickLeave float64,
	mraYears int,
	// mraMonths int, // For future precise MRA calculations
	retirementType string, // e.g., "ImmediateMRA30", "MRA+10Voluntary", "Deferred"
	expectedSSBenefitAt62 float64, // Monthly SS benefit estimate for SRS
	yearsActuallyWorkedUnderFERS float64, // For SRS, distinct from total creditable service
	prorationAppliedByCaller bool,
	prorationFactorUsedByCaller float64,
	totalCombinedActualServiceYears float64, // Added for SRS eligibility check
) models.FERSCalculationResult {
	var result models.FERSCalculationResult
	var notes []string
	// Validate for invalid retirement scenarios
	if retirementType == "Invalid" || ageAtRetirementYears < 30 {
		// For invalid retirement types or unrealistically young retirement ages,
		// return minimal/zero benefits with appropriate warnings
		result.MonthlyPension = 0.0
		result.MonthlyBasicAnnuity = 0.0
		result.AnnualPension = 0.0
		result.TotalServiceYears = totalCreditableServiceYearsBeforeSickLeave
		result.SickLeaveServiceCredit = 0.0
		result.RetirementType = retirementType
		result.IsEligibleForSupplement = false
		result.Notes = "Invalid retirement scenario - no pension benefit available"

		if retirementType == "Invalid" {
			result.Notes += "; Retirement type marked as Invalid"
		}
		if ageAtRetirementYears < 30 {
			result.Notes += "; Retirement age is unrealistically young for federal employment"
		}

		return result
	}
	// 1. Calculate Service Credit including Sick Leave
	sickLeaveYearsCredit := 0.0
	if input.UnusedSickLeaveHours > 0 {
		sickLeaveYearsCredit = float64(input.UnusedSickLeaveHours) / fersSickLeaveHoursPerYear
	}
	result.SickLeaveServiceCredit = sickLeaveYearsCredit
	totalServiceYearsWithSickLeave := totalCreditableServiceYearsBeforeSickLeave + sickLeaveYearsCredit
	fersServiceYearsWithSickLeave := yearsActuallyWorkedUnderFERS + sickLeaveYearsCredit
	result.TotalServiceYears = totalServiceYearsWithSickLeave
	notes = append(notes, fmt.Sprintf("Total creditable service: %.2f years (including %.2f years from sick leave).", totalServiceYearsWithSickLeave, sickLeaveYearsCredit))
	notes = append(notes, fmt.Sprintf("FERS service years for annuity calculation: %.2f years (including %.2f years from sick leave).", fersServiceYearsWithSickLeave, sickLeaveYearsCredit))
	// --- IRS Simplified Method Exclusion ---
	result.IRSSimplifiedMethodExclusion = 0.0
	if input.EmployeeContributions > 0 && ageAtRetirementYears > 0 {
		hasSurvivor := input.SurvivorBenefitElection != "None"
		annualExclusion := CalculateIRSSimplifiedMethodExclusion(input.EmployeeContributions, ageAtRetirementYears, hasSurvivor)
		result.IRSSimplifiedMethodExclusion = annualExclusion / 12.0 // Convert to monthly
		if result.IRSSimplifiedMethodExclusion > 0 {
			notes = append(notes, fmt.Sprintf("IRS Simplified Method monthly exclusion: $%.2f", result.IRSSimplifiedMethodExclusion))
		}
	}

	// 2. Determine Annuity Multiplier (1% or 1.1%)
	multiplier := 0.01
	// Age 62 or older with at least 20 years of service gets 1.1%
	if ageAtRetirementYears >= 62 && totalServiceYearsWithSickLeave >= 20 {
		multiplier = 0.011
		notes = append(notes, "Annuity multiplier: 1.1%.")
	} else {
		notes = append(notes, "Annuity multiplier: 1.0%.")
	}
	// 3. Calculate Monthly Basic Annuity (before any proration or reductions)
	// Calculate annual first, then convert to monthly for primary storage
	// Use FERS service years for the annuity calculation (important for CSRS transferees)
	grossAnnualAnnuity := input.High3Salary * fersServiceYearsWithSickLeave * multiplier
	monthlyBasicAnnuity := grossAnnualAnnuity / 12.0
	result.MonthlyBasicAnnuity = monthlyBasicAnnuity // Store monthly as primary value
	notes = append(notes, fmt.Sprintf("Monthly basic annuity: $%.2f (Annual: $%.2f).", monthlyBasicAnnuity, grossAnnualAnnuity))
	// 4. Handle Part-Time Proration (Note: High3Salary in 'input' is expected to be already prorated by app.go if applicable)
	// This section now primarily records the proration details for the result.
	monthlyAnnuityAfterProration := monthlyBasicAnnuity // Since input.High3Salary was pre-adjusted, this is effectively the prorated annuity.
	if prorationAppliedByCaller && prorationFactorUsedByCaller > 0 && prorationFactorUsedByCaller < 1.0 {
		result.ProrationApplied = true
		// result.MonthlyProratedPension should reflect the annuity based on the adjusted High3.
		// Since monthlyBasicAnnuity is calculated using the (potentially) adjusted High3, it's already the prorated value.
		result.MonthlyProratedPension = monthlyBasicAnnuity
		notes = append(notes, fmt.Sprintf("Part-time proration factor (%.3f) applied to High-3 salary. Monthly annuity after proration: $%.2f.", prorationFactorUsedByCaller, monthlyAnnuityAfterProration))
	} else {
		result.ProrationApplied = false
		result.MonthlyProratedPension = monthlyBasicAnnuity // If not part-time, this is the same as basic
	}
	// 5. Apply Early Retirement Reduction (if applicable)
	monthlyAnnuityAfterEarlyReduction := monthlyAnnuityAfterProration
	if retirementType == "MRA+10Voluntary" || retirementType == "MRA+20Voluntary" { // Assuming MRA+20 is similar rule if under 62
		if ageAtRetirementYears < 62 {
			yearsUnder62 := 62 - ageAtRetirementYears
			// OPM: "5 percent a year for each year you are under age 62 (this is 5/12 of 1 percent for each full month)."
			// For simplicity with integer years from UI for now:
			monthlyEarlyRetirementReduction := monthlyAnnuityAfterProration * 0.05 * float64(yearsUnder62)
			monthlyAnnuityAfterEarlyReduction -= monthlyEarlyRetirementReduction
			result.MonthlyEarlyRetirementReduction = monthlyEarlyRetirementReduction
			notes = append(notes, fmt.Sprintf("Early retirement reduction (%.0f%% for %d years under 62) applied: -$%.2f monthly.", 0.05*float64(yearsUnder62)*100, yearsUnder62, monthlyEarlyRetirementReduction))
		}
	} else if retirementType == "MRA+10Postponed" {
		notes = append(notes, "MRA+10 retirement with postponed annuity. No age-based reduction at postponed start date.")
	}
	// 6. Apply Survivor Benefit Cost (if applicable)
	monthlyAnnuityAfterSurvivorReduction := monthlyAnnuityAfterEarlyReduction
	survivorBenefitCostPercentage := 0.0
	switch input.SurvivorBenefitElection {
	case "Full50Percent": // Corresponds to 10% cost to retiree
		survivorBenefitCostPercentage = 0.10
		notes = append(notes, "Full survivor benefit (50%% to survivor) elected.")
	case "Partial25Percent": // Corresponds to 5% cost to retiree
		survivorBenefitCostPercentage = 0.05
		notes = append(notes, "Partial survivor benefit (25%% to survivor) elected.")
	case "None":
		survivorBenefitCostPercentage = 0.0
		notes = append(notes, "No survivor benefit elected.")
	default: // Includes "InsurableInterest" and any other unexpected values
		survivorBenefitCostPercentage = 0.0
		if strings.Contains(input.SurvivorBenefitElection, "InsurableInterest") {
			notes = append(notes, "Insurable Interest survivor benefit selected. Cost calculation not yet implemented in this module; assuming 0 cost for now.")
		} else if input.SurvivorBenefitElection != "" {
			notes = append(notes, fmt.Sprintf("Unknown survivor benefit election '%s'; assuming no cost.", input.SurvivorBenefitElection))
		}
	}
	if survivorBenefitCostPercentage > 0 {
		// Survivor benefit cost is applied to the unreduced basic annuity per OPM regulations.
		monthlySurvivorReduction := monthlyBasicAnnuity * survivorBenefitCostPercentage
		monthlyAnnuityAfterSurvivorReduction -= monthlySurvivorReduction
		result.MonthlySurvivorBenefitReduction = monthlySurvivorReduction
		notes = append(notes, fmt.Sprintf("Survivor benefit cost (%.0f%%) applied: -$%.2f monthly.", survivorBenefitCostPercentage*100, monthlySurvivorReduction))
	}
	result.MonthlyPension = math.Max(0, monthlyAnnuityAfterSurvivorReduction) // Ensure pension isn't negative
	result.AnnualPension = result.MonthlyPension * 12.0                       // Derive annual from monthly
	notes = append(notes, fmt.Sprintf("Final monthly pension: $%.2f (Annual: $%.2f).", result.MonthlyPension, result.AnnualPension))

	// Populate legacy annual fields for backward compatibility
	result.BasicAnnuity = result.MonthlyBasicAnnuity * 12.0
	result.EarlyRetirementReduction = result.MonthlyEarlyRetirementReduction * 12.0
	result.ProratedPension = result.MonthlyProratedPension * 12.0
	result.SurvivorBenefitReduction = result.MonthlySurvivorBenefitReduction * 12.0
	// 7. Calculate FERS Supplement (SRS) - if eligible
	effectiveTotalServiceForRetirementType := totalCombinedActualServiceYears + sickLeaveYearsCredit
	result.IsEligibleForSupplement = checkFersSupplementEligibility(ageAtRetirementYears, mraYears, retirementType, effectiveTotalServiceForRetirementType)
	if result.IsEligibleForSupplement {
		// SRS formula: (SS at 62 / 40) * Years of FERS service (actual FERS service, not necessarily total creditable)
		// OPM uses full years of FERS civilian service. For simplicity, we use yearsActuallyWorkedUnderFERS.
		yearsFactor := yearsActuallyWorkedUnderFERS / 40.0
		if yearsFactor > 1.0 {
			yearsFactor = 1.0 // Cap at 40 years
		}
		// This is a simplified calculation. Full SRS includes earnings test.
		result.FersSupplement = expectedSSBenefitAt62 * yearsFactor // This is monthly, as SS benefit is usually monthly
		notes = append(notes, fmt.Sprintf("Eligible for FERS Annuity Supplement. Estimated monthly supplement: $%.2f.", result.FersSupplement))
	} else {
		notes = append(notes, "Not eligible for FERS Annuity Supplement or supplement ends at age 62.")
	}

	// --- Set new fields ---
	result.RetirementType = retirementType
	result.ProrationFactor = prorationFactorUsedByCaller

	// 6. Notes
	result.Notes = strings.Join(notes, "\n")
	return result
}

// CalculateFERSWithAuditTrail performs the main FERS annuity calculation with detailed audit trail.
// This enhanced version provides step-by-step calculation breakdown for transparency and validation.
func CalculateFERSWithAuditTrail(
	input models.FERSCalculationInput,
	ageAtRetirementYears int,
	ageAtRetirementMonths int,
	totalCreditableServiceYearsBeforeSickLeave float64,
	mraYears int,
	retirementType string,
	expectedSSBenefitAt62 float64,
	yearsActuallyWorkedUnderFERS float64,
	prorationAppliedByCaller bool,
	prorationFactorUsedByCaller float64,
	totalCombinedActualServiceYears float64,
) models.FERSCalculationResult {

	// Initialize audit trail
	auditTrail := &models.CalculationAuditTrail{
		CalculationType: "FERS Pension Calculation",
		InputSummary: fmt.Sprintf("Employee retiring at age %d with %.2f years of service, High-3 salary of $%.2f",
			ageAtRetirementYears, totalCombinedActualServiceYears, input.High3Salary),
		Steps:    []models.AuditStep{},
		Warnings: []string{},
		OPMReferences: []string{
			"OPM FERS Handbook Chapter 50 - Computation",
			"5 USC 8415 - Computation of basic annuity",
		},
	}

	stepNum := 1

	// Step 1: Calculate Service Credit including Sick Leave
	sickLeaveYearsCredit := 0.0
	if input.UnusedSickLeaveHours > 0 {
		sickLeaveYearsCredit = float64(input.UnusedSickLeaveHours) / fersSickLeaveHoursPerYear
	}

	auditTrail.Steps = append(auditTrail.Steps, models.AuditStep{
		StepNumber:  stepNum,
		StepName:    "Sick Leave Service Credit",
		Description: "Convert unused sick leave hours to service credit years",
		Formula:     "Sick Leave Years = Unused Hours ÷ 2087 hours/year",
		Inputs: map[string]interface{}{
			"unusedSickLeaveHours": input.UnusedSickLeaveHours,
			"hoursPerYear":         fersSickLeaveHoursPerYear,
		},
		Calculation: fmt.Sprintf("%.0f ÷ %.0f = %.4f", float64(input.UnusedSickLeaveHours), fersSickLeaveHoursPerYear, sickLeaveYearsCredit),
		Result:      sickLeaveYearsCredit,
		Notes:       "Per OPM: 2,087 hours = 1 year of service credit for FERS",
	})
	stepNum++

	totalServiceYearsWithSickLeave := totalCreditableServiceYearsBeforeSickLeave + sickLeaveYearsCredit

	auditTrail.Steps = append(auditTrail.Steps, models.AuditStep{
		StepNumber:  stepNum,
		StepName:    "Total Creditable Service",
		Description: "Add sick leave credit to base service years",
		Formula:     "Total Service = Base Service + Sick Leave Credit",
		Inputs: map[string]interface{}{
			"baseServiceYears": totalCreditableServiceYearsBeforeSickLeave,
			"sickLeaveCredit":  sickLeaveYearsCredit,
		},
		Calculation: fmt.Sprintf("%.4f + %.4f = %.4f", totalCreditableServiceYearsBeforeSickLeave, sickLeaveYearsCredit, totalServiceYearsWithSickLeave),
		Result:      totalServiceYearsWithSickLeave,
	})
	stepNum++

	// Step 2: Determine Annuity Multiplier
	multiplier := 0.01
	multiplierReason := "Standard FERS multiplier (1.0%)"

	if ageAtRetirementYears >= 62 && totalServiceYearsWithSickLeave >= 20 {
		multiplier = 0.011
		multiplierReason = "Enhanced multiplier for age 62+ with 20+ years service (1.1%)"
	}

	auditTrail.Steps = append(auditTrail.Steps, models.AuditStep{
		StepNumber:  stepNum,
		StepName:    "Determine Annuity Multiplier",
		Description: "Select appropriate FERS multiplier based on age and service",
		Formula:     "If (Age ≥ 62 AND Service ≥ 20) then 1.1%, else 1.0%",
		Inputs: map[string]interface{}{
			"retirementAge": ageAtRetirementYears,
			"totalService":  totalServiceYearsWithSickLeave,
		},
		Calculation: fmt.Sprintf("Age %d ≥ 62? %v, Service %.2f ≥ 20? %v → %.1f%%",
			ageAtRetirementYears, ageAtRetirementYears >= 62,
			totalServiceYearsWithSickLeave, totalServiceYearsWithSickLeave >= 20,
			multiplier*100),
		Result: multiplier,
		Notes:  multiplierReason,
	})
	stepNum++

	// Step 3: Calculate Basic Annual Annuity
	grossAnnualAnnuity := input.High3Salary * totalServiceYearsWithSickLeave * multiplier
	monthlyBasicAnnuity := grossAnnualAnnuity / 12.0

	auditTrail.Steps = append(auditTrail.Steps, models.AuditStep{
		StepNumber:  stepNum,
		StepName:    "Basic Annuity Calculation",
		Description: "Calculate gross annual pension before any reductions",
		Formula:     "Annual Annuity = High-3 Salary × Total Service Years × Multiplier",
		Inputs: map[string]interface{}{
			"high3Salary":  input.High3Salary,
			"serviceYears": totalServiceYearsWithSickLeave,
			"multiplier":   multiplier,
		},
		Calculation: fmt.Sprintf("$%.2f × %.4f × %.3f = $%.2f annual ($%.2f monthly)",
			input.High3Salary, totalServiceYearsWithSickLeave, multiplier, grossAnnualAnnuity, monthlyBasicAnnuity),
		Result: monthlyBasicAnnuity,
		Notes:  "Monthly amount calculated as annual ÷ 12",
	})
	stepNum++

	// Continue with existing calculation logic but add audit steps...
	result := CalculateFERS(input, ageAtRetirementYears, ageAtRetirementMonths,
		totalCreditableServiceYearsBeforeSickLeave, mraYears, retirementType,
		expectedSSBenefitAt62, yearsActuallyWorkedUnderFERS, prorationAppliedByCaller,
		prorationFactorUsedByCaller, totalCombinedActualServiceYears)

	// Add survivor benefit reduction step if applicable
	if input.SurvivorBenefitElection != "None" {
		survivorReductionPct := 0.0
		switch input.SurvivorBenefitElection {
		case "Full50Percent":
			survivorReductionPct = 0.10
		case "Partial25Percent":
			survivorReductionPct = 0.05
		}

		auditTrail.Steps = append(auditTrail.Steps, models.AuditStep{
			StepNumber:  stepNum,
			StepName:    "Survivor Benefit Cost",
			Description: "Apply reduction for elected survivor benefit",
			Formula:     fmt.Sprintf("Reduction = Basic Annuity × %.0f%%", survivorReductionPct*100),
			Inputs: map[string]interface{}{
				"basicAnnuity":        monthlyBasicAnnuity,
				"reductionPercentage": survivorReductionPct,
				"electedBenefit":      input.SurvivorBenefitElection,
			},
			Calculation: fmt.Sprintf("$%.2f × %.1f%% = $%.2f monthly reduction",
				monthlyBasicAnnuity, survivorReductionPct*100, result.MonthlySurvivorBenefitReduction),
			Result: result.MonthlySurvivorBenefitReduction,
			Notes:  "Cost to provide survivor annuity benefit",
		})
		stepNum++
	}

	// Add FERS Supplement calculation step if eligible
	if result.IsEligibleForSupplement {
		auditTrail.Steps = append(auditTrail.Steps, models.AuditStep{
			StepNumber:  stepNum,
			StepName:    "FERS Supplement Calculation",
			Description: "Calculate FERS Annuity Supplement (SRS)",
			Formula:     "Monthly SRS = (SS Benefit at 62 ÷ 40) × FERS Service Years",
			Inputs: map[string]interface{}{
				"ssBenefitAt62":    expectedSSBenefitAt62,
				"fersServiceYears": yearsActuallyWorkedUnderFERS,
			},
			Calculation: fmt.Sprintf("($%.2f ÷ 40) × %.2f = $%.2f monthly",
				expectedSSBenefitAt62, yearsActuallyWorkedUnderFERS, result.FersSupplement),
			Result: result.FersSupplement,
			Notes:  "Supplement paid until age 62, subject to earnings test",
		})
		stepNum++
	}

	// Add IRS Simplified Method calculation if applicable
	if result.IRSSimplifiedMethodExclusion > 0 {
		auditTrail.Steps = append(auditTrail.Steps, models.AuditStep{
			StepNumber:  stepNum,
			StepName:    "IRS Simplified Method Exclusion",
			Description: "Calculate monthly tax-free portion of pension",
			Formula:     "Monthly Exclusion = (Total Contributions ÷ Life Expectancy Factor) ÷ 12",
			Inputs: map[string]interface{}{
				"employeeContributions": input.EmployeeContributions,
				"retirementAge":         ageAtRetirementYears,
			},
			Calculation: fmt.Sprintf("Tax-free portion: $%.2f monthly", result.IRSSimplifiedMethodExclusion),
			Result:      result.IRSSimplifiedMethodExclusion,
			Notes:       "Based on IRS Publication 721 tables",
		})
	}

	// Set final result and add warnings if applicable
	auditTrail.FinalResult = result.MonthlyPension

	if result.MonthlyEarlyRetirementReduction > 0 {
		auditTrail.Warnings = append(auditTrail.Warnings,
			fmt.Sprintf("Early retirement reduction of $%.2f monthly applied - this is a permanent reduction",
				result.MonthlyEarlyRetirementReduction))
	}

	if !result.IsEligibleForSupplement && ageAtRetirementYears < 62 {
		auditTrail.Warnings = append(auditTrail.Warnings,
			"Not eligible for FERS Supplement - consider impact on retirement income until age 62")
	}

	// Attach audit trail to result
	result.AuditTrail = auditTrail

	return result
}

// checkFersSupplementEligibility determines if the retiree meets the criteria for the FERS supplement.
// Based on OPM rules: supplement is only eligible for immediate, unreduced annuities before age 62.
// Specifically: MRA+30, Age 60+20, and VERA/DSR scenarios are eligible.
// MRA+10 (reduced annuity) is NOT eligible for supplement.
func checkFersSupplementEligibility(ageAtRetirementYears int, mraYears int, retirementType string, effectiveTotalServiceForRetirementType float64) bool {
	// Supplement is only payable if under age 62
	if ageAtRetirementYears >= 62 {
		return false // Supplement ends at 62, Social Security becomes available
	}

	switch retirementType {
	case "ImmediateMRA30":
		// MRA+30: immediate, unreduced annuity - supplement eligible
		return ageAtRetirementYears >= mraYears && effectiveTotalServiceForRetirementType >= 30
	case "Immediate60_20", "Immediate6020", "Immediate6220":
		// Age 60+20: immediate, unreduced annuity - supplement eligible
		return ageAtRetirementYears >= 60 && effectiveTotalServiceForRetirementType >= 20
	case "MRA+10Voluntary":
		// MRA+10: reduced annuity - NOT supplement eligible per OPM rules
		return false
	case "VERA", "DSR", "LEOImmediate":
		// VERA/DSR and special category immediate retirements - supplement eligible
		return true
	// Other retirement types like Deferred, Disability, etc. are not supplement eligible
	default:
		return false
	}
}
