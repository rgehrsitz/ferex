package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"time"
	"ferex/backend/models"
	"ferex/backend/calculation"
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

// PensionInput is a minimal struct for frontend pension calculation requests
// (Add more fields as needed for future expansion)
type PensionInput struct {
	System                string  `json:"system"`
	High3Salary           float64 `json:"high3Salary"`
	YearsOfService        float64 `json:"yearsOfService"`
	AgeAtRetirement       int     `json:"ageAtRetirement"`
	UnusedSickLeaveMonths int     `json:"unusedSickLeaveMonths"`
	SurvivorBenefitOption string  `json:"survivorBenefitOption"`
	IsPartTime            bool    `json:"isPartTime"`
	PartTimeProrationFactor float64 `json:"partTimeProrationFactor"`
}

// PensionResult is a minimal struct for frontend display
// (Add more fields as needed for future expansion)
type PensionResult struct {
	AnnualPension  float64 `json:"annualPension"`
	MonthlyPension float64 `json:"monthlyPension"`
	Notes          string  `json:"notes"`
}

// SocialSecurityInput contains data for social security calculations
type SocialSecurityInput struct {
	StartAge                int     `json:"startAge"`
	EstimatedMonthlyBenefit float64 `json:"estimatedMonthlyBenefit"`
	IsEligible              bool    `json:"isEligible"`
	BirthYear               int     `json:"birthYear"`
	BirthMonth              int     `json:"birthMonth"`
	CurrentAge              int     `json:"currentAge"`
	EstimatedAnnualSalary   float64 `json:"estimatedAnnualSalary"`
	YearsWorked             int     `json:"yearsWorked"`
	UserProvidedEstimate62  float64 `json:"userProvidedEstimate62"`
	UserProvidedEstimateFRA float64 `json:"userProvidedEstimateFRA"`
	UserProvidedEstimate70  float64 `json:"userProvidedEstimate70"`
}

// SocialSecurityResult contains the calculated Social Security benefit amounts
type SocialSecurityResult struct {
	EstimatedMonthlyAt62  float64 `json:"estimatedMonthlyAt62"`
	EstimatedMonthlyAtFRA float64 `json:"estimatedMonthlyAtFRA"`
	EstimatedMonthlyAt70  float64 `json:"estimatedMonthlyAt70"`
	EstimatedAnnualAt62   float64 `json:"estimatedAnnualAt62"`
	EstimatedAnnualAtFRA  float64 `json:"estimatedAnnualAtFRA"`
	EstimatedAnnualAt70   float64 `json:"estimatedAnnualAt70"`
	ClaimingAge           int     `json:"claimingAge"`
	ClaimingMonthlyAmount float64 `json:"claimingMonthlyAmount"`
	ClaimingAnnualAmount  float64 `json:"claimingAnnualAmount"`
	FullRetirementAge     int     `json:"fullRetirementAge"`
	Notes                 string  `json:"notes"`
}

// TSPInput contains data for thrift savings plan calculations
type TSPInput struct {
	CurrentBalance        float64 `json:"currentBalance"`
	TraditionalBalance    float64 `json:"traditionalBalance"`
	RothBalance           float64 `json:"rothBalance"`
	AnnualContribution    float64 `json:"annualContribution"`
	ExpectedReturnRate    float64 `json:"expectedReturnRate"`
	WithdrawalStrategy    string  `json:"withdrawalStrategy"`
	FixedWithdrawalAmount float64 `json:"fixedWithdrawalAmount"`
	WithdrawalPercentage  float64 `json:"withdrawalPercentage"`
	WithdrawalStartAge    int     `json:"withdrawalStartAge"`
	BirthYear             int     `json:"birthYear"`
	BirthMonth            int     `json:"birthMonth"`
	RetirementAge         int     `json:"retirementAge"`
}

// TSPYearData represents a single year in the TSP projection
type TSPYearData struct {
	Age              int     `json:"age"`
	Year             int     `json:"year"`
	StartingBalance  float64 `json:"startingBalance"`
	Contributions    float64 `json:"contributions"`
	Returns          float64 `json:"returns"`
	Withdrawals      float64 `json:"withdrawals"`
	EndingBalance    float64 `json:"endingBalance"`
}

// TSPProjectionResult contains the projected TSP growth and withdrawals
type TSPProjectionResult struct {
	YearlyData         []TSPYearData `json:"yearlyData"`
	MaxBalance         float64       `json:"maxBalance"`
	FinalBalance       float64       `json:"finalBalance"`
	TotalContributions float64       `json:"totalContributions"`
	TotalWithdrawals   float64       `json:"totalWithdrawals"`
	TotalReturns       float64       `json:"totalReturns"`
	Notes              string        `json:"notes"`
}

// TaxInput contains data for tax calculations
type TaxInput struct {
	FilingStatus          string  `json:"filingStatus"`
	StateOfResidence      string  `json:"stateOfResidence"`
	StateIncomeTaxRate    float64 `json:"stateIncomeTaxRate"`
	TotalIncome           float64 `json:"totalIncome"`
	PensionIncome         float64 `json:"pensionIncome"`
	SocialSecurityIncome  float64 `json:"socialSecurityIncome"`
	TSPWithdrawals        float64 `json:"tspWithdrawals"`
	OtherIncome           float64 `json:"otherIncome"`
	NonTaxableIncome      float64 `json:"nonTaxableIncome"`
	ItemizedDeductions    float64 `json:"itemizedDeductions"`
	FederalTaxCredits     float64 `json:"federalTaxCredits"`
	StateTaxCredits       float64 `json:"stateTaxCredits"`
	Age                   int     `json:"age"`
	SpouseAge             int     `json:"spouseAge"`
}

// TaxResult contains the calculated tax amounts and rates
type TaxResult struct {
	FederalTax          float64 `json:"federalTax"`
	StateTax            float64 `json:"stateTax"`
	TotalTax            float64 `json:"totalTax"`
	EffectiveFederalRate float64 `json:"effectiveFederalRate"`
	EffectiveStateRate  float64 `json:"effectiveStateRate"`
	EffectiveTotalRate  float64 `json:"effectiveTotalRate"`
	Notes               string  `json:"notes"`
}

// COLAInput contains data for cost of living adjustment calculations
type COLAInput struct {
	AssumedInflationRate      float64 `json:"assumedInflationRate"`
	ApplyCOLAToPension        bool    `json:"applyCOLAToPension"`
	ApplyColaToSocialSecurity bool    `json:"applyColaToSocialSecurity"`
	BaseAmount                float64 `json:"baseAmount"`
	RetirementSystem          string  `json:"retirementSystem"`
	RetirementAge             int     `json:"retirementAge"`
	IsSpecialProvision        bool    `json:"isSpecialProvision"`
	StartYear                 int     `json:"startYear"`
	ProjectionYears           int     `json:"projectionYears"`
	MonthsInFirstYear         int     `json:"monthsInFirstYear"`
}

// COLAYearData represents a single year in the COLA projection
type COLAYearData struct {
	Year              int     `json:"year"`
	InflationRate     float64 `json:"inflationRate"`
	StartingAmount    float64 `json:"startingAmount"`
	COLAPercentage    float64 `json:"colaPercentage"`
	AdjustedAmount    float64 `json:"adjustedAmount"`
	CumulativeGrowth  float64 `json:"cumulativeGrowth"`
}

// COLAResult contains the projected cost of living adjustments
type COLAResult struct {
	BaseAmount        float64        `json:"baseAmount"`
	FinalAmount       float64        `json:"finalAmount"`
	TotalGrowthRate   float64        `json:"totalGrowthRate"`
	EffectiveAnnualRate float64      `json:"effectiveAnnualRate"`
	YearlyAdjustments []COLAYearData `json:"yearlyAdjustments"`
	Notes             string         `json:"notes"`
}

// OtherIncomeSource represents an additional income stream
type OtherIncomeSource struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Amount     float64 `json:"amount"`
	Frequency  string  `json:"frequency"`
	StartAge   int     `json:"startAge"`
	EndAge     *int    `json:"endAge"`
	ApplyCola  bool    `json:"applyCola"`
}

// OtherIncomeInput contains additional income sources
type OtherIncomeInput struct {
	Sources []OtherIncomeSource `json:"sources"`
}

// RetirementScenarioInput combines all retirement income components
type RetirementScenarioInput struct {
	Pension        PensionInput       `json:"pension"`
	SocialSecurity SocialSecurityInput `json:"socialSecurity"`
	TSP            TSPInput           `json:"tsp"`
	Tax            TaxInput           `json:"tax"`
	COLA           COLAInput          `json:"cola"`
	OtherIncome    OtherIncomeInput   `json:"otherIncome"`
	ProjectionStartAge int             `json:"projectionStartAge"`
	ProjectionEndAge   int             `json:"projectionEndAge"`
}

// YearlyProjectionData contains calculated values for a specific year in retirement
type YearlyProjectionData struct {
	Age              int     `json:"age"`
	Year             int     `json:"year"`
	PensionIncome    float64 `json:"pensionIncome"`
	SocialSecurity   float64 `json:"socialSecurity"`
	TSPWithdrawal    float64 `json:"tspWithdrawal"`
	OtherIncome      float64 `json:"otherIncome"`
	TotalGrossIncome float64 `json:"totalGrossIncome"`
	FederalTax       float64 `json:"federalTax"`
	StateTax         float64 `json:"stateTax"`
	TotalTaxes       float64 `json:"totalTaxes"`
	NetIncome        float64 `json:"netIncome"`
	TSPBalance       float64 `json:"tspBalance"`
}

// RetirementProjectionResult contains the complete projection for a retirement scenario
type RetirementProjectionResult struct {
	YearlyData       []YearlyProjectionData `json:"yearlyData"`
	TotalGrossIncome float64                `json:"totalGrossIncome"`
	TotalNetIncome   float64                `json:"totalNetIncome"`
	TotalTaxes       float64                `json:"totalTaxes"`
	MaxTSPBalance    float64                `json:"maxTSPBalance"`
	Notes            string                 `json:"notes"`
}

// CalculatePension computes a FERS or CSRS pension based on user input
//export
func (a *App) CalculatePension(input PensionInput) PensionResult {
	if input.System == "FERS" {
		// Check if qualifies for 1.1% multiplier (age 62+ with 20+ years)
		isAge62With20Years := input.AgeAtRetirement >= 62 && input.YearsOfService >= 20
		
		// Map the survivor benefit option to the model format
		var survivorBenefitOption string
		switch input.SurvivorBenefitOption {
		case "full":
			survivorBenefitOption = "max"
		case "partial":
			survivorBenefitOption = "partial"
		case "none":
			survivorBenefitOption = "none"
		default:
			survivorBenefitOption = "none"
		}
		
		fersInput := models.FERSCalculationInput{
			High3Salary:             input.High3Salary,
			YearsOfService:          input.YearsOfService,
			AgeAtRetirement:         input.AgeAtRetirement,
			UnusedSickLeaveMonths:   input.UnusedSickLeaveMonths,
			IsAge62With20Years:      isAge62With20Years,
			SurvivorBenefitElection: survivorBenefitOption,
			IsPartTime:              input.IsPartTime,
			PartTimeProrationFactor: input.PartTimeProrationFactor,
		}
		
		fersResult := calculation.CalculateFERSPension(fersInput)
		return PensionResult{
			AnnualPension:  fersResult.AnnualPension,
			MonthlyPension: fersResult.MonthlyPension,
			Notes:          fersResult.Notes,
		}
	} else if input.System == "CSRS" || input.System == "CSRS Offset" {
		// Map the survivor benefit option to the model format
		var survivorBenefitOption string
		switch input.SurvivorBenefitOption {
		case "full":
			survivorBenefitOption = "max"
		case "partial":
			survivorBenefitOption = "partial"
		case "none":
			survivorBenefitOption = "none"
		default:
			survivorBenefitOption = "none"
		}
		
		csrsInput := models.CSRSCalculationInput{
			High3Salary:             input.High3Salary,
			YearsOfService:          input.YearsOfService,
			UnusedSickLeaveMonths:   input.UnusedSickLeaveMonths,
			IsCSRSOffset:            input.System == "CSRS Offset",
			AgeAtRetirement:         input.AgeAtRetirement,
			SurvivorBenefitElection: survivorBenefitOption,
			IsPartTime:              input.IsPartTime,
			PartTimeProrationFactor: input.PartTimeProrationFactor,
		}
		
		csrsResult := calculation.CalculateCSRS(csrsInput)
		return PensionResult{
			AnnualPension:  csrsResult.AnnualPension,
			MonthlyPension: csrsResult.MonthlyPension,
			Notes:          csrsResult.Notes,
		}
	}
	return PensionResult{Notes: "Unknown retirement system."}
}

// CalculateRetirementProjection generates a complete retirement income projection
//export
func (a *App) CalculateRetirementProjection(input RetirementScenarioInput) RetirementProjectionResult {
	// Initialize result
	result := RetirementProjectionResult{
		Notes: "",
	}
	
	// Calculate pension
	pensionResult := a.CalculatePension(input.Pension)
	
	// Get birth year and month for age calculations
	birthYear := input.SocialSecurity.BirthYear
	birthMonth := input.SocialSecurity.BirthMonth
	
	// Calculate the current date and user's age more accurately
	now := time.Now()
	currentYear := now.Year()
	currentMonth := int(now.Month())
	
	// Create annual projections
	var yearlyData []YearlyProjectionData
	var cumulativeGrossIncome, cumulativeTaxes, cumulativeNetIncome float64
	maxTSPBalance := input.TSP.CurrentBalance
	
	// Determine age range for projection
	startAge := input.ProjectionStartAge
	if startAge == 0 {
		startAge = input.Pension.AgeAtRetirement
	}
	
	endAge := input.ProjectionEndAge
	if endAge == 0 {
		endAge = 95 // Default to age 95 if not specified
	}
	
	// Current TSP balance at retirement start
	currentTSPBalance := input.TSP.CurrentBalance
	
	// Calculate more accurate starting year based on birth date
	startYear := currentYear
	
	// Calculate age more precisely (important for determining when withdrawals start)
	userCurrentAge := currentYear - birthYear
	if birthMonth > currentMonth {
		userCurrentAge--
	}
	
	// Calculate projections for each year
	for age := startAge; age <= endAge; age++ {
		// Calculate the correct year based on current age and projection age
		year := startYear + (age - userCurrentAge)
		
		// Create yearly data entry
		yearData := YearlyProjectionData{
			Age:  age,
			Year: year,
		}
		
		// Calculate pension income (with COLA if applicable)
		pensionIncome := pensionResult.AnnualPension
		if input.COLA.ApplyCOLAToPension && age > input.Pension.AgeAtRetirement {
			yearsSinceRetirement := age - input.Pension.AgeAtRetirement
			inflationFactor := math.Pow(1+input.COLA.AssumedInflationRate, float64(yearsSinceRetirement))
			pensionIncome *= inflationFactor
		}
		yearData.PensionIncome = pensionIncome
		
		// Calculate Social Security income
		ssIncome := 0.0
		if input.SocialSecurity.IsEligible && age >= input.SocialSecurity.StartAge {
			ssIncome = input.SocialSecurity.EstimatedMonthlyBenefit * 12
			if input.COLA.ApplyColaToSocialSecurity && age > input.SocialSecurity.StartAge {
				yearsSinceSsStart := age - input.SocialSecurity.StartAge
				inflationFactor := math.Pow(1+input.COLA.AssumedInflationRate, float64(yearsSinceSsStart))
				ssIncome *= inflationFactor
			}
		}
		yearData.SocialSecurity = ssIncome
		
		// Calculate TSP withdrawals
		tspWithdrawal := 0.0
		if age >= input.TSP.WithdrawalStartAge && currentTSPBalance > 0 {
			switch input.TSP.WithdrawalStrategy {
			case "fixed":
				tspWithdrawal = input.TSP.FixedWithdrawalAmount
				if tspWithdrawal > currentTSPBalance {
					tspWithdrawal = currentTSPBalance // Can't withdraw more than what's left
				}
			case "percentage":
				tspWithdrawal = currentTSPBalance * input.TSP.WithdrawalPercentage
			case "rmd":
				// Simplified RMD calculation - would need more complex table lookup in real implementation
				rmdFactor := float64(90 - age)
				if rmdFactor < 10 {
					rmdFactor = 10 // Minimum factor
				}
				tspWithdrawal = currentTSPBalance / rmdFactor
			}
		}
		yearData.TSPWithdrawal = tspWithdrawal
		
		// Calculate other income
		otherIncome := 0.0
		for _, source := range input.OtherIncome.Sources {
			if age >= source.StartAge && (source.EndAge == nil || age <= *source.EndAge) {
				amount := source.Amount
				if source.Frequency == "monthly" {
					amount *= 12
				}
				
				if source.ApplyCola && age > source.StartAge {
					yearsSinceStart := age - source.StartAge
					inflationFactor := math.Pow(1+input.COLA.AssumedInflationRate, float64(yearsSinceStart))
					amount *= inflationFactor
				}
				
				otherIncome += amount
			}
		}
		yearData.OtherIncome = otherIncome
		
		// Calculate total gross income
		yearData.TotalGrossIncome = pensionIncome + ssIncome + tspWithdrawal + otherIncome
		
		// Simplified tax calculation
		totalTaxableIncome := pensionIncome + ssIncome*0.85 + tspWithdrawal // Assume 85% of SS is taxable
		
		// Very simplified federal tax calculation (would need more complex bracketing in real implementation)
		var federalTaxRate float64
		switch input.Tax.FilingStatus {
		case "married_joint":
			if totalTaxableIncome < 20000 {
				federalTaxRate = 0.10
			} else if totalTaxableIncome < 80000 {
				federalTaxRate = 0.15
			} else if totalTaxableIncome < 170000 {
				federalTaxRate = 0.22
			} else {
				federalTaxRate = 0.24
			}
		default: // Single, etc.
			if totalTaxableIncome < 10000 {
				federalTaxRate = 0.10
			} else if totalTaxableIncome < 40000 {
				federalTaxRate = 0.15
			} else if totalTaxableIncome < 85000 {
				federalTaxRate = 0.22
			} else {
				federalTaxRate = 0.24
			}
		}
		
		yearData.FederalTax = (totalTaxableIncome - input.Tax.ItemizedDeductions) * federalTaxRate
		if yearData.FederalTax < 0 {
			yearData.FederalTax = 0
		}
		
		// State tax
		yearData.StateTax = totalTaxableIncome * input.Tax.StateIncomeTaxRate
		
		// Total taxes
		yearData.TotalTaxes = yearData.FederalTax + yearData.StateTax
		
		// Net income
		yearData.NetIncome = yearData.TotalGrossIncome - yearData.TotalTaxes
		
		// Update cumulative values
		cumulativeGrossIncome += yearData.TotalGrossIncome
		cumulativeTaxes += yearData.TotalTaxes
		cumulativeNetIncome += yearData.NetIncome
		
		// Update TSP balance
		if age < input.TSP.WithdrawalStartAge {
			// If still working/contributing before withdrawals start
			currentTSPBalance = currentTSPBalance * (1 + input.TSP.ExpectedReturnRate) + input.TSP.AnnualContribution
		} else {
			// During withdrawal phase
			currentTSPBalance = (currentTSPBalance - tspWithdrawal) * (1 + input.TSP.ExpectedReturnRate)
			if currentTSPBalance < 0 {
				currentTSPBalance = 0
			}
		}
		yearData.TSPBalance = currentTSPBalance
		
		// Track maximum TSP balance
		if currentTSPBalance > maxTSPBalance {
			maxTSPBalance = currentTSPBalance
		}
		
		// Add yearly data to result
		yearlyData = append(yearlyData, yearData)
	}
	
	// Set summary values
	result.YearlyData = yearlyData
	result.TotalGrossIncome = cumulativeGrossIncome
	result.TotalNetIncome = cumulativeNetIncome
	result.TotalTaxes = cumulativeTaxes
	result.MaxTSPBalance = maxTSPBalance
	
	return result
}

// CalculateSocialSecurity computes projected Social Security benefits based on user input
//export
func (a *App) CalculateSocialSecurity(input SocialSecurityInput) SocialSecurityResult {
	// Convert from API input to calculation model input
	ssInput := models.SocialSecurityCalculationInput{
		BirthYear:               input.BirthYear,
		CurrentAge:              input.CurrentAge,
		EstimatedAnnualSalary:   input.EstimatedAnnualSalary,
		YearsWorked:             input.YearsWorked,
		UserProvidedEstimate62:  input.UserProvidedEstimate62,
		UserProvidedEstimateFRA: input.UserProvidedEstimateFRA,
		UserProvidedEstimate70:  input.UserProvidedEstimate70,
		ClaimAge:                input.StartAge,
	}

	// Calculate the Social Security benefit
	result := calculation.CalculateSocialSecurity(ssInput)

	// Return the result in API format
	// Calculate full retirement age based on birth year (simplified)
	fullRetirementAge := 66
	if input.BirthYear >= 1960 {
		fullRetirementAge = 67
	}
	
	return SocialSecurityResult{
		EstimatedMonthlyAt62:  result.EstimatedAt62,
		EstimatedMonthlyAtFRA: result.EstimatedAtFRA,
		EstimatedMonthlyAt70:  result.EstimatedAt70,
		EstimatedAnnualAt62:   result.EstimatedAt62 * 12,
		EstimatedAnnualAtFRA:  result.EstimatedAtFRA * 12,
		EstimatedAnnualAt70:   result.EstimatedAt70 * 12,
		ClaimingAge:           result.ClaimingAge,
		ClaimingMonthlyAmount: result.ClaimingAmount,
		ClaimingAnnualAmount:  result.ClaimingAmount * 12,
		FullRetirementAge:     fullRetirementAge,
		Notes:                 result.Notes,
	}
}

// CalculateTSPProjection projects TSP growth and withdrawals
//export
func (a *App) CalculateTSPProjection(input TSPInput) TSPProjectionResult {
	// Initialize result
	result := TSPProjectionResult{
		YearlyData: []TSPYearData{},
		Notes:      "",
	}

	// Calculate current age based on birth year and month
	now := time.Now()
	currentYear := now.Year()
	currentMonth := int(now.Month())
	
	// More accurate age calculation that considers birth month
	currentAge := currentYear - input.BirthYear
	// Adjust age if we haven't reached their birthday month this year
	if input.BirthMonth > currentMonth {
		currentAge--
	}

	// Current balance
	currentBalance := input.CurrentBalance
	balance := currentBalance
	
	// First RMD age
	rmdStartAge := 73
	if input.BirthYear >= 1960 {
		rmdStartAge = 75
	}

	// Project year by year from current age to age 100
	maxAge := 100
	var notes string
	
	for age := currentAge; age <= maxAge; age++ {
		year := currentYear + (age - currentAge)
		yearData := TSPYearData{
			Age:              age,
			Year:             year,
			StartingBalance:  balance,
			Contributions:    0,
			Returns:          0,
			Withdrawals:      0,
			EndingBalance:    balance,
		}

		// Add contributions for working years
		if age < input.RetirementAge {
			yearData.Contributions = input.AnnualContribution
		}

		// Add investment returns
		yearData.Returns = balance * input.ExpectedReturnRate
		
		// Process withdrawals if in retirement
		if age >= input.WithdrawalStartAge {
			withdrawal := 0.0
			
			switch input.WithdrawalStrategy {
			case "fixed":
				withdrawal = input.FixedWithdrawalAmount
				if withdrawal > balance {
					withdrawal = balance // Can't withdraw more than balance
					if notes == "" {
						notes = fmt.Sprintf("Balance depleted at age %d.", age)
					}
				}
			case "percentage":
				withdrawal = balance * input.WithdrawalPercentage
			case "rmd":
				if age >= rmdStartAge {
					// Simplified RMD calculation
					lifeExpectancy := float64(90 - age)
					if lifeExpectancy < 10 {
						lifeExpectancy = 10
					}
					withdrawal = balance / lifeExpectancy
					
					// Add note if this is the first RMD
					if age == rmdStartAge && notes == "" {
						notes = fmt.Sprintf("Required Minimum Distributions begin at age %d.", rmdStartAge)
					}
				}
			}
			
			yearData.Withdrawals = withdrawal
		}
		
		// Calculate ending balance
		yearData.EndingBalance = yearData.StartingBalance + 
								  yearData.Contributions + 
								  yearData.Returns - 
								  yearData.Withdrawals
		
		if yearData.EndingBalance < 0 {
			yearData.EndingBalance = 0
			if notes == "" {
				notes = fmt.Sprintf("Balance depleted at age %d.", age)
			}
		}
		
		// Update balance for next year
		balance = yearData.EndingBalance
		
		// Add year data to results
		result.YearlyData = append(result.YearlyData, yearData)
		
		// If balance is depleted, no need to continue
		if balance <= 0.01 {
			break
		}
	}
	
	// Calculate summary stats
	lastYear := result.YearlyData[len(result.YearlyData)-1]
	result.FinalBalance = lastYear.EndingBalance
	result.MaxBalance = currentBalance
	result.TotalWithdrawals = 0
	result.TotalContributions = 0
	result.TotalReturns = 0
	
	for _, year := range result.YearlyData {
		if year.EndingBalance > result.MaxBalance {
			result.MaxBalance = year.EndingBalance
		}
		result.TotalWithdrawals += year.Withdrawals
		result.TotalContributions += year.Contributions
		result.TotalReturns += year.Returns
	}
	
	result.Notes = notes
	return result
}

// CalculateTaxLiability computes estimated federal and state tax liability
//export
func (a *App) CalculateTaxLiability(input TaxInput) TaxResult {
	// Initialize result
	result := TaxResult{
		FederalTax:         0,
		StateTax:           0,
		TotalTax:           0,
		EffectiveFederalRate: 0,
		EffectiveStateRate: 0,
		EffectiveTotalRate: 0,
		Notes:              "",
	}

	// Calculate taxable income (simple model)
	taxableIncome := input.TotalIncome
	
	// Subtract standard deduction based on filing status
	var standardDeduction float64
	switch input.FilingStatus {
	case "single":
		standardDeduction = 13850 // 2023 amount
		if input.Age >= 65 {
			standardDeduction += 1850 // Extra for seniors
		}
	case "married_joint":
		standardDeduction = 27700 // 2023 amount
		if input.Age >= 65 {
			standardDeduction += 1500 // Extra for seniors
		}
		if input.SpouseAge >= 65 {
			standardDeduction += 1500 // Extra for senior spouse
		}
	case "married_separate":
		standardDeduction = 13850 // 2023 amount
		if input.Age >= 65 {
			standardDeduction += 1500 // Extra for seniors
		}
	case "head_of_household":
		standardDeduction = 20800 // 2023 amount
		if input.Age >= 65 {
			standardDeduction += 1850 // Extra for seniors
		}
	}
	
	// Apply standard or itemized deduction, whichever is higher
	deduction := standardDeduction
	if input.ItemizedDeductions > standardDeduction {
		deduction = input.ItemizedDeductions
		result.Notes += "Using itemized deductions. "
	} else {
		result.Notes += "Using standard deduction. "
	}
	
	taxableIncome -= deduction
	if taxableIncome < 0 {
		taxableIncome = 0
	}
	
	// Apply Social Security taxation rules (simplified)
	if input.SocialSecurityIncome > 0 {
		// Simplified "provisional income" calculation
		provisionalIncome := taxableIncome + input.NonTaxableIncome + (input.SocialSecurityIncome * 0.5)
		
		// Determine how much of SS is taxable
		taxableSS := 0.0
		threshold1 := 25000.0
		threshold2 := 34000.0
		
		if input.FilingStatus == "married_joint" {
			threshold1 = 32000.0
			threshold2 = 44000.0
		}
		
		if provisionalIncome <= threshold1 {
			taxableSS = 0.0
		} else if provisionalIncome <= threshold2 {
			taxableSS = math.Min(input.SocialSecurityIncome * 0.5, (provisionalIncome - threshold1) * 0.5)
		} else {
			taxableSS = math.Min(
				input.SocialSecurityIncome * 0.85,
				(input.SocialSecurityIncome * 0.5) + ((provisionalIncome - threshold2) * 0.85),
			)
		}
		
		taxableIncome += taxableSS
		result.Notes += fmt.Sprintf("%.0f%% of Social Security is taxable. ", (taxableSS/input.SocialSecurityIncome)*100)
	}
	
	// Calculate federal tax (simplified 2023 brackets)
	federalTax := 0.0
	switch input.FilingStatus {
	case "married_joint":
		// 2023 Married Filing Jointly brackets
		if taxableIncome <= 22000 {
			federalTax = taxableIncome * 0.10
		} else if taxableIncome <= 89450 {
			federalTax = 2200 + ((taxableIncome - 22000) * 0.12)
		} else if taxableIncome <= 190750 {
			federalTax = 10294 + ((taxableIncome - 89450) * 0.22)
		} else if taxableIncome <= 364200 {
			federalTax = 32580 + ((taxableIncome - 190750) * 0.24)
		} else if taxableIncome <= 462500 {
			federalTax = 74208 + ((taxableIncome - 364200) * 0.32)
		} else if taxableIncome <= 693750 {
			federalTax = 105664 + ((taxableIncome - 462500) * 0.35)
		} else {
			federalTax = 186601.5 + ((taxableIncome - 693750) * 0.37)
		}
	case "single", "married_separate":
		// 2023 Single brackets
		if taxableIncome <= 11000 {
			federalTax = taxableIncome * 0.10
		} else if taxableIncome <= 44725 {
			federalTax = 1100 + ((taxableIncome - 11000) * 0.12)
		} else if taxableIncome <= 95375 {
			federalTax = 5147 + ((taxableIncome - 44725) * 0.22)
		} else if taxableIncome <= 182100 {
			federalTax = 16290 + ((taxableIncome - 95375) * 0.24)
		} else if taxableIncome <= 231250 {
			federalTax = 37104 + ((taxableIncome - 182100) * 0.32)
		} else if taxableIncome <= 578125 {
			federalTax = 52832 + ((taxableIncome - 231250) * 0.35)
		} else {
			federalTax = 174238.25 + ((taxableIncome - 578125) * 0.37)
		}
	case "head_of_household":
		// 2023 Head of Household brackets
		if taxableIncome <= 15700 {
			federalTax = taxableIncome * 0.10
		} else if taxableIncome <= 59850 {
			federalTax = 1570 + ((taxableIncome - 15700) * 0.12)
		} else if taxableIncome <= 95350 {
			federalTax = 6868 + ((taxableIncome - 59850) * 0.22)
		} else if taxableIncome <= 182100 {
			federalTax = 14678 + ((taxableIncome - 95350) * 0.24)
		} else if taxableIncome <= 231250 {
			federalTax = 35498 + ((taxableIncome - 182100) * 0.32)
		} else if taxableIncome <= 578100 {
			federalTax = 51226 + ((taxableIncome - 231250) * 0.35)
		} else {
			federalTax = 172623.5 + ((taxableIncome - 578100) * 0.37)
		}
	}
	
	// Calculate state tax
	stateTax := taxableIncome * input.StateIncomeTaxRate
	
	// Apply tax credits
	federalTax = math.Max(0, federalTax - input.FederalTaxCredits)
	stateTax = math.Max(0, stateTax - input.StateTaxCredits)
	
	// Calculate total tax and effective rates
	totalTax := federalTax + stateTax
	
	// Calculate effective rates
	if input.TotalIncome > 0 {
		result.EffectiveFederalRate = federalTax / input.TotalIncome
		result.EffectiveStateRate = stateTax / input.TotalIncome
		result.EffectiveTotalRate = totalTax / input.TotalIncome
	}
	
	// Set result values
	result.FederalTax = federalTax
	result.StateTax = stateTax
	result.TotalTax = totalTax
	
	return result
}

// CalculateCOLAAdjustment calculates cost of living adjustments
//export
func (a *App) CalculateCOLAAdjustment(input COLAInput) COLAResult {
	var result COLAResult
	result.YearlyAdjustments = make([]COLAYearData, input.ProjectionYears)
	
	// Start with the base amount
	currentAmount := input.BaseAmount
	result.BaseAmount = input.BaseAmount
	
	// Calculate COLA for each year
	for year := 0; year < input.ProjectionYears; year++ {
		yearData := COLAYearData{
			Year:            input.StartYear + year,
			InflationRate:   input.AssumedInflationRate,
			StartingAmount:  currentAmount,
			COLAPercentage:  0.0,
			AdjustedAmount:  currentAmount,
			CumulativeGrowth: 0.0,
		}
		
		// Apply system specific COLA rules
		var colaPercentage float64
		
		switch input.RetirementSystem {
		case "FERS":
			// FERS COLA rules: 
			// - If inflation <= 2%: Full inflation
			// - If inflation > 2% and <= 3%: 2%
			// - If inflation > 3%: Inflation - 1%
			if input.AssumedInflationRate <= 0.02 {
				colaPercentage = input.AssumedInflationRate
			} else if input.AssumedInflationRate <= 0.03 {
				colaPercentage = 0.02
			} else {
				colaPercentage = input.AssumedInflationRate - 0.01
			}
			
			// FERS retirees don't get COLA until age 62 unless special provisions
			if year == 0 && !input.IsSpecialProvision && input.RetirementAge < 62 {
				result.Notes += "FERS COLA not applied until age 62 for regular retirement. "
			}
			
			if !input.IsSpecialProvision && input.RetirementAge + year < 62 {
				colaPercentage = 0
			}
		case "CSRS", "CSRS Offset":
			// CSRS gets full COLA
			colaPercentage = input.AssumedInflationRate
		default:
			// Default to full inflation for other income sources
			colaPercentage = input.AssumedInflationRate
		}
		
		// Prorate first year COLA if applicable
		if year == 0 && input.MonthsInFirstYear > 0 && input.MonthsInFirstYear < 12 {
			colaPercentage = colaPercentage * float64(input.MonthsInFirstYear) / 12.0
			result.Notes += fmt.Sprintf("First year COLA prorated for %d months. ", input.MonthsInFirstYear)
		}
		
		// Apply the COLA
		adjustedAmount := currentAmount * (1 + colaPercentage)
		
		// Calculate cumulative growth from original base
		cumulativeGrowth := (adjustedAmount / input.BaseAmount) - 1.0
		
		// Set year data values
		yearData.COLAPercentage = colaPercentage
		yearData.AdjustedAmount = adjustedAmount
		yearData.CumulativeGrowth = cumulativeGrowth
		
		// Update current amount for next year
		currentAmount = adjustedAmount
		
		// Add year data to results
		result.YearlyAdjustments[year] = yearData
	}
	
	// Calculate final amount after all adjustments
	result.FinalAmount = currentAmount
	result.TotalGrowthRate = (result.FinalAmount / input.BaseAmount) - 1.0
	result.EffectiveAnnualRate = math.Pow(1+result.TotalGrowthRate, 1/float64(input.ProjectionYears)) - 1
	
	return result
}

// Scenario represents a single retirement scenario
type Scenario struct {
	ID   int         `json:"id"`
	Name string      `json:"name"`
	Data ScenarioData `json:"data"`
}

// ScenarioData contains all input data for a retirement scenario
type ScenarioData struct {
	Pension        PensionInput       `json:"pension"`
	SocialSecurity SocialSecurityInput `json:"socialSecurity"`
	TSP            TSPInput           `json:"tsp"`
	Tax            TaxInput           `json:"tax"`
	COLA           COLAInput          `json:"cola"`
	OtherIncome    OtherIncomeInput   `json:"otherIncome"`
}

// ScenariosCollection represents a collection of scenarios
type ScenariosCollection struct {
	Scenarios []Scenario `json:"scenarios"`
}

// SaveScenarioInput contains the data needed to save scenarios
type SaveScenarioInput struct {
	Scenarios []Scenario `json:"scenarios"`
	Filename  string     `json:"filename"`
}

// SaveScenarioResult contains the result of saving scenarios
type SaveScenarioResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Path    string `json:"path"`
}

// LoadScenarioInput contains the data needed to load scenarios
type LoadScenarioInput struct {
	Filename string `json:"filename"`
}

// LoadScenarioResult contains the result of loading scenarios
type LoadScenarioResult struct {
	Success   bool       `json:"success"`
	Message   string     `json:"message"`
	Scenarios []Scenario `json:"scenarios"`
}

// GetSavedScenariosResult contains a list of saved scenario files
type GetSavedScenariosResult struct {
	Files []string `json:"files"`
}

// SaveScenarios saves the current scenarios to a JSON file
//export
func (a *App) SaveScenarios(input SaveScenarioInput) SaveScenarioResult {
	// Create the scenarios directory if it doesn't exist
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return SaveScenarioResult{
			Success: false,
			Message: fmt.Sprintf("Failed to get user home directory: %v", err),
		}
	}
	
	scenariosDir := filepath.Join(userHomeDir, "FerexScenarios")
	if err := os.MkdirAll(scenariosDir, 0755); err != nil {
		return SaveScenarioResult{
			Success: false,
			Message: fmt.Sprintf("Failed to create scenarios directory: %v", err),
		}
	}
	
	// Make sure filename has .json extension
	filename := input.Filename
	if filepath.Ext(filename) != ".json" {
		filename += ".json"
	}
	
	// Create the full file path
	filePath := filepath.Join(scenariosDir, filename)
	
	// Create the scenarios collection
	collection := ScenariosCollection{
		Scenarios: input.Scenarios,
	}
	
	// Marshal the data to JSON
	jsonData, err := json.MarshalIndent(collection, "", "  ")
	if err != nil {
		return SaveScenarioResult{
			Success: false,
			Message: fmt.Sprintf("Failed to marshal scenarios to JSON: %v", err),
		}
	}
	
	// Write the JSON to the file
	if err := ioutil.WriteFile(filePath, jsonData, 0644); err != nil {
		return SaveScenarioResult{
			Success: false,
			Message: fmt.Sprintf("Failed to write scenarios to file: %v", err),
		}
	}
	
	return SaveScenarioResult{
		Success: true,
		Message: fmt.Sprintf("Successfully saved %d scenarios to %s", len(input.Scenarios), filename),
		Path:    filePath,
	}
}

// LoadScenarios loads scenarios from a JSON file
//export
func (a *App) LoadScenarios(input LoadScenarioInput) LoadScenarioResult {
	// Get the user's home directory
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return LoadScenarioResult{
			Success: false,
			Message: fmt.Sprintf("Failed to get user home directory: %v", err),
		}
	}
	
	// Create the full file path
	scenariosDir := filepath.Join(userHomeDir, "FerexScenarios")
	filePath := filepath.Join(scenariosDir, input.Filename)
	
	// Make sure the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return LoadScenarioResult{
			Success: false,
			Message: fmt.Sprintf("File not found: %s", filePath),
		}
	}
	
	// Read the file
	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return LoadScenarioResult{
			Success: false,
			Message: fmt.Sprintf("Failed to read scenarios file: %v", err),
		}
	}
	
	// Unmarshal the data
	var collection ScenariosCollection
	if err := json.Unmarshal(jsonData, &collection); err != nil {
		return LoadScenarioResult{
			Success: false,
			Message: fmt.Sprintf("Failed to parse scenarios file: %v", err),
		}
	}
	
	return LoadScenarioResult{
		Success:   true,
		Message:   fmt.Sprintf("Successfully loaded %d scenarios", len(collection.Scenarios)),
		Scenarios: collection.Scenarios,
	}
}

// GetSavedScenarios returns a list of saved scenario files
//export
func (a *App) GetSavedScenarios() GetSavedScenariosResult {
	// Get the user's home directory
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return GetSavedScenariosResult{
			Files: []string{},
		}
	}
	
	// Create the full path to the scenarios directory
	scenariosDir := filepath.Join(userHomeDir, "FerexScenarios")
	
	// If the directory doesn't exist, return an empty list
	if _, err := os.Stat(scenariosDir); os.IsNotExist(err) {
		return GetSavedScenariosResult{
			Files: []string{},
		}
	}
	
	// Get all JSON files in the directory
	files, err := filepath.Glob(filepath.Join(scenariosDir, "*.json"))
	if err != nil {
		return GetSavedScenariosResult{
			Files: []string{},
		}
	}
	
	// Extract just the filenames without the path
	var filenames []string
	for _, file := range files {
		filenames = append(filenames, filepath.Base(file))
	}
	
	return GetSavedScenariosResult{
		Files: filenames,
	}
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
