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
    getRetirementEvents,
    createCumulativeData,
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

  // Get active scenario and results
  let activeScenario = $derived(
    scenarios.find((s: Scenario) => s.id === activeScenarioId)
  );
  let results = $derived(activeScenario?.results);
  let hasResults = $derived(
    !!results &&
      (results.TotalRetirementIncome > 0 ||
        results.FERSResult ||
        results.CSRSResult ||
        results.SocialSecurityResult ||
        results.TSPResult)
  );

  // Chart instances
  let incomeChartCanvas = $state<HTMLCanvasElement>();
  let cumulativeChartCanvas = $state<HTMLCanvasElement>();
  let incomeChart: Chart | null = null;
  let cumulativeChart: Chart | null = null;

  // Data processing
  let monthlyData = $derived(() => {
    if (!hasResults || !results) return [];
    return extractMonthlyProjections(results);
  });

  let yearlyData = $derived(() => {
    const monthly = monthlyData();
    if (monthly.length === 0) return [];
    return aggregateToYearly(monthly);
  });

  let chartData = $derived(() => {
    const isMonthly = $displayMode === "monthly";
    const data = isMonthly ? monthlyData() : yearlyData();

    if (data.length === 0) return null;

    return {
      labels: createTimeLabels(data, isMonthly),
      datasets: createIncomeDatasets(data, isMonthly),
      isMonthly,
      data,
    };
  });

  let cumulativeData = $derived(() => {
    const monthly = monthlyData();
    if (monthly.length === 0) return null;

    const cumulative = createCumulativeData(monthly);
    const isMonthly = $displayMode === "monthly";
    const data = isMonthly ? cumulative : aggregateToYearly(cumulative);

    return {
      labels: createTimeLabels(data, isMonthly),
      datasets: createIncomeDatasets(data, isMonthly),
      isMonthly,
      data,
    };
  });

  // Chart configuration
  let incomeChartConfig = $derived(() => {
    const currentChartData = chartData();
    if (!currentChartData) return null;

    return {
      type: "line" as const,
      data: {
        labels: currentChartData.labels,
        datasets: currentChartData.datasets,
      },
      options: {
        ...defaultChartOptions,
        plugins: {
          ...defaultChartOptions.plugins,
          title: {
            display: true,
            text: `Monthly Income Streams Over Time (${
              currentChartData.isMonthly ? "Monthly" : "Annual"
            } Stacked)`,
          },
          tooltip: {
            ...defaultChartOptions.plugins?.tooltip,
            callbacks: {
              label: function (context: any) {
                const label = context.dataset.label || "";
                const value = context.parsed.y;
                return `${label}: ${formatChartCurrency(
                  value,
                  currentChartData.isMonthly
                )}`;
              },
            },
          },
        },
        scales: {
          ...defaultChartOptions.scales,
          x: {
            ...(defaultChartOptions.scales?.x || {}),
            stacked: true,
          },
          y: {
            ...defaultChartOptions.scales?.y,
            beginAtZero: true,
            stacked: true,
            ticks: {
              callback: function (value: any) {
                return formatChartCurrency(value, currentChartData.isMonthly);
              },
            },
          },
        },
        interaction: {
          mode: "index" as const,
          intersect: false,
        },
      },
    };
  });

  let cumulativeChartConfig = $derived(() => {
    const currentCumulativeData = cumulativeData();
    if (!currentCumulativeData) return null;

    return {
      type: "line" as const,
      data: {
        labels: currentCumulativeData.labels,
        datasets: currentCumulativeData.datasets.map((dataset) => ({
          ...dataset,
          fill: false, // Don't fill for cumulative chart
        })),
      },
      options: {
        ...defaultChartOptions,
        plugins: {
          ...defaultChartOptions.plugins,
          title: {
            display: true,
            text: `Cumulative Income Over Time (${
              currentCumulativeData.isMonthly ? "Monthly" : "Annual"
            })`,
          },
          tooltip: {
            ...defaultChartOptions.plugins?.tooltip,
            callbacks: {
              label: function (context: any) {
                const label = context.dataset.label || "";
                const value = context.parsed.y;
                return `${label}: ${formatChartCurrency(
                  value,
                  currentCumulativeData.isMonthly
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
                return formatChartCurrency(
                  value,
                  currentCumulativeData.isMonthly
                );
              },
            },
          },
        },
      },
    };
  });

  // Chart lifecycle management
  function createCharts() {
    if (incomeChartCanvas && incomeChartConfig()) {
      incomeChart = createChartInstance(
        incomeChartCanvas,
        incomeChartConfig()!
      );
    }
    if (cumulativeChartCanvas && cumulativeChartConfig()) {
      cumulativeChart = createChartInstance(
        cumulativeChartCanvas,
        cumulativeChartConfig()!
      );
    }
  }

  function updateCharts() {
    if (incomeChart && incomeChartConfig()) {
      incomeChart.data = incomeChartConfig()!.data;
      incomeChart.options = incomeChartConfig()!.options;
      incomeChart.update();
    }
    if (cumulativeChart && cumulativeChartConfig()) {
      cumulativeChart.data = cumulativeChartConfig()!.data;
      cumulativeChart.options = cumulativeChartConfig()!.options;
      cumulativeChart.update();
    }
  }

  function destroyCharts() {
    if (incomeChart) {
      incomeChart.destroy();
      incomeChart = null;
    }
    if (cumulativeChart) {
      cumulativeChart.destroy();
      cumulativeChart = null;
    }
  }

  // Effects for chart management
  onMount(() => {
    createCharts();
  });

  $effect(() => {
    // React to display mode changes
    const currentDisplayMode = $displayMode;
    if (incomeChart || cumulativeChart) {
      updateCharts();
    } else {
      createCharts();
    }
  });

  $effect(() => {
    // React to data changes
    const currentChartData = chartData();
    const currentCumulativeData = cumulativeData();
    if (currentChartData || currentCumulativeData) {
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

  // Key insights
  let keyInsights = $derived(() => {
    const data = $displayMode === "monthly" ? monthlyData() : yearlyData();
    if (data.length === 0) return null;

    const firstPeriod = data[0];
    const midPeriod = data[Math.floor(data.length / 2)] || firstPeriod;
    const lastPeriod = data[data.length - 1];

    return {
      first: firstPeriod,
      mid: midPeriod,
      last: lastPeriod,
      totalPeriods: data.length,
    };
  });
</script>

<div class="flex flex-col gap-6 w-full">
  <!-- Header -->
  <div class="flex justify-between items-center">
    <h1 class="text-3xl font-bold text-blue-900 dark:text-blue-200">
      Income Projections
    </h1>
    <div class="flex items-center gap-4">
      {#if activeScenario}
        <div class="text-sm text-gray-600 dark:text-gray-400">
          Scenario: <span class="font-semibold">{activeScenario.name}</span>
        </div>
      {/if}
      {#if hasResults}
        <DisplayToggle />
      {/if}
    </div>
  </div>

  {#if !hasResults}
    <!-- No Results State -->
    <div
      class="bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg p-6 text-center"
    >
      <div
        class="text-yellow-800 dark:text-yellow-200 text-lg font-semibold mb-2"
      >
        No Income Projections Available
      </div>
      <div class="text-yellow-700 dark:text-yellow-300">
        Please complete the scenario inputs and run calculations to see income
        projections.
      </div>
    </div>
  {:else}
    <!-- Charts Section -->
    <div class="grid grid-cols-1 gap-6">
      <!-- Income Over Time Chart -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-6">
        <div class="h-96">
          {#if incomeChartConfig()}
            <canvas bind:this={incomeChartCanvas}></canvas>
          {:else}
            <div
              class="h-full flex items-center justify-center text-gray-500 dark:text-gray-400"
            >
              <div class="text-center">
                <div class="text-lg font-medium mb-2">Loading Chart...</div>
                <div class="text-sm">Processing income projection data</div>
              </div>
            </div>
          {/if}
        </div>
      </div>

      <!-- Cumulative Income Chart -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-6">
        <div class="h-96">
          {#if cumulativeChartConfig()}
            <canvas bind:this={cumulativeChartCanvas}></canvas>
          {:else}
            <div
              class="h-full flex items-center justify-center text-gray-500 dark:text-gray-400"
            >
              <div class="text-center">
                <div class="text-lg font-medium mb-2">Loading Chart...</div>
                <div class="text-sm">Processing cumulative income data</div>
              </div>
            </div>
          {/if}
        </div>
      </div>
    </div>

    <!-- Key Insights -->
    {#if keyInsights()}
      {@const insights = keyInsights()}
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-6">
        <h3 class="text-xl font-bold text-gray-900 dark:text-gray-100 mb-4">
          Key Insights
        </h3>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div class="bg-blue-50 dark:bg-blue-900/20 p-4 rounded-lg">
            <div class="font-medium text-blue-800 dark:text-blue-200">
              At Retirement
            </div>
            <div class="text-sm text-blue-600 dark:text-blue-300 mt-1">
              Age {insights.first.ageYears}
            </div>
            <div
              class="text-lg font-bold text-blue-900 dark:text-blue-100 mt-2"
            >
              {formatChartCurrency(
                insights.first.totalPreTaxIncome,
                $displayMode === "monthly"
              )}
            </div>
            <div class="text-xs text-blue-600 dark:text-blue-400 mt-1">
              Pension: {formatChartCurrency(
                insights.first.pension,
                $displayMode === "monthly"
              )}
            </div>
            <div class="text-xs text-blue-600 dark:text-blue-400">
              Social Security: {formatChartCurrency(
                insights.first.socialSecurity,
                $displayMode === "monthly"
              )}
            </div>
            <div class="text-xs text-blue-600 dark:text-blue-400">
              TSP: {formatChartCurrency(
                insights.first.totalTspWithdrawal,
                $displayMode === "monthly"
              )}
            </div>
          </div>

          <div class="bg-green-50 dark:bg-green-900/20 p-4 rounded-lg">
            <div class="font-medium text-green-800 dark:text-green-200">
              Mid-Retirement
            </div>
            <div class="text-sm text-green-600 dark:text-green-300 mt-1">
              Age {insights.mid.ageYears}
            </div>
            <div
              class="text-lg font-bold text-green-900 dark:text-green-100 mt-2"
            >
              {formatChartCurrency(
                insights.mid.totalPreTaxIncome,
                $displayMode === "monthly"
              )}
            </div>
            <div class="text-xs text-green-600 dark:text-green-400 mt-1">
              Pension: {formatChartCurrency(
                insights.mid.pension,
                $displayMode === "monthly"
              )}
            </div>
            <div class="text-xs text-green-600 dark:text-green-400">
              Social Security: {formatChartCurrency(
                insights.mid.socialSecurity,
                $displayMode === "monthly"
              )}
            </div>
            <div class="text-xs text-green-600 dark:text-green-400">
              TSP: {formatChartCurrency(
                insights.mid.totalTspWithdrawal,
                $displayMode === "monthly"
              )}
            </div>
          </div>

          <div class="bg-purple-50 dark:bg-purple-900/20 p-4 rounded-lg">
            <div class="font-medium text-purple-800 dark:text-purple-200">
              Later Years
            </div>
            <div class="text-sm text-purple-600 dark:text-purple-300 mt-1">
              Age {insights.last.ageYears}
            </div>
            <div
              class="text-lg font-bold text-purple-900 dark:text-purple-100 mt-2"
            >
              {formatChartCurrency(
                insights.last.totalPreTaxIncome,
                $displayMode === "monthly"
              )}
            </div>
            <div class="text-xs text-purple-600 dark:text-purple-400 mt-1">
              Pension: {formatChartCurrency(
                insights.last.pension,
                $displayMode === "monthly"
              )}
            </div>
            <div class="text-xs text-purple-600 dark:text-purple-400">
              Social Security: {formatChartCurrency(
                insights.last.socialSecurity,
                $displayMode === "monthly"
              )}
            </div>
            <div class="text-xs text-purple-600 dark:text-purple-400">
              TSP: {formatChartCurrency(
                insights.last.totalTspWithdrawal,
                $displayMode === "monthly"
              )}
            </div>
          </div>
        </div>
      </div>
    {/if}

    <!-- Data Summary -->
    {#if monthlyData().length > 0}
      <div class="bg-gray-50 dark:bg-gray-900 rounded-lg p-4">
        <div class="text-sm text-gray-600 dark:text-gray-400">
          <strong>Data Source:</strong> Using real backend monthly projections ({monthlyData()
            .length} months of data)
          <br />
          <strong>Display Mode:</strong>
          {$displayMode === "monthly"
            ? "Monthly view"
            : "Yearly aggregated view"}
          <br />
          <strong>Coverage:</strong> Age {monthlyData()[0]?.ageYears} to Age {monthlyData()[
            monthlyData().length - 1
          ]?.ageYears}
        </div>
      </div>

      <!-- TSP Modeling Explanation -->
    {/if}
  {/if}
</div>
