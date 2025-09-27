<script lang="ts">
  import type {
    RetirementCalculationResult,
    ScenarioInput,
  } from "../../../types";
  import IncomeOverTimeChart from "./IncomeOverTimeChart.svelte";
  import TSPBalanceChart from "./TSPBalanceChart.svelte";
  import IncomeBreakdownChart from "./IncomeBreakdownChart.svelte";
  import TaxBurdenChart from "./TaxBurdenChart.svelte";
  import MonteCarloChart from "./MonteCarloChart.svelte";
  // import RiskAssessmentChart from "./RiskAssessmentChart.svelte";
  import ScenarioComparisonChart from "./ScenarioComparisonChart.svelte";
  import DisplayToggle from "../ui/DisplayToggle.svelte";
  import { displayMode } from "../../stores/displayPreferences";
  let {
    results,
    inputs,
    title = "Charts & Projections",
  } = $props<{
    results: RetirementCalculationResult | undefined;
    inputs: ScenarioInput | undefined;
    title?: string;
  }>();

  // Tab state
  let activeTab = $state("overview");

  // Tab definitions
  const tabs = [
    { id: "overview", label: "Overview", icon: "📊" },
    { id: "income", label: "Income Analysis", icon: "💰" },
    { id: "risk", label: "Risk Assessment", icon: "⚠️" },
    { id: "scenarios", label: "Scenario Comparison", icon: "🔄" },
  ];

  // Derive chart data from results
  let chartData = $derived(() => {
    if (!results || !inputs) {
      return null;
    }

    return {
      totalIncome: results.TotalRetirementIncome || 0,
      pensionAnnual: results.FERSResult?.GrossAnnualAnnuity || 0,
      socialSecurityAnnual: results.SocialSecurityResult?.AnnualBenefit || 0,
      tspBalance: results.TSPResult?.totalProjectedBalanceAtRetirement || 0,
      tspTraditional: results.TSPResult?.projectedTraditionalBalanceAtRetirement || 0,
      tspRoth: results.TSPResult?.projectedRothBalanceAtRetirement || 0,
      retirementAge: inputs.plannedRetirementDate
        ? Math.floor(
            (new Date(inputs.plannedRetirementDate).getTime() -
              new Date(inputs.dateOfBirth).getTime()) /
              (365.25 * 24 * 60 * 60 * 1000)
          )
        : 62,
      federalTax: results.TaxResult?.FederalTaxOwed || 0,
      stateTax: results.TaxResult?.StateTaxOwed || 0,
      netIncome: results.NetAfterTaxIncome || 0,
      effectiveTaxRate: results.EffectiveTaxRate || 0,
      monteCarloResults: results.MonteCarloResult,
    };
  });

  // Create a signal for chart data
  let chartDataSignal = $derived(() => chartData());
  let inputsSignal = $derived(() => inputs);

  $effect(() => {
    const chartDataSnapshot = $state.snapshot(chartData);
    // console.log(
    //   "ChartsSection: $effect triggered because chartData changed. Current chartData (snapshot):",
    //   chartDataSnapshot
    // );
    if (chartDataSnapshot) {
      // console.log(
      //   "ChartsSection: $effect - chartData is TRUTHY. Keys:",
      //   Object.keys(chartDataSnapshot)
      // );
    } else {
      // console.log(
      //   "ChartsSection: $effect - chartData is FALSY (null or undefined). This will cause 'No Data Available' message."
      // );
    }
  });
  // Effect to track display mode changes
  $effect(() => {
    // console.log("ChartsSection: Display mode changed to", $displayMode);
  });

  // Log initial prop values when the component instance is created
  // console.log(
  //   "ChartsSection: Component instance created. Initial props (snapshots) - results:",
  //   $state.snapshot(results),
  //   "inputs:",
  //   $state.snapshot(inputs)
  // );
</script>

<div class="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-6">
  <div class="flex justify-between items-center mb-6">
    <h2 class="text-2xl font-bold text-gray-900 dark:text-gray-100">
      {title}
    </h2>
    {#if chartData()}
      <DisplayToggle />
    {/if}
  </div>
  {#if chartData()}
    <!-- Tab Navigation -->
    <div class="border-b border-gray-200 dark:border-gray-600 mb-6">
      <nav class="-mb-px flex space-x-8">
        {#each tabs as tab}
          <button
            class="py-2 px-1 border-b-2 font-medium text-sm {activeTab ===
            tab.id
              ? 'border-blue-500 text-blue-600 dark:text-blue-400'
              : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300 dark:text-gray-400 dark:hover:text-gray-300'}"
            onclick={() => (activeTab = tab.id)}
          >
            <span class="mr-2">{tab.icon}</span>
            {tab.label}
          </button>
        {/each}
      </nav>
    </div>

    <!-- Tab Content -->
    <div class="space-y-8">
      {#if activeTab === "overview"}
        <div class="space-y-6">
          <IncomeOverTimeChart data={chartDataSignal} inputs={inputsSignal} />
          <TSPBalanceChart data={chartDataSignal} inputs={inputsSignal} />
        </div>
      {:else if activeTab === "income"}
        <div class="space-y-6">
          <IncomeBreakdownChart data={chartDataSignal} />
          <TaxBurdenChart data={chartDataSignal} inputs={inputsSignal} />
        </div>
      {:else if activeTab === "risk"}
        <div class="space-y-6">
          <MonteCarloChart data={chartDataSignal} inputs={inputsSignal} />
          <!-- <RiskAssessmentChart data={chartDataSignal} /> -->
        </div>
      {:else if activeTab === "scenarios"}
        <div class="space-y-6">
          <ScenarioComparisonChart scenarios={[]} />
        </div>
      {/if}
    </div>
  {:else}
    <div class="text-center text-gray-500 dark:text-gray-400 py-8">
      <div class="text-lg font-medium mb-2">No Data Available</div>
      <div class="text-sm">
        Complete scenario inputs and run calculations to view charts
      </div>
    </div>
  {/if}
</div>
