<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import {
    Chart,
    chartColors,
    formatCurrency,
    createChartInstance,
    convertToMonthly,
  } from "./chartUtils";
  import { displayMode } from "../../stores/displayPreferences";

  interface IncomeSource {
    name: string;
    value: number;
    percentage: number;
    color: string;
  }

  let { data: dataSignal } = $props<{ data: any }>(); // Renamed to dataSignal

  let canvasElement = $state<HTMLCanvasElement>();
  let chartInstance: Chart | null = null;
  // Calculate income source breakdown
  let incomeBreakdown = $derived(() => {
    const currentData = dataSignal(); // Call the signal
    if (!currentData) return [];

    const isMonthly = $displayMode === "monthly";

    // Base annual values
    let pension = currentData.pensionAnnual || 0;
    let socialSecurity = currentData.socialSecurityAnnual || 0;
    let tspWithdrawal = (currentData.tspBalance || 0) * 0.04; // 4% withdrawal

    // Convert to monthly if needed
    if (isMonthly) {
      pension = convertToMonthly(pension);
      socialSecurity = convertToMonthly(socialSecurity);
      tspWithdrawal = convertToMonthly(tspWithdrawal);
    }

    const total = pension + socialSecurity + tspWithdrawal;

    if (total === 0) return [];
    const sources: IncomeSource[] = [
      {
        name: "Federal Pension",
        value: pension,
        percentage: Math.round((pension / total) * 100),
        color: chartColors.pension,
      },
      {
        name: "Social Security",
        value: socialSecurity,
        percentage: Math.round((socialSecurity / total) * 100),
        color: chartColors.socialSecurity,
      },
      {
        name: "TSP Withdrawal",
        value: tspWithdrawal,
        percentage: Math.round((tspWithdrawal / total) * 100),
        color: chartColors.tsp,
      },
    ].filter((source) => source.value > 0);

    return sources;
  });
  // Create pie chart configuration
  let chartConfig = $derived(() => {
    const breakdown = incomeBreakdown();
    if (!breakdown || breakdown.length === 0) return null;

    return {
      type: "doughnut",
      data: {
        labels: breakdown.map((s: IncomeSource) => s.name),
        datasets: [
          {
            data: breakdown.map((s: IncomeSource) => s.value),
            backgroundColor: breakdown.map((s: IncomeSource) => s.color),
            borderColor: breakdown.map((s: IncomeSource) => s.color),
            borderWidth: 2,
            hoverOffset: 4,
          },
        ],
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          displayMode: $displayMode, // Add display mode to chart options
          title: {
            display: true,
            text: `Income Sources at Retirement (${$displayMode === "monthly" ? "Monthly" : "Annual"})`,
          },
          legend: {
            display: true,
            position: "bottom" as const,
          },
          tooltip: {
            callbacks: {
              label: function (context: any) {
                const currentBreakdown = incomeBreakdown(); // Ensure correct derived value is used
                const source = currentBreakdown[context.dataIndex];
                const isMonthly = $displayMode === "monthly";
                const value = isMonthly
                  ? convertToMonthly(source.value)
                  : source.value;
                return `${source.name}: ${formatCurrency(value)}${isMonthly ? "/mo" : "/yr"} (${source.percentage}%)`;
              },
            },
          },
        },
        cutout: "40%", // Creates doughnut hole
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
    // console.log(`IncomeBreakdownChart: Display mode changed to ${currentDisplayMode}`);

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
</script>

<div>
  <h3 class="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-4">
    Income Sources at Retirement
  </h3>

  {#if chartConfig()}
    <div
      class="h-64 bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-600 p-4"
    >
      <canvas bind:this={canvasElement}></canvas>
    </div>

    <!-- Income Source Details -->
    {#if incomeBreakdown() && incomeBreakdown().length > 0}
      {@const sources = incomeBreakdown()}
      {@const total = sources.reduce(
        (sum: number, s: IncomeSource) => sum + s.value,
        0
      )}
      <div class="mt-4 space-y-2">
        {#each sources as source (source.name)}
          <div
            class="flex items-center justify-between p-2 bg-gray-50 dark:bg-gray-700 rounded"
          >
            <div class="flex items-center">
              <span
                class="w-3 h-3 rounded-full mr-2"
                style="background-color: {source.color};"
              ></span>
              <span class="text-sm text-gray-700 dark:text-gray-300"
                >{source.name}</span
              >
            </div>
            <div class="text-sm text-gray-900 dark:text-gray-100">
              {$displayMode === "monthly"
                ? formatCurrency(convertToMonthly(source.value)) + "/mo"
                : formatCurrency(source.value) + "/yr"} ({source.percentage}%)
            </div>
          </div>
        {/each}
        <div
          class="flex items-center justify-between p-2 bg-gray-100 dark:bg-gray-600 rounded mt-2"
        >
          <span class="text-sm font-semibold text-gray-800 dark:text-gray-200"
            >Total Estimated Income</span
          >
          <span class="text-sm font-semibold text-gray-900 dark:text-gray-100"
            >{$displayMode === "monthly"
              ? formatCurrency(convertToMonthly(total)) + "/mo"
              : formatCurrency(total) + "/yr"}</span
          >
        </div>
      </div>
    {/if}
  {:else}
    <div
      class="h-64 bg-gray-50 dark:bg-gray-700 rounded-lg border-2 border-dashed border-gray-300 dark:border-gray-600 flex flex-col items-center justify-center p-4"
    >
      <div class="text-center text-gray-500 dark:text-gray-400">
        <div class="text-lg font-medium mb-2">No Income Data Available</div>
        <div class="text-sm">
          Complete scenario calculations to see income breakdown
        </div>
      </div>
    </div>
  {/if}
</div>
