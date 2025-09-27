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

  interface TaxBreakdown {
    federalTax: number;
    stateTax: number;
    totalTax: number;
    netIncome: number;
    federalRate: number;
    stateRate: number;
    effectiveRate: number;
  }

  interface TaxProjectionYear {
    age: number;
    grossIncome: number;
    federalTax: number;
    stateTax: number;
    totalTax: number;
    netIncome: number;
    effectiveRate: number;
  }
  let { data: dataSignal, inputs: inputsSignal } = $props<{ data: any; inputs: ScenarioInput | undefined }>();

  let canvasElement = $state<HTMLCanvasElement | undefined>();
  let chartInstance: Chart | null = $state(null);

  // Format percentage for display
  function formatPercent(value: number | undefined): string {
    if (value == null || isNaN(value)) return "N/A";
    return `${(value * 100).toFixed(1)}%`;
  }

  // Calculate tax breakdown
  let taxBreakdown = $derived(() => {
    const currentData = dataSignal();
    if (!currentData || !currentData.totalIncome) {
      return {
        federalTax: 0,
        stateTax: 0,
        totalTax: 0,
        netIncome: 0,
        federalRate: 0,
        stateRate: 0,
        effectiveRate: 0,
      };
    }

    const grossIncome = currentData.totalIncome;
    const federalTax = currentData.federalTax || 0;
    const stateTax = currentData.stateTax || 0;
    const totalTax = federalTax + stateTax;
    const netIncome = grossIncome - totalTax;

    return {
      federalTax,
      stateTax,
      totalTax,
      netIncome,
      federalRate: grossIncome > 0 ? federalTax / grossIncome : 0,
      stateRate: grossIncome > 0 ? stateTax / grossIncome : 0,
      effectiveRate: grossIncome > 0 ? totalTax / grossIncome : 0,
    };
  });

  // Extract tax breakdown for template use
  let taxValues = $derived(() => taxBreakdown() as TaxBreakdown); // Call taxBreakdown as a function
  // Enhanced tax projections with income growth and tax changes
  let taxProjectionData = $derived(() => {
    const currentData = dataSignal(); // Call the signal
    const currentInputs = inputsSignal();
    const currentTaxValues = taxValues(); // Call derived signal
    const isMonthly = $displayMode === "monthly";

    if (!currentData || !currentData.totalIncome) return [];

    const startAge = currentData.retirementAge || 62;
    const baseIncome = currentData.totalIncome;
    const baseFederalRate = currentTaxValues.federalRate;
    const baseStateRate = currentTaxValues.stateRate;
    const colaRate = 0.025; // 2.5% annual income growth from COLA
    
    // Use global calculationYears setting or default to 20
    const yearsToCalculate = currentInputs?.calculationYears || 20;

    return Array.from({ length: yearsToCalculate }, (_, i) => {
      const age = startAge + i;
      const adjustedIncome = baseIncome * Math.pow(1 + colaRate, i);

      const bracketCreepFactor = 1 + i * 0.001;
      const federalTax = adjustedIncome * baseFederalRate * bracketCreepFactor;

      // Will transform later when building chart data
      const stateTax = adjustedIncome * baseStateRate;
      const totalTax = federalTax + stateTax;
      const netIncome = adjustedIncome - totalTax;

      return {
        age,
        grossIncome: adjustedIncome,
        federalTax,
        stateTax,
        totalTax,
        netIncome,
        effectiveRate: adjustedIncome > 0 ? totalTax / adjustedIncome : 0,
      };
    });
  });
  // Create chart configuration for stacked bar chart showing tax burden over time
  let chartConfig = $derived(() => {
    const projectionData = taxProjectionData(); // Call derived signal
    if (!projectionData || projectionData.length === 0) return null;

    const isMonthly = $displayMode === "monthly";

    // Transform data for monthly display if needed
    const transformedData = isMonthly
      ? projectionData.map((year) => ({
          ...year,
          grossIncome: convertToMonthly(year.grossIncome),
          federalTax: convertToMonthly(year.federalTax),
          stateTax: convertToMonthly(year.stateTax),
          totalTax: convertToMonthly(year.totalTax),
          netIncome: convertToMonthly(year.netIncome),
        }))
      : projectionData;

    const ages = transformedData.map((p: TaxProjectionYear) => `Age ${p.age}`);

    return {
      type: "bar" as const,
      data: {
        labels: ages,
        datasets: [
          {
            label: "Net Income",
            data: transformedData.map((p: TaxProjectionYear) => p.netIncome),
            backgroundColor: chartColors.socialSecurity,
            borderColor: chartColors.socialSecurity,
            borderWidth: 1,
          },
          {
            label: "Federal Tax",
            data: transformedData.map((p: TaxProjectionYear) => p.federalTax),
            backgroundColor: chartColors.federal,
            borderColor: chartColors.federal,
            borderWidth: 1,
          },
          {
            label: "State Tax",
            data: transformedData.map((p: TaxProjectionYear) => p.stateTax),
            backgroundColor: chartColors.state,
            borderColor: chartColors.state,
            borderWidth: 1,
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
            text: `Tax Burden Over Time (${$displayMode === "monthly" ? "Monthly" : "Annual"} Stacked View)`,
          },
        },
        responsive: true,
        scales: {
          x: {
            stacked: true,
          },
          y: {
            stacked: true,
            beginAtZero: true,
            title: {
              display: true,
              text: "Amount ($)",
            },
            ticks: {
              callback: function (value: any) {
                return formatCurrency(value);
              },
            },
          },
        },
      },
    };
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
    // console.log(`TaxBurdenChart: Display mode changed to ${currentDisplayMode}`);

    const currentChartConfig = chartConfig();
    if (chartInstance && currentChartConfig) {
      chartInstance.data = currentChartConfig.data;
      chartInstance.options = currentChartConfig.options;
      chartInstance.update();
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

<div>
  <h3 class="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-4">
    Tax Burden Over Time
  </h3>

  {#if chartConfig()}
    <div
      class="h-80 bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-600 p-4"
    >
      <canvas bind:this={canvasElement}></canvas>
    </div>

    <!-- Tax Breakdown Details -->
    {@const breakdown = taxBreakdown()}
    <!-- Call derived signal -->
    {#if breakdown && breakdown.total > 0}
      <div class="mt-4 grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
        <div class="bg-red-50 dark:bg-red-900/20 p-3 rounded-lg">
          <div class="font-medium text-red-800 dark:text-red-200">
            Federal Tax
          </div>
          <div class="text-red-600 dark:text-red-300">
            {$displayMode === "monthly"
              ? formatCurrency(convertToMonthly(breakdown.federalTax)) + "/mo"
              : formatCurrency(breakdown.federalTax) + "/yr"} ({formatPercent(
              breakdown.federalRate
            )})
          </div>
        </div>
        <div class="bg-amber-50 dark:bg-amber-900/20 p-3 rounded-lg">
          <div class="font-medium text-amber-800 dark:text-amber-200">
            State Tax
          </div>
          <div class="text-amber-600 dark:text-amber-300">
            {$displayMode === "monthly"
              ? formatCurrency(convertToMonthly(breakdown.stateTax)) + "/mo"
              : formatCurrency(breakdown.stateTax) + "/yr"} ({formatPercent(
              breakdown.stateRate
            )})
          </div>
        </div>
        <div class="bg-gray-100 dark:bg-gray-700 p-3 rounded-lg">
          <div class="font-medium text-gray-800 dark:text-gray-200">
            Total Tax Burden
          </div>
          <div class="text-gray-600 dark:text-gray-300">
            {$displayMode === "monthly"
              ? formatCurrency(convertToMonthly(breakdown.totalTax)) + "/mo"
              : formatCurrency(breakdown.totalTax) + "/yr"} ({formatPercent(
              breakdown.effectiveRate
            )})
          </div>
        </div>
      </div>
    {/if}
  {:else}
    <div
      class="h-64 bg-gray-50 dark:bg-gray-700 rounded-lg border-2 border-dashed border-gray-300 dark:border-gray-600 flex flex-col items-center justify-center p-4"
    >
      <div class="text-center text-gray-500 dark:text-gray-400">
        <div class="text-lg font-medium mb-2">No Tax Data Available</div>
        <div class="text-sm">
          Complete scenario calculations to see tax analysis
        </div>
      </div>
    </div>
  {/if}
</div>
