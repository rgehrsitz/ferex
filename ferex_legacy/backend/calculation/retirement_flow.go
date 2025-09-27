package calculation

import (
	"ferex/backend/models"
	"log"
	"strings"
	"time"
)

// PerformRetirementCalculation is the main orchestrator for a full retirement scenario calculation.
// It integrates all modules, including the Monte Carlo simulation if requested.

func PerformRetirementCalculation(input models.RetirementCalculationInput) models.RetirementCalculationResult {
	// Log the incoming input for debugging
	// if inputBytes, err := json.MarshalIndent(input, "", "  "); err == nil {
	// 	log.Printf("DEBUG: Received RetirementCalculationInput:\n%s\n", string(inputBytes))
	// } else {
	// 	log.Printf("DEBUG: Failed to marshal RetirementCalculationInput: %v", err)
	// }

	var result models.RetirementCalculationResult
	var pensionAmount float64

	// (1) Pension Calculations (FERS, CSRS, SRS, etc.)
	// --- Derive FERS arguments ---
	// --- Derive FERS arguments ---
	// Notes collected throughout the calculation process
	notes := make([]string, 0)
	ageAtRetirementYears, ageAtRetirementMonths, errAge := CalculateAge(input.FERSInput.DateOfBirth, input.FERSInput.PlannedRetirementDate)
	if errAge != nil {
		notes = append(notes, "Error calculating age at retirement: "+errAge.Error())
		ageAtRetirementYears, ageAtRetirementMonths = 0, 0
	}
	mraYears, errMRA := MinimumRetirementAge(input.FERSInput.DateOfBirth)
	if errMRA != nil {
		notes = append(notes, "Error calculating MRA: "+errMRA.Error())
		mraYears = 0
	}
	totalCreditableServiceYearsBeforeSickLeave, errSvc := CalculateServiceYears(input.FERSInput.ServicePeriods)
	if errSvc != nil {
		notes = append(notes, "Error calculating service years: "+errSvc.Error())
		totalCreditableServiceYearsBeforeSickLeave = 0
	}
	// --- Determine retirementType and prorationFactor ---
	retirementType := DetermineRetirementType(ageAtRetirementYears, totalCreditableServiceYearsBeforeSickLeave, mraYears, input.FERSInput.PlannedRetirementDate, input.FERSInput.DateOfBirth)
	prorationFactor := ProrationFactor(input.FERSInput.ServicePeriods)
	prorationAppliedByCaller := prorationFactor != 1.0
	expectedSSBenefitAt62 := 0.0 // TODO: Implement logic
	yearsForSRSCalc := totalCreditableServiceYearsBeforeSickLeave // For regular FERS employees, use total service years
	yearsOfService := totalCreditableServiceYearsBeforeSickLeave

	result.FERSResult = CalculateFERS(
		input.FERSInput,
		ageAtRetirementYears,
		ageAtRetirementMonths,
		totalCreditableServiceYearsBeforeSickLeave,
		mraYears,
		retirementType,
		expectedSSBenefitAt62,
		yearsForSRSCalc,
		prorationAppliedByCaller,
		prorationFactor,
		yearsOfService,
	)
	// --- Derive CSRS arguments ---
	csrsAgeAtRetirement, _, errCSRSAge := CalculateAge(input.CSRSInput.DateOfBirth, input.CSRSInput.PlannedRetirementDate)
	if errCSRSAge != nil {
		notes = append(notes, "Error calculating CSRS age at retirement: "+errCSRSAge.Error())
		csrsAgeAtRetirement = 0
	}

	csrsYearsOfService, errCSRSSvc := CalculateServiceYears(input.FERSInput.ServicePeriods) // Use FERS periods for now (adjust if CSRS periods are separate)
	if errCSRSSvc != nil {
		notes = append(notes, "Error calculating CSRS service years: "+errCSRSSvc.Error())
		csrsYearsOfService = 0
	}
	result.CSRSResult = CalculateCSRS(float64(csrsAgeAtRetirement), csrsYearsOfService, input.CSRSInput)

	// Log FERS and CSRS results after calculation
	// if fersBytes, err := json.MarshalIndent(result.FERSResult, "", "  "); err == nil {
	// 	log.Printf("DEBUG: FERSResult after calculation:\n%s\n", string(fersBytes))
	// }
	// if csrsBytes, err := json.MarshalIndent(result.CSRSResult, "", "  "); err == nil {
	// 	log.Printf("DEBUG: CSRSResult after calculation:\n%s\n", string(csrsBytes))
	// }

	// Determine the primary pension amount for overall summary
	switch input.CalculationSystem {
	case "FERS":
		pensionAmount = result.FERSResult.AnnualPension
	case "CSRS":
		pensionAmount = result.CSRSResult.FinalAnnuity
	default:
		// If no specific system, or 'Both', use the larger of the two calculated pensions
		if result.CSRSResult.FinalAnnuity > result.FERSResult.AnnualPension {
			pensionAmount = result.CSRSResult.FinalAnnuity
		} else {
			pensionAmount = result.FERSResult.AnnualPension
		}
		// Add a note directly to result.Notes instead of the notes slice
		result.Notes += "Pension amount selected based on higher value between FERS and CSRS calculations as system was not specified or was 'Both'. "
	}
	log.Printf("DEBUG: Selected pensionAmount for summary: %.2f", pensionAmount)

	// --- SRS Calculation (TODO) ---
	result.SRSResult = models.SRSCalculationResult{} // TODO: Implement SRS logic

	// (2) TSP Calculation
	result.TSPResult = CalculateTSP(input.TSPInput)
	// if tspBytes, err := json.MarshalIndent(result.TSPResult, "", "  "); err == nil {
	// 	log.Printf("DEBUG: TSPResult after calculation:\n%s\n", string(tspBytes))
	// }

	// (3) Social Security, COLA, Survivor, Health (moved before tax calculation)
	result.SocialSecurityResult = CalculateSocialSecurity(input.SocialSecurityInput)
	// if ssBytes, err := json.MarshalIndent(result.SocialSecurityResult, "", "  "); err == nil {
	// 	log.Printf("DEBUG: SocialSecurityResult after calculation:\n%s\n", string(ssBytes))
	// }
	const layout = "2006-01-02"       // Standard Go time parsing layout for YYYY-MM-DD
	const maxProjectionAgeYears = 100 // TODO: Make this configurable

	// Determine precise retirement start date and birth date
	var retirementStartDate, birthDate time.Time
	var err error

	plannedRetirementDateStr := ""
	if input.CalculationSystem == "FERS" && input.FERSInput.PlannedRetirementDate != "" {
		plannedRetirementDateStr = input.FERSInput.PlannedRetirementDate
	} else if input.CalculationSystem == "CSRS" && input.CSRSInput.PlannedRetirementDate != "" {
		plannedRetirementDateStr = input.CSRSInput.PlannedRetirementDate
	} else if input.FERSInput.PlannedRetirementDate != "" { // Fallback if system is not specified but FERS date exists
		plannedRetirementDateStr = input.FERSInput.PlannedRetirementDate
	} else if input.CSRSInput.PlannedRetirementDate != "" { // Fallback if system is not specified but CSRS date exists
		plannedRetirementDateStr = input.CSRSInput.PlannedRetirementDate
	}

	if plannedRetirementDateStr != "" {
		retirementStartDate, err = time.Parse(layout, plannedRetirementDateStr)
		if err != nil || retirementStartDate.IsZero() {
			log.Printf("ERROR: Invalid PlannedRetirementDate '%s': %v. Cannot proceed with monthly projection.", plannedRetirementDateStr, err)
			result.Notes = "Error: Invalid PlannedRetirementDate. Monthly projection cannot be performed."
			return result // Early exit
		}
	} else {
		log.Printf("ERROR: PlannedRetirementDate is missing. Cannot proceed with monthly projection.")
		result.Notes = "Error: PlannedRetirementDate is missing. Monthly projection cannot be performed."
		return result // Early exit
	}

	dobStr := ""
	if input.CalculationSystem == "FERS" && input.FERSInput.DateOfBirth != "" {
		dobStr = input.FERSInput.DateOfBirth
	} else if input.CalculationSystem == "CSRS" && input.CSRSInput.DateOfBirth != "" {
		dobStr = input.CSRSInput.DateOfBirth
	} else if input.FERSInput.DateOfBirth != "" { // Fallback
		dobStr = input.FERSInput.DateOfBirth
	} else if input.CSRSInput.DateOfBirth != "" { // Fallback
		dobStr = input.CSRSInput.DateOfBirth
	}

	if dobStr != "" {
		birthDate, err = time.Parse(layout, dobStr)
		if err != nil || birthDate.IsZero() {
			log.Printf("WARN: Invalid DateOfBirth '%s': %v. Age calculation might be impacted.", dobStr, err)
			// Not returning, as we might still proceed with ageAtRetirementYears/Months if they were pre-calculated correctly
		}
	} else {
		log.Printf("WARN: DateOfBirth is missing. Age calculation might be impacted.")
	}
	log.Printf("DEBUG: Retirement Start Date: %s, DOB: %s, Initial Age at Retirement: %d years, %d months", retirementStartDate.Format(layout), birthDate.Format(layout), ageAtRetirementYears, ageAtRetirementMonths)

	// --- COLA Projections for Pension and Social Security (Annual Timelines) ---
	yearsToProjectForAnnualCOLA := 0
	if ageAtRetirementYears >= 0 && ageAtRetirementYears < maxProjectionAgeYears { // ageAtRetirementYears can be 0
		yearsToProjectForAnnualCOLA = maxProjectionAgeYears - ageAtRetirementYears + 1
	}

	var projectedPensionTimeline []float64 // Annual amounts
	var pensionCOLANotes string
	selectedPensionAmountForCOLA := 0.0
	selectedPensionTypeForCOLA := ""
	switch input.CalculationSystem {
	case "FERS":
		selectedPensionAmountForCOLA = result.FERSResult.AnnualPension
		selectedPensionTypeForCOLA = "FERS"
	case "CSRS":
		selectedPensionAmountForCOLA = result.CSRSResult.FinalAnnuity
		selectedPensionTypeForCOLA = "CSRS"
	default:
		if result.CSRSResult.FinalAnnuity > result.FERSResult.AnnualPension {
			selectedPensionAmountForCOLA = result.CSRSResult.FinalAnnuity
			selectedPensionTypeForCOLA = "CSRS"
		} else {
			selectedPensionAmountForCOLA = result.FERSResult.AnnualPension
			selectedPensionTypeForCOLA = "FERS"
		}
	}

	if yearsToProjectForAnnualCOLA > 0 && selectedPensionAmountForCOLA > 0 {
		pensionCOLAInput := models.PensionCOLAInput{
			InitialPension:     selectedPensionAmountForCOLA,
			PensionType:        selectedPensionTypeForCOLA,
			AssumedCPIWRate:    input.COLAInput.COLARate,
			ProjectionStartAge: ageAtRetirementYears,
			YearsToProject:     yearsToProjectForAnnualCOLA,
		}
		pensionCOLAResult := ProjectPensionWithCOLA(pensionCOLAInput)
		projectedPensionTimeline = pensionCOLAResult.ProjectedPension
		pensionCOLANotes = pensionCOLAResult.Notes
		// if len(projectedPensionTimeline) > 0 {
		// 	log.Printf("DEBUG: Annual Projected FERS/CSRS Pension with COLA for %d years. First year rate: %.2f.", yearsToProjectForAnnualCOLA, projectedPensionTimeline[0])
		// } else {
		// 	log.Printf("DEBUG: Annual Projected FERS/CSRS Pension with COLA for %d years. Timeline is empty.", yearsToProjectForAnnualCOLA)
		// }
	} else {
		projectedPensionTimeline = make([]float64, yearsToProjectForAnnualCOLA)
		// log.Printf("DEBUG: Pension COLA projection skipped. YearsToProject: %d, SelectedPension: %.2f", yearsToProjectForAnnualCOLA, selectedPensionAmountForCOLA)
	}

	var projectedSSTimeline []float64 // Annual amounts
	var ssCOLANotes string
	if yearsToProjectForAnnualCOLA > 0 && result.SocialSecurityResult.ClaimingAmount > 0 {
		ssCOLAInput := models.SocialSecurityCOLAInput{
			InitialAnnualBenefit: result.SocialSecurityResult.ClaimingAmount,
			AssumedCPIWRate:      input.COLAInput.COLARate,
			BenefitStartAge:      input.SocialSecurityInput.ClaimAge,
			ProjectionStartAge:   ageAtRetirementYears,
			YearsToProject:       yearsToProjectForAnnualCOLA,
		}
		ssCOLAResult := ProjectSocialSecurityWithCOLA(ssCOLAInput)
		projectedSSTimeline = ssCOLAResult.ProjectedBenefit
		result.SocialSecurityResult.ProjectedAnnualBenefitTimeline = projectedSSTimeline
		ssCOLANotes = ssCOLAResult.Notes
		// if len(projectedSSTimeline) > 0 {
		// 	log.Printf("DEBUG: Annual Projected Social Security with COLA for %d years. First year rate: %.2f.", yearsToProjectForAnnualCOLA, projectedSSTimeline[0])
		// } else {
		// 	log.Printf("DEBUG: Annual Projected Social Security with COLA for %d years. Timeline is empty.", yearsToProjectForAnnualCOLA)
		// }
	} else {
		projectedSSTimeline = make([]float64, yearsToProjectForAnnualCOLA)
		result.SocialSecurityResult.ProjectedAnnualBenefitTimeline = projectedSSTimeline
		// log.Printf("DEBUG: SS COLA projection skipped. YearsToProject: %d, SS Claiming Amount: %.2f", yearsToProjectForAnnualCOLA, result.SocialSecurityResult.ClaimingAmount)
	}

	// --- Health Premium Projections (Annual Timeline) --- // <--- YOUR NEW CODE BLOCK STARTS HERE
	var projectedHealthPremiumsTimeline []float64
	var healthPremiumNotes string

	if yearsToProjectForAnnualCOLA > 0 { // Reuse the same projection period as for pension/SS
		healthInputForProjection := models.HealthPremiumCalculationInput{
			FEHBPremium:        input.HealthInput.FEHBPremium,
			MedicarePremium:    input.HealthInput.MedicarePremium,
			IncludeFEHB:        input.HealthInput.IncludeFEHB,
			IncludeMedicare:    input.HealthInput.IncludeMedicare,
			COLARate:           input.COLAInput.COLARate, // Use the general COLA rate
			YearsToProject:     yearsToProjectForAnnualCOLA,
			OtherHealthPremium: input.HealthInput.OtherHealthPremium,
		}
		healthPremiumResult := CalculateHealthPremiums(healthInputForProjection)
		projectedHealthPremiumsTimeline = healthPremiumResult.ProjectedPremiums
		healthPremiumNotes = healthPremiumResult.Notes // Store notes if any
		result.HealthResult = healthPremiumResult      // Store the full result

		// if len(projectedHealthPremiumsTimeline) > 0 {
		// 	log.Printf("DEBUG: Annual Projected Health Premiums with COLA for %d years. First year rate: %.2f. Notes: %s", yearsToProjectForAnnualCOLA, projectedHealthPremiumsTimeline[0], healthPremiumNotes)
		// } else {
		// 	log.Printf("DEBUG: Annual Projected Health Premiums with COLA for %d years. Timeline is empty. Notes: %s", yearsToProjectForAnnualCOLA, healthPremiumNotes)
		// }
	} else {
		projectedHealthPremiumsTimeline = make([]float64, 0) // Ensure it's an empty slice, not nil
		// log.Printf("DEBUG: Health Premium projection skipped. YearsToProject: %d", yearsToProjectForAnnualCOLA)
	}

	result.COLAResult = models.COLACalculationResult{Notes: strings.TrimSpace(pensionCOLANotes + " " + ssCOLANotes + " " + healthPremiumNotes)}

	initialAnnualIRSExclusion := 0.0
	switch input.CalculationSystem {
	case "FERS":
		initialAnnualIRSExclusion = result.FERSResult.IRSSimplifiedMethodExclusion * 12.0
	case "CSRS":
		initialAnnualIRSExclusion = result.CSRSResult.IRSSimplifiedMethodExclusion * 12.0
	default:
		if result.CSRSResult.FinalAnnuity > result.FERSResult.AnnualPension {
			initialAnnualIRSExclusion = result.CSRSResult.IRSSimplifiedMethodExclusion * 12.0
		} else {
			initialAnnualIRSExclusion = result.FERSResult.IRSSimplifiedMethodExclusion * 12.0
		}
	}
	// log.Printf("DEBUG: Initial Annual IRS Exclusion for pension: $%.2f", initialAnnualIRSExclusion)

	// --- Detailed Monthly Projections ---
	result.DetailedMonthlyProjections = make([]models.MonthlyFinancialProjection, 0)
	// log.Printf("DEBUG: Starting detailed monthly projections from %s up to age %d.", retirementStartDate.Format(layout), maxProjectionAgeYears)

	currentProjectionTime := retirementStartDate
	maxLoopIterations := (maxProjectionAgeYears - ageAtRetirementYears + 5) * 12 // Safety break for loop
	loopCount := 0

	for loopCount < maxLoopIterations {
		loopCount++
		currentAgeYrs, currentAgeMths, _, _ := CalculateAgeAtDate(birthDate, currentProjectionTime)

		if currentAgeYrs > maxProjectionAgeYears || (currentAgeYrs == maxProjectionAgeYears && currentProjectionTime.Month() > birthDate.Month()) {
			// log.Printf("DEBUG: Reached max projection age (%d years, month %d). Stopping monthly projection at %s.", currentAgeYrs, currentProjectionTime.Month(), currentProjectionTime.Format(layout))
			break
		}

		monthlyProjection := models.MonthlyFinancialProjection{
			Year:      currentProjectionTime.Year(),
			Month:     int(currentProjectionTime.Month()),
			AgeYears:  currentAgeYrs,
			AgeMonths: currentAgeMths,
		}

		// Determine current year index relative to retirement for COLA timelines
		// This index is 0 for the first year of retirement, 1 for the second, etc.
		yearIndexOfProjection := currentAgeYrs - ageAtRetirementYears

		// Pension for the month
		if currentProjectionTime.Compare(retirementStartDate) >= 0 { // Pension starts from retirement date
			if yearIndexOfProjection >= 0 && yearIndexOfProjection < len(projectedPensionTimeline) {
				monthlyProjection.PensionForMonth = projectedPensionTimeline[yearIndexOfProjection] / 12.0
			} else if yearIndexOfProjection < 0 {
				// Pre-retirement, no pension
				monthlyProjection.PensionForMonth = 0
			} else {
				log.Printf("WARN: Pension timeline out of bounds. YearIndex: %d, TimelineLen: %d AgeYrs: %d RetAgeYrs: %d", yearIndexOfProjection, len(projectedPensionTimeline), currentAgeYrs, ageAtRetirementYears)
			}
		}

		// Social Security for the month
		if currentAgeYrs >= input.SocialSecurityInput.ClaimAge {
			// SS COLA timeline is also indexed starting from ageAtRetirementYears (ProjectionStartAge for COLA calc)
			if yearIndexOfProjection >= 0 && yearIndexOfProjection < len(projectedSSTimeline) {
				// Check if SS has actually started based on claim age vs current age
				// The projectedSSTimeline itself handles when benefits are $0 vs >$0 based on its internal logic and BenefitStartAge.
				monthlyProjection.SocialSecurityForMonth = projectedSSTimeline[yearIndexOfProjection] / 12.0
			} else {
				log.Printf("WARN: SS timeline out of bounds. YearIndex: %d, TimelineLen: %d", yearIndexOfProjection, len(projectedSSTimeline))
			}
		}

		// TSP Withdrawals for the month
		if currentAgeYrs >= input.TSPInput.WithdrawalStartAge {
			for _, tspYearDetail := range result.TSPResult.WithdrawalSchedule {
				if tspYearDetail.Age == currentAgeYrs {
					monthlyProjection.TSPWithdrawalTraditionalForMonth = tspYearDetail.TraditionalWithdrawn / 12.0
					monthlyProjection.TSPWithdrawalRothForMonth = tspYearDetail.RothWithdrawn / 12.0
					break
				}
			}
		}

		monthlyProjection.TotalPreTaxIncomeForMonth = monthlyProjection.PensionForMonth + monthlyProjection.SocialSecurityForMonth + monthlyProjection.TSPWithdrawalTraditionalForMonth

		log.Printf("DEBUG: PreTaxIncomeForMonth: Pension=%.2f, SS=%.2f, TSP=%.2f, Total=%.2f", monthlyProjection.PensionForMonth, monthlyProjection.SocialSecurityForMonth, monthlyProjection.TSPWithdrawalTraditionalForMonth, monthlyProjection.TotalPreTaxIncomeForMonth)

		// Health Premiums for the month
		if yearIndexOfProjection >= 0 && yearIndexOfProjection < len(projectedHealthPremiumsTimeline) {
			monthlyProjection.HealthPremiumsForMonth = projectedHealthPremiumsTimeline[yearIndexOfProjection] / 12.0
		} else if yearIndexOfProjection < 0 {
			// Pre-retirement, or if timeline doesn't cover, assume zero for now.
			// A more sophisticated approach might be needed if pre-retirement premiums are a factor.
			monthlyProjection.HealthPremiumsForMonth = 0.0
		} else {
			log.Printf("WARN: Health premium timeline out of bounds. YearIndex: %d, TimelineLen: %d", yearIndexOfProjection, len(projectedHealthPremiumsTimeline))
			monthlyProjection.HealthPremiumsForMonth = 0.0 // Default to zero if out of bounds post-retirement
		}

		// Tax Calculation (Simplified: Annualized income for the year, then 1/12th of tax)
		annualPensionForTaxYear := 0.0
		if yearIndexOfProjection >= 0 && yearIndexOfProjection < len(projectedPensionTimeline) {
			annualPensionForTaxYear = projectedPensionTimeline[yearIndexOfProjection]
		}
		annualSSForTaxYear := 0.0
		if yearIndexOfProjection >= 0 && yearIndexOfProjection < len(projectedSSTimeline) {
			annualSSForTaxYear = projectedSSTimeline[yearIndexOfProjection]
		}
		annualTSPTraditionalForTaxYear := 0.0
		annualTSPRothForTaxYear := 0.0
		for _, tspYearDetail := range result.TSPResult.WithdrawalSchedule {
			if tspYearDetail.Age == currentAgeYrs {
				annualTSPTraditionalForTaxYear = tspYearDetail.TraditionalWithdrawn
				annualTSPRothForTaxYear = tspYearDetail.RothWithdrawn // Though Roth not typically taxed
				break
			}
		}

		taxablePensionIncomeAnnual := annualPensionForTaxYear - initialAnnualIRSExclusion
		if taxablePensionIncomeAnnual < 0 {
			taxablePensionIncomeAnnual = 0
		}
		totalTaxableIncomeAnnual := taxablePensionIncomeAnnual + annualSSForTaxYear + annualTSPTraditionalForTaxYear

		taxInputForYear := input.TaxInput // Copy base tax input
		taxInputForYear.TaxYear = currentProjectionTime.Year()
		taxInputForYear.TaxablePension = taxablePensionIncomeAnnual
		taxInputForYear.SocialSecurity = annualSSForTaxYear
		taxInputForYear.TSPWithdrawal = annualTSPTraditionalForTaxYear
		taxInputForYear.TSPRothWithdrawal = annualTSPRothForTaxYear

		federalTaxResultForYear := CalculateFederalTax(taxInputForYear)
		stateTaxResultForYear := CalculateStateTax(taxInputForYear)

		monthlyProjection.AllocatedFederalTaxForMonth = federalTaxResultForYear.FederalTaxOwed / 12.0
		monthlyProjection.AllocatedStateTaxForMonth = stateTaxResultForYear.StateTaxOwed / 12.0
		monthlyProjection.NetCashFlowForMonth = monthlyProjection.TotalPreTaxIncomeForMonth - monthlyProjection.AllocatedFederalTaxForMonth - monthlyProjection.AllocatedStateTaxForMonth - monthlyProjection.HealthPremiumsForMonth

		// Add notes
		if currentProjectionTime.Compare(retirementStartDate) == 0 {
			monthlyProjection.Notes += "Retirement Month. "
		}
		if currentAgeYrs == input.SocialSecurityInput.ClaimAge && currentAgeMths == 0 && currentProjectionTime.Month() == birthDate.Month() { // Approx SS claim start month
			monthlyProjection.Notes += "Social Security Claim Month (approx). "
		}
		if currentAgeYrs == input.TSPInput.WithdrawalStartAge && currentAgeMths == 0 && currentProjectionTime.Month() == birthDate.Month() { // Approx TSP start month
			monthlyProjection.Notes += "TSP Withdrawals Start (approx). "
		}
		monthlyProjection.Notes = strings.TrimSpace(monthlyProjection.Notes)

		result.DetailedMonthlyProjections = append(result.DetailedMonthlyProjections, monthlyProjection)

		// Update overall summary fields (e.g., for the first full year of retirement)
		// This logic might need refinement for how to define the 'summary' year with monthly data.
		// For now, let's use the tax details from the 12th month of projection if available, or first month.
		if len(result.DetailedMonthlyProjections) == 1 || (currentProjectionTime.Year() == retirementStartDate.Year() && currentProjectionTime.Month() == 12) {
			result.NetAfterTaxIncome = monthlyProjection.NetCashFlowForMonth * 12           // Annualized for summary
			result.TotalRetirementIncome = monthlyProjection.TotalPreTaxIncomeForMonth * 12 // Annualized for summary
			if totalTaxableIncomeAnnual > 0 {
				result.EffectiveTaxRate = (federalTaxResultForYear.FederalTaxOwed + stateTaxResultForYear.StateTaxOwed) / totalTaxableIncomeAnnual
			} else {
				result.EffectiveTaxRate = 0
			}
			result.TaxResult = models.TaxCalculationResult{
				GrossIncome:       totalTaxableIncomeAnnual,
				TaxableIncome:     federalTaxResultForYear.TaxableIncome, // Assuming CalculateFederalTax result provides this
				FederalTaxOwed:    federalTaxResultForYear.FederalTaxOwed,
				StateTaxOwed:      stateTaxResultForYear.StateTaxOwed,
				NetAfterTaxIncome: result.NetAfterTaxIncome, // Store the annualized summary net income
				EffectiveTaxRate:  result.EffectiveTaxRate,
				Notes:             strings.TrimSpace(federalTaxResultForYear.Notes + " " + stateTaxResultForYear.Notes),
			}
		}

		currentProjectionTime = currentProjectionTime.AddDate(0, 1, 0) // Move to the next month
	}
	if loopCount >= maxLoopIterations {
		log.Printf("ERROR: Monthly projection loop hit max iterations (%d). Terminating early.", maxLoopIterations)
		result.Notes += " Error: Projection terminated due to exceeding max iterations."
	}

	log.Printf("DEBUG: Finished detailed monthly projections. Total months projected: %d", len(result.DetailedMonthlyProjections))

	// --- Survivor Benefit Calculation ---
	result.SurvivorResult = CalculateSurvivorBenefits(input.SurvivorInput, result.FERSResult, result.CSRSResult)

	// HealthResult is now populated by the annual projection block earlier.
	// The call previously here: result.HealthResult = CalculateHealthPremiums(input.HealthInput) is removed.

	// Consolidate general notes and finalize result.Notes
	// First, take notes collected at the top of the function (e.g., from pension system choice)
	generalNotesString := strings.TrimSpace(strings.Join(notes, " "))

	if generalNotesString != "" {
		if result.Notes != "" {
			// If result.Notes already has content (e.g., max iteration error), prepend general notes
			result.Notes = generalNotesString + " | " + result.Notes
		} else {
			result.Notes = generalNotesString
		}
	}

	// Now, apply default success message if no specific notes or errors are present
	if result.Notes == "" {
		result.Notes = "Retirement projection calculated successfully."
	} else {
		// If there are notes, but no "Error:", and it doesn't already say "calculated", add a general success indicator.
		if !strings.Contains(result.Notes, "Error:") &&
			!strings.Contains(result.Notes, "Retirement projection calculated successfully.") &&
			!strings.HasSuffix(result.Notes, "calculated.") { // Avoid double "calculated."
			result.Notes += " Retirement projection calculated."
		}
	}

	log.Printf("DEBUG: Finished detailed monthly projections. Total months projected: %d", len(result.DetailedMonthlyProjections))

	// (5) Monte Carlo Simulation (optional)
	if input.MonteCarloInput.NumSimulations > 0 && input.MonteCarloInput.Years > 0 {
		result.MonteCarloResult = RunMonteCarloSimulation(input.MonteCarloInput)
	}

	// (6) Aggregate summary fields (fill as needed)
	// This section is intended for future expansion or specific summary logic not yet implemented.

	// if finalResultBytes, err := json.MarshalIndent(result, "", "  "); err == nil {
	// 	log.Printf("INFO: PerformRetirementCalculation Final Result:\n%s", string(finalResultBytes))
	// } else {
	// 	log.Printf("ERROR: Failed to marshal final result: %v", err)
	// }
	return result
} // Closing brace for PerformRetirementCalculation function
