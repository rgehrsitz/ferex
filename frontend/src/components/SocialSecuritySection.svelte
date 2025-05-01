<script lang="ts">
  export let data: any = {
    monthlyBenefit: 2000,
    startAge: 62,
    eligibleForWEP: false,
    yearsSubstantialEarnings: 20,
    spouseBenefit: 0,
  };

  const ageOptions = [
    { value: 62, label: 'Age 62 (reduced)' },
    { value: 67, label: 'Age 67 (full retirement age)' },
    { value: 70, label: 'Age 70 (maximum benefit)' }
  ];

  function calculateAnnualBenefit() {
    return data.monthlyBenefit * 12;
  }

  $: annualBenefit = calculateAnnualBenefit();
</script>

<div>
  <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
    <div class="space-y-4">
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="monthlyBenefit">
          Monthly Social Security Benefit
        </label>
        <div class="relative rounded-md shadow-sm">
          <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <span class="text-gray-500 dark:text-gray-400 sm:text-sm">$</span>
          </div>
          <input
            id="monthlyBenefit"
            type="number"
            min="0"
            step="50"
            class="w-full pl-7 px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
            bind:value={data.monthlyBenefit}
          />
        </div>
        <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
          Enter the amount from your Social Security statement
        </p>
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
    </div>

    <div class="space-y-4">
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="spouseBenefit">
          Spouse Monthly Benefit (if applicable)
        </label>
        <div class="relative rounded-md shadow-sm">
          <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <span class="text-gray-500 dark:text-gray-400 sm:text-sm">$</span>
          </div>
          <input
            id="spouseBenefit"
            type="number"
            min="0"
            step="50"
            class="w-full pl-7 px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
            bind:value={data.spouseBenefit}
          />
        </div>
      </div>

      <div class="mt-1">
        <div class="flex items-center">
          <input
            id="eligibleForWEP"
            type="checkbox"
            class="h-4 w-4 text-primary-600 border-gray-300 rounded focus:ring-primary-500"
            bind:checked={data.eligibleForWEP}
          />
          <label for="eligibleForWEP" class="ml-2 block text-sm text-gray-700 dark:text-gray-300">
            Subject to Windfall Elimination Provision (WEP)
          </label>
        </div>
        <p class="mt-1 text-xs text-gray-500 dark:text-gray-400 pl-6">
          Note: WEP was repealed for benefits after January 2024
        </p>
      </div>

      {#if data.eligibleForWEP}
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="yearsSubstantialEarnings">
            Years of Substantial Earnings
          </label>
          <input
            id="yearsSubstantialEarnings"
            type="number"
            min="0"
            max="35"
            class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
            bind:value={data.yearsSubstantialEarnings}
          />
        </div>
      {/if}
    </div>
  </div>

  <div class="mt-8 bg-gray-50 dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-lg p-4">
    <h3 class="text-lg font-medium text-gray-800 dark:text-gray-300 mb-3">Social Security Summary</h3>
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
        <div class="text-sm text-gray-500 dark:text-gray-400">Annual Benefit</div>
        <div class="text-2xl font-semibold text-gray-800 dark:text-gray-200">${annualBenefit.toLocaleString()}</div>
      </div>
      <div class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
        <div class="text-sm text-gray-500 dark:text-gray-400">Benefit Start Age</div>
        <div class="text-2xl font-semibold text-gray-800 dark:text-gray-200">{data.startAge}</div>
      </div>
    </div>
  </div>
</div>