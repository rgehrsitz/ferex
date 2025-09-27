import type {
  ScenarioInput,
  TSPReturnAssumptions as FrontendTSPReturnAssumptions,
  TSPVolatility,
} from "../../types"; // Added TSPVolatility
import { models } from "../../../wailsjs/go/models";

/**
 * Maps frontend TSP data to backend TSP calculation input format.
 * @param inputs Frontend scenario inputs.
 * @returns Properly formatted TSP calculation input for the backend.
 */
export function mapTSPDataToBackend(
  inputs: ScenarioInput
): Partial<models.TSPCalculationInput> {
  const tspData: Partial<models.TSPCalculationInput> = {
    currentTraditionalBalance: inputs.tspBalanceTraditional || 0,
    currentRothBalance: inputs.tspBalanceRoth || 0,
    currentLoanBalance: inputs.tspLoanBalance || 0,
    traditionalContributionAllocationPct:
      inputs.tspContributionAllocationTraditionalPercent || 0,
    rothContributionAllocationPct:
      inputs.tspContributionAllocationRothPercent || 0,
    withdrawalStrategy: inputs.tspWithdrawalStrategy || "None",
    withdrawalFixedAmountValue: inputs.tspWithdrawalFixedAmountValue || 0,
    withdrawalPercentageValue: inputs.tspWithdrawalPercentageValue || 0,
    withdrawalStartAge: inputs.tspWithdrawalStartAge || 0,
    // Required fields by backend that need to be sourced or derived:
    // baseSalaryForContributions: number - (e.g., from inputs.high3Salary)
    // currentAgeYears: number - (derive from inputs.dateOfBirth)
    // birthYear: number - (derive from inputs.dateOfBirth)
    // retirementAgeYears: number - (derive from inputs.plannedRetirementDate and inputs.dateOfBirth)
    // catchUpContributionsEligible: boolean - (e.g., based on age)
    // expectedAnnualInflationRate: number - (e.g., from inputs.userAssumedGeneralInflationRate)
    // withdrawalStartDate: string - (e.g., from inputs.tspWithdrawalStartDate or derive from inputs.tspWithdrawalStartAge)
  };

  if (
    inputs.tspContributionPercentagePreRetirement !== null &&
    inputs.tspContributionPercentagePreRetirement !== undefined
  ) {
    tspData.isContributionPercentage = true;
    tspData.employeeContributionPercentage =
      inputs.tspContributionPercentagePreRetirement;
    tspData.employeeContributionAmount = undefined;
  } else if (
    inputs.tspAnnualContributionPreRetirement !== null &&
    inputs.tspAnnualContributionPreRetirement !== undefined
  ) {
    tspData.isContributionPercentage = false;
    tspData.employeeContributionAmount =
      inputs.tspAnnualContributionPreRetirement;
    tspData.employeeContributionPercentage = undefined;
  } else {
    tspData.isContributionPercentage = false;
    tspData.employeeContributionAmount = 0;
    tspData.employeeContributionPercentage = undefined;
  }

  if (
    inputs.tspContributeUntil === "SpecificAge" &&
    inputs.tspContributionStopAge !== null
  ) {
    tspData.contributeUntilRetirement = false;
    tspData.contributionStopAge = inputs.tspContributionStopAge;
  } else {
    tspData.contributeUntilRetirement = true;
    tspData.contributionStopAge = undefined;
  }

  if (inputs.tspCurrentAllocationToFunds) {
    tspData.currentAllocation = models.TSPFundAllocationPercentages.createFrom({
      g: inputs.tspCurrentAllocationToFunds.G || 0,
      f: inputs.tspCurrentAllocationToFunds.F || 0,
      c: inputs.tspCurrentAllocationToFunds.C || 0,
      s: inputs.tspCurrentAllocationToFunds.S || 0,
      i: inputs.tspCurrentAllocationToFunds.I || 0,
      lFundName: inputs.tspCurrentAllocationToFunds.LFundName || undefined,
    });
  }

  if (inputs.tspContributionAllocationToFunds) {
    tspData.contributionFundAllocation =
      models.TSPFundAllocationPercentages.createFrom({
        g: inputs.tspContributionAllocationToFunds.G || 0,
        f: inputs.tspContributionAllocationToFunds.F || 0,
        c: inputs.tspContributionAllocationToFunds.C || 0,
        s: inputs.tspContributionAllocationToFunds.S || 0,
        i: inputs.tspContributionAllocationToFunds.I || 0,
        lFundName:
          inputs.tspContributionAllocationToFunds.LFundName || undefined,
      });
  }

  if (inputs.tspFutureAllocationStrategy) {
    tspData.futureAllocationStrategy = inputs.tspFutureAllocationStrategy;
    if (
      inputs.tspFutureAllocationStrategy === "UsePostRetirementAllocation" &&
      inputs.tspPostRetirementAllocation
    ) {
      tspData.postRetirementAllocation =
        models.TSPFundAllocationPercentages.createFrom({
          g: inputs.tspPostRetirementAllocation.G || 0,
          f: inputs.tspPostRetirementAllocation.F || 0,
          c: inputs.tspPostRetirementAllocation.C || 0,
          s: inputs.tspPostRetirementAllocation.S || 0,
          i: inputs.tspPostRetirementAllocation.I || 0,
          lFundName: inputs.tspPostRetirementAllocation.LFundName || undefined,
        });
    }
  }

  if (inputs.userReturnAssumptionsTSP) {
    const backendAssumptions = models.TSPReturnAssumptions.createFrom({});
    backendAssumptions.g = inputs.userReturnAssumptionsTSP.G || 0;
    backendAssumptions.f = inputs.userReturnAssumptionsTSP.F || 0;
    backendAssumptions.c = inputs.userReturnAssumptionsTSP.C || 0;
    backendAssumptions.s = inputs.userReturnAssumptionsTSP.S || 0;
    backendAssumptions.i = inputs.userReturnAssumptionsTSP.I || 0;

    if (inputs.tspVolatilityAssumptions) {
      backendAssumptions.volatilityG =
        inputs.tspVolatilityAssumptions.G_stdDev || undefined;
      backendAssumptions.volatilityF =
        inputs.tspVolatilityAssumptions.F_stdDev || undefined;
      backendAssumptions.volatilityC =
        inputs.tspVolatilityAssumptions.C_stdDev || undefined;
      backendAssumptions.volatilityS =
        inputs.tspVolatilityAssumptions.S_stdDev || undefined;
      backendAssumptions.volatilityI =
        inputs.tspVolatilityAssumptions.I_stdDev || undefined;
    }
    tspData.userReturnAssumptions = backendAssumptions;
  }

  return tspData;
}

/**
 * Maps backend TSP calculation results back to frontend format.
 * @param backendResult Backend TSP calculation results (models.TSPCalculationResult).
 * @param currentInputs Current frontend scenario inputs.
 * @returns Updated frontend scenario inputs.
 */
export function mapTSPDataFromBackend(
  backendResult: models.TSPCalculationResult,
  currentInputs: ScenarioInput
): ScenarioInput {
  const updatedInputs = { ...currentInputs };

  if (backendResult && typeof backendResult === "object") {
    if (backendResult.projectedTraditionalBalanceAtRetirement !== undefined) {
      updatedInputs.tspBalanceTraditional =
        backendResult.projectedTraditionalBalanceAtRetirement;
    }
    if (backendResult.projectedRothBalanceAtRetirement !== undefined) {
      updatedInputs.tspBalanceRoth =
        backendResult.projectedRothBalanceAtRetirement;
    }
    // Note: Other fields from TSPCalculationInput are not part of TSPCalculationResult.
    // If the backend needs to return the input parameters or more detailed state,
    // the TSPCalculationResult model would need to be expanded or a different structure used.
  }
  return updatedInputs;
}
