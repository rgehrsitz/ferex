<script lang="ts">
  import type { Scenario } from "../../types";
  import ChartsSection from "./charts/ChartsSection.svelte";
  import RiskAssessmentChart from "./charts/RiskAssessmentChart.svelte";
  import ScenarioComparisonChart from "./charts/ScenarioComparisonChart.svelte";

  // Props - scenarios and activeScenarioId passed from parent
  let { scenarios = [], activeScenarioId } = $props<{
    scenarios?: Scenario[];
    activeScenarioId?: string | undefined;
  }>();

  // Analysis tab state
  let activeTab = $state<"overview" | "risk" | "comparison">("overview");

  // Get active scenario and results
  let activeScenario = $derived(
    scenarios.find((s: Scenario) => s.id === activeScenarioId)
  );

  // Simple reactive variables for results
  let results = $state<any>(null);
  let hasResults = $state(false);
  let keyMetrics = $state<any>(null);

  // Update results when activeScenario changes
  $effect(() => {
    const scenario = activeScenario;
    if (scenario?.results) {
      results = scenario.results;

      // Check if we have meaningful data
      hasResults = !!(
        results.TotalRetirementIncome > 0 ||
        results.FERSResult ||
        results.CSRSResult ||
        results.SocialSecurityResult ||
        results.TSPResult
      );

      // Update key metrics
      if (hasResults) {
        keyMetrics = {
          netIncome: results.NetAfterTaxIncome,
          totalIncome: results.TotalRetirementIncome,
          effectiveTaxRate: results.EffectiveTaxRate,
          pension:
            scenario.inputs.calculationSystem === "FERS"
              ? results.FERSResult?.annualPension
              : results.CSRSResult?.finalAnnuity,
          srs: results.SRSResult?.AnnualSRSAmount,
          socialSecurity: results.SocialSecurityResult?.ClaimingAmount,
          tspBalance: results.TSPResult?.totalProjectedBalanceAtRetirement,
          monthlyPension:
            scenario.inputs.calculationSystem === "FERS"
              ? results.FERSResult?.monthlyPension
              : results.CSRSResult?.monthlyFinalAnnuity,
        };
      } else {
        keyMetrics = null;
      }
    } else {
      results = null;
      hasResults = false;
      keyMetrics = null;
    }
  });

  // Helper function to format currency
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

  // Helper function to format percentage
  function formatPercent(
    value: number | null | undefined,
    defaultText = "N/A"
  ): string {
    if (value == null || isNaN(value)) return defaultText;
    return `${(value * 100).toFixed(1)}%`;
  }

  // Helper function to format TSP withdrawal strategy
  function formatTSPStrategy(
    strategy: string | undefined,
    inputs: any
  ): string {
    if (!strategy) return "Not specified";

    switch (strategy) {
      case "PercentageOfBalanceYearly":
        const percentage = inputs?.tspWithdrawalPercentageValue || "Unknown";
        return `${percentage}% of balance annually`;
      case "FixedAmountYearly":
      case "FixedDollarAmount":
        const amount = inputs?.tspWithdrawalFixedAmountValue || 0;
        return `Fixed $${amount.toLocaleString()} annually`;
      case "LifeExpectancy":
        return "Based on life expectancy";
      case "IRSMinimumRequiredDistribution":
        return "IRS Required Minimum Distribution";
      case "AsNeededToMeetIncomeGoal":
        return "As needed for income goal";
      case "None":
        return "No withdrawals planned";
      default:
        return strategy;
    }
  }
</script>

<div class="flex flex-col gap-6 w-full">
  <!-- Header -->
  <div class="flex justify-between items-center">
    <h1 class="text-3xl font-bold text-blue-900 dark:text-blue-200">
      Analysis & Results
    </h1>
    {#if activeScenario}
      <div class="text-sm text-gray-600 dark:text-gray-400">
        Scenario: <span class="font-semibold">{activeScenario.name}</span>
      </div>
    {/if}
  </div>

  {#if !hasResults}
    <!-- No Results State -->
    <div
      class="bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg p-6 text-center"
    >
      <div
        class="text-yellow-800 dark:text-yellow-200 text-lg font-semibold mb-2"
      >
        No Analysis Available
      </div>
      <div class="text-yellow-700 dark:text-yellow-300">
        Please complete the scenario inputs and run calculations to see results
        here.
      </div>
    </div>
  {:else}
    <!-- Results Display -->

    <!-- Top Level Summary Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <div
        class="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-6 border-l-4 border-blue-500"
      >
        <div
          class="text-sm font-semibold text-gray-600 dark:text-gray-400 uppercase tracking-wide"
        >
          Total Annual Income
        </div>
        <div class="text-3xl font-bold text-gray-900 dark:text-gray-100 mt-2">
          {formatCurrency(keyMetrics?.totalIncome)}
        </div>
        <div class="text-sm text-gray-500 dark:text-gray-400 mt-1">
          Gross retirement income
        </div>
      </div>

      <div
        class="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-6 border-l-4 border-green-500"
      >
        <div
          class="text-sm font-semibold text-gray-600 dark:text-gray-400 uppercase tracking-wide"
        >
          Net After-Tax Income
        </div>
        <div class="text-3xl font-bold text-gray-900 dark:text-gray-100 mt-2">
          {formatCurrency(keyMetrics?.netIncome)}
        </div>
        <div class="text-sm text-gray-500 dark:text-gray-400 mt-1">
          Take-home amount
        </div>
      </div>

      <div
        class="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-6 border-l-4 border-purple-500"
      >
        <div
          class="text-sm font-semibold text-gray-600 dark:text-gray-400 uppercase tracking-wide"
        >
          Effective Tax Rate
        </div>
        <div class="text-3xl font-bold text-gray-900 dark:text-gray-100 mt-2">
          {formatPercent(keyMetrics?.effectiveTaxRate)}
        </div>
        <div class="text-sm text-gray-500 dark:text-gray-400 mt-1">
          Combined fed/state taxes
        </div>
      </div>

      <div
        class="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-6 border-l-4 border-orange-500"
      >
        <div
          class="text-sm font-semibold text-gray-600 dark:text-gray-400 uppercase tracking-wide"
        >
          TSP at Retirement
        </div>
        <div class="text-3xl font-bold text-gray-900 dark:text-gray-100 mt-2">
          {formatCurrency(keyMetrics?.tspBalance)}
        </div>
        <div class="text-sm text-gray-500 dark:text-gray-400 mt-1">
          Projected balance
        </div>
      </div>
    </div>

    <!-- Analysis Tabs -->
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow-lg overflow-hidden">
      <!-- Tab Navigation -->
      <div class="border-b border-gray-200 dark:border-gray-600">
        <nav class="flex space-x-8 px-6" aria-label="Analysis Tabs">
          <button
            onclick={() => (activeTab = "overview")}
            class="py-4 px-1 border-b-2 font-medium text-sm transition-colors {activeTab ===
            'overview'
              ? 'border-blue-500 text-blue-600 dark:text-blue-400'
              : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300 dark:text-gray-400 dark:hover:text-gray-300'}"
          >
            Overview & Projections
          </button>
          <button
            onclick={() => (activeTab = "risk")}
            class="py-4 px-1 border-b-2 font-medium text-sm transition-colors {activeTab ===
            'risk'
              ? 'border-blue-500 text-blue-600 dark:text-blue-400'
              : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300 dark:text-gray-400 dark:hover:text-gray-300'}"
          >
            Risk Analysis
          </button>
          <button
            onclick={() => (activeTab = "comparison")}
            class="py-4 px-1 border-b-2 font-medium text-sm transition-colors {activeTab ===
            'comparison'
              ? 'border-blue-500 text-blue-600 dark:text-blue-400'
              : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300 dark:text-gray-400 dark:hover:text-gray-300'}"
          >
            Scenario Comparison
          </button>
        </nav>
      </div>

      <!-- Tab Content -->
      <div class="p-6">
        {#if activeTab === "overview"}
          <ChartsSection {results} inputs={activeScenario?.inputs} />
        {:else if activeTab === "risk"}
          <RiskAssessmentChart
            data={{
              retirementAge:
                activeScenario?.inputs?.plannedRetirementDate &&
                activeScenario?.inputs?.dateOfBirth
                  ? Math.floor(
                      (new Date(
                        activeScenario.inputs.plannedRetirementDate
                      ).getTime() -
                        new Date(activeScenario.inputs.dateOfBirth).getTime()) /
                        (365.25 * 24 * 60 * 60 * 1000)
                    )
                  : 65,
              tspBalance: results?.TSPResult?.totalProjectedBalanceAtRetirement,
              totalIncome: results?.TotalRetirementIncome,
            }}
          />
        {:else if activeTab === "comparison"}
          <ScenarioComparisonChart {scenarios} {activeScenarioId} />
        {/if}
      </div>
    </div>

    <!-- Charts & Projections Section -->
    <!-- <ChartsSection results={results} inputs={activeScenario?.inputs} /> -->

    <!-- Detailed Results Grid -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Pension Analysis -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-6">
        <h3 class="text-xl font-bold text-gray-900 dark:text-gray-100 mb-4">
          Pension Analysis
        </h3>
        <div class="space-y-4">
          {#if results.FERSResult && activeScenario?.inputs.calculationSystem === "FERS"}
            <div class="border-b border-gray-200 dark:border-gray-700 pb-4">
              <div class="flex justify-between items-center mb-2">
                <span class="font-semibold text-gray-700 dark:text-gray-300"
                  >FERS Annual Pension</span
                >
                <span class="text-lg font-bold text-blue-600 dark:text-blue-400"
                  >{formatCurrency(results.FERSResult.annualPension)}</span
                >
              </div>
              <div class="text-sm text-gray-600 dark:text-gray-400 space-y-1">
                <div>
                  Monthly: {formatCurrency(results.FERSResult.monthlyPension)}
                </div>
                <div>
                  Service Years: {results.FERSResult.totalServiceYears?.toFixed(
                    1
                  ) || "N/A"}
                </div>
                <div>
                  Retirement Type: {results.FERSResult.retirementType || "N/A"}
                </div>
              </div>
            </div>
          {/if}

          {#if results.CSRSResult && activeScenario?.inputs.calculationSystem === "CSRS"}
            <div class="border-b border-gray-200 dark:border-gray-700 pb-4">
              <div class="flex justify-between items-center mb-2">
                <span class="font-semibold text-gray-700 dark:text-gray-300"
                  >CSRS Annual Pension</span
                >
                <span class="text-lg font-bold text-blue-600 dark:text-blue-400"
                  >{formatCurrency(results.CSRSResult.finalAnnuity)}</span
                >
              </div>
              <div class="text-sm text-gray-600 dark:text-gray-400 space-y-1">
                <div>
                  Monthly: {formatCurrency(
                    results.CSRSResult.monthlyFinalAnnuity
                  )}
                </div>
                <div>
                  Service Years: {results.CSRSResult.totalServiceYears?.toFixed(
                    1
                  ) || "N/A"}
                </div>
                <div>
                  Retirement Type: {results.CSRSResult.retirementType || "N/A"}
                </div>
              </div>
            </div>
          {/if}

          {#if results.SRSResult && results.SRSResult.AnnualSRSAmount > 0}
            <div class="border-b border-gray-200 dark:border-gray-700 pb-4">
              <div class="flex justify-between items-center mb-2">
                <span class="font-semibold text-gray-700 dark:text-gray-300"
                  >FERS Supplement (SRS)</span
                >
                <span
                  class="text-lg font-bold text-purple-600 dark:text-purple-400"
                  >{formatCurrency(results.SRSResult.AnnualSRSAmount)}</span
                >
              </div>
              <div class="text-sm text-gray-600 dark:text-gray-400">
                Until age 62 • Monthly: {formatCurrency(
                  results.SRSResult.MonthlySRSAmount
                )}
              </div>
            </div>
          {/if}
        </div>
      </div>

      <!-- Social Security & Other Income -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-6">
        <h3 class="text-xl font-bold text-gray-900 dark:text-gray-100 mb-4">
          Social Security & Benefits
        </h3>
        <div class="space-y-4">
          {#if results.SocialSecurityResult}
            <div class="border-b border-gray-200 dark:border-gray-700 pb-4">
              <div class="flex justify-between items-center mb-2">
                <span class="font-semibold text-gray-700 dark:text-gray-300"
                  >Social Security (Annual)</span
                >
                <span
                  class="text-lg font-bold text-green-600 dark:text-green-400"
                  >{formatCurrency(
                    results.SocialSecurityResult.ClaimingAmount
                  )}</span
                >
              </div>
              <div class="text-sm text-gray-600 dark:text-gray-400 space-y-1">
                <div>
                  Claiming Age: {results.SocialSecurityResult.ClaimingAge ||
                    "N/A"}
                </div>
                <div>
                  Monthly: {formatCurrency(
                    (results.SocialSecurityResult.ClaimingAmount || 0) / 12
                  )}
                </div>
              </div>
            </div>
          {/if}

          {#if results.TSPResult}
            <div class="border-b border-gray-200 dark:border-gray-700 pb-4">
              <div class="flex justify-between items-center mb-2">
                <span class="font-semibold text-gray-700 dark:text-gray-300"
                  >TSP Strategy</span
                >
                <span class="text-sm text-gray-600 dark:text-gray-400"
                  >{formatTSPStrategy(
                    activeScenario?.inputs?.tspWithdrawalStrategy,
                    activeScenario?.inputs
                  )}</span
                >
              </div>
              <div class="text-sm text-gray-600 dark:text-gray-400 space-y-1">
                <div>
                  Balance at Retirement: {formatCurrency(
                    results.TSPResult.totalProjectedBalanceAtRetirement
                  )}
                </div>
                <div>
                  Traditional: {formatCurrency(
                    results.TSPResult.projectedTraditionalBalanceAtRetirement
                  )}
                </div>
                <div>
                  Roth: {formatCurrency(
                    results.TSPResult.projectedRothBalanceAtRetirement
                  )}
                </div>
              </div>
            </div>
          {/if}
        </div>
      </div>
    </div>

    <!-- Tax Analysis -->
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-6">
      <h3 class="text-xl font-bold text-gray-900 dark:text-gray-100 mb-4">
        Tax Analysis
      </h3>
      {#if results.TaxResult}
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div class="text-center">
            <div
              class="text-sm font-semibold text-gray-600 dark:text-gray-400 uppercase tracking-wide"
            >
              Federal Tax
            </div>
            <div
              class="text-2xl font-bold text-gray-900 dark:text-gray-100 mt-1"
            >
              {formatCurrency(results.TaxResult.FederalTaxOwed)}
            </div>
            <div class="text-sm text-gray-500 dark:text-gray-400">
              Annual liability
            </div>
          </div>
          <div class="text-center">
            <div
              class="text-sm font-semibold text-gray-600 dark:text-gray-400 uppercase tracking-wide"
            >
              State Tax
            </div>
            <div
              class="text-2xl font-bold text-gray-900 dark:text-gray-100 mt-1"
            >
              {formatCurrency(results.TaxResult.StateTaxOwed)}
            </div>
            <div class="text-sm text-gray-500 dark:text-gray-400">
              Annual liability
            </div>
          </div>
          <div class="text-center">
            <div
              class="text-sm font-semibold text-gray-600 dark:text-gray-400 uppercase tracking-wide"
            >
              Total Tax
            </div>
            <div
              class="text-2xl font-bold text-gray-900 dark:text-gray-100 mt-1"
            >
              {formatCurrency(
                (results.TaxResult.FederalTaxOwed || 0) +
                  (results.TaxResult.StateTaxOwed || 0)
              )}
            </div>
            <div class="text-sm text-gray-500 dark:text-gray-400">
              Combined annual
            </div>
          </div>
        </div>
      {:else}
        <div class="text-center text-gray-500 dark:text-gray-400">
          No tax calculation results available
        </div>
      {/if}
    </div>

    <!-- Debug Information -->
    <div
      class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-6"
    >
      <h3 class="text-lg font-bold text-red-900 dark:text-red-100 mb-2">
        🐛 Debug: Income Components Breakdown
      </h3>
      <div class="text-sm space-y-1 font-mono">
        <div>
          <strong>Backend Total:</strong>
          {formatCurrency(results.TotalRetirementIncome)}
        </div>
        <div>
          <strong>FERS Pension:</strong>
          {formatCurrency(results.FERSResult?.annualPension)}
        </div>
        <div>
          <strong>Social Security:</strong>
          {formatCurrency(results.SocialSecurityResult?.ClaimingAmount)}
        </div>
        <div>
          <strong>SRS Amount:</strong>
          {formatCurrency(results.SRSResult?.AnnualSRSAmount)}
        </div>
        <div>
          <strong>TSP Balance (not income):</strong>
          {formatCurrency(results.TSPResult?.totalProjectedBalanceAtRetirement)}
        </div>
        {#if results.TSPResult?.withdrawalSchedule && results.TSPResult.withdrawalSchedule.length > 0}
          <div>
            <strong>TSP First Year Withdrawal:</strong>
            {formatCurrency(
              results.TSPResult.withdrawalSchedule[0]?.totalWithdrawn
            )}
          </div>
          <div class="text-xs text-gray-600">
            • Traditional Withdrawn: {formatCurrency(
              results.TSPResult.withdrawalSchedule[0]?.traditionalWithdrawn
            )}
          </div>
          <div class="text-xs text-gray-600">
            • Roth Withdrawn: {formatCurrency(
              results.TSPResult.withdrawalSchedule[0]?.rothWithdrawn
            )}
          </div>
          <div class="text-xs text-gray-600">
            TSP Withdrawal Fields: {JSON.stringify(
              Object.keys(results.TSPResult.withdrawalSchedule[0] || {})
            )}
          </div>
          <div class="text-xs text-green-600">
            ✅ TSP Withdrawal Working: totalWithdrawn = {results.TSPResult
              .withdrawalSchedule[0]?.totalWithdrawn} (expected ~{formatCurrency(
              results.TSPResult.withdrawalSchedule[0]?.beginningBalanceTotal *
                0.04
            )})
          </div>
        {:else}
          <div class="text-red-600">No TSP withdrawal schedule found</div>
        {/if}
        {#if results.TSPResult}
          <div class="text-xs text-gray-600">
            TSP Result Fields: {JSON.stringify(Object.keys(results.TSPResult))}
          </div>
        {/if}
        <div class="pt-2 border-t border-red-300">
          {#if results}
            {@const pensionAmount = results.FERSResult?.annualPension || 0}
            {@const ssAmount =
              results.SocialSecurityResult?.ClaimingAmount || 0}
            {@const tspAnnualWithdrawal =
              results.TSPResult?.withdrawalSchedule?.[0]?.totalWithdrawn || 0}
            {@const manualSum = pensionAmount + ssAmount + tspAnnualWithdrawal}
            {@const backendTotal = results.TotalRetirementIncome || 0}
            {@const difference = backendTotal - manualSum}

            <strong>Income Component Analysis:</strong>
            <div class="ml-4 text-sm">
              • FERS Pension: {formatCurrency(pensionAmount)}
              • Social Security: {formatCurrency(ssAmount)}
              • TSP Annual Withdrawal: {formatCurrency(tspAnnualWithdrawal)}
              • SRS Amount: {formatCurrency(
                results.SRSResult?.AnnualSRSAmount || 0
              )}
              • <strong>Manual Sum:</strong>
              {formatCurrency(
                manualSum + (results.SRSResult?.AnnualSRSAmount || 0)
              )}
              • <strong>Backend Total:</strong>
              {formatCurrency(backendTotal)}
              {#if Math.abs(difference - (results.SRSResult?.AnnualSRSAmount || 0)) > 1}
                {@const updatedDifference =
                  backendTotal -
                  (manualSum + (results.SRSResult?.AnnualSRSAmount || 0))}
                •
                <span class="text-red-600"
                  >⚠️ Difference: {formatCurrency(updatedDifference)} (missing income
                  source?)</span
                >
                {#if activeScenario?.inputs}
                  <div class="mt-2 text-xs">
                    <strong>Checking other income sources:</strong>
                    <div>
                      • Other Taxable Income (Federal): {formatCurrency(
                        activeScenario.inputs.federalOtherTaxableIncomeAnnual ||
                          0
                      )}
                    </div>
                    {#if activeScenario.inputs}
                      {@const calculatedOtherIncome =
                        (activeScenario.inputs
                          .federalOtherTaxableIncomeAnnual || 0) +
                        (activeScenario.inputs.otherRecurringIncomeSources?.reduce(
                          (sum: number, source: any) => {
                            const amount = source.amount || 0;
                            return (
                              sum +
                              (source.frequency === "Monthly"
                                ? amount * 12
                                : amount)
                            );
                          },
                          0
                        ) || 0) +
                        (activeScenario.inputs.oneTimeIncomeEvents?.reduce(
                          (sum: number, event: any) => {
                            const eventYear = parseInt(
                              event.date.split("-")[0]
                            );
                            const currentYear = new Date().getFullYear();
                            return (
                              sum +
                              (eventYear === currentYear
                                ? event.amount || 0
                                : 0)
                            );
                          },
                          0
                        ) || 0)}
                      <div>
                        • Calculated Total Other Income: {formatCurrency(
                          calculatedOtherIncome
                        )}
                      </div>
                    {/if}
                    <div>
                      • Other Recurring Income Sources: {activeScenario.inputs
                        .otherRecurringIncomeSources?.length || 0} items
                    </div>
                    {#if activeScenario.inputs.otherRecurringIncomeSources && activeScenario.inputs.otherRecurringIncomeSources.length > 0}
                      {#each activeScenario.inputs.otherRecurringIncomeSources as incomeSource}
                        {@const annualAmount =
                          incomeSource.frequency === "Monthly"
                            ? (incomeSource.amount || 0) * 12
                            : incomeSource.amount || 0}
                        <div class="ml-4">
                          - {incomeSource.name}: {formatCurrency(annualAmount)} ({incomeSource.frequency})
                        </div>
                      {/each}
                    {/if}
                    <div>
                      • One-Time Income Events: {activeScenario.inputs
                        .oneTimeIncomeEvents?.length || 0} items
                    </div>
                    {#if activeScenario.inputs.oneTimeIncomeEvents && activeScenario.inputs.oneTimeIncomeEvents.length > 0}
                      {#each activeScenario.inputs.oneTimeIncomeEvents as event}
                        {@const eventYear = event.date.split("-")[0]}
                        <div class="ml-4">
                          - {event.name}: {formatCurrency(event.amount || 0)} in
                          {eventYear}
                        </div>
                      {/each}
                    {/if}
                  </div>
                {/if}
              {:else}
                • <span class="text-green-600">✅ Components match total</span>
              {/if}
            </div>
          {/if}
        </div>
        <div class="pt-2 border-t border-red-300">
          <strong>Backend Income Investigation:</strong>
        </div>
        <div>
          Backend TotalRetirementIncome: {formatCurrency(
            results.TotalRetirementIncome
          )}
        </div>
        <div>
          Tax Calculation Gross Income: {formatCurrency(
            results.TaxResult?.GrossIncome
          )}
        </div>
        <div>
          Tax Calculation Taxable Income: {formatCurrency(
            results.TaxResult?.TaxableIncome
          )}
        </div>
        <div class="text-xs text-red-600 mt-2">
          <strong>🔍 Income Source Breakdown:</strong>
          <div>
            Expected: FERS ($84,826) + SS ($33,540) + TSP ($76,724) = {formatCurrency(
              84826 + 33540 + 76724
            )}
          </div>
          <div>
            Backend Total: {formatCurrency(results.TotalRetirementIncome)}
          </div>
          <div>
            Missing: {formatCurrency(
              (results.TotalRetirementIncome || 0) - (84826 + 33540 + 76724)
            )}
            {#if Math.abs((results.TotalRetirementIncome || 0) - (84826 + 33540 + 76724) - 62278) < 100}
              ⚠️ Consistent $62k difference suggests backend auto-calculation
            {/if}
          </div>
          <div class="text-xs mt-2">
            <strong>🔬 Backend Result Fields Debug:</strong>
            <div>
              All Result Fields: {JSON.stringify(Object.keys(results || {}))}
            </div>
            {#if results.TaxResult}
              <div>
                Tax Result Fields: {JSON.stringify(
                  Object.keys(results.TaxResult)
                )}
              </div>
            {/if}
          </div>
        </div>
        {#if results.TaxResult?.GrossIncome && results.TotalRetirementIncome}
          {@const grossVsRetirement =
            (results.TaxResult.GrossIncome || 0) -
            (results.TotalRetirementIncome || 0)}
          <div class="text-sm text-red-600">
            • Difference (Tax Gross - Retirement Total): {formatCurrency(
              grossVsRetirement
            )}
          </div>
        {/if}
        <div>
          Federal Tax: {formatCurrency(results.TaxResult?.FederalTaxOwed)}
        </div>
        <div>State Tax: {formatCurrency(results.TaxResult?.StateTaxOwed)}</div>
        <div>
          Net After Tax: {formatCurrency(results.TaxResult?.NetAfterTaxIncome)}
        </div>
        <div>Effective Rate: {formatPercent(results.EffectiveTaxRate)}</div>
        <div class="text-xs text-red-600 mt-2">
          {#if results.TotalRetirementIncome && results.TaxResult?.FederalTaxOwed !== undefined && results.TaxResult?.StateTaxOwed !== undefined}
            {@const totalTax =
              (results.TaxResult.FederalTaxOwed || 0) +
              (results.TaxResult.StateTaxOwed || 0)}
            {@const calculatedNetIncome =
              results.TotalRetirementIncome - totalTax}
            {@const actualEffectiveRate =
              results.TotalRetirementIncome > 0
                ? totalTax / results.TotalRetirementIncome
                : 0}
            Math Check: ${formatCurrency(results.TotalRetirementIncome)} - ${formatCurrency(
              totalTax
            )} = {formatCurrency(calculatedNetIncome)}
            <br />Actual Effective Rate: {formatPercent(actualEffectiveRate)}
            <br />Backend Net Income: {formatCurrency(
              results.NetAfterTaxIncome
            )}
            {#if Math.abs(calculatedNetIncome - (results.NetAfterTaxIncome || 0)) > 1}
              <span class="text-red-700">
                ⚠️ MISMATCH: Difference of {formatCurrency(
                  Math.abs(
                    calculatedNetIncome - (results.NetAfterTaxIncome || 0)
                  )
                )}</span
              >
            {:else}
              <span class="text-green-700"> ✅ MATCH</span>
            {/if}
          {:else}
            Math Check: Missing tax calculation data
          {/if}
        </div>
      </div>
    </div>

    <!-- Notes and Additional Information -->
    {#if results.Notes}
      <div
        class="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-6"
      >
        <h3 class="text-lg font-bold text-blue-900 dark:text-blue-200 mb-2">
          Calculation Notes
        </h3>
        <div class="text-blue-800 dark:text-blue-300 whitespace-pre-wrap">
          {results.Notes}
        </div>
      </div>
    {/if}
  {/if}
</div>
