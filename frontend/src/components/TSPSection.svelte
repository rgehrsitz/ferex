<script lang="ts">
  export let data: any = {
    traditionalBalance: 500000,
    rothBalance: 100000,
    contributionRate: 5,
    contributionRateRoth: 5,
    expectedReturn: 7,
    withdrawalRate: 4,
    withdrawalStartAge: 62,
    fixedMonthlyWithdrawal: 0, // If user wants fixed amount instead of percentage
    withdrawalStrategy: 'percentage', // 'percentage', 'fixed', 'rmd'
  };

  const withdrawalStrategies = [
    { value: 'percentage', label: 'Percentage of Balance' },
    { value: 'fixed', label: 'Fixed Monthly Amount' },
    { value: 'rmd', label: 'Required Minimum Distribution' }
  ];

  function calculateAnnualWithdrawal() {
    const totalBalance = data.traditionalBalance + data.rothBalance;
    if (data.withdrawalStrategy === 'percentage') {
      return totalBalance * (data.withdrawalRate / 100);
    } else if (data.withdrawalStrategy === 'fixed') {
      return data.fixedMonthlyWithdrawal * 12;
    } else {
      // RMD simplified approximation (actual calculation depends on age)
      return totalBalance / 25; // Rough estimate, proper RMD uses life expectancy tables
    }
  }

  $: annualWithdrawal = calculateAnnualWithdrawal();
  $: monthlyWithdrawal = annualWithdrawal / 12;
  $: totalBalance = data.traditionalBalance + data.rothBalance;
</script>

<div>
  <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
    <div class="space-y-4">
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="traditionalBalance">
          TSP Traditional Balance
        </label>
        <div class="relative rounded-md shadow-sm">
          <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <span class="text-gray-500 dark:text-gray-400 sm:text-sm">$</span>
          </div>
          <input
            id="traditionalBalance"
            type="number"
            min="0"
            step="10000"
            class="w-full pl-7 px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
            bind:value={data.traditionalBalance}
          />
        </div>
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="rothBalance">
          TSP Roth Balance
        </label>
        <div class="relative rounded-md shadow-sm">
          <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <span class="text-gray-500 dark:text-gray-400 sm:text-sm">$</span>
          </div>
          <input
            id="rothBalance"
            type="number"
            min="0"
            step="10000"
            class="w-full pl-7 px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
            bind:value={data.rothBalance}
          />
        </div>
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="contributionRate">
          Current Traditional Contribution (% of salary)
        </label>
        <div class="relative rounded-md shadow-sm">
          <input
            id="contributionRate"
            type="number"
            min="0"
            max="100"
            class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
            bind:value={data.contributionRate}
          />
          <div class="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none">
            <span class="text-gray-500 dark:text-gray-400 sm:text-sm">%</span>
          </div>
        </div>
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="contributionRateRoth">
          Current Roth Contribution (% of salary)
        </label>
        <div class="relative rounded-md shadow-sm">
          <input
            id="contributionRateRoth"
            type="number"
            min="0"
            max="100"
            class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
            bind:value={data.contributionRateRoth}
          />
          <div class="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none">
            <span class="text-gray-500 dark:text-gray-400 sm:text-sm">%</span>
          </div>
        </div>
      </div>
    </div>

    <div class="space-y-4">
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="expectedReturn">
          Expected Annual Return
        </label>
        <div class="relative rounded-md shadow-sm">
          <input
            id="expectedReturn"
            type="number"
            min="0"
            max="20"
            step="0.1"
            class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
            bind:value={data.expectedReturn}
          />
          <div class="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none">
            <span class="text-gray-500 dark:text-gray-400 sm:text-sm">%</span>
          </div>
        </div>
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="withdrawalStartAge">
          Withdrawal Start Age
        </label>
        <input
          id="withdrawalStartAge"
          type="number"
          min="55"
          max="75"
          class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
          bind:value={data.withdrawalStartAge}
        />
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="withdrawalStrategy">
          Withdrawal Strategy
        </label>
        <select 
          id="withdrawalStrategy"
          class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
          bind:value={data.withdrawalStrategy}
        >
          {#each withdrawalStrategies as strategy}
            <option value={strategy.value}>{strategy.label}</option>
          {/each}
        </select>
      </div>

      {#if data.withdrawalStrategy === 'percentage'}
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="withdrawalRate">
            Annual Withdrawal Rate
          </label>
          <div class="relative rounded-md shadow-sm">
            <input
              id="withdrawalRate"
              type="number"
              min="0"
              max="20"
              step="0.1"
              class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
              bind:value={data.withdrawalRate}
            />
            <div class="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none">
              <span class="text-gray-500 dark:text-gray-400 sm:text-sm">%</span>
            </div>
          </div>
        </div>
      {:else if data.withdrawalStrategy === 'fixed'}
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="fixedMonthlyWithdrawal">
            Fixed Monthly Withdrawal
          </label>
          <div class="relative rounded-md shadow-sm">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <span class="text-gray-500 dark:text-gray-400 sm:text-sm">$</span>
            </div>
            <input
              id="fixedMonthlyWithdrawal"
              type="number"
              min="0"
              step="100"
              class="w-full pl-7 px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
              bind:value={data.fixedMonthlyWithdrawal}
            />
          </div>
        </div>
      {/if}
    </div>
  </div>

  <div class="mt-8 bg-gray-50 dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-lg p-4">
    <h3 class="text-lg font-medium text-gray-800 dark:text-gray-300 mb-3">TSP Summary</h3>
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <div class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
        <div class="text-sm text-gray-500 dark:text-gray-400">Total TSP Balance</div>
        <div class="text-2xl font-semibold text-gray-800 dark:text-gray-200">${totalBalance.toLocaleString()}</div>
      </div>
      <div class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
        <div class="text-sm text-gray-500 dark:text-gray-400">Annual Withdrawal</div>
        <div class="text-2xl font-semibold text-gray-800 dark:text-gray-200">${annualWithdrawal.toLocaleString(undefined, { maximumFractionDigits: 0 })}</div>
      </div>
      <div class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
        <div class="text-sm text-gray-500 dark:text-gray-400">Monthly Withdrawal</div>
        <div class="text-2xl font-semibold text-gray-800 dark:text-gray-200">${monthlyWithdrawal.toLocaleString(undefined, { maximumFractionDigits: 0 })}</div>
      </div>
    </div>
  </div>
</div>