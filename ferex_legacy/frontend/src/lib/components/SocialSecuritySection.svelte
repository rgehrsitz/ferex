<script lang="ts">
  import type { ScenarioInput } from '../../types';

  let { inputs = $bindable() } = $props<{ inputs: ScenarioInput }>();

  // Standard SS claiming ages
  const standardAges = [62, 63, 64, 65, 66, 67, 68, 69, 70];

  // Calculate benefit for any age based on the 3 key data points
  function calculateBenefitForAge(targetAge: number): number | null {
    // Get the three key values
    const age62Benefit = inputs.userProvidedSSBenefitAmount1;
    const age67Benefit = inputs.socialSecurityEstimate; // FRA benefit
    const age70Benefit = inputs.userProvidedSSBenefitAmount2;
    
    // Create array of known points (convert to number to handle string inputs)
    const knownPoints: {age: number, benefit: number}[] = [];
    if (age62Benefit && Number(age62Benefit) > 0) knownPoints.push({age: 62, benefit: Number(age62Benefit)});
    if (age67Benefit && Number(age67Benefit) > 0) knownPoints.push({age: 67, benefit: Number(age67Benefit)});
    if (age70Benefit && Number(age70Benefit) > 0) knownPoints.push({age: 70, benefit: Number(age70Benefit)});
    
    // If we have only one data point, use simple scaling for interpolation
    if (knownPoints.length === 1) {
      const singlePoint = knownPoints[0];
      // Use standard Social Security scaling factors
      const scaleFactors = {
        62: 0.75,  // ~75% of FRA benefit
        63: 0.80,
        64: 0.87,
        65: 0.93,
        66: 0.97,
        67: 1.00,  // Full Retirement Age
        68: 1.08,
        69: 1.16,
        70: 1.24   // ~124% of FRA benefit
      };
      
      // Calculate the implied FRA benefit from the single data point
      const fraMultiplier = scaleFactors[singlePoint.age as keyof typeof scaleFactors] || 1.0;
      const impliedFRABenefit = singlePoint.benefit / fraMultiplier;
      
      // Calculate benefit for target age
      const targetMultiplier = scaleFactors[targetAge as keyof typeof scaleFactors];
      if (targetMultiplier) {
        const calculatedBenefit = Math.round(impliedFRABenefit * targetMultiplier);
        return calculatedBenefit;
      }
    }
    
    // Need at least 2 points to interpolate
    if (knownPoints.length < 2) {
      return null;
    }
    
    // Sort by age
    knownPoints.sort((a, b) => a.age - b.age);
    
    // Check for exact match
    const exactMatch = knownPoints.find(point => point.age === targetAge);
    if (exactMatch) return exactMatch.benefit;
    
    // Linear interpolation/extrapolation
    if (knownPoints.length >= 2) {
      // Find the two points to use for interpolation
      let lowerPoint = knownPoints[0];
      let upperPoint = knownPoints[1];
      
      // If we have 3 points, choose the best pair
      if (knownPoints.length === 3) {
        if (targetAge <= knownPoints[1].age) {
          // Use first two points
          lowerPoint = knownPoints[0];
          upperPoint = knownPoints[1];
        } else {
          // Use last two points
          lowerPoint = knownPoints[1];
          upperPoint = knownPoints[2];
        }
      }
      
      // Calculate slope and interpolate
      const slope = (upperPoint.benefit - lowerPoint.benefit) / (upperPoint.age - lowerPoint.age);
      const calculatedBenefit = lowerPoint.benefit + slope * (targetAge - lowerPoint.age);
      
      return Math.round(Math.max(0, calculatedBenefit)); // Ensure non-negative
    }
    
    return null;
  }

  // Auto-set the claiming age based on retirement age if not already set
  function getRecommendedClaimingAge(): number {
    if (inputs.userSSClaimingAge) return inputs.userSSClaimingAge;
    
    // Calculate retirement age from birth date and planned retirement date
    if (inputs.dateOfBirth && inputs.plannedRetirementDate) {
      const birthDate = new Date(inputs.dateOfBirth);
      const retirementDate = new Date(inputs.plannedRetirementDate);
      const retirementAge = Math.floor((retirementDate.getTime() - birthDate.getTime()) / (365.25 * 24 * 60 * 60 * 1000));
      
      // Default to retirement age, but cap between 62-70
      return Math.max(62, Math.min(70, retirementAge));
    }
    
    return 67; // Default FRA
  }

  let recommendedClaimingAge = $derived(getRecommendedClaimingAge());
  let projectedBenefit = $derived(calculateBenefitForAge(inputs.userSSClaimingAge || recommendedClaimingAge));

  // Set claiming age to recommended if not already set
  $effect(() => {
    if (!inputs.userSSClaimingAge) {
      inputs.userSSClaimingAge = recommendedClaimingAge;
    }
  });
</script>

<section class="space-y-6 p-4 border border-gray-200 rounded-md shadow-sm">
  <h3 class="text-lg font-semibold text-gray-700 border-b pb-2 mb-3">Social Security Benefits</h3>
  
  <!-- SS Benefit Table Input -->
  <div>
    <h4 class="text-md font-medium text-gray-700 mb-3">Social Security Statement Benefits</h4>
    <p class="text-sm text-gray-600 mb-4">Enter just the key benefit amounts from your Social Security Statement. The app will calculate all other ages automatically.</p>
    
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <!-- Key ages that users typically have on their SS statement -->
      <div class="flex items-center space-x-2">
        <label for="ss-age-62" class="block text-sm font-medium text-gray-700 w-16">Age 62:</label>
        <div class="flex-1 relative">
          <span class="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-500">$</span>
          <input 
            id="ss-age-62"
            type="number" 
            bind:value={inputs.userProvidedSSBenefitAmount1}
            class="pl-8 pr-3 py-2 block w-full border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
            placeholder="2,795"
            min="0"
            step="1"
          />
          <span class="absolute right-3 top-1/2 transform -translate-y-1/2 text-xs text-gray-500">/month</span>
        </div>
      </div>

      <div class="flex items-center space-x-2">
        <label for="ss-age-67" class="block text-sm font-medium text-gray-700 w-16">Age 67:</label>
        <div class="flex-1 relative">
          <span class="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-500">$</span>
          <input 
            id="ss-age-67"
            type="number" 
            bind:value={inputs.socialSecurityEstimate}
            class="pl-8 pr-3 py-2 block w-full border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
            placeholder="4,012"
            min="0"
            step="1"
          />
          <span class="absolute right-3 top-1/2 transform -translate-y-1/2 text-xs text-gray-500">/month</span>
        </div>
      </div>

      <div class="flex items-center space-x-2">
        <label for="ss-age-70" class="block text-sm font-medium text-gray-700 w-16">Age 70:</label>
        <div class="flex-1 relative">
          <span class="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-500">$</span>
          <input 
            id="ss-age-70"
            type="number" 
            bind:value={inputs.userProvidedSSBenefitAmount2}
            class="pl-8 pr-3 py-2 block w-full border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
            placeholder="5,000"
            min="0"
            step="1"
          />
          <span class="absolute right-3 top-1/2 transform -translate-y-1/2 text-xs text-gray-500">/month</span>
        </div>
      </div>
    </div>
    
    <p class="text-xs text-gray-500 mt-2">
      <strong>Tip:</strong> These are the three most common ages shown on SS statements. Enter any 2-3 values you have.
    </p>
  </div>

  <!-- Claiming Strategy -->
  <div class="bg-blue-50 p-4 rounded-md">
    <h4 class="text-md font-medium text-gray-700 mb-3">Your Claiming Strategy</h4>
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div>
        <label for="userSSClaimingAge" class="block text-sm font-medium text-gray-700 mb-1">Social Security Claiming Age:</label>
        <select 
          id="userSSClaimingAge" 
          bind:value={inputs.userSSClaimingAge}
          class="mt-1 block w-full p-2 bg-white text-gray-900 border border-gray-300 dark:bg-slate-700 dark:text-slate-200 dark:border-slate-600 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
        >
          {#each standardAges as age}
            <option value={age}>Age {age}</option>
          {/each}
        </select>
        <p class="mt-1 text-xs text-gray-500">
          {#if recommendedClaimingAge !== inputs.userSSClaimingAge}
            Recommended: Age {recommendedClaimingAge} (based on retirement date)
          {:else}
            Matches your retirement age
          {/if}
        </p>
      </div>

      <div>
        <div class="block text-sm font-medium text-gray-700 mb-1">Projected Monthly Benefit:</div>
        <div class="mt-1 block w-full p-2 border border-gray-300 rounded-md bg-gray-50 text-lg font-semibold {projectedBenefit ? 'text-green-700' : 'text-gray-500'}">
          {#if projectedBenefit}
            ${projectedBenefit.toLocaleString()}
          {:else}
            {#if !inputs.userProvidedSSBenefitAmount1 && !inputs.socialSecurityEstimate && !inputs.userProvidedSSBenefitAmount2}
              Enter benefits above
            {:else}
              Calculating...
            {/if}
          {/if}
        </div>
        <p class="mt-1 text-xs text-gray-500">
          {#if projectedBenefit}
            Based on your claiming age selection
          {:else if inputs.userProvidedSSBenefitAmount1 || inputs.socialSecurityEstimate || inputs.userProvidedSSBenefitAmount2}
            Need at least one benefit amount to calculate
          {:else}
            Enter benefit amounts from your Social Security Statement
          {/if}
        </p>
      </div>
    </div>
  </div>

  <!-- Advanced Options -->
  <div class="space-y-4">
    <div>
      <label for="userAssumedSSCOLA" class="block text-sm font-medium text-gray-700 mb-1">Assumed Annual COLA (%):</label>
      <input 
        type="number" 
        id="userAssumedSSCOLA" 
        bind:value={inputs.userAssumedSSCOLA} 
        class="mt-1 block w-32 p-2 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm bg-white" 
        placeholder="2.0" 
        step="0.1" 
        min="0"
        max="10"
      />
      <p class="mt-1 text-xs text-gray-500">Historical average is around 2-3%</p>
    </div>
  </div>
</section>