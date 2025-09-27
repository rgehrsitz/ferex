package scenario

import (
	"time"

	"ferex/backend/models" // Corrected import path
)

// UserProfile contains data truly global to the person using Ferex for this file.
type UserProfile struct {
	EmployeeName string `json:"employeeName,omitempty"`
	BirthDate    string `json:"birthDate"` // YYYY-MM-DD, for consistency
	MRAYears     int    `json:"mraYears"`
	MRAMonths    int    `json:"mraMonths"`
}

// ScenarioVariant holds a complete set of inputs for one distinct scenario.
// It is designed to be self-contained and fully independent after creation.
// All fields are directly included for simplicity in save/load operations.
type ScenarioVariant struct {
	VariantID    string    `json:"variantId"`
	VariantName  string    `json:"variantName"`
	CreatedAt    time.Time `json:"createdAt"`    // wails:type string
	LastModified time.Time `json:"lastModified"` // wails:type string
	Notes        string    `json:"notes,omitempty"`

	// Fields from CoreScenarioInputs
	CalculationSystem                  string                              `json:"calculationSystem"` // 'FERS' | 'CSRS' | ''
	High3Salary                        *float64                            `json:"high3Salary,omitempty"`
	ServiceComputationDate             string                              `json:"serviceComputationDate"` // YYYY-MM-DD
	DateOfBirth                        string                              `json:"dateOfBirth"`            // YYYY-MM-DD (Matches UserProfile.BirthDate, but stored per variant for independence)
	PlannedRetirementDate              string                              `json:"plannedRetirementDate"`  // YYYY-MM-DD
	UnusedSickLeaveHoursAtRetirement   *float64                            `json:"unusedSickLeaveHoursAtRetirement,omitempty"`
	EmployeeContributions              *float64                            `json:"employeeContributions,omitempty"`
	ServicePeriods                     []models.ServicePeriod              `json:"servicePeriods,omitempty"`
	LwopPeriods                        []models.LWOPPeriod                 `json:"lwopPeriods,omitempty"`
	IsVeraDsRetirement                 bool                                `json:"isVeraDsRetirement,omitempty"`
	MraPlus10PostponeAnnuityStartDate  *string                             `json:"mraPlus10PostponeAnnuityStartDate,omitempty"` // YYYY-MM-DD
	IsDeferredRetirement               bool                                `json:"isDeferredRetirement,omitempty"`
	DeferredRetirementAnnuityStartDate *string                             `json:"deferredRetirementAnnuityStartDate,omitempty"` // YYYY-MM-DD
	IsDisabilityRetirement             bool                                `json:"isDisabilityRetirement,omitempty"`
	IsMarriedAtRetirement              bool                                `json:"isMarriedAtRetirement,omitempty"`
	HasFormerSpouseEntitlement         bool                                `json:"hasFormerSpouseEntitlement,omitempty"`
	AgeOfInsurableInterestBeneficiary  *int                                `json:"ageOfInsurableInterestBeneficiary,omitempty"` // Relevant if InsurableInterest is chosen for FERS
	FormerSpouseSurvivorBenefitDetails *models.FormerSpouseSurvivorDetails `json:"formerSpouseSurvivorBenefitDetails,omitempty"`

	// Fields from FersScenarioInputs (Applicable if CalculationSystem == "FERS")
	FersCoverageType             *string                          `json:"fersCoverageType,omitempty"` // 'FERS' | 'FERS_RAE' | 'FERS_FRAE' | null
	EstimatedSSBenefitAt62ForSRS *float64                         `json:"estimatedSSBenefitAt62ForSRS,omitempty"`
	DidSwitchFromCSRS            bool                             `json:"didSwitchFromCSRS,omitempty"`
	SwitchedToFERSDate           *string                          `json:"switchedToFERSDate,omitempty"` // YYYY-MM-DD
	SurvivorBenefitFERS          *models.SurvivorBenefitFERSInput `json:"survivorBenefitFERS,omitempty"`

	// Fields from CsrsScenarioInputs (Applicable if CalculationSystem == "CSRS")
	IsCSRSOffset         bool                             `json:"isCSRSOffset,omitempty"`
	YearsOfOffsetService *float64                         `json:"yearsOfOffsetService,omitempty"`
	SsAt62WithOffset     *float64                         `json:"ssAt62WithOffset,omitempty"`
	SsAt62WithoutOffset  *float64                         `json:"ssAt62WithoutOffset,omitempty"`
	SurvivorBenefitCSRS  *models.SurvivorBenefitCSRSInput `json:"survivorBenefitCSRS,omitempty"`

	// Fields from SocialSecurityScenarioInputs
	UserProvidedSSBenefitAmount1      *float64                               `json:"userProvidedSSBenefitAmount1,omitempty"`
	UserProvidedSSBenefitClaimingAge1 *float64                               `json:"userProvidedSSBenefitClaimingAge1,omitempty"`
	UserProvidedSSBenefitAmount2      *float64                               `json:"userProvidedSSBenefitAmount2,omitempty"`
	UserProvidedSSBenefitClaimingAge2 *float64                               `json:"userProvidedSSBenefitClaimingAge2,omitempty"`
	SSBenefitSpousalOption            *string                                `json:"ssBenefitSpousalOption,omitempty"`
	SSBenefitSpousalAmount            *float64                               `json:"ssBenefitSpousalAmount,omitempty"`
	SSBenefitSurvivorOption           *string                                `json:"ssBenefitSurvivorOption,omitempty"`
	SSBenefitSurvivorAmount           *float64                               `json:"ssBenefitSurvivorAmount,omitempty"`
	UserSSClaimingAge                 *float64                               `json:"userSSClaimingAge,omitempty"`
	UserAssumedSSCOLA                 *float64                               `json:"userAssumedSSCOLA,omitempty"`
	CalculateHistoricalWEPGPO         bool                                   `json:"calculateHistoricalWEPGPO,omitempty"`
	SocialSecurityEstimateSource      *string                                `json:"socialSecurityEstimateSource,omitempty"` // 'pia' | 'userProvided' | 'ssaStatement'
	SocialSecurityEstimate            *float64                               `json:"socialSecurityEstimate,omitempty"`
	SocialSecurityFRA                 *float64                               `json:"socialSecurityFRA,omitempty"`
	SocialSecurityBenefitStartDate    *string                                `json:"socialSecurityBenefitStartDate,omitempty"`
	SocialSecurityCreditedEarnings    []models.SocialSecurityCreditedEarning `json:"socialSecurityCreditedEarnings,omitempty"`
	SocialSecurityBenefitTable        []models.SocialSecurityCreditedEarning `json:"socialSecurityBenefitTable,omitempty"`

	// Spouse's Social Security Details
	SpouseBirthDate                         *string  `json:"spouseBirthDate,omitempty"` // YYYY-MM-DD
	SpouseSSClaimingAge                     *float64 `json:"spouseSSClaimingAge,omitempty"`
	SpouseSocialSecurityEstimate            *float64 `json:"spouseSocialSecurityEstimate,omitempty"`
	IsSpouseSubjectToWEP                    *bool    `json:"isSpouseSubjectToWEP,omitempty"`
	SpouseWEPNonCoveredPensionAmount        *float64 `json:"spouseWEPNonCoveredPensionAmount,omitempty"`
	IsSpouseSubjectToGPO                    *bool    `json:"isSpouseSubjectToGPO,omitempty"`
	SpouseGPONonCoveredSpousalBenefitAmount *float64 `json:"spouseGPONonCoveredSpousalBenefitAmount,omitempty"`

	// Fields from TSPScenarioInputs
	TSPBalanceTraditional                       *float64                               `json:"tspBalanceTraditional,omitempty"`
	TSPBalanceRoth                              *float64                               `json:"tspBalanceRoth,omitempty"`
	TSPLoanBalance                              *float64                               `json:"tspLoanBalance,omitempty"`
	TSPAnnualContributionPreRetirement          *float64                               `json:"tspAnnualContributionPreRetirement,omitempty"`
	TSPContributionPercentagePreRetirement      *float64                               `json:"tspContributionPercentagePreRetirement,omitempty"`
	TSPContributeUntil                          *string                                `json:"tspContributeUntil,omitempty"`
	TSPContributionStopAge                      *int                                   `json:"tspContributionStopAge,omitempty"`
	TSPCatchUpContributionsEligible             bool                                   `json:"tspCatchUpContributionsEligible,omitempty"`
	TSPContributionAllocationTraditionalPercent *float64                               `json:"tspContributionAllocationTraditionalPercent,omitempty"`
	TSPContributionAllocationRothPercent        *float64                               `json:"tspContributionAllocationRothPercent,omitempty"`
	TSPContributionAllocationToFunds            *models.TSPFundAllocation              `json:"tspContributionAllocationToFunds,omitempty"`
	TSPCurrentAllocationToFunds                 *models.TSPFundAllocation              `json:"tspCurrentAllocationToFunds,omitempty"`
	TSPPostRetirementAllocation                 *models.TSPFundAllocation              `json:"tspPostRetirementAllocation,omitempty"`
	TSPFutureAllocationStrategy                 *string                                `json:"tspFutureAllocationStrategy,omitempty"` // "MaintainCurrent", "ReallocateTo", "LifecycleFund", "UsePostRetirementAllocation"
	TSPFutureAllocationTarget                   *models.TSPFundAllocation              `json:"tspFutureAllocationTarget,omitempty"`   // For 'ReallocateTo' strategy
	TSPLifecycleFundTargetDate                  *string                                `json:"tspLifecycleFundTargetDate,omitempty"`  // YYYY, e.g. "2050" for L2050 if 'LifecycleFund' strategy
	UserReturnAssumptionsTSP                    *models.TSPIndividualReturnAssumptions `json:"userReturnAssumptionsTSP,omitempty"`
	TSPExpenseRatio                             *float64                               `json:"tspExpenseRatio,omitempty"`
	TSPExpectedAnnualGrowthRatePreRetirement    *float64                               `json:"tspExpectedAnnualGrowthRatePreRetirement,omitempty"`
	TSPWithdrawalStrategy                       *string                                `json:"tspWithdrawalStrategy,omitempty"`
	TSPWithdrawalFixedAmountValue               *float64                               `json:"tspWithdrawalFixedAmountValue,omitempty"`
	TSPWithdrawalPercentageValue                *float64                               `json:"tspWithdrawalPercentageValue,omitempty"`
	TSPWithdrawalStartAge                       *int                                   `json:"tspWithdrawalStartAge,omitempty"`
	TSPWithdrawalStartDate                      *string                                `json:"tspWithdrawalStartDate,omitempty"`
	TSPWithdrawalOrder                          *string                                `json:"tspWithdrawalOrder,omitempty"`
	ApplyRMDsToTSP                              bool                                   `json:"applyRMDsToTSP,omitempty"`
	TSPExpectedAnnualGrowthRatePostRetirement   *float64                               `json:"tspExpectedAnnualGrowthRatePostRetirement,omitempty"`
	TSPVolatilityAssumptions                    *models.TSPVolatilityRates             `json:"tspVolatilityAssumptions,omitempty"`

	// TSP Annuity Details
	TSPAnnuityType              *string  `json:"tspAnnuityType,omitempty"` // 'SingleLife' | 'JointLifeSpouse' | 'JointLifeNonSpouse'
	TSPAnnuitySpouseAge         *int     `json:"tspAnnuitySpouseAge,omitempty"`
	TSPAnnuityFeatures          []string `json:"tspAnnuityFeatures,omitempty"` // ('IncreasingPayments' | 'CashRefund' | 'TenYearCertain')
	TSPAnnuityAmountForPurchase *float64 `json:"tspAnnuityAmountForPurchase,omitempty"`
	TSPAnnuityPurchaseAge       *int     `json:"tspAnnuityPurchaseAge,omitempty"`

	// Fields from InsuranceScenarioInputs
	FegliBasicCoverage               bool     `json:"fegliBasicCoverage,omitempty"`
	FegliOptionAAmount               *float64 `json:"fegliOptionAAmount,omitempty"`
	FegliOptionBCoverageMultiples    *int     `json:"fegliOptionBCoverageMultiples,omitempty"`
	FegliOptionCCoverageMultiples    *int     `json:"fegliOptionCCoverageMultiples,omitempty"`
	FegliPost65Reduction             *string  `json:"fegliPost65Reduction,omitempty"`
	FegliBasicRetiredCoverage        *string  `json:"fegliBasicRetiredCoverage,omitempty"`
	FegliOptionARetiredCoverage      *string  `json:"fegliOptionARetiredCoverage,omitempty"`
	FegliOptionBRetiredCoverage      *string  `json:"fegliOptionBRetiredCoverage,omitempty"`
	FegliOptionCRetiredCoverage      *string  `json:"fegliOptionCRetiredCoverage,omitempty"`
	FehbPlanName                     *string  `json:"fehbPlanName,omitempty"`
	FehbPlanCode                     *string  `json:"fehbPlanCode,omitempty"`
	FehbCoverageType                 *string  `json:"fehbCoverageType,omitempty"`
	FehbPremiumPreRetirementMonthly  *float64 `json:"fehbPremiumPreRetirementMonthly,omitempty"`
	FehbPremiumPostRetirementMonthly *float64 `json:"fehbPremiumPostRetirementMonthly,omitempty"`
	FehbContinueInRetirement         bool     `json:"fehbContinueInRetirement,omitempty"`
	FehbBiweeklyDeduction            *float64 `json:"fehbBiweeklyDeduction,omitempty"`
	FehbContinuedInRetirement        bool     `json:"fehbContinuedInRetirement,omitempty"`
	FehbProjectedPremiumIncreaseRate *float64 `json:"fehbProjectedPremiumIncreaseRate,omitempty"`
	UserAssumedFehbPremiumGrowthRate *float64 `json:"userAssumedFehbPremiumGrowthRate,omitempty"`
	FltcipEnrolled                   bool     `json:"fltcipEnrolled,omitempty"`
	FltcipCurrentStatus              *string  `json:"fltcipCurrentStatus,omitempty"` // 'enrolledPaying' | 'enrolledPaidUp' | 'notEnrolled'
	FltcipDailyBenefitAmount         *float64 `json:"fltcipDailyBenefitAmount,omitempty"`
	FltcipBenefitPeriod              *string  `json:"fltcipBenefitPeriod,omitempty"`
	FltcipInflationProtection        *string  `json:"fltcipInflationProtection,omitempty"`
	FltcipWaitingPeriod              *string  `json:"fltcipWaitingPeriod,omitempty"` // e.g., '30days', '90days'
	FltcipPremiumMonthly             *float64 `json:"fltcipPremiumMonthly,omitempty"`
	FltcipContinueInRetirement       bool     `json:"fltcipContinueInRetirement,omitempty"`

	// FEDVIP Details
	FedvipDentalPlanName *string  `json:"fedvipDentalPlanName,omitempty"`
	FedvipDentalPremium  *float64 `json:"fedvipDentalPremium,omitempty"` // Monthly
	FedvipVisionPlanName *string  `json:"fedvipVisionPlanName,omitempty"`
	FedvipVisionPremium  *float64 `json:"fedvipVisionPremium,omitempty"` // Monthly

	OtherLifeInsuranceCoverageAmount   *float64 `json:"otherLifeInsuranceCoverageAmount,omitempty"`
	OtherLifeInsurancePremiumAnnual    *float64 `json:"otherLifeInsurancePremiumAnnual,omitempty"`
	DisabilityInsuranceCoverageAmount  *float64 `json:"disabilityInsuranceCoverageAmount,omitempty"`
	DisabilityInsurancePremiumAnnual   *float64 `json:"disabilityInsurancePremiumAnnual,omitempty"`
	LongTermCareInsurancePolicyDetails *string  `json:"longTermCareInsurancePolicyDetails,omitempty"`

	// Fields from TaxScenarioInputs
	FederalTaxFilingStatus                *string  `json:"federalTaxFilingStatus,omitempty"`
	FederalTaxNumberOfDependents          *int     `json:"federalTaxNumberOfDependents,omitempty"`
	FederalTaxExtraWithholding            *float64 `json:"federalTaxExtraWithholding,omitempty"`
	FederalTaxAdjustmentsToIncome         *float64 `json:"federalTaxAdjustmentsToIncome,omitempty"`
	FederalTaxItemizedDeductions          *float64 `json:"federalTaxItemizedDeductions,omitempty"`
	FederalOtherTaxableIncomeAnnual       *float64 `json:"federalOtherTaxableIncomeAnnual,omitempty"`
	FederalTaxCreditsAnnual               *float64 `json:"federalTaxCreditsAnnual,omitempty"`
	FederalTaxLawAssumption               *string  `json:"federalTaxLawAssumption,omitempty"`         // 'CurrentLaw' | 'ExtendsTCJA' | 'Custom'
	UserAssumedFederalTaxRateChange       *string  `json:"userAssumedFederalTaxRateChange,omitempty"` // 'none' | 'increase' | 'decrease' | 'custom'
	FederalTaxUserAssumedRateChangeAmount *float64 `json:"federalTaxUserAssumedRateChangeAmount,omitempty"`
	FederalTaxUserAssumedRateChangeYear   *int     `json:"federalTaxUserAssumedRateChangeYear,omitempty"`
	StateOfResidenceForTax                *string  `json:"stateOfResidenceForTax,omitempty"`
	RetirementEffectiveStateTaxRate       *float64 `json:"retirementEffectiveStateTaxRate,omitempty"`
	StateTaxFilingStatus                  *string  `json:"stateTaxFilingStatus,omitempty"`
	StateTaxNumberOfDependents            *int     `json:"stateTaxNumberOfDependents,omitempty"`
	StateTaxExtraWithholding              *float64 `json:"stateTaxExtraWithholding,omitempty"`
	StateTaxAdjustmentsToIncome           *float64 `json:"stateTaxAdjustmentsToIncome,omitempty"`
	StateTaxItemizedDeductions            *float64 `json:"stateTaxItemizedDeductions,omitempty"`
	UserAssumedStateTaxRateChange         *string  `json:"userAssumedStateTaxRateChange,omitempty"` // 'none' | 'increase' | 'decrease' | 'custom'
	StateTaxUserAssumedRateChangeAmount   *float64 `json:"stateTaxUserAssumedRateChangeAmount,omitempty"`
	StateTaxUserAssumedRateChangeYear     *int     `json:"stateTaxUserAssumedRateChangeYear,omitempty"`
	TaxTreatmentOfMilitaryRetirementPay   *string  `json:"taxTreatmentOfMilitaryRetirementPay,omitempty"` // 'fullyTaxable' | 'partiallyTaxable' | 'notTaxable' | 'na'

	// Fields from OtherFinancialsScenarioInputs
	OtherTaxableAccountBalance                  *float64                              `json:"otherTaxableAccountBalance,omitempty"`
	OtherTaxableAccountAllocation               *models.TaxableAccountAssetAllocation `json:"otherTaxableAccountAllocation,omitempty"`
	OtherTaxableAccountAnnualContribution       *float64                              `json:"otherTaxableAccountAnnualContribution,omitempty"`
	OtherTaxableAccountExpectedGrowth           *float64                              `json:"otherTaxableAccountExpectedGrowth,omitempty"`
	OtherTaxAdvantagedAccountBalance            *float64                              `json:"otherTaxAdvantagedAccountBalance,omitempty"`
	OtherTaxAdvantagedAccountExpectedGrowth     *float64                              `json:"otherTaxAdvantagedAccountExpectedGrowth,omitempty"`
	OtherTaxAdvantagedAccountAnnualContribution *float64                              `json:"otherTaxAdvantagedAccountAnnualContribution,omitempty"`
	OtherRecurringIncomeSources                 []models.OtherRecurringIncomeSource   `json:"otherRecurringIncomeSources,omitempty"`
	OtherRecurringMonthlyExpensesPreRetirement  *float64                              `json:"otherRecurringMonthlyExpensesPreRetirement,omitempty"`
	OtherRecurringMonthlyExpensesPostRetirement *float64                              `json:"otherRecurringMonthlyExpensesPostRetirement,omitempty"`
	OneTimeIncomeEvents                         []models.OneTimeIncomeEvent           `json:"oneTimeIncomeEvents,omitempty"`
	OneTimeExpenseEvents                        []models.OneTimeExpenseEvent          `json:"oneTimeExpenseEvents,omitempty"`
	PrimaryHomeValue                            *float64                              `json:"primaryHomeValue,omitempty"`
	PrimaryHomeMortgageBalance                  *float64                              `json:"primaryHomeMortgageBalance,omitempty"`
	PrimaryHomeMortgageInterestRate             *float64                              `json:"primaryHomeMortgageInterestRate,omitempty"`
	PrimaryHomeMortgageRemainingTermYears       *float64                              `json:"primaryHomeMortgageRemainingTermYears,omitempty"`
	PrimaryHomePropertyTaxRate                  *float64                              `json:"primaryHomePropertyTaxRate,omitempty"`
	PrimaryHomeInsuranceRate                    *float64                              `json:"primaryHomeInsuranceRate,omitempty"`
	PrimaryHomeHoaMonthly                       *float64                              `json:"primaryHomeHoaMonthly,omitempty"`
	SellPrimaryHomeInRetirement                 bool                                  `json:"sellPrimaryHomeInRetirement,omitempty"`
	SellPrimaryHomeAge                          *int                                  `json:"sellPrimaryHomeAge,omitempty"`
	DownsizeOrMoveInRetirement                  bool                                  `json:"downsizeOrMoveInRetirement,omitempty"`
	DownsizeOrMoveAge                           *int                                  `json:"downsizeOrMoveAge,omitempty"`
	EstimatedNetProceedsFromSale                *float64                              `json:"estimatedNetProceedsFromSale,omitempty"`
	CostOfNewHomeOrRental                       *float64                              `json:"costOfNewHomeOrRental,omitempty"` // Annual if renting, total if buying new
	OtherMajorAssets                            *string                               `json:"otherMajorAssets,omitempty"`
	OtherMajorLiabilities                       *string                               `json:"otherMajorLiabilities,omitempty"`

	// Fields from GeneralAssumptions
	GeneralInflationRate            *float64 `json:"generalInflationRate,omitempty"`
	UserAssumedGeneralInflationRate *float64 `json:"userAssumedGeneralInflationRate,omitempty"`
	UserAssumedPensionCOLA_FERS     *float64 `json:"userAssumedPensionCOLA_FERS,omitempty"`
	UserAssumedPensionCOLA_CSRS     *float64 `json:"userAssumedPensionCOLA_CSRS,omitempty"`
	COLAAnnuity                     *float64 `json:"colaAnnuity,omitempty"`
	COLASocialSecurity              *float64 `json:"colaSocialSecurity,omitempty"`
	WageInflationRate               *float64 `json:"wageInflationRate,omitempty"` // For projecting future salary/High-3
	UserLifeExpectancy              *int     `json:"userLifeExpectancy,omitempty"`
	SpouseDateOfBirth               *string  `json:"spouseDateOfBirth,omitempty"`
	SpouseLifeExpectancy            *int     `json:"spouseLifeExpectancy,omitempty"`
	MonteCarloEnabled               bool     `json:"monteCarloEnabled,omitempty"`
	MonteCarloNumSimulations        *int     `json:"monteCarloNumSimulations,omitempty"`
	MonteCarloConfidenceLevel       *float64 `json:"monteCarloConfidenceLevel,omitempty"`

	// Retirement Spending Goals
	DesiredAnnualSpendingInRetirement *float64 `json:"desiredAnnualSpendingInRetirement,omitempty"`
	IsSpendingGoalInflationAdjusted   *bool    `json:"isSpendingGoalInflationAdjusted,omitempty"`

	// Healthcare Cost Assumptions
	HealthcareCostInflationRate *float64 `json:"healthcareCostInflationRate,omitempty"`
	UserLongTermCareNeed        *string  `json:"userLongTermCareNeed,omitempty"`   // 'none' | 'low' | 'medium' | 'high'
	SpouseLongTermCareNeed      *string  `json:"spouseLongTermCareNeed,omitempty"` // 'none' | 'low' | 'medium' | 'high'
	LongTermCareCostStartAge    *int     `json:"longTermCareCostStartAge,omitempty"`
}

// FerexFile is the top-level structure saved to a .ferex file.
type FerexFile struct {
	FileFormatVersion    string            `json:"fileFormatVersion"` // e.g., "1.0"
	LastOpenedAppVersion string            `json:"lastOpenedAppVersion,omitempty"`
	UserProfile          UserProfile       `json:"userProfile"`
	Variants             []ScenarioVariant `json:"variants"`
	LastViewedVariantID  string            `json:"lastViewedVariantID,omitempty"`
}
