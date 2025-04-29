package tests

import (
	"ferex/backend/calculation"
	"ferex/backend/models"
	"testing"
)

func TestRetirementCalculationIntegration(t *testing.T) {
	input := models.RetirementCalculationInput{
		FERSInput: models.FERSCalculationInput{
			High3Salary:             90000,
			YearsOfService:          30,
			AgeAtRetirement:         62,
			UnusedSickLeaveMonths:   0,
			SurvivorBenefitElection: "max",
		},
		CSRSInput: models.CSRSCalculationInput{}, // Not used in this scenario
		SRSInput: models.SRSCalculationInput{
			EstimatedSocialSecurityAt62: 0, // Not used in this scenario
			YearsOfFERSService:          30,
			RetirementAge:               62,
			MRA:                         57,
			IsImmediateUnreducedAnnuity: true,
			ProjectedEarnedIncome:       0,
			RetirementYear:              2025,
		},
		TSPInput: models.TSPCalculationInput{
			CurrentBalance:             500000,
			AnnualEmployeeContribution: 0,
			AnnualAgencyMatch:          0,
			YearsUntilRetirement:       0,
			ExpectedAnnualReturnRate:   0,
			WithdrawalStartAge:         62,
			WithdrawalMethod:           "fixed",
			WithdrawalAmountOrPercent:  20000,
			ProjectedWithdrawalYears:   30,
			RothBalance:                0,
			TraditionalBalance:         0,
		},
		TaxInput: models.TaxCalculationInput{
			FilingStatus:       "single",
			TaxYear:            2025,
			GrossPension:       40000,
			TaxablePension:     40000,
			TSPWithdrawal:      20000,
			TSPRothWithdrawal:  0,
			SocialSecurity:     0,
			OtherTaxableIncome: 0,
			Deductions:         14600,
			TaxCredits:         0,
		},
		SocialSecurityInput: models.SocialSecurityCalculationInput{
			BirthYear:               1963,
			CurrentAge:              62,
			UserProvidedEstimate62:  1800,
			UserProvidedEstimateFRA: 2400,
			UserProvidedEstimate70:  3100,
			ClaimAge:                62,
		},
		COLAInput: models.COLACalculationInput{
			InitialAmount: 40000,
			COLARate:      0.02,
			Years:         5,
			COLAPolicy:    "FERS",
		},
		SurvivorInput: models.SurvivorBenefitCalculationInput{
			PensionType:      "FERS",
			InitialAnnuity:   40000,
			SurvivorElection: "max",
			COLARate:         0.02,
			YearsToProject:   5,
		},
		HealthInput: models.HealthPremiumCalculationInput{
			FEHBPremium:     5000,
			IncludeFEHB:     true,
			IncludeMedicare: false,
			COLARate:        0.0,
			YearsToProject:  1,
		},
	}

	result := calculation.CalculateRetirementProjection(input)

	if result.NetAfterTaxIncome <= 0 {
		t.Errorf("NetAfterTaxIncome should be positive, got %.2f", result.NetAfterTaxIncome)
	}
	if result.TotalRetirementIncome <= 0 {
		t.Errorf("TotalRetirementIncome should be positive, got %.2f", result.TotalRetirementIncome)
	}
	if result.FERSResult.AnnualPension <= 0 {
		t.Errorf("FERSResult.AnnualPension should be positive, got %.2f", result.FERSResult.AnnualPension)
	}
	if result.SRSResult.AnnualSRSAmount < 0 {
		t.Errorf("SRSResult.AnnualSRSAmount should not be negative, got %.2f", result.SRSResult.AnnualSRSAmount)
	}
	if result.TSPResult.AnnualWithdrawalIncome < 0 {
		t.Errorf("TSPResult.AnnualWithdrawalIncome should not be negative, got %.2f", result.TSPResult.AnnualWithdrawalIncome)
	}
	if result.TaxResult.FederalTaxOwed < 0 {
		t.Errorf("TaxResult.FederalTaxOwed should not be negative, got %.2f", result.TaxResult.FederalTaxOwed)
	}
	if result.HealthResult.ProjectedPremiums[0] != 5000 {
		t.Errorf("HealthResult.ProjectedPremiums[0] should be 5000, got %.2f", result.HealthResult.ProjectedPremiums[0])
	}
	if result.Notes == "" {
		t.Errorf("Notes should not be empty")
	}
}
