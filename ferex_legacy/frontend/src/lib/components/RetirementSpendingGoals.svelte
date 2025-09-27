<script lang="ts">
  import type { ScenarioInput } from '../../types';
  import { createEventDispatcher } from 'svelte';

  let { inputs = $bindable() } = $props<{ inputs: ScenarioInput }>();
  const dispatch = createEventDispatcher();

  function handleChange() {
    dispatch('change');
  }

  // Initialize Retirement Spending Goals values if they don't exist
  $effect(() => {
    inputs.desiredAnnualSpendingInRetirement = inputs.desiredAnnualSpendingInRetirement === undefined ? null : inputs.desiredAnnualSpendingInRetirement;
    inputs.isSpendingGoalInflationAdjusted = inputs.isSpendingGoalInflationAdjusted === undefined ? true : inputs.isSpendingGoalInflationAdjusted;
  });
</script>

<section class="space-y-4 p-4 border border-gray-200 rounded-md shadow-sm">
  <h3 class="text-lg font-semibold text-gray-700 border-b pb-2 mb-3">
    Retirement Spending Goals
  </h3>
  <div class="grid grid-cols-1 md:grid-cols-2 gap-x-6 gap-y-4 items-start">
    <div>
      <label for="desiredAnnualSpendingInRetirement" class="block text-sm font-medium text-gray-700 mb-1">Desired Annual Spending (Today's Dollars)</label>
      <input type="number" id="desiredAnnualSpendingInRetirement" bind:value={inputs.desiredAnnualSpendingInRetirement} onblur={handleChange} onkeydown={(e) => e.key === 'Enter' && handleChange()} min="0" step="any" class="mt-1 block w-full p-2 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm" placeholder="e.g., 60000" />
      <p class="mt-1 text-xs text-gray-500">Enter your target annual pre-tax spending in retirement, expressed in today's dollars.</p>
    </div>
    <div class="pt-6 md:pt-7"> <!-- Adjusted padding for alignment with label --> 
      <label class="flex items-center text-sm font-medium text-gray-700">
        <input type="checkbox" bind:checked={inputs.isSpendingGoalInflationAdjusted} onblur={handleChange} onkeydown={(e) => e.key === 'Enter' && handleChange()} class="mr-2 h-4 w-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500" />
        <span>Adjust for inflation</span>
      </label>
      <p class="mt-1 ml-6 text-xs text-gray-500">If checked, this spending goal will increase annually with inflation.</p>
    </div>
  </div>
</section>
