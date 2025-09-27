package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"ferex/backend/calculation"
	"ferex/backend/models"
	"ferex/backend/scenario" // Corrected import path for scenario package
	"log"                    // Added for logging

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// PerformFERSCalculation calls the backend FERS calculation logic,
// deriving values from dates.
func (a *App) PerformFERSCalculation(input models.FERSCalculationInput) (models.FERSCalculationResult, error) {
	// --- Date Derivations ---
	var dob, scd, retireDate, switchDate time.Time
	var ageAtRetirement int
	var yearsOfService float64
	var mra int
	var csrsServiceYearsActual float64 = 0
	var fersServiceYearsActual float64 = 0
	var csrsAnnuityComponent float64 = 0
	var err error
	parseLayout := "2006-01-02" // YYYY-MM-DD

	// Core dates are required by frontend validation now
	dob, err = time.Parse(parseLayout, input.DateOfBirth)
	if err != nil {
		return models.FERSCalculationResult{}, fmt.Errorf("invalid Date of Birth format: %w", err)
	}
	retireDate, err = time.Parse(parseLayout, input.PlannedRetirementDate)
	if err != nil {
		return models.FERSCalculationResult{}, fmt.Errorf("invalid Planned Retirement Date format: %w", err)
	}
	scd, err = time.Parse(parseLayout, input.ServiceComputationDate)
	if err != nil {
		return models.FERSCalculationResult{}, fmt.Errorf("invalid Service Computation Date format: %w", err)
	}

	if retireDate.Before(dob) {
		return models.FERSCalculationResult{}, errors.New("retirement date cannot be before date of birth")
	}
	if retireDate.Before(scd) {
		return models.FERSCalculationResult{}, errors.New("retirement date cannot be before service computation date")
	}

	// Derive core values
	ageAtRetirement = calculateAge(dob, retireDate)
	yearsOfService = calculateYearsOfService(scd, retireDate)
	mra = determineMRA(dob)

	// Calculate CSRS and FERS service years if switch date is provided
	// totalCreditableServiceForFERSCalc will be the service years used for FERS component calculation (excluding sick leave)
	// yearsForSRSCalc will be actual FERS service for supplement calculation
	totalCreditableServiceForFERSCalc := yearsOfService // Default to total service if not a transferee
	yearsForSRSCalc := yearsOfService                   // Default for SRS calculation

	if input.SwitchedToFERSDate != "" {
		switchDate, err = time.Parse(parseLayout, input.SwitchedToFERSDate)
		if err != nil {
			return models.FERSCalculationResult{}, fmt.Errorf("invalid Switched to FERS Date format: %w", err)
		}
		if retireDate.Before(switchDate) {
			return models.FERSCalculationResult{}, errors.New("retirement date cannot be before switch date")
		}
		if switchDate.Before(scd) {
			return models.FERSCalculationResult{}, errors.New("switch date cannot be before service computation date")
		}
		csrsServiceYearsActual = calculateYearsOfService(scd, switchDate)
		fersServiceYearsActual = calculateYearsOfService(switchDate, retireDate)
		totalCreditableServiceForFERSCalc = fersServiceYearsActual
		yearsForSRSCalc = fersServiceYearsActual
	} else {
		// If no switch date, all service is FERS
		fersServiceYearsActual = yearsOfService
		totalCreditableServiceForFERSCalc = fersServiceYearsActual
		yearsForSRSCalc = fersServiceYearsActual
	}

	// Calculate age at retirement in years and months for more precise calculations
	var ageAtRetirementMonths int = 0
	// Use the refined calculateAgeParts function to get more precise age information
	ageAtRetirement, ageAtRetirementMonths, _ = calculateAgeParts(dob, retireDate)

	// Comprehensive retirement type determination based on OPM rules
	var retirementScenario string

	// Check for Special Provisions (Law Enforcement, Firefighter, Air Traffic Controller)
	// Note: This would require additional input fields to be fully implemented
	// For now, we'll assume no special provisions apply

	// 1. Check for Immediate Retirement eligibility (no age reduction)
	if ageAtRetirement >= mra && yearsOfService >= 30 {
		// MRA with 30 years of service
		retirementScenario = "ImmediateMRA30"
	} else if ageAtRetirement >= 60 && yearsOfService >= 20 {
		// Age 60 with 20 years of service
		retirementScenario = "Immediate6020"
	} else if ageAtRetirement >= 62 && yearsOfService >= 5 {
		// Age 62 with 5 years of service
		retirementScenario = "Immediate6205"
	} else if ageAtRetirement >= mra && yearsOfService >= 10 && yearsOfService < 30 {
		// MRA with at least 10 but less than 30 years of service
		// This has an age-based reduction if taken before age 62
		retirementScenario = "MRA+10Voluntary"
	} else if ageAtRetirement >= mra && yearsOfService >= 20 && yearsOfService < 30 {
		// This is a special case that's technically covered by MRA+10 but with 20+ years
		// We distinguish it for clarity in some calculations
		retirementScenario = "MRA+20Voluntary"
	} else if ageAtRetirement < mra && yearsOfService >= 10 {
		// Less than MRA but with at least 10 years service
		// This is a deferred retirement that begins when MRA is reached
		retirementScenario = "Deferred"
	} else if ageAtRetirement >= mra && ageAtRetirement < 62 && yearsOfService >= 10 {
		// MRA with 10+ years, but choosing to postpone receipt to avoid or reduce the age penalty
		// Note: This would require additional input to know if they're postponing
		// For now, we'll assume they're taking it immediately with the reduction
		// retirementScenario = "MRA+10Postponed"
		retirementScenario = "MRA+10Voluntary"
	} else if yearsOfService >= 5 {
		// At least 5 years of service but not meeting any immediate retirement criteria
		// This is a deferred retirement that begins at age 62
		retirementScenario = "Deferred62"
	} else {
		// Less than 5 years of service - not eligible for retirement
		retirementScenario = "NotEligible"
	}

	// 2. Check for Early Retirement (VERA) eligibility
	// Note: This would require additional input fields to be fully implemented
	// For now, we'll assume VERA is not available
	// if veraOffered && ageAtRetirement >= 50 && yearsOfService >= 20 {
	//     retirementScenario = "VERAAge50Service20"
	// } else if veraOffered && yearsOfService >= 25 {
	//     retirementScenario = "VERAService25"
	// }

	// 3. Check for Disability Retirement eligibility
	// Note: This would require additional input fields to be fully implemented
	// For now, we'll assume disability retirement is not applicable
	// if disabilityEligible && yearsOfService >= 18 {
	//     retirementScenario = "DisabilityRetirement"
	// }

	// 4. Special case for Discontinued Service Retirement (RIF, etc.)
	// Note: This would require additional input fields to be fully implemented
	// For now, we'll assume DSR is not applicable
	// if discontinuedService && ageAtRetirement >= 50 && yearsOfService >= 20 {
	//     retirementScenario = "DSRAge50Service20"
	// } else if discontinuedService && yearsOfService >= 25 {
	//     retirementScenario = "DSRService25"
	// }

	// --- Calculate Part-Time Proration Factor ---
	var totalRawServiceDays float64 = 0
	var totalCreditedServiceDays float64 = 0
	actualProrationFactor := 1.0
	prorationAppliedByCaller := false
	const standardFullTimeHoursPerWeek = 40.0

	if len(input.ServicePeriods) > 0 {
		for _, sp := range input.ServicePeriods {
			spStartDate, err := time.Parse(parseLayout, sp.StartDate)
			if err != nil {
				return models.FERSCalculationResult{}, fmt.Errorf("invalid ServicePeriod StartDate '%s': %w", sp.StartDate, err)
			}
			spEndDate, err := time.Parse(parseLayout, sp.EndDate)
			if err != nil {
				return models.FERSCalculationResult{}, fmt.Errorf("invalid ServicePeriod EndDate '%s': %w", sp.EndDate, err)
			}

			if spEndDate.Before(spStartDate) {
				return models.FERSCalculationResult{}, fmt.Errorf("service period end date '%s' cannot be before start date '%s'", sp.EndDate, sp.StartDate)
			}

			// OPM typically includes the end date, so add 1 day to the difference if just subtracting.
			// However, to be consistent with calculateYearsOfService which uses simple subtraction:
			periodDuration := spEndDate.Sub(spStartDate)
			periodDays := periodDuration.Hours() / 24.0

			totalRawServiceDays += periodDays

			if sp.ServiceCategory == "Civilian" && sp.IsPartTime && sp.HoursPerWeekIfPartTime != nil && *sp.HoursPerWeekIfPartTime > 0 && *sp.HoursPerWeekIfPartTime < standardFullTimeHoursPerWeek {
				periodCreditFactor := *sp.HoursPerWeekIfPartTime / standardFullTimeHoursPerWeek
				totalCreditedServiceDays += periodDays * periodCreditFactor
			} else {
				// Treat CivilianFullTime (IsPartTime == false), Military, etc., as fully credited for this period's duration
				totalCreditedServiceDays += periodDays
			}
		}

		if totalRawServiceDays > 0 {
			actualProrationFactor = totalCreditedServiceDays / totalRawServiceDays
			if actualProrationFactor < 0 {
				actualProrationFactor = 0
			} // Should not happen
			if actualProrationFactor > 1.0 {
				actualProrationFactor = 1.0
			} // Cap at 100%
		}
	} // else, no service periods, proration factor remains 1.0

	// Determine if proration was effectively applied
	if actualProrationFactor > 0 && actualProrationFactor < 1.0 { // Check > 0 to avoid issues with zero service
		prorationAppliedByCaller = true
	}

	// Adjust High3Salary if proration is applied
	// Create a copy of the input to modify High3Salary, or modify directly if input is already a copy for this func scope.
	// Since 'input' is a value type (not a pointer) for this function, we can modify it directly.
	originalHigh3 := input.High3Salary // For potential reference, though not strictly needed if input is a copy
	if prorationAppliedByCaller {
		input.High3Salary = originalHigh3 * actualProrationFactor
	}

	// If FERS Transferee, calculate CSRS component using the (potentially prorated) High3
	if input.SwitchedToFERSDate != "" && csrsServiceYearsActual > 0 {
		csrsAnnuityComponent = calculation.CalculateCSRSComponentAnnuity(input.High3Salary, csrsServiceYearsActual)
	}

	// --- Call Core Calculation ---
	// Pass the input struct (with potentially modified High3Salary) and other derived values
	result := calculation.CalculateFERS(
		input,                             // models.FERSCalculationInput (High3Salary may be adjusted)
		ageAtRetirement,                   // int (ageAtRetirementYears)
		ageAtRetirementMonths,             // int
		totalCreditableServiceForFERSCalc, // float64 (FERS service years before sick leave)
		mra,                               // int (mraYears)
		retirementScenario,                // string (retirementType)
		input.ExpectedSSBenefitAt62,       // float64
		yearsForSRSCalc,                   // float64 (Actual FERS service for SRS)
		prorationAppliedByCaller,          // bool
		actualProrationFactor,             // float64
		yearsOfService,                    // float64 (totalCombinedActualServiceYears)
	)

	// Add CSRS component to the FERS result if applicable
	if csrsAnnuityComponent > 0 {
		monthlyCsrsComponent := csrsAnnuityComponent / 12.0
		result.Notes = fmt.Sprintf("CSRS Annuity Component: $%.2f monthly ($%.2f annual).\n", monthlyCsrsComponent, csrsAnnuityComponent) + result.Notes
		result.MonthlyBasicAnnuity += monthlyCsrsComponent
		result.MonthlyProratedPension += monthlyCsrsComponent // Assuming proration applies to High3 which affects both components
		result.MonthlyPension += monthlyCsrsComponent

		// Update derived annual values
		result.BasicAnnuity = result.MonthlyBasicAnnuity * 12.0
		result.ProratedPension = result.MonthlyProratedPension * 12.0
		result.AnnualPension = result.MonthlyPension * 12.0

		// Recalculate survivor benefit reduction on the combined annuity for transferees
		combinedMonthlyAnnuity := result.MonthlyBasicAnnuity // This now includes both FERS and CSRS components
		if input.SurvivorBenefitElection != "None" {
			// Use FERS survivor cost rules on the combined annuity
			monthlySurvivorReduction := 0.0
			switch input.SurvivorBenefitElection {
			case "Full50Percent":
				monthlySurvivorReduction = combinedMonthlyAnnuity * 0.10
			case "Partial25Percent":
				monthlySurvivorReduction = combinedMonthlyAnnuity * 0.05
			}
			result.MonthlySurvivorBenefitReduction = monthlySurvivorReduction
			result.MonthlyPension -= monthlySurvivorReduction

			// Update derived annual values
			result.SurvivorBenefitReduction = result.MonthlySurvivorBenefitReduction * 12.0
			result.AnnualPension = result.MonthlyPension * 12.0
			result.Notes += "\nSurvivor benefit reduction recalculated on combined FERS+CSRS annuity for transferee."
		}
	}

	return result, nil
}

// PerformCSRSCalculation calls the backend CSRS calculation logic,
// deriving values from dates.
func (a *App) PerformCSRSCalculation(input models.CSRSCalculationInput) (models.CSRSCalculationResult, error) {
	// --- Date Derivations ---
	var dob, scd, retireDate time.Time
	var ageAtRetirement int
	var yearsOfService float64
	var err error
	parseLayout := "2006-01-02" // YYYY-MM-DD

	dob, err = time.Parse(parseLayout, input.DateOfBirth)
	if err != nil {
		return models.CSRSCalculationResult{}, fmt.Errorf("invalid Date of Birth format: %w", err)
	}
	retireDate, err = time.Parse(parseLayout, input.PlannedRetirementDate)
	if err != nil {
		return models.CSRSCalculationResult{}, fmt.Errorf("invalid Planned Retirement Date format: %w", err)
	}
	scd, err = time.Parse(parseLayout, input.ServiceComputationDate)
	if err != nil {
		return models.CSRSCalculationResult{}, fmt.Errorf("invalid Service Computation Date format: %w", err)
	}

	if retireDate.Before(dob) {
		return models.CSRSCalculationResult{}, errors.New("retirement date cannot be before date of birth")
	}
	if retireDate.Before(scd) {
		return models.CSRSCalculationResult{}, errors.New("retirement date cannot be before service computation date")
	}

	ageAtRetirement = calculateAge(dob, retireDate)
	yearsOfService = calculateYearsOfService(scd, retireDate)

	// Call backend logic from calculation package
	// Note: calculation.CalculateCSRS expects retirementAge as float64
	return calculation.CalculateCSRS(float64(ageAtRetirement), yearsOfService, input), nil
}

// CalculateTSP accepts a ScenarioVariant, extracts TSP-related information, derives necessary age data,
// populates a models.TSPCalculationInput, validates it, and calls the backend TSP calculation logic.
func (a *App) CalculateTSP(variant scenario.ScenarioVariant) (models.TSPCalculationResult, error) {
	log.Println("!!!!!!!!!! APP.GO CALCULATETSP FUNCTION ENTERED !!!!!!!!!!")
	var tspInput models.TSPCalculationInput
	var err error

	// --- Date Derivations and Age Calculations ---
	dob, err := time.Parse("2006-01-02", variant.DateOfBirth)
	if err != nil {
		return models.TSPCalculationResult{}, fmt.Errorf("invalid date of birth format for '%s': %w", variant.DateOfBirth, err)
	}
	retireDate, err := time.Parse("2006-01-02", variant.PlannedRetirementDate)
	if err != nil {
		return models.TSPCalculationResult{}, fmt.Errorf("invalid planned retirement date format for '%s': %w", variant.PlannedRetirementDate, err)
	}

	// For TSPCalculationInput, CurrentAgeYears often means the age at the start of the projection period.
	// If projecting for future retirement, this would be the retirement age.
	// If already retired, it might be the current actual age or retirement age depending on context.
	// The backend calculation.CalculateTSP handles the timeline internally based on these inputs.
	// Let's set CurrentAgeYears to retirementAgeYears for projections starting at retirement.
	// BirthYear is crucial for RMD calculations.
	retirementAgeYears := calculateAge(dob, retireDate)
	birthYearVal, err := strconv.Atoi(strings.Split(variant.DateOfBirth, "-")[0])
	if err != nil {
		return models.TSPCalculationResult{}, fmt.Errorf("could not parse birth year from DOB '%s': %w", variant.DateOfBirth, err)
	}

	tspInput.CurrentAgeYears = retirementAgeYears // Age at the start of the main TSP projection period (retirement)
	tspInput.RetirementAgeYears = retirementAgeYears
	tspInput.BirthYear = birthYearVal

	// --- Basic Validation ---
	if tspInput.RetirementAgeYears < 18 { // Arbitrary minimum sensible retirement age
		log.Printf("Warning: TSP Calculation using a low RetirementAgeYears: %d", tspInput.RetirementAgeYears)
		// Depending on rules, might return error: return models.TSPCalculationResult{}, errors.New("retirement age is too low")
	}
	if retireDate.Before(dob) {
		return models.TSPCalculationResult{}, errors.New("planned retirement date cannot be before date of birth")
	}

	// --- Map ScenarioVariant TSP fields to models.TSPCalculationInput ---
	// BaseSalaryForContributions: Use High3Salary from ScenarioVariant. Adjust if a different salary basis is needed for TSP contributions.
	tspInput.BaseSalaryForContributions = derefFloat64(variant.High3Salary, 0.0)

	tspInput.CurrentTraditionalBalance = derefFloat64(variant.TSPBalanceTraditional, 0.0)
	tspInput.CurrentRothBalance = derefFloat64(variant.TSPBalanceRoth, 0.0)
	tspInput.CurrentLoanBalance = derefFloat64(variant.TSPLoanBalance, 0.0)
	tspInput.CurrentAllocation = convertTSPFundAllocation(variant.TSPCurrentAllocationToFunds)

	// Contributions Pre-Retirement
	if variant.TSPContributionPercentagePreRetirement != nil {
		tspInput.EmployeeContributionPercentage = derefFloat64(variant.TSPContributionPercentagePreRetirement, 0.0) / 100.0 // Convert to decimal
		tspInput.IsContributionPercentage = true
	} else if variant.TSPAnnualContributionPreRetirement != nil {
		tspInput.EmployeeContributionAmount = derefFloat64(variant.TSPAnnualContributionPreRetirement, 0.0)
		tspInput.IsContributionPercentage = false
	} // Else, no contributions specified (IsContributionPercentage defaults to false, Amount defaults to 0)

	if variant.TSPContributeUntil != nil {
		contributeUntil := strings.ToLower(derefString(variant.TSPContributeUntil, "retirement"))
		if contributeUntil == "retirement" {
			tspInput.ContributeUntilRetirement = true
		} else if contributeUntil == "specificage" { // Assuming "SpecificAge" from ScenarioVariant
			tspInput.ContributeUntilRetirement = false
			tspInput.ContributionStopAge = derefInt(variant.TSPContributionStopAge, tspInput.RetirementAgeYears)
		} else {
			tspInput.ContributeUntilRetirement = true // Default to retirement if unrecognized
		}
	} else {
		tspInput.ContributeUntilRetirement = true // Default to contributing until retirement if not specified
	}

	tspInput.CatchUpContributionsEligible = variant.TSPCatchUpContributionsEligible
	// Percentages are 0-100 in TSPCalculationInput for Roth/Trad allocation of contributions
	tspInput.TraditionalContributionAllocationPct = derefFloat64(variant.TSPContributionAllocationTraditionalPercent, 100.0)
	tspInput.RothContributionAllocationPct = derefFloat64(variant.TSPContributionAllocationRothPercent, 0.0)
	tspInput.ContributionFundAllocation = convertTSPFundAllocation(variant.TSPContributionAllocationToFunds)

	// Growth Assumptions
	tspInput.UserReturnAssumptions = mapReturnAssumptions(variant)         // Uses helper
	tspInput.ExpenseRatio = derefFloat64(variant.TSPExpenseRatio, 0.00043) // Default based on recent TSP (e.g., 0.043% -> 0.00043)

	// TODO: ExpectedAnnualInflationRate should ideally come from ScenarioVariant general assumptions or global config.
	tspInput.ExpectedAnnualInflationRate = 0.025 // Placeholder: 2.5%

	tspInput.FutureAllocationStrategy = derefString(variant.TSPFutureAllocationStrategy, "MaintainCurrent")
	tspInput.PostRetirementAllocation = convertTSPFundAllocation(variant.TSPPostRetirementAllocation)

	// Withdrawals Post-Retirement
	tspInput.WithdrawalStrategy = derefString(variant.TSPWithdrawalStrategy, "None")
	tspInput.WithdrawalFixedAmountValue = derefFloat64(variant.TSPWithdrawalFixedAmountValue, 0.0)
	// Ensure percentage is decimal for backend (e.g., 4% -> 0.04)
	tspInput.WithdrawalPercentageValue = derefFloat64(variant.TSPWithdrawalPercentageValue, 0.0)

	if variant.TSPWithdrawalStartDate != nil {
		withdrawalStartDateOpt := strings.ToLower(derefString(variant.TSPWithdrawalStartDate, "retirement"))
		if withdrawalStartDateOpt == "retirement" {
			tspInput.WithdrawalStartDate = "Retirement" // Backend expects this string literal or "SpecificAge"
			tspInput.WithdrawalStartAge = tspInput.RetirementAgeYears
		} else if withdrawalStartDateOpt == "specificage" { // Assuming "SpecificAge" from ScenarioVariant
			tspInput.WithdrawalStartDate = "SpecificAge"
			tspInput.WithdrawalStartAge = derefInt(variant.TSPWithdrawalStartAge, tspInput.RetirementAgeYears)
		} else {
			tspInput.WithdrawalStartDate = "Retirement" // Default
			tspInput.WithdrawalStartAge = tspInput.RetirementAgeYears
		}
	} else {
		tspInput.WithdrawalStartDate = "Retirement" // Default if not specified
		tspInput.WithdrawalStartAge = tspInput.RetirementAgeYears
	}

	tspInput.WithdrawalOrder = derefString(variant.TSPWithdrawalOrder, "ProRata") // Default to ProRata

	// TODO: YearsToProjectWithdrawals should ideally come from ScenarioVariant user preferences or a sensible default.
	// Projecting until age 95 is a common heuristic.
	yearsToProject := 95 - tspInput.RetirementAgeYears
	if yearsToProject < 0 { // Already past 95 or retiring at/after 95
		yearsToProject = 0
	}
	if yearsToProject == 0 && tspInput.RetirementAgeYears <= 95 { // If retiring at 95, project for that year
		yearsToProject = 1
	}
	tspInput.YearsToProjectWithdrawals = yearsToProject
	if tspInput.YearsToProjectWithdrawals == 0 && tspInput.WithdrawalStrategy != "None" && tspInput.RetirementAgeYears < 95 {
		log.Printf("Warning: YearsToProjectWithdrawals is 0 for TSP, but a withdrawal strategy ('%s') is set and retirement age (%d) is less than 95.", tspInput.WithdrawalStrategy, tspInput.RetirementAgeYears)
	}

	log.Printf("Calling calculation.CalculateTSP with mapped input: CurrentAgeYears=%d, RetirementAgeYears=%d, BirthYear=%d, TradBalance=%.2f, RothBalance=%.2f, WithdrawalStrategy='%s', WithdrawalPercentageValue=%.8f, YearsToProjectWithdrawals=%d",
		tspInput.CurrentAgeYears, tspInput.RetirementAgeYears, tspInput.BirthYear, tspInput.CurrentTraditionalBalance, tspInput.CurrentRothBalance, tspInput.WithdrawalStrategy, tspInput.WithdrawalPercentageValue, tspInput.YearsToProjectWithdrawals)

	// Call the backend calculation function
	result := calculation.CalculateTSP(tspInput)
	// calculation.CalculateTSP does not return an error per its current definition in tsp_calc.go
	return result, nil
}

// --- Integrated Retirement Calculation ---

// PerformRetirementCalculation runs the full retirement scenario calculation, including Monte Carlo, and returns the result.
func (a *App) PerformRetirementCalculation(input models.RetirementCalculationInput) (models.RetirementCalculationResult, error) {
	// Call the orchestrator in calculation/retirement_flow.go
	result := calculation.PerformRetirementCalculation(input)
	return result, nil
}

// --- Helper functions for TSP Calculation Data Mapping ---

func derefFloat64(f *float64, defaultValue float64) float64 {
	if f != nil {
		return *f
	}
	return defaultValue
}

func derefInt(i *int, defaultValue int) int {
	if i != nil {
		return *i
	}
	return defaultValue
}

func derefString(s *string, defaultValue string) string {
	if s != nil {
		return *s
	}
	return defaultValue
}

// convertTSPFundAllocation takes scenario.TSPFundAllocation and returns models.TSPFundAllocationPercentages
func convertTSPFundAllocation(scenarioAlloc *models.TSPFundAllocation) models.TSPFundAllocationPercentages {
	if scenarioAlloc == nil {
		return models.TSPFundAllocationPercentages{} // Return empty if input is nil
	}
	return models.TSPFundAllocationPercentages{
		G:         scenarioAlloc.G,
		F:         scenarioAlloc.F,
		C:         scenarioAlloc.C,
		S:         scenarioAlloc.S,
		I:         scenarioAlloc.I,
		LFundName: scenarioAlloc.LFundName,
	}
}

// mapReturnAssumptions converts ScenarioVariant's TSP return/volatility inputs to models.TSPReturnAssumptions
func mapReturnAssumptions(variantInput scenario.ScenarioVariant) models.TSPReturnAssumptions {
	calcReturnAssumptions := models.TSPReturnAssumptions{}

	// Map individual fund return assumptions from UserReturnAssumptionsTSP
	if variantInput.UserReturnAssumptionsTSP != nil {
		calcReturnAssumptions.G = variantInput.UserReturnAssumptionsTSP.G
		calcReturnAssumptions.F = variantInput.UserReturnAssumptionsTSP.F
		calcReturnAssumptions.C = variantInput.UserReturnAssumptionsTSP.C
		calcReturnAssumptions.S = variantInput.UserReturnAssumptionsTSP.S
		calcReturnAssumptions.I = variantInput.UserReturnAssumptionsTSP.I
	}

	// Map volatility assumptions from TSPVolatilityAssumptions
	if variantInput.TSPVolatilityAssumptions != nil {
		calcReturnAssumptions.VolatilityG = variantInput.TSPVolatilityAssumptions.GStdDev
		calcReturnAssumptions.VolatilityF = variantInput.TSPVolatilityAssumptions.FStdDev
		calcReturnAssumptions.VolatilityC = variantInput.TSPVolatilityAssumptions.CStdDev
		calcReturnAssumptions.VolatilityS = variantInput.TSPVolatilityAssumptions.SStdDev
		calcReturnAssumptions.VolatilityI = variantInput.TSPVolatilityAssumptions.IStdDev
	}

	// Handle overall growth rates if UserReturnAssumptionsTSP is not provided or for specific phases
	// This logic prioritizes individual fund returns if UserReturnAssumptionsTSP is present.
	// If TSPExpectedAnnualGrowthRatePreRetirement is meant to be an 'Overall' rate when individual rates are absent:
	if variantInput.UserReturnAssumptionsTSP == nil && variantInput.TSPExpectedAnnualGrowthRatePreRetirement != nil {
		calcReturnAssumptions.Overall = variantInput.TSPExpectedAnnualGrowthRatePreRetirement
		// Potentially set UseOverallForPre = true, depending on backend logic needs
	}
	// Similarly for TSPExpectedAnnualGrowthRatePostRetirement for post-retirement phase.
	// The backend calculation will ultimately decide how to use these fields.

	return calcReturnAssumptions
}

// --- Scenario File Operations ---

// SaveScenarioFile saves the provided scenario data to the specified file path.
// The 'data' parameter will be unmarshalled by Wails from the frontend JSON object.
func (a *App) SaveScenarioFile(filePath string, data scenario.FerexFile) error {
	// The scenario.SaveFerexFile function expects a pointer, so we pass the address of data.
	return scenario.SaveFerexFile(filePath, &data)
}

// LoadScenarioFile loads scenario data from the specified file path.
func (a *App) LoadScenarioFile(filePath string) (*scenario.FerexFile, error) {
	return scenario.LoadFerexFile(filePath)
}

// OpenFileAndLoadScenario displays a file open dialog and loads the selected Ferex scenario file.
func (a *App) OpenFileAndLoadScenario() (*scenario.FerexFile, error) {
	dialogOpts := runtime.OpenDialogOptions{
		Title: "Open Ferex Scenario File",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "FeReX Files (*.ferex)",
				Pattern:     "*.ferex",
			},
		},
	}

	filePath, err := runtime.OpenFileDialog(a.ctx, dialogOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to open file dialog: %w", err)
	}

	if filePath == "" {
		// User cancelled the dialog
		return nil, nil
	}

	return scenario.LoadFerexFile(filePath)
}

// SaveScenarioWithDialog displays a file save dialog and saves the provided FerexFile data.
// The 'currentFerexFile' parameter will be unmarshalled by Wails from the frontend JSON object.
func (a *App) SaveScenarioWithDialog(currentFerexFile scenario.FerexFile) error {
	dialogOpts := runtime.SaveDialogOptions{
		Title: "Save FeReX Scenario File",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "FeReX Files (*.ferex)",
				Pattern:     "*.ferex",
			},
		},
		DefaultFilename: "MyFerexScenario.ferex",
	}

	filePath, err := runtime.SaveFileDialog(a.ctx, dialogOpts)
	if err != nil {
		return fmt.Errorf("failed to open save file dialog: %w", err)
	}

	if filePath == "" {
		// User cancelled the dialog
		return nil
	}

	// The scenario.SaveFerexFile function expects a pointer.
	return scenario.SaveFerexFile(filePath, &currentFerexFile)
}

// --- Helper Functions ---

// calculateAge determines the age in whole years based on birth and reference dates.
func calculateAge(birthDate, referenceDate time.Time) int {
	age := referenceDate.Year() - birthDate.Year()

	// Handle leap year birthday edge case (Feb 29)
	birthMonth, birthDay := birthDate.Month(), birthDate.Day()
	refMonth, refDay := referenceDate.Month(), referenceDate.Day()

	if birthMonth == time.February && birthDay == 29 && !isLeapYear(referenceDate.Year()) {
		// For leap year birthdays (Feb 29), consider March 1 in non-leap years
		if refMonth < time.March || (refMonth == time.March && refDay < 1) {
			age--
		}
	} else if refMonth < birthMonth || (refMonth == birthMonth && refDay < birthDay) {
		age--
	}

	return age
}

// calculateAgeParts determines age in years, months, and days.
// This follows OPM's method of calculating age more precisely.
func calculateAgeParts(birthDate, referenceDate time.Time) (years, months, days int) {
	// Extract components
	birthYear, birthMonth, birthDay := birthDate.Date()
	refYear, refMonth, refDay := referenceDate.Date()

	// Calculate years
	years = refYear - birthYear

	// Handle leap year birthday edge case (Feb 29)
	if birthMonth == time.February && birthDay == 29 && !isLeapYear(refYear) {
		// For leap year birthdays (Feb 29), consider March 1 in non-leap years
		if refMonth < time.March || (refMonth == time.March && refDay < 1) {
			years--
			months = int(refMonth) + 12 - int(time.March)
			days = refDay
			return
		}
	}

	// Adjust years if needed based on month and day
	if refMonth < birthMonth || (refMonth == birthMonth && refDay < birthDay) {
		years--
		// Calculate months remaining in the year
		if refMonth < birthMonth {
			months = int(refMonth) + 12 - int(birthMonth)
		} else {
			months = 0
		}

		// Calculate days
		if refDay < birthDay {
			// Get days in the previous month
			var daysInPrevMonth int
			prevMonth := refMonth - 1
			if prevMonth == 0 {
				prevMonth = 12
				daysInPrevMonth = 31 // December
			} else {
				// This is a simplification - could use time.Date to get exact days
				daysInMonth := map[time.Month]int{
					time.January: 31, time.February: 28, time.March: 31,
					time.April: 30, time.May: 31, time.June: 30,
					time.July: 31, time.August: 31, time.September: 30,
					time.October: 31, time.November: 30, time.December: 31,
				}
				// Adjust for leap year if February
				if prevMonth == time.February && isLeapYear(refYear) {
					daysInPrevMonth = 29
				} else {
					daysInPrevMonth = daysInMonth[prevMonth]
				}
			}
			days = refDay + daysInPrevMonth - birthDay
			// If we borrowed days from previous month, decrement months
			if months > 0 {
				months--
			} else {
				months = 11
				years--
			}
		} else {
			days = refDay - birthDay
		}
	} else {
		// Calculate months and days when no year adjustment needed
		months = int(refMonth) - int(birthMonth)
		if refDay < birthDay {
			// Similar logic as above for days calculation
			var daysInPrevMonth int
			prevMonth := refMonth - 1
			if prevMonth == 0 {
				prevMonth = 12
				daysInPrevMonth = 31 // December
			} else {
				daysInMonth := map[time.Month]int{
					time.January: 31, time.February: 28, time.March: 31,
					time.April: 30, time.May: 31, time.June: 30,
					time.July: 31, time.August: 31, time.September: 30,
					time.October: 31, time.November: 30, time.December: 31,
				}
				if prevMonth == time.February && isLeapYear(refYear) {
					daysInPrevMonth = 29
				} else {
					daysInPrevMonth = daysInMonth[prevMonth]
				}
			}
			days = refDay + daysInPrevMonth - birthDay
			if months > 0 {
				months--
			} else {
				months = 11
				years--
			}
		} else {
			days = refDay - birthDay
		}
	}

	return years, months, days
}

// isLeapYear determines if a year is a leap year
func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// calculateYearsOfService determines the precise years of service (including fractional) between two dates.
// To maintain compatibility with existing tests, this uses the 365.25 days per year method.
// For OPM calculations, the days are typically calculated using a 360-day year (30 days per month),
// but we're maintaining the existing implementation for consistency.
func calculateYearsOfService(startDate, endDate time.Time) float64 {
	if endDate.Before(startDate) {
		return 0.0
	}
	// Calculate difference in days
	diff := endDate.Sub(startDate)
	days := diff.Hours() / 24
	// Approximate years (using 365.25 for average days per year)
	years := days / 365.25
	return years
}

// determineMRA finds the FERS Minimum Retirement Age based on birth year.
func determineMRA(birthDate time.Time) int {
	birthYear := birthDate.Year()
	switch {
	case birthYear < 1948:
		return 55
	case birthYear == 1948:
		return 55 // 55 and 2 months
	case birthYear == 1949:
		return 55 // 55 and 4 months
	case birthYear == 1950:
		return 55 // 55 and 6 months
	case birthYear == 1951:
		return 55 // 55 and 8 months
	case birthYear == 1952:
		return 55 // 55 and 10 months
	case birthYear >= 1953 && birthYear <= 1964:
		return 56
	case birthYear == 1965:
		return 56 // 56 and 2 months
	case birthYear == 1966:
		return 56 // 56 and 4 months
	case birthYear == 1967:
		return 56 // 56 and 6 months
	case birthYear == 1968:
		return 56 // 56 and 8 months
	case birthYear == 1969:
		return 56 // 56 and 10 months
	case birthYear >= 1970:
		return 57
	default:
		return 57 // Default or handle error? Returning 57 for safety.
	}
	// NOTE: This MRA calculation returns only the whole year part.
	// OPM MRA includes months (e.g., 56 and 2 months).
	// The exact MRA with months might be needed for precise eligibility checks later.
}
