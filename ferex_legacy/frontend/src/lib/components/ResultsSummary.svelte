<script lang="ts">
  export let results: any = null;

  // Helper for formatting currency
  function formatCurrency(value: number | null | undefined): string {
    if (value == null) return "N/A";
    return new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: "USD",
    }).format(value);
  }

  // Helper for formatting years (to 2 decimals)
  function formatYears(value: number | null | undefined): string {
    if (value == null) return "N/A";
    return value.toFixed(2) + " years";
  }
</script>

{#if !results}
  <div class="text-gray-500 italic">No results available.</div>
{:else if typeof results !== "object"}
  <div class="text-red-500">Results data is invalid.</div>
{:else}
  <div class="space-y-6">
    <!-- Summary Section -->
    <div class="bg-green-50 border border-green-200 rounded p-4 mt-4">
      <div class="font-semibold text-green-700 text-lg mb-2">
        Calculation Results:
      </div>
      <!-- Debug output for results object -->
      <!--
      <pre
        class="bg-gray-100 text-xs rounded p-2 mb-4 overflow-x-auto">{JSON.stringify(
          {
            TotalRetirementIncome: results.TotalRetirementIncome,
            NetAfterTaxIncome: results.NetAfterTaxIncome,
            EffectiveTaxRate: results.EffectiveTaxRate,
            Notes: results.Notes
          },
          null,
          2
        )}</pre>
      -->
      <h3 class="text-xl font-bold text-blue-900 mb-2">Retirement Summary</h3>
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <div>
          <div class="font-semibold text-blue-700">Total Retirement Income</div>
          <div>{formatCurrency(results?.TotalRetirementIncome)}</div>
        </div>
        <div>
          <div class="font-semibold text-blue-700">Net After-Tax Income</div>
          <div>{formatCurrency(results?.NetAfterTaxIncome)}</div>
        </div>
        <div>
          <div class="font-semibold text-blue-700">Effective Tax Rate</div>
          <div>
            {results?.EffectiveTaxRate != null
              ? (results.EffectiveTaxRate * 100).toFixed(2) + "%"
              : "N/A"}
          </div>
        </div>
      </div>
      {#if results?.Notes}
        <div
          class="mt-3 text-yellow-900 bg-yellow-50 border-l-4 border-yellow-400 p-2 rounded"
        >
          <span class="font-semibold">Notes:</span>
          {results.Notes}
        </div>
      {/if}
    </div>

    <!-- FERS Results -->
    {#if results?.FERSResult}
      <details
        class="bg-white rounded shadow p-4 border border-gray-200 mb-2"
        open
      >
        <summary class="font-semibold text-blue-800 text-lg cursor-pointer mb-2"
          >FERS Pension Details</summary
        >
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <span class="font-semibold">Retirement Type:</span>
            {results?.FERSResult?.retirementType ?? "N/A"}
          </div>
          <div>
            <span class="font-semibold">Annual Pension:</span>
            {formatCurrency(results?.FERSResult?.annualPension)}
          </div>
          <div>
            <span class="font-semibold">Monthly Pension:</span>
            {formatCurrency(results?.FERSResult?.monthlyPension)}
          </div>
          <div>
            <span class="font-semibold">Total Service Years:</span>
            {formatYears(results?.FERSResult?.totalServiceYears)}
          </div>
          <div>
            <span class="font-semibold">Basic Annuity:</span>
            {formatCurrency(results?.FERSResult?.basicAnnuity)}
          </div>
          <div>
            <span class="font-semibold">FERS Supplement:</span>
            {formatCurrency(results?.FERSResult?.fersSupplement)}
          </div>
          <div>
            <span class="font-semibold">Prorated Pension:</span>
            {formatCurrency(results?.FERSResult?.ProratedPension)}
          </div>
          <div>
            <span class="font-semibold">Sick Leave Service Credit:</span>
            {formatYears(results?.FERSResult?.SickLeaveServiceCredit)}
          </div>
          <div>
            <span class="font-semibold">Early Retirement Reduction:</span>
            {formatCurrency(results?.FERSResult?.EarlyRetirementReduction)}
          </div>
          <div>
            <span class="font-semibold">Survivor Benefit Reduction:</span>
            {formatCurrency(results?.FERSResult?.SurvivorBenefitReduction)}
          </div>
        </div>
        {#if results?.FERSResult?.notes}
          <div
            class="mt-2 text-yellow-900 bg-yellow-50 border-l-4 border-yellow-400 p-2 rounded"
          >
            <span class="font-semibold">Notes:</span>
            {results.FERSResult.notes}
          </div>
        {/if}
      </details>
    {/if}

    <!-- CSRS Results -->
    {#if results?.CSRSResult}
      <details class="bg-white rounded shadow p-4 border border-gray-200 mb-2">
        <summary class="font-semibold text-blue-800 text-lg cursor-pointer mb-2"
          >CSRS Pension Details</summary
        >
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <span class="font-semibold">Retirement Type:</span>
            {results?.CSRSResult?.retirementType ?? "N/A"}
          </div>
          <div>
            <span class="font-semibold">Annual Pension:</span>
            {formatCurrency(results?.CSRSResult?.grossAnnuity)}
          </div>
          <div>
            <span class="font-semibold">Monthly Pension:</span>
            {formatCurrency(results?.CSRSResult?.monthlyGrossAnnuity)}
          </div>
          <div>
            <span class="font-semibold">Total Service Years:</span>
            {formatYears(results?.CSRSResult?.totalServiceYears)}
          </div>
        </div>
        {#if results?.CSRSResult?.notes}
          <div
            class="mt-2 text-yellow-900 bg-yellow-50 border-l-4 border-yellow-400 p-2 rounded"
          >
            <span class="font-semibold">Notes:</span>
            {results.CSRSResult.notes}
          </div>
        {/if}
      </details>
    {/if}

    <!-- SRS Results -->
    {#if results?.SRSResult}
      <details class="bg-white rounded shadow p-4 border border-gray-200 mb-2">
        <summary class="font-semibold text-blue-800 text-lg cursor-pointer mb-2"
          >FERS Annuity Supplement (SRS)</summary
        >
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <span class="font-semibold">Annual SRS Amount:</span>
            {formatCurrency(results?.SRSResult?.AnnualSRSAmount)}
          </div>
          <div>
            <span class="font-semibold">Monthly SRS Amount:</span>
            {formatCurrency(results?.SRSResult?.MonthlySRSAmount)}
          </div>
          <div>
            <span class="font-semibold">Earnings Test Reduction:</span>
            {formatCurrency(results?.SRSResult?.EarningsTestReduction)}
          </div>
          <div>
            <span class="font-semibold">Eligible:</span>
            {results?.SRSResult?.IsEligible ? "Yes" : "No"}
          </div>
        </div>
        {#if results?.SRSResult?.Notes}
          <div
            class="mt-2 text-yellow-900 bg-yellow-50 border-l-4 border-yellow-400 p-2 rounded"
          >
            <span class="font-semibold">Notes:</span>
            {results.SRSResult.Notes}
          </div>
        {/if}
      </details>
    {/if}

    <!-- TSP Results -->
    {#if results?.TSPResult}
      <details class="bg-white rounded shadow p-4 border border-gray-200 mb-2">
        <summary class="font-semibold text-blue-800 text-lg cursor-pointer mb-2"
          >TSP Projections</summary
        >
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <span class="font-semibold">Traditional Balance at Retirement:</span
            >
            {formatCurrency(
              results?.TSPResult?.projectedTraditionalBalanceAtRetirement
            )}
          </div>
          <div>
            <span class="font-semibold">Roth Balance at Retirement:</span>
            {formatCurrency(
              results?.TSPResult?.projectedRothBalanceAtRetirement
            )}
          </div>
          <div>
            <span class="font-semibold"
              >Total Projected Balance at Retirement:</span
            >
            {formatCurrency(
              results?.TSPResult?.totalProjectedBalanceAtRetirement
            )}
          </div>
        </div>
        {#if results?.TSPResult?.withdrawalSchedule && results.TSPResult.withdrawalSchedule.length > 0}
          <!-- Table rendering commented out for crash isolation -->
          <!--
          <div class="overflow-x-auto mt-4">
            <table class="min-w-full text-xs">
              <thead>
                <tr class="bg-gray-100">
                  <th class="px-2 py-1">Year</th>
                  <th class="px-2 py-1">Age</th>
                  <th class="px-2 py-1">Begin Bal</th>
                  <th class="px-2 py-1">Withdrawn</th>
                  <th class="px-2 py-1">End Bal</th>
                </tr>
              </thead>
              <tbody>
                {#each results?.TSPResult?.withdrawalSchedule?.slice(0, 20) as row}
                  <tr>
                    <td class="px-2 py-1">{row.year}</td>
                    <td class="px-2 py-1">{row.age}</td>
                    <td class="px-2 py-1"
                      >{formatCurrency(row.beginningBalanceTotal)}</td
                    >
                    <td class="px-2 py-1"
                      >{formatCurrency(row.totalWithdrawn)}</td
                    >
                    <td class="px-2 py-1"
                      >{formatCurrency(row.endingBalanceTotal)}</td
                    >
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
          -->
        {/if}
        {#if results?.TSPResult?.notes}
          <div
            class="mt-2 text-yellow-900 bg-yellow-50 border-l-4 border-yellow-400 p-2 rounded"
          >
            <span class="font-semibold">Notes:</span>
            {results.TSPResult.notes}
          </div>
        {/if}
      </details>
    {/if}

    <!-- Monte Carlo Results -->
    {#if results?.MonteCarloResult}
      <details
        class="bg-white rounded shadow p-4 border border-gray-200 mb-2"
        open
      >
        <summary class="font-semibold text-blue-800 text-lg cursor-pointer mb-2"
          >Monte Carlo Simulation</summary
        >
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <span class="font-semibold">Success Rate:</span>
            {results?.MonteCarloResult?.SuccessRate != null
              ? (results.MonteCarloResult.SuccessRate * 100).toFixed(1) + "%"
              : "N/A"}
          </div>
        </div>
        {#if results?.MonteCarloResult?.Percentiles}
          <!--
          <div class="mt-2">
            <span class="font-semibold">Percentiles:</span>
            <ul class="list-disc ml-6">
              {#each Object.entries(results?.MonteCarloResult?.Percentiles).slice(0, 20) as [percentile, value]}
                <li>{percentile}th: {formatCurrency(Number(value))}</li>
              {/each}
            </ul>
          </div>
          -->
        {/if}
        {#if results?.MonteCarloResult?.DepletionProbabilities}
          <!--
          <div class="mt-2">
            <span class="font-semibold">Depletion Probabilities:</span>
            <ul class="list-disc ml-6">
              {#each results?.MonteCarloResult?.DepletionProbabilities?.slice(0, 20) as prob, i}
                <li>Year {i + 1}: {(prob * 100).toFixed(1)}%</li>
              {/each}
            </ul>
          </div>
          -->
        {/if}
      </details>
    {/if}

    <!-- Social Security Results -->
    {#if results?.SocialSecurityResult}
      <details class="bg-white rounded shadow p-4 border border-gray-200 mb-2">
        <summary class="font-semibold text-blue-800 text-lg cursor-pointer mb-2"
          >Social Security</summary
        >
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <span class="font-semibold">Estimated at 62:</span>
            {formatCurrency(results?.SocialSecurityResult?.EstimatedAt62)}
          </div>
          <div>
            <span class="font-semibold">Estimated at FRA:</span>
            {formatCurrency(results?.SocialSecurityResult?.EstimatedAtFRA)}
          </div>
          <div>
            <span class="font-semibold">Estimated at 70:</span>
            {formatCurrency(results?.SocialSecurityResult?.EstimatedAt70)}
          </div>
          <div>
            <span class="font-semibold">Claiming Age:</span>
            {results?.SocialSecurityResult?.ClaimingAge ?? "N/A"}
          </div>
          <div>
            <span class="font-semibold">Claiming Amount:</span>
            {formatCurrency(results?.SocialSecurityResult?.ClaimingAmount)}
          </div>
        </div>
        {#if results?.SocialSecurityResult?.Notes}
          <div
            class="mt-2 text-yellow-900 bg-yellow-50 border-l-4 border-yellow-400 p-2 rounded"
          >
            <span class="font-semibold">Notes:</span>
            {results.SocialSecurityResult.Notes}
          </div>
        {/if}
      </details>
    {/if}

    <!-- Tax Results -->
    {#if results?.TaxResult}
      <details class="bg-white rounded shadow p-4 border border-gray-200 mb-2">
        <summary class="font-semibold text-blue-800 text-lg cursor-pointer mb-2"
          >Tax Results</summary
        >
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <span class="font-semibold">Federal Tax Owed:</span>
            {formatCurrency(results?.TaxResult?.FederalTaxOwed)}
          </div>
          <div>
            <span class="font-semibold">State Tax Owed:</span>
            {formatCurrency(results?.TaxResult?.StateTaxOwed)}
          </div>
          <div>
            <span class="font-semibold">Net After-Tax Income:</span>
            {formatCurrency(results?.TaxResult?.NetAfterTaxIncome)}
          </div>
          <div>
            <span class="font-semibold">Effective Tax Rate:</span>
            {results?.TaxResult?.EffectiveTaxRate != null
              ? (results.TaxResult.EffectiveTaxRate * 100).toFixed(2) + "%"
              : "N/A"}
          </div>
        </div>
        {#if results?.TaxResult?.Notes}
          <div
            class="mt-2 text-yellow-900 bg-yellow-50 border-l-4 border-yellow-400 p-2 rounded"
          >
            <span class="font-semibold">Notes:</span>
            {results.TaxResult.Notes}
          </div>
        {/if}
      </details>
    {/if}

    <!-- COLA Results -->
    {#if results?.COLAResult}
      <details class="bg-white rounded shadow p-4 border border-gray-200 mb-2">
        <summary class="font-semibold text-blue-800 text-lg cursor-pointer mb-2"
          >COLA Projections</summary
        >
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <span class="font-semibold">Final Amount:</span>
            {formatCurrency(results?.COLAResult?.FinalAmount)}
          </div>
          <div>
            <span class="font-semibold">Cumulative COLA:</span>
            {formatCurrency(results?.COLAResult?.CumulativeCOLA)}
          </div>
        </div>
        {#if results?.COLAResult?.ProjectedAmounts && results.COLAResult.ProjectedAmounts.length > 0}
          <!--
          <div class="overflow-x-auto mt-4">
            <table class="min-w-full text-xs">
              <thead>
                <tr class="bg-gray-100">
                  <th class="px-2 py-1">Year</th>
                  <th class="px-2 py-1">Projected Amount</th>
                </tr>
              </thead>
              <tbody>
                {#each results?.COLAResult?.ProjectedAmounts?.slice(0, 20) as amt, i}
                  <tr>
                    <td class="px-2 py-1">{i + 1}</td>
                    <td class="px-2 py-1">{formatCurrency(amt)}</td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
          -->
        {/if}
        {#if results?.COLAResult?.Notes}
          <div
            class="mt-2 text-yellow-900 bg-yellow-50 border-l-4 border-yellow-400 p-2 rounded"
          >
            <span class="font-semibold">Notes:</span>
            {results.COLAResult.Notes}
          </div>
        {/if}
      </details>
    {/if}

    <!-- Survivor Benefit Results -->
    {#if results?.SurvivorResult}
      <details class="bg-white rounded shadow p-4 border border-gray-200 mb-2">
        <summary class="font-semibold text-blue-800 text-lg cursor-pointer mb-2"
          >Survivor Benefit</summary
        >
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <span class="font-semibold">Initial Survivor Annuity:</span>
            {formatCurrency(results?.SurvivorResult?.InitialSurvivorAnnuity)}
          </div>
          <div>
            <span class="font-semibold">Total Survivor Income:</span>
            {formatCurrency(results?.SurvivorResult?.TotalSurvivorIncome)}
          </div>
        </div>
        {#if results?.SurvivorResult?.ProjectedAnnuities && results.SurvivorResult.ProjectedAnnuities.length > 0}
          <!--
          <div class="overflow-x-auto mt-4">
            <table class="min-w-full text-xs">
              <thead>
                <tr class="bg-gray-100">
                  <th class="px-2 py-1">Year</th>
                  <th class="px-2 py-1">Projected Annuity</th>
                </tr>
              </thead>
              <tbody>
                {#each results?.SurvivorResult?.ProjectedAnnuities?.slice(0, 20) as amt, i}
                  <tr>
                    <td class="px-2 py-1">{i + 1}</td>
                    <td class="px-2 py-1">{formatCurrency(amt)}</td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
          -->
        {/if}
        {#if results?.SurvivorResult?.Notes}
          <div
            class="mt-2 text-yellow-900 bg-yellow-50 border-l-4 border-yellow-400 p-2 rounded"
          >
            <span class="font-semibold">Notes:</span>
            {results.SurvivorResult.Notes}
          </div>
        {/if}
      </details>
    {/if}

    <!-- Health Premium Results -->
    {#if results?.HealthResult}
      <details class="bg-white rounded shadow p-4 border border-gray-200 mb-2">
        <summary class="font-semibold text-blue-800 text-lg cursor-pointer mb-2"
          >Health Premiums</summary
        >
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <span class="font-semibold">Total Premiums:</span>
            {formatCurrency(results?.HealthResult?.TotalPremiums)}
          </div>
          <div>
            <span class="font-semibold">Medicare Premiums:</span>
            {formatCurrency(results?.HealthResult?.MedicarePremiums)}
          </div>
        </div>
        {#if results?.HealthResult?.ProjectedPremiums && results.HealthResult.ProjectedPremiums.length > 0}
          <!--
          <div class="overflow-x-auto mt-4">
            <table class="min-w-full text-xs">
              <thead>
                <tr class="bg-gray-100">
                  <th class="px-2 py-1">Year</th>
                  <th class="px-2 py-1">Projected Premium</th>
                </tr>
              </thead>
              <tbody>
                {#each results?.HealthResult?.ProjectedPremiums?.slice(0, 20) as amt, i}
                  <tr>
                    <td class="px-2 py-1">{i + 1}</td>
                    <td class="px-2 py-1">{formatCurrency(amt)}</td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
          -->
        {/if}
        {#if results?.HealthResult?.Notes}
          <div
            class="mt-2 text-yellow-900 bg-yellow-50 border-l-4 border-yellow-400 p-2 rounded"
          >
            <span class="font-semibold">Notes:</span>
            {results.HealthResult.Notes}
          </div>
        {/if}
      </details>
    {/if}
  </div>
{/if}

<style>
  /* All styling is handled by Tailwind classes */
</style>
