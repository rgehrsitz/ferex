<script lang="ts">
  import type { Scenario } from "../../types";
  import { displayMode } from "../stores/displayPreferences";
  import DisplayToggle from "./ui/DisplayToggle.svelte";
  import {
    extractMonthlyProjections,
    aggregateToYearly,
    formatChartCurrency,
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
  let riskChartCanvas = $state<HTMLCanvasElement>();
  let sensitivityChartCanvas = $state<HTMLCanvasElement>();
  let riskChart: Chart | null = null;
  let sensitivityChart: Chart | null = null;

  // Risk analysis data
  let riskMetrics = $derived(() => {
    if (!hasResults || !results) return null;

    const monthlyData = extractMonthlyProjections(results);
    if (monthlyData.length === 0) return null;

    // Calculate income volatility and risk metrics
    const incomes = monthlyData.map((d) => d.totalPreTaxIncome);
    const avgIncome =
      incomes.reduce((sum, val) => sum + val, 0) / incomes.length;
    const variance =
      incomes.reduce((sum, val) => sum + Math.pow(val - avgIncome, 2), 0) /
      incomes.length;
    const stdDev = Math.sqrt(variance);
    const coefficientOfVariation = stdDev / avgIncome;

    // Income source diversification
    const avgPension =
      monthlyData.reduce((sum, d) => sum + d.pension, 0) / monthlyData.length;
    const avgSS =
      monthlyData.reduce((sum, d) => sum + d.socialSecurity, 0) /
      monthlyData.length;
    const avgTSP =
      monthlyData.reduce((sum, d) => sum + d.totalTspWithdrawal, 0) /
      monthlyData.length;

    const totalAvg = avgPension + avgSS + avgTSP;
    const pensionPct = totalAvg > 0 ? avgPension / totalAvg : 0;
    const ssPct = totalAvg > 0 ? avgSS / totalAvg : 0;
    const tspPct = totalAvg > 0 ? avgTSP / totalAvg : 0;

    // Longevity risk (years of coverage)
    const yearsOfData = monthlyData.length / 12;
    const retirementAge = monthlyData[0]?.ageYears || 60;
    const projectedLifespan = 85; // Assumption
    const coverageRatio = yearsOfData / (projectedLifespan - retirementAge);

    // Inflation risk
    const inputs = activeScenario?.inputs;
    const inflationRate =
      (inputs?.userAssumedGeneralInflationRate || 2.5) / 100;
    const realReturnRisk =
      inflationRate > 0.03 ? "High" : inflationRate > 0.02 ? "Medium" : "Low";

    return {
      avgIncome,
      stdDev,
      coefficientOfVariation,
      pensionPct,
      ssPct,
      tspPct,
      yearsOfData,
      coverageRatio,
      inflationRate,
      realReturnRisk,
      diversificationScore: 1 - Math.max(pensionPct, ssPct, tspPct), // Higher when more diversified
    };
  });

  // Risk score calculation
  let overallRiskScore = $derived(() => {
    const metrics = riskMetrics();
    if (!metrics) return null;

    let score = 0;
    let maxScore = 0;

    // Income volatility (0-30 points)
    const volatilityScore = Math.min(30, metrics.coefficientOfVariation * 1000);
    score += volatilityScore;
    maxScore += 30;

    // Diversification (0-25 points, inverted - lower is better)
    const diversificationScore = (1 - metrics.diversificationScore) * 25;
    score += diversificationScore;
    maxScore += 25;

    // Longevity coverage (0-25 points)
    const longevityScore =
      metrics.coverageRatio < 0.8 ? 25 : metrics.coverageRatio < 1.0 ? 15 : 5;
    score += longevityScore;
    maxScore += 25;

    // Inflation risk (0-20 points)
    const inflationScore =
      metrics.inflationRate > 0.03 ? 20 : metrics.inflationRate > 0.02 ? 12 : 5;
    score += inflationScore;
    maxScore += 20;

    const normalizedScore = (score / maxScore) * 100;
    const riskLevel =
      normalizedScore > 70 ? "High" : normalizedScore > 40 ? "Medium" : "Low";

    return {
      score: normalizedScore,
      level: riskLevel,
      components: {
        volatility: volatilityScore,
        diversification: diversificationScore,
        longevity: longevityScore,
        inflation: inflationScore,
      },
    };
  });

  // Sensitivity analysis data
  let sensitivityData = $derived(() => {
    if (!hasResults || !results || !activeScenario?.inputs) return null;

    const baseIncome = results.TotalRetirementIncome || 0;

    // Simulate different scenarios
    const scenarios = [
      { label: "Base Case", change: 0, income: baseIncome },
      { label: "Inflation +1%", change: -5, income: baseIncome * 0.95 },
      { label: "Inflation +2%", change: -10, income: baseIncome * 0.9 },
      { label: "Market Down 20%", change: -15, income: baseIncome * 0.85 },
      { label: "Early Retirement", change: -12, income: baseIncome * 0.88 },
      { label: "Delayed Retirement", change: 8, income: baseIncome * 1.08 },
      { label: "Max SS Benefits", change: 15, income: baseIncome * 1.15 },
    ];

    return scenarios;
  });

  // Chart configurations
  let riskChartConfig = $derived(() => {
    const metrics = riskMetrics();
    if (!metrics) return null;

    return {
      type: "doughnut" as const,
      data: {
        labels: ["Pension", "Social Security", "TSP Withdrawal"],
        datasets: [
          {
            data: [
              metrics.pensionPct * 100,
              metrics.ssPct * 100,
              metrics.tspPct * 100,
            ],
            backgroundColor: [
              "rgba(59, 130, 246, 0.8)",
              "rgba(16, 185, 129, 0.8)",
              "rgba(139, 92, 246, 0.8)",
            ],
            borderColor: [
              "rgb(59, 130, 246)",
              "rgb(16, 185, 129)",
              "rgb(139, 92, 246)",
            ],
            borderWidth: 2,
          },
        ],
      },
      options: {
        ...defaultChartOptions,
        plugins: {
          ...defaultChartOptions.plugins,
          title: {
            display: true,
            text: "Income Source Diversification",
          },
          tooltip: {
            callbacks: {
              label: function (context: any) {
                const label = context.label || "";
                const value = context.parsed;
                return `${label}: ${value.toFixed(1)}%`;
              },
            },
          },
        },
      },
    };
  });

  let sensitivityChartConfig = $derived(() => {
    const data = sensitivityData();
    if (!data) return null;

    return {
      type: "bar" as const,
      data: {
        labels: data.map((d) => d.label),
        datasets: [
          {
            label: "Income Change (%)",
            data: data.map((d) => d.change),
            backgroundColor: data.map((d) =>
              d.change >= 0
                ? "rgba(16, 185, 129, 0.8)"
                : "rgba(239, 68, 68, 0.8)"
            ),
            borderColor: data.map((d) =>
              d.change >= 0 ? "rgb(16, 185, 129)" : "rgb(239, 68, 68)"
            ),
            borderWidth: 2,
          },
        ],
      },
      options: {
        ...defaultChartOptions,
        plugins: {
          ...defaultChartOptions.plugins,
          title: {
            display: true,
            text: "Sensitivity Analysis - Income Impact",
          },
          tooltip: {
            callbacks: {
              label: function (context: any) {
                const value = context.parsed.y;
                const sign = value >= 0 ? "+" : "";
                return `Impact: ${sign}${value}%`;
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
                return `${sign}${value}%`;
              },
            },
          },
        },
      },
    };
  });

  // Chart lifecycle management
  function createCharts() {
    if (riskChartCanvas && riskChartConfig()) {
      riskChart = createChartInstance(riskChartCanvas, riskChartConfig()!);
    }
    if (sensitivityChartCanvas && sensitivityChartConfig()) {
      sensitivityChart = createChartInstance(
        sensitivityChartCanvas,
        sensitivityChartConfig()!
      );
    }
  }

  function updateCharts() {
    if (riskChart && riskChartConfig()) {
      riskChart.data = riskChartConfig()!.data;
      riskChart.options = riskChartConfig()!.options;
      riskChart.update();
    }
    if (sensitivityChart && sensitivityChartConfig()) {
      sensitivityChart.data = sensitivityChartConfig()!.data;
      sensitivityChart.options = sensitivityChartConfig()!.options;
      sensitivityChart.update();
    }
  }

  function destroyCharts() {
    if (riskChart) {
      riskChart.destroy();
      riskChart = null;
    }
    if (sensitivityChart) {
      sensitivityChart.destroy();
      sensitivityChart = null;
    }
  }

  // Effects for chart management
  onMount(() => {
    createCharts();
  });

  $effect(() => {
    // React to data changes
    const currentRiskConfig = riskChartConfig();
    const currentSensitivityConfig = sensitivityChartConfig();
    if (currentRiskConfig || currentSensitivityConfig) {
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

  function formatPercentage(value: number): string {
    return `${(value * 100).toFixed(1)}%`;
  }

  function getRiskColor(level: string): string {
    switch (level) {
      case "Low":
        return "text-green-600 dark:text-green-400";
      case "Medium":
        return "text-yellow-600 dark:text-yellow-400";
      case "High":
        return "text-red-600 dark:text-red-400";
      default:
        return "text-gray-600 dark:text-gray-400";
    }
  }

  function getRiskBgColor(level: string): string {
    switch (level) {
      case "Low":
        return "bg-green-50 dark:bg-green-900/20";
      case "Medium":
        return "bg-yellow-50 dark:bg-yellow-900/20";
      case "High":
        return "bg-red-50 dark:bg-red-900/20";
      default:
        return "bg-gray-50 dark:bg-gray-900/20";
    }
  }
</script>

<div class="flex flex-col gap-6 w-full">
  <!-- Header -->
  <div class="flex justify-between items-center">
    <h1 class="text-3xl font-bold text-blue-900 dark:text-blue-200">
      Risk Analysis
    </h1>
    <div class="flex items-center gap-4">
      {#if activeScenario}
        <div class="text-sm text-gray-600 dark:text-gray-400">
          Scenario: <span class="font-semibold">{activeScenario.name}</span>
        </div>
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
        No Risk Analysis Available
      </div>
      <div class="text-yellow-700 dark:text-yellow-300">
        Please complete the scenario inputs and run calculations to see risk
        analysis.
      </div>
    </div>
  {:else}
    <!-- Overall Risk Score -->
    {#if overallRiskScore()}
      {@const risk = overallRiskScore()}
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-6">
        <h3 class="text-xl font-bold text-gray-900 dark:text-gray-100 mb-4">
          Overall Risk Assessment
        </h3>
        <div class="flex items-center justify-between mb-4">
          <div>
            <div class="text-3xl font-bold {getRiskColor(risk.level)}">
              {risk.level} Risk
            </div>
            <div class="text-sm text-gray-600 dark:text-gray-400">
              Risk Score: {risk.score.toFixed(0)}/100
            </div>
          </div>
          <div class="w-32 h-32 relative">
            <svg class="w-32 h-32 transform -rotate-90" viewBox="0 0 36 36">
              <path
                class="text-gray-300 dark:text-gray-600"
                stroke="currentColor"
                stroke-width="3"
                fill="none"
                d="M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831"
              />
              <path
                class={risk.level === "Low"
                  ? "text-green-500"
                  : risk.level === "Medium"
                    ? "text-yellow-500"
                    : "text-red-500"}
                stroke="currentColor"
                stroke-width="3"
                stroke-linecap="round"
                fill="none"
                stroke-dasharray="{risk.score}, 100"
                d="M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831"
              />
            </svg>
            <div class="absolute inset-0 flex items-center justify-center">
              <span class="text-xl font-bold text-gray-900 dark:text-gray-100">
                {risk.score.toFixed(0)}
              </span>
            </div>
          </div>
        </div>

        <!-- Risk Components -->
        <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
          <div class="text-center p-3 bg-gray-50 dark:bg-gray-900 rounded-lg">
            <div class="text-sm font-medium text-gray-700 dark:text-gray-300">
              Volatility
            </div>
            <div class="text-lg font-bold text-gray-900 dark:text-gray-100">
              {risk.components.volatility.toFixed(0)}/30
            </div>
          </div>
          <div class="text-center p-3 bg-gray-50 dark:bg-gray-900 rounded-lg">
            <div class="text-sm font-medium text-gray-700 dark:text-gray-300">
              Diversification
            </div>
            <div class="text-lg font-bold text-gray-900 dark:text-gray-100">
              {risk.components.diversification.toFixed(0)}/25
            </div>
          </div>
          <div class="text-center p-3 bg-gray-50 dark:bg-gray-900 rounded-lg">
            <div class="text-sm font-medium text-gray-700 dark:text-gray-300">
              Longevity
            </div>
            <div class="text-lg font-bold text-gray-900 dark:text-gray-100">
              {risk.components.longevity.toFixed(0)}/25
            </div>
          </div>
          <div class="text-center p-3 bg-gray-50 dark:bg-gray-900 rounded-lg">
            <div class="text-sm font-medium text-gray-700 dark:text-gray-300">
              Inflation
            </div>
            <div class="text-lg font-bold text-gray-900 dark:text-gray-100">
              {risk.components.inflation.toFixed(0)}/20
            </div>
          </div>
        </div>
      </div>
    {/if}

    <!-- Charts Section -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Income Diversification Chart -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-6">
        <div class="h-80">
          {#if riskChartConfig()}
            <canvas bind:this={riskChartCanvas}></canvas>
          {:else}
            <div
              class="h-full flex items-center justify-center text-gray-500 dark:text-gray-400"
            >
              <div class="text-center">
                <div class="text-lg font-medium mb-2">Loading Chart...</div>
                <div class="text-sm">Processing risk analysis data</div>
              </div>
            </div>
          {/if}
        </div>
      </div>

      <!-- Sensitivity Analysis Chart -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-6">
        <div class="h-80">
          {#if sensitivityChartConfig()}
            <canvas bind:this={sensitivityChartCanvas}></canvas>
          {:else}
            <div
              class="h-full flex items-center justify-center text-gray-500 dark:text-gray-400"
            >
              <div class="text-center">
                <div class="text-lg font-medium mb-2">Loading Chart...</div>
                <div class="text-sm">Processing sensitivity analysis</div>
              </div>
            </div>
          {/if}
        </div>
      </div>
    </div>

    <!-- Risk Metrics Details -->
    {#if riskMetrics()}
      {@const metrics = riskMetrics()}
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-6">
        <h3 class="text-xl font-bold text-gray-900 dark:text-gray-100 mb-4">
          Detailed Risk Metrics
        </h3>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <div>
            <h4 class="font-semibold text-gray-900 dark:text-gray-100 mb-2">
              Income Stability
            </h4>
            <div class="space-y-2 text-sm">
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400"
                  >Average Income:</span
                >
                <span class="font-medium"
                  >{formatCurrency(metrics.avgIncome)}</span
                >
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400"
                  >Standard Deviation:</span
                >
                <span class="font-medium">{formatCurrency(metrics.stdDev)}</span
                >
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400">Volatility:</span
                >
                <span class="font-medium"
                  >{formatPercentage(metrics.coefficientOfVariation)}</span
                >
              </div>
            </div>
          </div>

          <div>
            <h4 class="font-semibold text-gray-900 dark:text-gray-100 mb-2">
              Income Sources
            </h4>
            <div class="space-y-2 text-sm">
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400">Pension:</span>
                <span class="font-medium"
                  >{formatPercentage(metrics.pensionPct)}</span
                >
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400"
                  >Social Security:</span
                >
                <span class="font-medium"
                  >{formatPercentage(metrics.ssPct)}</span
                >
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400">TSP:</span>
                <span class="font-medium"
                  >{formatPercentage(metrics.tspPct)}</span
                >
              </div>
            </div>
          </div>

          <div>
            <h4 class="font-semibold text-gray-900 dark:text-gray-100 mb-2">
              Longevity & Inflation
            </h4>
            <div class="space-y-2 text-sm">
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400"
                  >Coverage Years:</span
                >
                <span class="font-medium">{metrics.yearsOfData.toFixed(1)}</span
                >
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400"
                  >Coverage Ratio:</span
                >
                <span class="font-medium"
                  >{formatPercentage(metrics.coverageRatio)}</span
                >
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600 dark:text-gray-400"
                  >Inflation Rate:</span
                >
                <span class="font-medium"
                  >{formatPercentage(metrics.inflationRate)}</span
                >
              </div>
            </div>
          </div>
        </div>
      </div>
    {/if}

    <!-- Risk Mitigation Recommendations -->
    {#if overallRiskScore()}
      {@const risk = overallRiskScore()}
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-6">
        <h3 class="text-xl font-bold text-gray-900 dark:text-gray-100 mb-4">
          Risk Mitigation Recommendations
        </h3>
        <div class="space-y-4">
          {#if risk.level === "High"}
            <div
              class="p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg"
            >
              <h4 class="font-semibold text-red-800 dark:text-red-200 mb-2">
                High Risk Detected
              </h4>
              <ul class="text-sm text-red-700 dark:text-red-300 space-y-1">
                <li>
                  • Consider delaying retirement to increase pension benefits
                </li>
                <li>• Diversify income sources beyond current plan</li>
                <li>
                  • Review TSP allocation for better risk-adjusted returns
                </li>
                <li>• Consider part-time work in early retirement</li>
              </ul>
            </div>
          {:else if risk.level === "Medium"}
            <div
              class="p-4 bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg"
            >
              <h4
                class="font-semibold text-yellow-800 dark:text-yellow-200 mb-2"
              >
                Moderate Risk
              </h4>
              <ul
                class="text-sm text-yellow-700 dark:text-yellow-300 space-y-1"
              >
                <li>• Monitor inflation impact on fixed income sources</li>
                <li>• Consider inflation-protected investments</li>
                <li>• Review withdrawal strategy for TSP</li>
                <li>• Plan for healthcare cost increases</li>
              </ul>
            </div>
          {:else}
            <div
              class="p-4 bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-lg"
            >
              <h4 class="font-semibold text-green-800 dark:text-green-200 mb-2">
                Low Risk Profile
              </h4>
              <ul class="text-sm text-green-700 dark:text-green-300 space-y-1">
                <li>• Well-diversified income sources</li>
                <li>• Good coverage for expected lifespan</li>
                <li>• Consider optimizing for tax efficiency</li>
                <li>• Regular review of assumptions recommended</li>
              </ul>
            </div>
          {/if}
        </div>
      </div>
    {/if}
  {/if}
</div>
