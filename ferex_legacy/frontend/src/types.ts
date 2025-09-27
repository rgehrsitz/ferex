import type { models } from "../wailsjs/go/models";

// Type alias for clarity
export type RetirementCalculationResult = models.RetirementCalculationResult;

// Helper interface for a single period of service
export interface ServicePeriod {
  id: string; // auto-generated unique identifier
  serviceCategory: "Civilian" | "Military"; // Broad category of service
  civilianServiceType: string | null; // Applicable if serviceCategory is 'Civilian'. e.g., 'RegularDeductionFERS', 'NonDeductionPre10_82CSRS', etc.
  militaryServiceType: string | null; // Applicable if serviceCategory is 'Military'. e.g., 'Post1956MilitaryServiceDepositPaid', 'Post1956MilitaryServiceDepositOwed'.
  depositRedepositPaymentStatus:
    | "PaidInFull"
    | "OwedOrPartiallyPaid"
    | "NotApplicable"
    | "AwaitingDetermination"; // Status of any required deposits or redeposits
  systemDuringService:
    | "FERS"
    | "CSRS"
    | "CSRS_Offset"
    | "SocialSecurityOnly"
    | "None"
    | "Other"; // The retirement system the employee was covered by during this specific period of service.
  startDate: string; // YYYY-MM-DD format, start date of the service period
  endDate: string; // YYYY-MM-DD format, end date of the service period
  isPartTime: boolean; // Indicates if the service was performed on a part-time basis.
  hoursPerWeekIfPartTime: number | null; // Number of hours worked per week if isPartTime is true, based on a standard full-time work week, e.g., 20 for half-time if standard is 40 hours.
  notes: string | null; // Optional user-provided notes for this service period.
}

// Helper interface for Leave Without Pay (LWOP) periods
export interface LWOPPeriod {
  id: string; // For UI list management
  startDate: string; // YYYY-MM-DD
  endDate: string; // YYYY-MM-DD
  type: "PersonalNonMilitary" | "MilitaryLWOP" | "OWCP"; // Type of LWOP
}

// Helper interface for Insurable Interest details
export interface InsurableInterestDetails {
  relationshipToEmployee: string; // e.g., Parent, Sibling, Financially Dependent Child
  dateOfBirth: string; // YYYY-MM-DD
}

// TSP Fund Allocation Helper Types
export interface TSPFundAllocation {
  G?: number | null;
  F?: number | null;
  C?: number | null;
  S?: number | null;
  I?: number | null;
  LFundName?: string | null; // e.g., 'L2050', 'LIncome'
}

export interface TSPReturnAssumptions {
  G?: number | null;
  F?: number | null;
  C?: number | null;
  S?: number | null;
  I?: number | null;
}

// Helper interface for TSP Fund Volatility (Standard Deviations)
export interface TSPVolatility {
  G_stdDev?: number | null;
  F_stdDev?: number | null;
  C_stdDev?: number | null;
  S_stdDev?: number | null;
  I_stdDev?: number | null;
}

// Helper interface for asset allocation in other taxable accounts
export interface TaxableAccountAssetAllocation {
  stocksPercent?: number | null;
  bondsPercent?: number | null;
  cashPercent?: number | null;
}

// Helper interface for Social Security benefit table (from SS statement)
export interface SocialSecurityBenefitEntry {
  age: number;
  monthlyBenefit: number | null;
}

// Helper interface for Other Recurring Income Sources
export interface OtherRecurringIncomeSource {
  id: string;
  name: string;
  amount: number | null;
  frequency: "Monthly" | "Annually";
  startDate: string; // YYYY-MM-DD
  endDate?: string | null; // YYYY-MM-DD, optional
  isInflationAdjusted?: boolean;
}

// Helper interface for One-Time Income Events
export interface OneTimeIncomeEvent {
  id: string;
  name: string;
  amount: number | null;
  date: string; // YYYY-MM-DD
  type: "Inheritance" | "SaleOfAsset" | "Bonus" | "Other";
}

// Helper interface for One-Time Expense Events
export interface OneTimeExpenseEvent {
  id: string;
  name: string;
  amount: number | null;
  date: string; // YYYY-MM-DD
  type:
    | "HomePurchase"
    | "VehiclePurchase"
    | "Education"
    | "Travel"
    | "Medical"
    | "Other";
}

// Defines the structure for all user-configurable inputs for a retirement scenario.
export interface TSPFundVolatilityRates {
  G_stdDev?: number | null;
  F_stdDev?: number | null;
  C_stdDev?: number | null;
  S_stdDev?: number | null;
  I_stdDev?: number | null;
}

export interface ScenarioInput {
  calculationSystem: "FERS" | "CSRS" | "";
  calculationYears?: number | null; // Global setting to restrict all calculations to N years (for debugging)
  high3Salary: number | null;
  serviceComputationDate: string; // YYYY-MM-DD
  notes?: string; // User notes for the scenario variant
  createdAt?: string; // ISO date string, when the variant was first created
  lastModified?: string; // ISO date string, when the variant was last saved
  dateOfBirth: string; // YYYY-MM-DD
  plannedRetirementDate: string; // YYYY-MM-DD
  unusedSickLeaveHoursAtRetirement: number | null;
  employeeContributions: number | null;

  // Creditable Service Periods
  servicePeriods?: ServicePeriod[];
  lwopPeriods?: LWOPPeriod[]; // Leave Without Pay periods

  // Retirement Eligibility / Special Conditions
  fersCoverageType?: "FERS" | "FERS_RAE" | "FERS_FRAE" | null; // FERS specific coverage variant
  isVeraDsRetirement?: boolean; // VERA/Deferred Service Retirement
  mraPlus10PostponeAnnuityStartDate?: string | null; // YYYY-MM-DD, for MRA+10 postponed annuity
  isDeferredRetirement?: boolean;
  deferredRetirementAnnuityStartDate?: string | null; // YYYY-MM-DD
  isDisabilityRetirement?: boolean;

  // FERS Specific
  estimatedSSBenefitAt62ForSRS?: number | null;
  didSwitchFromCSRS?: boolean;
  switchedToFERSDate?: string | null; // YYYY-MM-DD
  survivorBenefitFERS?: SurvivorBenefitFERSInput | null;

  // CSRS Specific
  isCSRSOffset: boolean;
  yearsOfOffsetService: number | null;
  ssAt62WithOffset: number | null;
  ssAt62WithoutOffset: number | null;
  survivorBenefitCSRS?: SurvivorBenefitCSRSInput | null;

  // Survivor Benefits / Spouse Details
  isMarriedAtRetirement?: boolean;
  spouseBirthDate?: string; // YYYY-MM-DD, required if isMarriedAtRetirement is true
  ageOfInsurableInterestBeneficiary?: number | null; // Relevant if InsurableInterest is chosen for FERS
  hasFormerSpouseEntitlement?: boolean;
  formerSpouseSurvivorBenefitDetails?: models.FormerSpouseSurvivorDetails | null;

  // Social Security & TSP Tab
  userProvidedSSBenefitAmount1?: number | null;
  userProvidedSSBenefitClaimingAge1?: number | null;
  userProvidedSSBenefitAmount2?: number | null;
  userProvidedSSBenefitClaimingAge2?: number | null;
  ssBenefitSpousalOption?:
    | "None"
    | "OwnRecord"
    | "SpousalDefault"
    | "SpousalMax";
  ssBenefitSpousalAmount?: number | null;
  ssBenefitSurvivorOption?:
    | "None"
    | "OwnRecord"
    | "SurvivorDefault"
    | "SurvivorMax";
  ssBenefitSurvivorAmount?: number | null;
  userSSClaimingAge?: number | null; // User's own SS claiming age if different from default
  userAssumedSSCOLA?: number | null; // User's assumption for SS COLA, percentage
  calculateHistoricalWEPGPO?: boolean; // Option to use historical WEP/GPO calculation methods
  socialSecurityEstimate?: number | null; // Monthly PIA at FRA
  socialSecurityFRA?: number | null; // Full Retirement Age (e.g., 67, 66.5)
  socialSecurityBenefitStartDate?: string | null; // YYYY-MM-DD or 'FRA'
  socialSecurityCreditedEarnings?: { year: number; earnings: number }[]; // For detailed SS calculation
  socialSecurityBenefitTable?: SocialSecurityBenefitEntry[]; // Benefits by claiming age from SS statement

  // TSP Fund Allocation Helper Types
  tspBalanceTraditional?: number | null;
  tspBalanceRoth?: number | null;
  tspLoanBalance?: number | null;
  tspAnnualContributionPreRetirement?: number | null;
  tspContributionPercentagePreRetirement?: number | null;
  tspContributeUntil?: "Retirement" | "SpecificAge"; // When TSP contributions stop
  tspContributionStopAge?: number | null; // Specific age to stop TSP contributions
  tspCatchUpContributionsEligible?: boolean; // Eligible for catch-up contributions
  tspContributionAllocationTraditionalPercent?: number | null; // Percentage to Traditional TSP
  tspContributionAllocationRothPercent?: number | null; // Percentage to Roth TSP
  tspContributionAllocationToFunds?: TSPFundAllocation | null;
  tspCurrentAllocationToFunds?: TSPFundAllocation | null;
  tspPostRetirementAllocation?: TSPFundAllocation | null; // Allocation to use after retirement
  tspFutureAllocationStrategy?:
    | "MaintainCurrent"
    | "MatchContributionAllocation"
    | "CustomFutureAllocation"
    | "UsePostRetirementAllocation";
  userReturnAssumptionsTSP?: TSPReturnAssumptions | null;
  tspExpenseRatio?: number | null;
  tspExpectedAnnualGrowthRatePreRetirement?: number | null;
  tspWithdrawalStrategy?:
    | "None"
    | "FixedAmountYearly"
    | "FixedAmountMonthly"
    | "PercentageOfBalanceYearly"
    | "IRSMinimumRequiredDistribution"
    | "AsNeededToMeetIncomeGoal"
    | null;
  tspWithdrawalFixedAmountValue?: number | null; // For 'FixedAmountYearly' or 'FixedAmountMonthly' strategies
  tspWithdrawalPercentageValue?: number | null; // For 'PercentageOfBalanceYearly' strategy
  tspWithdrawalStartAge?: number | null;
  tspWithdrawalStartDate?: "Retirement" | "SpecificAge";
  tspWithdrawalOrder?:
    | "ProRata"
    | "TraditionalFirst"
    | "RothFirst"
    | "CustomSequence"; // Order of TSP withdrawals
  applyRMDsToTSP?: boolean; // Whether Required Minimum Distributions apply to TSP

  // COLA / Inflation Assumptions
  pensionAnnuityColaRate?: number | null; // Assumed annual COLA for FERS/CSRS pension, percentage
  tspExpectedAnnualGrowthRatePostRetirement?: number | null;
  tspVolatilityAssumptions?: TSPVolatility | null;

  // Insurance Tab
  // FEGLI (Federal Employees' Group Life Insurance) Fields
  /** Whether the employee has FEGLI Basic coverage */
  fegliBasicCoverage?: boolean;
  /** Amount of Option A - Standard Insurance coverage in dollars */
  fegliOptionAAmount?: number | null;
  /** Number of multiples for Option B - Additional Insurance */
  fegliOptionBCoverageMultiples?: number | null;
  /** Number of multiples for Option C - Family Insurance */
  fegliOptionCCoverageMultiples?: number | null;
  /** Type of reduction for FEGLI coverage after age 65 */
  fegliPost65Reduction?: "FullReduction" | "NoReduction" | "PartialReduction75";
  /** Whether Basic insurance is continued into retirement */
  fegliBasicHeldIntoRetirement?: boolean;
  /** Type of reduction for Basic insurance after retirement */
  fegliBasicPostRetirementReductionChoice?:
    | "75Reduce"
    | "50Reduce"
    | "NoReduce";

  // FEHB (Federal Employees Health Benefits) Fields
  /** Biweekly FEHB deduction amount from LES */
  fehbBiweeklyDeduction?: number | null;
  /** Projected annual premium increase rate as a percentage */
  fehbProjectedPremiumIncreaseRate?: number | null;
  /** Whether to continue FEHB coverage into retirement */
  fehbContinuedInRetirement?: boolean;

  // FLTCIP (Federal Long Term Care Insurance Program) Fields
  /** Whether the employee has FLTCIP coverage */
  fltcipHasCoverage?: boolean;
  /** Daily benefit amount for FLTCIP coverage */
  fltcipDailyBenefitAmount?: number | null;
  /** Duration of the FLTCIP benefit period */
  fltcipBenefitPeriod?: "2years" | "3years" | "5years" | "unlimited";
  /** Type of inflation protection for FLTCIP coverage */
  fltcipInflationProtectionType?: "Automatic" | "FuturePurchase" | "None";
  /** Current annual premium for FLTCIP coverage */
  fltcipCurrentAnnualPremium?: number | null;

  // Financial & Tax Planning Tab
  // General Inflation & COLA Assumptions
  colaAnnuity?: number | null; // Assumed COLA for FERS/CSRS annuity, percentage
  colaSocialSecurity?: number | null; // Assumed COLA for Social Security, percentage
  userAssumedGeneralInflationRate?: number | null; // General inflation rate, percentage
  userAssumedPensionCOLA_FERS?: number | null; // Specific FERS COLA assumption, overrides general if provided
  userAssumedPensionCOLA_CSRS?: number | null; // Specific CSRS COLA assumption, overrides general if provided

  // Other Income & Assets
  otherIncomeMonthly?: number | null; // Simplified single monthly other income field (legacy, consider removing or migrating)
  otherRecurringIncomeSources?: OtherRecurringIncomeSource[];
  oneTimeIncomeEvents?: OneTimeIncomeEvent[];
  otherTaxableAccountBalance?: number | null;
  otherTaxableAccountAllocation?: TaxableAccountAssetAllocation | null;
  otherTaxableAccountExpectedGrowth?: number | null;

  // Retirement Spending Goals
  desiredAnnualSpendingInRetirement?: number | null; // In today's dollars, pre-tax
  isSpendingGoalInflationAdjusted?: boolean; // Default true

  // TSP Volatility
  volatilityTSPFunds?: TSPFundVolatilityRates | null; // Standard deviations for TSP funds

  // Home, Liabilities & Major Expenses
  homeValue?: number | null;
  mortgageBalance?: number | null;
  otherSignificantDebtsTotal?: number | null;
  oneTimeExpenseEvents?: OneTimeExpenseEvent[];

  // Retirement Spending & Legacy Goals
  desiredRetirementIncomeYearly?: number | null; // Target annual income in retirement
  legacyGoalAmount?: number | null; // Amount to leave as legacy

  // Federal Tax Assumptions - Using backend field names
  federalTaxRate?: number | null; // Effective federal tax rate percentage
  federalTaxFilingStatus?:
    | "Single"
    | "MarriedFilingJointly"
    | "MarriedFilingSeparately"
    | "HeadOfHousehold"
    | "QualifyingWidow(er)";
  federalTaxNumberOfDependents?: number | null;
  federalTaxItemizedDeductions?: number | null;
  federalOtherTaxableIncomeAnnual?: number | null;
  federalTaxCreditsAnnual?: number | null;
  federalTaxLawAssumption?: "CurrentLaw" | "ExtendsTCJA" | "Custom";
  stateOfResidenceForTax?: string; // e.g., 'VA', 'MD', 'DC'
  retirementEffectiveStateTaxRate?: number | null;

  // Spouse & Household Planning (if applicable)
  planAsHousehold?: boolean; // If true, spouse's details below are considered
  spouseDateOfBirth?: string | null; // YYYY-MM-DD
  spousePlannedRetirementDate?: string | null; // YYYY-MM-DD
  spouseProvidedSSBenefitAmount1?: number | null;
  spouseProvidedSSBenefitClaimingAge1?: number | null;
  spouseSSClaimingAge?: number | null;
  spouseHasPension?: boolean;
  spousePensionAnnualAmount?: number | null;
  spousePensionCOLA?: number | null;
  spouseOtherRetirementAccountBalance?: number | null;
  spouseOtherRetirementAccountAnnualWithdrawal?: number | null;

  // Monte Carlo Simulation Settings
  monteCarloNumberOfIterations?: number | null;
  monteCarloSimulationEndAge?: number | null;
  monteCarloSuccessDefinition?:
    | "AssetsRemainPositive"
    | "MeetSpendingGoals"
    | "Custom";
  monteCarloMinEssentialSpendingLevel?: number | null; // Used if success definition is 'MeetSpendingGoals' or 'Custom'
}

export interface SurvivorBenefitFERSInput {
  spouseElection?:
    | "Full50Percent"
    | "Partial25Percent"
    | "None"
    | "InsurableInterest"
    | null;
  formerSpouseElection?:
    | "Full50Percent"
    | "Partial25Percent"
    | "None"
    | "InsurableInterest"
    | null;
  formerSpouseConsent?: "Yes" | "No" | "NotApplicable" | null; // Consent from former spouse if their benefit is reduced/eliminated
  currentSpouseConsentForFormer?: "Yes" | "No" | "NotApplicable" | null; // Consent from current spouse if former spouse gets benefit
  currentSpouseWaiverForSelf?: "Yes" | "No" | "NotApplicable" | null; // Current spouse waives their own right to a survivor annuity
  insurableInterestDetails?: InsurableInterestDetails | null;
}

export interface SurvivorBenefitCSRSInput {
  election?:
    | "Full55PercentMax"
    | "PartialCustomBase"
    | "None"
    | "InsurableInterest"
    | null;
  customBaseAmountForPartial?: number | null; // Dollar amount for partial survivor annuity base
  // Mirroring FERS for former spouse, assuming court order dictates
  formerSpouseElection?:
    | "Full55PercentMax"
    | "PartialCustomBase"
    | "None"
    | "InsurableInterest"
    | null;
  formerSpouseCustomBaseAmount?: number | null;
  formerSpouseConsent?: "Yes" | "No" | "NotApplicable" | null;
  currentSpouseConsentForFormer?: "Yes" | "No" | "NotApplicable" | null;
  currentSpouseWaiverForSelf?: "Yes" | "No" | "NotApplicable" | null;
  insurableInterestDetails?: InsurableInterestDetails | null;
}

// Defines the structure for a single retirement planning scenario.
export interface Scenario {
  id: string;
  name: string;
  inputs: ScenarioInput;
  results: RetirementCalculationResult | null;
  error: string | null;
  calculating: boolean;
}
