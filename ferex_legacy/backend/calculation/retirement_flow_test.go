package calculation

import (
	"ferex/backend/models"
	"testing"
)

func TestPensionSystemSelection(t *testing.T) {
	// Create base test input with realistic service periods
	civilianServiceType := "RegularDeductionFERS"
	servicePeriods := []models.ServicePeriod{
		{
			ID:                            "test-1",
			ServiceCategory:               "Civilian",
			CivilianServiceType:           &civilianServiceType,
			DepositRedepositPaymentStatus: "NotApplicable",
			SystemDuringService:           "FERS",
			StartDate:                     "1995-01-01",
			EndDate:                       "2025-12-31",
			IsPartTime:                    false,
		},
	}

	baseInput := models.RetirementCalculationInput{
		CalculationSystem: "", // Will be set in each test
		FERSInput: models.FERSCalculationInput{
			High3Salary:             100000,
			UnusedSickLeaveHours:    0,
			EmployeeContributions:   15000,
			ServicePeriods:          servicePeriods,
			SurvivorBenefitElection: "None",
			DateOfBirth:             "1965-01-01",
			ServiceComputationDate:  "1995-01-01",
			PlannedRetirementDate:   "2025-12-31",
		},
		CSRSInput: models.CSRSCalculationInput{
			High3Salary:             100000,
			UnusedSickLeaveHours:    0,
			SurvivorBenefitType:     "None",
			IsPartTime:              false,
			PartTimeProrationFactor: 1.0,
			EmployeeContributions:   20000,
			IsCSRSOffset:            false,
			DateOfBirth:             "1965-01-01",
			ServiceComputationDate:  "1995-01-01",
			PlannedRetirementDate:   "2025-12-31",
		},
		TSPInput: models.TSPCalculationInput{
			BaseSalaryForContributions:   100000,
			CurrentAgeYears:              60,
			BirthYear:                    1965,
			RetirementAgeYears:           60,
			CurrentTraditionalBalance:    200000,
			CurrentRothBalance:           50000,
			IsContributionPercentage:     true,
			ContributeUntilRetirement:    true,
			CatchUpContributionsEligible: true,
			WithdrawalStrategy:           "PercentageOfBalanceYearly",
			WithdrawalPercentageValue:    0.04,
			WithdrawalStartDate:          "2025-12-31",
			YearsToProjectWithdrawals:    30,
		},
		SocialSecurityInput: models.SocialSecurityCalculationInput{
			BirthYear:               1965,
			EstimatedAnnualSalary:   100000,
			YearsWorked:             30,
			UserProvidedEstimateFRA: 24000,
			ClaimAge:                67,
		},
		TaxInput: models.TaxCalculationInput{
			FilingStatus:     "Single",
			TaxYear:          2025,
			Age:              60,
			StateOfResidence: "VA",
		},
		COLAInput:       models.COLACalculationInput{},
		SurvivorInput:   models.SurvivorBenefitCalculationInput{},
		HealthInput:     models.HealthPremiumCalculationInput{},
		MonteCarloInput: models.MonteCarloInput{},
	}

	t.Run("FERS system selected", func(t *testing.T) {
		input := baseInput
		input.CalculationSystem = "FERS"

		result := PerformRetirementCalculation(input)

		// Check that the tax calculation used the FERS pension amount
		// We need to verify the tax input was set correctly, so we check GrossIncome contains the FERS amount
		expectedPension := result.FERSResult.AnnualPension

		// The GrossIncome should include the FERS pension plus other income sources
		if result.TaxResult.GrossIncome < expectedPension {
			t.Errorf("Expected gross income to include FERS pension amount %.2f, but gross income was %.2f",
				expectedPension, result.TaxResult.GrossIncome)
		}

		t.Logf("FERS pension: $%.2f", result.FERSResult.AnnualPension)
		t.Logf("CSRS pension: $%.2f", result.CSRSResult.FinalAnnuity)
		t.Logf("Tax calculation gross income: $%.2f", result.TaxResult.GrossIncome)
	})

	t.Run("CSRS system selected", func(t *testing.T) {
		input := baseInput
		input.CalculationSystem = "CSRS"

		result := PerformRetirementCalculation(input)

		// Check that the tax calculation used the CSRS pension amount
		expectedPension := result.CSRSResult.FinalAnnuity

		// The GrossIncome should include the CSRS pension plus other income sources
		if result.TaxResult.GrossIncome < expectedPension {
			t.Errorf("Expected gross income to include CSRS pension amount %.2f, but gross income was %.2f",
				expectedPension, result.TaxResult.GrossIncome)
		}

		t.Logf("FERS pension: $%.2f", result.FERSResult.AnnualPension)
		t.Logf("CSRS pension: $%.2f", result.CSRSResult.FinalAnnuity)
		t.Logf("Tax calculation gross income: $%.2f", result.TaxResult.GrossIncome)
	})

	t.Run("No system specified - fallback to higher amount", func(t *testing.T) {
		input := baseInput
		input.CalculationSystem = ""

		result := PerformRetirementCalculation(input)

		fersAmount := result.FERSResult.AnnualPension
		csrsAmount := result.CSRSResult.FinalAnnuity
		expectedPension := fersAmount
		systemUsed := "FERS"
		if csrsAmount > fersAmount {
			expectedPension = csrsAmount
			systemUsed = "CSRS"
		}

		// The GrossIncome should include the higher pension amount plus other income sources
		if result.TaxResult.GrossIncome < expectedPension {
			t.Errorf("Expected gross income to include %s pension amount %.2f, but gross income was %.2f",
				systemUsed, expectedPension, result.TaxResult.GrossIncome)
		}

		t.Logf("FERS pension: $%.2f", result.FERSResult.AnnualPension)
		t.Logf("CSRS pension: $%.2f", result.CSRSResult.FinalAnnuity)
		t.Logf("Auto-selected: %s ($%.2f)", systemUsed, expectedPension)
		t.Logf("Tax calculation gross income: $%.2f", result.TaxResult.GrossIncome)
	})
}
