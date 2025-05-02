<script lang="ts">
  import { CalculateSocialSecurity } from '../../wailsjs/go/main/App';
  import { main } from '../../wailsjs/go/models';
  import type { SocialSecurityData } from '../types/scenario';
  
  export let data: SocialSecurityData;
  
  // Initialize with defaults if not set and ensure proper types
  if (!data.birthYear) data.birthYear = 1970;
  // Make sure birthMonth is stored as a number
  if (!data.birthMonth) {
    data.birthMonth = 1;
  } else {
    data.birthMonth = parseInt(data.birthMonth, 10);
  }
  if (!data.startAge) data.startAge = 67; // Default to full retirement age
  if (data.isEligible === undefined) data.isEligible = true;
  
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
  
  async function calculateSocialSecurity() {
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
      ssInput.startAge = parseInt(data.startAge || 62, 10);
      ssInput.estimatedMonthlyBenefit = parseFloat(data.estimatedMonthlyBenefit || 0);
      ssInput.isEligible = data.isEligible !== undefined ? data.isEligible : true;
      ssInput.birthYear = parseInt(data.birthYear || 1960, 10);
      ssInput.birthMonth = parseInt(data.birthMonth || 1, 10);
      ssInput.currentAge = currentYear - parseInt(data.birthYear || 1960, 10);
      
      // For estimation - default values if not available
      ssInput.estimatedAnnualSalary = 60000; // Default annual salary for estimation
      ssInput.yearsWorked = 35; // Default years worked for estimation
      
      // If SSA estimates are provided, use them - ensure they're numbers
      if (data.ssaEstimateAt62 && data.ssaEstimateAt62 > 0) {
        ssInput.userProvidedEstimate62 = parseFloat(data.ssaEstimateAt62);
      }
      if (data.ssaEstimateAtFRA && data.ssaEstimateAtFRA > 0) {
        ssInput.userProvidedEstimateFRA = parseFloat(data.ssaEstimateAtFRA);
      }
      if (data.ssaEstimateAt70 && data.ssaEstimateAt70 > 0) {
        ssInput.userProvidedEstimate70 = parseFloat(data.ssaEstimateAt70);
      }
      
      console.log("Social Security calculation input:", JSON.stringify(ssInput, null, 2));
      
      // Call backend API
      calculationResult = await CalculateSocialSecurity(ssInput);
      
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
          data.birthMonth = parseInt(data.birthMonth, 10);
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

  // Track previous values to prevent infinite loops
  let prevValues = {
    startAge: data.startAge,
    estimatedMonthlyBenefit: data.estimatedMonthlyBenefit,
    birthYear: data.birthYear,
    birthMonth: data.birthMonth,
    ssaEstimateAt62: data.ssaEstimateAt62,
    ssaEstimateAtFRA: data.ssaEstimateAtFRA,
    ssaEstimateAt70: data.ssaEstimateAt70
  };
  
  // We've removed the automatic calculation trigger
  // Now we'll only calculate when the user clicks the calculate button
  // This code still updates the prevValues for tracking purposes
  $: {
    // Check if any relevant inputs have changed
    const hasChanged = 
      data.startAge !== prevValues.startAge ||
      data.estimatedMonthlyBenefit !== prevValues.estimatedMonthlyBenefit ||
      data.birthYear !== prevValues.birthYear ||
      data.birthMonth !== prevValues.birthMonth ||
      data.ssaEstimateAt62 !== prevValues.ssaEstimateAt62 ||
      data.ssaEstimateAtFRA !== prevValues.ssaEstimateAtFRA ||
      data.ssaEstimateAt70 !== prevValues.ssaEstimateAt70;
    
    if (hasChanged) {
      // Update previous values
      prevValues = {
        startAge: data.startAge,
        estimatedMonthlyBenefit: data.estimatedMonthlyBenefit,
        birthYear: data.birthYear,
        birthMonth: data.birthMonth,
        ssaEstimateAt62: data.ssaEstimateAt62,
        ssaEstimateAtFRA: data.ssaEstimateAtFRA,
        ssaEstimateAt70: data.ssaEstimateAt70
      };
      
      // No automatic calculation - we'll let the user click the button
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
          on:change={() => {
            // Ensure stored as number
            data.birthYear = parseInt(data.birthYear, 10);
          }}
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
          on:change={() => {
            // Ensure birthMonth is stored as a number
            data.birthMonth = parseInt(data.birthMonth, 10);
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
          Benefit Start Age
        </label>
        <select 
          id="startAge"
          class="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
          bind:value={data.startAge}
          on:change={() => {
            // Ensure startAge is stored as a number
            data.startAge = parseInt(data.startAge, 10);
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
            on:change={() => {
              if (data.ssaEstimateAt62) {
                data.ssaEstimateAt62 = parseFloat(data.ssaEstimateAt62);
              }
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
              if (data.ssaEstimateAtFRA) {
                data.ssaEstimateAtFRA = parseFloat(data.ssaEstimateAtFRA);
              }
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
              if (data.ssaEstimateAt70) {
                data.ssaEstimateAt70 = parseFloat(data.ssaEstimateAt70);
              }
            }}
          />
        </div>
      </div>
      
      <!-- Calculate Button -->
      <button
        class="mt-4 w-full inline-flex justify-center items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 disabled:opacity-50 disabled:cursor-not-allowed"
        on:click={calculateSocialSecurity}
        disabled={loading || !data.birthYear || !data.startAge}
      >
        {#if loading}
          <span class="mr-2 animate-spin rounded-full h-4 w-4 border-b-2 border-white"></span>
          Calculating...
        {:else}
          Calculate Social Security Benefits
        {/if}
      </button>
    </div>
  </div>

  <div class="mt-8 bg-gray-50 dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-lg p-4">
    <h3 class="text-lg font-medium text-gray-800 dark:text-gray-300 mb-3">Social Security Summary</h3>
    
    {#if !calculationResult || (!calculationResult.estimatedMonthlyAt62 && !error)}
      <div class="bg-blue-100 dark:bg-blue-900 border border-blue-200 dark:border-blue-700 text-blue-700 dark:text-blue-300 p-4 rounded-lg mb-4">
        To calculate your Social Security benefits, please fill in your birth information and claiming age, then click the "Calculate Social Security Benefits" button.
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