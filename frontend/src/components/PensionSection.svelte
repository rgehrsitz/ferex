<script lang="ts">
  import { CalculatePension } from '../../wailsjs/go/main/App.js';
  import { main } from '../../wailsjs/go/models.js';
  import type { PensionData } from '../types/scenario.js';
  import SectionHeader from './SectionHeader.svelte';

  // Default data object
  const defaultData = {
    system: 'FERS',
    high3Salary: 100000,
    yearsOfService: 30,
    ageAtRetirement: 62,
    unusedSickLeaveMonths: 6,
    isPartTime: false,
    partTimeProrationFactor: 1.0,
    survivorBenefitOption: 'full',
    csrsOffset: false,
    militaryService: 0
  };

  // Use defaulted props pattern
  const { 
    data: propData = defaultData, 
    scenarioName, 
    onUpdate = (data: any) => {} 
  } = $props<{
    data?: PensionData;
    scenarioName: string;
    onUpdate?: (data: any) => void;
  }>();
  
  // Create a local copy we can work with
  const inputData = propData ? { ...propData } : { ...defaultData };

  // Set defaults for any missing fields if we have data from props
  if (propData) {
    inputData.system = propData.system || defaultData.system;
    inputData.high3Salary = propData.high3Salary || defaultData.high3Salary;
    inputData.yearsOfService = propData.yearsOfService || defaultData.yearsOfService;
    inputData.ageAtRetirement = propData.ageAtRetirement || defaultData.ageAtRetirement;
    inputData.unusedSickLeaveMonths = propData.unusedSickLeaveMonths || defaultData.unusedSickLeaveMonths;
    inputData.isPartTime = propData.isPartTime ?? defaultData.isPartTime;
    inputData.partTimeProrationFactor = propData.partTimeProrationFactor || defaultData.partTimeProrationFactor;
    inputData.survivorBenefitOption = propData.survivorBenefitOption || defaultData.survivorBenefitOption;
    inputData.csrsOffset = propData.csrsOffset ?? defaultData.csrsOffset;
    inputData.militaryService = propData.militaryService || defaultData.militaryService;
  }

  // State with Svelte 5 runes
  let retirementSystem = $state(inputData.system);
  let highThreeSalary = $state(inputData.high3Salary);
  let yearsOfService = $state(inputData.yearsOfService);
  let retirementAge = $state(inputData.ageAtRetirement);
  let unusedSickLeaveMonths = $state(inputData.unusedSickLeaveMonths);
  let survivorBenefitOption = $state(inputData.survivorBenefitOption);
  let isPartTime = $state(inputData.isPartTime);
  let partTimeProrationFactor = $state(inputData.partTimeProrationFactor);
  let csrsOffset = $state(inputData.csrsOffset);
  let militaryService = $state(inputData.militaryService);
  
  let calculationResult = $state({
    annualPension: 0,
    monthlyPension: 0,
    notes: ''
  });
  
  let loading = $state(false);
  let error = $state('');

  // Sync changes to a new data object and notify parent
  $effect(() => {
    // Create a new data object with current values
    const updatedData = {
      system: retirementSystem,
      high3Salary: highThreeSalary,
      yearsOfService: yearsOfService,
      ageAtRetirement: retirementAge,
      unusedSickLeaveMonths: unusedSickLeaveMonths,
      survivorBenefitOption: survivorBenefitOption,
      isPartTime: isPartTime,
      partTimeProrationFactor: partTimeProrationFactor,
      militaryService: militaryService,
      csrsOffset: retirementSystem === 'CSRS_OFFSET' ? csrsOffset : false
    };
    
    // Notify parent component
    onUpdate(updatedData);
  });

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

  // Derived values from calculation result
  const annualPension = $derived(calculationResult.annualPension);
  const monthlyPension = $derived(calculationResult.monthlyPension);

  // Handle all field changes in a single function
  function handleFieldChange() {
    // Calculate pension with updated values
    calculatePension();
  }

  // Calculate pension when data changes
  async function calculatePension() {
    try {
      loading = true;
      error = '';
      
      // Map the frontend retirement system to the backend format
      let systemValue = retirementSystem;
      if (systemValue === 'CSRS_OFFSET') {
        systemValue = 'CSRS Offset';
      }
      
      // Create input object for backend
      const pensionInput = new main.PensionInput();
      pensionInput.system = systemValue ?? '';
      pensionInput.high3Salary = highThreeSalary ?? 0;
      pensionInput.yearsOfService = yearsOfService ?? 0;
      pensionInput.ageAtRetirement = retirementAge ?? 0;
      pensionInput.unusedSickLeaveMonths = unusedSickLeaveMonths ?? 0;
      pensionInput.survivorBenefitOption = survivorBenefitOption ?? '';
      pensionInput.isPartTime = isPartTime ?? false;
      pensionInput.partTimeProrationFactor = partTimeProrationFactor ?? 1;
      pensionInput.militaryService = militaryService ?? 0;
      
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
</script>

<!-- Capture input events to bubble across shadow boundaries -->
<div>
  <SectionHeader sectionName="Pension Calculator" {scenarioName} />
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
          onchange={handleFieldChange}
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
            onchange={handleFieldChange}
            oninput={() => {
              // For number inputs, update on input too for more responsive feel
              if (typeof highThreeSalary === 'string') {
                highThreeSalary = parseFloat(highThreeSalary);
              }
              handleFieldChange();
            }}
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
          bind:value={yearsOfService}
          onchange={handleFieldChange}
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
          onchange={handleFieldChange}
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
          onchange={handleFieldChange}
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
          bind:value={militaryService}
          onchange={handleFieldChange}
        />
      </div>
      
      <div>
        <div class="flex items-center mb-2">
          <input
            id="isPartTime"
            type="checkbox"
            class="h-4 w-4 text-primary-600 border-gray-300 rounded focus:ring-primary-500"
            bind:checked={isPartTime}
            onchange={handleFieldChange}
          />
          <label for="isPartTime" class="ml-2 block text-sm text-gray-700 dark:text-gray-300">
            Include Part-time Service
          </label>
        </div>
        
        {#if isPartTime}
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
            bind:value={partTimeProrationFactor}
            onchange={handleFieldChange}
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
          onchange={handleFieldChange}
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
            onchange={handleFieldChange}
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