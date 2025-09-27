<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import {
    Chart,
    chartColors,
    defaultChartOptions,
    formatCurrency,
    createChartInstance,
    convertToMonthly,
  } from "./chartUtils";
  import { displayMode } from "../../stores/displayPreferences";
  import type { models } from "../../../../wailsjs/go/models"; // Import backend models
  import type { ScenarioInput } from "../../../types";

  // Interface for a single year's data points across percentiles
  interface MonteCarloProjectionYear {
    age: number;
    p10: number;
    p25: number;
    p50: number; // Median
    p75: number;
    p90: number;
  }

  // Keep existing interfaces if they are used elsewhere or for other chart types
  interface TSPBreakdown {
    traditionalPct: number;
    rothPct: number;
  }

  interface ProjectionYear {
    // This might be deprecated or used for a non-Monte Carlo view
    age: number;
    balance: number;
    traditional: number;
    roth: number;
    withdrawn: number;
    totalWithdrawn: number;
  }

  // Expecting propsDataSignal to contain a field like 'monteCarloResults' of type models.MonteCarloResult
  // Also, include other potential fields from the original data structure that might be needed.
  type ChartDataProps = {
    monteCarloResults?: models.MonteCarloResult;
    tspBalance?: number;
    tspTraditional?: number;
    tspRoth?: number;
    retirementAge?: number;
    // Add other fields from the original 'data: any' if they are used by breakdown or deterministic projection
    // For example, if monteCarloResults also contains input parameters:
    monteCarloInput?: models.MonteCarloInput;
  };
  let { data: propsDataSignal, inputs: inputsSignal } = $props<{ data: ChartDataProps; inputs: ScenarioInput | undefined }>();

  $effect(() => {
    // console.log("TSPBalanceChart: propsDataSignal received:", propsDataSignal);
    const currentData = propsDataSignal();
    if (!currentData) {
      // console.log(
      //   "TSPBalanceChart: currentData from propsDataSignal() is null or undefined."
      // );
    } else {
      // console.log(
      //   "TSPBalanceChart: currentData.monteCarloResults (from propsDataSignal()):",
      //   currentData.monteCarloResults
      // );
      // if (!currentData.monteCarloResults) {
      //   console.log("TSPBalanceChart: currentData.monteCarloResults is falsy.");
      // }
    }
  });

  let canvasElement = $state<HTMLCanvasElement | undefined>();
  let chartInstance: Chart | null = $state(null);

  // Helper function to calculate percentile from a sorted array
  function getPercentileValue(
    sortedData: number[],
    percentile: number
  ): number {
    if (!sortedData || sortedData.length === 0) return 0;
    const index = (percentile / 100) * (sortedData.length - 1);
    if (Number.isInteger(index)) {
      return sortedData[index];
    }
    const lower = Math.floor(index);
    const upper = Math.ceil(index);
    // Guard against lower and upper being the same due to small array length
    if (lower === upper) return sortedData[lower];
    return (
      sortedData[lower] +
      (sortedData[upper] - sortedData[lower]) * (index - lower)
    );
  }

  // Calculate Monte Carlo projection data for the fan chart
  let monteCarloProjectionData = $derived(() => {
    const currentData = propsDataSignal();
    const mcResults = currentData?.monteCarloResults;

    if (
      !mcResults ||
      !mcResults.YearlyBalances ||
      mcResults.YearlyBalances.length === 0
    ) {
      // console.log(
      //   "TSPBalanceChart: Monte Carlo data is not available or empty. YearlyBalances:",
      //   mcResults?.YearlyBalances
      // );
      return [];
    }

    const numSimulations = mcResults.YearlyBalances.length;
    if (numSimulations === 0) return [];
    const numYears = mcResults.YearlyBalances[0]?.length || 0;
    if (numYears === 0) return [];

    const startAge = currentData?.retirementAge || 62;

    const projection: MonteCarloProjectionYear[] = [];

    for (let yearIndex = 0; yearIndex < numYears; yearIndex++) {
      const yearlyBalancesForThisYear: number[] = [];
      for (let simIndex = 0; simIndex < numSimulations; simIndex++) {
        // Ensure YearlyBalances[simIndex] is valid and has enough elements
        if (
          mcResults.YearlyBalances[simIndex] &&
          mcResults.YearlyBalances[simIndex].length > yearIndex
        ) {
          yearlyBalancesForThisYear.push(
            mcResults.YearlyBalances[simIndex][yearIndex]
          );
        } else {
          // Handle case where a simulation might not have data for this year (should ideally not happen)
          yearlyBalancesForThisYear.push(0);
        }
      }
      yearlyBalancesForThisYear.sort((a, b) => a - b);

      projection.push({
        age: startAge + yearIndex,
        p10: getPercentileValue(yearlyBalancesForThisYear, 10),
        p25: getPercentileValue(yearlyBalancesForThisYear, 25),
        p50: getPercentileValue(yearlyBalancesForThisYear, 50),
        p75: getPercentileValue(yearlyBalancesForThisYear, 75),
        p90: getPercentileValue(yearlyBalancesForThisYear, 90),
      });
    }
    // console.log(
    //   "TSPBalanceChart: monteCarloProjectionData successfully derived:",
    //   projection
    // );
    return projection;
  });

  // Calculate percentage breakdown - this might still be useful for other parts of the UI
  let breakdown = $derived(() => {
    const currentData = propsDataSignal();
    if (!currentData) return { traditionalPct: 0, rothPct: 0 };

    const total = currentData.tspBalance || 0;
    const traditional = currentData.tspTraditional || 0;
    const roth = currentData.tspRoth || 0;

    if (total === 0) return { traditionalPct: 0, rothPct: 0 };

    return {
      traditionalPct: Math.round((traditional / total) * 100),
      rothPct: Math.round((roth / total) * 100),
    };
  });

  // This deterministic projectionData might be kept for comparison or removed if MC is primary
  let projectionData = $derived(() => {
    const currentData = propsDataSignal();
    const currentInputs = inputsSignal();
    // console.log(
    //   "TSPBalanceChart: Deriving deterministic projectionData. currentData?.tspBalance:",
    //   currentData?.tspBalance
    // );
    if (!currentData || !currentData.tspBalance) {
      // console.log(
      //   "TSPBalanceChart: deterministic projectionData returning [] because currentData or currentData.tspBalance is falsy."
      // );
      return [];
    }

    const startAge = currentData.retirementAge || 62;
    const startBalance = currentData.tspBalance;
    const withdrawalRate = 0.04;
    const growthRate = 0.06;
    const currentBreakdown = breakdown();

    let cumulativeWithdrawn = 0;
    // Use global calculationYears setting or default to 35
    const yearsToCalculate = currentInputs?.calculationYears || 35;

    const years = Array.from({ length: yearsToCalculate }, (_, i) => {
      const age = startAge + i;
      const yearsSinceRetirement = i;
      const netGrowthRate = growthRate - withdrawalRate;
      const remainingBalance =
        startBalance * Math.pow(1 + netGrowthRate, yearsSinceRetirement);
      const finalBalance = Math.max(0, remainingBalance);
      const currentWithdrawal =
        finalBalance > 0
          ? Math.min(
              finalBalance,
              startBalance *
                withdrawalRate *
                Math.pow(1 + netGrowthRate, yearsSinceRetirement - 1) // Simplified withdrawal calc
            )
          : 0;
      cumulativeWithdrawn += currentWithdrawal;

      return {
        age,
        balance: finalBalance,
        traditional: Math.max(
          0,
          finalBalance * (currentBreakdown.traditionalPct / 100)
        ),
        roth: Math.max(0, finalBalance * (currentBreakdown.rothPct / 100)),
        withdrawn: currentWithdrawal,
        totalWithdrawn: cumulativeWithdrawn,
      };
    });

    // console.log(
    //   "TSPBalanceChart: deterministic projectionData successfully derived:",
    //   years
    // );
    return years;
  });

  // Create chart configuration for Monte Carlo Fan Chart
  let chartConfig = $derived(() => {
    const mcProjection = monteCarloProjectionData();
    const isMonthly = $displayMode === "monthly";

    // console.log(
    //   "TSPBalanceChart: Deriving chartConfig. mcProjection length:",
    //   mcProjection?.length
    // );
    if (!mcProjection || mcProjection.length === 0) {
      // console.log(
      //   "TSPBalanceChart: chartConfig returning null because mcProjection is empty."
      // );
      return null;
    }

    const ages = mcProjection.map((p) => `Age ${p.age}`);
    const config = {
      type: "line" as const,
      data: {
        labels: ages,
        datasets: [
          // P10-P90 Band (order 0)
          {
            // Dataset 0: P90 (Upper line)
            label: "10th-90th Percentile",
            data: mcProjection.map((p) =>
              isMonthly ? convertToMonthly(p.p90) : p.p90
            ),
            borderColor: "transparent",
            backgroundColor: `${chartColors.tsp}33`, // Lighter color for wider band
            fill: "1", // Fill to dataset at index 1 (P10)
            order: 0,
            pointRadius: 0,
            borderWidth: 0,
            yAxisID: "y",
          },
          {
            // Dataset 1: P10 (Lower line - target for fill)
            label: "P10 (helper)",
            data: mcProjection.map((p) =>
              isMonthly ? convertToMonthly(p.p10) : p.p10
            ),
            borderColor: "transparent",
            pointRadius: 0,
            fill: false, // This line itself is not filled
            showInLegend: false,
            borderWidth: 0,
            yAxisID: "y",
          },
          // P25-P75 Band (order 1)
          {
            // Dataset 2: P75 (Upper line)
            label: "25th-75th Percentile",
            data: mcProjection.map((p) =>
              isMonthly ? convertToMonthly(p.p75) : p.p75
            ),
            borderColor: "transparent",
            backgroundColor: `${chartColors.pension}4D`, // Darker color for narrower band
            fill: "3", // Fill to dataset at index 3 (P25)
            order: 1,
            pointRadius: 0,
            borderWidth: 0,
            yAxisID: "y",
          },
          {
            // Dataset 3: P25 (Lower line - target for fill)
            label: "P25 (helper)",
            data: mcProjection.map((p) =>
              isMonthly ? convertToMonthly(p.p25) : p.p25
            ),
            borderColor: "transparent",
            pointRadius: 0,
            fill: false,
            showInLegend: false,
            borderWidth: 0,
            yAxisID: "y",
          },
          // Median Line (order 2)
          {
            // Dataset 4: P50
            label: "Median (50th Percentile)",
            data: mcProjection.map((p) =>
              isMonthly ? convertToMonthly(p.p50) : p.p50
            ),
            borderColor: chartColors.total,
            borderWidth: 2.5, // Prominent line
            fill: false,
            order: 2,
            pointRadius: 0, // No points on the line for cleaner look
            yAxisID: "y",
          },
        ],
      },
      options: {
        ...defaultChartOptions,
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          ...defaultChartOptions.plugins,
          legend: {
            display: true,
            labels: {
              filter: function (legendItem: any, chartData: any) {
                const dataset = chartData.datasets[legendItem.datasetIndex];
                return dataset.showInLegend !== false;
              },
            },
          },
          title: {
            display: true,
            text: `TSP Balance Projections (${$displayMode === "monthly" ? "Monthly" : "Annual"} Withdrawals)`,
          },
        },
        interaction: {
          mode: "index" as const,
          intersect: false,
        },
        scales: {
          y: {
            type: "linear" as const,
            display: true,
            position: "left" as const,
            title: {
              display: true,
              text: "Projected Balance ($)",
            },
            ticks: {
              callback: function (value: any) {
                return isMonthly
                  ? formatCurrency(convertToMonthly(value))
                  : formatCurrency(value);
              },
            },
          },
        },
      },
    };

    // console.log(
    //   "TSPBalanceChart: Monte Carlo chartConfig successfully derived."
    // );
    return config;
    // console.log(
    //   "TSPBalanceChart: Monte Carlo chartConfig successfully derived."
    // );
    return config;
  });

  onMount(() => {
    const currentChartConfig = chartConfig();
    if (canvasElement && currentChartConfig) {
      chartInstance = createChartInstance(canvasElement, currentChartConfig);
    }
  });
  $effect(() => {
    // Explicitly depend on displayMode to trigger reactivity
    const currentDisplayMode = $displayMode;
    console.log(
      `TSPBalanceChart: Display mode changed to ${currentDisplayMode}`
    );

    const currentChartConfig = chartConfig();
    if (chartInstance && currentChartConfig) {
      if (canvasElement) {
        chartInstance.destroy();
        chartInstance = createChartInstance(canvasElement, currentChartConfig);
      } else {
        chartInstance.destroy();
        chartInstance = null;
      }
    } else if (canvasElement && currentChartConfig && !chartInstance) {
      chartInstance = createChartInstance(canvasElement, currentChartConfig);
    }
  });

  onDestroy(() => {
    if (chartInstance) {
      chartInstance.destroy();
      chartInstance = null;
    }
  });
</script>

<div class="mb-8">
  <h3 class="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-4">
    Risk Analysis & Monte Carlo Projections
  </h3>

  {#if chartConfig()}
    <div
      class="h-96 bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-600 p-4"
    >
      <canvas bind:this={canvasElement}></canvas>
    </div>

    <!-- Key Insights for Monte Carlo -->
    {@const currentProps = propsDataSignal()}
    {@const mcSimResults = currentProps?.monteCarloResults}
    {@const mcInputParams = currentProps?.monteCarloInput}
    <!-- Assuming input params are passed via props -->
    {@const mcProjection = monteCarloProjectionData()}
    {#if mcSimResults && mcProjection && mcProjection.length > 0}
      {@const medianEndingBalance = mcProjection[mcProjection.length - 1]?.p50}
      {@const lastAge = mcProjection[mcProjection.length - 1]?.age}

      {@const medianDepletionYearInfo = (() => {
        const firstDepletedYear = mcProjection.find((p) => p.p50 < 1000);
        if (firstDepletedYear)
          return { age: firstDepletedYear.age, depleted: true };
        return { age: lastAge, depleted: false };
      })()}

      <div class="mt-4 grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
        <div class="bg-blue-50 dark:bg-blue-900/20 p-3 rounded-lg">
          <div class="font-medium text-blue-800 dark:text-blue-200">
            10-Year Outlook (Median)
          </div>
          {#if mcProjection.length > 9 && mcProjection[9]}
            <div class="text-blue-600 dark:text-blue-300">
              {$displayMode === "monthly"
                ? formatCurrency(convertToMonthly(mcProjection[9].p50)) + "/mo"
                : formatCurrency(mcProjection[9].p50) + "/yr"} (Age {mcProjection[9]
                .age})
            </div>
            <div class="text-xs text-blue-500 dark:text-blue-400 mt-1">
              Range (10th-90th): {$displayMode === "monthly"
                ? formatCurrency(convertToMonthly(mcProjection[9].p10))
                : formatCurrency(mcProjection[9].p10)} -
              {$displayMode === "monthly"
                ? formatCurrency(convertToMonthly(mcProjection[9].p90))
                : formatCurrency(mcProjection[9].p90)}
              {$displayMode === "monthly" ? "/mo" : "/yr"}
            </div>
          {:else}
            <div class="text-xs text-blue-500 dark:text-blue-400 mt-1">
              Data unavailable (less than 10 years).
            </div>
          {/if}
        </div>
        <div class="bg-green-50 dark:bg-green-900/20 p-3 rounded-lg">
          <div class="font-medium text-green-800 dark:text-green-200">
            20-Year Outlook (Median)
          </div>
          {#if mcProjection.length > 19 && mcProjection[19]}
            <div class="text-green-600 dark:text-green-300">
              {$displayMode === "monthly"
                ? formatCurrency(convertToMonthly(mcProjection[19].p50)) + "/mo"
                : formatCurrency(mcProjection[19].p50) + "/yr"} (Age {mcProjection[19]
                .age})
            </div>
            <div class="text-xs text-green-500 dark:text-green-400 mt-1">
              Range (10th-90th): {$displayMode === "monthly"
                ? formatCurrency(convertToMonthly(mcProjection[19].p10))
                : formatCurrency(mcProjection[19].p10)} -
              {$displayMode === "monthly"
                ? formatCurrency(convertToMonthly(mcProjection[19].p90))
                : formatCurrency(mcProjection[19].p90)}
              {$displayMode === "monthly" ? "/mo" : "/yr"}
            </div>
          {:else}
            <div class="text-xs text-green-500 dark:text-green-400 mt-1">
              Data unavailable (less than 20 years).
            </div>
          {/if}
        </div>
        <div class="bg-red-50 dark:bg-red-900/20 p-3 rounded-lg">
          <div class="font-medium text-red-800 dark:text-red-200">
            Depletion Risk (Median)
          </div>
          <div class="text-red-600 dark:text-red-300">
            {#if medianDepletionYearInfo.depleted}
              Potentially around Age {medianDepletionYearInfo.age}
            {:else}
              Very Low Risk (Median balance remains above $1k at Age {lastAge})
            {/if}
          </div>
          {#if mcSimResults.SuccessRate != null}
            <div class="text-xs text-red-500 dark:text-red-400 mt-1">
              Overall Success Rate: {(mcSimResults.SuccessRate * 100).toFixed(
                0
              )}%
            </div>
          {/if}
        </div>
      </div>

      <div class="mt-6 bg-gray-50 dark:bg-gray-700/30 p-4 rounded-lg">
        <h4 class="text-md font-semibold text-gray-700 dark:text-gray-300 mb-2">
          Monte Carlo Assumptions
        </h4>
        <ul
          class="text-xs text-gray-600 dark:text-gray-400 list-disc list-inside space-y-1"
        >
          {#if mcInputParams?.NumSimulations}
            <li>
              Simulations: {mcInputParams.NumSimulations.toLocaleString()}
            </li>
          {/if}
          <li>
            Expected Return: {mcInputParams?.ExpectedReturn != null
              ? (mcInputParams.ExpectedReturn * 100).toFixed(1) + "%"
              : "N/A"} annually
          </li>
          <li>
            Volatility (Std Dev): {mcInputParams?.ReturnStdDev != null
              ? (mcInputParams.ReturnStdDev * 100).toFixed(1) + "%"
              : "N/A"}
          </li>
          <li>
            Withdrawal Strategy: {mcInputParams?.AnnualWithdrawal != null
              ? ($displayMode === "monthly"
                  ? formatCurrency(
                      convertToMonthly(mcInputParams.AnnualWithdrawal)
                    ) + "/mo"
                  : formatCurrency(mcInputParams.AnnualWithdrawal) + "/yr") +
                " (inflation adj.)"
              : "N/A"}
          </li>
          <li>
            Years Simulated: {mcInputParams?.Years != null
              ? mcInputParams.Years
              : "N/A"}
          </li>
          <li>Rebalancing: Annual (assumed)</li>
        </ul>
      </div>
    {/if}
  {:else}
    <div
      class="h-96 bg-gray-50 dark:bg-gray-700 rounded-lg border-2 border-dashed border-gray-300 dark:border-gray-600 flex flex-col items-center justify-center p-4"
    >
      <div class="text-center text-gray-500 dark:text-gray-400">
        <div class="text-lg font-medium mb-2">
          No Monte Carlo Data Available
        </div>
        <div class="text-sm">
          Ensure Monte Carlo simulation is enabled and input data is available.
        </div>
      </div>
    </div>
  {/if}
</div>
