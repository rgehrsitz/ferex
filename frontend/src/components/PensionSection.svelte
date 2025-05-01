<script lang="ts">
  export let data: any = {
    retirementSystem: 'FERS',
    highThreeSalary: 100000,
    yearsOfService: 30,
    retirementAge: 62,
    unusedSickLeave: 1000, // in hours
    militaryService: 0, // in years
    partTimeService: 0, // in years
    csrsOffset: false,
    survivorBenefit: 'full', // 'none', 'partial', 'full'
  };

  const retirementSystems = [
    { value: 'FERS', label: 'FERS' },
    { value: 'CSRS', label: 'CSRS' },
    { value: 'CSRS_OFFSET', label: 'CSRS Offset' }
  ];

  const survivorBenefitOptions = [
    { value: 'none', label: 'None (0%)' },
    { value: 'partial', label: 'Partial (25%)' },
    { value: 'full', label: 'Full (50%)' }
  ];

  function calculatePension() {
    // This will be implemented with actual calculation logic
    return data.highThreeSalary * 0.01 * data.yearsOfService;
  }

  $: annualPension = calculatePension();
  $: monthlyPension = annualPension / 12;
</script>

<div>
  <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
    <div class="space-y-4">
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="retirementSystem">
          Retirement System
        </label>
        <select 
          id="retirementSystem"
          class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
          bind:value={data.retirementSystem}
        >
          {#each retirementSystems as system}
            <option value={system.value}>{system.label}</option>
          {/each}
        </select>
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="highThreeSalary">
          High-3 Average Salary
        </label>
        <div class="relative rounded-md shadow-sm">
          <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <span class="text-gray-500 dark:text-gray-400 sm:text-sm">$</span>
          </div>
          <input
            id="highThreeSalary"
            type="number"
            min="0"
            step="1000"
            class="w-full pl-7 px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
            bind:value={data.highThreeSalary}
          />
        </div>
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="yearsOfService">
          Years of Service
        </label>
        <input
          id="yearsOfService"
          type="number"
          min="0"
          step="0.5"
          class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
          bind:value={data.yearsOfService}
        />
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="retirementAge">
          Age at Retirement
        </label>
        <input
          id="retirementAge"
          type="number"
          min="55"
          max="75"
          class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
          bind:value={data.retirementAge}
        />
      </div>
    </div>

    <div class="space-y-4">
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="unusedSickLeave">
          Unused Sick Leave (hours)
        </label>
        <input
          id="unusedSickLeave"
          type="number"
          min="0"
          class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
          bind:value={data.unusedSickLeave}
        />
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="militaryService">
          Military Service (years, with deposit paid)
        </label>
        <input
          id="militaryService"
          type="number"
          min="0"
          step="0.5"
          class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
          bind:value={data.militaryService}
        />
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="survivorBenefit">
          Survivor Benefit Election
        </label>
        <select 
          id="survivorBenefit"
          class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
          bind:value={data.survivorBenefit}
        >
          {#each survivorBenefitOptions as option}
            <option value={option.value}>{option.label}</option>
          {/each}
        </select>
      </div>

      {#if data.retirementSystem === 'CSRS_OFFSET'}
        <div class="flex items-center">
          <input
            id="csrsOffset"
            type="checkbox"
            class="h-4 w-4 text-primary-600 border-gray-300 rounded focus:ring-primary-500"
            bind:checked={data.csrsOffset}
          />
          <label for="csrsOffset" class="ml-2 block text-sm text-gray-700 dark:text-gray-300">
            Apply CSRS Offset at age 62
          </label>
        </div>
      {/if}
    </div>
  </div>

  <div class="mt-8 bg-gray-50 dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-lg p-4">
    <h3 class="text-lg font-medium text-gray-800 dark:text-gray-300 mb-3">Estimated Pension</h3>
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
        <div class="text-sm text-gray-500 dark:text-gray-400">Annual</div>
        <div class="text-2xl font-semibold text-gray-800 dark:text-gray-200">${annualPension.toLocaleString()}</div>
      </div>
      <div class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
        <div class="text-sm text-gray-500 dark:text-gray-400">Monthly</div>
        <div class="text-2xl font-semibold text-gray-800 dark:text-gray-200">${monthlyPension.toLocaleString(undefined, { maximumFractionDigits: 2 })}</div>
      </div>
    </div>
  </div>
</div>