<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import {
    Chart,
    chartColors,
    formatCurrency,
    createChartInstance,
  } from "./chartUtils";

  interface RiskMetrics {
    successProbabilityAt85: number;
    successProbabilityAt90: number;
    successProbabilityAt95: number;
    medianDepletionAge: number | null;
    volatilityScore: number;
    sustainabilityRating: string;
  }

  let { data } = $props<{ data: any }>();

  let canvasElement = $state<HTMLCanvasElement>();
  let chartInstance: Chart | null = null;

  // Calculate risk metrics
  let riskMetrics = $derived(() => {
    if (!data || !data.tspBalance) {
      return {
        successProbabilityAt85: 0,
        successProbabilityAt90: 0,
        successProbabilityAt95: 0,
        medianDepletionAge: null,
        volatilityScore: 0,
        sustainabilityRating: "Insufficient Data",
      };
    }

    const balance = data.tspBalance;
    const withdrawalRate = 0.04; // 4% rule
    const annualWithdrawal = balance * withdrawalRate;

    // Simplified risk calculations
    const balanceToWithdrawalRatio = balance / annualWithdrawal;

    // Success probability estimates based on withdrawal rate and balance
    const successAt85 = Math.min(
      0.95,
      Math.max(0.3, balanceToWithdrawalRatio / 30)
    );
    const successAt90 = Math.min(
      0.9,
      Math.max(0.2, balanceToWithdrawalRatio / 35)
    );
    const successAt95 = Math.min(
      0.8,
      Math.max(0.1, balanceToWithdrawalRatio / 40)
    );

    // Estimated depletion age
    const depletionAge =
      withdrawalRate > 0.05 ? 75 + balanceToWithdrawalRatio / 5 : null;

    // Volatility score (simplified)
    const volatilityScore = withdrawalRate > 0.045 ? 0.7 : 0.4;

    return {
      successProbabilityAt85: successAt85,
      successProbabilityAt90: successAt90,
      successProbabilityAt95: successAt95,
      medianDepletionAge: depletionAge,
      volatilityScore: volatilityScore,
      sustainabilityRating:
        successAt85 > 0.8 ? "Strong" : successAt85 > 0.6 ? "Moderate" : "Weak",
    };
  });

  // Get risk level assessment
  function getRiskLevel(metrics: RiskMetrics) {
    const avgSuccess =
      (metrics.successProbabilityAt85 +
        metrics.successProbabilityAt90 +
        metrics.successProbabilityAt95) /
      3;

    if (avgSuccess > 0.8) {
      return {
        level: "Low Risk",
        color: "text-green-600 dark:text-green-400",
        description: "Strong sustainability outlook",
      };
    } else if (avgSuccess > 0.6) {
      return {
        level: "Moderate Risk",
        color: "text-yellow-600 dark:text-yellow-400",
        description: "Generally sustainable with some uncertainty",
      };
    } else {
      return {
        level: "High Risk",
        color: "text-red-600 dark:text-red-400",
        description: "Consider adjusting withdrawal strategy",
      };
    }
  }

  // Create radar chart configuration
  let chartConfig = $derived(() => {
    const metrics = riskMetrics();
    if (!metrics) return null;

    return {
      type: "radar",
      data: {
        labels: [
          "Age 85 Success",
          "Age 90 Success",
          "Age 95 Success",
          "Balance Stability",
          "Low Volatility",
          "Sustainability",
        ],
        datasets: [
          {
            label: "Risk Profile",
            data: [
              metrics.successProbabilityAt85 * 100,
              metrics.successProbabilityAt90 * 100,
              metrics.successProbabilityAt95 * 100,
              metrics.medianDepletionAge
                ? Math.min(100, (metrics.medianDepletionAge - 65) * 2.5)
                : 50,
              (1 - metrics.volatilityScore) * 100,
              metrics.sustainabilityRating === "Strong"
                ? 90
                : metrics.sustainabilityRating === "Moderate"
                  ? 65
                  : 40,
            ],
            backgroundColor: "rgba(59, 130, 246, 0.2)",
            borderColor: "rgba(59, 130, 246, 1)",
            borderWidth: 2,
            pointBackgroundColor: "rgba(59, 130, 246, 1)",
            pointBorderColor: "#fff",
            pointHoverBackgroundColor: "#fff",
            pointHoverBorderColor: "rgba(59, 130, 246, 1)",
          },
        ],
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          title: {
            display: true,
            text: "Portfolio Risk Assessment",
          },
          legend: {
            display: false,
          },
        },
        scales: {
          r: {
            angleLines: {
              display: true,
            },
            suggestedMin: 0,
            suggestedMax: 100,
            ticks: {
              stepSize: 20,
              callback: function (value: any) {
                return value + "%";
              },
            },
          },
        },
      },
    };
  });

  // Chart lifecycle management
  onMount(() => {
    const config = chartConfig();
    if (canvasElement && config) {
      chartInstance = createChartInstance(canvasElement, config);
    }
  });

  $effect(() => {
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
  <h3 class="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-4">
    Risk Assessment
  </h3>

  {#if chartConfig()}
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Radar Chart -->
      <div
        class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-600 p-4"
      >
        <div class="h-80">
          <canvas bind:this={canvasElement}></canvas>
        </div>
      </div>

      <!-- Risk Metrics Summary -->
      {#if riskMetrics()}
        {@const metrics = riskMetrics()}
        {@const riskAssessment = getRiskLevel(metrics)}
        <div class="space-y-4">
          <!-- Overall Risk Level -->
          <div class="bg-gray-50 dark:bg-gray-700 rounded-lg p-4">
            <div class="font-semibold {riskAssessment.color} text-lg">
              {riskAssessment.level}
            </div>
            <div class="text-sm text-gray-600 dark:text-gray-400 mt-1">
              {riskAssessment.description}
            </div>
          </div>

          <!-- Success Probabilities -->
          <div
            class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-600 p-4"
          >
            <h4 class="font-semibold text-gray-800 dark:text-gray-200 mb-3">
              Success Probabilities
            </h4>
            <div class="space-y-2 text-sm">
              <div class="flex justify-between">
                <span>Age 85:</span>
                <span class="font-medium"
                  >{(metrics.successProbabilityAt85 * 100).toFixed(1)}%</span
                >
              </div>
              <div class="flex justify-between">
                <span>Age 90:</span>
                <span class="font-medium"
                  >{(metrics.successProbabilityAt90 * 100).toFixed(1)}%</span
                >
              </div>
              <div class="flex justify-between">
                <span>Age 95:</span>
                <span class="font-medium"
                  >{(metrics.successProbabilityAt95 * 100).toFixed(1)}%</span
                >
              </div>
            </div>
          </div>

          <!-- Additional Metrics -->
          <div
            class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-600 p-4"
          >
            <h4 class="font-semibold text-gray-800 dark:text-gray-200 mb-3">
              Key Indicators
            </h4>
            <div class="space-y-2 text-sm">
              <div class="flex justify-between">
                <span>Sustainability:</span>
                <span class="font-medium">{metrics.sustainabilityRating}</span>
              </div>
              <div class="flex justify-between">
                <span>Volatility Score:</span>
                <span class="font-medium"
                  >{(metrics.volatilityScore * 100).toFixed(0)}%</span
                >
              </div>
              {#if metrics.medianDepletionAge}
                <div class="flex justify-between">
                  <span>Est. Depletion Age:</span>
                  <span class="font-medium"
                    >{metrics.medianDepletionAge.toFixed(0)}</span
                  >
                </div>
              {/if}
            </div>
          </div>

          <!-- Recommendations -->
          <div
            class="bg-blue-50 dark:bg-blue-900/20 rounded-lg p-4 border border-blue-200 dark:border-blue-700"
          >
            <div class="font-medium text-blue-800 dark:text-blue-200 mb-2">
              Recommendations
            </div>
            <div class="text-blue-700 dark:text-blue-300 space-y-1 text-sm">
              {#if metrics.successProbabilityAt90 < 0.7}
                <div>• Consider reducing withdrawal rate to 3.5%</div>
              {/if}
              {#if metrics.volatilityScore > 0.5}
                <div>• Review asset allocation for better stability</div>
              {/if}
              {#if metrics.medianDepletionAge && metrics.medianDepletionAge < 85}
                <div>• Explore additional income sources</div>
              {/if}
              {#if metrics.successProbabilityAt85 > 0.9}
                <div>• Strong position - consider legacy planning</div>
              {/if}
            </div>
          </div>
        </div>
      {/if}
    </div>
  {:else}
    <div
      class="h-64 bg-gray-50 dark:bg-gray-700 rounded-lg border-2 border-dashed border-gray-300 dark:border-gray-600 flex flex-col items-center justify-center p-4"
    >
      <div class="text-center text-gray-500 dark:text-gray-400">
        <div class="text-lg font-medium mb-2">No Risk Data Available</div>
        <div class="text-sm">Complete TSP inputs to see risk assessment</div>
      </div>
    </div>
  {/if}
</div>
