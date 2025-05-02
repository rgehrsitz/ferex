<script lang="ts">
  import { CalculateSocialSecurity } from '../../wailsjs/go/main/App';
  import { main } from '../../wailsjs/go/models';
  import type { SocialSecurityData } from '../types/scenario';
  
  export let data: SocialSecurityData;
  
  let loading = false;
  let error = '';
  let calculationResult: any = {
    estimatedMonthlyAt62: 0,
    estimatedMonthlyAtFRA: 0,
    estimatedMonthlyAt70: 0,
    estimatedAnnualAt62: 0,
    estimatedAnnualAtFRA: 0,
    estimatedAnnualAt70: 0,
    claimingAge: 0,
    claimingMonthlyAmount: 0,
    claimingAnnualAmount: 0,
    fullRetirementAge: 0,
    notes: ''
  };

  const ageOptions = [
    { value: 62, label: 'Age 62 (reduced)' },
    { value: 67, label: 'Age 67 (full retirement age)' },
    { value: 70, label: 'Age 70 (maximum benefit)' }
  ];
  
  const currentYear = new Date().getFullYear();
  
  async function calculateSocialSecurity() {
    try {
      loading = true;
      error = '';
      
      // Create input object for backend
      const ssInput = new main.SocialSecurityInput();
      ssInput.startAge = data.startAge;
      ssInput.estimatedMonthlyBenefit = data.estimatedMonthlyBenefit;
      ssInput.isEligible = data.isEligible;
      ssInput.birthYear = data.birthYear;
      ssInput.birthMonth = data.birthMonth;
      ssInput.currentAge = currentYear - data.birthYear;
      
      // If SSA estimates are provided, use them
      if (data.ssaEstimateAt62) {
        ssInput.userProvidedEstimate62 = data.ssaEstimateAt62;
      }
      if (data.ssaEstimateAtFRA) {
        ssInput.userProvidedEstimateFRA = data.ssaEstimateAtFRA;
      }
      if (data.ssaEstimateAt70) {
        ssInput.userProvidedEstimate70 = data.ssaEstimateAt70;
      }
      
      // Call backend API
      calculationResult = await CalculateSocialSecurity(ssInput);
      
      // Update data with calculated results
      data.estimatedMonthlyBenefit = calculationResult.claimingMonthlyAmount;
      
      return calculationResult.claimingAnnualAmount;
    } catch (err) {
      console.error("Error calculating Social Security:", err);
      error = 'Failed to calculate Social Security benefit. Please try again.';
      return 0;
    } finally {
      loading = false;
    }
  }

  // Trigger calculation when relevant inputs change
  $: {
    if (data.startAge || data.estimatedMonthlyBenefit || data.birthYear || 
        data.ssaEstimateAt62 || data.ssaEstimateAtFRA || data.ssaEstimateAt70) {
      calculateSocialSecurity();
    }
  }

  $: annualBenefit = calculationResult.claimingAnnualAmount || (data.estimatedMonthlyBenefit * 12);
</script>

<div>
  <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
    <div class="space-y-4">
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="birthYear">
          Birth Year
        </label>
        <input
          id="birthYear"
          type="number"
          min="1930"
          max="2000"
          class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
          bind:value={data.birthYear}
        />
      </div>
      
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="birthMonth">
          Birth Month
        </label>
        <select 
          id="birthMonth"
          class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
          bind:value={data.birthMonth}
        >
          <option value="1">January</option>
          <option value="2">February</option>
          <option value="3">March</option>
          <option value="4">April</option>
          <option value="5">May</option>
          <option value="6">June</option>
          <option value="7">July</option>
          <option value="8">August</option>
          <option value="9">September</option>
          <option value="10">October</option>
          <option value="11">November</option>
          <option value="12">December</option>
        </select>
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="startAge">
          Benefit Start Age
        </label>
        <select 
          id="startAge"
          class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
          bind:value={data.startAge}
        >
          {#each ageOptions as option}
            <option value={option.value}>{option.label}</option>
          {/each}
        </select>
      </div>
      
      <div class="flex items-center">
        <input
          id="isEligible"
          type="checkbox"
          class="h-4 w-4 text-primary-600 border-gray-300 rounded focus:ring-primary-500"
          bind:checked={data.isEligible}
        />
        <label for="isEligible" class="ml-2 block text-sm text-gray-700 dark:text-gray-300">
          Eligible for Social Security
        </label>
      </div>
    </div>

    <div class="space-y-4">
      <h4 class="font-medium text-gray-700 dark:text-gray-300">SSA Benefit Estimates</h4>
      <p class="text-xs text-gray-500 dark:text-gray-400">
        Enter values from your Social Security statement for the most accurate projection
      </p>
      
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="ssaEstimateAt62">
          Age 62 Estimate
        </label>
        <div class="relative rounded-md shadow-sm">
          <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <span class="text-gray-500 dark:text-gray-400 sm:text-sm">$</span>
          </div>
          <input
            id="ssaEstimateAt62"
            type="number"
            min="0"
            step="50"
            class="w-full pl-7 px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
            bind:value={data.ssaEstimateAt62}
          />
        </div>
      </div>
      
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="ssaEstimateAtFRA">
          Full Retirement Age Estimate
        </label>
        <div class="relative rounded-md shadow-sm">
          <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <span class="text-gray-500 dark:text-gray-400 sm:text-sm">$</span>
          </div>
          <input
            id="ssaEstimateAtFRA"
            type="number"
            min="0"
            step="50"
            class="w-full pl-7 px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
            bind:value={data.ssaEstimateAtFRA}
          />
        </div>
      </div>
      
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="ssaEstimateAt70">
          Age 70 Estimate
        </label>
        <div class="relative rounded-md shadow-sm">
          <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <span class="text-gray-500 dark:text-gray-400 sm:text-sm">$</span>
          </div>
          <input
            id="ssaEstimateAt70"
            type="number"
            min="0"
            step="50"
            class="w-full pl-7 px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
            bind:value={data.ssaEstimateAt70}
          />
        </div>
      </div>
    </div>
  </div>

  <div class="mt-8 bg-gray-50 dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-lg p-4">
    <h3 class="text-lg font-medium text-gray-800 dark:text-gray-300 mb-3">Social Security Summary</h3>
    
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
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
        <div class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
          <div class="text-sm text-gray-500 dark:text-gray-400">Annual Benefit at Age {data.startAge}</div>
          <div class="text-2xl font-semibold text-gray-800 dark:text-gray-200">
            ${annualBenefit.toLocaleString(undefined, { maximumFractionDigits: 0 })}
          </div>
        </div>
        <div class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
          <div class="text-sm text-gray-500 dark:text-gray-400">Monthly Benefit at Age {data.startAge}</div>
          <div class="text-2xl font-semibold text-gray-800 dark:text-gray-200">
            ${(annualBenefit/12).toLocaleString(undefined, { maximumFractionDigits: 0 })}
          </div>
        </div>
      </div>
      
      {#if calculationResult.fullRetirementAge}
        <div class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow mb-4">
          <h4 class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Benefit Estimates by Claiming Age</h4>
          <div class="grid grid-cols-3 gap-4">
            <div>
              <div class="text-xs text-gray-500 dark:text-gray-400">Age 62</div>
              <div class="text-sm font-medium text-gray-800 dark:text-gray-200">
                ${calculationResult.estimatedMonthlyAt62.toLocaleString(undefined, { maximumFractionDigits: 0 })}/mo
              </div>
            </div>
            <div>
              <div class="text-xs text-gray-500 dark:text-gray-400">Full Age ({calculationResult.fullRetirementAge})</div>
              <div class="text-sm font-medium text-gray-800 dark:text-gray-200">
                ${calculationResult.estimatedMonthlyAtFRA.toLocaleString(undefined, { maximumFractionDigits: 0 })}/mo
              </div>
            </div>
            <div>
              <div class="text-xs text-gray-500 dark:text-gray-400">Age 70</div>
              <div class="text-sm font-medium text-gray-800 dark:text-gray-200">
                ${calculationResult.estimatedMonthlyAt70.toLocaleString(undefined, { maximumFractionDigits: 0 })}/mo
              </div>
            </div>
          </div>
        </div>
      {/if}
      
      {#if calculationResult.notes}
        <div class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
          <h4 class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Notes</h4>
          <div class="text-sm text-gray-600 dark:text-gray-400">
            {calculationResult.notes}
          </div>
        </div>
      {/if}
    {/if}
  </div>
</div>