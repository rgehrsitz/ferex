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

  interface ProjectionYear {
    age: number;
    year: number;
    pension: number;
    socialSecurity: number;
    tsp: number;
    total: number;
  }
  let { data: dataSignal, inputs: inputsSignal } = $props<{ data: any; inputs: ScenarioInput | undefined }>();

  $effect(() => {
    const currentSignalValue = dataSignal();
    // console.log(
    //   "IncomeOverTimeChart: propsDataSignal (dataSignal) received:",
    //   currentSignalValue
    // );
    // if (currentSignalValue) {
    //   console.log(
    //     "IncomeOverTimeChart: currentSignalValue.totalIncome:",
    //     currentSignalValue.totalIncome
    //   );
    // }
  });

  let canvasElement = $state<HTMLCanvasElement>();
  let chartInstance: Chart | null = null;

  // Enhanced projection data for visualization with COLA adjustments
  let projectionData = $derived(() => {
    const currentData = dataSignal();
    const currentInputs = inputsSignal();
    // console.log(
    //   "IncomeOverTimeChart: Deriving projectionData. currentData:",
    //   currentData
    // );
    if (!currentData || !currentData.totalIncome) {
      // console.log(
      //   "IncomeOverTimeChart: projectionData returning [] because currentData or currentData.totalIncome is falsy."
      // );
      return [];
    }

    const startAge = currentData.retirementAge || 62;
    const colaRate = 0.025; // 2.5% annual COLA
    const basePension = currentData.pensionAnnual || 0;
    const baseSocialSecurity = currentData.socialSecurityAnnual || 0;
    const initialTSPBalance = currentData.tspBalance || 0;
    const withdrawalRate = 0.04;
    
    // Use global calculationYears setting or default to 30
    const yearsToCalculate = currentInputs?.calculationYears || 30;

    const years = Array.from({ length: yearsToCalculate }, (_, i) => {
      const age = startAge + i;
      const year = new Date().getFullYear() + i;

      // Apply COLA to pension and SS (typically starts at 62 for SS)
      const pensionWithCOLA = basePension * Math.pow(1 + colaRate, i);
      const ssWithCOLA =
        age >= 62 && baseSocialSecurity > 0
          ? baseSocialSecurity *
            Math.pow(1 + colaRate, Math.max(0, i - Math.max(0, 62 - startAge)))
          : 0;

      // TSP balance decreases with withdrawals but may have some growth
      const remainingTSPBalance = Math.max(
        0,
        initialTSPBalance * Math.pow(1.06 - withdrawalRate, i)
      ); // 6% growth minus 4% withdrawal
      const tspWithdrawal =
        remainingTSPBalance > 0 ? remainingTSPBalance * withdrawalRate : 0;

      return {
        age,
        year,
        pension: pensionWithCOLA,
        socialSecurity: ssWithCOLA,
        tsp: tspWithdrawal,
        total: pensionWithCOLA + ssWithCOLA + tspWithdrawal,
      };
    });

    // console.log(
    //   "IncomeOverTimeChart: projectionData successfully derived:",
    //   years
    // );
    return years;
  }); // Create chart configuration
  let chartConfig = $derived(() => {
    const currentProjectionData = projectionData(); // Call derived signal
    const isMonthly = $displayMode === "monthly";

    // Transform data for monthly display if needed
    const transformedData = isMonthly
      ? currentProjectionData.map((year) => ({
          ...year,
          pension: convertToMonthly(year.pension),
          socialSecurity: convertToMonthly(year.socialSecurity),
          tsp: convertToMonthly(year.tsp),
          total: convertToMonthly(year.total),
        }))
      : currentProjectionData;

    // console.log(
    //   "IncomeOverTimeChart: Deriving chartConfig. currentProjectionData (length: " +
    //     (currentProjectionData?.length ?? "undefined") +
    //     "):",
    //   currentProjectionData
    // );

    // Log the imported utilities to check their values
    // console.log(
    //   "IncomeOverTimeChart: Checking imported chartUtils: chartColors:",
    //   chartColors
    // );
    // console.log(
    //   "IncomeOverTimeChart: Checking imported chartUtils: defaultChartOptions:",
    //   defaultChartOptions
    // );
    // if (defaultChartOptions) {
    //   console.log(
    //     "IncomeOverTimeChart: Checking imported chartUtils: defaultChartOptions.plugins:",
    //     defaultChartOptions.plugins
    //   );
    // } else {
    //   console.log(
    //     "IncomeOverTimeChart: defaultChartOptions is null or undefined."
    //   );
    // }

    if (!currentProjectionData || currentProjectionData.length === 0) {
      // console.log(
      //   "IncomeOverTimeChart: chartConfig returning null because currentProjectionData is empty or null."
      // );
      return null;
    }
    // console.log(
    //   "IncomeOverTimeChart: Passed currentProjectionData length check."
    // );

    let ages;
    try {
      ages = currentProjectionData.map((p: any) => `Age ${p.age}`);
      // console.log(
      //   "IncomeOverTimeChart: 'ages' array created successfully:",
      //   ages
      // );
    } catch (e) {
      // console.error("IncomeOverTimeChart: ERROR creating 'ages' array:", e);
      return null;
    }

    try {
      const config = {
        type: "line" as const,
        data: {
          labels: ages,
          datasets: [
            {
              label: "Pension",
              data: transformedData.map((p: any) => p.pension),
              backgroundColor: `${chartColors.pension}33`,
              borderColor: chartColors.pension,
              borderWidth: 2,
              fill: "+1" as const, // Ensure 'fill' is a valid Chart.js type if specific strings are expected
            },
            {
              label: "Social Security",
              data: transformedData.map((p: any) => p.socialSecurity),
              backgroundColor: `${chartColors.socialSecurity}33`,
              borderColor: chartColors.socialSecurity,
              borderWidth: 2,
              fill: "+1" as const,
            },
            {
              label: "TSP Withdrawal",
              data: transformedData.map((p: any) => p.tsp),
              backgroundColor: `${chartColors.tsp}33`,
              borderColor: chartColors.tsp,
              borderWidth: 2,
              fill: "origin" as const,
            },
          ],
        },
        options: {
          ...defaultChartOptions,
          plugins: {
            ...(defaultChartOptions?.plugins || {}), // Guard against undefined plugins            displayMode: $displayMode, // Add the display mode to the chart options
            title: {
              display: true,
              text: `Income Streams Over Time (${$displayMode === "monthly" ? "Monthly" : "Annual"} with COLA adjustments)`,
            },
          },
          interaction: {
            mode: "index" as const,
            intersect: false,
          },
          elements: {
            point: {
              radius: 3,
              hoverRadius: 6,
            },
          },
        },
      };
      // console.log(
      //   "IncomeOverTimeChart: chartConfig object CREATED successfully:",
      //   config
      // );
      return config;
    } catch (error) {
      // console.error(
      //   "IncomeOverTimeChart: ERROR during chart config object creation:",
      //   error
      // );
      return null; // Explicitly return null on error
    }
  });

  // Create/update chart when data changes
  onMount(() => {
    if (canvasElement && chartConfig()) {
      chartInstance = createChartInstance(canvasElement, chartConfig()!);
    }
  });
  // Update chart when config changes or display mode changes
  $effect(() => {
    // Explicitly depend on displayMode to trigger reactivity
    const currentDisplayMode = $displayMode;
    console.log(
      `IncomeOverTimeChart: Display mode changed to ${currentDisplayMode}`
    );

    if (chartInstance && chartConfig()) {
      chartInstance.data = chartConfig()!.data;
      chartInstance.options = chartConfig()!.options;
      chartInstance.update();
    } else if (canvasElement && chartConfig() && !chartInstance) {
      chartInstance = createChartInstance(canvasElement, chartConfig()!);
    }
  });

  // Cleanup on destroy
  onDestroy(() => {
    if (chartInstance) {
      chartInstance.destroy();
      chartInstance = null;
    }
  });

  // Extract the actual array for template use
  let projections = $derived(projectionData as unknown as ProjectionYear[]);
</script>

<div class="mb-8">
  <h3 class="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-4">
    Income Streams Over Time
  </h3>

  {#if chartConfig()}
    <div
      class="h-80 bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-600 p-4"
    >
      <canvas bind:this={canvasElement}></canvas>
    </div>

    <!-- Key Insights -->
    {#if projections && projections.length > 0}
      {@const firstYear = projections[0]}
      {@const tenthYear = projections[9] || firstYear}
      {@const lastYear = projections[projections.length - 1]}

      <div class="mt-4 grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
        <div class="bg-blue-50 dark:bg-blue-900/20 p-3 rounded-lg">
          <div class="font-medium text-blue-800 dark:text-blue-200">
            At Retirement (Age {firstYear.age})
          </div>
          <div class="text-blue-600 dark:text-blue-300">
            {$displayMode === "monthly"
              ? formatCurrency(convertToMonthly(firstYear.total)) + "/mo"
              : formatCurrency(firstYear.total) + "/yr"} total
          </div>
        </div>
        <div class="bg-green-50 dark:bg-green-900/20 p-3 rounded-lg">
          <div class="font-medium text-green-800 dark:text-green-200">
            10 Years Later (Age {tenthYear.age})
          </div>
          <div class="text-green-600 dark:text-green-300">
            {$displayMode === "monthly"
              ? formatCurrency(convertToMonthly(tenthYear.total)) + "/mo"
              : formatCurrency(tenthYear.total) + "/yr"} total
          </div>
        </div>
        <div class="bg-purple-50 dark:bg-purple-900/20 p-3 rounded-lg">
          <div class="font-medium text-purple-800 dark:text-purple-200">
            At Age {lastYear.age}
          </div>
          <div class="text-purple-600 dark:text-purple-300">
            {$displayMode === "monthly"
              ? formatCurrency(convertToMonthly(lastYear.total)) + "/mo"
              : formatCurrency(lastYear.total) + "/yr"} total
          </div>
        </div>
      </div>
    {/if}
  {:else}
    <div
      class="h-64 bg-gray-50 dark:bg-gray-700 rounded-lg border-2 border-dashed border-gray-300 dark:border-gray-600 flex flex-col items-center justify-center p-4"
    >
      <div class="text-center text-gray-500 dark:text-gray-400">
        <div class="text-lg font-medium mb-2">No Chart Data Available</div>
        <div class="text-sm">
          Complete scenario inputs to see income projections
        </div>
      </div>
    </div>
  {/if}
</div>
