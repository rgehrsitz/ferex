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
  import type { ScenarioInput } from "../../../types";

  interface MonteCarloResult {
    age: number;
    percentile10: number;
    percentile25: number;
    percentile50: number;
    percentile75: number;
    percentile90: number;
    failureRate: number;
  }
  let { data: dataSignal, inputs: inputsSignal } = $props<{ data: any; inputs: ScenarioInput | undefined }>();

  let canvasElement = $state<HTMLCanvasElement>();
  let chartInstance: Chart | null = null;
  // Generate Monte Carlo simulation data
  let monteCarloData = $derived(() => {
    if (!dataSignal() || !dataSignal().tspBalance) return [];

    const currentInputs = inputsSignal();
    const startAge = dataSignal().retirementAge || 62;
    const initialBalance = dataSignal().tspBalance;
    const withdrawalRate = 0.04;
    const meanReturn = 0.07;
    const volatility = 0.15;
    const isMonthly = $displayMode === "monthly";
    
    // Use global calculationYears setting or default to 30
    const yearsToCalculate = currentInputs?.calculationYears || 30;

    // Simulate portfolio outcomes with different return scenarios
    return Array.from({ length: yearsToCalculate }, (_, i) => {
      const age = startAge + i;
      const years = i + 1;

      // Simplified Monte Carlo percentiles
      const worstCase =
        initialBalance *
        Math.pow(1 + meanReturn - 2 * volatility - withdrawalRate, years);
      const poorCase =
        initialBalance *
        Math.pow(1 + meanReturn - volatility - withdrawalRate, years);
      const medianCase =
        initialBalance * Math.pow(1 + meanReturn - withdrawalRate, years);
      const goodCase =
        initialBalance *
        Math.pow(1 + meanReturn + volatility - withdrawalRate, years);
      const bestCase =
        initialBalance *
        Math.pow(1 + meanReturn + 2 * volatility - withdrawalRate, years);

      // Failure rate increases as balance depletes
      const failureRate = Math.max(
        0,
        Math.min(1, ((yearsToCalculate - years) / yearsToCalculate) * (worstCase < 0 ? 1 : 0))
      );

      return {
        age,
        percentile10: Math.max(0, worstCase),
        percentile25: Math.max(0, poorCase),
        percentile50: Math.max(0, medianCase),
        percentile75: Math.max(0, goodCase),
        percentile90: Math.max(0, bestCase),
        failureRate,
      };
    });
  });
  // Create fan chart configuration showing percentile bands
  let chartConfig = $derived(() => {
    const simData = monteCarloData();
    if (!simData || simData.length === 0) return null;

    const ages = simData.map((p: MonteCarloResult) => `Age ${p.age}`);

    return {
      type: "line",
      data: {
        labels: ages,
        datasets: [
          {
            label: "10th-90th Percentile",
            data: simData.map((p: MonteCarloResult) => p.percentile90),
            backgroundColor: `${chartColors.background}20`,
            borderColor: `${chartColors.border}40`,
            borderWidth: 1,
            fill: "+3",
            pointRadius: 0,
          },
          {
            label: "25th-75th Percentile",
            data: simData.map((p: MonteCarloResult) => p.percentile75),
            backgroundColor: `${chartColors.background}40`,
            borderColor: `${chartColors.border}60`,
            borderWidth: 1,
            fill: "+1",
            pointRadius: 0,
          },
          {
            label: "Median (50th)",
            data: simData.map((p: MonteCarloResult) => p.percentile50),
            backgroundColor: chartColors.pension,
            borderColor: chartColors.pension,
            borderWidth: 3,
            fill: false,
            pointRadius: 2,
          },
          {
            label: "25th Percentile",
            data: simData.map((p: MonteCarloResult) => p.percentile25),
            backgroundColor: "transparent",
            borderColor: "transparent",
            borderWidth: 0,
            fill: false,
            pointRadius: 0,
          },
          {
            label: "10th Percentile",
            data: simData.map((p: MonteCarloResult) => p.percentile10),
            backgroundColor: "transparent",
            borderColor: "transparent",
            borderWidth: 0,
            fill: false,
            pointRadius: 0,
          },
        ],
      },
      options: {
        ...defaultChartOptions,
        plugins: {
          ...defaultChartOptions.plugins,
          displayMode: $displayMode, // Add display mode to chart options
          title: {
            display: true,
            text: `Portfolio Balance Projections (${$displayMode === "monthly" ? "Monthly" : "Annual"} Monte Carlo Simulation)`,
          },
          legend: {
            display: true,
            position: "bottom" as const,
            labels: {
              filter: (legendItem: any) => {
                // Only show meaningful legend items
                return (
                  !legendItem.text.includes("Percentile") ||
                  legendItem.text.includes("Median")
                );
              },
            },
          },
        },
        interaction: {
          mode: "index" as const,
          intersect: false,
        },
        elements: {
          point: {
            radius: 1,
            hoverRadius: 4,
          },
        },
        scales: {
          y: {
            beginAtZero: true,
            ticks: {
              callback: function (value: any) {
                const isMonthly = $displayMode === "monthly";
                return isMonthly
                  ? formatCurrency(convertToMonthly(value))
                  : formatCurrency(value);
              },
            },
          },
        },
      },
    };
  });
  // Create/update chart when data changes
  onMount(() => {
    const config = chartConfig();
    if (canvasElement && config) {
      chartInstance = createChartInstance(canvasElement, config);
    }
  }); // Update chart when config changes or display mode changes
  $effect(() => {
    // Explicitly depend on displayMode to trigger reactivity
    const currentDisplayMode = $displayMode;
    // console.log(`MonteCarloChart: Display mode changed to ${currentDisplayMode}`);

    const config = chartConfig();
    if (chartInstance && config) {
      chartInstance.data = config.data;
      chartInstance.options = config.options;
      chartInstance.update();
    } else if (canvasElement && config && !chartInstance) {
      chartInstance = createChartInstance(canvasElement, config);
    }
  });

  // Cleanup on destroy
  onDestroy(() => {
    if (chartInstance) {
      chartInstance.destroy();
      chartInstance = null;
    }
  });

  // Extract data for template use
  let projections = $derived(() => monteCarloData() as MonteCarloResult[]);
</script>

<div class="mb-8">
  <h3 class="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-4">
    Risk Analysis & Monte Carlo Projections
  </h3>
  {#if chartConfig()}
    <div
      class="h-80 bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-600 p-4"
    >
      <canvas bind:this={canvasElement}></canvas>
    </div>

    <!-- Risk Metrics -->
    {#if projections() && projections().length > 0}
      {@const tenYearProjection = projections()[9] || projections()[0]}
      {@const twentyYearProjection =
        projections()[19] || projections()[projections().length - 1]}
      {@const depletionRisk = projections().find((p) => p.percentile25 < 10000)}

      <div class="mt-4 grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
        <div class="bg-blue-50 dark:bg-blue-900/20 p-3 rounded-lg">
          <div class="font-medium text-blue-800 dark:text-blue-200">
            10-Year Outlook
          </div>
          <div class="text-blue-600 dark:text-blue-300 space-y-1">
            <div>
              Median: {$displayMode === "monthly"
                ? formatCurrency(
                    convertToMonthly(tenYearProjection.percentile50)
                  ) + "/mo"
                : formatCurrency(tenYearProjection.percentile50) + "/yr"}
            </div>
            <div class="text-xs">
              Range: {$displayMode === "monthly"
                ? formatCurrency(
                    convertToMonthly(tenYearProjection.percentile25)
                  )
                : formatCurrency(tenYearProjection.percentile25)} -
              {$displayMode === "monthly"
                ? formatCurrency(
                    convertToMonthly(tenYearProjection.percentile75)
                  )
                : formatCurrency(tenYearProjection.percentile75)}
              {$displayMode === "monthly" ? "/mo" : "/yr"}
            </div>
          </div>
        </div>
        <div class="bg-yellow-50 dark:bg-yellow-900/20 p-3 rounded-lg">
          <div class="font-medium text-yellow-800 dark:text-yellow-200">
            20-Year Outlook
          </div>
          <div class="text-yellow-600 dark:text-yellow-300 space-y-1">
            <div>
              Median: {$displayMode === "monthly"
                ? formatCurrency(
                    convertToMonthly(twentyYearProjection.percentile50)
                  ) + "/mo"
                : formatCurrency(twentyYearProjection.percentile50) + "/yr"}
            </div>
            <div class="text-xs">
              Range: {$displayMode === "monthly"
                ? formatCurrency(
                    convertToMonthly(twentyYearProjection.percentile25)
                  )
                : formatCurrency(twentyYearProjection.percentile25)} -
              {$displayMode === "monthly"
                ? formatCurrency(
                    convertToMonthly(twentyYearProjection.percentile75)
                  )
                : formatCurrency(twentyYearProjection.percentile75)}
              {$displayMode === "monthly" ? "/mo" : "/yr"}
            </div>
          </div>
        </div>
        <div class="bg-red-50 dark:bg-red-900/20 p-3 rounded-lg">
          <div class="font-medium text-red-800 dark:text-red-200">
            Depletion Risk
          </div>
          <div class="text-red-600 dark:text-red-300">
            {#if depletionRisk}
              <div>Possible by age {depletionRisk.age}</div>
              <div class="text-xs">25th percentile scenario</div>
            {:else}
              <div>Very Low Risk</div>
              <div class="text-xs">Beyond 30-year horizon</div>
            {/if}
          </div>
        </div>
      </div>

      <!-- Assumption Details -->
      <div class="mt-4 p-3 bg-gray-50 dark:bg-gray-700 rounded-lg">
        <div class="text-xs text-gray-600 dark:text-gray-300">
          <div class="font-medium mb-1">Monte Carlo Assumptions:</div>
          <div class="space-y-1">
            <div>• Expected Return: 7% annually</div>
            <div>• Volatility: 15% (standard deviation)</div>
            <div>
              • Withdrawal Rate: 4% annually ({$displayMode === "monthly"
                ? "~0.33% monthly"
                : "4% annually"})
            </div>
            <div>• Rebalancing: Annual</div>
          </div>
        </div>
      </div>
    {/if}
  {:else}
    <div
      class="h-64 bg-gray-50 dark:bg-gray-700 rounded-lg border-2 border-dashed border-gray-300 dark:border-gray-600 flex flex-col items-center justify-center p-4"
    >
      <div class="text-center text-gray-500 dark:text-gray-400">
        <div class="text-lg font-medium mb-2">No Portfolio Data Available</div>
        <div class="text-sm">Complete TSP inputs to see risk analysis</div>
      </div>
    </div>
  {/if}
</div>
