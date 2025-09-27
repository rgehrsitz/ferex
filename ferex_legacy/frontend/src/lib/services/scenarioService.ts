import type {
  Scenario,
  ScenarioInput,
  TSPFundAllocation,
  TSPReturnAssumptions,
  TSPVolatility,
} from "../../types";
// Potentially other Wails imports will be needed here later for save/load functionality
// import { scenario as wailsScenarioModels } from '../../../wailsjs/go/models';

// Ensures all required fields for ScenarioInput are initialized for Svelte 5 binding safety
export function sanitizeScenarioInputs (inputs: ScenarioInput): ScenarioInput {
  // Arrays that must always be present
  if (!inputs.servicePeriods) inputs.servicePeriods = [];
  if (!inputs.lwopPeriods) inputs.lwopPeriods = [];
  if (!inputs.socialSecurityCreditedEarnings)
    inputs.socialSecurityCreditedEarnings = [];
  if (!inputs.socialSecurityBenefitTable)
    inputs.socialSecurityBenefitTable = [];
  if (!inputs.otherRecurringIncomeSources)
    inputs.otherRecurringIncomeSources = [];
  if (!inputs.oneTimeIncomeEvents) inputs.oneTimeIncomeEvents = [];
  if (!inputs.oneTimeExpenseEvents) inputs.oneTimeExpenseEvents = [];

  // Objects that must always be present
  if (!inputs.survivorBenefitFERS) {
    inputs.survivorBenefitFERS = {
      spouseElection: null,
      formerSpouseElection: null,
      formerSpouseConsent: "NotApplicable",
      currentSpouseConsentForFormer: "NotApplicable",
      currentSpouseWaiverForSelf: "NotApplicable",
      insurableInterestDetails: null,
    };
  }
  if (!inputs.survivorBenefitCSRS) {
    inputs.survivorBenefitCSRS = {
      election: null,
      customBaseAmountForPartial: null,
      formerSpouseElection: null,
      formerSpouseCustomBaseAmount: null,
      formerSpouseConsent: "NotApplicable",
      currentSpouseConsentForFormer: "NotApplicable",
      currentSpouseWaiverForSelf: "NotApplicable",
      insurableInterestDetails: null,
    };
  }
  if (!inputs.otherTaxableAccountAssetAllocation) {
    inputs.otherTaxableAccountAssetAllocation = {
      stocksPercent: 60,
      bondsPercent: 30,
      cashPercent: 10,
    };
  }
  if (
    !inputs.tspContributionAllocationToFunds ||
    Object.keys(inputs.tspContributionAllocationToFunds).length === 0
  ) {
    inputs.tspContributionAllocationToFunds = {
      G: 20,
      F: 20,
      C: 20,
      S: 20,
      I: 20,
      LFundName: null,
    };
  }
  if (
    !inputs.tspCurrentAllocationToFunds ||
    Object.keys(inputs.tspCurrentAllocationToFunds).length === 0
  ) {
    inputs.tspCurrentAllocationToFunds = {
      G: 20,
      F: 20,
      C: 20,
      S: 20,
      I: 20,
      LFundName: null,
    };
  }
  if (
    !inputs.tspPostRetirementAllocation ||
    Object.keys(inputs.tspPostRetirementAllocation).length === 0
  ) {
    inputs.tspPostRetirementAllocation = {
      G: 100,
      F: 0,
      C: 0,
      S: 0,
      I: 0,
      LFundName: null,
    };
  }
  if (
    !inputs.userReturnAssumptionsTSP ||
    Object.keys(inputs.userReturnAssumptionsTSP).length === 0
  ) {
    inputs.userReturnAssumptionsTSP = {
      G: 2.0,
      F: 2.5,
      C: 7.0,
      S: 7.5,
      I: 6.5,
    };
  }
  if (
    !inputs.tspVolatilityAssumptions ||
    Object.keys(inputs.tspVolatilityAssumptions).length === 0
  ) {
    inputs.tspVolatilityAssumptions = {
      G_stdDev: 0.5,
      F_stdDev: 3.0,
      C_stdDev: 15.0,
      S_stdDev: 18.0,
      I_stdDev: 16.0,
    };
  }

  // Ensure COLA/Inflation rates are present
  if (inputs.pensionAnnuityColaRate === undefined)
    inputs.pensionAnnuityColaRate = null;
  if (inputs.userAssumedSSCOLA === undefined) inputs.userAssumedSSCOLA = null; // Initialized in createNewScenario, but good for loaded data
  if (inputs.userAssumedGeneralInflationRate === undefined)
    inputs.userAssumedGeneralInflationRate = null; // Initialized in createNewScenario, but good for loaded data

  return inputs;
}

export function deleteScenario (
  currentScenarios: Scenario[],
  idToDelete: string
): { updatedScenarios: Scenario[] } {
  let updatedScenarios = currentScenarios.filter((s) => s.id !== idToDelete);
  if (updatedScenarios.length === 0) {
    updatedScenarios = [createNewScenario("Default Scenario (auto-created)")];
  }
  return { updatedScenarios };
}

export function addScenario (currentScenarios: Scenario[]): {
  updatedScenarios: Scenario[];
  newScenarioId: string;
} {
  const newScenario = createNewScenario(
    `Scenario ${currentScenarios.length + 1}`
  );
  const updatedScenarios = [...currentScenarios, newScenario];
  return { updatedScenarios, newScenarioId: newScenario.id };
}

export function createNewScenario (name: string): Scenario {
  return {
    id: crypto.randomUUID(),
    name: name,
    inputs: sanitizeScenarioInputs({
      // Core Details
      calculationSystem: "", // User must select FERS or CSRS
      dateOfBirth: "",
      plannedRetirementDate: "",
      serviceComputationDate: "",
      high3Salary: null,
      unusedSickLeaveHoursAtRetirement: null,
      employeeContributions: null,

      // Service & LWOP
      servicePeriods: [],
      lwopPeriods: [],

      // Retirement Eligibility / Special Conditions
      fersCoverageType: null, // FERS specific coverage variant (FERS, FERS_RAE, FERS_FRAE)
      isVeraDsRetirement: false, // VERA/Deferred Service Retirement
      mraPlus10PostponeAnnuityStartDate: null,
      isDeferredRetirement: false,
      deferredRetirementAnnuityStartDate: null,
      isDisabilityRetirement: false,

      // FERS Specific
      estimatedSSBenefitAt62ForSRS: 0, // Default to 0, cannot be negative
      didSwitchFromCSRS: false,
      switchedToFERSDate: "",
      // survivorBenefitFERS is initialized below

      // CSRS Specific
      isCSRSOffset: false,
      yearsOfOffsetService: null,
      ssAt62WithOffset: null,
      ssAt62WithoutOffset: null,
      // survivorBenefitCSRS is initialized below

      // Survivor Benefits / Spouse Details
      isMarriedAtRetirement: false,
      spouseBirthDate: "", // Correctly initialized as empty string
      ageOfInsurableInterestBeneficiary: null,
      hasFormerSpouseEntitlement: false,
      formerSpouseSurvivorBenefitDetails: null, // This comes from wails models.FormerSpouseSurvivorDetails

      // FERS Survivor Benefit Details (initialized even if not FERS, for UI consistency)
      survivorBenefitFERS: {
        spouseElection: null,
        formerSpouseElection: null,
        formerSpouseConsent: "NotApplicable",
        currentSpouseConsentForFormer: "NotApplicable",
        currentSpouseWaiverForSelf: "NotApplicable",
        insurableInterestDetails: null, // models.InsurableInterestDetails
      },

      // CSRS Survivor Benefit Details (initialized even if not CSRS, for UI consistency)
      survivorBenefitCSRS: {
        election: null,
        customBaseAmountForPartial: null,
        formerSpouseElection: null,
        formerSpouseCustomBaseAmount: null,
        formerSpouseConsent: "NotApplicable",
        currentSpouseConsentForFormer: "NotApplicable",
        currentSpouseWaiverForSelf: "NotApplicable",
        insurableInterestDetails: null, // models.InsurableInterestDetails
      },

      // Social Security
      userProvidedSSBenefitAmount1: null,
      userProvidedSSBenefitClaimingAge1: null,
      userProvidedSSBenefitAmount2: null,
      userProvidedSSBenefitClaimingAge2: null,
      userSSClaimingAge: null, // User's own SS claiming age
      userAssumedSSCOLA: 2.0, // Default assumption
      pensionAnnuityColaRate: 2.0, // Default assumption for FERS/CSRS pension COLA
      calculateHistoricalWEPGPO: false,
      socialSecurityEstimate: null,
      socialSecurityFRA: null,
      socialSecurityBenefitStartDate: null,
      socialSecurityCreditedEarnings: [],
      // Spouse SS moved further down for grouping

      // TSP
      tspBalanceTraditional: null,
      tspBalanceRoth: null,
      tspLoanBalance: null,
      tspAnnualContributionPreRetirement: null,
      tspContributionPercentagePreRetirement: null,
      tspContributeUntil: "Retirement",
      tspContributionStopAge: null,
      tspCatchUpContributionsEligible: false,
      tspContributionAllocationTraditionalPercent: 100, // Default to 100% Traditional
      tspContributionAllocationRothPercent: 0,
      tspContributionAllocationToFunds: {
        G: 20,
        F: 20,
        C: 20,
        S: 20,
        I: 20,
        LFundName: null,
      },
      tspCurrentAllocationToFunds: {
        G: 20,
        F: 20,
        C: 20,
        S: 20,
        I: 20,
        LFundName: null,
      },
      tspPostRetirementAllocation: {
        G: 100,
        F: 0,
        C: 0,
        S: 0,
        I: 0,
        LFundName: null,
      }, // Default to 100% G in retirement
      tspFutureAllocationStrategy: "MaintainCurrent",
      userReturnAssumptionsTSP: { G: 2.0, F: 2.5, C: 7.0, S: 7.5, I: 6.5 },
      tspExpenseRatio: 0.05, // Example default, ensure this is reasonable
      tspExpectedAnnualGrowthRatePreRetirement: null, // Calculated if not set
      tspWithdrawalStrategy: "None",
      tspWithdrawalFixedAmountValue: null,
      tspWithdrawalPercentageValue: null,
      tspWithdrawalStartDate: "Retirement",
      tspWithdrawalStartAge: null,
      tspWithdrawalOrder: "ProRata",
      applyRMDsToTSP: true,
      tspExpectedAnnualGrowthRatePostRetirement: null, // Calculated if not set
      tspVolatilityAssumptions: {
        G_stdDev: 0.5,
        F_stdDev: 3.0,
        C_stdDev: 15.0,
        S_stdDev: 18.0,
        I_stdDev: 16.0,
      },

      // FEGLI
      fegliBasicCoverage: false,
      fegliBasicHeldIntoRetirement: false,
      fegliBasicPostRetirementReductionChoice: "75Reduce",
      fegliOptionAAmount: 0,
      fegliOptionBCoverageMultiples: 0,
      fegliOptionCCoverageMultiples: 0,

      // FEHB
      fehbBiweeklyDeduction: null,
      fehbProjectedPremiumIncreaseRate: 5.0,
      fehbContinuedInRetirement: true,

      // FLTCIP
      fltcipHasCoverage: false,
      fltcipDailyBenefitAmount: null,
      fltcipBenefitPeriod: undefined, // As per previous setup
      fltcipInflationProtectionType: undefined, // As per previous setup
      fltcipCurrentAnnualPremium: null,

      // Financial & Tax Planning - Using backend field names
      federalTaxFilingStatus: "Single",
      federalTaxNumberOfDependents: 0,
      federalTaxItemizedDeductions: null,
      federalOtherTaxableIncomeAnnual: null,
      federalTaxCreditsAnnual: null,
      federalTaxLawAssumption: "CurrentLaw",
      stateOfResidenceForTax: "",
      retirementEffectiveStateTaxRate: null,
      userAssumedGeneralInflationRate: 2.5,
      userAssumedPensionCOLA_FERS: 2.0, // Specific to FERS
      userAssumedPensionCOLA_CSRS: 2.0, // Specific to CSRS
      otherRecurringIncomeSources: [],
      oneTimeIncomeEvents: [],
      otherTaxableAccountBalance: null,
      otherTaxableAccountAssetAllocation: {
        stocksPercent: 60,
        bondsPercent: 30,
        cashPercent: 10,
      },
      otherTaxableAccountExpectedReturn: 5.0,
      homeValue: null,
      mortgageBalance: null,
      otherSignificantDebtsTotal: null,
      desiredAnnualRetirementSpending: null,
      applyInflationToDesiredSpending: true,
      oneTimeExpenseEvents: [],

      // Spouse Financial Info (if applicable, often linked to isMarriedAtRetirement)
      spouseDateOfBirth: null, // This is distinct from spouseBirthDate for survivor benefits; used for joint financial planning
      spousePlannedRetirementDate: null,
      spouseProvidedSSBenefitAmount1: null,
      spouseProvidedSSBenefitClaimingAge1: null,
      spouseSSClaimingAge: null,
      spouseHasPension: false,
      spousePensionAnnualAmount: null,
      spousePensionCOLA: 2.0,
      spouseOtherRetirementAccountBalance: null,
      spouseOtherRetirementAccountAnnualWithdrawal: null,

      // Monte Carlo
      monteCarloNumberOfIterations: 5000,
      monteCarloSimulationEndAge: 95,
      monteCarloSuccessDefinition: "AssetsRemainPositive",
      monteCarloMinEssentialSpendingLevel: null,

      // User Notes
      notes: "",
      createdAt: new Date().toISOString(),
      lastModified: new Date().toISOString(),
    }),
    results: null,
    error: null,
    calculating: false,
  };
}

/**
 * Duplicates a scenario.
 * @param currentScenarios The current array of scenarios.
 * @param idToDuplicate The ID of the scenario to duplicate.
 * @returns Object containing the updated scenarios array and the ID of the newly duplicated scenario.
 *          Returns null if the scenario to duplicate is not found.
 */
export function duplicateScenario (
  currentScenarios: Scenario[],
  idToDuplicate: string | undefined
): { updatedScenarios: Scenario[]; newScenarioId: string } | null {
  const scenarioToDuplicate = currentScenarios.find(
    (s) => s.id === idToDuplicate
  );
  if (!scenarioToDuplicate) {
    return null; // Scenario to duplicate not found
  }

  const newId = crypto.randomUUID();
  const newName = `${scenarioToDuplicate.name} (Copy)`;

  // Deep copy of inputs. Results, error, calculating are reset.
  const copiedInputs = JSON.parse(
    JSON.stringify(scenarioToDuplicate.inputs)
  ) as ScenarioInput;

  const duplicatedScenario: Scenario = {
    id: newId,
    name: newName,
    inputs: sanitizeScenarioInputs({ ...copiedInputs }),
    results: null, // Reset results
    error: null, // Reset error
    calculating: false, // Reset calculating state
  };

  const updatedScenarios = [...currentScenarios, duplicatedScenario];
  return { updatedScenarios, newScenarioId: duplicatedScenario.id };
}

// ... (rest of the code remains the same)

// --- Wails File Operations ---
import {
  OpenFileAndLoadScenario,
  SaveScenarioWithDialog,
} from "../../../wailsjs/go/main/App";
import { scenario, models } from "../../../wailsjs/go/models";

// -- Helper conversions for Wails models --
function toWailsTSPFundAllocation (
  alloc: TSPFundAllocation | null | undefined
): models.TSPFundAllocation | undefined {
  // Use default values if nothing is provided
  if (!alloc) {
    return models.TSPFundAllocation.createFrom({
      g: 20,
      f: 20,
      c: 20,
      s: 20,
      i: 20,
      lFundName: null,
    });
  }

  try {
    // Ensure we're not sending undefined values by using default percentages
    return models.TSPFundAllocation.createFrom({
      g: alloc.G ?? 20,
      f: alloc.F ?? 20,
      c: alloc.C ?? 20,
      s: alloc.S ?? 20,
      i: alloc.I ?? 20,
      lFundName: alloc.LFundName ?? null,
    });
  } catch (error) {
    console.error("Error converting TSPFundAllocation to Wails format:", error);
    // console.debug("Original allocation:", alloc);
    // Return a default object rather than failing
    return models.TSPFundAllocation.createFrom({
      g: 20,
      f: 20,
      c: 20,
      s: 20,
      i: 20,
      lFundName: null,
    });
  }
}

function toWailsTSPReturnAssumptions (
  ra: TSPReturnAssumptions | null | undefined
): models.TSPIndividualReturnAssumptions | undefined {
  // Use default values if nothing is provided
  if (!ra) {
    return models.TSPIndividualReturnAssumptions.createFrom({
      g: 2.0,
      f: 2.5,
      c: 7.0,
      s: 7.5,
      i: 6.5,
    });
  }

  try {
    return models.TSPIndividualReturnAssumptions.createFrom({
      g: ra.G ?? 2.0,
      f: ra.F ?? 2.5,
      c: ra.C ?? 7.0,
      s: ra.S ?? 7.5,
      i: ra.I ?? 6.5,
    });
  } catch (error) {
    console.error(
      "Error converting TSPReturnAssumptions to Wails format:",
      error
    );
    // Return default values instead of undefined
    return models.TSPIndividualReturnAssumptions.createFrom({
      g: 2.0,
      f: 2.5,
      c: 7.0,
      s: 7.5,
      i: 6.5,
    });
  }
}

function toWailsTSPVolatility (
  vol: TSPVolatility | null | undefined
): models.TSPVolatilityRates | undefined {
  // Use default values if nothing is provided
  if (!vol) {
    return models.TSPVolatilityRates.createFrom({
      gStdDev: 0.5,
      fStdDev: 3.0,
      cStdDev: 15.0,
      sStdDev: 18.0,
      iStdDev: 16.0,
    });
  }

  try {
    return models.TSPVolatilityRates.createFrom({
      gStdDev: vol.G_stdDev ?? 0.5,
      fStdDev: vol.F_stdDev ?? 3.0,
      cStdDev: vol.C_stdDev ?? 15.0,
      sStdDev: vol.S_stdDev ?? 18.0,
      iStdDev: vol.I_stdDev ?? 16.0,
    });
  } catch (error) {
    console.error("Error converting TSPVolatility to Wails format:", error);
    // Return default values instead of undefined
    return models.TSPVolatilityRates.createFrom({
      gStdDev: 0.5,
      fStdDev: 3.0,
      cStdDev: 15.0,
      sStdDev: 18.0,
      iStdDev: 16.0,
    });
  }
}

function fromWailsTSPFundAllocation (
  alloc: models.TSPFundAllocation | null | undefined
): TSPFundAllocation | undefined {
  if (!alloc) return undefined;
  return {
    G: alloc.g ?? null,
    F: alloc.f ?? null,
    C: alloc.c ?? null,
    S: alloc.s ?? null,
    I: alloc.i ?? null,
    LFundName: alloc.lFundName ?? null,
  };
}

function fromWailsTSPReturnAssumptions (
  ra: models.TSPIndividualReturnAssumptions | null | undefined
): TSPReturnAssumptions | undefined {
  if (!ra) return undefined;
  return {
    G: ra.g ?? null,
    F: ra.f ?? null,
    C: ra.c ?? null,
    S: ra.s ?? null,
    I: ra.i ?? null,
  };
}

function fromWailsTSPVolatility (
  vol: models.TSPVolatilityRates | null | undefined
): TSPVolatility | undefined {
  if (!vol) return undefined;
  return {
    G_stdDev: vol.gStdDev ?? null,
    F_stdDev: vol.fStdDev ?? null,
    C_stdDev: vol.cStdDev ?? null,
    S_stdDev: vol.sStdDev ?? null,
    I_stdDev: vol.iStdDev ?? null,
  };
}

/**
 * Loads scenarios from a file using a Wails dialog.
 * @returns A promise that resolves to an object containing the loaded scenarios,
 *          the active scenario ID, and the user profile, or null if the operation was cancelled or failed.
 */
export async function loadScenariosFromFile (): Promise<{
  loadedScenarios: Scenario[];
  activeScenarioId: string | undefined;
  userProfile: scenario.UserProfile | null;
} | null> {
  try {
    const loadedFile: scenario.FerexFile | null =
      await OpenFileAndLoadScenario();
    if (!loadedFile) {
      // console.log("File open cancelled by user.");
      return null;
    }

    // console.log("File loaded:", loadedFile);
    let loadedScenarios: Scenario[] = [];
    let newActiveScenarioId: string | undefined;

    if (loadedFile.variants && loadedFile.variants.length > 0) {
      loadedScenarios = loadedFile.variants.map(
        (variant: scenario.ScenarioVariant, index: number) => {
          // The createNewScenario function sets up the basic structure including results, error, calculating flags.
          // We then override with loaded data.
          const baseScenario = createNewScenario(
            variant.variantName || `Loaded Scenario ${index + 1}`
          );

          // Directly map fields from backend variant to frontend ScenarioInput.
          // The backend 'variant' object should largely match the structure of frontend's ScenarioInput.
          // Ensure date fields are correctly handled (e.g., converted to ISO strings if necessary).
          // Construct an object representing the data *as it is in the file*, but with necessary transformations.
          // Fields not present in 'variant' will not be in 'fileData'.
          // The 'Omit' helps with type safety if variant's type isn't perfectly ScenarioInput for these specific fields before transformation.
          const fileData: Partial<ScenarioInput> = {
            ...(variant as unknown as Omit<
              ScenarioInput,
              | "dateOfBirth"
              | "plannedRetirementDate"
              | "serviceComputationDate"
              | "switchedToFERSDate"
              | "spouseBirthDate"
              | "spouseDateOfBirth"
              | "notes"
              | "createdAt"
              | "lastModified"
            >),
          };

          // Convert allocation objects from backend (lowercase) to frontend format
          fileData.tspContributionAllocationToFunds =
            fromWailsTSPFundAllocation(
              variant.tspContributionAllocationToFunds
            );
          fileData.tspCurrentAllocationToFunds = fromWailsTSPFundAllocation(
            variant.tspCurrentAllocationToFunds
          );
          fileData.tspPostRetirementAllocation = fromWailsTSPFundAllocation(
            variant.tspPostRetirementAllocation
          );
          fileData.userReturnAssumptionsTSP = fromWailsTSPReturnAssumptions(
            variant.userReturnAssumptionsTSP
          );
          fileData.tspVolatilityAssumptions = fromWailsTSPVolatility(
            variant.tspVolatilityAssumptions
          );

          if (variant.dateOfBirth)
            fileData.dateOfBirth = new Date(variant.dateOfBirth)
              .toISOString()
              .split("T")[0];
          if (variant.plannedRetirementDate)
            fileData.plannedRetirementDate = new Date(
              variant.plannedRetirementDate
            )
              .toISOString()
              .split("T")[0];
          if (variant.serviceComputationDate)
            fileData.serviceComputationDate = new Date(
              variant.serviceComputationDate
            )
              .toISOString()
              .split("T")[0];
          if (variant.switchedToFERSDate)
            fileData.switchedToFERSDate = new Date(variant.switchedToFERSDate)
              .toISOString()
              .split("T")[0];
          if (variant.spouseBirthDate)
            fileData.spouseBirthDate = new Date(variant.spouseBirthDate)
              .toISOString()
              .split("T")[0]; // For survivor benefits

          // For joint financial planning spouse DOB, assuming 'variant' might have a distinct field.
          // Use 'as any' if the property isn't strictly typed on wailsScenarioModels.ScenarioVariant.
          if ((variant as any).spouseDateOfBirth) {
            fileData.spouseDateOfBirth = new Date(
              (variant as any).spouseDateOfBirth
            )
              .toISOString()
              .split("T")[0];
          }

          // Handle notes: use from file if it's a string, otherwise default will be used from baseScenario.inputs
          if (typeof variant.notes === "string") {
            fileData.notes = variant.notes;
          } else if (variant.notes === null) {
            fileData.notes = undefined; // Assign undefined instead of null to match string | undefined
          }

          if (variant.createdAt)
            fileData.createdAt = new Date(variant.createdAt).toISOString();
          // Always set lastModified to now when loading, or use from file if present.
          fileData.lastModified = variant.lastModified
            ? new Date(variant.lastModified).toISOString()
            : new Date().toISOString();

          // Merge file data with base defaults. Fields in fileData will override baseScenario.inputs.
          // Then sanitize the fully merged inputs.
          const finalInputs = sanitizeScenarioInputs({
            ...baseScenario.inputs, // Start with all defaults
            ...fileData, // Override with (transformed) data from the file
          });

          console.log("TSP Post-Retirement Allocation at calculation time:", finalInputs.tspPostRetirementAllocation);

          return {
            ...baseScenario, // Includes default .results, .error, .calculating flags
            id: variant.variantId || baseScenario.id, // Prioritize loaded ID
            name: variant.variantName || baseScenario.name,
            inputs: finalInputs,
          };
        }
      );
      newActiveScenarioId =
        loadedFile.lastViewedVariantID &&
          loadedScenarios.find((s) => s.id === loadedFile.lastViewedVariantID)
          ? loadedFile.lastViewedVariantID
          : loadedScenarios[0]?.id;
    } else {
      // If file has no variants, create a default scenario
      const defaultScenario = createNewScenario(
        "Default Scenario (from empty file)"
      );
      loadedScenarios = [defaultScenario];
      newActiveScenarioId = defaultScenario.id;
    }

    return {
      loadedScenarios,
      activeScenarioId: newActiveScenarioId,
      userProfile: loadedFile.userProfile || null,
    };
  } catch (error) {
    console.error("Error opening file in service:", error);
    throw error; // Re-throw to be caught by the caller in App.svelte
  }
}

/**
 * Saves the current scenarios to a file using a Wails dialog.
 * @param scenariosToSave The array of scenarios to save.
 * @param currentActiveScenarioId The ID of the currently active scenario.
 * @param currentUserProfile The current user profile.
 * @returns A promise that resolves when the save operation is attempted.
 */
export async function saveScenariosToFile (
  scenariosToSave: Scenario[],
  currentActiveScenarioId: string | undefined,
  currentUserProfile: scenario.UserProfile | null
): Promise<void> {
  try {
    // console.log(
    //   "Starting save operation with",
    //   scenariosToSave.length,
    //   "scenario(s)"
    // );

    const variantsToSave: scenario.ScenarioVariant[] = scenariosToSave.map(
      (s) => {
        // Sanitize inputs before saving to ensure proper default values
        const sanitizedInputs = sanitizeScenarioInputs(s.inputs);

        // Debug TSP allocation values before conversion
        // console.debug("Pre-conversion TSP allocations for", s.name, ":", {
        //   contribution: sanitizedInputs.tspContributionAllocationToFunds,
        //   current: sanitizedInputs.tspCurrentAllocationToFunds,
        //   postRetirement: sanitizedInputs.tspPostRetirementAllocation,
        // });

        // Extract TSP fields that need special conversion
        const {
          tspContributionAllocationToFunds,
          tspCurrentAllocationToFunds,
          tspPostRetirementAllocation,
          userReturnAssumptionsTSP,
          tspVolatilityAssumptions,
          ...restOfInputs
        } = sanitizedInputs; // Convert TSP fields with special handling for empty/null values
        const convertedContributionAllocation = toWailsTSPFundAllocation(
          tspContributionAllocationToFunds
        );
        const convertedCurrentAllocation = toWailsTSPFundAllocation(
          tspCurrentAllocationToFunds
        );
        const convertedPostRetirementAllocation = toWailsTSPFundAllocation(
          tspPostRetirementAllocation
        );
        const convertedReturnAssumptions = toWailsTSPReturnAssumptions(
          userReturnAssumptionsTSP
        );
        const convertedVolatilityAssumptions = toWailsTSPVolatility(
          tspVolatilityAssumptions
        );

        // Verify TSP data has been properly converted
        // console.debug("Post-conversion TSP allocations:", {
        //   contribution: convertedContributionAllocation,
        //   current: convertedCurrentAllocation,
        //   postRetirement: convertedPostRetirementAllocation,
        // }); // Convert restOfInputs to make it compatible with ScenarioVariant types
        // We need to explicitly handle null values to match the backend expectations
        const processedInputs: Record<string, unknown> = {};

        // Process each field to convert null to undefined where needed
        // Cast restOfInputs to a properly indexed type to avoid TypeScript errors
        const typedRestOfInputs = restOfInputs as Record<string, unknown>;

        Object.keys(typedRestOfInputs).forEach((key) => {
          if (typedRestOfInputs[key] === null) {
            // Skip null values or set to undefined
            processedInputs[key] = undefined;
          } else {
            processedInputs[key] = typedRestOfInputs[key];
          }
        });

        // Map frontend Scenario to backend ScenarioVariant with proper Wails model conversions
        const variantData: Partial<scenario.ScenarioVariant> = {
          ...processedInputs, // Include everything except TSP fields that need special conversion
          variantId: s.id,
          variantName: s.name,
          createdAt: sanitizedInputs.createdAt || new Date().toISOString(),
          lastModified: new Date().toISOString(),
          notes: sanitizedInputs.notes || "",
          // Add the converted TSP fields
          tspContributionAllocationToFunds: convertedContributionAllocation,
          tspCurrentAllocationToFunds: convertedCurrentAllocation,
          tspPostRetirementAllocation: convertedPostRetirementAllocation,
          userReturnAssumptionsTSP: convertedReturnAssumptions,
          tspVolatilityAssumptions: convertedVolatilityAssumptions,
        };
        // We need to use a type assertion here because TypeScript can't fully verify
        // the compatibility between our frontend types and the generated Wails types
        return variantData as unknown as scenario.ScenarioVariant;
      }
    ); // Create file data structure, ensure variants are properly set up
    const userProfileData =
      currentUserProfile ||
      scenario.UserProfile.createFrom({
        employeeName: "",
        birthDate: new Date().toISOString().split("T")[0],
        mraYears: 0,
        mraMonths: 0,
      });

    // Ensure all variants have valid IDs and basic structure
    const validatedVariantsToSave = variantsToSave.map((variant) => {
      if (!variant.variantId) {
        variant.variantId = crypto.randomUUID(); // Ensure ID is present
      }
      return variant;
    });

    const fileData = scenario.FerexFile.createFrom({
      fileFormatVersion: "1.0.1", // Bump version if format changes
      lastOpenedAppVersion: "0.2.0", // TODO: Get this from actual app version
      userProfile: userProfileData,
      variants: validatedVariantsToSave,
      lastViewedVariantID:
        currentActiveScenarioId ||
        (validatedVariantsToSave.length > 0
          ? validatedVariantsToSave[0].variantId
          : undefined),
    }); // The Wails Go types might need explicit creation if direct assignment is problematic.
    // For FerexFile, if it has methods or specific Go-side initialization, ensure it's handled.
    // However, for plain data structs, direct object construction like above is often fine.

    // Validate that variants have correct structure
    if (!fileData.variants || !Array.isArray(fileData.variants)) {
      console.error(
        "Invalid fileData structure: variants missing or not an array"
      );
      throw new Error("Invalid scenario data structure: variants missing");
    }

    try {
      await SaveScenarioWithDialog(fileData); // Pass the structured object
      // console.log("File save dialog initiated via service.");
    } catch (saveError) {
      console.error("Error in SaveScenarioWithDialog:", saveError);
      // console.debug("fileData structure:", JSON.stringify(fileData, null, 2));
      throw new Error(`Error saving scenario: ${saveError}`);
    }
  } catch (error) {
    console.error("Error saving file in service:", error);
    throw error; // Re-throw to be caught by the caller in App.svelte
  }
}
