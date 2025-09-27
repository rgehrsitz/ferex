package calculation

import (
	"ferex/backend/config"
	"ferex/backend/models"
	"fmt"
	"log"
	"math"
	"time"
)

// getTSPConfig returns the TSP configuration or panics if it cannot be loaded
// This is used internally by the calculation package to ensure the config is available
func getTSPConfig() *config.TSPConfig {
	tspConfig, err := config.GetTSPConfig()
	if err != nil {
		log.Printf("Error loading TSP configuration: %v", err)
		// Return default values as fallback
		return &config.TSPConfig{
			ContributionLimits: struct {
				MaxAgencyMatchPercentage                float64 `json:"maxAgencyMatchPercentage"`
				AgencyAutomaticContributionPercentage   float64 `json:"agencyAutomaticContributionPercentage"`
				StandardEmployeeContributionForMaxMatch float64 `json:"standardEmployeeContributionForMaxMatch"`
				CatchUpContributionLimit                float64 `json:"catchUpContributionLimit"`
				MaxRegularContributionLimit             float64 `json:"maxRegularContributionLimit"`
				Year                                    int     `json:"year"`
			}{
				MaxAgencyMatchPercentage:                0.04,
				AgencyAutomaticContributionPercentage:   0.01,
				StandardEmployeeContributionForMaxMatch: 0.05,
				CatchUpContributionLimit:                7500,
				MaxRegularContributionLimit:             23000,
				Year:                                    2024,
			},
		}
	}
	return tspConfig
}

// CalculateTSP projects TSP balance growth and simulates withdrawals.
func CalculateTSP(input models.TSPCalculationInput) models.TSPCalculationResult {
	// log.Printf("Enter calculation.CalculateTSP: Received WithdrawalPercentageValue = %.8f", input.WithdrawalPercentageValue)
	result := models.TSPCalculationResult{
		Notes:              "",
		WithdrawalSchedule: []models.YearlyTSPWithdrawalDetail{},
	}

	// --- Part 1: Project TSP Growth until Retirement ---
	currentTradBalance := input.CurrentTraditionalBalance
	currentRothBalance := input.CurrentRothBalance
	yearsToGrow := input.RetirementAgeYears - input.CurrentAgeYears

	if !input.ContributeUntilRetirement && input.ContributionStopAge > input.CurrentAgeYears {
		yearsToGrow = input.ContributionStopAge - input.CurrentAgeYears
	} else if !input.ContributeUntilRetirement && input.ContributionStopAge <= input.CurrentAgeYears {
		yearsToGrow = 0 // Contributions already stopped or stop immediately
	}

	// log.Printf("Starting TSP Growth Projection: Years to Grow: %d, Current Trad: %.2f, Current Roth: %.2f", yearsToGrow, currentTradBalance, currentRothBalance)

	for i := 0; i < yearsToGrow; i++ {
		// --- Calculate Annual Contributions ---
		annualEmployeeContributionTotal := 0.0
		if input.IsContributionPercentage {
			annualEmployeeContributionTotal = input.BaseSalaryForContributions * (input.EmployeeContributionPercentage / 100.0)
		} else {
			annualEmployeeContributionTotal = input.EmployeeContributionAmount
		}

		actualAnnualEmployeeContributionTotal := annualEmployeeContributionTotal
		currentAgeInLoop := input.CurrentAgeYears + i

		// Get the current TSP configuration
		tspConfig := getTSPConfig()

		// Regular Contribution Limit
		if actualAnnualEmployeeContributionTotal > tspConfig.ContributionLimits.MaxRegularContributionLimit {
			actualAnnualEmployeeContributionTotal = tspConfig.ContributionLimits.MaxRegularContributionLimit
		}

		// Catch-Up Contributions
		potentialCatchUpAmount := 0.0
		if input.CatchUpContributionsEligible && currentAgeInLoop >= 50 {
			// If their intended contribution was already above regular limit, the difference is potential catch-up
			if annualEmployeeContributionTotal > tspConfig.ContributionLimits.MaxRegularContributionLimit {
				potentialCatchUpAmount = annualEmployeeContributionTotal - tspConfig.ContributionLimits.MaxRegularContributionLimit
			} else {
				// If their intended contribution was below or at regular limit, they can add full catch-up amount if desired
				// For simplicity, assume they intend to maximize catch-up if eligible and their base contribution doesn't already use it up.
				potentialCatchUpAmount = tspConfig.ContributionLimits.CatchUpContributionLimit // Default to max catch-up
			}

			// Cap at the catch-up limit
			if potentialCatchUpAmount > tspConfig.ContributionLimits.CatchUpContributionLimit {
				potentialCatchUpAmount = tspConfig.ContributionLimits.CatchUpContributionLimit
			}

			// Add catch-up to total contributions, but respect combined limits
			totalWithCatchUp := actualAnnualEmployeeContributionTotal + potentialCatchUpAmount
			if totalWithCatchUp > tspConfig.ContributionLimits.MaxRegularContributionLimit+tspConfig.ContributionLimits.CatchUpContributionLimit {
				totalWithCatchUp = tspConfig.ContributionLimits.MaxRegularContributionLimit + tspConfig.ContributionLimits.CatchUpContributionLimit
			}
			actualAnnualEmployeeContributionTotal = totalWithCatchUp
			if totalWithCatchUp == tspConfig.ContributionLimits.MaxRegularContributionLimit+tspConfig.ContributionLimits.CatchUpContributionLimit {
				result.Notes += fmt.Sprintf("Note: Year %d (Age %d): Total contribution (regular + catch-up) capped at %.2f. ", input.CurrentAgeYears+i, currentAgeInLoop, tspConfig.ContributionLimits.MaxRegularContributionLimit+tspConfig.ContributionLimits.CatchUpContributionLimit)
			}
		}

		employeeTradContribution := actualAnnualEmployeeContributionTotal * (input.TraditionalContributionAllocationPct / 100.0)
		employeeRothContribution := actualAnnualEmployeeContributionTotal * (input.RothContributionAllocationPct / 100.0)

		// Agency Contributions (always Traditional)
		// Get the agency contribution percentages from configuration
		agencyAutoContribPct := 0.01               // Default 1% automatic contribution
		standardEmployeeContribForMaxMatch := 0.05 // Default 5% employee contribution for max match
		maxAgencyMatchPct := 0.04                  // Default 4% maximum agency matching

		// Use configuration values if available
		if tspConfig.ContributionLimits.AgencyAutomaticContributionPercentage > 0 {
			agencyAutoContribPct = tspConfig.ContributionLimits.AgencyAutomaticContributionPercentage
		}
		if tspConfig.ContributionLimits.MaxAgencyMatchPercentage > 0 {
			maxAgencyMatchPct = tspConfig.ContributionLimits.MaxAgencyMatchPercentage
		}
		if tspConfig.ContributionLimits.StandardEmployeeContributionForMaxMatch > 0 {
			standardEmployeeContribForMaxMatch = tspConfig.ContributionLimits.StandardEmployeeContributionForMaxMatch
		}

		agencyAutomaticContribution := input.BaseSalaryForContributions * agencyAutoContribPct
		agencyMatchingContribution := 0.0
		employeeContribPctForMatch := 0.0
		if input.BaseSalaryForContributions > 0 {
			employeeContribPctForMatch = annualEmployeeContributionTotal / input.BaseSalaryForContributions
		}

		// Calculate matching contribution based on employee contribution percentage
		// TSP matching formula: 100% of first 3%, then 50% of next 2%
		if employeeContribPctForMatch >= standardEmployeeContribForMaxMatch {
			// Full match - employee contributes at least 5%
			agencyMatchingContribution = input.BaseSalaryForContributions * maxAgencyMatchPct
		} else if employeeContribPctForMatch >= 0.04 {
			// Employee contributes 4% - match is 3.5% (3% fully matched + 0.5% for the 4th percent)
			agencyMatchingContribution = input.BaseSalaryForContributions * (0.03 + (employeeContribPctForMatch-0.03)*0.5)
		} else if employeeContribPctForMatch >= 0.03 {
			// Employee contributes 3% - match is 3% (3% fully matched)
			agencyMatchingContribution = input.BaseSalaryForContributions * 0.03
		} else if employeeContribPctForMatch >= 0.02 {
			// Employee contributes 2% - match is 2%
			agencyMatchingContribution = input.BaseSalaryForContributions * 0.02
		} else if employeeContribPctForMatch >= 0.01 {
			// Employee contributes 1% - match is 1%
			agencyMatchingContribution = input.BaseSalaryForContributions * 0.01
		} else {
			// No match for 0% contribution
			agencyMatchingContribution = 0.0
		}

		totalAgencyContribution := agencyAutomaticContribution + agencyMatchingContribution

		currentTradBalance += employeeTradContribution + totalAgencyContribution
		currentRothBalance += employeeRothContribution

		// log.Printf("Year %d (Age %d): EmpCont: %.2f (T: %.2f, R: %.2f), AgencyCont: %.2f (Auto: %.2f, Match: %.2f), TradBal: %.2f, RothBal: %.2f", currentYear+i, input.CurrentAgeYears+i, annualEmployeeContributionTotal, employeeTradContribution, employeeRothContribution, totalAgencyContribution, agencyAutomaticContribution, agencyMatchingContribution, currentTradBalance, currentRothBalance)

		// --- Calculate Growth for the Year ---
		// Determine allocation for growth for the current year
		var allocationForGrowth models.TSPFundAllocationPercentages

		// Determine which allocation strategy to use for pre-retirement growth
		switch input.FutureAllocationStrategy {
		case "MatchContributionAllocation":
			allocationForGrowth = input.ContributionFundAllocation
			// If it's the first year (i=0) and strategy is MatchContributionAllocation,
			// this means the initial balance (input.CurrentAllocation) should immediately be considered
			// rebalanced to input.ContributionFundAllocation before the first year's growth.
			// If i > 0, it means rebalancing happens at the start of each subsequent year.
		case "UsePostRetirementAllocation":
			// For pre-retirement growth, we'll still use the current allocation
			// Post-retirement allocation will only be used after retirement
			allocationForGrowth = input.CurrentAllocation
		default:
			// "MaintainCurrent" or other (unspecified) strategy means we stick with input.CurrentAllocation
			// for the duration of the pre-retirement growth projection.
			allocationForGrowth = input.CurrentAllocation
		}

		// If an L-Fund is selected, update its composition for the current year based on the glide path
		if allocationForGrowth.LFundName != nil && *allocationForGrowth.LFundName != "" {
			// Calculate the projection year (current year + years into projection)
			currentYear := time.Now().Year() + i

			// Get the updated L-Fund composition for the current year
			allocationForGrowth = getLFundComposition(*allocationForGrowth.LFundName, currentYear)
		}

		weightedAvgReturn := calculateWeightedAverageReturn(allocationForGrowth, input.UserReturnAssumptions, true /* isPreRetirement */)

		tradGrowth := currentTradBalance * weightedAvgReturn
		rothGrowth := currentRothBalance * weightedAvgReturn
		currentTradBalance += tradGrowth
		currentRothBalance += rothGrowth

		// --- Deduct Expense Ratio (Placeholder) ---
		if input.ExpenseRatio > 0 {
			currentTradBalance *= (1.0 - input.ExpenseRatio)
			currentRothBalance *= (1.0 - input.ExpenseRatio)
		}

		// --- Handle L-fund rebalancing based on glide path ---
		// If an L-Fund is selected, update its composition for the next year based on the glide path
		if allocationForGrowth.LFundName != nil && *allocationForGrowth.LFundName != "" {
			// Calculate the projection year (current year + years into projection)
			currentYear := time.Now().Year() + i + 1 // +1 for next year's allocation

			// Get the updated L-Fund composition for the next year
			if input.FutureAllocationStrategy == "MaintainCurrent" {
				// For MaintainCurrent, we still need to adjust the L-Fund composition according to its glide path
				allocationForGrowth = getLFundComposition(*allocationForGrowth.LFundName, currentYear)
			}
		}
	}

	result.ProjectedTraditionalBalanceAtRetirement = currentTradBalance
	result.ProjectedRothBalanceAtRetirement = currentRothBalance
	result.TotalProjectedBalanceAtRetirement = currentTradBalance + currentRothBalance

	// --- Part 2: Simulate TSP Withdrawals Post-Retirement ---
	tradBalanceForWithdrawal := result.ProjectedTraditionalBalanceAtRetirement
	rothBalanceForWithdrawal := result.ProjectedRothBalanceAtRetirement
	totalBalanceForWithdrawal := result.TotalProjectedBalanceAtRetirement

	if input.WithdrawalStrategy != "None" && input.YearsToProjectWithdrawals > 0 && totalBalanceForWithdrawal > 0 {
		startWithdrawalAge := input.RetirementAgeYears // Default to retirement age
		if input.WithdrawalStartDate == "SpecificAge" && input.WithdrawalStartAge > input.RetirementAgeYears {
			startWithdrawalAge = input.WithdrawalStartAge
		} else if input.WithdrawalStartDate == "SpecificAge" && input.WithdrawalStartAge <= input.RetirementAgeYears {
			startWithdrawalAge = input.RetirementAgeYears // Can't start before retirement
		}

		// Determine calendar year of retirement for Year field in YearlyTSPWithdrawalDetail
		// Assuming input.CurrentAgeYears corresponds to a known calendar year (e.g., current real year)
		// This part might need adjustment if we need precise calendar years rather than just 'Year 1, Year 2' of withdrawal
		// For now, 'yearCounter' will be 1-indexed for withdrawal years.

		// Growth between retirement and withdrawal start age if withdrawal is deferred
		yearsDeferred := startWithdrawalAge - input.RetirementAgeYears
		for y := 0; y < yearsDeferred; y++ {
			allocationForPostRetirementGrowth := input.CurrentAllocation         // TODO: Allow separate post-retirement allocation strategy
			if input.FutureAllocationStrategy == "MatchContributionAllocation" { // Or a new specific post-retirement one
				allocationForPostRetirementGrowth = input.ContributionFundAllocation
			}

			// If an L-Fund is selected, update its composition for the current year based on the glide path
			if allocationForPostRetirementGrowth.LFundName != nil && *allocationForPostRetirementGrowth.LFundName != "" {
				// Calculate the projection year (current year + retirement age + years deferred)
				currentYear := time.Now().Year() + (input.RetirementAgeYears - input.CurrentAgeYears) + y

				// Get the updated L-Fund composition for the current year
				allocationForPostRetirementGrowth = getLFundComposition(*allocationForPostRetirementGrowth.LFundName, currentYear)
			}

			postRetirementReturn := calculateWeightedAverageReturn(allocationForPostRetirementGrowth, input.UserReturnAssumptions, false /*isPreRetirement*/)

			tradBalanceForWithdrawal *= (1 + postRetirementReturn)
			rothBalanceForWithdrawal *= (1 + postRetirementReturn)
			if input.ExpenseRatio > 0 {
				tradBalanceForWithdrawal *= (1.0 - input.ExpenseRatio)
				rothBalanceForWithdrawal *= (1.0 - input.ExpenseRatio)
			}
		}
		totalBalanceForWithdrawal = tradBalanceForWithdrawal + rothBalanceForWithdrawal

		for yearNum := 0; yearNum < input.YearsToProjectWithdrawals; yearNum++ {
			if tradBalanceForWithdrawal+rothBalanceForWithdrawal <= 0 {
				result.Notes += "TSP balance exhausted. "
				break
			}

			detail := models.YearlyTSPWithdrawalDetail{
				Year:                  input.RetirementAgeYears + yearsDeferred + yearNum, // This is Age, effectively
				Age:                   input.RetirementAgeYears + yearsDeferred + yearNum,
				BeginningBalanceTrad:  tradBalanceForWithdrawal,
				BeginningBalanceRoth:  rothBalanceForWithdrawal,
				BeginningBalanceTotal: tradBalanceForWithdrawal + rothBalanceForWithdrawal,
			}

			targetAnnualWithdrawal := 0.0
			switch input.WithdrawalStrategy {
			case "FixedAmountYearly":
				targetAnnualWithdrawal = input.WithdrawalFixedAmountValue
			case "FixedAmountMonthly": // Treated as annual for this model
				targetAnnualWithdrawal = input.WithdrawalFixedAmountValue * 12
			case "PercentageOfBalanceYearly":
				// log.Printf("DEBUG TSP Percentage Calc: BeginningBalanceTotal=%.2f, WithdrawalPercentageValue=%.4f", detail.BeginningBalanceTotal, input.WithdrawalPercentageValue)
				targetAnnualWithdrawal = detail.BeginningBalanceTotal * input.WithdrawalPercentageValue
			case "IRSMinimumRequiredDistribution":
				rmdFactor, rmdApplicable := getRmdFactor(input.BirthYear, detail.Age)
				if rmdApplicable && rmdFactor > 0 {
					// RMD is based on the balance at the beginning of the year (which is the EOY balance of the previous year)
					balanceForRmdCalc := detail.BeginningBalanceTotal

					// Calculate RMD by dividing the balance by the distribution factor
					targetAnnualWithdrawal = balanceForRmdCalc / rmdFactor
					detail.RmdAmount = targetAnnualWithdrawal // Store the calculated RMD amount

					// Add a note if this is the first RMD year
					if yearNum == 0 || (yearNum > 0 && !wasRmdApplicableLastYear(input.BirthYear, detail.Age-1)) {
						detail.Notes = fmt.Sprintf("First year of RMD at age %d", detail.Age)
					}
				} else if !rmdApplicable {
					// RMD not yet applicable for this age based on birth year
					targetAnnualWithdrawal = 0
					detail.RmdAmount = 0
				} else { // rmdApplicable but factor is 0 (error case)
					result.Notes += fmt.Sprintf("RMD factor not found for age %d. ", detail.Age)
					targetAnnualWithdrawal = 0
					detail.RmdAmount = 0
				}
			case "InflationAdjustedFixedAmount":
				// Calculate inflation-adjusted withdrawal amount
				baseAmount := input.WithdrawalFixedAmountValue
				inflationFactor := math.Pow(1+input.ExpectedAnnualInflationRate, float64(yearNum))
				targetAnnualWithdrawal = baseAmount * inflationFactor

				// Add note about inflation adjustment
				if yearNum > 0 {
					detail.Notes = fmt.Sprintf("Inflation-adjusted from $%.2f to $%.2f", baseAmount, targetAnnualWithdrawal)
				}
			default:
				targetAnnualWithdrawal = 0
			}

			// Ensure withdrawal doesn't exceed total balance
			if targetAnnualWithdrawal > detail.BeginningBalanceTotal {
				targetAnnualWithdrawal = detail.BeginningBalanceTotal
			}

			var withdrawnTrad, withdrawnRoth float64
			if targetAnnualWithdrawal > 0 {
				switch input.WithdrawalOrder {
				case "TraditionalFirst":
					if detail.BeginningBalanceTrad >= targetAnnualWithdrawal {
						withdrawnTrad = targetAnnualWithdrawal
						withdrawnRoth = 0
					} else {
						withdrawnTrad = detail.BeginningBalanceTrad
						withdrawnRoth = math.Min(detail.BeginningBalanceRoth, targetAnnualWithdrawal-withdrawnTrad)
					}
				case "RothFirst":
					if detail.BeginningBalanceRoth >= targetAnnualWithdrawal {
						withdrawnTrad = 0
						withdrawnRoth = targetAnnualWithdrawal
					} else {
						withdrawnRoth = detail.BeginningBalanceRoth
						withdrawnTrad = math.Min(detail.BeginningBalanceTrad, targetAnnualWithdrawal-withdrawnRoth)
					}
				case "ProRata":
					// ProRata withdraws from Traditional and Roth accounts proportionally based on their balances
					if detail.BeginningBalanceTotal > 0 { // Avoid division by zero
						// Calculate the proportion of each account type relative to the total balance
						tradProportion := detail.BeginningBalanceTrad / detail.BeginningBalanceTotal
						rothProportion := detail.BeginningBalanceRoth / detail.BeginningBalanceTotal

						// Calculate withdrawal amounts based on proportions, ensuring we don't withdraw more than available
						withdrawnTrad = math.Min(detail.BeginningBalanceTrad, targetAnnualWithdrawal*tradProportion)
						withdrawnRoth = math.Min(detail.BeginningBalanceRoth, targetAnnualWithdrawal*rothProportion)

						// Adjust if the sum of prorated amounts doesn't perfectly match the target
						// This can happen due to rounding or if one account is depleted
						if withdrawnTrad+withdrawnRoth < targetAnnualWithdrawal {
							remainingNeeded := targetAnnualWithdrawal - (withdrawnTrad + withdrawnRoth)

							// Check if we need to take additional funds from one account after the other is depleted
							if detail.BeginningBalanceTrad-withdrawnTrad > 0 && detail.BeginningBalanceRoth-withdrawnRoth <= 0 {
								// Roth is empty but Traditional has funds, take the remainder from Traditional
								canTakeMoreFromTrad := detail.BeginningBalanceTrad - withdrawnTrad
								additionalFromTrad := math.Min(remainingNeeded, canTakeMoreFromTrad)
								withdrawnTrad += additionalFromTrad
								detail.Notes += fmt.Sprintf(" Additional $%.2f taken from Traditional after Roth depletion.", additionalFromTrad)
							} else if detail.BeginningBalanceRoth-withdrawnRoth > 0 && detail.BeginningBalanceTrad-withdrawnTrad <= 0 {
								// Traditional is empty but Roth has funds, take the remainder from Roth
								canTakeMoreFromRoth := detail.BeginningBalanceRoth - withdrawnRoth
								additionalFromRoth := math.Min(remainingNeeded, canTakeMoreFromRoth)
								withdrawnRoth += additionalFromRoth
								detail.Notes += fmt.Sprintf(" Additional $%.2f taken from Roth after Traditional depletion.", additionalFromRoth)
							}
						}
					} else { // Total balance is zero, can't withdraw
						withdrawnTrad = 0
						withdrawnRoth = 0
						detail.Notes = "No funds available for withdrawal"
					}
				default: // Default to ProRata or consider it an error/note
					result.Notes += "Unknown withdrawal order, defaulting to ProRata. "
					if detail.BeginningBalanceTotal > 0 {
						tradProportion := detail.BeginningBalanceTrad / detail.BeginningBalanceTotal
						rothProportion := detail.BeginningBalanceRoth / detail.BeginningBalanceTotal
						withdrawnTrad = math.Min(detail.BeginningBalanceTrad, targetAnnualWithdrawal*tradProportion)
						withdrawnRoth = math.Min(detail.BeginningBalanceRoth, targetAnnualWithdrawal*rothProportion)
					}
				}
			}
			detail.TraditionalWithdrawn = withdrawnTrad
			detail.RothWithdrawn = withdrawnRoth
			detail.TotalWithdrawn = withdrawnTrad + withdrawnRoth

			// CORRECTED ORDER OF OPERATIONS:
			// Calculate growth BEFORE withdrawing money (growth happens on full beginning balance)

			// Calculate growth for the year based on allocation and withdrawal strategy
			allocationForWithdrawalGrowth := input.CurrentAllocation // Default to current allocation

			// Determine which allocation strategy to use for post-retirement growth
			switch input.FutureAllocationStrategy {
			case "MatchContributionAllocation":
				allocationForWithdrawalGrowth = input.ContributionFundAllocation
			case "UsePostRetirementAllocation":
				// Use the specified post-retirement allocation if provided
				if input.PostRetirementAllocation.G != nil ||
					input.PostRetirementAllocation.F != nil ||
					input.PostRetirementAllocation.C != nil ||
					input.PostRetirementAllocation.S != nil ||
					input.PostRetirementAllocation.I != nil ||
					input.PostRetirementAllocation.LFundName != nil {
					allocationForWithdrawalGrowth = input.PostRetirementAllocation
					detail.Notes += "Using post-retirement allocation for growth calculations. "
				}
			default:
				// MaintainCurrent or any other value defaults to current allocation
				allocationForWithdrawalGrowth = input.CurrentAllocation
			}

			// If an L-Fund is selected, update its composition for the current year based on the glide path
			if allocationForWithdrawalGrowth.LFundName != nil && *allocationForWithdrawalGrowth.LFundName != "" {
				// Calculate the projection year (current year + retirement age + years deferred + years into withdrawal)
				currentYear := time.Now().Year() + (input.RetirementAgeYears - input.CurrentAgeYears) + yearsDeferred + yearNum

				// Get the updated L-Fund composition for the current year
				allocationForWithdrawalGrowth = getLFundComposition(*allocationForWithdrawalGrowth.LFundName, currentYear)
			}

			postRetirementReturn := calculateWeightedAverageReturn(allocationForWithdrawalGrowth, input.UserReturnAssumptions, false /*isPostRetirement*/)

			// Calculate growth on the FULL beginning balance (before withdrawal)
			tradGrowth := tradBalanceForWithdrawal * postRetirementReturn
			rothGrowth := rothBalanceForWithdrawal * postRetirementReturn
			detail.GrowthAmount = tradGrowth + rothGrowth

			// Apply growth first
			tradBalanceForWithdrawal += tradGrowth
			rothBalanceForWithdrawal += rothGrowth

			// THEN subtract withdrawals
			tradBalanceForWithdrawal -= withdrawnTrad
			rothBalanceForWithdrawal -= withdrawnRoth

			// Finally, deduct expense ratio from remaining balance
			if input.ExpenseRatio > 0 {
				expenseDeduction := (tradBalanceForWithdrawal + rothBalanceForWithdrawal) * input.ExpenseRatio
				tradBalanceForWithdrawal *= (1.0 - input.ExpenseRatio)
				rothBalanceForWithdrawal *= (1.0 - input.ExpenseRatio)
				detail.GrowthAmount -= expenseDeduction // Adjust growth amount to reflect expense impact
			}

			detail.EndingBalanceTrad = tradBalanceForWithdrawal
			detail.EndingBalanceRoth = rothBalanceForWithdrawal
			detail.EndingBalanceTotal = tradBalanceForWithdrawal + rothBalanceForWithdrawal

			result.WithdrawalSchedule = append(result.WithdrawalSchedule, detail)
		}
	}

	if input.WithdrawalStrategy == "None" {
		result.Notes += "No TSP withdrawals selected. "
	} else if len(result.WithdrawalSchedule) == 0 && totalBalanceForWithdrawal > 0 {
		result.Notes += "TSP withdrawal conditions not met (e.g. zero years to project or zero balance). "
	}

	// Clear placeholder note if we have some results or a 'None' strategy
	if result.Notes == "TSP calculation is a simplified placeholder." && (input.WithdrawalStrategy == "None" || len(result.WithdrawalSchedule) > 0) {
		result.Notes = ""
	} else if result.Notes == "TSP calculation is a simplified placeholder." {
		// If it's still the placeholder, but we did some work, append rather than replace
		result.Notes = "Initial TSP projection complete. " + result.Notes
	}
	return result
}

// float64Ptr returns a pointer to a float64 value.
func float64Ptr(v float64) *float64 {
	return &v
}

// getLFundBaseCompositions returns the L Fund base compositions from the configuration
func getLFundBaseCompositions() map[string]models.TSPFundAllocationPercentages {
	tspConfig := getTSPConfig()

	// Convert the JSON config to the expected format with pointers
	result := make(map[string]models.TSPFundAllocationPercentages)
	for fundName, composition := range tspConfig.LFundBaseCompositions {
		// Create a new composition with pointers
		newComp := models.TSPFundAllocationPercentages{
			LFundName: &fundName,
		}

		// Convert each value to a pointer
		if composition.G != nil {
			val := *composition.G
			newComp.G = &val
		}
		if composition.F != nil {
			val := *composition.F
			newComp.F = &val
		}
		if composition.C != nil {
			val := *composition.C
			newComp.C = &val
		}
		if composition.S != nil {
			val := *composition.S
			newComp.S = &val
		}
		if composition.I != nil {
			val := *composition.I
			newComp.I = &val
		}

		result[fundName] = newComp
	}

	// If the config is empty or invalid, return default values
	if len(result) == 0 {
		// Default fallback values
		return map[string]models.TSPFundAllocationPercentages{
			"LIncome": {G: float64Ptr(71.22), F: float64Ptr(6.51), C: float64Ptr(11.72), S: float64Ptr(2.81), I: float64Ptr(7.74)},
			"L2025":   {G: float64Ptr(69.32), F: float64Ptr(6.64), C: float64Ptr(12.71), S: float64Ptr(3.17), I: float64Ptr(8.16)},
			"L2030":   {G: float64Ptr(42.00), F: float64Ptr(5.00), C: float64Ptr(27.00), S: float64Ptr(9.00), I: float64Ptr(17.00)},
			"L2035":   {G: float64Ptr(34.00), F: float64Ptr(5.00), C: float64Ptr(31.00), S: float64Ptr(11.00), I: float64Ptr(19.00)},
			"L2040":   {G: float64Ptr(26.00), F: float64Ptr(5.00), C: float64Ptr(35.00), S: float64Ptr(13.00), I: float64Ptr(21.00)},
			"L2045":   {G: float64Ptr(18.00), F: float64Ptr(5.00), C: float64Ptr(39.00), S: float64Ptr(15.00), I: float64Ptr(23.00)},
			"L2050":   {G: float64Ptr(10.00), F: float64Ptr(5.00), C: float64Ptr(43.00), S: float64Ptr(17.00), I: float64Ptr(25.00)},
			"L2055":   {G: float64Ptr(6.00), F: float64Ptr(4.00), C: float64Ptr(47.00), S: float64Ptr(19.00), I: float64Ptr(24.00)},
			"L2060":   {G: float64Ptr(6.00), F: float64Ptr(4.00), C: float64Ptr(47.00), S: float64Ptr(19.00), I: float64Ptr(24.00)},
			"L2065":   {G: float64Ptr(6.00), F: float64Ptr(4.00), C: float64Ptr(47.00), S: float64Ptr(19.00), I: float64Ptr(24.00)},
		}
	}

	return result
}

// getLFundTargetComposition returns the L Fund target composition from the configuration
func getLFundTargetComposition() models.TSPFundAllocationPercentages {
	tspConfig := getTSPConfig()

	// Convert the JSON config to the expected format with pointers
	result := models.TSPFundAllocationPercentages{}

	// Convert each value to a pointer
	if tspConfig.LFundTargetComposition.G != nil {
		val := *tspConfig.LFundTargetComposition.G
		result.G = &val
	}
	if tspConfig.LFundTargetComposition.F != nil {
		val := *tspConfig.LFundTargetComposition.F
		result.F = &val
	}
	if tspConfig.LFundTargetComposition.C != nil {
		val := *tspConfig.LFundTargetComposition.C
		result.C = &val
	}
	if tspConfig.LFundTargetComposition.S != nil {
		val := *tspConfig.LFundTargetComposition.S
		result.S = &val
	}
	if tspConfig.LFundTargetComposition.I != nil {
		val := *tspConfig.LFundTargetComposition.I
		result.I = &val
	}

	// If the config is empty or invalid, return default values
	if result.G == nil && result.F == nil && result.C == nil && result.S == nil && result.I == nil {
		// Default fallback to LIncome values
		return models.TSPFundAllocationPercentages{
			G: float64Ptr(71.22), F: float64Ptr(6.51), C: float64Ptr(11.72), S: float64Ptr(2.81), I: float64Ptr(7.74),
		}
	}

	return result
}

// getLFundTargetYear extracts the target year from an L Fund name
// Returns 0 for LIncome or invalid fund names
func getLFundTargetYear(lFundName string) int {
	if lFundName == "LIncome" {
		return 0 // LIncome has no target year
	}

	// Extract year from L20XX format
	var targetYear int
	if _, err := fmt.Sscanf(lFundName, "L%d", &targetYear); err != nil {
		return 0 // Invalid format
	}

	return targetYear
}

// calculateLFundComposition determines the L Fund composition for a specific projection year
// based on the glide path approach, where funds gradually shift to more conservative allocations
// as they approach their target date.
func calculateLFundComposition(lFundName string, currentYear int) models.TSPFundAllocationPercentages {
	// Get the L-Fund compositions from configuration
	lFundBaseCompositions := getLFundBaseCompositions()
	lFundTargetComposition := getLFundTargetComposition()

	// For LIncome, return the fixed composition (already at target)
	if lFundName == "LIncome" {
		return lFundBaseCompositions["LIncome"]
	}

	// Get the target year for this L Fund
	targetYear := getLFundTargetYear(lFundName)
	if targetYear == 0 || targetYear <= currentYear {
		// If we've reached or passed the target year, return the target composition (LIncome)
		return lFundTargetComposition
	}

	// Get the base composition for this L Fund
	baseComposition, ok := lFundBaseCompositions[lFundName]
	if !ok {
		// If the fund name is not found, return the LIncome composition as a fallback
		return lFundTargetComposition
	}

	// Calculate the percentage of the glide path that has been completed
	// We assume the glide path starts 25 years before the target date
	glidePathStartYear := targetYear - 25
	if currentYear <= glidePathStartYear {
		// If we're before the glide path starts, return the base composition
		return baseComposition
	}

	// Calculate how far along the glide path we are (0.0 to 1.0)
	glidePathProgress := float64(currentYear-glidePathStartYear) / float64(targetYear-glidePathStartYear)
	if glidePathProgress > 1.0 {
		glidePathProgress = 1.0
	}

	// Interpolate between the base composition and the target composition
	result := models.TSPFundAllocationPercentages{
		LFundName: baseComposition.LFundName,
	}

	// Helper function to interpolate between two values
	interpolate := func(start, end *float64, progress float64) *float64 {
		if start == nil || end == nil {
			return nil
		}
		value := *start + ((*end - *start) * progress)
		return &value
	}

	// Interpolate each fund allocation
	result.G = interpolate(baseComposition.G, lFundTargetComposition.G, glidePathProgress)
	result.F = interpolate(baseComposition.F, lFundTargetComposition.F, glidePathProgress)
	result.C = interpolate(baseComposition.C, lFundTargetComposition.C, glidePathProgress)
	result.S = interpolate(baseComposition.S, lFundTargetComposition.S, glidePathProgress)
	result.I = interpolate(baseComposition.I, lFundTargetComposition.I, glidePathProgress)

	return result
}

// getLFundComposition returns the L Fund composition for a given fund name and year
// This is the main function that should be used to get L Fund compositions
func getLFundComposition(lFundName string, currentYear int) models.TSPFundAllocationPercentages {
	// Get the L-Fund compositions from configuration
	lFundBaseCompositions := getLFundBaseCompositions()

	// If the current year is 0 or not provided, use the base composition (for backward compatibility)
	if currentYear <= 0 {
		if comp, ok := lFundBaseCompositions[lFundName]; ok {
			return comp
		}
		// If not found, return LIncome as fallback
		return lFundBaseCompositions["LIncome"]
	}

	// Calculate the dynamic composition based on the glide path
	return calculateLFundComposition(lFundName, currentYear)
}

// calculateWeightedAverageReturn calculates the blended return rate.
// The allocation parameter should already have L-Fund compositions resolved if applicable.
func calculateWeightedAverageReturn(allocation models.TSPFundAllocationPercentages, returns models.TSPReturnAssumptions, isPreRetirement bool) float64 {
	// Helper to safely get float value from pointer, defaulting to 0 if nil
	getVal := func(p *float64) float64 {
		if p != nil {
			return *p
		}
		return 0.0
	}

	// Handle overall return override
	overallReturnVal := getVal(returns.Overall)
	if isPreRetirement && returns.UseOverallForPre && overallReturnVal != 0 {
		return overallReturnVal
	}
	if !isPreRetirement && returns.UseOverallForPost && overallReturnVal != 0 {
		return overallReturnVal
	}

	actualAllocation := allocation // This will be a copy

	// If LFundName is specified, use its composition
	if allocation.LFundName != nil && *allocation.LFundName != "" {
		lFundNameValue := *allocation.LFundName

		// Get the current year for dynamic L-Fund composition calculation
		currentYear := time.Now().Year()

		// Use the dynamic L-Fund composition based on the current year
		actualAllocation = getLFundComposition(lFundNameValue, currentYear)

		// If the L-Fund name is not recognized, use fallbacks
		if actualAllocation.G == nil && actualAllocation.F == nil && actualAllocation.C == nil &&
			actualAllocation.S == nil && actualAllocation.I == nil {
			// Default to overall if available, then G fund return, otherwise 0.
			if overallReturnVal != 0 { // Use already fetched overallReturnVal
				return overallReturnVal
			}
			return getVal(returns.G) // Default to G fund return
		}
	}

	// Calculate weighted average based on individual fund allocations and returns
	calculatedReturn := 0.0
	totalWeight := 0.0

	gAlloc := getVal(actualAllocation.G)
	fAlloc := getVal(actualAllocation.F)
	cAlloc := getVal(actualAllocation.C)
	sAlloc := getVal(actualAllocation.S)
	iAlloc := getVal(actualAllocation.I)

	gReturn := getVal(returns.G)
	fReturn := getVal(returns.F)
	cReturn := getVal(returns.C)
	sReturn := getVal(returns.S)
	iReturn := getVal(returns.I)

	if gAlloc > 0 {
		calculatedReturn += (gAlloc / 100.0) * gReturn
		totalWeight += gAlloc
	}
	if fAlloc > 0 {
		calculatedReturn += (fAlloc / 100.0) * fReturn
		totalWeight += fAlloc
	}
	if cAlloc > 0 {
		calculatedReturn += (cAlloc / 100.0) * cReturn
		totalWeight += cAlloc
	}
	if sAlloc > 0 {
		calculatedReturn += (sAlloc / 100.0) * sReturn
		totalWeight += sAlloc
	}
	if iAlloc > 0 {
		calculatedReturn += (iAlloc / 100.0) * iReturn
		totalWeight += iAlloc
	}

	if totalWeight > 0 {
		// The `calculatedReturn` already represents the sum of (percentage_allocation * rate_of_return) for each fund.
		// E.g., if G is 50% of portfolio and returns 2%, it contributes (0.50 * 2.0) = 1.0 to the overall portfolio return.
		// So, `calculatedReturn` is the final blended rate.
		return calculatedReturn
	}

	// If no allocation weights (totalWeight is 0), or all returns/allocations are zero/nil,
	// default to overall return if available and non-zero, otherwise G fund's return, or finally a reasonable default.
	if overallReturnVal != 0 {
		return overallReturnVal
	}

	gReturnFallback := getVal(returns.G)
	if gReturnFallback != 0 {
		return gReturnFallback // Use G fund return if available
	}

	// CRITICAL FIX: If no returns are specified at all, use a reasonable default
	// This prevents the 0% growth + 5% expense ratio = -5% annual decline issue
	log.Printf("WARNING: No TSP return assumptions provided. Using default 5%% annual return.")
	return 0.05 // Default to 5% annual return when no assumptions are provided
}

// getUniformLifetimeTable returns the RMD factors from the configuration
func getUniformLifetimeTable() map[int]float64 {
	tspConfig := getTSPConfig()
	return tspConfig.GetUniformLifetimeTableAsMap()
}

// getRmdStartAge returns the RMD start age for a given birth year from the configuration
func getRmdStartAge(birthYear int) int {
	tspConfig := getTSPConfig()
	return tspConfig.GetRMDStartAge(birthYear)
}

// wasRmdApplicableLastYear checks if RMD was applicable for the previous year
// This is used to determine if the current year is the first RMD year
func wasRmdApplicableLastYear(birthYear int, previousAge int) bool {
	_, applicable := getRmdFactor(birthYear, previousAge)
	return applicable
}

func getRmdFactor(birthYear int, currentAge int) (factor float64, applicable bool) {
	// Get the RMD start age for this birth year from the configuration
	statutoryRmdStartAge := getRmdStartAge(birthYear)
	if statutoryRmdStartAge == 0 {
		// If no rule matches, default to a conservative approach (highest age)
		statutoryRmdStartAge = 75
	}

	if currentAge < statutoryRmdStartAge {
		return 0, false // RMD not applicable yet
	}

	// Get the uniform lifetime table from the configuration
	uniformLifetimeTable := getUniformLifetimeTable()

	if factor, ok := uniformLifetimeTable[currentAge]; ok {
		return factor, true
	}

	// Handle ages outside the explicit table range but at or above statutoryRmdStartAge
	if currentAge >= statutoryRmdStartAge {
		// Find the lowest age in the table
		lowestAge := 120 // Start with a high value
		for age := range uniformLifetimeTable {
			if age < lowestAge {
				lowestAge = age
			}
		}

		if currentAge < lowestAge {
			// If current age is below the lowest age in the table but RMD is applicable
			if f, exists := uniformLifetimeTable[lowestAge]; exists {
				return f, true // Default to lowest factor in table
			}
		}

		// Find the highest age in the table
		highestAge := 0
		for age := range uniformLifetimeTable {
			if age > highestAge {
				highestAge = age
			}
		}

		if currentAge > highestAge {
			// If current age is above the highest age in the table
			if f, exists := uniformLifetimeTable[highestAge]; exists {
				return f, true // Default to highest factor in table
			}
		}
	}

	// If RMD is applicable but age is not in table
	return 0, true // Indicates RMD is applicable, but factor couldn't be found (effectively an error for calculation)
}

// Helper function to calculate agency contributions (to be implemented)
// func calculateAgencyContributions(employeeContributionTotal float64, baseSalary float64) (automatic float64, matching float64) {
// 	 // ... logic ...
// 	 return 0.0, 0.0
// }
