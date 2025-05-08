<script lang="ts">
  import { main } from '../../wailsjs/go/models.js';
  import type { SocialSecurityData } from '../types/scenario.js';
  import { api } from '../api.js';
  import { onMount } from 'svelte';
  import SectionHeader from './SectionHeader.svelte';
  
  // Default data values
  const defaultData = {
    startAge: 67,
    estimatedMonthlyBenefit: 2000,
    isEligible: true,
    birthYear: 1970,
    birthMonth: 1
  };
  
  const { 
    data: propData = defaultData, 
    scenarioId, 
    scenarioName, 
    onUpdate = (data: any) => {} 
  } = $props<{
    data?: SocialSecurityData;
    scenarioId: number;
    scenarioName: string;
    onUpdate?: (data: any) => void;
  }>();
  
  // Using $state for local data, initialized with defaults or incoming data
  let localData = $state<SocialSecurityData>(
    propData ? { ...defaultData, ...propData } : { ...defaultData }
  );
  
  // Store a previous version of localData as a string for comparison
  let previousLocalDataString = $state("");
  
  // Initialize previousLocalDataString after localData is defined
  $effect(() => {
    if (previousLocalDataString === "") {
      previousLocalDataString = JSON.stringify(localData);
    }
  });
  
  // Track if we're updating from props to prevent loops
  let updatingFromProps = $state(false);

  // Manual function to update local data when props change
  function updateLocalDataFromProps() {
    if (!propData) return;
    
    try {
      const dataString = JSON.stringify(propData);
      // Only update if different
      if (dataString !== previousLocalDataString) {
        console.log('SocialSecuritySection: Props changed, updating local data');
        updatingFromProps = true;
        localData = { ...defaultData, ...propData };
        
        // Make sure certain fields are the correct type
        if (typeof localData.birthMonth === 'string') {
          localData.birthMonth = parseInt(String(localData.birthMonth), 10);
        }
        
        // Ensure these fields exist with defaults
        localData.ssaEstimateAt62 = localData.ssaEstimateAt62 || 0;
        localData.ssaEstimateAtFRA = localData.ssaEstimateAtFRA || 0;
        localData.ssaEstimateAt70 = localData.ssaEstimateAt70 || 0;
        
        previousLocalDataString = JSON.stringify(localData);
      }
    } finally {
      updatingFromProps = false;
    }
  }

  // Update local data when props change - but only on initial props or when they actually change
  let previousPropString = $state('');
  
  $effect(() => {
    if (propData) {
      const propString = JSON.stringify(propData);
      if (propString !== previousPropString) {
        previousPropString = propString;
        updateLocalDataFromProps();
        // Log the local data for debugging
        console.log('SocialSecuritySection localData:', JSON.stringify(localData));
      }
    }
  });

  // Manual function to update the props from localData
  function updateParentData() {
    if (updatingFromProps) return; // Skip if we're currently updating from props
    
    const currentString = JSON.stringify(localData);
    
    // Only update if the data has actually changed
    if (currentString !== previousLocalDataString) {
      console.log('SocialSecuritySection: localData changed, updating parent');
      previousLocalDataString = currentString;
      
      // Create a deep copy to avoid reference issues
      const dataForParent = JSON.parse(currentString);
      
      // Notify parent via callback
      onUpdate(dataForParent);
    }
  }

  // Handle all field changes in a single function
  function handleFieldChange() {
    console.log('SocialSecuritySection updating data - field changed');
    
    // Update parent with new data
    updateParentData();
    
    // Try to calculate on change
    checkAndCalculate();
  }
  
  // Create a local calculation cache instead of using the global store
  let calculationCache = $state(new Map());
  
  // Helper to get calculation result
  function getCalculationResult(scenarioId, type) {
    const key = `${scenarioId}-${type}`;
    return calculationCache.get(key);
  }
  
  let loading = $state(false);
  let error = $state('');
  let calculationResult = $state({
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
  });
  
  // Default values at age 62 - need to be $state because they're updated in an effect
  let monthlyBenefit = $state(0);
  let annualBenefit = $state(0);

  const ageOptions = [
    { value: 62, label: 'Age 62 (reduced)' },
    { value: 67, label: 'Age 67 (full retirement age)' },
    { value: 70, label: 'Age 70 (maximum benefit)' }
  ];
  
  const currentYear = new Date().getFullYear();
  
  // Make the function accessible from outside
  // We need to find and fix all references to data in the template

  // Debugging function is no longer needed as we've resolved the issue
// Remove to prevent unnecessary effects

export async function calculateSocialSecurity() {
    try {
      loading = true;
      error = '';
      
      // Validate required inputs
      if (!localData.birthYear) {
        throw new Error("Birth year is required");
      }
      
      if (localData.birthYear < 1900 || localData.birthYear > currentYear - 20) {
        throw new Error(`Birth year must be between 1900 and ${currentYear - 20}`);
      }
      
      if (!localData.startAge) {
        throw new Error("Benefit start age is required");
      }
      
      // Create input object for backend
      const ssInput = new main.SocialSecurityInput();
      // Ensure all numeric values are properly converted to numbers
      ssInput.startAge = parseInt(String(localData.startAge || 62), 10);
      ssInput.estimatedMonthlyBenefit = parseFloat(String(localData.estimatedMonthlyBenefit || 0));
      ssInput.isEligible = localData.isEligible !== undefined ? localData.isEligible : true;
      ssInput.birthYear = parseInt(String(localData.birthYear || 1960), 10);
      ssInput.birthMonth = parseInt(String(localData.birthMonth || 1), 10);
      ssInput.currentAge = currentYear - parseInt(String(localData.birthYear || 1960), 10);
      
      // For estimation - default values if not available
      ssInput.estimatedAnnualSalary = 60000; // Default annual salary for estimation
      ssInput.yearsWorked = 35; // Default years worked for estimation
      
      // If SSA estimates are provided, use them - ensure they're numbers
      if (localData.ssaEstimateAt62 && localData.ssaEstimateAt62 > 0) {
        ssInput.userProvidedEstimate62 = parseFloat(String(localData.ssaEstimateAt62));
      }
      if (localData.ssaEstimateAtFRA && localData.ssaEstimateAtFRA > 0) {
        ssInput.userProvidedEstimateFRA = parseFloat(String(localData.ssaEstimateAtFRA));
      }
      if (localData.ssaEstimateAt70 && localData.ssaEstimateAt70 > 0) {
        ssInput.userProvidedEstimate70 = parseFloat(String(localData.ssaEstimateAt70));
      }
      
      // Call backend API via the store
      calculationResult = await api.calculateSocialSecurity(ssInput);
      
      // Store the calculation locally for reuse
      calculationCache.set(`${scenarioId}-socialSecurity`, { 
        ...calculationResult,
        timestamp: new Date().toISOString()
      });
      
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
        return (localData.estimatedMonthlyBenefit || 0) * 12;
      }
    } catch (err) {
      console.error("Error calculating Social Security:", err);
      // Provide more specific error messages
      if (err instanceof Error) {
        // Check for the specific birthMonth type error from the Go backend
        if (err.message.includes("birthMonth") && err.message.includes("type")) {
          // Fix the type issue and retry
          localData.birthMonth = parseInt(String(localData.birthMonth), 10);
          updateParentData();
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

  // Load calculation results when component mounts
  onMount(() => {
    try {
      // Retrieve previously calculated result from our local cache
      const savedResult = getCalculationResult(scenarioId, 'socialSecurity');
      console.log('Retrieved saved Social Security calculation for scenario', scenarioId, savedResult);
      
      if (savedResult) {
        // Update the calculation result with the stored data
        calculationResult = savedResult;
        
        // Also update the estimated monthly benefit if it was empty
        if (!localData.estimatedMonthlyBenefit && calculationResult.claimingMonthlyAmount) {
          localData.estimatedMonthlyBenefit = calculationResult.claimingMonthlyAmount;
          updateParentData();
        }
      } else {
        // If no saved result and we have enough data, calculate
        checkAndCalculate();
      }
    } catch (err) {
      console.error('Error loading saved calculation:', err);
    }
  });
  
  // For template references that might use data directly
  // Define a static compatibility object - not reactive to avoid update loops
  const data = {
    startAge: localData.startAge,
    birthYear: localData.birthYear
  };
  
  // No timers to clean up
  
  // Handle scenario changes - load saved calculation when scenario changes
  $effect(() => {
    if (scenarioId) {
      // Get saved calculation from our local cache
      const savedResult = getCalculationResult(scenarioId, 'socialSecurity');
      if (savedResult) {
        console.log('Scenario changed - loading calculation for new scenario:', scenarioId, savedResult);
        calculationResult = savedResult;
      }
    }
  });
  
  // Simplified auto-calculation approach
  function checkAndCalculate() {
    // We can call this on any field change
    if (localData.birthYear && localData.startAge) {
      calculateSocialSecurity();
    }
  }

  // Calculate annual and monthly benefit amounts based on data/calculation state
  $effect(() => {
    let annual = 0;
    let monthly = 0;
    
    if (calculationResult && calculationResult.claimingMonthlyAmount > 0) {
      // Use result from calculation if available
      monthly = calculationResult.claimingMonthlyAmount;
      annual = calculationResult.claimingAnnualAmount || monthly * 12;
    } else if (localData.estimatedMonthlyBenefit && localData.estimatedMonthlyBenefit > 0) {
      // Fallback to data directly
      monthly = localData.estimatedMonthlyBenefit;
      annual = monthly * 12;
    } else if (localData.ssaEstimateAt62 && localData.startAge === 62) {
      // Use SSA estimate if we have it and it matches claim age
      monthly = localData.ssaEstimateAt62;
      annual = monthly * 12;
    } else if (localData.ssaEstimateAtFRA && localData.startAge === 67) {
      // Use FRA estimate
      monthly = localData.ssaEstimateAtFRA;
      annual = monthly * 12;
    } else if (localData.ssaEstimateAt70 && localData.startAge === 70) {
      // Use age 70 estimate
      monthly = localData.ssaEstimateAt70;
      annual = monthly * 12;
    }
    
    // Update local state variables
    annualBenefit = annual;
    monthlyBenefit = monthly;
    
    // Update estimate if it's not already set
    if (monthly > 0 && !localData.estimatedMonthlyBenefit) {
      localData.estimatedMonthlyBenefit = monthly;
      updateParentData();
    }
  });
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
          bind:value={localData.birthYear}
          onchange={() => {
            // Ensure stored as number
            if (typeof localData.birthYear === 'string') {
              localData.birthYear = parseInt(localData.birthYear, 10);
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
          bind:value={localData.birthMonth}
          onchange={() => {
            // Ensure birthMonth is stored as a number
            if (typeof localData.birthMonth === 'string') {
              localData.birthMonth = parseInt(localData.birthMonth, 10);
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
          bind:value={localData.startAge}
          onchange={() => {
            // Ensure startAge is stored as a number
            if (typeof localData.startAge === 'string') {
              localData.startAge = parseInt(localData.startAge, 10);
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
          bind:checked={localData.isEligible}
          onchange={handleFieldChange}
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
            bind:value={localData.ssaEstimateAt62}
          onchange={() => {
              if (typeof localData.ssaEstimateAt62 === 'string') {
                localData.ssaEstimateAt62 = parseFloat(localData.ssaEstimateAt62);
              }
              handleFieldChange();
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
            bind:value={localData.ssaEstimateAtFRA}
          onchange={() => {
              if (typeof localData.ssaEstimateAtFRA === 'string') {
                localData.ssaEstimateAtFRA = parseFloat(localData.ssaEstimateAtFRA);
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
            bind:value={localData.ssaEstimateAt70}
          onchange={() => {
              if (typeof localData.ssaEstimateAt70 === 'string') {
                localData.ssaEstimateAt70 = parseFloat(localData.ssaEstimateAt70);
              }
              handleFieldChange();
            }}
          />
        </div>
      </div>
      
      <!-- Calculate Button (now secondary since we have a global calculate button) -->
      <button
        class="mt-4 w-full inline-flex justify-center items-center px-4 py-2 border border-gray-300 dark:border-gray-600 text-sm font-medium rounded-md shadow-sm text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-800 hover:bg-gray-50 dark:hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 disabled:opacity-50 disabled:cursor-not-allowed"
        onclick={calculateSocialSecurity}
        disabled={loading || !localData.birthYear || !localData.startAge}
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
          <div class="text-sm text-gray-500 dark:text-gray-400">Annual Benefit at Age {localData.startAge || 62}</div>
          <div class="text-2xl font-semibold text-gray-800 dark:text-gray-200">
            ${(isNaN(annualBenefit) ? 0 : annualBenefit).toLocaleString(undefined, { maximumFractionDigits: 0 })}
          </div>
        </div>
        <div class="bg-white dark:bg-gray-800 p-4 rounded-lg shadow">
          <div class="text-sm text-gray-500 dark:text-gray-400">Monthly Benefit at Age {localData.startAge || 62}</div>
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