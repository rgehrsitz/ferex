package main

import (
	"context"
	"fmt"
	"math"
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
	StartAge              int     `json:"startAge"`
	EstimatedMonthlyBenefit float64 `json:"estimatedMonthlyBenefit"`
	IsEligible            bool    `json:"isEligible"`
	BirthYear             int     `json:"birthYear"`
	BirthMonth            int     `json:"birthMonth"`
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
}

// TaxInput contains data for tax calculations
type TaxInput struct {
	FilingStatus          string  `json:"filingStatus"`
	StateOfResidence      string  `json:"stateOfResidence"`
	AdditionalIncome      float64 `json:"additionalIncome"`
	AdditionalDeductions  float64 `json:"additionalDeductions"`
	StateIncomeTaxRate    float64 `json:"stateIncomeTaxRate"`
}

// COLAInput contains data for cost of living adjustment calculations
type COLAInput struct {
	AssumedInflationRate  float64 `json:"assumedInflationRate"`
	ApplyCOLAToPension    bool    `json:"applyCOLAToPension"`
	ApplyColaToSocialSecurity bool `json:"applyColaToSocialSecurity"`
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
	
	// Determine year of birth for age calculations
	birthYear := input.SocialSecurity.BirthYear
	currentYear := time.Now().Year()
	
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
	
	// Calculate projections for each year
	for age := startAge; age <= endAge; age++ {
		year := birthYear + age
		
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
		
		yearData.FederalTax = (totalTaxableIncome - input.Tax.AdditionalDeductions) * federalTaxRate
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

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
