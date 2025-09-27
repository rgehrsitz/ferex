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

  interface ScenarioComparison {
    name: string;
    totalIncome: number;
    netIncome: number;
    effectiveTaxRate: number;
    tspBalance: number;
    depletionAge: number | null;
  }

  let { scenarios = [], activeScenarioId } = $props<{
    scenarios?: any[];
    activeScenarioId?: string;
  }>();

  let canvasElement = $state<HTMLCanvasElement>();
  let chartInstance: Chart | null = null;
  let chartType = $state<"income" | "tax" | "sustainability">("income");

  // Process scenarios for comparison
  let scenarioData = $derived(() => {
    const validScenarios = scenarios.filter(
      (s: any) => s.results && s.results.TotalRetirementIncome > 0
    );

    return validScenarios.map((scenario: any): ScenarioComparison => {
      const results = scenario.results;
      const isActive = scenario.id === activeScenarioId;

      // Calculate depletion age estimate
      const tspBalance =
        results.TSPResult?.totalProjectedBalanceAtRetirement || 0;
      const withdrawalRate = 0.04;
      const growthRate = 0.06;
      const netRate = growthRate - withdrawalRate;
      const depletionAge =
        tspBalance > 0 && netRate < 0
          ? Math.log(0.1) / Math.log(1 + netRate) +
            (scenario.inputs?.plannedRetirementDate
              ? Math.floor(
                  (new Date(scenario.inputs.plannedRetirementDate).getTime() -
                    new Date(scenario.inputs.dateOfBirth).getTime()) /
                    (365.25 * 24 * 60 * 60 * 1000)
                )
              : 65)
          : null;

      return {
        name: isActive ? `${scenario.name} (Current)` : scenario.name,
        totalIncome: results.TotalRetirementIncome || 0,
        netIncome: results.NetAfterTaxIncome || 0,
        effectiveTaxRate: results.EffectiveTaxRate || 0,
        tspBalance: tspBalance,
        depletionAge:
          depletionAge && depletionAge > 0 && depletionAge < 100
            ? Math.round(depletionAge)
            : null,
      };
    });
  });

  // Create comparison chart based on selected type
  let chartConfig = $derived(() => {
    const data = scenarioData();
    if (!data || data.length === 0) return null;

    const labels = data.map((s: any) => s.name);
    const isActive = (name: string) => name.includes("(Current)");

    const getColor = (name: string, alpha: string = "") => {
      return isActive(name)
        ? `${chartColors.pension}${alpha}`
        : `${chartColors.total}${alpha}`;
    };

    switch (chartType) {
      case "income":
        return {
          type: "bar",
          data: {
            labels,
            datasets: [
              {
                label: "Total Income",
                data: data.map((s: any) => s.totalIncome),
                backgroundColor: labels.map((name: any) =>
                  getColor(name, "60")
                ),
                borderColor: labels.map((name: any) => getColor(name)),
                borderWidth: 2,
              },
              {
                label: "Net After-Tax Income",
                data: data.map((s: any) => s.netIncome),
                backgroundColor: labels.map((name: any) =>
                  getColor(name, "30")
                ),
                borderColor: labels.map((name: any) => getColor(name)),
                borderWidth: 2,
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
                text: `Income Comparison (${$displayMode === "monthly" ? "Monthly" : "Annual"})`,
              },
              legend: {
                display: true,
                position: "bottom" as const,
              },
            },
            scales: {
              y: {
                beginAtZero: true,
                ticks: {
                  callback: function (value: any) {
                    return $displayMode === "monthly"
                      ? formatCurrency(convertToMonthly(value))
                      : formatCurrency(value);
                  },
                },
              },
            },
          },
        };

      case "tax":
        return {
          type: "bar",
          data: {
            labels,
            datasets: [
              {
                label: "Tax Burden",
                data: data.map((s: any) => s.totalIncome - s.netIncome),
                backgroundColor: labels.map((name: any) =>
                  getColor(name, "60")
                ),
                borderColor: labels.map((name: any) => getColor(name)),
                borderWidth: 2,
              },
              {
                label: "Effective Tax Rate (%)",
                data: data.map((s: any) => s.effectiveTaxRate * 100),
                type: "line" as const,
                backgroundColor: chartColors.federal,
                borderColor: chartColors.federal,
                borderWidth: 3,
                yAxisID: "y1",
                fill: false,
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
                text: `Tax Burden Comparison (${$displayMode === "monthly" ? "Monthly" : "Annual"})`,
              },
              legend: {
                display: true,
                position: "bottom" as const,
              },
            },
            scales: {
              y: {
                type: "linear" as const,
                display: true,
                position: "left" as const,
                beginAtZero: true,
                title: {
                  display: true,
                  text: "Tax Amount ($)",
                },
                ticks: {
                  callback: function (value: any) {
                    return $displayMode === "monthly"
                      ? formatCurrency(convertToMonthly(value))
                      : formatCurrency(value);
                  },
                },
              },
              y1: {
                type: "linear" as const,
                display: true,
                position: "right" as const,
                title: {
                  display: true,
                  text: "Effective Rate (%)",
                },
                grid: {
                  drawOnChartArea: false,
                },
                ticks: {
                  callback: function (value: any) {
                    return value.toFixed(1) + "%";
                  },
                },
              },
            },
          },
        };

      case "sustainability":
        return {
          type: "scatter",
          data: {
            datasets: data.map((scenario: any, index: any) => ({
              label: scenario.name,
              data: [
                {
                  x: scenario.tspBalance / 1000, // TSP Balance in thousands
                  y: scenario.depletionAge || 100, // Depletion age or 100 if never
                },
              ],
              backgroundColor: isActive(scenario.name)
                ? chartColors.pension
                : chartColors.total,
              borderColor: isActive(scenario.name)
                ? chartColors.pension
                : chartColors.total,
              pointRadius: 8,
              pointHoverRadius: 10,
            })),
          },
          options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
              displayMode: $displayMode, // Add display mode to chart options
              title: {
                display: true,
                text: `TSP Balance vs. Sustainability (${$displayMode === "monthly" ? "Monthly" : "Annual"})`,
              },
              legend: {
                display: true,
                position: "bottom" as const,
              },
              tooltip: {
                callbacks: {
                  label: function (context: any) {
                    const scenario = data[context.datasetIndex];
                    return [
                      `${scenario.name}`,
                      `TSP Balance: ${formatCurrency(scenario.tspBalance)}`,
                      `Estimated Depletion: ${scenario.depletionAge ? `Age ${scenario.depletionAge}` : "Beyond Age 100"}`,
                      `Annual Withdrawal: ${$displayMode === "monthly" ? formatCurrency(convertToMonthly(scenario.tspBalance * 0.04)) + "/mo" : formatCurrency(scenario.tspBalance * 0.04) + "/yr"}`,
                    ];
                  },
                },
              },
            },
            scales: {
              x: {
                title: {
                  display: true,
                  text: "TSP Balance ($000s)",
                },
                ticks: {
                  callback: function (value: any) {
                    return "$" + value + "k";
                  },
                },
              },
              y: {
                title: {
                  display: true,
                  text: "Estimated Depletion Age",
                },
                min: 60,
                max: 100,
                ticks: {
                  callback: function (value: any) {
                    return "Age " + value;
                  },
                },
              },
            },
          },
        };

      default:
        return null;
    }
  });

  // Chart lifecycle
  onMount(() => {
    const config = chartConfig();
    if (canvasElement && config) {
      chartInstance = createChartInstance(canvasElement, config);
    }
  });

  $effect(() => {
    // Explicitly depend on displayMode to trigger reactivity
    const currentDisplayMode = $displayMode;
    // console.log(`ScenarioComparisonChart: Display mode changed to ${currentDisplayMode}`);

    const config = chartConfig();
    if (chartInstance && config) {
      chartInstance.data = config.data;
      chartInstance.options = config.options;
      chartInstance.update();
    } else if (canvasElement && config && !chartInstance) {
      chartInstance = createChartInstance(canvasElement, config);
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
  {#if scenarioData().length > 1}
    <div class="flex justify-between items-center mb-4">
      <h3 class="text-lg font-semibold text-gray-800 dark:text-gray-200">
        Scenario Comparison
      </h3>

      <!-- Chart Type Selector -->
      <div class="flex bg-gray-100 dark:bg-gray-700 rounded-lg p-1">
        <button
          onclick={() => (chartType = "income")}
          class="px-3 py-1 text-sm rounded {chartType === 'income'
            ? 'bg-white dark:bg-gray-600 shadow'
            : ''} transition-all"
        >
          Income
        </button>
        <button
          onclick={() => (chartType = "tax")}
          class="px-3 py-1 text-sm rounded {chartType === 'tax'
            ? 'bg-white dark:bg-gray-600 shadow'
            : ''} transition-all"
        >
          Tax Burden
        </button>
        <button
          onclick={() => (chartType = "sustainability")}
          class="px-3 py-1 text-sm rounded {chartType === 'sustainability'
            ? 'bg-white dark:bg-gray-600 shadow'
            : ''} transition-all"
        >
          Sustainability
        </button>
      </div>
    </div>

    <div
      class="h-80 bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-600 p-4"
    >
      <canvas bind:this={canvasElement}></canvas>
    </div>

    <!-- Comparison Summary Table -->
    <div class="mt-4 overflow-x-auto">
      <table class="min-w-full text-sm">
        <thead class="bg-gray-50 dark:bg-gray-700">
          <tr>
            <th
              class="px-4 py-2 text-left font-medium text-gray-700 dark:text-gray-300"
              >Scenario</th
            >
            <th
              class="px-4 py-2 text-right font-medium text-gray-700 dark:text-gray-300"
              >Total Income ({$displayMode === "monthly"
                ? "Monthly"
                : "Annual"})</th
            >
            <th
              class="px-4 py-2 text-right font-medium text-gray-700 dark:text-gray-300"
              >Net Income ({$displayMode === "monthly"
                ? "Monthly"
                : "Annual"})</th
            >
            <th
              class="px-4 py-2 text-right font-medium text-gray-700 dark:text-gray-300"
              >Tax Rate</th
            >
            <th
              class="px-4 py-2 text-right font-medium text-gray-700 dark:text-gray-300"
              >TSP Balance</th
            >
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200 dark:divide-gray-600">
          {#each scenarioData() as scenario}
            <tr
              class="bg-white dark:bg-gray-800 {scenario.name.includes(
                '(Current)'
              )
                ? 'ring-2 ring-blue-500 ring-opacity-50'
                : ''}"
            >
              <td
                class="px-4 py-2 font-medium text-gray-900 dark:text-gray-100"
              >
                {scenario.name}
              </td>
              <td class="px-4 py-2 text-right text-gray-700 dark:text-gray-300">
                {$displayMode === "monthly"
                  ? formatCurrency(convertToMonthly(scenario.totalIncome)) +
                    "/mo"
                  : formatCurrency(scenario.totalIncome) + "/yr"}
              </td>
              <td class="px-4 py-2 text-right text-gray-700 dark:text-gray-300">
                {$displayMode === "monthly"
                  ? formatCurrency(convertToMonthly(scenario.netIncome)) + "/mo"
                  : formatCurrency(scenario.netIncome) + "/yr"}
              </td>
              <td class="px-4 py-2 text-right text-gray-700 dark:text-gray-300">
                {(scenario.effectiveTaxRate * 100).toFixed(1)}%
              </td>
              <td class="px-4 py-2 text-right text-gray-700 dark:text-gray-300">
                {formatCurrency(scenario.tspBalance)}
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {:else if scenarioData().length === 1}
    <div
      class="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-4 text-center"
    >
      <div class="text-blue-800 dark:text-blue-200 font-medium">
        Single Scenario
      </div>
      <div class="text-blue-600 dark:text-blue-400 text-sm">
        Create additional scenarios to see comparisons
      </div>
    </div>
  {:else}
    <div
      class="bg-gray-50 dark:bg-gray-700 rounded-lg border-2 border-dashed border-gray-300 dark:border-gray-600 p-8 text-center"
    >
      <div class="text-gray-500 dark:text-gray-400">
        <div class="text-lg font-medium mb-2">No Scenarios to Compare</div>
        <div class="text-sm">
          Complete calculations for multiple scenarios to see comparisons
        </div>
      </div>
    </div>
  {/if}
</div>
