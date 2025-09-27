<script lang="ts">
  import type { Scenario } from '../../types';
  export let scenario: Scenario | null = null;

  // Helper to format currency
  function fmtCurrency(val: number | null | undefined): string {
    if (val == null || isNaN(val)) return 'N/A';
    return `$${val.toLocaleString(undefined, { maximumFractionDigits: 0 })}`;
  }

  // Extract key metrics from scenario.results, fallback to N/A if missing
  $: netIncome = scenario?.results?.NetAfterTaxIncome ?? null;
  $: tspBalance = scenario?.results?.TSPResult?.totalProjectedBalanceAtRetirement ?? null;
  $: pension = scenario?.results?.FERSResult?.annualPension ?? scenario?.results?.CSRSResult?.finalAnnuity ?? null;
  $: ss = scenario?.results?.SocialSecurityResult?.ClaimingAmount ?? null;
  $: taxRate = scenario?.results?.EffectiveTaxRate ?? null;

</script>

<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mt-2">
  <div class="bg-white dark:bg-blue-900/30 rounded-lg shadow p-4 flex flex-col items-start border border-blue-100 dark:border-blue-800">
    <div class="text-xs text-blue-700 dark:text-blue-200 font-semibold">Net After-Tax Income</div>
    <div class="text-xl font-bold text-blue-900 dark:text-blue-100">{fmtCurrency(netIncome)}</div>
  </div>
  <div class="bg-white dark:bg-blue-900/30 rounded-lg shadow p-4 flex flex-col items-start border border-blue-100 dark:border-blue-800">
    <div class="text-xs text-blue-700 dark:text-blue-200 font-semibold">Projected TSP Balance</div>
    <div class="text-xl font-bold text-blue-900 dark:text-blue-100">{fmtCurrency(tspBalance)}</div>
  </div>
  <div class="bg-white dark:bg-blue-900/30 rounded-lg shadow p-4 flex flex-col items-start border border-blue-100 dark:border-blue-800">
    <div class="text-xs text-blue-700 dark:text-blue-200 font-semibold">Estimated Pension</div>
    <div class="text-xl font-bold text-blue-900 dark:text-blue-100">{fmtCurrency(pension)}</div>
  </div>
  <div class="bg-white dark:bg-blue-900/30 rounded-lg shadow p-4 flex flex-col items-start border border-blue-100 dark:border-blue-800">
    <div class="text-xs text-blue-700 dark:text-blue-200 font-semibold">Estimated Social Security</div>
    <div class="text-xl font-bold text-blue-900 dark:text-blue-100">{fmtCurrency(ss)}</div>
  </div>
</div>
<div class="mt-2 text-xs text-gray-500 dark:text-gray-300">
  <span>Effective Tax Rate: </span>
  <span>{taxRate != null ? (taxRate * 100).toFixed(1) + '%' : 'N/A'}</span>
</div>
