package scenario

import (
	"ferex/backend/models"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

// Helper functions to get pointers to literals
func ptrFloat64(v float64) *float64 { return &v }
func ptrInt(v int) *int          { return &v }
func ptrString(v string) *string { return &v }
func ptrBool(v bool) *bool       { return &v }

func TestFerexFileRoundTrip(t *testing.T) {
	// 1. Construct a FerexFile object with comprehensive data
	originalFerexFile := FerexFile{
		FileFormatVersion:    "1.0.1",
		LastOpenedAppVersion: "0.2.0-test",
		UserProfile: UserProfile{
			EmployeeName: "Test User Jr.",
			BirthDate:    "1985-02-20",
			MRAYears:     57,
			MRAMonths:    0,
		},
		Variants: []ScenarioVariant{
			{
				VariantID:    "variant-xyz-789",
				VariantName:  "Comprehensive Round-Trip Test Variant",
				CreatedAt:    time.Date(2024, 3, 10, 12, 0, 0, 0, time.UTC),
				LastModified: time.Date(2024, 3, 11, 14, 30, 0, 0, time.UTC),
				Notes:        "This variant tests all (or most) fields for save/load.",

				// CoreScenarioInputs
				CalculationSystem:                "FERS",
				High3Salary:                      ptrFloat64(120000.75),
				ServiceComputationDate:           "2005-07-15",
				DateOfBirth:                      "1985-02-20", // Matches UserProfile for this test
				PlannedRetirementDate:            "2035-07-15",
				UnusedSickLeaveHoursAtRetirement: ptrFloat64(600.00),
				EmployeeContributions:            ptrFloat64(85000.50),
				ServicePeriods: []models.ServicePeriod{
					{
						ID: "sp-a", StartDate: "2005-07-15", EndDate: "2015-07-14",
						ServiceCategory: "Civilian",
						CivilianServiceType: ptrString("RegularDeductionFERS"),
						MilitaryServiceType: nil,
						DepositRedepositPaymentStatus: "NotApplicable",
						SystemDuringService: "FERS",
						IsPartTime: false,
						HoursPerWeekIfPartTime: nil,
						Notes: nil,
					},
					{
						ID: "sp-b", StartDate: "2015-07-15", EndDate: "2025-07-14",
						ServiceCategory: "Civilian",
						CivilianServiceType: ptrString("RegularDeductionFERS"), // Assuming FERS part-time
						MilitaryServiceType: nil,
						DepositRedepositPaymentStatus: "PaidInFull", // "Yes" mapped to "PaidInFull"
						SystemDuringService: "None",
						IsPartTime: true,
						HoursPerWeekIfPartTime: ptrFloat64(20),
						Notes: nil,
					},
				},
				LwopPeriods: []models.LWOPPeriod{
					{ID: "lwop-a", StartDate: "2010-01-01", EndDate: "2010-06-30", Type: "PersonalNonMilitary"},
				},
				IsVeraDsRetirement:                 true,
				MraPlus10PostponeAnnuityStartDate:  ptrString("2037-02-20"), 
				IsDeferredRetirement:               false,
				DeferredRetirementAnnuityStartDate: nil, // Example: ptrString("2045-02-20") if IsDeferredRetirement is true
				IsDisabilityRetirement:             false,
				IsMarriedAtRetirement:            true,
				HasFormerSpouseEntitlement:         true,
				FormerSpouseSurvivorBenefitDetails: &models.FormerSpouseSurvivorDetails{
					CourtOrderNumber:    ptrString("CO-12345"),
					DateOfCourtOrder:    ptrString("2018-05-01"),
					PortionAwarded:      ptrString("Percentage"), // Example: 'Percentage', 'FixedAmount', 'ProRataShare'
					CustomAmountAwarded: ptrFloat64(25.0),      // Value if PortionAwarded indicates a specific amount or percentage
				},
				AgeOfInsurableInterestBeneficiary: ptrInt(60), // Relevant if InsurableInterest is chosen

				// FersScenarioInputs
				FersCoverageType:             ptrString("FERS_FRAE"),
				EstimatedSSBenefitAt62ForSRS: ptrFloat64(1800.00),
				DidSwitchFromCSRS:            true,
				SwitchedToFERSDate:           ptrString("2003-01-01"),
				SurvivorBenefitFERS: &models.SurvivorBenefitFERSInput{
					SpouseElection:                ptrString("FullSpouse"), 
					FormerSpouseElection:          ptrString("FullFormerSpouse"),
					FormerSpouseConsent:           ptrString("Yes"),
					CurrentSpouseConsentForFormer: ptrString("Yes"),
					CurrentSpouseWaiverForSelf:    ptrString("No"),
					InsurableInterestDetails:      &models.InsurableInterestDetails{
						RelationshipToEmployee: "Parent",
						DateOfBirth:            "1960-05-15",
					},
				},

				// SocialSecurityScenarioInputs
				UserProvidedSSBenefitAmount1:      ptrFloat64(2200.00),
				UserProvidedSSBenefitClaimingAge1: ptrFloat64(67.5),
				UserProvidedSSBenefitAmount2:      ptrFloat64(2350.00),
				UserProvidedSSBenefitClaimingAge2: ptrFloat64(70.0),
				SSBenefitSpousalOption:            ptrString("ClaimOwnThenSpousal"),
				SSBenefitSpousalAmount:            ptrFloat64(1100.00),
				SSBenefitSurvivorOption:           ptrString("ClaimOwnThenSurvivor"),
				SSBenefitSurvivorAmount:           ptrFloat64(2000.00),
				UserSSClaimingAge:                 ptrFloat64(70.0),
				UserAssumedSSCOLA:                 ptrFloat64(0.0225),
				CalculateHistoricalWEPGPO:         true,
				SocialSecurityEstimateSource:      ptrString("pia"),
				SocialSecurityEstimate:            ptrFloat64(2400.00),
				SocialSecurityFRA:                 ptrFloat64(67.0),
				SocialSecurityBenefitStartDate:    ptrString("2055-02-20"),
				SocialSecurityCreditedEarnings: []models.SocialSecurityCreditedEarning{
					{Year: 2022, Earnings: 90000},
					{Year: 2023, Earnings: 92000},
				},
				SocialSecurityBenefitTable: []models.SocialSecurityCreditedEarning{ // Example, often PIA calculation output
					{Year: 62, Earnings: 1800},
					{Year: 67, Earnings: 2400},
					{Year: 70, Earnings: 2976},
				},
				SpouseBirthDate:      ptrString("1987-07-10"),
				SpouseSSClaimingAge:  ptrFloat64(68.0),
				SpouseSocialSecurityEstimate:          ptrFloat64(1200.00),
				IsSpouseSubjectToWEP:                  ptrBool(false),
				SpouseWEPNonCoveredPensionAmount:      ptrFloat64(0.00),
				IsSpouseSubjectToGPO:                  ptrBool(true),
				SpouseGPONonCoveredSpousalBenefitAmount: ptrFloat64(400.00),

				// TSPScenarioInputs
				TSPBalanceTraditional:                   ptrFloat64(200000.00),
				TSPBalanceRoth:                          ptrFloat64(75000.00),
				TSPLoanBalance:                          ptrFloat64(5000.00),
				TSPAnnualContributionPreRetirement:      ptrFloat64(22500.00),
				TSPContributionPercentagePreRetirement:  ptrFloat64(15.0), // 15% of salary
				TSPContributeUntil:                      ptrString("RetirementDate"), // "RetirementDate" or "SpecificAge"
				TSPContributionStopAge:                  ptrInt(65),
				TSPCatchUpContributionsEligible:         true,
				TSPContributionAllocationTraditionalPercent: ptrFloat64(70.0), // 70% to Traditional
				TSPContributionAllocationRothPercent:    ptrFloat64(30.0), // 30% to Roth
				TSPContributionAllocationToFunds: &models.TSPFundAllocation{
					G: ptrFloat64(10.0), F: ptrFloat64(10.0), C: ptrFloat64(30.0), S: ptrFloat64(30.0), I: ptrFloat64(20.0),
				},
				TSPCurrentAllocationToFunds: &models.TSPFundAllocation{ 
					G: ptrFloat64(15.0), C: ptrFloat64(50.0), S: ptrFloat64(35.0), F: ptrFloat64(0.0), I: ptrFloat64(0.0), LFundName: nil,
				},
				TSPPostRetirementAllocation: &models.TSPFundAllocation{
					G: ptrFloat64(50.0), F: ptrFloat64(30.0), C: ptrFloat64(10.0), S: ptrFloat64(5.0), I: ptrFloat64(5.0), LFundName: ptrString("LIncome"),
				},
				TSPFutureAllocationStrategy:             ptrString("UsePostRetirementAllocation"), // "MaintainCurrent", "ReallocateTo", "LifecycleFund", "UsePostRetirementAllocation"
				TSPFutureAllocationTarget:             &models.TSPFundAllocation{ // For 'ReallocateTo' strategy
					G: ptrFloat64(20.0), F: ptrFloat64(20.0), C: ptrFloat64(20.0), S: ptrFloat64(20.0), I: ptrFloat64(20.0),
				},
				UserReturnAssumptionsTSP: &models.TSPIndividualReturnAssumptions{
					G: ptrFloat64(0.025), F: ptrFloat64(0.035), C: ptrFloat64(0.075), S: ptrFloat64(0.085), I: ptrFloat64(0.065),
				},
				TSPExpenseRatio:                         ptrFloat64(0.00051), // 0.051%
				TSPExpectedAnnualGrowthRatePreRetirement:  ptrFloat64(0.065), // Overall blended rate
				TSPWithdrawalStrategy: ptrString("PercentageOfBalance"),
				TSPWithdrawalFixedAmountValue:           ptrFloat64(3000.00), // For "FixedAmount..." strategy
				TSPWithdrawalPercentageValue: ptrFloat64(4.0),
				TSPWithdrawalStartAge: ptrInt(62),
				TSPWithdrawalStartDate:                  ptrString("SpecificAge"), // "RetirementDate" or "SpecificAge"
				TSPWithdrawalOrder:                      ptrString("ProRata"),     // "ProRata", "TraditionalFirst", "RothFirst"
				ApplyRMDsToTSP:        true,
				TSPExpectedAnnualGrowthRatePostRetirement: ptrFloat64(0.045),
				TSPVolatilityAssumptions: &models.TSPVolatilityRates{
					GStdDev: ptrFloat64(0.005), FStdDev: ptrFloat64(0.02), CStdDev: ptrFloat64(0.15), SStdDev: ptrFloat64(0.18), IStdDev: ptrFloat64(0.16),
				},
				TSPLifecycleFundTargetDate: ptrString("2050"),

				// TSP Annuity Details
				TSPAnnuityType:                      ptrString("JointLifeSpouse"), // 'SingleLife' | 'JointLifeSpouse' | 'JointLifeNonSpouse'
				TSPAnnuitySpouseAge:                 ptrInt(60),
				TSPAnnuityFeatures:                  []string{"IncreasingPayments", "CashRefund"}, // ('IncreasingPayments' | 'CashRefund' | 'TenYearCertain')
				TSPAnnuityAmountForPurchase:         ptrFloat64(100000.00),
				TSPAnnuityPurchaseAge:               ptrInt(65),

				// InsuranceScenarioInputs
				FegliBasicCoverage:          true,
				FegliOptionAAmount:          ptrFloat64(10000.00),
				FegliOptionBCoverageMultiples: ptrInt(3),
				FegliOptionCCoverageMultiples: ptrInt(2),
				FegliPost65Reduction:        ptrString("No Reduction"),
				FegliBasicRetiredCoverage:   ptrString("Full"), // "Full", "50%", "None"
				FegliOptionARetiredCoverage: ptrString("Full"), // "Full", "None"
				FegliOptionBRetiredCoverage: ptrString("Full"), // "Full", "None"
				FegliOptionCRetiredCoverage: ptrString("Full"), // "Full", "None"
				FehbPlanName:                       ptrString("GEHA Standard"),
				FehbPlanCode:                       ptrString("GEHA123"),
				FehbCoverageType:                   ptrString("SelfPlusOne"), // "Self", "SelfPlusOne", "SelfAndFamily"
				FehbPremiumPreRetirementMonthly:    ptrFloat64(450.00),
				FehbPremiumPostRetirementMonthly:   ptrFloat64(500.00),
				FehbContinueInRetirement:           true,
				UserAssumedFehbPremiumGrowthRate: ptrFloat64(0.045),
				FltcipEnrolled:             true,
				FltcipCurrentStatus:        ptrString("enrolledPaying"), // 'enrolledPaying' | 'enrolledPaidUp' | 'notEnrolled'
				FltcipDailyBenefitAmount:   ptrFloat64(200.00),
				FltcipBenefitPeriod:        ptrString("5years"),
				FltcipInflationProtection:  ptrString("AutomaticCompound"),
				FltcipWaitingPeriod:        ptrString("90days"),
				FltcipPremiumMonthly:       ptrFloat64(150.00),
				FltcipContinueInRetirement: true,

				// FEDVIP Details
				FedvipDentalPlanName:              ptrString("MetLife Dental High"),
				FedvipDentalPremium:               ptrFloat64(80.00),
				FedvipVisionPlanName:              ptrString("VSP Vision Plus"),
				FedvipVisionPremium:               ptrFloat64(30.00),

				OtherLifeInsuranceCoverageAmount:  ptrFloat64(250000.00),
				OtherLifeInsurancePremiumAnnual:   ptrFloat64(600.00),
				DisabilityInsuranceCoverageAmount: ptrFloat64(5000.00), // Monthly benefit
				DisabilityInsurancePremiumAnnual:  ptrFloat64(1200.00),
				LongTermCareInsurancePolicyDetails: ptrString("Private LTC policy, details..."),

				// TaxScenarioInputs
				FederalTaxFilingStatus:        ptrString("Single"),
				FederalTaxNumberOfDependents:  ptrInt(0),
				FederalTaxExtraWithholding:    ptrFloat64(100.00),
				FederalTaxAdjustmentsToIncome: ptrFloat64(2000.00),
				FederalTaxItemizedDeductions:  ptrFloat64(13850.00),
				FederalOtherTaxableIncomeAnnual: ptrFloat64(5000.00),
				FederalTaxCreditsAnnual:       ptrFloat64(500.00),
				FederalTaxLawAssumption:       ptrString("CurrentLaw"), // 'CurrentLaw' | 'ExtendsTCJA' | 'Custom'
				UserAssumedFederalTaxRateChange: ptrString("increase"), // 'none' | 'increase' | 'decrease' | 'custom'
				FederalTaxUserAssumedRateChangeAmount: ptrFloat64(0.02), // 2% increase
				FederalTaxUserAssumedRateChangeYear: ptrInt(2028),
				StateOfResidenceForTax:       ptrString("MD"),
				RetirementEffectiveStateTaxRate: ptrFloat64(0.0575),
				StateTaxFilingStatus:          ptrString("Single"),
				StateTaxNumberOfDependents:    ptrInt(0),
				StateTaxExtraWithholding:      ptrFloat64(50.00),
				StateTaxAdjustmentsToIncome:   ptrFloat64(1000.00),
				StateTaxItemizedDeductions:    ptrFloat64(5000.00),
				UserAssumedStateTaxRateChange: ptrString("none"),
				StateTaxUserAssumedRateChangeAmount: ptrFloat64(0.00),
				StateTaxUserAssumedRateChangeYear: ptrInt(0),
				TaxTreatmentOfMilitaryRetirementPay: ptrString("partiallyTaxable"),

				// OtherFinancialsScenarioInputs
				OtherTaxableAccountBalance:          ptrFloat64(75000.00),
				OtherTaxableAccountAllocation:     &models.TaxableAccountAssetAllocation{
					StocksPercent: ptrFloat64(60.0), BondsPercent: ptrFloat64(30.0), CashPercent: ptrFloat64(10.0),
				},
				OtherTaxableAccountAnnualContribution: ptrFloat64(5000.00),
				OtherTaxableAccountExpectedGrowth: ptrFloat64(0.06),
				OtherTaxAdvantagedAccountBalance:  ptrFloat64(150000.00),
				OtherTaxAdvantagedAccountExpectedGrowth: ptrFloat64(0.07),
				OtherTaxAdvantagedAccountAnnualContribution: ptrFloat64(6000.00),
				OtherRecurringIncomeSources: []models.OtherRecurringIncomeSource{
					{ID: "income-x", Name: "Consulting Gig", Amount: ptrFloat64(500.00), Frequency: "Monthly", StartDate: "2030-01-01", EndDate: ptrString("2035-12-31"), IsInflationAdjusted: ptrBool(false), AnnualIncreaseRate: ptrFloat64(0.02), SubjectToFICA: ptrBool(true)},
					{ID: "income-y", Name: "Rental Income", Amount: ptrFloat64(1200.00), Frequency: "Monthly", StartDate: "2025-01-01", EndDate: nil, IsInflationAdjusted: ptrBool(true), AnnualIncreaseRate: ptrFloat64(0.03), SubjectToFICA: ptrBool(false)},
				},
			},
		},
		LastViewedVariantID: "variant-xyz-789",
	}

	// 2. Create a temporary file
	tmpDir := t.TempDir()
	tmpFilePath := filepath.Join(tmpDir, "test_ferex_file_roundtrip.json")

	// 3. Save the FerexFile
	err := SaveFerexFile(tmpFilePath, &originalFerexFile)
	if err != nil {
		t.Fatalf("SaveFerexFile failed: %v", err)
	}

	// 4. Load the FerexFile
	loadedFerexFile, err := LoadFerexFile(tmpFilePath)
	if err != nil {
		t.Fatalf("LoadFerexFile failed: %v", err)
	}

	// 5. Compare the original and loaded FerexFile objects
	if originalFerexFile.FileFormatVersion != loadedFerexFile.FileFormatVersion {
		t.Errorf("FileFormatVersion mismatch. Original: %s, Loaded: %s", originalFerexFile.FileFormatVersion, loadedFerexFile.FileFormatVersion)
	}
	if originalFerexFile.LastOpenedAppVersion != loadedFerexFile.LastOpenedAppVersion {
		t.Errorf("LastOpenedAppVersion mismatch. Original: %s, Loaded: %s", originalFerexFile.LastOpenedAppVersion, loadedFerexFile.LastOpenedAppVersion)
	}
	if !reflect.DeepEqual(originalFerexFile.UserProfile, loadedFerexFile.UserProfile) {
		t.Errorf("UserProfile mismatch.\nOriginal: %+v\nLoaded:   %+v", originalFerexFile.UserProfile, loadedFerexFile.UserProfile)
	}
	if originalFerexFile.LastViewedVariantID != loadedFerexFile.LastViewedVariantID {
		t.Errorf("LastViewedVariantID mismatch. Original: %s, Loaded: %s", originalFerexFile.LastViewedVariantID, loadedFerexFile.LastViewedVariantID)
	}

	if len(originalFerexFile.Variants) != len(loadedFerexFile.Variants) {
		t.Fatalf("Number of variants differs. Original: %d, Loaded: %d", len(originalFerexFile.Variants), len(loadedFerexFile.Variants))
	}

	for i := range originalFerexFile.Variants {
		originalVariant := originalFerexFile.Variants[i]
		loadedVariant := loadedFerexFile.Variants[i]

		// Store original time.Time fields
		originalCreatedAt := originalVariant.CreatedAt
		originalLastModified := originalVariant.LastModified

		// Nil out time.Time fields for DeepEqual comparison of the rest of the struct
		// This is because JSON marshalling/unmarshalling can affect monotonic clock readings
		// or sub-second precision in a way that makes DeepEqual fail.
		originalVariant.CreatedAt = time.Time{}
		originalVariant.LastModified = time.Time{}
		loadedVariant.CreatedAt = time.Time{}
		loadedVariant.LastModified = time.Time{}

		if !reflect.DeepEqual(originalVariant, loadedVariant) {
			t.Errorf("Loaded ScenarioVariant at index %d does not match original (excluding times).\nOriginal: %+v\nLoaded:   %+v", i, originalVariant, loadedVariant)
		}

		// Restore original time.Time fields for originalVariant for individual comparison
		originalVariant.CreatedAt = originalCreatedAt
		originalVariant.LastModified = originalLastModified

		// Compare time.Time fields using .Equal() and ensuring UTC for consistency
		if !originalCreatedAt.Equal(loadedFerexFile.Variants[i].CreatedAt.In(time.UTC)) {
			t.Errorf("Variant %d CreatedAt mismatch. Original: %v, Loaded: %v", i, originalCreatedAt, loadedFerexFile.Variants[i].CreatedAt.In(time.UTC))
		}
		if !originalLastModified.Equal(loadedFerexFile.Variants[i].LastModified.In(time.UTC)) {
			t.Errorf("Variant %d LastModified mismatch. Original: %v, Loaded: %v", i, originalLastModified, loadedFerexFile.Variants[i].LastModified.In(time.UTC))
		}
		
		// Restore times in loadedFerexFile for the final DeepEqual of the whole structure if needed
		// Ensure they are in UTC as well for consistency if original times were UTC.
		loadedFerexFile.Variants[i].CreatedAt = loadedFerexFile.Variants[i].CreatedAt.In(time.UTC)
		loadedFerexFile.Variants[i].LastModified = loadedFerexFile.Variants[i].LastModified.In(time.UTC)
	}

	// Final check on the whole structure. This might be redundant if individual checks are thorough
	// but can catch issues with the overall FerexFile structure or slice handling.
	// Ensure originalFerexFile also has its variant times restored before this final comparison.
	if !reflect.DeepEqual(originalFerexFile, *loadedFerexFile) {
	    t.Errorf("Loaded FerexFile does not deeply match original after individual checks and time adjustments.\nOriginal: %+v\nLoaded:   %+v", originalFerexFile, *loadedFerexFile)
	}
}
