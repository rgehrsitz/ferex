<script lang="ts">
  import type { ScenarioInput } from '../../types';
  import { createEventDispatcher } from 'svelte';

  let { inputs = $bindable() } = $props<{ inputs?: ScenarioInput }>();
  const dispatch = createEventDispatcher();

  function handleChange() {
    // console.log('[DEBUG] COLA Settings change:', {
    //   pensionCola: inputs.pensionAnnuityColaRate,
    //   ssCola: inputs.userAssumedSSCOLA,
    //   generalInflation: inputs.generalInflationRate
    // });
    dispatch('change');
  }

  // Initialize COLA values if they don't exist or are undefined
  $effect(() => {
    if (inputs) { // Add this check
      inputs.pensionAnnuityColaRate = inputs.pensionAnnuityColaRate === undefined ? null : inputs.pensionAnnuityColaRate;
      inputs.userAssumedSSCOLA = inputs.userAssumedSSCOLA === undefined ? null : inputs.userAssumedSSCOLA;
      inputs.userAssumedGeneralInflationRate = inputs.userAssumedGeneralInflationRate === undefined ? null : inputs.userAssumedGeneralInflationRate;
    }
  });

  const defaultPensionCOLA = 2.0; // Example default, user can override
  const defaultSSCOLA = 2.0;      // Example default
  const defaultGeneralInflation = 2.5; // Example default

</script>

{#if inputs}
<div class="flex flex-col gap-6">
  <h1 class="text-2xl font-bold text-gray-800 dark:text-gray-100">Application Settings</h1>

  <!-- Calculation Settings Section -->
  <section class="space-y-4 p-4 border border-gray-200 dark:border-gray-700 rounded-md shadow-sm bg-white dark:bg-gray-800">
    <h3 class="text-lg font-semibold text-gray-700 dark:text-gray-200 border-b border-gray-300 dark:border-gray-600 pb-2 mb-4">
      Calculation Settings
    </h3>
    <div class="grid grid-cols-1 md:grid-cols-2 gap-x-6 gap-y-4">
      <div>
        <label for="calculationYears" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Years to Calculate (Debug Mode)</label>
        <input 
          type="number" 
          id="calculationYears" 
          bind:value={inputs.calculationYears} 
          oninput={handleChange} 
          min="1" 
          max="50" 
          step="1" 
          class="mt-1 block w-full p-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100" 
          placeholder="e.g., 30" 
        />
        <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">Restrict all calculations to N years for debugging. Leave blank for full projection (30 years).</p>
      </div>
    </div>
  </section>

  <section class="space-y-4 p-4 border border-gray-200 dark:border-gray-700 rounded-md shadow-sm bg-white dark:bg-gray-800">
    <h3 class="text-lg font-semibold text-gray-700 dark:text-gray-200 border-b border-gray-300 dark:border-gray-600 pb-2 mb-4">
      Cost of Living Adjustments (COLA) & Inflation Rates (%)
    </h3>
    <div class="grid grid-cols-1 md:grid-cols-3 gap-x-6 gap-y-4">
      <div>
        <label for="pensionColaRate" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Pension COLA Rate</label>
        <input type="number" id="pensionColaRate" bind:value={inputs.pensionAnnuityColaRate} oninput={handleChange} step="0.1" class="mt-1 block w-full p-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100" placeholder="e.g., {defaultPensionCOLA}" />
        <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">Annual COLA for FERS/CSRS pension.</p>
      </div>

      <div>
        <label for="ssColaRate" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Social Security COLA Rate</label>
        <input type="number" id="ssColaRate" bind:value={inputs.userAssumedSSCOLA} oninput={handleChange} step="0.1" class="mt-1 block w-full p-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100" placeholder="e.g., {defaultSSCOLA}" />
        <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">Annual COLA for Social Security benefits.</p>
      </div>

      <div>
        <label for="generalInflationRate" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">General Inflation Rate</label>
        <input type="number" id="generalInflationRate" step="0.1" class="mt-1 block w-full p-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100" placeholder="e.g., {defaultGeneralInflation}" bind:value={inputs.userAssumedGeneralInflationRate} oninput={handleChange} />
        <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">For other expenses or income streams.</p>
      </div>
    </div>
  </section>

  <!-- Other settings sections can be added here -->
  <!-- For example, Monte Carlo settings, UI preferences, etc. -->

</div>
{:else}
  <div class="p-4 text-gray-500 dark:text-gray-400">Loading settings or inputs not available...</div>
{/if}
