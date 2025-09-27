<script lang="ts">
  import { onMount } from "svelte";
  import { scenario, models } from "../wailsjs/go/models";
  import {
    PerformFERSCalculation,
    PerformCSRSCalculation,
    OpenFileAndLoadScenario,
    SaveScenarioWithDialog,
  } from "../wailsjs/go/main/App";
  import { TSPService } from "./lib/services/tspService";
  import {
    performCalculation,
    type CalculationServiceResult,
  } from "./lib/services/calculationService";
  import { validateScenarioInputs } from "./lib/scenarioValidation";
  import type { Scenario, ScenarioInput } from "./types";
  import {
    createNewScenario,
    addScenario as addScenarioService,
    deleteScenario as deleteScenarioService,
    duplicateScenario as duplicateScenarioService,
    loadScenariosFromFile,
    saveScenariosToFile,
  } from "./lib/services/scenarioService";
  import AnnuityServiceDetailsTab from "./lib/components/AnnuityServiceDetailsTab.svelte";
  import ScenarioManager from "./lib/components/ScenarioManager.svelte";
  import TabNavigation from "./lib/components/TabNavigation.svelte";
  import SocialSecurityTSPTab from "./lib/components/SocialSecurityTSPTab.svelte";
  import InsuranceTab from "./lib/components/InsuranceTab.svelte";
  import FinancialTaxPlanningTab from "./lib/components/FinancialTaxPlanningTab.svelte";
  import ResultsSummary from "./lib/components/ResultsSummary.svelte";
  import Sidebar from "./lib/components/layout/Sidebar.svelte";
  import TopBar from "./lib/components/layout/TopBar.svelte";
  import QuickPanel from "./lib/components/layout/QuickPanel.svelte";
  import HomePage from "./lib/components/HomePage.svelte";
  import ScenarioEditorPage from "./lib/components/ScenarioEditorPage.svelte";
  import AnalysisPage from "./lib/components/AnalysisPage.svelte";
  import IncomeProjectionsPage from "./lib/components/IncomeProjectionsPage.svelte";
  import ScenarioComparisonPage from "./lib/components/ScenarioComparisonPage.svelte";
  import RiskAnalysisPage from "./lib/components/RiskAnalysisPage.svelte";
  import HouseholdViewPage from "./lib/components/HouseholdViewPage.svelte";
  import ReportsExportsPage from "./lib/components/ReportsExportsPage.svelte";
  import SettingsPage from "./lib/components/SettingsPage.svelte";

  // --- Global App State (Wails specific setup or other app-wide logic) ---

  function renameScenario(id: string, newName: string) {
    scenarios = scenarios.map((s) =>
      s.id === id ? { ...s, name: newName } : s
    );
  }
  onMount(() => {
    // Placeholder for any Wails event listeners or other onMount logic
  });

  // --- Scenario Data Structures & Creation ---

  // --- Main Scenario State ---
  const _initialDefaultScenario = createNewScenario("Default Scenario");
  let scenarios = $state<Scenario[]>([_initialDefaultScenario]);
  let activeScenarioId = $state<string | undefined>(_initialDefaultScenario.id);

  let activeScenario = $derived(
    scenarios.find((s) => s.id === activeScenarioId)
  );

  // Tab Management
  const TABS = [
    "Annuity & Service Details",
    "Social Security & TSP",
    "Insurance",
    "Financial & Tax Planning",
  ] as const;
  type TabName = (typeof TABS)[number];
  let activeTab = $state<TabName>(TABS[0]);

  // --- User Profile State (from loaded file) ---
  // 'scenario.UserProfile' from Wails generated types
  let userProfile = $state<scenario.UserProfile | null>(null);

  // --- File Open/Save Handlers ---

  async function handleOpenFile() {
    try {
      const result = await loadScenariosFromFile();
      if (result && result.loadedScenarios) {
        scenarios = result.loadedScenarios;
        activeScenarioId = result.activeScenarioId;
        userProfile = result.userProfile;
      }
    } catch (err) {
      console.error("Error loading scenarios:", err);
      // Optionally show user notification here
    }
  }

  async function handleSaveFile() {
    try {
      await saveScenariosToFile(scenarios, activeScenarioId, userProfile);
    } catch (err) {
      console.error("Error saving scenarios:", err);
      // Optionally show user notification here
    }
  }

  let validationError = $state<string | null>(null);

  // --- Helper Reactive Derived State (overall app logic) ---
  let canCalculate = $derived(
    activeScenario &&
      activeScenario.inputs.calculationSystem &&
      activeScenario.inputs.high3Salary &&
      activeScenario.inputs.dateOfBirth &&
      activeScenario.inputs.serviceComputationDate &&
      activeScenario.inputs.plannedRetirementDate
  );

  // Debounced calculation trigger using Svelte 5 idioms
  let calculationTimeout: number | null = null;

  function triggerCalculation() {
    if (!canCalculate) return;

    // Clear existing timeout
    if (calculationTimeout) {
      clearTimeout(calculationTimeout);
    }

    // Debounce for 500ms to prevent excessive API calls
    calculationTimeout = setTimeout(() => {
      handleCalculation(activeScenario);
    }, 500);
  }

  // Auto-calculate when key scenario inputs change (restored original logic)
  $effect(() => {
    if (canCalculate) {
      // Track changes to core calculation inputs
      activeScenario?.inputs.calculationSystem;
      activeScenario?.inputs.high3Salary;
      activeScenario?.inputs.dateOfBirth;
      activeScenario?.inputs.serviceComputationDate;
      activeScenario?.inputs.plannedRetirementDate;
      activeScenario?.inputs.unusedSickLeaveHoursAtRetirement;
      activeScenario?.inputs.employeeContributions;
      activeScenario?.inputs.switchedToFERSDate;
      // Track key pension-specific inputs
      activeScenario?.inputs.survivorBenefitFERS?.spouseElection;
      activeScenario?.inputs.survivorBenefitCSRS?.election;
      activeScenario?.inputs.estimatedSSBenefitAt62ForSRS;
      // Track service periods changes
      activeScenario?.inputs.servicePeriods?.length;
      // Track Social Security inputs
      activeScenario?.inputs.userProvidedSSBenefitAmount1;
      activeScenario?.inputs.userProvidedSSBenefitClaimingAge1;
      activeScenario?.inputs.userProvidedSSBenefitAmount2;
      activeScenario?.inputs.userProvidedSSBenefitClaimingAge2;
      activeScenario?.inputs.userSSClaimingAge;
      activeScenario?.inputs.socialSecurityEstimate;
      activeScenario?.inputs.userAssumedSSCOLA;
      // Track TSP inputs
      activeScenario?.inputs.tspBalanceTraditional;
      activeScenario?.inputs.tspBalanceRoth;
      activeScenario?.inputs.tspAnnualContributionPreRetirement;
      activeScenario?.inputs.tspContributionPercentagePreRetirement;
      activeScenario?.inputs.tspContributeUntil;
      activeScenario?.inputs.tspContributionStopAge;
      activeScenario?.inputs.tspWithdrawalStrategy;
      activeScenario?.inputs.tspWithdrawalFixedAmountValue;
      activeScenario?.inputs.tspWithdrawalStartAge;
      activeScenario?.inputs.tspExpectedAnnualGrowthRatePreRetirement;
      activeScenario?.inputs.tspWithdrawalPercentageValue;
      // Track TSP allocation and return assumption objects
      activeScenario?.inputs.tspCurrentAllocationToFunds;
      activeScenario?.inputs.tspCurrentAllocationToFunds?.G;
      activeScenario?.inputs.tspCurrentAllocationToFunds?.F;
      activeScenario?.inputs.tspCurrentAllocationToFunds?.C;
      activeScenario?.inputs.tspCurrentAllocationToFunds?.S;
      activeScenario?.inputs.tspCurrentAllocationToFunds?.I;
      activeScenario?.inputs.tspCurrentAllocationToFunds?.LFundName;
      activeScenario?.inputs.tspPostRetirementAllocation;
      activeScenario?.inputs.tspPostRetirementAllocation?.G;
      activeScenario?.inputs.tspPostRetirementAllocation?.F;
      activeScenario?.inputs.tspPostRetirementAllocation?.C;
      activeScenario?.inputs.tspPostRetirementAllocation?.S;
      activeScenario?.inputs.tspPostRetirementAllocation?.I;
      activeScenario?.inputs.tspPostRetirementAllocation?.LFundName;
      activeScenario?.inputs.userReturnAssumptionsTSP;
      activeScenario?.inputs.userReturnAssumptionsTSP?.G;
      activeScenario?.inputs.userReturnAssumptionsTSP?.F;
      activeScenario?.inputs.userReturnAssumptionsTSP?.C;
      activeScenario?.inputs.userReturnAssumptionsTSP?.S;
      activeScenario?.inputs.userReturnAssumptionsTSP?.I;
      activeScenario?.inputs.tspContributionAllocationToFunds;
      activeScenario?.inputs.tspFutureAllocationStrategy;
      // Track tax-related inputs
      activeScenario?.inputs.federalTaxFilingStatus;
      activeScenario?.inputs.federalTaxNumberOfDependents;
      activeScenario?.inputs.federalTaxItemizedDeductions;
      activeScenario?.inputs.federalOtherTaxableIncomeAnnual;
      activeScenario?.inputs.federalTaxCreditsAnnual;
      activeScenario?.inputs.federalTaxLawAssumption;
      activeScenario?.inputs.stateOfResidenceForTax;
      activeScenario?.inputs.retirementEffectiveStateTaxRate;
      // Track insurance inputs (FEGLI, FEHB, FLTCIP)
      activeScenario?.inputs.fegliBasicCoverage;
      activeScenario?.inputs.fegliOptionAAmount;
      activeScenario?.inputs.fegliOptionBCoverageMultiples;
      activeScenario?.inputs.fegliOptionCCoverageMultiples;
      activeScenario?.inputs.fegliPost65Reduction;
      activeScenario?.inputs.fegliBasicHeldIntoRetirement;
      activeScenario?.inputs.fegliBasicPostRetirementReductionChoice;
      activeScenario?.inputs.fehbBiweeklyDeduction;
      activeScenario?.inputs.fehbProjectedPremiumIncreaseRate;
      activeScenario?.inputs.fehbContinuedInRetirement;
      activeScenario?.inputs.fltcipHasCoverage;
      activeScenario?.inputs.fltcipDailyBenefitAmount;
      activeScenario?.inputs.fltcipBenefitPeriod;
      activeScenario?.inputs.fltcipCurrentAnnualPremium;
      // Track financial planning inputs
      activeScenario?.inputs.userAssumedGeneralInflationRate;
      activeScenario?.inputs.desiredAnnualSpendingInRetirement;
      activeScenario?.inputs.desiredAnnualRetirementSpending;
      activeScenario?.inputs.otherTaxableAccountBalance;
      activeScenario?.inputs.otherTaxableAccountExpectedReturn;
      activeScenario?.inputs.otherTaxableAccountAssetAllocation;
      activeScenario?.inputs.homeValue;
      activeScenario?.inputs.mortgageBalance;
      activeScenario?.inputs.otherSignificantDebtsTotal;
      activeScenario?.inputs.otherRecurringIncomeSources;
      activeScenario?.inputs.oneTimeIncomeEvents;
      activeScenario?.inputs.oneTimeExpenseEvents;
      triggerCalculation();
    }
  });

  // --- Scenario Management Functions ---
  function addNewScenario() {
    const result = addScenarioService(scenarios);
    scenarios = result.updatedScenarios;
    activeScenarioId = result.newScenarioId;
  }

  function deleteActiveScenario() {
    if (!activeScenarioId) return; // Guard: No active scenario to delete

    // UI should ideally prevent deleting the very last scenario or provide a reset option.
    // The service's deleteScenario ensures the list is never empty by adding a default if needed.

    const currentActiveIdBeforeDeletion = activeScenarioId;
    const { updatedScenarios } = deleteScenarioService(
      scenarios,
      activeScenarioId
    );
    scenarios = updatedScenarios;

    // Determine the new activeScenarioId:
    if (scenarios.length === 1) {
      // If only one scenario remains (either the one not deleted or a new default from the service),
      // make it active.
      activeScenarioId = scenarios[0].id;
    } else if (scenarios.find((s) => s.id === currentActiveIdBeforeDeletion)) {
      // If the scenario that was active before this operation is still in the list
      // (e.g., it wasn't the one targeted for deletion, or deletion failed to remove it),
      // keep it as active. This branch might be less common if activeScenarioId is always the target.
      activeScenarioId = currentActiveIdBeforeDeletion;
    } else {
      // If the previously active scenario was successfully deleted and there are multiple scenarios left,
      // default to making the first one in the updated list active.
      activeScenarioId = scenarios[0]?.id; // The ?.id handles if scenarios somehow became empty, though service prevents this.
    }
  }

  function duplicateActiveScenario() {
    if (!activeScenarioId) return; // activeScenario is derived, so check activeScenarioId
    const result = duplicateScenarioService(scenarios, activeScenarioId);
    if (result) {
      scenarios = result.updatedScenarios;
      activeScenarioId = result.newScenarioId;
    } else {
      // This case (scenario to duplicate not found) should be rare if activeScenarioId is valid.
      console.error(
        "Failed to duplicate scenario: Active scenario not found by service."
      );
    }
  }

  function updateActiveScenarioName(newName: string) {
    if (activeScenario) {
      // Update the name in the scenarios array directly
      const scenarioIndex = scenarios.findIndex(
        (s) => s.id === activeScenarioId
      );
      if (scenarioIndex !== -1) {
        scenarios[scenarioIndex].name = newName;
        // Svelte's reactivity should pick this up if 'scenarios' is $state
        // and activeScenario is $derived. Force re-assignment if necessary for deep updates in older Svelte or complex cases.
        // scenarios = [...scenarios]; // Uncomment if reactivity isn't triggered for name change
      }
    }
  }

  // --- File Operation Handler Functions ---

  // --- Formatting ---
  function formatCurrency(
    value: number | null | undefined,
    defaultString: string = "N/A"
  ): string {
    if (value == null) return defaultString;
    return new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: "USD",
    }).format(value);
  }

  // --- Main Calculation Logic ---
  async function handleCalculation(currentScenario: Scenario | undefined) {
    if (!currentScenario) return;

    validationError = null; // Clear global validation error message
    currentScenario.error = null; // Clear scenario-specific error
    currentScenario.calculating = true;
    currentScenario.results = null; // Clear previous results

    try {
      // DEBUG: Log scenario inputs before calculation
      // console.log(
      //   "DEBUG: Scenario inputs before performCalculation:",
      //   JSON.stringify(currentScenario.inputs, null, 2)
      // );
      const calcResult: CalculationServiceResult = await performCalculation(
        currentScenario.inputs
      );

      if (calcResult.validationMessages.length > 0) {
        validationError = calcResult.validationMessages.join("\n");
        currentScenario.error = validationError; // Also set on scenario for context
      } else if (calcResult.error) {
        currentScenario.error = calcResult.error;
        // Optionally set global validationError too if it's a general calculation error
        // validationError = calcResult.error;
      } else {
        // Success
        currentScenario.results = calcResult.results;
      }
    } catch (e: any) {
      // This catch is for unexpected errors during the service call itself or result processing
      console.error("Error in handleCalculation:", e);
      currentScenario.error = e.message || e.toString();
      validationError = currentScenario.error; // Show in UI
    } finally {
      currentScenario.calculating = false;
    }
  }

  // TabButton remains as it was - it's Svelte markup, not TSX/JSX, so it's defined outside typical TS function returns
  // Its usage is in the template directly.
  // --- View Switching State ---
  let activeView = $state<string>("home");
</script>

<div class="flex h-screen bg-gray-50 dark:bg-gray-950 font-sans">
  <Sidebar {activeView} on:changeView={(e) => (activeView = e.detail.view)} />
  <div class="flex-1 flex flex-col min-w-0">
    <TopBar
      {scenarios}
      {activeScenarioId}
      on:openFile={handleOpenFile}
      on:saveFile={handleSaveFile}
      on:selectScenario={(e) => (activeScenarioId = e.detail)}
      on:addScenario={addNewScenario}
      on:duplicateScenario={duplicateActiveScenario}
      on:deleteScenario={deleteActiveScenario}
      on:renameScenario={(e) => renameScenario(e.detail.id, e.detail.name)}
    />
    <main class="flex-1 overflow-y-auto p-4 md:p-6 lg:p-8">
      {#if activeView === "home"}
        <!-- Pass scenarios and activeScenarioId to HomePage to avoid undefined errors -->
        <HomePage {scenarios} {activeScenarioId} />
      {:else if activeView === "scenario"}
        <ScenarioEditorPage
          bind:scenarios
          {activeScenarioId}
          setActiveScenarioId={(id) => (activeScenarioId = id)}
          {activeTab}
          setActiveTab={(tab) => (activeTab = tab)}
          {addNewScenario}
          {deleteActiveScenario}
          {duplicateActiveScenario}
          {TABS}
        />
      {:else if activeView === "results"}
        <AnalysisPage {scenarios} {activeScenarioId} />
      {:else if activeView === "income-projections"}
        <IncomeProjectionsPage {scenarios} {activeScenarioId} />
      {:else if activeView === "scenario-comparison"}
        <ScenarioComparisonPage {scenarios} {activeScenarioId} />
      {:else if activeView === "risk-analysis"}
        <RiskAnalysisPage {scenarios} {activeScenarioId} />
      {:else if activeView === "household"}
        <HouseholdViewPage />
      {:else if activeView === "reports"}
        <ReportsExportsPage />
      {:else if activeView === "settings"}
        <SettingsPage inputs={activeScenario?.inputs} />
      {:else}
        <div class="text-gray-500 dark:text-gray-400">
          This section is under construction.
        </div>
      {/if}
    </main>
  </div>
  <QuickPanel />
</div>

<style>
  /* Global styles for desktop application */
  /* Full width layout for desktop - no artificial constraints */
</style>
