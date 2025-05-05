<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { CalculatePension } from '../../wailsjs/go/main/App.js';
  import { main } from '../../wailsjs/go/models.js';
  import type { PensionData } from '../types/scenario.js';

  const dispatch = createEventDispatcher();
  export let data: PensionData;

  // Ensure data is initialized with defaults
  if (!data) {
    data = {};
  }
  
  // Set defaults for any missing fields
  data.system = data.system || 'FERS';
  data.high3Salary = data.high3Salary || 100000;
  data.yearsOfService = data.yearsOfService || 30;
  data.ageAtRetirement = data.ageAtRetirement || 62;
  data.unusedSickLeaveMonths = data.unusedSickLeaveMonths || 6;
  data.isPartTime = data.isPartTime || false;
  data.partTimeProrationFactor = data.partTimeProrationFactor || 1.0;
  data.survivorBenefitOption = data.survivorBenefitOption || 'full';
  data.csrsOffset = data.csrsOffset || false;

  // For UI display, map from backend field names to frontend-friendly names where needed
  let retirementSystem = data.system;
  let highThreeSalary = data.high3Salary;
  let retirementAge = data.ageAtRetirement;
  let unusedSickLeaveMonths = data.unusedSickLeaveMonths;
  let survivorBenefitOption = data.survivorBenefitOption;

  // Add custom field for display
  let csrsOffset = data.csrsOffset;

  // Options for dropdowns
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

  // Update retirementSystem when csrsOffset changes
  $: {
    // If we're in CSRS_OFFSET mode, update the system when offset changes
    if (retirementSystem === 'CSRS_OFFSET') {
      // This is a UI-only value, doesn't directly map to backend
      data.csrsOffset = csrsOffset;
    }
  }
  
  // Update the data object whenever UI fields change
  $: {
    data.system = retirementSystem;
    data.high3Salary = highThreeSalary;
    data.ageAtRetirement = retirementAge;
    data.unusedSickLeaveMonths = unusedSickLeaveMonths;
    data.survivorBenefitOption = survivorBenefitOption;
    
    console.log('PensionSection updating data, pension data now:', data);
    
    // Notify parent component of changes
    dispatch('update', data);
  }

  // Calculate pension when data changes
  async function calculatePension() {
    if (!data) return;
    
    try {
      loading = true;
      error = '';
      
      // Map the frontend retirement system to the backend format
      let systemValue = data.system;
      if (systemValue === 'CSRS_OFFSET') {
        systemValue = 'CSRS Offset';
      }
      
      // Create input object for backend
      const pensionInput = new main.PensionInput();
      pensionInput.system = systemValue;
      pensionInput.high3Salary = data.high3Salary;
      pensionInput.yearsOfService = data.yearsOfService;
      pensionInput.ageAtRetirement = data.ageAtRetirement;
      pensionInput.unusedSickLeaveMonths = data.unusedSickLeaveMonths;
      pensionInput.survivorBenefitOption = data.survivorBenefitOption;
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
  
  // Debounce the pension calculation
  let calculationTimer;
  $: {
    // Reset the timer whenever any data changes
    clearTimeout(calculationTimer);
    calculationTimer = setTimeout(() => {
      calculatePension();
    }, 300);
  }

  $: annualPension = calculationResult.annualPension;
  $: monthlyPension = calculationResult.monthlyPension;
</script>

<!-- Capture input events to bubble across shadow boundaries -->
<div>
  <!-- Required fields legend -->
  <div class="mb-4 text-sm text-gray-600 dark:text-gray-400 flex items-center">
    <span class="text-red-500 mr-1">*</span> Required fields for accurate calculations
  </div>
  
  <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
    <div class="space-y-4">
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="retirementSystem">
          Retirement System <span class="text-red-500">*</span>
        </label>
        <select 
          id="retirementSystem"
          class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
          bind:value={retirementSystem}
        >
          {#each retirementSystems as system}
            <option value={system.value}>{system.label}</option>
          {/each}
        </select>
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="highThreeSalary">
          High-3 Average Salary <span class="text-red-500">*</span>
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
            bind:value={highThreeSalary}
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
          bind:value={retirementAge}
        />
      </div>
    </div>

    <div class="space-y-4">
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="unusedSickLeave">
          Unused Sick Leave (months)
        </label>
        <input
          id="unusedSickLeave"
          type="number"
          min="0"
          class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
          bind:value={unusedSickLeaveMonths}
        />
        <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
          Note: 1 month ≈ 174 hours of sick leave
        </p>
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
          bind:value={survivorBenefitOption}
        >
          {#each survivorBenefitOptions as option}
            <option value={option.value}>{option.label}</option>
          {/each}
        </select>
      </div>

      {#if retirementSystem === 'CSRS_OFFSET'}
        <div class="flex items-center">
          <input
            id="csrsOffset"
            type="checkbox"
            class="h-4 w-4 text-primary-600 border-gray-300 rounded focus:ring-primary-500"
            bind:checked={csrsOffset}
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