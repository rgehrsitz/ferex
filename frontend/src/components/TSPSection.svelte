<script lang="ts">
  import { CalculateTSPProjection } from '../../wailsjs/go/main/App';
  import { main } from '../../wailsjs/go/models';
  import type { TSPData } from '../types/scenario';
  
  export let data: TSPData;
  
  let loading = false;
  let error = '';
  let projectionResult: any = {
    yearlyData: [],
    maxBalance: 0,
    finalBalance: 0,
    totalContributions: 0,
    totalWithdrawals: 0,
    totalReturns: 0,
    notes: ""
  };
  
  // Defaults for visualization
  let projectionYears = 5; // Show first 5 years by default

  const withdrawalStrategies = [
    { value: 'fixed', label: 'Fixed Monthly Amount' },
    { value: 'percentage', label: 'Percentage of Balance' },
    { value: 'rmd', label: 'Required Minimum Distribution' }
  ];
  
  async function calculateTSPProjection() {
    try {
      loading = true;
      error = '';
      
      // Calculate annual contribution from percentages (estimate)
      const estimatedSalary = 100000; // Default estimate for annual salary
      const annualContribution = (
        (data.contributionRate / 100) * estimatedSalary + 
        (data.contributionRateRoth / 100) * estimatedSalary
      );
      
      // Create input object for backend
      const tspInput = new main.TSPInput();
      tspInput.currentBalance = data.traditionalBalance + data.rothBalance;
      tspInput.traditionalBalance = data.traditionalBalance;
      tspInput.rothBalance = data.rothBalance;
      tspInput.annualContribution = annualContribution;
      tspInput.expectedReturnRate = data.expectedReturn / 100;
      tspInput.withdrawalStrategy = data.withdrawalStrategy;
      tspInput.fixedWithdrawalAmount = data.fixedMonthlyWithdrawal * 12; // Convert to annual
      tspInput.withdrawalPercentage = data.withdrawalRate / 100;
      tspInput.withdrawalStartAge = data.withdrawalStartAge;
      tspInput.birthYear = new Date().getFullYear() - 40; // Default assumption
      tspInput.retirementAge = data.withdrawalStartAge;
      
      // Call backend API
      projectionResult = await CalculateTSPProjection(tspInput);
      
      return {
        annualWithdrawal: getAnnualWithdrawal(),
        totalBalance: data.traditionalBalance + data.rothBalance
      };
    } catch (err) {
      console.error("Error calculating TSP projection:", err);
      error = 'Failed to calculate TSP projection. Please try again.';
      return {
        annualWithdrawal: 0,
        totalBalance: data.traditionalBalance + data.rothBalance
      };
    } finally {
      loading = false;
    }
  }
  
  function getAnnualWithdrawal() {
    // Check for valid projection data first
    if (projectionResult && projectionResult.yearlyData && projectionResult.yearlyData.length > 0) {
      // Find the first year with withdrawals
      for (const year of projectionResult.yearlyData) {
        if (year.withdrawals > 0) {
          return year.withdrawals;
        }
      }
      
      // If no withdrawals found but we have data, find first year after withdrawal start age
      const withdrawalStartYear = projectionResult.yearlyData.find(year => year.age >= data.withdrawalStartAge);
      if (withdrawalStartYear) {
        // Calculate based on strategy
        if (data.withdrawalStrategy === 'percentage') {
          return withdrawalStartYear.startingBalance * (data.withdrawalRate / 100);
        } else if (data.withdrawalStrategy === 'fixed') {
          return data.fixedMonthlyWithdrawal * 12;
        } else {
          // RMD simplified approximation
          return withdrawalStartYear.startingBalance / 25;
        }
      }
    }
    
    // Fallback to simple calculation if no projection data
    const totalBalance = data.traditionalBalance + data.rothBalance || 0;
    
    if (data.withdrawalStrategy === 'percentage' && data.withdrawalRate) {
      return totalBalance * (data.withdrawalRate / 100);
    } else if (data.withdrawalStrategy === 'fixed' && data.fixedMonthlyWithdrawal) {
      return data.fixedMonthlyWithdrawal * 12;
    } else {
      // RMD simplified approximation
      return totalBalance / 25;
    }
  }

  // Trigger calculation when relevant inputs change
  $: {
    if (data.traditionalBalance !== undefined || 
        data.rothBalance !== undefined || 
        data.withdrawalStrategy || 
        data.withdrawalRate || 
        data.withdrawalStartAge ||
        data.fixedMonthlyWithdrawal ||
        data.expectedReturn) {
      calculateTSPProjection();
    }
  }

  $: annualWithdrawal = getAnnualWithdrawal() || 0;
  $: monthlyWithdrawal = isNaN(annualWithdrawal) ? 0 : annualWithdrawal / 12;
  $: totalBalance = (data.traditionalBalance || 0) + (data.rothBalance || 0);
  $: yearsUntilDepletion = projectionResult.notes?.match(/Balance depleted at age (\d+)/) 
                            ? parseInt(projectionResult.notes.match(/Balance depleted at age (\d+)/)[1]) - data.withdrawalStartAge
                            : '30+';
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
    
    {#if error}
      <div class="bg-red-100 dark:bg-red-900 border border-red-200 dark:border-red-700 text-red-700 dark:text-red-300 p-4 rounded-lg mb-4">
        {error}
      </div>
    {/if}
    
    {#if loading}
      <div class="flex justify-center items-center py-8">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600 dark:border-primary-400"></div>
      </div>
    {:else}
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
        <div class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
          <div class="text-sm text-gray-500 dark:text-gray-400">Current TSP Balance</div>
          <div class="text-2xl font-semibold text-gray-800 dark:text-gray-200">
            ${totalBalance.toLocaleString(undefined, { maximumFractionDigits: 0 })}
          </div>
        </div>
        <div class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
          <div class="text-sm text-gray-500 dark:text-gray-400">Annual Withdrawal</div>
          <div class="text-2xl font-semibold text-gray-800 dark:text-gray-200">
            ${(isNaN(annualWithdrawal) ? 0 : annualWithdrawal).toLocaleString(undefined, { maximumFractionDigits: 0 })}
          </div>
        </div>
        <div class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
          <div class="text-sm text-gray-500 dark:text-gray-400">Monthly Withdrawal</div>
          <div class="text-2xl font-semibold text-gray-800 dark:text-gray-200">
            ${(isNaN(monthlyWithdrawal) ? 0 : monthlyWithdrawal).toLocaleString(undefined, { maximumFractionDigits: 0 })}
          </div>
        </div>
      </div>
      
      {#if projectionResult && projectionResult.yearlyData && projectionResult.yearlyData.length > 0}
        <div class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow mb-6">
          <h4 class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-4">TSP Projection Summary</h4>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <div class="text-xs text-gray-500 dark:text-gray-400">Years Until Depletion</div>
              <div class="text-base font-medium text-gray-800 dark:text-gray-200">
                {typeof yearsUntilDepletion === 'string' ? yearsUntilDepletion : `${yearsUntilDepletion} years`}
              </div>
            </div>
            <div>
              <div class="text-xs text-gray-500 dark:text-gray-400">Maximum Balance</div>
              <div class="text-base font-medium text-gray-800 dark:text-gray-200">
                ${projectionResult.maxBalance.toLocaleString(undefined, { maximumFractionDigits: 0 })}
              </div>
            </div>
            <div>
              <div class="text-xs text-gray-500 dark:text-gray-400">Total Future Contributions</div>
              <div class="text-base font-medium text-gray-800 dark:text-gray-200">
                ${projectionResult.totalContributions.toLocaleString(undefined, { maximumFractionDigits: 0 })}
              </div>
            </div>
            <div>
              <div class="text-xs text-gray-500 dark:text-gray-400">Total Expected Withdrawals</div>
              <div class="text-base font-medium text-gray-800 dark:text-gray-200">
                ${projectionResult.totalWithdrawals.toLocaleString(undefined, { maximumFractionDigits: 0 })}
              </div>
            </div>
          </div>
        </div>
        
        <!-- Year-by-Year Projection Table (first few years) -->
        <div class="overflow-x-auto">
          <table class="min-w-full bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700">
            <thead>
              <tr class="bg-gray-50 dark:bg-gray-700">
                <th class="py-2 px-3 border-b border-gray-200 dark:border-gray-700 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">Age</th>
                <th class="py-2 px-3 border-b border-gray-200 dark:border-gray-700 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">Year</th>
                <th class="py-2 px-3 border-b border-gray-200 dark:border-gray-700 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">Starting Balance</th>
                <th class="py-2 px-3 border-b border-gray-200 dark:border-gray-700 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">Contributions</th>
                <th class="py-2 px-3 border-b border-gray-200 dark:border-gray-700 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">Growth</th>
                <th class="py-2 px-3 border-b border-gray-200 dark:border-gray-700 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">Withdrawals</th>
                <th class="py-2 px-3 border-b border-gray-200 dark:border-gray-700 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">Ending Balance</th>
              </tr>
            </thead>
            <tbody>
              {#each projectionResult.yearlyData.slice(0, projectionYears) as yearData}
                <tr class="hover:bg-gray-50 dark:hover:bg-gray-700">
                  <td class="py-2 px-3 border-b border-gray-200 dark:border-gray-700 text-sm text-gray-800 dark:text-gray-200">{yearData.age}</td>
                  <td class="py-2 px-3 border-b border-gray-200 dark:border-gray-700 text-sm text-gray-800 dark:text-gray-200">{yearData.year}</td>
                  <td class="py-2 px-3 border-b border-gray-200 dark:border-gray-700 text-sm text-gray-800 dark:text-gray-200 text-right">
                    ${yearData.startingBalance.toLocaleString(undefined, { maximumFractionDigits: 0 })}
                  </td>
                  <td class="py-2 px-3 border-b border-gray-200 dark:border-gray-700 text-sm text-gray-800 dark:text-gray-200 text-right">
                    ${yearData.contributions.toLocaleString(undefined, { maximumFractionDigits: 0 })}
                  </td>
                  <td class="py-2 px-3 border-b border-gray-200 dark:border-gray-700 text-sm text-gray-800 dark:text-gray-200 text-right">
                    ${yearData.returns.toLocaleString(undefined, { maximumFractionDigits: 0 })}
                  </td>
                  <td class="py-2 px-3 border-b border-gray-200 dark:border-gray-700 text-sm text-gray-800 dark:text-gray-200 text-right">
                    ${yearData.withdrawals.toLocaleString(undefined, { maximumFractionDigits: 0 })}
                  </td>
                  <td class="py-2 px-3 border-b border-gray-200 dark:border-gray-700 text-sm text-gray-800 dark:text-gray-200 text-right">
                    ${yearData.endingBalance.toLocaleString(undefined, { maximumFractionDigits: 0 })}
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
        
        <div class="mt-4 flex justify-between items-center">
          <button 
            class="px-3 py-1 text-sm bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-gray-200 rounded hover:bg-gray-300 dark:hover:bg-gray-600"
            on:click={() => projectionYears = Math.max(5, projectionYears - 5)}>
            Show Less
          </button>
          <span class="text-sm text-gray-500 dark:text-gray-400">
            Showing {projectionYears} years of {projectionResult.yearlyData.length} total
          </span>
          <button 
            class="px-3 py-1 text-sm bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-gray-200 rounded hover:bg-gray-300 dark:hover:bg-gray-600"
            on:click={() => projectionYears = Math.min(projectionResult.yearlyData.length, projectionYears + 5)}>
            Show More
          </button>
        </div>
      {/if}
      
      {#if projectionResult && projectionResult.notes}
        <div class="mt-4 bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
          <h4 class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Notes</h4>
          <div class="text-sm text-gray-600 dark:text-gray-400">
            {projectionResult.notes}
          </div>
        </div>
      {/if}
    {/if}
  </div>
</div>