<script lang="ts">
  import { main } from '../../wailsjs/go/models.js';
  import type { SocialSecurityData } from '../types/scenario.js';
  import { api } from '../stores/apiStore.js';
  import { storeCalculationResult, getCalculationResult } from '../stores/calculationStore.js';
  import { onMount } from 'svelte';
  import SectionHeader from './SectionHeader.svelte';
  
  export const onUpdate: (data: any) => void = () => {};
  
  export let data: SocialSecurityData;
  export let scenarioId: number;
  export let scenarioName: string;
  
  // Ensure data is initialized with defaults
  if (!data) {
    data = {
      startAge: 67,
      estimatedMonthlyBenefit: 2000,
      isEligible: true,
      birthYear: 1970,
      birthMonth: 1
    };
  }
  
  // Set defaults for any missing fields
  data.startAge = data.startAge ?? 67; // Default to full retirement age
  data.estimatedMonthlyBenefit = data.estimatedMonthlyBenefit ?? 2000;
  data.isEligible = data.isEligible ?? true;
  data.birthYear = data.birthYear ?? 1970;
  data.birthMonth = data.birthMonth ?? 1;
  
  // Make sure birthMonth is stored as a number
  if (typeof data.birthMonth === 'string') {
    data.birthMonth = parseInt(String(data.birthMonth), 10);
  }

  // Create local variables for UI binding
  let startAge = data.startAge;
  let estimatedMonthlyBenefit = data.estimatedMonthlyBenefit;
  let isEligible = data.isEligible;
  let birthYear = data.birthYear;
  let birthMonth = data.birthMonth;

  // Svelte 5 idiom: $: for reactivity
  $: data.startAge = startAge;
  $: data.estimatedMonthlyBenefit = estimatedMonthlyBenefit;
  $: data.isEligible = isEligible;
  $: data.birthYear = birthYear;
  $: data.birthMonth = birthMonth;
  $: onUpdate({ ...data });

  // Handle all field changes in a single function
  function handleFieldChange() {
    // Update data object directly
    data.startAge = startAge;
    data.estimatedMonthlyBenefit = estimatedMonthlyBenefit;
    data.isEligible = isEligible;
    data.birthYear = birthYear;
    data.birthMonth = birthMonth;
    
    console.log('SocialSecuritySection updating data:', data);
    
    
    onUpdate({ ...data });
    
    // Try to calculate on change
    checkAndCalculate();
  }
  
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
    claimingAmount: 0,
    fullRetirementAge: 0,
    notes: ''
  };
  
  // Default values at age 62
  let monthlyBenefit = 0;
  let annualBenefit = 0;

  const ageOptions = [
    { value: 62, label: 'Age 62 (reduced)' },
    { value: 67, label: 'Age 67 (full retirement age)' },
    { value: 70, label: 'Age 70 (maximum benefit)' }
  ];
  
  const currentYear = new Date().getFullYear();
  
  // Make the function accessible from outside
  export async function calculateSocialSecurity() {
    try {
      loading = true;
      error = '';
      
      // Validate required inputs
      if (!data.birthYear) {
        throw new Error("Birth year is required");
      }
      
      if (data.birthYear < 1900 || data.birthYear > currentYear - 20) {
        throw new Error(`Birth year must be between 1900 and ${currentYear - 20}`);
      }
      
      if (!data.startAge) {
        throw new Error("Benefit start age is required");
      }
      
      // Create input object for backend
      const ssInput = new main.SocialSecurityInput();
      // Ensure all numeric values are properly converted to numbers
      ssInput.startAge = parseInt(String(data.startAge || 62), 10);
      ssInput.estimatedMonthlyBenefit = parseFloat(String(data.estimatedMonthlyBenefit || 0));
      ssInput.isEligible = data.isEligible !== undefined ? data.isEligible : true;
      ssInput.birthYear = parseInt(String(data.birthYear || 1960), 10);
      ssInput.birthMonth = parseInt(String(data.birthMonth || 1), 10);
      ssInput.currentAge = currentYear - parseInt(String(data.birthYear || 1960), 10);
      
      // For estimation - default values if not available
      ssInput.estimatedAnnualSalary = 60000; // Default annual salary for estimation
      ssInput.yearsWorked = 35; // Default years worked for estimation
      
      // If SSA estimates are provided, use them - ensure they're numbers
      if (data.ssaEstimateAt62 && data.ssaEstimateAt62 > 0) {
        ssInput.userProvidedEstimate62 = parseFloat(String(data.ssaEstimateAt62));
      }
      if (data.ssaEstimateAtFRA && data.ssaEstimateAtFRA > 0) {
        ssInput.userProvidedEstimateFRA = parseFloat(String(data.ssaEstimateAtFRA));
      }
      if (data.ssaEstimateAt70 && data.ssaEstimateAt70 > 0) {
        ssInput.userProvidedEstimate70 = parseFloat(String(data.ssaEstimateAt70));
      }
      
      // Call backend API via the store
      calculationResult = await api.calculateSocialSecurity(ssInput);
      
      // Store the calculation result for potential reuse by other components
      storeCalculationResult(scenarioId, 'socialSecurity', calculationResult);
      
      console.log("Social Security calculation result:", JSON.stringify(calculationResult, null, 2));
      
      // Validate result
      if (!calculationResult) {
        throw new Error("No calculation result returned from server");
      }
      
      // Instead of updating data directly, we'll use the calculationResult in our reactive statement
      if (calculationResult.claimingAnnualAmount) {
        return calculationResult.claimingAnnualAmount;
      } else if (calculationResult.claimingMonthlyAmount) {
        return calculationResult.claimingMonthlyAmount * 12;
      } else {
        // Fallback calculation based on monthly benefit estimate
        return (data.estimatedMonthlyBenefit || 0) * 12;
      }
    } catch (err) {
      console.error("Error calculating Social Security:", err);
      // Provide more specific error messages
      if (err instanceof Error) {
        // Check for the specific birthMonth type error from the Go backend
        if (err.message.includes("birthMonth") && err.message.includes("type")) {
          // Fix the type issue and retry
          data.birthMonth = parseInt(String(data.birthMonth), 10);
          error = "Birth month format fixed. Please try calculating again.";
        } else {
          error = `Error: ${err.message}`;
        }
      } else {
        error = 'Failed to calculate Social Security benefit. Please try again.';
      }
      return 0;
    } finally {
      loading = false;
    }
  }

  // Load saved calculation results when component mounts or scenario changes
  onMount(async () => {
    try {
      // Retrieve previously calculated result from the store
      const savedResult = await getCalculationResult(scenarioId, 'socialSecurity');
      console.log('Retrieved saved Social Security calculation for scenario', scenarioId, savedResult);
      
      if (savedResult) {
        // Update the calculation result with the stored data
        calculationResult = savedResult;
        
        // Also update the estimated monthly benefit if it was empty
        if (!data.estimatedMonthlyBenefit && calculationResult.claimingMonthlyAmount) {
          estimatedMonthlyBenefit = calculationResult.claimingMonthlyAmount;
          data.estimatedMonthlyBenefit = estimatedMonthlyBenefit;
          onUpdate({ ...data });
        }
      } else {
        // If no saved result and we have enough data, calculate
        checkAndCalculate();
      }
    } catch (err) {
      console.error('Error loading saved calculation:', err);
    }
  });
  
  // No timers to clean up
  
  // Handle scenario changes - load saved calculation when scenario changes
  $: if (scenarioId) {
    // Asynchronously get the saved calculation for the new scenario
    getCalculationResult(scenarioId, 'socialSecurity').then(savedResult => {
      if (savedResult) {
        console.log('Scenario changed - loading calculation for new scenario:', scenarioId, savedResult);
        calculationResult = savedResult;
      }
    });
  }
  
  // Simplified auto-calculation approach
  function checkAndCalculate() {
    // We can call this on any field change
    if (data.birthYear && data.startAge) {
      calculateSocialSecurity();
    }
  }

  $: {
    // Calculate annual and monthly benefit amounts
    let annual = 0;
    let monthly = 0;
    
    if (calculationResult && calculationResult.claimingMonthlyAmount > 0) {
      // Use result from calculation if available
      monthly = calculationResult.claimingMonthlyAmount;
      annual = calculationResult.claimingAnnualAmount || monthly * 12;
    } else if (data.estimatedMonthlyBenefit && data.estimatedMonthlyBenefit > 0) {
      // Fallback to data directly
      monthly = data.estimatedMonthlyBenefit;
      annual = monthly * 12;
    } else if (data.ssaEstimateAt62 && data.startAge === 62) {
      // Use SSA estimate if we have it and it matches claim age
      monthly = data.ssaEstimateAt62;
      annual = monthly * 12;
    } else if (data.ssaEstimateAtFRA && data.startAge === 67) {
      // Use FRA estimate
      monthly = data.ssaEstimateAtFRA;
      annual = monthly * 12;
    } else if (data.ssaEstimateAt70 && data.startAge === 70) {
      // Use age 70 estimate
      monthly = data.ssaEstimateAt70;
      annual = monthly * 12;
    }
    
    // Update estimate in our data if it's not already set
    if (!data.estimatedMonthlyBenefit && monthly > 0) {
      data.estimatedMonthlyBenefit = monthly;
    }
    
    annualBenefit = annual;
    monthlyBenefit = monthly;
  }
</script>

<div data-section="social-security">
  <SectionHeader sectionName="Social Security Calculator" {scenarioName} />
  <!-- Required fields legend -->
  <div class="mb-4 text-sm text-gray-600 dark:text-gray-400 flex items-center">
    <span class="text-red-500 mr-1">*</span> Required fields for accurate calculations
  </div>

  <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
    <div class="space-y-4">
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="birthYear">
          Birth Year *
        </label>
        <input
          id="birthYear"
          type="number"
          min="1930"
          max="2000"
          class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
          bind:value={birthYear}
          on:change={() => {
            // Ensure stored as number
            if (typeof birthYear === 'string') {
              birthYear = parseInt(birthYear, 10);
            }
            handleFieldChange();
          }}
        />
      </div>
      
      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="birthMonth">
          Birth Month *
        </label>
        <select 
          id="birthMonth"
          class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
          bind:value={birthMonth}
          on:change={() => {
            // Ensure birthMonth is stored as a number
            if (typeof birthMonth === 'string') {
              birthMonth = parseInt(birthMonth, 10);
            }
            handleFieldChange();
          }}
        >
          <option value={1}>January</option>
          <option value={2}>February</option>
          <option value={3}>March</option>
          <option value={4}>April</option>
          <option value={5}>May</option>
          <option value={6}>June</option>
          <option value={7}>July</option>
          <option value={8}>August</option>
          <option value={9}>September</option>
          <option value={10}>October</option>
          <option value={11}>November</option>
          <option value={12}>December</option>
        </select>
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" for="startAge">
          Benefit Start Age *
        </label>
        <select 
          id="startAge"
          class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
          bind:value={startAge}
          on:change={() => {
            // Ensure startAge is stored as a number
            if (typeof startAge === 'string') {
              startAge = parseInt(startAge, 10);
            }
            handleFieldChange();
          }}
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
          bind:checked={isEligible}
          on:change={handleFieldChange}
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
          on:change={() => {
              if (data.ssaEstimateAt62 && typeof data.ssaEstimateAt62 === 'string') {
                data.ssaEstimateAt62 = parseFloat(data.ssaEstimateAt62);
              }
              handleFieldChange();
              onUpdate?.(data);
            }}
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
          on:change={() => {
              if (data.ssaEstimateAtFRA && typeof data.ssaEstimateAtFRA === 'string') {
                data.ssaEstimateAtFRA = parseFloat(data.ssaEstimateAtFRA);
              }
              handleFieldChange();
            }}
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
          on:change={() => {
              if (data.ssaEstimateAt70 && typeof data.ssaEstimateAt70 === 'string') {
                data.ssaEstimateAt70 = parseFloat(data.ssaEstimateAt70);
              }
              handleFieldChange();
            }}
          />
        </div>
      </div>
      
      <!-- Calculate Button (now secondary since we have a global calculate button) -->
      <button
        class="mt-4 w-full inline-flex justify-center items-center px-4 py-2 border border-gray-300 dark:border-gray-600 text-sm font-medium rounded-md shadow-sm text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-800 hover:bg-gray-50 dark:hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 disabled:opacity-50 disabled:cursor-not-allowed"
        on:click={calculateSocialSecurity}
        disabled={loading || !data.birthYear || !data.startAge}
      >
        {#if loading}
          <span class="mr-2 animate-spin rounded-full h-4 w-4 border-b-2 border-primary-600 dark:border-primary-400"></span>
          Calculating...
        {:else}
          <span class="mr-1">📊</span> Update Estimates
        {/if}
      </button>
      <div class="mt-2 text-xs text-gray-500 dark:text-gray-400 text-center">
        Or use the "Calculate All" button in the header for full scenario calculations
      </div>
    </div>
  </div>

  <div class="mt-8 bg-gray-50 dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-lg p-4">
    <h3 class="text-lg font-medium text-gray-800 dark:text-gray-300 mb-3">Social Security Summary</h3>
    
    {#if !calculationResult || (!calculationResult.estimatedMonthlyAt62 && !error)}
      <div class="bg-blue-100 dark:bg-blue-900 border border-blue-200 dark:border-blue-700 text-blue-700 dark:text-blue-300 p-4 rounded-lg mb-4">
        To calculate your Social Security benefits, please fill in your birth information and claiming age (marked with * above), then click the "Update Estimates" button or use the "Calculate All" button in the header.
      </div>
    {/if}
    
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
          <div class="text-sm text-gray-500 dark:text-gray-400">Annual Benefit at Age {data.startAge || 62}</div>
          <div class="text-2xl font-semibold text-gray-800 dark:text-gray-200">
            ${(isNaN(annualBenefit) ? 0 : annualBenefit).toLocaleString(undefined, { maximumFractionDigits: 0 })}
          </div>
        </div>
        <div class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
          <div class="text-sm text-gray-500 dark:text-gray-400">Monthly Benefit at Age {data.startAge || 62}</div>
          <div class="text-2xl font-semibold text-gray-800 dark:text-gray-200">
            ${(isNaN(monthlyBenefit) ? 0 : monthlyBenefit).toLocaleString(undefined, { maximumFractionDigits: 0 })}
          </div>
        </div>
      </div>
      
      {#if calculationResult && calculationResult.fullRetirementAge}
        <div class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow mb-4">
          <h4 class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Benefit Estimates by Claiming Age</h4>
          <div class="grid grid-cols-3 gap-4">
            <div>
              <div class="text-xs text-gray-500 dark:text-gray-400">Age 62</div>
              <div class="text-sm font-medium text-gray-800 dark:text-gray-200">
                ${(calculationResult.estimatedMonthlyAt62 || 0).toLocaleString(undefined, { maximumFractionDigits: 0 })}/mo
              </div>
            </div>
            <div>
              <div class="text-xs text-gray-500 dark:text-gray-400">Full Age ({calculationResult.fullRetirementAge})</div>
              <div class="text-sm font-medium text-gray-800 dark:text-gray-200">
                ${(calculationResult.estimatedMonthlyAtFRA || 0).toLocaleString(undefined, { maximumFractionDigits: 0 })}/mo
              </div>
            </div>
            <div>
              <div class="text-xs text-gray-500 dark:text-gray-400">Age 70</div>
              <div class="text-sm font-medium text-gray-800 dark:text-gray-200">
                ${(calculationResult.estimatedMonthlyAt70 || 0).toLocaleString(undefined, { maximumFractionDigits: 0 })}/mo
              </div>
            </div>
          </div>
        </div>
      {/if}
      
      {#if calculationResult && calculationResult.notes}
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