<script lang="ts">
  import { CalculatePension } from '../../wailsjs/go/main/App';
  import { main } from '../../wailsjs/go/models';
  import type { PensionData } from '../types/scenario';

  export let data: PensionData;

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

  let calculationResult = {
    annualPension: 0,
    monthlyPension: 0,
    notes: ''
  };
  
  let loading = false;
  let error = '';

  async function calculatePension() {
    try {
      loading = true;
      error = '';
      
      // Convert hours to months for sick leave (assuming 174 hours = 1 month)
      const sickLeaveMonths = Math.round(data.unusedSickLeave / 174);
      
      // Map the frontend retirement system to the backend format
      let systemValue = data.retirementSystem;
      if (systemValue === 'CSRS_OFFSET') {
        systemValue = 'CSRS Offset';
      }
      
      // Create input object for backend
      const pensionInput = new main.PensionInput();
      pensionInput.system = systemValue;
      pensionInput.high3Salary = data.highThreeSalary;
      pensionInput.yearsOfService = data.yearsOfService;
      pensionInput.ageAtRetirement = data.retirementAge;
      pensionInput.unusedSickLeaveMonths = sickLeaveMonths;
      pensionInput.survivorBenefitOption = data.survivorBenefit;
      pensionInput.isPartTime = data.isPartTime;
      pensionInput.partTimeProrationFactor = data.partTimeProrationFactor;
      
      // Call backend API
      calculationResult = await CalculatePension(pensionInput);
      
      return calculationResult.annualPension;
    } catch (err) {
      console.error("Error calculating pension:", err);
      error = 'Failed to calculate pension. Please try again.';
      return 0;
    } finally {
      loading = false;
    }
  }

  // Trigger calculation when relevant inputs change
  $: {
    if (data.retirementSystem || data.highThreeSalary || data.yearsOfService || 
        data.retirementAge || data.unusedSickLeave || data.survivorBenefit) {
      calculatePension();
    }
  }

  $: annualPension = calculationResult.annualPension;
  $: monthlyPension = calculationResult.monthlyPension;
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
        <div class="flex items-center mb-2">
          <input
            id="isPartTime"
            type="checkbox"
            class="h-4 w-4 text-primary-600 border-gray-300 rounded focus:ring-primary-500"
            bind:checked={data.isPartTime}
          />
          <label for="isPartTime" class="ml-2 block text-sm text-gray-700 dark:text-gray-300">
            Include Part-time Service
          </label>
        </div>
        
        {#if data.isPartTime}
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="partTimeProrationFactor">
            Part-time Proration Factor (0.1 to 1.0)
          </label>
          <input
            id="partTimeProrationFactor"
            type="number"
            min="0.1"
            max="1.0"
            step="0.01"
            class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
            bind:value={data.partTimeProrationFactor}
          />
          <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
            Example: 0.5 for half-time service
          </p>
        {/if}
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
      
      {#if calculationResult.notes}
        <div class="mt-4 bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
          <h4 class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Calculation Notes</h4>
          <div class="text-sm text-gray-600 dark:text-gray-400 whitespace-pre-line">
            {calculationResult.notes}
          </div>
        </div>
      {/if}
    {/if}
  </div>
</div>