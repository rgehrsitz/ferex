import type { ScenarioInput, RetirementCalculationResult } from "../../types";
import { validateScenarioInputs } from "../scenarioValidation";
import { TSPService } from "./tspService";
import { PerformRetirementCalculation } from "../../../wailsjs/go/main/App";

export interface CalculationServiceResult {
  results: RetirementCalculationResult | null;
  tspProjections: any | null;
  error: string | null;
  validationMessages: string[];
}

// Helper function to convert percentage values (like 7.5) to decimal values (like 0.075)
function convertPercentagesToDecimals(percentageObj: any): any {
  if (!percentageObj) return {};

  const result: any = {};
  for (const [key, value] of Object.entries(percentageObj)) {
    if (typeof value === "number" && value !== null) {
      result[key] = value / 100; // Convert 7.5 to 0.075
    } else {
      result[key] = value; // Keep non-number values as-is
    }
  }
  return result;
}

// Helper function to convert TSPFundAllocation to TSPFundAllocationPercentages format
function convertTSPAllocation(allocation: any): any {
  if (!allocation) return {};

  const result: any = {};
  const fields = ["G", "F", "C", "S", "I"];

  for (const field of fields) {
    const value = allocation[field];
    if (value !== null && value !== undefined) {
      result[field.toLowerCase()] = value;
    }
  }

  if (allocation.LFundName) {
    result.lFundName = allocation.LFundName;
  }

  return result;
}

// Helper function to calculate age from date of birth
function calculateAge(dateOfBirth: string): number {
  if (!dateOfBirth) return 0;

  const birthDate = new Date(dateOfBirth);
  const today = new Date();
  let age = today.getFullYear() - birthDate.getFullYear();
  const monthDiff = today.getMonth() - birthDate.getMonth();

  if (
    monthDiff < 0 ||
    (monthDiff === 0 && today.getDate() < birthDate.getDate())
  ) {
    age--;
  }

  return age;
}

// Helper function to calculate age at retirement
function calculateRetirementAge(
  dateOfBirth: string,
  retirementDate: string
): number {
  if (!dateOfBirth || !retirementDate) return 0;

  const birthDate = new Date(dateOfBirth);
  const retireDate = new Date(retirementDate);
  let age = retireDate.getFullYear() - birthDate.getFullYear();
  const monthDiff = retireDate.getMonth() - birthDate.getMonth();

  if (
    monthDiff < 0 ||
    (monthDiff === 0 && retireDate.getDate() < birthDate.getDate())
  ) {
    age--;
  }

  return age;
}

// Helper function to calculate birth year
function getBirthYear(dateOfBirth: string): number {
  if (!dateOfBirth) return 0;
  return parseInt(dateOfBirth.split("-")[0], 10);
}

// Helper function to calculate total other taxable income
function calculateOtherTaxableIncome(inputs: ScenarioInput): number {
  let totalOtherIncome = 0;

  // Add federal other taxable income
  totalOtherIncome += inputs.federalOtherTaxableIncomeAnnual || 0;

  // Add recurring income sources
  if (inputs.otherRecurringIncomeSources) {
    for (const source of inputs.otherRecurringIncomeSources) {
      const amount = source.amount || 0;
      // Convert to annual amount based on frequency
      if (source.frequency === "Monthly") {
        totalOtherIncome += amount * 12;
      } else {
        totalOtherIncome += amount;
      }
    }
  }

  // Add one-time income events for the current year (simplified - using current year)
  const currentYear = new Date().getFullYear();
  if (inputs.oneTimeIncomeEvents) {
    for (const event of inputs.oneTimeIncomeEvents) {
      const eventYear = parseInt(event.date.split("-")[0]);
      if (eventYear === currentYear) {
        totalOtherIncome += event.amount || 0;
      }
    }
  }

  return totalOtherIncome;
}

// Helper function to calculate years of FERS service for SRS
function calculateFERSYearsForSRS(inputs: ScenarioInput): number {
  if (!inputs.servicePeriods || inputs.servicePeriods.length === 0) {
    return 0;
  }

  let fersYears = 0;
  for (const period of inputs.servicePeriods) {
    if (period.systemDuringService === "FERS") {
      const startDate = new Date(period.startDate);
      const endDate = new Date(period.endDate);
      const years =
        (endDate.getTime() - startDate.getTime()) /
        (365.25 * 24 * 60 * 60 * 1000);
      fersYears += years;
    }
  }

  return fersYears;
}

// Helper function to determine if immediate unreduced annuity
function isImmediateUnreducedAnnuity(inputs: ScenarioInput): boolean {
  const retirementAge = calculateRetirementAge(
    inputs.dateOfBirth,
    inputs.plannedRetirementDate
  );
  const mra = getMRA(inputs.dateOfBirth);

  // Calculate total service years
  let totalServiceYears = 0;
  if (inputs.servicePeriods) {
    for (const period of inputs.servicePeriods) {
      const startDate = new Date(period.startDate);
      const endDate = new Date(period.endDate);
      const years =
        (endDate.getTime() - startDate.getTime()) /
        (365.25 * 24 * 60 * 60 * 1000);
      totalServiceYears += years;
    }
  }

  // Check for immediate retirement criteria
  return (
    (retirementAge >= mra && totalServiceYears >= 30) || // MRA + 30
    (retirementAge >= 60 && totalServiceYears >= 20) || // Age 60 + 20
    (retirementAge >= 62 && totalServiceYears >= 5)
  ); // Age 62 + 5
}

// Helper function to get Minimum Retirement Age
function getMRA(dateOfBirth: string): number {
  const birthYear = getBirthYear(dateOfBirth);

  if (birthYear <= 1947) return 55;
  if (birthYear <= 1952) return 56;
  if (birthYear <= 1957) return 57;
  if (birthYear <= 1964) return 58;
  if (birthYear <= 1969) return 59;
  return 60;
}

// Helper function to get retirement year
function getRetirementYear(retirementDate: string): number {
  if (!retirementDate) return 0;
  return parseInt(retirementDate.split("-")[0], 10);
}

// Main mapping function with comprehensive error handling
function mapScenarioInputToRetirementCalculationInput(
  inputs: ScenarioInput,
  tspProjections: any | null
): any {
  try {
    const currentAge = calculateAge(inputs.dateOfBirth);
    const retirementAge = calculateRetirementAge(
      inputs.dateOfBirth,
      inputs.plannedRetirementDate
    );
    const birthYear = getBirthYear(inputs.dateOfBirth);
    const retirementYear = getRetirementYear(inputs.plannedRetirementDate);
    const mra = getMRA(inputs.dateOfBirth);
    const fersYearsForSRS = calculateFERSYearsForSRS(inputs);
    const isImmediateUnreduced = isImmediateUnreducedAnnuity(inputs);
    const otherTaxableIncome = calculateOtherTaxableIncome(inputs);

    // Calculate TSP withdrawal amount for tax calculations
    let tspWithdrawalAmount = 0;
    if (
      tspProjections?.withdrawalSchedule?.[0]?.traditionalWithdrawn !==
      undefined
    ) {
      tspWithdrawalAmount =
        tspProjections.withdrawalSchedule[0].traditionalWithdrawn;
    } else if (
      inputs.tspWithdrawalStrategy === "PercentageOfBalanceYearly" &&
      inputs.tspWithdrawalPercentageValue &&
      inputs.tspBalanceTraditional
    ) {
      tspWithdrawalAmount =
        inputs.tspBalanceTraditional *
        (inputs.tspWithdrawalPercentageValue / 100);
    } else {
      tspWithdrawalAmount = inputs.tspWithdrawalFixedAmountValue || 0;
    }

    // Determine the number of years to calculate based on user setting or default
    const yearsToCalculate =
      inputs.calculationYears != null ? inputs.calculationYears : 30;

    return {
      CalculationSystem: inputs.calculationSystem || "",

      FERSInput: {
        High3Salary: inputs.high3Salary || 0,
        UnusedSickLeaveHours: inputs.unusedSickLeaveHoursAtRetirement || 0,
        EmployeeContributions: inputs.employeeContributions || 0,
        ServicePeriods: inputs.servicePeriods || [],
        SurvivorBenefitElection:
          inputs.survivorBenefitFERS?.spouseElection || "",
        ExpectedSSBenefitAt62: inputs.estimatedSSBenefitAt62ForSRS || 0,
        DateOfBirth: inputs.dateOfBirth || "",
        ServiceComputationDate: inputs.serviceComputationDate || "",
        PlannedRetirementDate: inputs.plannedRetirementDate || "",
        SwitchedToFERSDate: inputs.switchedToFERSDate || "",
      },

      CSRSInput: {
        High3Salary: inputs.high3Salary || 0,
        UnusedSickLeaveHours: inputs.unusedSickLeaveHoursAtRetirement || 0,
        SurvivorBenefitType: inputs.survivorBenefitCSRS?.election || "",
        IsPartTime: false, // TODO: Map from service periods if available
        PartTimeProrationFactor: 1, // TODO: Calculate from service periods
        EmployeeContributions: inputs.employeeContributions || 0,
        IsCSRSOffset: inputs.isCSRSOffset || false,
        YearsOfOffsetService: inputs.yearsOfOffsetService || 0,
        SSAt62WithOffset: inputs.ssAt62WithOffset || 0,
        SSAt62WithoutOffset: inputs.ssAt62WithoutOffset || 0,
        DateOfBirth: inputs.dateOfBirth || "",
        ServiceComputationDate: inputs.serviceComputationDate || "",
        PlannedRetirementDate: inputs.plannedRetirementDate || "",
      },

      SRSInput: {
        EstimatedSocialSecurityAt62: inputs.estimatedSSBenefitAt62ForSRS || 0,
        YearsOfFERSService: fersYearsForSRS,
        RetirementAge: retirementAge,
        MRA: mra,
        IsImmediateUnreducedAnnuity: isImmediateUnreduced,
        ProjectedEarnedIncome: 0, // TODO: Add to UI if needed
        RetirementYear: retirementYear,
      },

      TSPInput: {
        BaseSalaryForContributions: inputs.high3Salary || 0,
        CurrentAgeYears: currentAge,
        BirthYear: birthYear,
        RetirementAgeYears: retirementAge,
        CurrentTraditionalBalance: inputs.tspBalanceTraditional || 0,
        CurrentRothBalance: inputs.tspBalanceRoth || 0,
        CurrentLoanBalance: inputs.tspLoanBalance || 0,
        CurrentAllocation: convertTSPAllocation(
          inputs.tspCurrentAllocationToFunds
        ),
        EmployeeContributionAmount:
          inputs.tspAnnualContributionPreRetirement || 0,
        EmployeeContributionPercentage:
          inputs.tspContributionPercentagePreRetirement || 0,
        IsContributionPercentage:
          inputs.tspContributionPercentagePreRetirement !== null &&
          inputs.tspContributionPercentagePreRetirement !== undefined,
        ContributeUntilRetirement: inputs.tspContributeUntil === "Retirement",
        ContributionStopAge: inputs.tspContributionStopAge || 0,
        CatchUpContributionsEligible:
          inputs.tspCatchUpContributionsEligible || false,
        TraditionalContributionAllocationPct:
          inputs.tspContributionAllocationTraditionalPercent || 0,
        RothContributionAllocationPct:
          inputs.tspContributionAllocationRothPercent || 0,
        ContributionFundAllocation: convertTSPAllocation(
          inputs.tspContributionAllocationToFunds
        ),
        UserReturnAssumptions: convertPercentagesToDecimals(
          inputs.userReturnAssumptionsTSP
        ),
        ExpenseRatio: inputs.tspExpenseRatio || 0.00051, // Default TSP expense ratio
        ExpectedAnnualInflationRate:
          (inputs.userAssumedGeneralInflationRate || 2.5) / 100, // Convert percentage to decimal
        FutureAllocationStrategy:
          inputs.tspFutureAllocationStrategy || "MaintainCurrent",
        PostRetirementAllocation: convertTSPAllocation(
          inputs.tspPostRetirementAllocation
        ),
        WithdrawalStrategy: inputs.tspWithdrawalStrategy || "None",
        WithdrawalFixedAmountValue: inputs.tspWithdrawalFixedAmountValue || 0,
        WithdrawalPercentageValue: inputs.tspWithdrawalPercentageValue || 0,
        WithdrawalStartDate: inputs.tspWithdrawalStartDate || "Retirement",
        WithdrawalStartAge: inputs.tspWithdrawalStartAge || 0,
        WithdrawalOrder: inputs.tspWithdrawalOrder || "ProRata",
        YearsToProjectWithdrawals: yearsToCalculate,
      },

      TaxInput: {
        IRSSimplifiedMethodExclusion: 0, // Will be calculated by backend
        FilingStatus:
          inputs.federalTaxFilingStatus ||
          (inputs.isMarriedAtRetirement ? "MarriedFilingJointly" : "Single"),
        TaxYear: new Date().getFullYear(),
        Age: currentAge,
        NumberOfDependents: inputs.federalTaxNumberOfDependents || 0,
        IsBlind: false, // TODO: Add to UI if needed
        GrossPension: 0, // Will be calculated by backend
        TaxablePension: 0, // Will be calculated by backend
        TSPWithdrawal: tspWithdrawalAmount,
        TSPRothWithdrawal:
          tspProjections?.withdrawalSchedule?.[0]?.rothWithdrawn || 0,
        SocialSecurity: 0, // Will be calculated by backend
        OtherTaxableIncome: otherTaxableIncome,
        StateOfResidence: inputs.stateOfResidenceForTax || "",
        StateTaxableIncome: 0, // Will be calculated by backend
        Deductions: inputs.federalTaxItemizedDeductions || 0,
        TaxCredits: inputs.federalTaxCreditsAnnual || 0,
        TaxLawAssumption: inputs.federalTaxLawAssumption || "CurrentLaw",
      },

      SocialSecurityInput: {
        BirthYear: birthYear,
        CurrentAge: currentAge,
        EarningsHistory:
          inputs.socialSecurityCreditedEarnings?.map((e) => e.earnings) || [],
        EstimatedAnnualSalary: inputs.high3Salary || 0,
        YearsWorked: 0, // TODO: Calculate from service periods
        UserProvidedEstimate62: (inputs.userProvidedSSBenefitAmount1 || 0) * 12, // Convert monthly to annual
        UserProvidedEstimateFRA: (inputs.socialSecurityEstimate || 0) * 12, // Convert monthly to annual
        UserProvidedEstimate70: (inputs.userProvidedSSBenefitAmount2 || 0) * 12, // Convert monthly to annual
        ClaimAge: inputs.userSSClaimingAge || 67, // Default to FRA
      },

      COLAInput: {
        InitialAmount: 0, // Will be calculated by backend
        COLARate: (inputs.userAssumedGeneralInflationRate || 2.5) / 100, // Convert percentage to decimal
        Years: yearsToCalculate, // Use global setting
        COLAPolicy: inputs.calculationSystem || "",
        StartYear: retirementYear,
      },

      SurvivorInput: {
        PensionType: inputs.calculationSystem || "",
        InitialAnnuity: 0, // Will be calculated by backend
        SurvivorElection:
          inputs.survivorBenefitFERS?.spouseElection ||
          inputs.survivorBenefitCSRS?.election ||
          "",
        SpouseAge: inputs.spouseBirthDate
          ? calculateAge(inputs.spouseBirthDate)
          : 0,
        RetireeAgeAtDeath: 0, // TODO: Add to UI if needed
        COLARate: (inputs.userAssumedGeneralInflationRate || 2.5) / 100, // Convert percentage to decimal
        YearsToProject: yearsToCalculate, // Use global setting
        IncludeSSSurvivor: false, // TODO: Add to UI if needed
        SSSurvivorAmount: 0, // TODO: Add to UI if needed
        IncludeTSP: false, // TODO: Add to UI if needed
        TSPBalanceAtDeath: 0, // TODO: Add to UI if needed
        OtherSurvivorIncome: 0, // TODO: Add to UI if needed
      },

      HealthInput: {
        FEHBPremium: inputs.fehbBiweeklyDeduction
          ? inputs.fehbBiweeklyDeduction * 26
          : 0, // Convert biweekly to annual
        MedicarePremium: 0, // TODO: Add to UI if needed
        IncludeFEHB: inputs.fehbContinuedInRetirement || false,
        IncludeMedicare: false, // TODO: Add to UI if needed
        COLARate: (inputs.fehbProjectedPremiumIncreaseRate || 2.5) / 100, // Convert percentage to decimal
        YearsToProject: yearsToCalculate, // Use global setting
        OtherHealthPremium: 0, // TODO: Add to UI if needed
      },

      MonteCarloInput: {
        NumSimulations: inputs.monteCarloNumberOfIterations || 1000,
        Years: inputs.monteCarloSimulationEndAge
          ? Math.min(
              inputs.monteCarloSimulationEndAge - retirementAge,
              yearsToCalculate
            )
          : yearsToCalculate, // Use global setting
        InitialBalance: inputs.tspBalanceTraditional || 0,
        AnnualWithdrawal: tspWithdrawalAmount,
        ExpectedReturn: 0.06, // Default 6%
        ReturnStdDev: 0.15, // Default 15%
        InflationMean: (inputs.userAssumedGeneralInflationRate || 2.5) / 100, // Convert percentage to decimal
        InflationStdDev: 0.02, // Default 2%
        Seed: Math.floor(Math.random() * 1000000), // Random seed
      },
    };
  } catch (error) {
    console.error(
      "Error in mapScenarioInputToRetirementCalculationInput:",
      error
    );
    throw new Error(
      `Data mapping error: ${
        error instanceof Error ? error.message : "Unknown error"
      }`
    );
  }
}

// Main calculation function
export async function performCalculation(
  inputs: ScenarioInput
): Promise<CalculationServiceResult> {
  // console.log('Starting calculation with inputs:', inputs);

  // Validate inputs first
  const validationMessages = validateScenarioInputs(inputs);
  if (validationMessages.length > 0) {
    console.warn("Validation failed:", validationMessages);
    return {
      results: null,
      tspProjections: null,
      error: validationMessages.join("\n"),
      validationMessages: validationMessages,
    };
  }

  // FINAL SAFETY CHECK: Ensure servicePeriods is always present before backend call
  if (
    (!inputs.servicePeriods || inputs.servicePeriods.length === 0) &&
    inputs.serviceComputationDate &&
    inputs.plannedRetirementDate &&
    inputs.calculationSystem
  ) {
    const defaultServicePeriod = {
      id: globalThis.crypto?.randomUUID() || Math.random().toString(36),
      startDate: inputs.serviceComputationDate,
      endDate: inputs.plannedRetirementDate,
      serviceCategory: "Civilian" as const,
      civilianServiceType:
        inputs.calculationSystem === "CSRS"
          ? "NonDeductionPre10_82CSRS"
          : "RegularDeductionFERS",
      militaryServiceType: null,
      depositRedepositPaymentStatus: "NotApplicable" as const,
      systemDuringService: inputs.calculationSystem,
      isPartTime: false,
      hoursPerWeekIfPartTime: null,
      notes: "Auto-created to span SCD to retirement date.",
    };
    inputs.servicePeriods = [defaultServicePeriod];
  }

  let results: RetirementCalculationResult | null = null;
  let tspProjections: any | null = null;
  let calculationError: string | null = null;

  try {
    // Step 1: Calculate TSP projections if applicable
    if (
      inputs.tspBalanceTraditional ||
      inputs.tspBalanceRoth ||
      inputs.tspAnnualContributionPreRetirement ||
      inputs.tspContributionPercentagePreRetirement
    ) {
      try {
        // console.log('Calculating TSP projections...');
        tspProjections = await TSPService.calculateTSP(inputs);
        // console.log('TSP projections calculated:', tspProjections);
      } catch (tspError: any) {
        console.error("TSP calculation error:", tspError);
        // Continue with main calculation even if TSP fails
      }
    }

    // Step 2: Perform main retirement calculation
    // console.log('Mapping inputs to backend format...');
    const retirementInput = mapScenarioInputToRetirementCalculationInput(
      inputs,
      tspProjections
    );

    // Debug: Log service periods being sent to backend
    console.log(
      "DEBUG: ServicePeriods sent to backend:",
      inputs.servicePeriods
    );
    console.log(
      "DEBUG: FERSInput.ServicePeriods:",
      retirementInput.FERSInput.ServicePeriods
    );

    // console.log('Sending to backend calculation...');
    // console.log('Retirement input structure:', {
    //   CalculationSystem: retirementInput.CalculationSystem,
    //   FERSInput: { High3Salary: retirementInput.FERSInput.High3Salary, ServicePeriods: retirementInput.FERSInput.ServicePeriods.length },
    //   TSPInput: { CurrentTraditionalBalance: retirementInput.TSPInput.CurrentTraditionalBalance },
    //   TaxInput: { FilingStatus: retirementInput.TaxInput.FilingStatus, OtherTaxableIncome: retirementInput.TaxInput.OtherTaxableIncome }
    // });

    results = await PerformRetirementCalculation(retirementInput);
    // console.log('Backend calculation completed successfully');
  } catch (e: any) {
    console.error("Calculation Service Error:", e);
    calculationError = e.message || e.toString();
    results = null;
    tspProjections = null;
  }

  return {
    results,
    tspProjections,
    error: calculationError,
    validationMessages: [],
  };
}
