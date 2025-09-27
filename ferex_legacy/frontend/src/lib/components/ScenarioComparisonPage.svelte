<script lang="ts">
  import type { Scenario } from "../../types";
  import { displayMode } from "../stores/displayPreferences";
  import DisplayToggle from "./ui/DisplayToggle.svelte";
  import {
    extractMonthlyProjections,
    aggregateToYearly,
    createIncomeDatasets,
    createTimeLabels,
    formatChartCurrency,
    calculateScenarioDifferences,
    type MonthlyDataPoint,
    type YearlyDataPoint,
  } from "../utils/chartDataUtils";
  import {
    Chart,
    createChartInstance,
    chartColors,
    defaultChartOptions,
  } from "./charts/chartUtils";
  import { onMount, onDestroy } from "svelte";

  let { scenarios = [], activeScenarioId } = $props<{
    scenarios?: Scenario[];
    activeScenarioId?: string | undefined;
  }>();

  // Scenario selection for comparison
  let scenarioAId = $state<string | undefined>(activeScenarioId);
  let scenarioBId = $state<string | undefined>(undefined);

  // Get scenarios and results
  let scenarioA = $derived(
    scenarios.find((s: Scenario) => s.id === scenarioAId)
  );
  let scenarioB = $derived(
    scenarios.find((s: Scenario) => s.id === scenarioBId)
  );

  let resultsA = $derived(scenarioA?.results);
  let resultsB = $derived(scenarioB?.results);

  let hasResultsA = $derived(
    !!resultsA &&
      (resultsA.TotalRetirementIncome > 0 ||
        resultsA.FERSResult ||
        resultsA.CSRSResult ||
        resultsA.SocialSecurityResult ||
        resultsA.TSPResult)
  );

  let hasResultsB = $derived(
    !!resultsB &&
      (resultsB.TotalRetirementIncome > 0 ||
        resultsB.FERSResult ||
        resultsB.CSRSResult ||
        resultsB.SocialSecurityResult ||
        resultsB.TSPResult)
  );

  let canCompare = $derived(hasResultsA && hasResultsB);

  // Chart instances
  let comparisonChartCanvas = $state<HTMLCanvasElement>();
  let deltaChartCanvas = $state<HTMLCanvasElement>();
  let comparisonChart: Chart | null = null;
  let deltaChart: Chart | null = null;

  // Data processing
  let monthlyDataA = $derived(() => {
    if (!hasResultsA || !resultsA) return [];
    return extractMonthlyProjections(resultsA);
  });

  let monthlyDataB = $derived(() => {
    if (!hasResultsB || !resultsB) return [];
    return extractMonthlyProjections(resultsB);
  });

  let yearlyDataA = $derived(() => {
    const monthly = monthlyDataA();
    if (monthly.length === 0) return [];
    return aggregateToYearly(monthly);
  });

  let yearlyDataB = $derived(() => {
    const monthly = monthlyDataB();
    if (monthly.length === 0) return [];
    return aggregateToYearly(monthly);
  });

  let comparisonData = $derived(() => {
    if (!canCompare) return null;

    const isMonthly = $displayMode === "monthly";
    const dataA = isMonthly ? monthlyDataA() : yearlyDataA();
    const dataB = isMonthly ? monthlyDataB() : yearlyDataB();

    if (dataA.length === 0 || dataB.length === 0) return null;

    return {
      labels: createTimeLabels(dataA, isMonthly),
      datasetsA: createIncomeDatasets(dataA, isMonthly).map((dataset) => ({
        ...dataset,
        label: `${dataset.label} (${scenarioA?.name || "Scenario A"})`,
        borderDash: [],
      })),
      datasetsB: createIncomeDatasets(dataB, isMonthly).map((dataset) => ({
        ...dataset,
        label: `${dataset.label} (${scenarioB?.name || "Scenario B"})`,
        borderDash: [5, 5], // Dashed lines for scenario B
        backgroundColor: dataset.backgroundColor?.replace("0.2", "0.1"), // More transparent
      })),
      isMonthly,
      dataA,
      dataB,
    };
  });

  let deltaData = $derived(() => {
    if (!canCompare) return null;

    const monthlyA = monthlyDataA();
    const monthlyB = monthlyDataB();

    if (monthlyA.length === 0 || monthlyB.length === 0) return null;

    const differences = calculateScenarioDifferences(monthlyA, monthlyB);
    const isMonthly = $displayMode === "monthly";
    const data = isMonthly ? differences : aggregateToYearly(differences);

    return {
      labels: createTimeLabels(data, isMonthly),
      datasets: [
        {
          label: "Pension Difference",
          data: data.map((d) => d.pension),
          backgroundColor: "rgba(59, 130, 246, 0.3)",
          borderColor: "rgb(59, 130, 246)",
          borderWidth: 2,
        },
        {
          label: "Social Security Difference",
          data: data.map((d) => d.socialSecurity),
          backgroundColor: "rgba(16, 185, 129, 0.3)",
          borderColor: "rgb(16, 185, 129)",
          borderWidth: 2,
        },
        {
          label: "TSP Difference",
          data: data.map((d) => d.totalTspWithdrawal),
          backgroundColor: "rgba(139, 92, 246, 0.3)",
          borderColor: "rgb(139, 92, 246)",
          borderWidth: 2,
        },
        {
          label: "Total Income Difference",
          data: data.map((d) => d.totalPreTaxIncome),
          backgroundColor: "rgba(239, 68, 68, 0.3)",
          borderColor: "rgb(239, 68, 68)",
          borderWidth: 3,
          type: "line" as const,
        },
      ],
      isMonthly,
      data,
    };
  });

  // Chart configurations
  let comparisonChartConfig = $derived(() => {
    const currentData = comparisonData();
    if (!currentData) return null;

    return {
      type: "line" as const,
      data: {
        labels: currentData.labels,
        datasets: [...currentData.datasetsA, ...currentData.datasetsB],
      },
      options: {
        ...defaultChartOptions,
        plugins: {
          ...defaultChartOptions.plugins,
          title: {
            display: true,
            text: `Income Comparison: ${scenarioA?.name || "Scenario A"} vs ${
              scenarioB?.name || "Scenario B"
            } (${currentData.isMonthly ? "Monthly" : "Annual"})`,
          },
          tooltip: {
            ...defaultChartOptions.plugins?.tooltip,
            callbacks: {
              label: function (context: any) {
                const label = context.dataset.label || "";
                const value = context.parsed.y;
                return `${label}: ${formatChartCurrency(
                  value,
                  currentData.isMonthly
                )}`;
              },
            },
          },
        },
        scales: {
          ...defaultChartOptions.scales,
          y: {
            ...defaultChartOptions.scales?.y,
            beginAtZero: true,
            ticks: {
              callback: function (value: any) {
                return formatChartCurrency(value, currentData.isMonthly);
              },
            },
          },
        },
      },
    };
  });

  let deltaChartConfig = $derived(() => {
    const currentData = deltaData();
    if (!currentData) return null;

    return {
      type: "bar" as const,
      data: {
        labels: currentData.labels,
        datasets: currentData.datasets,
      },
      options: {
        ...defaultChartOptions,
        plugins: {
          ...defaultChartOptions.plugins,
          title: {
            display: true,
            text: `Income Differences: ${scenarioB?.name || "Scenario B"} - ${
              scenarioA?.name || "Scenario A"
            } (${currentData.isMonthly ? "Monthly" : "Annual"})`,
          },
          tooltip: {
            ...defaultChartOptions.plugins?.tooltip,
            callbacks: {
              label: function (context: any) {
                const label = context.dataset.label || "";
                const value = context.parsed.y;
                const sign = value >= 0 ? "+" : "";
                return `${label}: ${sign}${formatChartCurrency(
                  value,
                  currentData.isMonthly
                )}`;
              },
            },
          },
        },
        scales: {
          ...defaultChartOptions.scales,
          y: {
            ...defaultChartOptions.scales?.y,
            beginAtZero: true,
            ticks: {
              callback: function (value: any) {
                const sign = value >= 0 ? "+" : "";
                return `${sign}${formatChartCurrency(value, currentData.isMonthly)}`;
              },
            },
          },
        },
      },
    };
  });

  // Chart lifecycle management
  function createCharts() {
    if (comparisonChartCanvas && comparisonChartConfig()) {
      comparisonChart = createChartInstance(
        comparisonChartCanvas,
        comparisonChartConfig()!
      );
    }
    if (deltaChartCanvas && deltaChartConfig()) {
      deltaChart = createChartInstance(deltaChartCanvas, deltaChartConfig()!);
    }
  }

  function updateCharts() {
    if (comparisonChart && comparisonChartConfig()) {
      comparisonChart.data = comparisonChartConfig()!.data;
      comparisonChart.options = comparisonChartConfig()!.options;
      comparisonChart.update();
    }
    if (deltaChart && deltaChartConfig()) {
      deltaChart.data = deltaChartConfig()!.data;
      deltaChart.options = deltaChartConfig()!.options;
      deltaChart.update();
    }
  }

  function destroyCharts() {
    if (comparisonChart) {
      comparisonChart.destroy();
      comparisonChart = null;
    }
    if (deltaChart) {
      deltaChart.destroy();
      deltaChart = null;
    }
  }

  // Effects for chart management
  onMount(() => {
    createCharts();
  });

  $effect(() => {
    // React to display mode changes
    const currentDisplayMode = $displayMode;
    if (comparisonChart || deltaChart) {
      updateCharts();
    } else {
      createCharts();
    }
  });

  $effect(() => {
    // React to data changes
    const currentComparisonData = comparisonData();
    const currentDeltaData = deltaData();
    if (currentComparisonData || currentDeltaData) {
      updateCharts();
    }
  });

  onDestroy(() => {
    destroyCharts();
  });

  // Helper functions
  function formatCurrency(
    value: number | null | undefined,
    defaultText = "N/A"
  ): string {
    if (value == null || isNaN(value)) return defaultText;
    return new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: "USD",
      maximumFractionDigits: 0,
    }).format(value);
  }

  // Summary comparison
  let summaryComparison = $derived(() => {
    if (!canCompare || !resultsA || !resultsB) return null;

    return {
      totalIncomeDiff: (resultsB.TotalRetirementIncome || 0) - (resultsA.TotalRetirementIncome || 0),
      netIncomeDiff: (resultsB.NetAfterTaxIncome || 0) - (resultsA.NetAfterTaxIncome || 0),
      taxRateDiff: (resultsB.EffectiveTaxRate || 0) - (resultsA.EffectiveTaxRate || 0),
      tspBalanceDiff: 
        (resultsB.TSPResult?.totalProjectedBalanceAtRetirement || 0) - 
        (resultsA.TSPResult?.totalProjectedBalanceAtRetirement || 0),
    };
  });

  // Auto-select second scenario if only one is available
  $effect(() => {
    if (scenarios.length >= 2 && !scenarioBId) {
      const otherScenario = scenarios.find(s => s.id !== scenarioAId);
      if (otherScenario) {
        scenarioBId = otherScenario.id;
      }
    }
  });
</script>

<div class="flex flex-col gap-6 w-full">
  <!-- Header -->
  <div class="flex justify-between items-center">
    <h1 class="text-3xl font-bold text-blue-900 dark:text-blue-200">
      Scenario Comparison
    </h1>
    <div class="flex items-center gap-4">
      {#if canCompare}
        <DisplayToggle />
      {/if}
    </div>
  </div>

  <!-- Scenario Selection -->
  <div class="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-6">
    <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">
      Select Scenarios to Compare
    </h3>
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          Scenario A (Baseline)
        </label>
        <select
          bind:value={scenarioAId}
          class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-gray-100"
        >
          <option value={undefined}>Select a scenario...</option>
          {#each scenarios as scenario}
            <option value={scenario.id}>{scenario.name}</option>
          {/each}
        </select>
      </div>
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          Scenario B (Comparison)
        </label>
        <select
          bind:value={scenarioBId}
          class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-gray-100"
        >
          <option value={undefined}>Select a scenario...</option>
          {#each scenarios as scenario}
            {#if scenario.id !== scenarioAId}
              <option value={scenario.id}>{scenario.name}</option>
            {/if}
          {/each}
        </select>
      </div>
    </div>
  </div>

  {#if !canCompare}
    <!-- No Comparison State -->
    <div
      class="bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg p-6 text-center"
    >
      <div
        class="text-yellow-800 dark:text-yellow-200 text-lg font-semibold mb-2"
      >
        {scenarios.length < 2 
          ? "Need More Scenarios" 
          : !scenarioAId || !scenarioBId 
            ? "Select Scenarios to Compare"
            : "Missing Calculation Results"}
      </div>
      <div class="text-yellow-700 dark:text-yellow-300">
        {scenarios.length < 2
          ? "Create at least two scenarios with completed calculations to compare them."
          : !scenarioAId || !scenarioBId
            ? "Please select two scenarios from the dropdowns above."
            : "Both scenarios need completed calculations to enable comparison."}
      </div>
    </div>
  {:else}
    <!-- Summary Comparison -->
    {#if summaryComparison()}
      {@const summary = summaryComparison()}
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-6">
        <h3 class="text-xl font-bold text-gray-900 dark:text-gray-100 mb-4">
          Key Differences Summary
        </h3>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          <div class="text-center p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg">
            <div class="text-sm font-medium text-blue-800 dark:text-blue-200">
              Total Income Difference
            </div>
            <div class="text-2xl font-bold text-blue-900 dark:text-blue-100 mt-1">
              {summary.totalIncomeDiff >= 0 ? "+" : ""}{formatCurrency(summary.totalIncomeDiff)}
            </div>
            <div class="text-xs text-blue-600 dark:text-blue-400 mt-1">
              Annual difference
            </div>
          </div>

          <div class="text-center p-4 bg-green-50 dark:bg-green-900/20 rounded-lg">
            <div class="text-sm font-medium text-green-800 dark:text-green-200">
              Net Income Difference
            </div>
            <div class="text-2xl font-bold text-green-900 dark:text-green-100 mt-1">
              {summary.netIncomeDiff >= 0 ? "+" : ""}{formatCurrency(summary.netIncomeDiff)}
            </div>
            <div class="text-xs text-green-600 dark:text-green-400 mt-1">
              After-tax difference
            </div>
          </div>

          <div class="text-center p-4 bg-purple-50 dark:bg-purple-900/20 rounded-lg">
            <div class="text-sm font-medium text-purple-800 dark:text-purple-200">
              Tax Rate Difference
            </div>
            <div class="text-2xl font-bold text-purple-900 dark:text-purple-100 mt-1">
              {summary.taxRateDiff >= 0 ? "+" : ""}{(summary.taxRateDiff * 100).toFixed(1)}%
            </div>
            <div class="text-xs text-purple-600 dark:text-purple-400 mt-1">
              Effective rate change
            </div>
          </div>

          <div class="text-center p-4 bg-orange-50 dark:bg-orange-900/20 rounded-lg">
            <div class="text-sm font-medium text-orange-800 dark:text-orange-200">
              TSP Balance Difference
            </div>
            <div class="text-2xl font-bold text-orange-900 dark:text-orange-100 mt-1">
              {summary.tspBalanceDiff >= 0 ? "+" : ""}{formatCurrency(summary.tspBalanceDiff)}
            </div>
            <div class="text-xs text-orange-600 dark:text-orange-400 mt-1">
              At retirement
            </div>
          </div>
        </div>
      </div>
    {/if}

    <!-- Charts Section -->
    <div class="grid grid-cols-1 gap-6">
      <!-- Side-by-Side Comparison Chart -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-6">
        <div class="h-96">
          {#if comparisonChartConfig()}
            <canvas bind:this={comparisonChartCanvas}></canvas>
          {:else}
            <div
              class="h-full flex items-center justify-center text-gray-500 dark:text-gray-400"
            >
              <div class="text-center">
                <div class="text-lg font-medium mb-2">Loading Comparison Chart...</div>
                <div class="text-sm">Processing scenario comparison data</div>
              </div>
            </div>
          {/if}
        </div>
      </div>

      <!-- Delta/Difference Chart -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-6">
        <div class="h-96">
          {#if deltaChartConfig()}
            <canvas bind:this={deltaChartCanvas}></canvas>
          {:else}
            <div
              class="h-full flex items-center justify-center text-gray-500 dark:text-gray-400"
            >
              <div class="text-center">
                <div class="text-lg font-medium mb-2">Loading Delta Chart...</div>
                <div class="text-sm">Processing difference calculations</div>
              </div>
            </div>
          {/if}
        </div>
      </div>
    </div>

    <!-- Data Summary -->
    <div class="bg-gray-50 dark:bg-gray-900 rounded-lg p-4">
      <div class="text-sm text-gray-600 dark:text-gray-400">
        <strong>Comparison:</strong> {scenarioA?.name || "Scenario A"} vs {scenarioB?.name || "Scenario B"}
        <br />
        <strong>Data Source:</strong> Using real backend monthly projections
        <br />
        <strong>Display Mode:</strong>
        {$displayMode === "monthly" ? "Monthly view" : "Yearly aggregated view"}
      </div>
    </div>
  {/if}
</div>
